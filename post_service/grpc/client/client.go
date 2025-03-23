package client

import (
	"fmt"
	"post_service/config"
	"post_service/genproto/user_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceManagerI interface {
	UserService() user_service.UserServiceClient
}

type grpcClients struct {
	userService user_service.UserServiceClient
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {

	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.UserServicePort, cfg.UserServiceHost),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(52428800), grpc.MaxCallSendMsgSize(52428800)),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		userService: user_service.NewUserServiceClient(connUserService),
	}, nil
}

func (g *grpcClients) UserService() user_service.UserServiceClient {
	return g.userService
}
