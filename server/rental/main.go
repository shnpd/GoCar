package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	carpb "coolcar/car/api/gen/v1"
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
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"log"
	"time"

	"github.com/namsral/flag"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", ":8082", "address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo uri")
var blobAddr = flag.String("blob_addr", "localhost:8083", "address for blob service")
var aiAddr = flag.String("ai_addr", "localhost:18001", "address for ai service")
var carAddr = flag.String("car_addr", "localhost:8084", "address for car service")
var authPublicKeyFile = flag.String("auth_public_key_file", "../shared/auth/public.key", "public key file for auth")
var redisAddr = flag.String("redis_addr", "localhost:6379", "address for redis")

func main() {
	flag.Parse()
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	ac, err := grpc.Dial(*aiAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot connect aiservice", zap.Error(err))
	}

	db := mongoClient.Database("coolcar")

	blobConn, err := grpc.Dial(*blobAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot connect blob service", zap.Error(err))
	}

	profService := &profile.Service{
		BlobClient:        blobpb.NewBlobServiceClient(blobConn),
		PhotoGetExpire:    5 * time.Second,
		PhotoUploadExpire: 10 * time.Second,
		Mongo:             profiledao.NewMongo(db),
		Logger:            logger,
	}

	carConn, err := grpc.Dial(*carAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot connect car service", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              *addr,
		AuthPublicKeyFile: *authPublicKeyFile,
		Logger:            logger,
		ResigterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager: &car.Manager{
					CarService: carpb.NewCarServiceClient(carConn),
				},
				ProfileManager: &profClient.Manager{
					Fetcher: profService,
				},
				POIManager: &poi.Manager{},
				DistanceCalc: &ai.Client{
					AIClient: coolenvpb.NewAIServiceClient(ac),
				},
				Mongo:  tripdao.NewMongo(db),
				Logger: logger,
				RedisLock: redsync.New(goredis.NewPool(goredislib.NewClient(&goredislib.Options{
					Addr: *redisAddr,
				}))),
			})
			rentalpb.RegisterProfileServiceServer(s, profService)
		},
	}))
}
