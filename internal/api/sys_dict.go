package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DictApi struct{}

// CreateDictColumn
// @Summary Create dict column
// @Description Create dict column
// @Tags dict
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param dict_column body model.CreateDictReq true "Dict column"
// @Success 200 {object} model.CreateDictRes
// @Router   /api/v1/dict/column [post]
func (*DictApi) CreateDictColumn(c *gin.Context) {

	var createDictReq model.CreateDictReq
	if !BindAndValidate(c, &createDictReq) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.Dict.CreateDictColumn(&createDictReq, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// CreateDictLanguage
// @Summary Create dict language
// @Description Create dict language
// @Tags dict
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param dict_language body model.CreateDictLanguageReq true "Dict language"
// @Success 200 {object} model.CreateDictLanguageRes
// @Router   /api/v1/dict/language [post]
func (*DictApi) CreateDictLanguage(c *gin.Context) {

	var createDictLanguageReq model.CreateDictLanguageReq
	if !BindAndValidate(c, &createDictLanguageReq) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.Dict.CreateDictLanguage(&createDictLanguageReq, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// DeleteDictColumn
// @Summary Delete dict column
// @Description Delete dict column
// @Tags dict
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Dict column id"
// @Success 200 {object} model.DeleteDictRes
// @Router   /api/v1/dict/column/{id} [delete]
func (*DictApi) DeleteDictColumn(c *gin.Context) {
	id := c.Param("id")
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.Dict.DeleteDict(id, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// DeleteDictLanguage
// @Summary Delete dict language
// @Description Delete dict language
// @Tags dict
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Dict language id"
// @Success 200 {object} model.DeleteDictLanguageRes
// @Router   /api/v1/dict/language/{id} [delete]
func (*DictApi) DeleteDictLanguage(c *gin.Context) {
	id := c.Param("id")
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.Dict.DeleteDictLanguage(id, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// CreateDictColumn
// @Summary Get dict
// @Description Get dict
// @Tags dict
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param dict_enum body model.DictListReq true "Dict enum"
// @Success 200 {object} model.DictListRes
// @Router   /api/v1/dict/enum [get]
func (*DictApi) HandleDict(c *gin.Context) {
	var dictEnum model.DictListReq
	if !BindAndValidate(c, &dictEnum) {
		return
	}
	lang := c.GetHeader("Accept-Language")
	list, err := service.GroupApp.Dict.GetDict(&dictEnum, lang)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", list)
}

// Protocol service dropdown menu query API
// @Summary Protocol service dropdown menu query API
// @Description Protocol service dropdown menu query API
// @Tags dict
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.ProtocolMenuRes
// @Router   /api/v1/dict/protocol/service [get]
func (*DictApi) HandleProtocolAndService(c *gin.Context) {
	var protocolMenuReq model.ProtocolMenuReq
	if !BindAndValidate(c, &protocolMenuReq) {
		return
	}
	list, err := service.GroupApp.Dict.GetProtocolMenu(&protocolMenuReq)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// GetDictLanguage
// @Summary Get dict language
// @Description Get dict language
// @Tags dict
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Dict language id"
// @Success 200 {object} model.GetDictLanguageRes
// @Router   /api/v1/dict/language/{id} [get]
func (*DictApi) HandleDictLanguage(c *gin.Context) {
	id := c.Param("id")
	data, err := service.GroupApp.Dict.GetDictLanguageListById(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// GetDictLisyByPage
// @Summary Get dict list by page
// @Description Get dict list by page
// @Tags dict
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetDictLisyByPageRes
// @Router   /api/v1/dict [get]
func (*DictApi) HandleDictLisyByPage(c *gin.Context) {
	var byList model.GetDictLisyByPageReq
	if !BindAndValidate(c, &byList) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	logrus.Info("byList", byList)
	list, err := service.GroupApp.Dict.GetDictListByPage(&byList, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", list)
}
