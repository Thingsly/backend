// internal/api/open_api_keys.go
package api

import (
	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type OpenAPIKeyApi struct{}

// CreateOpenAPIKey
// @Summary Create open api key
// @Description Create open api key
// @Tags open_api_keys
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param open_api_key body model.CreateOpenAPIKeyReq true "Open api key"
// @Success 200 {object} model.CreateOpenAPIKeyRes
// @Router /api/v1/open/keys [post]
func (*OpenAPIKeyApi) CreateOpenAPIKey(c *gin.Context) {
	var createReq model.CreateOpenAPIKeyReq
	if !BindAndValidate(c, &createReq) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.OpenAPIKey.CreateOpenAPIKey(&createReq, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// GetOpenAPIKeyList
// @Summary Get open api key list
// @Description Get open api key list
// @Tags open_api_keys
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param open_api_key body model.OpenAPIKeyListReq true "Open api key"
// @Success 200 {object} model.OpenAPIKeyListRes
// @Router /api/v1/open/keys [get]
func (*OpenAPIKeyApi) GetOpenAPIKeyList(c *gin.Context) {
	var listReq model.OpenAPIKeyListReq
	if !BindAndValidate(c, &listReq) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)

	list, err := service.GroupApp.OpenAPIKey.GetOpenAPIKeyList(&listReq, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", list)
}

// UpdateOpenAPIKey
// @Summary Update open api key
// @Description Update open api key
// @Tags open_api_keys
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param open_api_key body model.UpdateOpenAPIKeyReq true "Open api key"
// @Success 200 {object} model.UpdateOpenAPIKeyRes
// @Router /api/v1/open/keys [put]
func (*OpenAPIKeyApi) UpdateOpenAPIKey(c *gin.Context) {
	var updateReq model.UpdateOpenAPIKeyReq
	if !BindAndValidate(c, &updateReq) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.OpenAPIKey.UpdateOpenAPIKey(&updateReq, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// DeleteOpenAPIKey
// @Summary Delete open api key
// @Description Delete open api key
// @Tags open_api_keys
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Open api key id"
// @Success 200 {object} model.DeleteOpenAPIKeyRes
// @Router /api/v1/open/keys/{id} [delete]
func (*OpenAPIKeyApi) DeleteOpenAPIKey(c *gin.Context) {
	id := c.Param("id")

	var userClaims = c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.OpenAPIKey.DeleteOpenAPIKey(id, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}
