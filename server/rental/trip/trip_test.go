package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/client/poi"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"coolcar/shared/server"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestCreateTrip(t *testing.T) {
	c := context.Background()
	pm := &profileManager{}
	cm := &carManager{}
	s := newService(c, t, pm, cm)

	nowFunc = func() int64 {
		return 1734440597
	}

	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  32.123,
			Longitude: 114.2525,
		},
	}

	pm.iID = "identity1"
	golden := `{"account_id":%q,"car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"五道口","timestamp_sec":1734440597},"current":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"五道口","timestamp_sec":1734440597},"status":1,"identity_id":"identity1"}`
	cases := []struct {
		name         string
		accountID    string
		tripID       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name:      "normal_create",
			accountID: "account1",
			tripID:    "675d3904556c6773ffbd1718",
			want:      fmt.Sprintf(golden, "account1"),
		},
		{
			name:       "profile_err",
			accountID:  "account2",
			tripID:     "675d3904556c6773ffbd1719",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		},
		{
			name:         "car_verify_err",
			accountID:    "account3",
			tripID:       "675d3904556c6773ffbd171a",
			carVerifyErr: fmt.Errorf("verify"),
			wantErr:      true,
		},
		{
			name:         "car_unlock_err",
			accountID:    "account4",
			tripID:       "675d3904556c6773ffbd171b",
			carUnlockErr: fmt.Errorf("unlock"),
			wantErr:      false, //解锁失败不影响行程创建
			want:         fmt.Sprintf(golden, "account4"),
		},
	}

	for _, cc := range cases {
		// 没有先后顺序，所以可以使用subtest
		t.Run(cc.name, func(t *testing.T) {
			mgutil.NewObjIDWithValue(id.TripID(cc.tripID))
			// 结构体s中保存的是pm和cm的指针，所以修改pm和cm的值会直接修改s
			pm.err = cc.profileErr
			cm.unlockError = cc.carUnlockErr
			cm.verifyError = cc.carVerifyErr
			c := auth.ContextWithAccountID(context.Background(), id.AccountID(cc.accountID))
			res, err := s.CreateTrip(c, req)
			if cc.wantErr {
				if err == nil {
					t.Errorf("want error; got none")
				} else {
					return
				}
			}
			// 希望没错误，但是有错误
			if err != nil {
				t.Errorf("error creating trip: %v", err)
				return
			}
			// 希望无错实际无错，进一步验证结果
			if res.Id != cc.tripID {
				t.Errorf("incorrect id: want %q, got %q", cc.tripID, res.Id)
			}
			b, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("failed to marshal trip: %v", err)
			}
			if cc.want != string(b) {
				t.Errorf("incorrect trip: want trip %s, got %s", cc.want, string(b))
			}
		})
	}

}

func TestTripLifeCycle(t *testing.T) {
	c := auth.ContextWithAccountID(context.Background(), id.AccountID("account_for_lifecycle"))
	s := newService(c, t, &profileManager{}, &carManager{})

	tid := id.TripID("675d3904556c6783ffbd171a")
	// 固定tripid
	mgutil.NewObjIDWithValue(tid)

	cases := []struct {
		name    string
		now     int64
		op      func() (*rentalpb.Trip, error)
		want    string
		wantErr bool
	}{
		{
			name: "create_trip",
			now:  10000,
			op: func() (*rentalpb.Trip, error) {
				e, err := s.CreateTrip(c, &rentalpb.CreateTripRequest{
					CarId: "car1",
					Start: &rentalpb.Location{
						Latitude:  32.123,
						Longitude: 114.2525,
					},
				})
				if err != nil {
					return nil, err
				}
				return e.Trip, nil
			},
			want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"五道口","timestamp_sec":10000},"current":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"五道口","timestamp_sec":10000},"status":1}`,
		},
		{
			name: "update_trip",
			now:  20000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id: tid.String(),
					Current: &rentalpb.Location{
						Latitude:  35.123,
						Longitude: 116.2525,
					},
				})
			},
			want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"五道口","timestamp_sec":10000},"current":{"location":{"latitude":35.123,"longitude":116.2525},"fee_cent":7968,"km_driven":100,"poi_name":"上地","timestamp_sec":20000},"status":1}`,
		},
		{
			name: "finish_trip",
			now:  30000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id:      tid.String(),
					EndTrip: true,
				})
			},
			want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"五道口","timestamp_sec":10000},"current":{"location":{"latitude":35.123,"longitude":116.2525},"fee_cent":11825,"km_driven":100,"poi_name":"上地","timestamp_sec":30000},"end":{"location":{"latitude":35.123,"longitude":116.2525},"fee_cent":11825,"km_driven":100,"poi_name":"上地","timestamp_sec":30000},"status":2}`,
		},
		{
			name: "query_trip",
			now:  40000,
			op: func() (*rentalpb.Trip, error) {
				return s.GetTrip(c, &rentalpb.GetTripRequest{
					Id: tid.String(),
				})
			},
			want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"五道口","timestamp_sec":10000},"current":{"location":{"latitude":35.123,"longitude":116.2525},"fee_cent":11825,"km_driven":100,"poi_name":"上地","timestamp_sec":30000},"end":{"location":{"latitude":35.123,"longitude":116.2525},"fee_cent":11825,"km_driven":100,"poi_name":"上地","timestamp_sec":30000},"status":2}`,
		},
		{
			name: "update_after_finished",
			now:  50000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id: tid.String(),
				})
			},
			want:    "",
			wantErr: true,
		},
	}
	// 固定随机数
	rand.Seed(1345)

	for _, cc := range cases {
		nowFunc = func() int64 {
			return cc.now
		}
		trip, err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want error; got none", cc.name)
			} else {
				continue
			}
		}

		if err != nil {
			t.Errorf("%s: operation failed: %v", cc.name, err)
			continue
		}
		b, err := json.Marshal(trip)
		if err != nil {
			t.Errorf("%s: failed to marshal trip: %v", cc.name, err)
		}
		got := string(b)
		if cc.want != got {
			t.Errorf("%s: incorrect trip: want %s, got %s", cc.name, cc.want, got)
		}
	}
}

func newService(c context.Context, t *testing.T, pm ProfileManager, cm CarManager) *Service {
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("failed to create mongo client: %v", err)
	}
	logger, err := server.NewZapLogger()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	db := mc.Database("coolcar")
	mongotesting.SetupIndexes(c, db)
	return &Service{
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		DistanceCalc:   &distCalc{},
		Mongo:          dao.NewMongo(db),
		Logger:         logger,
	}
}

type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(context.Context, id.AccountID) (id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	verifyError error
	unlockError error
}

func (m *carManager) Verify(context.Context, id.CarID, *rentalpb.Location) error {
	return m.verifyError
}

func (m *carManager) Unlock(c context.Context, cid id.CarID, aid id.AccountID, tid id.TripID, avatarURL string) error {
	return m.unlockError
}
func (m *carManager) Lock(c context.Context, cid id.CarID) error {
	return nil
}

type distCalc struct{}

func (*distCalc) DistanceKm(c context.Context, from *rentalpb.Location, to *rentalpb.Location) (float64, error) {
	if from.Latitude == to.Latitude && from.Longitude == to.Longitude {
		return 0, nil
	}
	return 100, nil
}
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
