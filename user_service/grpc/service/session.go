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

type SessionService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewSessionService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *SessionService {
	return &SessionService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *SessionService) Create(ctx context.Context, req *user_service.Session) (*user_service.Session, error) {
	s.log.Info("---CreateSession--->>>", logger.Any("req", &req))

	resp, err := s.strg.Session().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreateSession--->>>", logger.Error(err))
		return &user_service.Session{}, err
	}

	return resp, nil
}

func (s *SessionService) GetSingle(ctx context.Context, req *user_service.SessionSingleRequest) (*user_service.Session, error) {
	s.log.Info("---GetSingleSession--->>>", logger.Any("req", req))

	resp, err := s.strg.Session().GetSingle(ctx, req)
	if err != nil {
		s.log.Error("---GetSingleSession--->>>", logger.Error(err))
		return &user_service.Session{}, err
	}

	return resp, nil
}

func (s *SessionService) GetList(ctx context.Context, req *user_service.GetListSessionRequest) (*user_service.GetListSessionResponse, error) {
	s.log.Info("---GetAllSessions--->>>", logger.Any("req", req))

	resp, err := s.strg.Session().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllSessions--->>>", logger.Error(err))
		return &user_service.GetListSessionResponse{}, err
	}

	return resp, nil
}

func (s *SessionService) Update(ctx context.Context, req *user_service.Session) (*user_service.Session, error) {
	s.log.Info("---UpdateSession--->>>", logger.Any("req", req))

	resp, err := s.strg.Session().Update(ctx, req)
	if err != nil {
		s.log.Error("---UpdateSession--->>>", logger.Error(err))
		return &user_service.Session{}, err
	}

	return resp, nil
}

func (s *SessionService) Delete(ctx context.Context, req *user_service.SessionSingleRequest) (*emptypb.Empty, error) {
	s.log.Info("---DeleteSession--->>>", logger.Any("req", req))

	_, err := s.strg.Session().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeleteSession--->>>", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
