package api

import (
	"net/http"
	"user_api_gateway/api/handler"
	"user_api_gateway/config"
	"user_api_gateway/pkg/grpc_client"
	"user_api_gateway/pkg/logger"

	_ "user_api_gateway/api/docs" //for swagger

	"github.com/casbin/casbin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	rediscache "github.com/golanguzb70/redis-cache"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Config ...
type Config struct {
	Logger     logger.Logger
	GrpcClient *grpc_client.GrpcClient
	Cfg        config.Config
	Redis      rediscache.RedisCache
}

// NewRouter -.
// Swagger spec:
// @title       Go Microservice API
// @description This is a Go Microservice API
// @version     1.0
// @host        localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(cnf Config) *gin.Engine {
	r := gin.New()

	r.Static("/images", "./static/images")

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "*")
	r.Use(cors.New(config))

	handler := handler.New(
		&handler.HandlerConfig{
			Logger:     cnf.Logger,
			GrpcClient: cnf.GrpcClient,
			Cfg:        cnf.Cfg,
			Redis:      cnf.Redis,
		},
	)

	// Initialize Casbin enforcer
	e := casbin.NewEnforcer("config/rbac.conf", "config/policy.csv")

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Api gateway"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
		auth.POST("/verify-email", handler.VerifyEmail)
	}

	protected := r.Group("/")
	protected.Use(handler.AuthMiddleware(e))

	user := protected.Group("/user")
	{
		user.POST("/", handler.CreateUser)
		user.GET("/list", handler.GetUsers)
		user.GET("/:id", handler.GetUser)
		user.PUT("/", handler.UpdateUser)
		user.DELETE("/:id", handler.DeleteUser)
	}

	authProtected := protected.Group("/auth")
	{
		authProtected.POST("/logout", handler.Logout)

	}

	session := protected.Group("/session")
	{
		session.GET("/list", handler.GetSessions)
		session.GET("/:id", handler.GetSession)
		session.PUT("/", handler.UpdateSession)
		session.DELETE("/:id", handler.DeleteSession)
	}

	post := protected.Group("/post")
	{

		post.POST("/", handler.CreatePost)
		post.GET("/:id", handler.GetPost)
		post.GET("/list", handler.GetPosts)
		post.PUT("/", handler.UpdatePost)
		post.DELETE("/:id", handler.DeletePost)
	}

	// Swagger endpoint
	url := ginSwagger.URL("swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return r
}
