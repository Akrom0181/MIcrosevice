package service

import (
	"context"
	"user_service/config"
	"user_service/genproto/user_service"
	"user_service/grpc/client"

	"user_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewUserService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *UserService {
	return &UserService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *UserService) Create(ctx context.Context, req *user_service.User) (*user_service.User, error) {
	s.log.Info("---CreateUser--->>>", logger.Any("req", req))

	resp, err := s.strg.User().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreateUser--->>>", logger.Error(err))
		return &user_service.User{}, err
	}

	return resp, nil
}

func (s *UserService) GetSingle(ctx context.Context, req *user_service.UserSingleRequest) (*user_service.User, error) {
	s.log.Info("---GetSingleUser--->>>", logger.Any("req", req))

	resp, err := s.strg.User().GetSingle(ctx, req)
	if err != nil {
		s.log.Error("---GetSingleUser--->>>", logger.Error(err))
		return &user_service.User{}, err
	}

	return resp, nil
}

func (s *UserService) GetList(ctx context.Context, req *user_service.GetListUserRequest) (*user_service.GetListUserResponse, error) {
	s.log.Info("---GetAllUsers--->>>", logger.Any("req", req))

	resp, err := s.strg.User().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllUsers--->>>", logger.Error(err))
		return &user_service.GetListUserResponse{}, err
	}

	return resp, nil
}

func (s *UserService) Update(ctx context.Context, req *user_service.User) (*user_service.User, error) {
	s.log.Info("---UpdateUser--->>>", logger.Any("req", req))

	resp, err := s.strg.User().Update(ctx, req)
	if err != nil {
		s.log.Error("---UpdateUser--->>>", logger.Error(err))
		return &user_service.User{}, err
	}

	return resp, nil
}

func (s *UserService) Delete(ctx context.Context, req *user_service.UserPrimaryKey) (*emptypb.Empty, error) {
	s.log.Info("---DeleteUser--->>>", logger.Any("req", req))

	_, err := s.strg.User().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeleteUser--->>>", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
