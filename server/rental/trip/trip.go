package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	"coolcar/shared/mongo/objid"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service implements a trip service
type Service struct {
	ProfileManager ProfileManager
	CarManager     CarManager
	POIManager     POIManager
	DistanceCalc   DistanceCalc
	Mongo          *dao.Mongo
	Logger         *zap.Logger
}

// ProfileManager defines the ACL(Anti Corruption Layer) for profile verification logic.
type ProfileManager interface {
	Verify(context.Context, id.AccountID) (id.IdentityID, error)
}

// CarManager defines the ACL(Anti Corruption Layer) for car management.
type CarManager interface {
	Verify(c context.Context, cid id.CarID, loc *rentalpb.Location) error
	Unlock(c context.Context, cid id.CarID, aid id.AccountID, tid id.TripID, avatarURL string) error
	Lock(c context.Context, cid id.CarID) error
}

// POIManager resolve POI(point of interest) information.
type POIManager interface {
	Resolve(context.Context, *rentalpb.Location) (string, error)
}

// DistanceCalc calculates the distance between two locations.
type DistanceCalc interface {
	DistanceKm(context.Context, *rentalpb.Location, *rentalpb.Location) (float64, error)
}

// CreateTrip creates a trip
func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	// 这里的c已经是经过auth.Interceptor处理过的
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}

	if req.CarId == "" || req.Start == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	// 验证驾驶者身份
	iID, err := s.ProfileManager.Verify(c, aid)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	// 检查车辆状态
	carID := id.CarID(req.CarId)
	err = s.CarManager.Verify(c, carID, req.Start)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	// 创建行程：写入数据库，开始计费
	ls := s.calcCurrentStatus(c, &rentalpb.LocationStatus{
		Location:     req.Start,
		TimestampSec: nowFunc(),
	}, req.Start)
	tr, err := s.Mongo.CreateTrip(c, &rentalpb.Trip{
		AccountId:  aid.String(),
		CarId:      carID.String(),
		IdentityId: iID.String(),
		Status:     rentalpb.TripStatus_IN_PROGRESS,
		Start:      ls,
		Current:    ls,
	})
	if err != nil {
		// 将错误打印日志返回给我们自己看，对外只返回AlreadyExists，不返回具体的错误信息
		s.Logger.Warn("cannot create trip", zap.Error(err))
		return nil, status.Error(codes.AlreadyExists, "")
	}

	// 车辆开锁，使用后台任务，返回行程创建成功的同时进行开锁操作
	go func() {
		err = s.CarManager.Unlock(context.Background(), carID, aid, objid.ToTripID(tr.ID), req.AvatarUrl)
		if err != nil {
			s.Logger.Error("cannot unlock car", zap.Error(err))
		}
	}()

	return &rentalpb.TripEntity{
		Id:   tr.ID.Hex(),
		Trip: tr.Trip,
	}, nil
}

func (s *Service) GetTrip(c context.Context, req *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	tr, err := s.Mongo.GetTrip(c, id.TripID(req.Id), aid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}

	return tr.Trip, nil
}

func (s *Service) GetTrips(c context.Context, req *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	trips, err := s.Mongo.GetTrips(c, aid, req.Status)
	if err != nil {
		s.Logger.Error("cannot get trips", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	res := &rentalpb.GetTripsResponse{}
	for _, tr := range trips {
		res.Trips = append(res.Trips, &rentalpb.TripEntity{
			Id:   tr.ID.Hex(),
			Trip: tr.Trip,
		})
	}
	return res, nil
}

func (s *Service) UpdateTrip(c context.Context, req *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {

	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}

	// 获取trip
	tid := id.TripID(req.Id)
	tr, err := s.Mongo.GetTrip(c, tid, aid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}

	if tr.Trip.Status == rentalpb.TripStatus_FINISHED {
		return nil, status.Error(codes.FailedPrecondition, "cannot update a finished trip")
	}

	// 业务计算
	// 如果有新的位置信息，则更新当前位置，如果没有，则使用之前的位置
	if tr.Trip.Current == nil {
		s.Logger.Error("trip without current set", zap.String("id", tid.String()))
		return nil, status.Error(codes.Internal, "")
	}
	cur := tr.Trip.Current.Location
	if req.Current != nil {
		cur = req.Current
	}
	tr.Trip.Current = s.calcCurrentStatus(c, tr.Trip.Current, cur)

	if req.EndTrip {
		tr.Trip.End = tr.Trip.Current
		tr.Trip.Status = rentalpb.TripStatus_FINISHED
		err := s.CarManager.Lock(c, id.CarID(tr.Trip.CarId))
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "cannot lock car: %v", err)
		}
	}

	// 更新trip
	err = s.Mongo.UpdateTrip(c, tid, aid, tr.UpdatedAt, tr.Trip)
	if err != nil {
		return nil, status.Error(codes.Aborted, "")
	}
	return tr.Trip, nil
}

var nowFunc = func() int64 {
	return time.Now().Unix()
}

const centsPerSecond = 0.7

// 根据当前位置返回当前状态（计算距离，花费）
func (s *Service) calcCurrentStatus(c context.Context, last *rentalpb.LocationStatus, cur *rentalpb.Location) *rentalpb.LocationStatus {
	now := nowFunc()
	elapsedSec := float64(now - last.TimestampSec)

	dist, err := s.DistanceCalc.DistanceKm(c, last.Location, cur)
	if err != nil {
		s.Logger.Warn("cannot calculate distance", zap.Error(err))
	}

	poi, err := s.POIManager.Resolve(c, cur)
	if err != nil {
		s.Logger.Info("cannot resolve poi", zap.Stringer("location", cur), zap.Error(err))
	}

	return &rentalpb.LocationStatus{
		Location:     cur,
		FeeCent:      last.FeeCent + int32(centsPerSecond*elapsedSec*2*rand.Float64()),
		KmDriven:     last.KmDriven + dist,
		TimestampSec: now,
		PoiName:      poi,
	}

}
