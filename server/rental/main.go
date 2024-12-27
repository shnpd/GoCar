package main

import (
	"context"
	"coolcar/rental/ai"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/profile"
	profiledao "coolcar/rental/profile/dao"
	"coolcar/rental/trip"
	"coolcar/rental/trip/client/car"
	"coolcar/rental/trip/client/poi"
	profClient "coolcar/rental/trip/client/profile"
	tripdao "coolcar/rental/trip/dao"
	coolenvpb "coolcar/shared/coolenv"
	"coolcar/shared/server"
	"log"

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

	ac, err := grpc.Dial("localhost:18001", grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot connect aiservice", zap.Error(err))
	}

	db := mongoClient.Database("coolcar")

	profService := &profile.Service{
		Mongo:  profiledao.NewMongo(db),
		Logger: logger,
	}
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: "../shared/auth/public.key",
		Logger:            logger,
		ResigterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager: &car.Manager{},
				ProfileManager: &profClient.Manager{
					Fetcher: profService,
				},
				POIManager: &poi.Manager{},
				DistanceCalc: &ai.Client{
					AIClient: coolenvpb.NewAIServiceClient(ac),
				},
				Mongo:  tripdao.NewMongo(db),
				Logger: logger,
			})
			rentalpb.RegisterProfileServiceServer(s, profService)
		},
	}))
}
