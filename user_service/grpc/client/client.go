package client

import (
	"fmt"
	"user_service/config"
	"user_service/genproto/post_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceManagerI interface {
	PostService() post_service.PostServiceClient
}

type grpcClients struct {
	postService post_service.PostServiceClient
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {

	connPostService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.PostServicePort, cfg.PostServiceHost),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(52428800), grpc.MaxCallSendMsgSize(52428800)),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		postService: post_service.NewPostServiceClient(connPostService),
	}, nil
}

func (g *grpcClients) PostService() post_service.PostServiceClient {
	return g.postService
}
