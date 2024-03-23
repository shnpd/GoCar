package main

import (
	"context"
	trippb "coolcar/proto/gen/go"
	trip "coolcar/tripservice"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(log.Lshortfile)
	go startGRPCGateway()
	// 监听端口
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		// 不使用panic，Fatalf：输出之后程序退出
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	trippb.RegisterTripServiceServer(s, &trip.Service{})
	// 如果s.Serve没有出错的话会一直监听不会退出
	log.Fatal(s.Serve(lis))
}

func startGRPCGateway() {
	// context在后面仔细讲，这里生成了一个没有具体内容的上下文，在这个上下文中连接后端的grpc服务
	c := context.Background()
	// 为上下文添加cancel的能力，cancel是一个函数，调用cancel连接就会被断开
	c, cancel := context.WithCancel(c)
	// 服务结束断开连接
	defer cancel()

	// mux: multiplexer 一对多，分发器
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			EnumsAsInts: true,
			OrigName:    true,
		},
	))
	// 通过context c连接,c在函数返回后就会被cancel掉,这个连接注册在NewServeMux上:8081,连接方式为insecure通过tcp明文连接
	err := trippb.RegisterTripServiceHandlerFromEndpoint(
		c,
		mux,
		"localhost:8081",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}

	// http监听地址,8081是tripservice的地址与gateway的地址是不同的
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("cannot listen and server: %v", err)
	}
}
