package main

import (
	"fmt"
	"user_api_gateway/api"
	"user_api_gateway/config"
	"user_api_gateway/pkg/grpc_client"
	"user_api_gateway/pkg/logger"

	rediscache "github.com/golanguzb70/redis-cache"
)

var (
	log        logger.Logger
	cfg        config.Config
	grpcClient *grpc_client.GrpcClient
	redis      rediscache.RedisCache
)

// initDeps initializes dependencies like config, logger, Redis, and gRPC client
func initDeps() {
	var err error
	cfg = config.Load()

	log = logger.New(cfg.LogLevel, "user-api-gateway")

	redis, err = rediscache.New(&rediscache.Config{
		RedisHost: cfg.RedisHost,
		RedisPort: cfg.RedisPort,
	})
	if err != nil {
		log.Fatal("redis error", logger.Error(err))
	}

	grpcClient, err = grpc_client.New(cfg, redis)
	if err != nil {
		log.Fatal("grpc dial error", logger.Error(err))
	}
}

func main() {
	initDeps()

	server := api.New(api.Config{
		Logger:     log,
		GrpcClient: grpcClient,
		Cfg:        cfg,
		Redis:      redis,
	})

	fmt.Println("Starting server on port", cfg.HTTPPort)
	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("Server failed", logger.Error(err))
	}
}
