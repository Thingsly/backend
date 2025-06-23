package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type DeviceApi struct{}

// @Summary Create device
// @Description Create a new device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device body model.CreateDeviceReq true "Device information"
// @Param device.name body string true "Device name"
// @Param device.label body string false "Device label"
// @Param device.device_config_id body string false "Device config ID"
// @Param device.access_way body string true "Access way (A for direct access)"
// @Success 200 {object} model.CreateDeviceRes
// @Success 400 {object} model.ErrorResponse
// @Success 401 {object} model.ErrorResponse
// @Success 403 {object} model.ErrorResponse
// @Success 500 {object} model.ErrorResponse
// @Router /api/v1/device [post]
// @Example request - Create device
//
//	{
//	    "name": "test-create-device-dodat",
//	    "label": "",
//	    "device_config_id": "",
//	    "access_way": "A"
//	}
//
// @Example response - 200
//
//	{
//	    "code": 200,
//	    "message": "Operation successful",
//	    "data": {
//	        "id": "4ee2d599-cbac-ca9a-0491-f3a25b3057b3",
//	        "name": "test-create-device-dodat",
//	        "voucher": "{\"username\":\"fe0533d9-6fc9-dbe4-414\",\"password\":\"ae84399\"}",
//	        "tenant_id": "d616bcbb",
//	        "is_enabled": "",
//	        "activate_flag": "active",
//	        "created_at": "2025-04-27T11:12:35.972366502Z",
//	        "update_at": "2025-04-27T11:12:35.972366502Z",
//	        "device_number": "4ee2d599-cbac-ca9a-0491-f3a25b3057b3",
//	        "product_id": null,
//	        "parent_id": null,
//	        "protocol": null,
//	        "label": "",
//	        "location": null,
//	        "sub_device_addr": null,
//	        "current_version": null,
//	        "additional_info": "{}",
//	        "protocol_config": "{}",
//	        "remark1": null,
//	        "remark2": null,
//	        "remark3": null,
//	        "device_config_id": null,
//	        "batch_number": null,
//	        "activate_at": null,
//	        "is_online": 0,
//	        "access_way": "A",
//	        "description": null,
//	        "service_access_id": null
//	    }
//	}
//
// @Example response - 400
//
//	{
//	    "code": 400,
//	    "message": "Invalid request parameters",
//	    "data": null
//	}
//
// @Example response - 401
//
//	{
//	    "code": 401,
//	    "message": "Unauthorized",
//	    "data": null
//	}
//
// @Example response - 403
//
//	{
//	    "code": 403,
//	    "message": "Forbidden",
//	    "data": null
//	}
//
// @Example response - 500
//
//	{
//	    "code": 500,
//	    "message": "Internal server error",
//	    "data": null
//	}
func (*DeviceApi) CreateDevice(c *gin.Context) {
	var req model.CreateDeviceReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.Device.CreateDevice(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// @Summary Batch create devices
// @Description Create multiple devices at service access points
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param devices body model.BatchCreateDeviceReq true "List of devices to create"
// @Success 200 {object} model.BatchCreateDeviceRes
// @Router /api/v1/device/service/access/batch [post]
func (*DeviceApi) CreateDeviceBatch(c *gin.Context) {
	var req model.BatchCreateDeviceReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.Device.CreateDeviceBatch(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Delete device
// @Description Delete a device by ID
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Device ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/{id} [delete]
func (*DeviceApi) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.Device.DeleteDevice(id, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Update device
// @Description Update an existing device's information
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device body model.UpdateDeviceReq true "Device information"
// @Success 200 {object} model.UpdateDeviceRes
// @Router /api/v1/device [put]
func (*DeviceApi) UpdateDevice(c *gin.Context) {
	var req model.UpdateDeviceReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.Device.UpdateDevice(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// @Summary Activate device
// @Description Activate a device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device body model.ActiveDeviceReq true "Device activation information"
// @Success 200 {object} model.Device
// @Router /api/v1/device/active [put]
func (*DeviceApi) ActiveDevice(c *gin.Context) {
	var req model.ActiveDeviceReq
	if !BindAndValidate(c, &req) {
		return
	}

	device, err := service.GroupApp.Device.ActiveDevice(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", device)
}

// @Summary Get device by ID
// @Description Get detailed information about a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Device ID"
// @Success 200 {object} model.Device
// @Router /api/v1/device/detail/{id} [get]
func (*DeviceApi) HandleDeviceByID(c *gin.Context) {
	id := c.Param("id")
	device, err := service.GroupApp.Device.GetDeviceByIDV1(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", device)
}

// @Summary Get device list by page
// @Description Retrieve a paginated list of devices with optional filtering
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} model.DeviceListResponse
// @Router /api/v1/device [get]
func (*DeviceApi) HandleDeviceListByPage(c *gin.Context) {
	var req model.GetDeviceListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	list, err := service.GroupApp.Device.GetDeviceListByPage(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// @Summary Check device number availability
// @Description Check if a device number is available
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param deviceNumber path string true "Device number to check"
// @Success 200 {object} map[string]bool
// @Router /api/v1/device/check/{deviceNumber} [get]
func (*DeviceApi) CheckDeviceNumber(c *gin.Context) {
	deviceNumber := c.Param("deviceNumber")
	ok, _ := service.GroupApp.Device.CheckDeviceNumber(deviceNumber)
	data := map[string]interface{}{"is_available": ok}
	c.Set("data", data)
}

// @Summary Create device template
// @Description Create a new device template
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param template body model.CreateDeviceTemplateReq true "Device template information"
// @Success 200 {object} model.CreateDeviceTemplateRes
// @Router /api/v1/device/template [post]
func (*DeviceApi) CreateDeviceTemplate(c *gin.Context) {
	var req model.CreateDeviceTemplateReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.DeviceTemplate.CreateDeviceTemplate(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// @Summary Update device template
// @Description Update an existing device template
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param template body model.UpdateDeviceTemplateReq true "Device template information"
// @Success 200 {object} model.UpdateDeviceTemplateRes
// @Router /api/v1/device/template [put]
func (*DeviceApi) UpdateDeviceTemplate(c *gin.Context) {
	var req model.UpdateDeviceTemplateReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.DeviceTemplate.UpdateDeviceTemplate(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// @Summary Get device template list
// @Description Get a paginated list of device templates
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} model.GetDeviceTemplateListRes
// @Router /api/v1/device/template [get]
func (*DeviceApi) HandleDeviceTemplateListByPage(c *gin.Context) {
	var req model.GetDeviceTemplateListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.DeviceTemplate.GetDeviceTemplateListByPage(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	serilizedData, err := utils.SerializeData(data, GetDeviceTemplateListData{})
	if err != nil {
		c.Error(errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"error": err.Error(),
		}))
		return
	}

	c.Set("data", serilizedData)
}

// @Summary Get device template menu
// @Description Get device template menu structure
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.GetDeviceTemplateMenuRes
// @Router /api/v1/device/template/menu [get]
func (*DeviceApi) HandleDeviceTemplateMenu(c *gin.Context) {
	var req model.GetDeviceTemplateMenuReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.DeviceTemplate.GetDeviceTemplateMenu(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// @Summary Delete device template
// @Description Delete a device template by ID
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Template ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/template/{id} [delete]
func (*DeviceApi) DeleteDeviceTemplate(c *gin.Context) {
	id := c.Param("id")
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.DeviceTemplate.DeleteDeviceTemplate(id, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Get device template by ID
// @Description Get detailed information about a specific device template
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Template ID"
// @Success 200 {object} model.DeviceTemplateReadSchema
// @Router /api/v1/device/template/detail/{id} [get]
func (*DeviceApi) HandleDeviceTemplateById(c *gin.Context) {
	id := c.Param("id")
	data, err := service.GroupApp.DeviceTemplate.GetDeviceTemplateById(id)
	if err != nil {
		c.Error(err)
		return
	}
	serilizedData, err := utils.SerializeData(data, DeviceTemplateReadSchema{})
	if err != nil {
		c.Error(errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"error": err.Error(),
		}))
		return
	}
	c.Set("data", serilizedData)
}

// @Summary Get device template by device ID
// @Description Get device template details for a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device_id query string true "Device ID"
// @Success 200 {object} model.DeviceTemplateChartRes
// @Router /api/v1/device/template/chart [get]
func (*DeviceApi) HandleDeviceTemplateByDeviceId(c *gin.Context) {
	deviceId := c.Query("device_id")
	if deviceId == "" {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"device_id": deviceId,
			"msg":       "device_id is required",
		}))
		return
	}
	data, err := service.GroupApp.DeviceTemplate.GetDeviceTemplateByDeviceId(deviceId)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Create device group
// @Description Create a new device group
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param group body model.CreateDeviceGroupReq true "Device group information"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/group [post]
func (*DeviceApi) CreateDeviceGroup(c *gin.Context) {
	var req model.CreateDeviceGroupReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.DeviceGroup.CreateDeviceGroup(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Delete device group
// @Description Delete a device group by ID
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Group ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/group/{id} [delete]
func (*DeviceApi) DeleteDeviceGroup(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.DeviceGroup.DeleteDeviceGroup(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Update device group
// @Description Update an existing device group
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param group body model.UpdateDeviceGroupReq true "Device group information"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/group [put]
func (*DeviceApi) UpdateDeviceGroup(c *gin.Context) {
	var req model.UpdateDeviceGroupReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.DeviceGroup.UpdateDeviceGroup(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Get device group list
// @Description Get a paginated list of device groups
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} model.GetDeviceGroupsListRes
// @Router /api/v1/device/group [get]
func (*DeviceApi) HandleDeviceGroupByPage(c *gin.Context) {
	var req model.GetDeviceGroupsListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.DeviceGroup.GetDeviceGroupListByPage(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Get device group tree
// @Description Get device group tree structure
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.GetDeviceGroupTreeRes
// @Router /api/v1/device/group/tree [get]
func (*DeviceApi) HandleDeviceGroupByTree(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.DeviceGroup.GetDeviceGroupByTree(userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Get device group detail
// @Description Get detailed information about a specific device group
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Group ID"
// @Success 200 {object} model.GetDeviceGroupDetailRes
// @Router /api/v1/device/group/detail/{id} [get]
func (*DeviceApi) HandleDeviceGroupByDetail(c *gin.Context) {
	id := c.Param("id")
	data, err := service.GroupApp.DeviceGroup.GetDeviceGroupDetail(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Create device group relation
// @Description Create a relation between device and group
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param relation body model.CreateDeviceGroupRelationReq true "Device group relation information"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/group/relation [post]
func (*DeviceApi) CreateDeviceGroupRelation(c *gin.Context) {
	var req model.CreateDeviceGroupRelationReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.DeviceGroup.CreateDeviceGroupRelation(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Delete device group relation
// @Description Delete a relation between device and group
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param group_id query string true "Group ID"
// @Param device_id query string true "Device ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/group/relation [delete]
func (*DeviceApi) DeleteDeviceGroupRelation(c *gin.Context) {
	var req model.DeleteDeviceGroupRelationReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.DeviceGroup.DeleteDeviceGroupRelation(req.GroupId, req.DeviceId)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Get device group relation list
// @Description Get list of devices in a group
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param group_id query string true "Group ID"
// @Success 200 {object} model.GetDeviceListByGroupRes
// @Router /api/v1/device/group/relation/list [get]
func (*DeviceApi) HandleDeviceGroupRelation(c *gin.Context) {
	var req model.GetDeviceListByGroup
	if !BindAndValidate(c, &req) {
		return
	}
	data, err := service.GroupApp.DeviceGroup.GetDeviceGroupRelation(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Get device group list by device ID
// @Description Get list of groups a device belongs to
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device_id query string true "Device ID"
// @Success 200 {object} model.GetDeviceGroupListByDeviceIdRes
// @Router /api/v1/device/group/relation [get]
func (*DeviceApi) HandleDeviceGroupListByDeviceId(c *gin.Context) {
	var req model.GetDeviceGroupListByDeviceIdReq
	if !BindAndValidate(c, &req) {
		return
	}
	data, err := service.GroupApp.DeviceGroup.GetDeviceGroupByDeviceId(req.DeviceId)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Remove sub-device
// @Description Remove a sub-device from its parent device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param sub_device_id query string true "Sub-device ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/sub-remove [post]
func (*DeviceApi) RemoveSubDevice(c *gin.Context) {
	var req model.RemoveSonDeviceReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.Device.RemoveSubDevice(req.SubDeviceId, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Get tenant device list
// @Description Get list of devices for a specific tenant
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.DeviceListResponse
// @Router /api/v1/device/tenant/list [get]
func (*DeviceApi) HandleTenantDeviceList(c *gin.Context) {
	var req model.GetDeviceMenuReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.Device.GetTenantDeviceList(&req, userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Get device list
// @Description Get list of all devices
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.DeviceListResponse
// @Router /api/v1/device/list [get]
func (*DeviceApi) HandleDeviceList(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.Device.GetDeviceList(c, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Create sub device
// @Description Create a new sub-device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device body model.CreateSonDeviceRes true "Sub-device information"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/son/add [post]
func (*DeviceApi) CreateSonDevice(c *gin.Context) {
	var param model.CreateSonDeviceRes
	if !BindAndValidate(c, &param) {
		return
	}

	err := service.GroupApp.Device.CreateSonDevice(c, &param)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Get device connect form
// @Description Get device connection form data
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device_id query string true "Device ID"
// @Success 200 {object} model.DeviceConnectFormRes
// @Router /api/v1/device/connect/form [get]
func (*DeviceApi) DeviceConnectForm(c *gin.Context) {
	var param model.DeviceConnectFormReq
	if !BindAndValidate(c, &param) {
		return
	}
	list, err := service.GroupApp.Device.DeviceConnectForm(c, &param)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// @Summary Get device connect info
// @Description Get device connection information
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device_id query string true "Device ID"
// @Success 200 {object} model.DeviceConnectInfoRes
// @Router /api/v1/device/connect/info [get]
func (*DeviceApi) DeviceConnect(c *gin.Context) {
	var param model.DeviceConnectFormReq
	if !BindAndValidate(c, &param) {
		return
	}
	// Get language settings
	lang := c.Request.Header.Get("Accept-Language")
	if lang == "" {
		lang = "vi_VN"
	}

	list, err := service.GroupApp.Device.DeviceConnect(c, &param, lang)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// @Summary Update device voucher
// @Description Update device authentication voucher
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param voucher body model.UpdateDeviceVoucherReq true "Device voucher information"
// @Success 200 {object} model.UpdateDeviceVoucherRes
// @Router /api/v1/device/update/voucher [post]
func (*DeviceApi) UpdateDeviceVoucher(c *gin.Context) {
	var param model.UpdateDeviceVoucherReq
	if !BindAndValidate(c, &param) {
		return
	}
	voucher, err := service.GroupApp.Device.UpdateDeviceVoucher(c, &param)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", voucher)
}

// @Summary Get sub-list
// @Description Get a list of sub-devices
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Parent device ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} model.SubDeviceListResponse
// @Router /api/v1/device/sub-list/{id} [get]
func (*DeviceApi) HandleSubList(c *gin.Context) {
	var req model.PageReq
	parent_id := c.Param("id")
	if parent_id == "" {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"msg": "no parent_id",
		}))
		return
	}
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	list, total, err := service.GroupApp.Device.GetSubList(c, parent_id, int64(req.Page), int64(req.PageSize), userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", map[string]interface{}{
		"total": total,
		"list":  list,
	})
}

// @Summary Get device metrics
// @Description Get metrics for a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/device/metrics/{id} [get]
func (*DeviceApi) HandleMetrics(c *gin.Context) {
	id := c.Param("id")
	list, err := service.GroupApp.Device.GetMetrics(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// @Summary Get device action menu
// @Description Get action menu for a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device_id query string true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/device/metrics/menu [get]
func (*DeviceApi) HandleActionByDeviceID(c *gin.Context) {
	var param model.GetActionByDeviceIDReq
	if !BindAndValidate(c, &param) {
		return
	}
	list, err := service.GroupApp.Device.GetActionByDeviceID(param.DeviceID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// @Summary Get device condition menu
// @Description Get condition menu for a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device_id query string true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/device/metrics/condition/menu [get]
func (*DeviceApi) HandleConditionByDeviceID(c *gin.Context) {
	var param model.GetActionByDeviceIDReq
	if !BindAndValidate(c, &param) {
		return
	}
	list, err := service.GroupApp.Device.GetConditionByDeviceID(param.DeviceID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// @Summary Get device telemetry map
// @Description Get telemetry map for a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/device/map/telemetry/{id} [get]
func (*DeviceApi) HandleMapTelemetry(c *gin.Context) {
	id := c.Param("id")
	data, err := service.GroupApp.Device.GetMapTelemetry(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Get device template chart select
// @Description Get device template chart selection options
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.DeviceTemplateChartSelectRes
// @Router /api/v1/device/template/chart/select [get]
func (*DeviceApi) HandleDeviceTemplateChartSelect(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	list, err := service.GroupApp.Device.GetDeviceTemplateChartSelect(userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// @Summary Update device configuration
// @Description Update configuration for a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param config body model.ChangeDeviceConfigReq true "Device configuration"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/update/config [put]
func (*DeviceApi) UpdateDeviceConfig(c *gin.Context) {
	var param model.ChangeDeviceConfigReq
	if !BindAndValidate(c, &param) {
		return
	}
	err := service.GroupApp.Device.UpdateDeviceConfig(&param)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// @Summary Get device online status
// @Description Get the online status of a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/device/online/status/{id} [get]
func (*DeviceApi) HandleDeviceOnlineStatus(c *gin.Context) {
	id := c.Param("id")
	data, err := service.GroupApp.Device.GetDeviceOnlineStatus(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Gateway register
// @Description Register a gateway device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param gateway body model.GatewayRegisterReq true "Gateway registration information"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/gateway/register [post]
func (*DeviceApi) GatewayRegister(c *gin.Context) {
	var req model.GatewayRegisterReq
	if !BindAndValidate(c, &req) {
		return
	}
	data, err := service.GroupApp.Device.GatewayRegister(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// @Summary Gateway sub register
// @Description Register a sub-device of a gateway
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param gateway body model.GatewaySubRegisterReq true "Gateway sub-device registration information"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/device/gateway/sub/register [post]
func (*DeviceApi) GatewaySubRegister(c *gin.Context) {
	var req model.DeviceRegisterReq
	if !BindAndValidate(c, &req) {
		logrus.Warningf("GatewaySubRegister:%#v", req)
		return
	}
	logrus.Warningf("GatewaySubRegister:%#v", req)
	data, err := service.GroupApp.Device.GatewayDeviceRegister(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// @Summary Get device metrics chart
// @Description Get metrics chart data for a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device_id query string true "Device ID"
// @Param metric_key query string true "Metric key"
// @Param start_time query string false "Start time"
// @Param end_time query string false "End time"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/device/metrics/chart [get]
func (*DeviceApi) HandleDeviceMetricsChart(c *gin.Context) {
	var param model.GetDeviceMetricsChartReq
	if !BindAndValidate(c, &param) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	data, err := service.GroupApp.Device.GetDeviceMetricsChart(&param, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// @Summary Get device selector
// @Description Get device selector options
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param device_id query string false "Device ID"
// @Param device_type query string false "Device type"
// @Success 200 {object} model.DeviceSelectorRes
// @Router /api/v1/device/selector [get]
func (*DeviceApi) HandleDeviceSelector(c *gin.Context) {
	var req model.DeviceSelectorReq
	if !BindAndValidate(c, &req) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)
	list, err := service.GroupApp.Device.GetDeviceSelector(req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// @Summary Get latest telemetry data
// @Description Get the latest telemetry data for a specific device
// @Tags device
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/device/telemetry/latest [get]
func (*DeviceApi) HandleTenantTelemetryData(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	data, err := service.GroupApp.Device.GetTenantTelemetryData(userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}