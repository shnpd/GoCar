package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCreateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	db := mc.Database("coolcar")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("cannot setup indexes: %v", err)
	}
	m := NewMongo(db)

	// 表格驱动测试
	cases := []struct {
		name       string
		tripID     string
		accountID  string
		tripStatus rentalpb.TripStatus
		wantErr    bool
	}{
		{
			name:       "finished",
			tripID:     "675d37ff04c057f533ce5656",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "another_finished",
			tripID:     "675d37ff04c057f533ce5657",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "In_Progress",
			tripID:     "675d37ff04c057f533ce5658",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			name:       "another_in_progress",
			tripID:     "675d37ff04c057f533ce5659",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    true,
		},
		{
			name:       "in_progress_by_another_account",
			tripID:     "675d37ff04c057f533ce565a",
			accountID:  "account2",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	for _, cc := range cases {
		mgutil.NewObjIDWithValue(id.TripID(cc.tripID))

		tr, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: cc.accountID,
			Status:    cc.tripStatus,
		})

		if cc.wantErr {
			if err == nil {
				t.Errorf("%s:expect error, but got nil", cc.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s:cannot create trip: %v", cc.name, err)
			continue
		}
		if tr.ID.Hex() != cc.tripID {
			t.Errorf("%s:want trip id %q, got %q", cc.name, cc.tripID, tr.ID.Hex())
		}
	}

}
func TestGetTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	m := NewMongo(mc.Database("coolcar"))
	acct := id.AccountID("account2")
	mgutil.NewObjID = primitive.NewObjectID
	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: acct.String(),
		CarId:     "car1",
		Start: &rentalpb.LocationStatus{
			PoiName: "startpoint",
			Location: &rentalpb.Location{
				Latitude:  30,
				Longitude: 120,
			},
		},
		End: &rentalpb.LocationStatus{
			PoiName:  "endpoint",
			FeeCent:  10000,
			KmDriven: 35,
			Location: &rentalpb.Location{
				Latitude:  35,
				Longitude: 115,
			},
		},
		Status: rentalpb.TripStatus_FINISHED,
	})
	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}

	got, err := m.GetTrip(c, objid.ToTripID(tr.ID), acct)
	if err != nil {
		t.Errorf("cannot get trip: %v", err)
	}

	if diff := cmp.Diff(tr, got, protocmp.Transform()); diff != "" {
		t.Errorf("diff differs; -want +got: %s", diff)
	}
}

func TestGetTrips(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	m := NewMongo(mc.Database("coolcar"))

	// 插入数据
	rows := []struct {
		id        id.TripID
		accountID id.AccountID
		status    rentalpb.TripStatus
	}{
		{
			id:        "675d3904556c6773ffbd1718",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "675d3904556c6773ffbd1719",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "675d3904556c6773ffbd171a",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "675d3904556c6773ffbd171b",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			id:        "675d3904556c6773ffbd171c",
			accountID: "account_id_for_get_trips_1",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	for _, r := range rows {
		mgutil.NewObjIDWithValue(id.TripID(r.id))
		_, err = m.CreateTrip(c, &rentalpb.Trip{
			AccountId: r.accountID.String(),
			Status:    r.status,
		})
		if err != nil {
			t.Fatalf("cannot create trip: %v", err)
		}
	}

	// 测试方法
	cases := []struct {
		name       string
		accountID  string
		status     rentalpb.TripStatus
		wantCount  int
		wantOnlyID string
	}{
		{
			name:      "get_all",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		},
		{
			name:       "get_in_progress",
			accountID:  "account_id_for_get_trips",
			status:     rentalpb.TripStatus_IN_PROGRESS,
			wantCount:  1,
			wantOnlyID: "675d3904556c6773ffbd171b",
		},
	}
	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			res, err := m.GetTrips(context.Background(), id.AccountID(cc.accountID), cc.status)
			if err != nil {
				t.Errorf("cannot get trips: %v", err)
			}
			if len(res) != cc.wantCount {
				t.Errorf("incorrect result count: want %d trips, got %d", cc.wantCount, len(res))
			}
			if cc.wantOnlyID != "" && len(res) > 0 {
				if res[0].ID.Hex() != cc.wantOnlyID {
					t.Errorf("want trip id %q, got %q", cc.wantOnlyID, res[0].ID.Hex())
				}
			}
		})
	}
}

func TestUpdateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	tid := id.TripID("675d3914556c6773ffbd171b")
	aid := id.AccountID("account_for_update")
	mgutil.NewObjIDWithValue(tid)

	var now int64 = 10000
	mgutil.UpdatedAt = func() int64 {
		return now
	}

	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi",
		},
	})
	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}
	if tr.UpdatedAt != 10000 {
		t.Fatalf("want updateAt 10000, got %d", tr.UpdatedAt)
	}

	update := &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi_updated",
		},
	}

	cases := []struct {
		name string
		// 当前实际更新的时间
		now int64
		// 获取trip的时间
		withUpdatedAt int64
		wantErr       bool
	}{
		// 第一次更新可以通过
		{
			name:          "normal_update",
			now:           20000,
			withUpdatedAt: 10000,
		},
		// 第二次更新获取时间为10000的数据，但是10000的数据在20000更新过，此次在30000更新会报错
		{
			name:          "update_with_old_timestamp",
			now:           30000,
			withUpdatedAt: 10000,
			wantErr:       true,
		},
		// 第三次更新获取时间为20000的数据，此次在40000更新可以通过
		{
			name:          "update_with_refetch",
			now:           40000,
			withUpdatedAt: 20000,
		},
	}

	for _, cc := range cases {
		now = cc.now
		err := m.UpdateTrip(c, tid, aid, update)
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: expect error, but got nil", cc.name)
			} else {
				continue
			}
		} else {
			if err != nil {
				t.Errorf("%s: cannot update: %v", cc.name, err)
			}
		}
		updatedTrip, err := m.GetTrip(c, tid, aid)
		if err != nil {
			t.Errorf("%s: cannot get trip after update: %v", cc.name, err)
		}

		if cc.now != updatedTrip.UpdatedAt {
			t.Errorf("%s: want updated at %d, got %d", cc.name, cc.now, updatedTrip.UpdatedAt)
		}
	}
}

// TestMain 用于测试前的初始化
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
