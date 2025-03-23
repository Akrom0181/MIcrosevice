package service

import (
	"context"
	"post_service/config"
	"post_service/genproto/post_service"
	"post_service/grpc/client"
	"post_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PostService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewPostService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *PostService {
	return &PostService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *PostService) Create(ctx context.Context, req *post_service.Post) (*post_service.Post, error) {
	s.log.Info("---CreatePost--->>>", logger.Any("req", req))

	resp, err := s.strg.Post().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreatePost--->>>", logger.Error(err))
		return &post_service.Post{}, err
	}

	return resp, nil
}

func (s *PostService) GetSingle(ctx context.Context, req *post_service.PostSingleRequest) (*post_service.Post, error) {
	s.log.Info("---GetSinglePost--->>>", logger.Any("req", req))

	resp, err := s.strg.Post().GetSingle(ctx, req)
	if err != nil {
		s.log.Error("---GetSinglePost--->>>", logger.Error(err))
		return &post_service.Post{}, err
	}

	return resp, nil
}

func (s *PostService) GetList(ctx context.Context, req *post_service.GetListPostRequest) (*post_service.PostList, error) {
	s.log.Info("---GetAllPosts--->>>", logger.Any("req", req))

	resp, err := s.strg.Post().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllPosts--->>>", logger.Error(err))
		return &post_service.PostList{}, err
	}

	return resp, nil
}

func (s *PostService) Update(ctx context.Context, req *post_service.Post) (*post_service.Post, error) {
	s.log.Info("---UpdatePost--->>>", logger.Any("req", req))

	resp, err := s.strg.Post().Update(ctx, req)
	if err != nil {
		s.log.Error("---UpdatePost--->>>", logger.Error(err))
		return &post_service.Post{}, err
	}

	return resp, nil
}

func (s *PostService) Delete(ctx context.Context, req *post_service.PostSingleRequest) (*emptypb.Empty, error) {
	s.log.Info("---DeletePost--->>>", logger.Any("req", req))

	_, err := s.strg.Post().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeletePost--->>>", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PostService) GetDefaultTags(ctx context.Context, req *post_service.GetDefaultTagsRequest) (*post_service.GetDefaultTagsResponse, error) {
	s.log.Info("---GetDefaultTags--->>>", logger.Any("req", req))

	resp, err := s.strg.PostAttachment().GetDefaultTags(ctx, req)
	if err != nil {
		s.log.Error("---GetDefaultTags--->>>", logger.Error(err))
		return &post_service.GetDefaultTagsResponse{}, err
	}

	return resp, nil
}
