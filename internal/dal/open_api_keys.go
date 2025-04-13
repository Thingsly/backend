// internal/dal/open_api_keys.go
package dal

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gen"

	model "github.com/HustIoTPlatform/backend/internal/model"
	query "github.com/HustIoTPlatform/backend/internal/query"
	global "github.com/HustIoTPlatform/backend/pkg/global"
)

// CreateOpenAPIKey 
func CreateOpenAPIKey(key *model.OpenAPIKey) error {
	return query.OpenAPIKey.Create(key)
}

// GetOpenAPIKeyByID 
func GetOpenAPIKeyByID(id string) (*model.OpenAPIKey, error) {
	return query.OpenAPIKey.Where(query.OpenAPIKey.ID.Eq(id)).First()
}

// GetOpenAPIKeyByAppKey
func GetOpenAPIKeyByAppKey(appKey string) (*model.OpenAPIKey, error) {
	return query.OpenAPIKey.Where(query.OpenAPIKey.APIKey.Eq(appKey)).First()
}

// GetOpenAPIKeyListByPage - Paginate to get the OpenAPI key list
// param listReq - Query parameters
// param tenantID - Tenant ID, used for permission filtering
// return - Total count, data list, error message
func GetOpenAPIKeyListByPage(listReq *model.OpenAPIKeyListReq, tenantID string) (int64, interface{}, error) {
	q := query.OpenAPIKey
	queryBuilder := q.WithContext(context.Background())

	// Add tenant filtering
	if tenantID != "" {
		queryBuilder = queryBuilder.Where(q.TenantID.Eq(tenantID))
	}

	// Add query conditions
	if listReq.Status != nil {
		queryBuilder = queryBuilder.Where(q.Status.Eq(*listReq.Status))
	}

	// Get total count
	count, err := queryBuilder.Count()
	if err != nil {
		return 0, nil, err
	}

	// Pagination handling
	if listReq.Page != 0 && listReq.PageSize != 0 {
		queryBuilder = queryBuilder.Limit(listReq.PageSize)
		queryBuilder = queryBuilder.Offset((listReq.Page - 1) * listReq.PageSize)
	}

	// Execute query
	keys, err := queryBuilder.Order(q.CreatedAt.Desc()).Find()
	if err != nil {
		return 0, nil, err
	}

	return count, keys, nil
}

// UpdateOpenAPIKey - Update OpenAPI key information
// param id - Key ID
// param updates - Fields to be updated
func UpdateOpenAPIKey(id string, updates map[string]interface{}) error {
	q := query.OpenAPIKey
	updates["updated_at"] = time.Now()
	_, err := q.Where(q.ID.Eq(id)).Updates(updates)
	return err
}

// DeleteOpenAPIKey - Delete OpenAPI key
// param id - Key ID
// @note - Deleting will also clean the Redis cache
func DeleteOpenAPIKey(id string) error {
	// Delete database record
	_, err := query.OpenAPIKey.Where(query.OpenAPIKey.ID.Eq(id)).Delete()
	if err != nil {
		return err
	}

	// Clean up cache
	cacheKey := "openapi:key:" + id
	err = global.REDIS.Del(context.Background(), cacheKey).Err()
	if err != nil {
		logrus.Warnf("Failed to delete OpenAPI key cache: %v", err)
	}

	return nil
}

// OpenAPIKeyQuery - OpenAPI key query structure
type OpenAPIKeyQuery struct{}

// Count - Get the total count of OpenAPI keys
func (OpenAPIKeyQuery) Count(ctx context.Context, option ...gen.Condition) (count int64, err error) {
	count, err = query.OpenAPIKey.WithContext(ctx).Where(option...).Count()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

// Select - Query the OpenAPI key list based on conditions
func (OpenAPIKeyQuery) Select(ctx context.Context, option ...gen.Condition) (list []*model.OpenAPIKey, err error) {
	list, err = query.OpenAPIKey.WithContext(ctx).Where(option...).Find()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

// VerifyOpenAPIKey - Verify if the OpenAPI key is valid and return the tenant ID
// Redis cache structure: key: "apikey:{api_key}", value: tenantID
// Valid for 1 hour
func VerifyOpenAPIKey(ctx context.Context, appKey string) (string, error) {
	// Get tenant ID from Redis cache
	cacheKey := "apikey:" + appKey
	tenantID, err := global.REDIS.Get(ctx, cacheKey).Result()
	if err != nil {
		// If the key does not exist in the cache, query from the database
		apiKey, err := query.OpenAPIKey.WithContext(ctx).Where(query.OpenAPIKey.APIKey.Eq(appKey), query.OpenAPIKey.Status.Eq(1)).First()
		if err != nil {
			return "", err
		}
		// Store the result in Redis cache with a 1-hour expiry
		tenantID = apiKey.TenantID
		err = global.REDIS.Set(ctx, cacheKey, tenantID, time.Hour).Err()
		if err != nil {
			logrus.Warnf("Failed to set OpenAPI key cache: %v", err)
		}
	}
	return tenantID, nil
}
