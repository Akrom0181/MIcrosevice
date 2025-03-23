package grpc_client

import (
	"fmt"
	"log"
	"user_api_gateway/config"

	ps "user_api_gateway/genproto/post_service"
	us "user_api_gateway/genproto/user_service"

	rediscache "github.com/golanguzb70/redis-cache"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClientI ..90.
type GrpcClientI interface {
	UserService() us.UserServiceClient
	PostService() ps.PostServiceClient
	SessionService() us.SessionServiceClient
	PostAttachment() ps.PostAttachmentServiceClient
}

// GrpcClient ...
type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

// New ...
func New(cfg config.Config, redis rediscache.RedisCache) (*GrpcClient, error) {

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("user service dial host: %s port:%s err: %s",
			cfg.UserServiceHost, cfg.UserServicePort, err)
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("user service dial host: %s port:%s err: %s",
			cfg.UserServiceHost, cfg.UserServicePort, err)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"user_service":           us.NewUserServiceClient(connUser),
			"session_service":        us.NewSessionServiceClient(connUser),
			"post_service":           ps.NewPostServiceClient(connPost),
			"postattachment_service": ps.NewPostAttachmentServiceClient(connPost),
		},
	}, nil
}

func (g *GrpcClient) UserService() us.UserServiceClient {
	client, ok := g.connections["user_service"].(us.UserServiceClient)
	if !ok {
		log.Println("failed to assert type for user")
		return nil
	}
	return client
}

func (g *GrpcClient) SessionService() us.SessionServiceClient {
	client, ok := g.connections["session_service"].(us.SessionServiceClient)
	if !ok {
		log.Println("failed to assert type for session")
		return nil
	}
	return client
}

func (g *GrpcClient) PostService() ps.PostServiceClient {
	client, ok := g.connections["post_service"].(ps.PostServiceClient)
	if !ok {
		log.Println("failed to assert type for post")
		return nil
	}
	return client
}

func (g *GrpcClient) PostAttachmentService() ps.PostAttachmentServiceClient {
	client, ok := g.connections["postattachment_service"].(ps.PostAttachmentServiceClient)
	if !ok {
		log.Println("failed to assert type for post_attachment")
		return nil
	}
	return client
}

func (g *GrpcClient) CloseConnections() {
	for key, conn := range g.connections {
		if c, ok := conn.(*grpc.ClientConn); ok {
			err := c.Close()
			if err != nil {
				log.Printf("failed to close connection for %s: %v", key, err)
			}
		}
	}
}
