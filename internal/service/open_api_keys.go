// internal/service/open_api_keys.go
package service

import (
	"time"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"

	"github.com/HustIoTPlatform/backend/internal/dal"
	"github.com/HustIoTPlatform/backend/internal/model"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	"github.com/HustIoTPlatform/backend/pkg/utils"
)

type OpenAPIKey struct{}

// CreateOpenAPIKey
func (o *OpenAPIKey) CreateOpenAPIKey(req *model.CreateOpenAPIKeyReq, claims *utils.UserClaims) error {

	if claims.Authority != "SYS_ADMIN" && claims.Authority != "TENANT_ADMIN" {
		return errcode.WithVars(errcode.CodeNoPermission, map[string]interface{}{
			"required_role": "SYS_ADMIN or TENANT_ADMIN",
			"current_role":  claims.Authority,
		})
	}

	if claims.Authority == "TENANT_ADMIN" && claims.TenantID != req.TenantID {
		return errcode.WithVars(errcode.CodeNoPermission, map[string]interface{}{
			"required_tenant": req.TenantID,
			"current_tenant":  claims.TenantID,
		})
	}

	apikey, err := utils.GenerateAPIKey()
	if err != nil {
		logrus.Errorf("Failed to generate AppSecret: %v", err)
		return errcode.New(errcode.CodeSystemError)
	}

	status := int16(1)

	key := &model.OpenAPIKey{
		ID:       uuid.New(),
		TenantID: req.TenantID,
		APIKey:   apikey,
		Status:   &status,
		Name:     req.Name,
	}

	t := time.Now().UTC()
	key.CreatedAt = &t
	key.UpdatedAt = &t

	if err := dal.CreateOpenAPIKey(key); err != nil {
		logrus.Errorf("Failed to create OpenAPI key: %v", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return nil
}

// GetOpenAPIKeyList
func (o *OpenAPIKey) GetOpenAPIKeyList(req *model.OpenAPIKeyListReq, claims *utils.UserClaims) (map[string]interface{}, error) {
	var tenantID string

	if claims.Authority == "TENANT_ADMIN" || claims.Authority == "TENANT_USER" {
		tenantID = claims.TenantID
	}

	total, list, err := dal.GetOpenAPIKeyListByPage(req, tenantID)
	if err != nil {
		logrus.Errorf("Failed to query the OpenAPI key list: %v", err)
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	result := make(map[string]interface{})
	result["total"] = total
	result["list"] = list
	return result, nil
}

// UpdateOpenAPIKey
func (o *OpenAPIKey) UpdateOpenAPIKey(req *model.UpdateOpenAPIKeyReq, claims *utils.UserClaims) error {

	key, err := dal.GetOpenAPIKeyByID(req.ID)
	if err != nil {
		logrus.Errorf("Failed to retrieve OpenAPI key information: %v", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		})
	}

	if claims.Authority != "SYS_ADMIN" {
		if claims.Authority != "TENANT_ADMIN" || key.TenantID != claims.TenantID {
			return errcode.WithVars(errcode.CodeNoPermission, map[string]interface{}{
				"required_role": "SYS_ADMIN or TENANT_ADMIN",
				"current_role":  claims.Authority,
			})
		}
	}

	updates := make(map[string]interface{})
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if err := dal.UpdateOpenAPIKey(req.ID, updates); err != nil {
		logrus.Errorf("Failed to update OpenAPI key: %v", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		})
	}

	return nil
}

// DeleteOpenAPIKey
func (o *OpenAPIKey) DeleteOpenAPIKey(id string, claims *utils.UserClaims) error {

	key, err := dal.GetOpenAPIKeyByID(id)
	if err != nil {
		logrus.Errorf("Failed to retrieve OpenAPI key information: %v", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
	}

	if claims.Authority != "SYS_ADMIN" {
		if claims.Authority != "TENANT_ADMIN" || key.TenantID != claims.TenantID {
			return errcode.WithVars(errcode.CodeNoPermission, map[string]interface{}{
				"required_role": "SYS_ADMIN or TENANT_ADMIN",
				"current_role":  claims.Authority,
			})
		}
	}

	if err := dal.DeleteOpenAPIKey(id); err != nil {
		logrus.Errorf("Failed to delete OpenAPI key: %v", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
	}

	return nil
}
