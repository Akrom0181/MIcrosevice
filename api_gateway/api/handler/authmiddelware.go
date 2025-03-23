package handler

import (
	"fmt"
	"net/http"
	"strings"
	"user_api_gateway/genproto/user_service"
	"user_api_gateway/pkg/jwt"

	"go.uber.org/zap"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

func (h *handler) AuthMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			userRole  string
			act       = c.Request.Method
			obj       = c.FullPath()
			sessionID string
		)

		token := c.GetHeader("Authorization")
		if token == "" {
			userRole = "unauthorized"
		}

		if userRole == "" {
			token = strings.TrimPrefix(token, "Bearer ")

			claims, err := jwt.ParseJWT(token, h.cfg.JWT)
			if err != nil {
				h.log.Error("Error parsing JWT", zap.Error(err))
				userRole = "unauthorized"
			} else {
				v, ok := claims["user_role"].(string)
				if !ok {
					userRole = "unauthorized"
				} else {
					userRole = v
				}

				// Extract session_id from claims and store it
				if sid, ok := claims["session_id"].(string); ok {
					sessionID = sid
					// Store in context and header
					c.Set("session_id", sessionID)
					c.Request.Header.Set("session_id", sessionID)
				} else {
					h.log.Error("Missing session_id in JWT claims")
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID", "code": "BAD_REQUEST"})
					return
				}

				// Set other claims as headers
				for key, value := range claims {
					c.Request.Header.Set(key, fmt.Sprintf("%v", value))
				}
			}
		}

		// Only verify session if user is authenticated
		if userRole != "unauthorized" && sessionID != "" {
			// Use the session ID extracted from JWT claims
			session, err := h.grpcClient.SessionService().GetSingle(c, &user_service.SessionSingleRequest{Id: sessionID})
			if err != nil {
				h.log.Error("Error getting session", zap.Error(err), zap.String("session_id", sessionID))
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID", "code": "BAD_REQUEST"})
				return
			}

			if !session.IsActive {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is not active"})
				return
			}
		}

		// Check permissions
		ok, err := e.EnforceSafe(userRole, obj, act)
		if err != nil {
			h.log.Error("Error enforcing", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}
