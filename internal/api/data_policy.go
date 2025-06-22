package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type DataPolicyApi struct{}

// UpdateDataPolicy
// @Router   /api/v1/datapolicy [put]
// @Summary Update data policy
// @Description Update the data policy
// @Tags Data Policy
// @Accept json
// @Produce json
// @Param update_data_policy_req body model.UpdateDataPolicyReq true "Data policy details"
// @Success 200 {object} model.UpdateDataPolicyReq "Data policy updated successfully"
func (*DataPolicyApi) UpdateDataPolicy(c *gin.Context) {
	var req model.UpdateDataPolicyReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.DataPolicy.UpdateDataPolicy(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// GetDataPolicyListByPage
// @Router   /api/v1/datapolicy [get]
// @Summary Get data policy list by page
// @Description Get the list of data policies by page
// @Tags Data Policy
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {object} model.GetDataPolicyListByPageReq "Data policy list retrieved successfully"
func (*DataPolicyApi) HandleDataPolicyListByPage(c *gin.Context) {
	var req model.GetDataPolicyListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	datapolicyList, err := service.GroupApp.DataPolicy.GetDataPolicyListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", datapolicyList)
}
