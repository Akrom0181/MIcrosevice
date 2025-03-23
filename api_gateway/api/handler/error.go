package handler

import (
	"net/http"
	"strings"
	"user_api_gateway/config"
	"user_api_gateway/genproto/user_service"
	"user_api_gateway/pkg/logger"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func (h handler) HandleDbError(c *gin.Context, err error, message string) bool {
	if err == nil {
		return false
	}

	h.log.Error(message, logger.Error(err))
	var errorResponse user_service.ErrorResponse
	statusCode := http.StatusInternalServerError

	if err == pgx.ErrNoRows {
		errorResponse = user_service.ErrorResponse{
			Message: "The requested resource was not found.",
			Code:    config.ErrorNotFound,
		}
		c.JSON(http.StatusNotFound, &errorResponse)
		return true
	}

	switch e := err.(type) {
	case *pgconn.PgError:
		// Handle PostgreSQL-specific errors
		switch e.Code {
		case "23505":
			// Unique constraint violation
			errorResponse = user_service.ErrorResponse{
				Message: "Duplicate key error (unique constraint violation).",
				Code:    config.ErrorDuplicateKey,
			}
			statusCode = http.StatusBadRequest
		case "23503":
			// Foreign key violation
			errorResponse = user_service.ErrorResponse{
				Message: "The record could not be deleted because it is used in other records.",
				Code:    config.ErrorConflict,
			}
			statusCode = http.StatusBadRequest
		case "22001":
			// Value too long for column
			errorResponse = user_service.ErrorResponse{
				Message: "Value too long for column.",
				Code:    config.ErrorInvalidRequest,
			}
			statusCode = http.StatusBadRequest
		case "22005":
			// No rows in result set
			errorResponse = user_service.ErrorResponse{
				Message: "No rows in result set.",
				Code:    config.ErrorNotFound,
			}
			statusCode = http.StatusNotFound
		default:
			// General PostgreSQL error
			errorResponse = user_service.ErrorResponse{
				Message: "Ooops! Something went wrong.",
				Code:    config.ErrorInternalServer,
			}
		}
	default:
		if strings.Contains(err.Error(), "BAD_REQUEST") {
			errorResponse = user_service.ErrorResponse{
				Message: strings.TrimPrefix(err.Error(), "BAD_REQUEST"),
				Code:    config.ErrorBadRequest,
			}
		} else {
			// General PostgreSQL error
			errorResponse = user_service.ErrorResponse{
				Message: "Ooops! Something went wrong.",
				Code:    config.ErrorInternalServer,
			}
		}
	}

	c.JSON(statusCode, &errorResponse)
	return true
}

func (h handler) ReturnError(c *gin.Context, code string, message string, statusCode int) {
	h.log.Error(message, logger.String("code", code))
	errorResponse := user_service.ErrorResponse{
		Message: message,
		Code:    code,
	}
	c.JSON(statusCode, &errorResponse)
}
