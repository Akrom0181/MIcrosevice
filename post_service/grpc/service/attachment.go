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

type PostAttachmentService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewPostAttachmentService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *PostAttachmentService {
	return &PostAttachmentService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *PostAttachmentService) MultipleUpsert(ctx context.Context, req *post_service.AttachmentMultipleInsertRequest) (*post_service.AttachmentList, error) {
	s.log.Info("---CreatePostAttachment--->>>", logger.Any("req", req))

	resp, err := s.strg.PostAttachment().MultipleUpsert(ctx, req)
	if err != nil {
		s.log.Error("---CreatePostAttachment--->>>", logger.Error(err))
		return &post_service.AttachmentList{}, err
	}

	return resp, nil
}

func (s *PostAttachmentService) Create(ctx context.Context, req *post_service.Attachment) (*post_service.Attachment, error) {
	s.log.Info("---CreatePostAttachment--->>>", logger.Any("req", req))

	resp, err := s.strg.PostAttachment().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreatePostAttachment--->>>", logger.Error(err))
		return &post_service.Attachment{}, err
	}

	return resp, nil
}

func (s *PostAttachmentService) GetSingle(ctx context.Context, req *post_service.AttachmentSingleRequest) (*post_service.Attachment, error) {
	s.log.Info("---GetSinglePostAttachment--->>>", logger.Any("req", req))

	resp, err := s.strg.PostAttachment().GetSingle(ctx, req)
	if err != nil {
		s.log.Error("---GetSinglePostAttachment--->>>", logger.Error(err))
		return &post_service.Attachment{}, err
	}

	return resp, nil
}

func (s *PostAttachmentService) GetList(ctx context.Context, req *post_service.GetListAttachmentRequest) (*post_service.AttachmentList, error) {
	s.log.Info("---GetAllPostAttachments--->>>", logger.Any("req", req))

	resp, err := s.strg.PostAttachment().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllPostAttachments--->>>", logger.Error(err))
		return &post_service.AttachmentList{}, err
	}

	return resp, nil
}

func (s *PostAttachmentService) Delete(ctx context.Context, req *post_service.AttachmentSingleRequest) (*emptypb.Empty, error) {
	s.log.Info("---DeletePostAttachment--->>>", logger.Any("req", req))

	_, err := s.strg.PostAttachment().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeletePostAttachment--->>>", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PostAttachmentService) GetDefaultTags(ctx context.Context, req *post_service.GetDefaultTagsRequest) (*post_service.GetDefaultTagsResponse, error) {
	s.log.Info("---GetDefaultTags--->>>", logger.Any("req", req))

	resp, err := s.strg.PostAttachment().GetDefaultTags(ctx, req)
	if err != nil {
		s.log.Error("---GetDefaultTags--->>>", logger.Error(err))
		return &post_service.GetDefaultTagsResponse{}, err
	}

	return resp, nil
}
