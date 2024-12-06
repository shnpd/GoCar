package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar"))
	if err != nil {
		panic(err)
	}
	// 获取collection
	col := mc.Database("coolcar").Collection("account")

	findRows(c, col)
}
func findRows(c context.Context, col *mongo.Collection) {
	cur, err := col.Find(c, bson.M{})
	if err != nil {
		panic(err)
	}
	for cur.Next(c) {
		var row struct {
			ID     primitive.ObjectID `bson:"_id"`
			OpenID string             `bson:"open_id"`
		}
		err := cur.Decode(&row)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", row)
	}
}

func insertRows(c context.Context, col *mongo.Collection) {
	// 插入数据
	res, err := col.InsertMany(c, []interface{}{
		// bson是mongodb的二进制json格式可以节省空间
		bson.M{
			"open_id": "123",
		},
		bson.M{
			"open_id": "456",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}

// func main() {
// 	c := context.Background()
// 	// 建立数据库连接
// 	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	// 获取操作表
// 	col := mc.Database("coolcar").Collection("account")
// 	// insertRows(c, col)
// 	findRows(c, col)
// }

// func findRows(c context.Context, col *mongo.Collection) {
// 	res := col.FindOne(c, bson.M{
// 		"open_id": "123",
// 	})

// 	var row struct {
// 		ID     primitive.ObjectID `bson:"_id"`
// 		OpenID string             `bson:"open_id"`
// 	}
// 	err := res.Decode(&row)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%+v\n", row)
// }
