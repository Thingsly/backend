package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/HustIoTPlatform/backend/internal/dal"
	"github.com/HustIoTPlatform/backend/pkg/global"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Error code constants
const (
	ErrCodeNoAuth         = 40100 // Missing authentication information
	ErrCodeInvalidToken   = 40101 // Invalid Token
	ErrCodeTokenExpired   = 40102 // Token has expired
	ErrCodeInvalidAPIKey  = 40103 // Invalid APIKey
	ErrCodeAPIKeyDisabled = 40104 // APIKey has been disabled
)

// Unified error response structure
type ErrorResponse struct {
	Code      int    `json:"code"`                 // Error code
	Message   string `json:"message"`              // Error description
	RequestID string `json:"request_id,omitempty"` // Request ID, for easier tracking
}

// JWTAuth middleware, checks token and APIKey
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. First check the JWT token
		token := c.Request.Header.Get("x-token")
		if token != "" {
			// JWT Token exists, validate JWT
			if isValidJWT(c, token) {
				c.Next()
				return
			}
			// JWT validation failed, continue to try APIKey
		}

		// 2. Try APIKey validation
		if !OpenAPIKeyAuth(c) {
			// APIKey validation also failed, error response is already set in OpenAPIKeyAuth
			return
		}

		// APIKey validation successful
		c.Next()
	}
}

// isValidJWT Validates the JWT token
func isValidJWT(c *gin.Context, token string) bool {
	requestID := c.GetString("X-Request-ID")

	// Validate the token in Redis
	if global.REDIS.Get(context.Background(), token).Val() != "1" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:      ErrCodeTokenExpired,
			Message:   "token has expired",
			RequestID: requestID,
		})
		c.Abort()
		return false
	}

	// Refresh the token expiration time
	timeout := viper.GetInt("session.timeout")
	global.REDIS.Set(context.Background(), token, "1", time.Duration(timeout)*time.Minute)

	// Validate the JWT token
	key := viper.GetString("jwt.key")
	j := utils.NewJWT([]byte(key))
	claims, err := j.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:      ErrCodeInvalidToken,
			Message:   "invalid token format",
			RequestID: requestID,
		})
		c.Abort()
		return false
	}

	// Set claims in the context
	c.Set("claims", claims)
	return true
}

// OpenAPIKeyAuth APIKey validation
func OpenAPIKeyAuth(c *gin.Context) bool {
	requestID := c.GetString("X-Request-ID")

	appKey := c.Request.Header.Get("x-api-key")
	if appKey == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:      ErrCodeNoAuth,
			Message:   "missing authentication (x-token or x-api-key required)",
			RequestID: requestID,
		})
		c.Abort()
		return false
	}

	tenantID, err := dal.VerifyOpenAPIKey(context.Background(), appKey)
	if err != nil {
		errCode := ErrCodeInvalidAPIKey
		errMsg := "api key verification failed"

		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:      errCode,
			Message:   errMsg,
			RequestID: requestID,
		})
		c.Abort()
		return false
	}

	// Set claims in the context
	claims := utils.UserClaims{
		TenantID:  tenantID,
		Authority: "TENANT_ADMIN",
	}
	c.Set("claims", &claims)
	return true
}
