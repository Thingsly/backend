package api

import (
	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SceneApi struct{}

// Create scene
// @Summary Create scene
// @Description Create scene
// @Tags scene
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param scene body model.CreateSceneReq true "Scene"
// @Success 200 {object} model.CreateSceneRes
// @Router   /api/v1/scene [post]
func (*SceneApi) CreateScene(c *gin.Context) {
	var req model.CreateSceneReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	id, err := service.GroupApp.Scene.CreateScene(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", map[string]interface{}{"scene_id": id})
}

// Delete scene
// @Summary Delete scene
// @Description Delete scene
// @Tags scene
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Scene id"
// @Success 200 {object} model.DeleteSceneRes
// @Router   /api/v1/scene [delete]
func (*SceneApi) DeleteScene(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.Scene.DeleteScene(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// Update scene
// @Summary Update scene
// @Description Update scene
// @Tags scene
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param scene body model.UpdateSceneReq true "Scene"
// @Success 200 {object} model.UpdateSceneRes
// @Router   /api/v1/scene [put]
func (*SceneApi) UpdateScene(c *gin.Context) {
	var req model.UpdateSceneReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	id, err := service.GroupApp.Scene.UpdateScene(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", map[string]interface{}{"scene_id": id})
}

// Get scene by id
// @Summary Get scene by id
// @Description Get scene by id
// @Tags scene
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Scene id"
// @Success 200 {object} model.GetSceneRes
// @Router   /api/v1/scene/detail/{id} [get]
func (*SceneApi) HandleScene(c *gin.Context) {
	id := c.Param("id")
	data, err := service.GroupApp.Scene.GetScene(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// Get scene by page
// @Summary Get scene by page
// @Description Get scene by page
// @Tags scene
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetSceneListByPageRes
// @Router   /api/v1/scene [get]
func (*SceneApi) HandleSceneByPage(c *gin.Context) {
	var req model.GetSceneListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.Scene.GetSceneListByPage(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// Active scene
// @Summary Active scene
// @Description Active scene
// @Tags scene
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Scene id"
// @Success 200 {object} model.ActiveSceneRes
// @Router   /api/v1/scene/active/{id} [post]
// todo: Incomplete
func (*SceneApi) ActiveScene(c *gin.Context) {
	id := c.Param("id")

	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.Scene.ActiveScene(id, userClaims.ID, userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// Get scene log by page
// @Summary Get scene log by page
// @Description Get scene log by page
// @Tags scene
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetSceneLogListByPageRes
// @Router   /api/v1/scene/log [get]
func (*SceneApi) HandleSceneLog(c *gin.Context) {
	var req model.GetSceneLogListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	data, err := service.GroupApp.Scene.GetSceneLog(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}
