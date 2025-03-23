package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"
)

// Config struct holds all configuration settings
type Config struct {
	ServiceName string
	Environment string
	Version     string

	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	UserServiceHost string
	UserServicePort string

	PostServiceHost string
	PostServicePort string

	LogLevel string
	HTTPPort string

	GmailHost     string
	GmailPort     string
	GmailUser     string
	GmailPassword string

	JWT string

	RedisHost     string
	RedisPort     int
	RedisPassword string

	PostgresMaxConnections int32
}

// Load reads environment variables and returns a Config instance
func Load() Config {
	// Load .env file if it exists
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}

	return Config{
		ServiceName: cast.ToString(os.Getenv("SERVICE_NAME")),
		Environment: cast.ToString(os.Getenv("ENVIRONMENT")),
		Version:     cast.ToString(os.Getenv("VERSION")),

		PostgresHost:           cast.ToString(os.Getenv("POSTGRES_HOST")),
		PostgresPort:           cast.ToInt(os.Getenv("POSTGRES_PORT")),
		PostgresUser:           cast.ToString(os.Getenv("POSTGRES_USER")),
		PostgresPassword:       cast.ToString(os.Getenv("POSTGRES_PASSWORD")),
		PostgresDatabase:       cast.ToString(os.Getenv("POSTGRES_DATABASE")),
		PostgresMaxConnections: cast.ToInt32(os.Getenv("POSTGRES_MAX_CONNECTIONS")),

		RedisHost:     cast.ToString(os.Getenv("REDIS_HOST")),
		RedisPort:     cast.ToInt(os.Getenv("REDIS_PORT")),
		RedisPassword: cast.ToString(os.Getenv("REDIS_PASSWORD")),

		JWT: cast.ToString(os.Getenv("JWT")),

		GmailHost:     cast.ToString(os.Getenv("GMAIL_HOST")),
		GmailPort:     cast.ToString(os.Getenv("GMAIL_PORT")),
		GmailUser:     cast.ToString(os.Getenv("GMAIL_USER")),
		GmailPassword: cast.ToString(os.Getenv("GMAIL_PASSWORD")),

		UserServiceHost: cast.ToString(os.Getenv("USER_SERVICE_HOST")),
		UserServicePort: cast.ToString(os.Getenv("USER_SERVICE_PORT")),

		PostServiceHost: cast.ToString(os.Getenv("POST_SERVICE_HOST")),
		PostServicePort: cast.ToString(os.Getenv("POST_SERVICE_PORT")),

		LogLevel: cast.ToString(os.Getenv("LOG_LEVEL")),
		HTTPPort: cast.ToString(os.Getenv("HTTP_PORT")),
	}
}
