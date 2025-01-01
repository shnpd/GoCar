package main

import (
	"context"
	"coolcar/car/amqpclt"
	carpb "coolcar/car/api/gen/v1"
	"coolcar/car/car"
	"coolcar/car/dao"
	"coolcar/shared/server"
	"log"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	db := mongoClient.Database("coolcar")

	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Fatal("cannot dial amqp", zap.Error(err))
	}
	pub, err := amqpclt.NewPublisher(amqpConn, "coolcar")
	if err != nil {
		logger.Fatal("cannot create publisher", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "car",
		Addr:   ":8084",
		Logger: logger,
		ResigterFunc: func(s *grpc.Server) {
			carpb.RegisterCarServiceServer(s, &car.Service{
				Mongo:     dao.NewMongo(db),
				Logger:    logger,
				Publisher: pub,
			})
		},
	}))
}
