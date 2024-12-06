package dao

import (
	"context"
	mgo "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI string

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	m := NewMongo(mc.Database("coolcar"))
	// 初始化插入两条数据以及对应的objectID
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgo.IDField: mustObjID("674fe7b3f846790fba023980"),
			openIDField: "openid_1",
		},
		bson.M{
			mgo.IDField: mustObjID("674fe7b3f846790fba023981"),
			openIDField: "openid_2",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values: %v", err)
	}

	// 测试环境中每次都使用以下函数生成固定的ObjectID
	m.newObjID = func() primitive.ObjectID {
		return mustObjID("674fe7b3f846790fba023982")
	}

	// 表格驱动测试，创建了3个测试用例，openid_1和openid_2是已经存在的用户希望resolve得到插入的objectid，openid_3是新用户希望得到固定函数生成的objectid
	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "674fe7b3f846790fba023980",
		},
		{
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "674fe7b3f846790fba023981",
		},
		{
			name:   "new_user",
			openID: "openid_3",
			want:   "674fe7b3f846790fba023982",
		},
	}

	for _, cc := range cases {
		// t.Run()启动一个子测试
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.openID)
			if err != nil {
				t.Errorf("faild resolve account id for %q: %v", cc.openID, err)
			}
			if id != cc.want {
				t.Errorf("resolve account id: want: %q, got:%q", cc.want, id)
			}
		})
	}
}

func mustObjID(hex string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return objID
}

// TestMain 用于测试前的初始化
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}
