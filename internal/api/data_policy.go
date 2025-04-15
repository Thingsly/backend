package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type DataPolicyApi struct{}

// UpdateDataPolicy
// @Router   /api/v1/datapolicy [put]
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
