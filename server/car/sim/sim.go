package sim

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"time"

	"go.uber.org/zap"
)

type Subscriber interface {
	Subscribe(context.Context) (ch chan *carpb.CarEntity, cleanUp func(), err error)
}
type Controller struct {
	CarService carpb.CarServiceClient
	Subscriber Subscriber
	Logger     *zap.Logger
}

func (c *Controller) RunSimulations(ctx context.Context) {
	var cars []*carpb.CarEntity
	// 调用GetCars时car service可能还未启动，如果获取失败就等3s再重试
	for {
		time.Sleep(3 * time.Second)
		res, err := c.CarService.GetCars(ctx, &carpb.GetCarsRequest{})
		if err != nil {
			c.Logger.Error("cannot get cars", zap.Error(err))
			continue
		}
		cars = res.Cars
		break
	}
	c.Logger.Info("Running car simulations", zap.Int("car_count", len(cars)))

	res, err := c.CarService.GetCars(ctx, &carpb.GetCarsRequest{})
	if err != nil {
		c.Logger.Error("cannot get cars", zap.Error(err))
		return
	}

	msgCh, cleanUp, err := c.Subscriber.Subscribe(ctx)
	defer cleanUp()

	if err != nil {
		c.Logger.Error("cannot subscribe", zap.Error(err))
		return
	}

	// 创建车辆id到channel的映射，启动goroutine
	carChans := make(map[string]chan *carpb.Car)
	for _, car := range res.Cars {
		ch := make(chan *carpb.Car)
		carChans[car.Id] = ch
		go c.SimulateCar(context.Background(), car, ch)
	}

	// 从消息队列中读取消息，根据消息中的车辆id找到对应的channel，将消息发送到channel中
	for carUpdate := range msgCh {
		ch := carChans[carUpdate.Id]
		if ch != nil {
			ch <- carUpdate.Car
		}
	}
}

func (c *Controller) SimulateCar(ctx context.Context, initial *carpb.CarEntity, ch chan *carpb.Car) {
	carID := initial.Id
	c.Logger.Info("Simulating car", zap.String("car_id", carID))
	// 从channel中读取消息，根据消息的状态更新车辆状态
	for update := range ch {
		if update.Status == carpb.CarStatus_UNLOCKING {
			// 硬件开锁过程，假设开锁成功
			_, err := c.CarService.UpdateCar(ctx, &carpb.UpdateCarRequest{
				Id:     carID,
				Status: carpb.CarStatus_UNLOCKED,
			})
			if err != nil {
				c.Logger.Error("cannot unlock car", zap.Error(err))
			}
		} else if update.Status == carpb.CarStatus_LOCKING {
			// 硬件开锁过程，假设开锁成功
			_, err := c.CarService.UpdateCar(ctx, &carpb.UpdateCarRequest{
				Id:     carID,
				Status: carpb.CarStatus_LOCKED,
			})
			if err != nil {
				c.Logger.Error("cannot lock car", zap.Error(err))
			}
		}
	}
}
