package api

import (
	"fmt"

	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AlarmApi struct{}

// /api/v1/alarm/config [post]
// @Summary Create alarm config
// @Description Create a new alarm configuration
// @Tags Alarm
// @Accept json
// @Produce json
// @Param alarm_config body model.CreateAlarmConfigReq true "Alarm configuration details"
// @Success 200 {object} model.AlarmConfig "Alarm configuration created successfully"
func (*AlarmApi) CreateAlarmConfig(c *gin.Context) {
	var req model.CreateAlarmConfigReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	req.TenantID = userClaims.TenantID
	data, err := service.GroupApp.Alarm.CreateAlarmConfig(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// /api/v1/alarm/config/{id} [Delete]
// @Summary Delete alarm config
// @Description Delete an existing alarm configuration
// @Tags Alarm
// @Accept json
// @Produce json
// @Param id path string true "Alarm configuration ID"
// @Success 200 {object} model.AlarmConfig "Alarm configuration deleted successfully"
func (*AlarmApi) DeleteAlarmConfig(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"err": fmt.Sprintf("id is %s", id),
		}))
		return
	}

	err := service.GroupApp.Alarm.DeleteAlarmConfig(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// /api/v1/alarm/config [PUT]
// @Summary Update alarm config
// @Description Update an existing alarm configuration
// @Tags Alarm
// @Accept json
// @Produce json
// @Param alarm_config body model.UpdateAlarmConfigReq true "Alarm configuration details"
// @Success 200 {object} model.AlarmConfig "Alarm configuration updated successfully"
func (*AlarmApi) UpdateAlarmConfig(c *gin.Context) {
	var req model.UpdateAlarmConfigReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	req.TenantID = &userClaims.TenantID
	data, err := service.GroupApp.Alarm.UpdateAlarmConfig(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// /api/v1/alarm/config [GET]
// @Summary Get alarm config list
// @Description Get a list of alarm configurations
// @Tags Alarm
// @Accept json
// @Produce json
// @Param page query int false "Page number"
func (*AlarmApi) ServeAlarmConfigListByPage(c *gin.Context) {
	var req model.GetAlarmConfigListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	req.TenantID = userClaims.TenantID

	data, err := service.GroupApp.Alarm.GetAlarmConfigListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// /api/v1/alarm/info [put]
// @Summary Update alarm info
// @Description Update an existing alarm information
// @Tags Alarm
// @Accept json
// @Produce json
// @Param alarm_info body model.UpdateAlarmInfoReq true "Alarm information details"
// @Success 200 {object} model.AlarmInfo "Alarm information updated successfully"
func (*AlarmApi) UpdateAlarmInfo(c *gin.Context) {
	var req model.UpdateAlarmInfoReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	data, err := service.GroupApp.Alarm.UpdateAlarmInfo(&req, userClaims.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)

}

// /api/v1/alarm/info/batch [put]
// @Summary Update alarm info batch
// @Description Update multiple alarm information
// @Tags Alarm
// @Accept json
// @Produce json
// @Param alarm_info body model.UpdateAlarmInfoBatchReq true "Alarm information details"
// @Success 200 {object} model.AlarmInfo "Alarm information updated successfully"
func (*AlarmApi) BatchUpdateAlarmInfo(c *gin.Context) {
	var req model.UpdateAlarmInfoBatchReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.Alarm.UpdateAlarmInfoBatch(&req, userClaims.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// /api/v1/alarm/info [get]
// @Summary Get alarm info list
// @Description Get a list of alarm information
// @Tags Alarm
// @Accept json
// @Produce json
// @Param page query int false "Page number"
func (*AlarmApi) HandleAlarmInfoListByPage(c *gin.Context) {
	var req model.GetAlarmInfoListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	req.TenantID = userClaims.TenantID

	data, err := service.GroupApp.Alarm.GetAlarmInfoListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// /api/v1/alarm/info/history [get]
// @Summary Get alarm history list
// @Description Get a list of alarm history
// @Tags Alarm
// @Accept json
// @Produce json
// @Param page query int false "Page number"
func (*AlarmApi) HandleAlarmHisttoryListByPage(c *gin.Context) {
	//
	var req model.GetAlarmHisttoryListByPage
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	data, err := service.GroupApp.Alarm.GetAlarmHisttoryListByPage(&req, userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// /api/v1/alarm/info/history [put]
// @Summary Update alarm history description
// @Description Update the description of an alarm history
// @Tags Alarm
// @Accept json
// @Produce json
// @Param alarm_history body model.AlarmHistoryDescUpdateReq true "Alarm history details"
// @Success 200 {object} model.AlarmInfo "Alarm history updated successfully"
func (*AlarmApi) AlarmHistoryDescUpdate(c *gin.Context) {
	//
	var req model.AlarmHistoryDescUpdateReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.Alarm.AlarmHistoryDescUpdate(&req, userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

func (*AlarmApi) HandleDeviceAlarmStatus(c *gin.Context) {
	var req model.GetDeviceAlarmStatusReq
	if !BindAndValidate(c, &req) {
		return
	}
	// var userClaims = c.MustGet("claims").(*utils.UserClaims)

	ok := service.GroupApp.Alarm.GetDeviceAlarmStatus(&req)
	c.Set("data", map[string]bool{
		"alarm": ok,
	})
}

// /api/v1/alarm/info/config/device [get]
// @Summary Get alarm config by device
// @Description Get alarm config by device
// @Tags Alarm
// @Accept json
// @Produce json
// @Param device_id query string true "Device ID"
// @Success 200 {object} model.AlarmConfig "Alarm config retrieved successfully"
func (*AlarmApi) HandleConfigByDevice(c *gin.Context) {
	//
	var req model.GetDeviceAlarmStatusReq
	if !BindAndValidate(c, &req) {
		return
	}
	// var userClaims = c.MustGet("claims").(*utils.UserClaims)

	list, err := service.GroupApp.Alarm.GetConfigByDevice(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// /api/v1/alarm/info/history/{id} [GET]
// @Summary Get alarm history by ID
// @Description Get alarm history by its ID
// @Tags Alarm
// @Accept json
// @Produce json
// @Param id path string true "Alarm history ID"
// @Success 200 {object} model.AlarmInfo "Alarm history retrieved successfully"
func (*AlarmApi) HandleAlarmInfoHistory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"err": fmt.Sprintf("id is %s", id),
		}))
		return
	}

	data, err := service.GroupApp.Alarm.GetAlarmInfoHistoryByID(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// GetAlarmDeviceCountsByTenant 
// @Summary 
// @Description 
// @Tags 
// @Accept json
// @Produce json
// @Success 200 {object} model.AlarmDeviceCountsResponse
// @Router /api/v1/alarm/device/counts [get]
func (api *AlarmApi) GetAlarmDeviceCountsByTenant(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	counts, err := service.GroupApp.Alarm.GetAlarmDeviceCountsByTenant(userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", counts)
}
