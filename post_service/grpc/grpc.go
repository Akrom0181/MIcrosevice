package grpc

import (
	"post_service/config"
	"post_service/genproto/post_service"
	"post_service/grpc/client"
	"post_service/grpc/service"
	"post_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvc client.ServiceManagerI) (grpcServer *grpc.Server) {

	grpcServer = grpc.NewServer()

	post_service.RegisterPostServiceServer(grpcServer, service.NewPostService(cfg, log, strg, srvc))
	post_service.RegisterPostAttachmentServiceServer(grpcServer, service.NewPostAttachmentService(cfg, log, strg, srvc))
	reflection.Register(grpcServer)
	return
}
