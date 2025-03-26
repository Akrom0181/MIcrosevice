package handler

import (
	"strconv"
	"user_api_gateway/config"
	"user_api_gateway/pkg/grpc_client"
	"user_api_gateway/pkg/logger"

	rediscache "github.com/golanguzb70/redis-cache"

	"github.com/gin-gonic/gin"
)

type handler struct {
	log        logger.Logger
	grpcClient *grpc_client.GrpcClient
	cfg        config.Config
	redis      rediscache.RedisCache
}

// HandlerV1Config ...
type HandlerConfig struct {
	Logger     logger.Logger
	GrpcClient *grpc_client.GrpcClient
	Cfg        config.Config
	Redis      rediscache.RedisCache
}

const (
	// ErrorCodeInvalidURL ...
	ErrorCodeInvalidURL = "INVALID_URL"
	// ErrorCodeInvalidJSON ...
	ErrorCodeInvalidJSON = "INVALID_JSON"
	// ErrorCodeInternal ...
	ErrorCodeInternal = "INTERNAL"
	// ErrorCodeUnauthorized ...
	ErrorCodeUnauthorized = "UNAUTHORIZED"
	// ErrorCodeAlreadyExists ...
	ErrorCodeAlreadyExists = "ALREADY_EXISTS"
	// ErrorCodeNotFound ...
	ErrorCodeNotFound = "NOT_FOUND"
	// ErrorCodeInvalidCode ...
	ErrorCodeInvalidCode = "INVALID_CODE"
	// ErrorBadRequest ...
	ErrorBadRequest = "BAD_REQUEST"
	// ErrorCodeForbidden ...
	ErrorCodeForbidden = "FORBIDDEN"
	// ErrorCodeNotApproved ...
	ErrorCodeNotApproved = "NOT_APPROVED"
	// ErrorCodeWrongClub ...
	ErrorCodeWrongClub = "WRONG_CLUB"
	// ErrorCodePasswordsNotEqual ...
	ErrorCodePasswordsNotEqual = "PASSWORDS_NOT_EQUAL"
)

var (
	SigningKey = []byte("FfLbN7pIEYe8@!EqrttOLiwa(H8)7Ddo")
)

// New ...
func New(c *HandlerConfig) *handler {
	return &handler{
		log:        c.Logger,
		grpcClient: c.GrpcClient,
		cfg:        c.Cfg,
		redis:      c.Redis,
	}
}

func ParsePageQueryParam(c *gin.Context) (uint64, error) {
	pageStr := c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.ParseUint(pageStr, 10, 30)
	if err != nil {
		return 0, err
	}

	if page == 0 {
		return 1, nil
	}

	return page, nil
}

func ParseLimitQueryParam(c *gin.Context) (uint64, error) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "2"
	}

	limit, err := strconv.ParseUint(limitStr, 10, 30)
	if err != nil {
		return 0, err
	}

	if limit == 0 {
		return 2, nil
	}

	return limit, nil
}
