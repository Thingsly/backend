package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type DataScriptApi struct{}

// CreateDataScript
// @Router   /api/v1/data_script [post]
// @Summary Create data script
// @Description Create a new data script
// @Tags Data Script
// @Accept json
// @Produce json
// @Param create_data_script_req body model.CreateDataScriptReq true "Data script details"
// @Success 200 {object} model.CreateDataScriptReq "Data script created successfully"
func (*DataScriptApi) CreateDataScript(c *gin.Context) {
	var req model.CreateDataScriptReq
	if !BindAndValidate(c, &req) {
		return
	}
	data, err := service.GroupApp.DataScript.CreateDataScript(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// UpdateDataScript
// @Router   /api/v1/data_script [put]
// @Summary Update data script
// @Description Update the data script
// @Tags Data Script
// @Accept json
// @Produce json
// @Param update_data_script_req body model.UpdateDataScriptReq true "Data script details"
// @Success 200 {object} model.UpdateDataScriptReq "Data script updated successfully"
func (*DataScriptApi) UpdateDataScript(c *gin.Context) {
	var req model.UpdateDataScriptReq
	if !BindAndValidate(c, &req) {
		return
	}

	if req.Description == nil && req.Name == "" {
		c.Error(errcode.WithData(errcode.CodeParamError, "description and name can not be empty at the same time"))
		return
	}

	err := service.GroupApp.DataScript.UpdateDataScript(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// DeleteDataScript
// @Router   /api/v1/data_script/{id} [delete]
// @Summary Delete data script
// @Description Delete the data script
// @Tags Data Script
// @Accept json
// @Produce json
// @Param id path string true "Data script ID"
// @Success 200 {object} model.DataScript "Data script deleted successfully"
func (*DataScriptApi) DeleteDataScript(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.DataScript.DeleteDataScript(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GetDataScriptListByPage - Data Script List Pagination
// @Router   /api/v1/data_script [get]
// @Summary Get data script list by page
// @Description Get the list of data scripts by page
// @Tags Data Script
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {object} model.GetDataScriptListByPageReq "Data script list retrieved successfully"
func (*DataScriptApi) HandleDataScriptListByPage(c *gin.Context) {
	var req model.GetDataScriptListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	data_scriptList, err := service.GroupApp.DataScript.GetDataScriptListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data_scriptList)
}

// /api/v1/data_script/quiz
// @Summary Quiz data script
// @Description Quiz the data script
// @Tags Data Script
// @Accept json
// @Produce json
// @Param quiz_data_script_req body model.QuizDataScriptReq true "Data script details"
// @Success 200 {object} model.QuizDataScriptReq "Data script quiz retrieved successfully"
func (*DataScriptApi) QuizDataScript(c *gin.Context) {
	var req model.QuizDataScriptReq
	if !BindAndValidate(c, &req) {
		return
	}

	data, err := service.GroupApp.DataScript.QuizDataScript(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// /api/v1/data_script/enable [put]
// @Summary Enable data script
// @Description Enable the data script
// @Tags Data Script
// @Accept json
// @Produce json
// @Param enable_data_script_req body model.EnableDataScriptReq true "Data script details"
// @Success 200 {object} model.EnableDataScriptReq "Data script enabled successfully"
func (*DataScriptApi) EnableDataScript(c *gin.Context) {
	var req model.EnableDataScriptReq
	if !BindAndValidate(c, &req) {
		return
	}

	err := service.GroupApp.DataScript.EnableDataScript(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}
