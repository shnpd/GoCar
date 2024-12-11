package server

import (
	"coolcar/shared/auth"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ServerConfig defines a server configuration
type GRPCConfig struct {
	Name              string
	Addr              string
	AuthPublicKeyFile string
	ResigterFunc      func(*grpc.Server)
	Logger            *zap.Logger
}

// RunGRPCServer runs a gRPC server
func RunGRPCServer(c *GRPCConfig) error {
	nameFiled := zap.String("name", c.Name)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		c.Logger.Fatal("cannot listen", nameFiled, zap.Error(err))
	}

	var opts []grpc.ServerOption
	if c.AuthPublicKeyFile != "" {
		// 输入公钥文件，返回一个拦截器interceptor（登录请求不需要过拦截器）
		in, err := auth.Interceptor(c.AuthPublicKeyFile)
		if err != nil {
			c.Logger.Fatal("cannot create auth interceptor", nameFiled, zap.Error(err))
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}
	s := grpc.NewServer(opts...)

	// 注册服务比较复杂每个服务以及注册的方法都不同，所以将这部分提取出一个函数，参数只有s，函数由调用者传入
	c.ResigterFunc(s)

	c.Logger.Info("grpc server started", nameFiled, zap.String("addr", c.Addr))
	return s.Serve(lis)
}
