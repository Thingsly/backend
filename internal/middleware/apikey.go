package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrInvalidAPIKey  = errors.New("invalid api key")
	ErrAPIKeyDisabled = errors.New("api key is disabled")
	ErrAPIKeyNotFound = errors.New("api key not found")
	ErrInternalServer = errors.New("internal server error")
)

// APIKey
type OpenAPIKey struct {
	ID        string    `gorm:"column:id;primary_key"`
	TenantID  string    `gorm:"column:tenant_id"`
	APIKey    string    `gorm:"column:api_key"`
	Status    int       `gorm:"column:status"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (OpenAPIKey) TableName() string {
	return "open_api_keys"
}

// APIKeyInfo
type APIKeyInfo struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Status   int    `json:"status"`
	Name     string `json:"name"`
}

// APIKeyValidator
type APIKeyValidator struct {
	db          *gorm.DB
	redisClient *redis.Client
	ctx         context.Context
}

func NewAPIKeyValidator(db *gorm.DB, redisClient *redis.Client) *APIKeyValidator {
	return &APIKeyValidator{
		db:          db,
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// Validate APIKey
func (v *APIKeyValidator) ValidateAPIKey(apiKey string) (*APIKeyInfo, error) {
	// 1. Get from Redis cache
	info, err := v.getFromCache(apiKey)
	if err == nil {
		if info.Status != 1 {
			return nil, ErrAPIKeyDisabled
		}
		return info, nil
	}

	// 2. Cache miss, query from the database
	var key OpenAPIKey
	err = v.db.Where("api_key = ?", apiKey).First(&key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAPIKeyNotFound
		}
		return nil, fmt.Errorf("query database error: %w", err)
	}

	// 3. Check APIKey status
	if key.Status != 1 {
		return nil, ErrAPIKeyDisabled
	}

	// 4. Construct cache information
	info = &APIKeyInfo{
		ID:       key.ID,
		TenantID: key.TenantID,
		Status:   key.Status,
		Name:     key.Name,
	}

	// 5. Update cache
	if err := v.setCache(apiKey, info); err != nil {
		// Cache update failure only logs the error, does not affect the validation result
		fmt.Printf("update cache error: %v\n", err)
	}

	return info, nil
}

// Retrieve APIKey information from the cache
func (v *APIKeyValidator) getFromCache(apiKey string) (*APIKeyInfo, error) {
	key := fmt.Sprintf("apikey:%s", apiKey)
	data, err := v.redisClient.Get(v.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var info APIKeyInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// Set APIKey cache
func (v *APIKeyValidator) setCache(apiKey string, info *APIKeyInfo) error {
	key := fmt.Sprintf("apikey:%s", apiKey)
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	// Set cache expiration to 5 minutes
	return v.redisClient.Set(v.ctx, key, data, 5*time.Minute).Err()
}

// Delete APIKey cache
func (v *APIKeyValidator) DeleteCache(apiKey string) error {
	key := fmt.Sprintf("apikey:%s", apiKey)
	return v.redisClient.Del(v.ctx, key).Err()
}

// ValidateAPIKeyMiddleware is a middleware that validates the API key
func ValidateAPIKeyMiddleware(validator *APIKeyValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(401, gin.H{"code": 40001, "message": "Missing API Key"})
			c.Abort()
			return
		}

		info, err := validator.ValidateAPIKey(apiKey)
		if err != nil {
			switch err {
			case ErrAPIKeyNotFound:
				c.JSON(401, gin.H{"code": 40001, "message": "Invalid API Key"})
			case ErrAPIKeyDisabled:
				c.JSON(401, gin.H{"code": 40002, "message": "API Key Disabled"})
			default:
				c.JSON(500, gin.H{"code": 50001, "message": "Internal Server Error"})
			}
			c.Abort()
			return
		}

		// Store APIKey information in the context
		c.Set("apikey_info", info)
		c.Next()
	}
}
