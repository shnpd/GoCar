package mq

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
)

type Subscriber interface {
	Subscribe(context.Context) (ch chan *carpb.CarEntity, cleanUp func(), err error)
}

type Publisher interface {
	Publish(context.Context, *carpb.CarEntity) error
}
