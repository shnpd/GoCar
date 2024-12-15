package mgutil

import (
	"coolcar/shared/mongo/objid"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Common field names.
const (
	IDFieldName        = "_id"
	UpdatedAtFieldName = "updatedat"
)

// IDField defines the object id filed.
type IDField struct {
	ID primitive.ObjectID `bson:"_id"`
}

// UpdatedAtField defines the updatedat  field.
type UpdatedAtField struct {
	UpdatedAt int64 `bson:"updatedat`
}

// NewObjectID returns a new object id.
// 不直接使用primitive.NewObjectID是为了方便测试，测试时可以通过修改变量的值来固定的ObjectID
var NewObjID = primitive.NewObjectID

func NewObjIDWithValue(id fmt.Stringer) {
	NewObjID = func() primitive.ObjectID {
		return objid.MustFromID(id)
	}
}

// UpdatedAt returns the current time in unix nano.
var UpdatedAt = func() int64 {
	return time.Now().UnixNano()
}

// Set returns a $set update document.
func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}

func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}
