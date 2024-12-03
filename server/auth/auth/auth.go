package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	OpenIdResolver OpenIdResolver
	Mongo          *dao.Mongo
	Logger         *zap.Logger
}

type OpenIdResolver interface {
	Resolve(code string) (string, error)
}

func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	openID, err := s.OpenIdResolver.Resolve(req.Code)
	if err != nil {
		// 将err返回给用户
		return nil, status.Errorf(codes.Unavailable, "cannot resolve openid: %v", err)
	}
	accountID, err := s.Mongo.ResolveAccountID(c, openID)
	if err != nil {
		// 这个Login的报错是很外层的直接给用户看，不希望把err直接给用户，同时记日志
		s.Logger.Error("cannot resolve account id", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	s.Logger.Info("received code", zap.String("code", req.Code))
	return &authpb.LoginResponse{
		AccessToken: "token for account id: " + accountID,
		ExpiresIn:   7200,
	}, nil
}
