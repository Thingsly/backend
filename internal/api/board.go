package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BoardApi struct{}

// CreateBoard
// @Router   /api/v1/board [post]
func (*BoardApi) CreateBoard(c *gin.Context) {
	var req model.CreateBoardReq
	if !BindAndValidate(c, &req) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)
	req.TenantID = userClaims.TenantID

	boardInfo, err := service.GroupApp.Board.CreateBoard(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", boardInfo)
}

// UpdateBoard
// @Router   /api/v1/board [put]
func (*BoardApi) UpdateBoard(c *gin.Context) {
	var req model.UpdateBoardReq
	if !BindAndValidate(c, &req) {
		return
	}

	// if req.Description == nil && req.Name == "" && req.HomeFlag == "" {
	// 	c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Content to be modified cannot be empty})
	// 	return
	// }

	userClaims := c.MustGet("claims").(*utils.UserClaims)
	req.TenantID = userClaims.TenantID

	d, err := service.GroupApp.Board.UpdateBoard(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", d)
}

// DeleteBoard
// @Router   /api/v1/board/{id} [delete]
func (*BoardApi) DeleteBoard(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.Board.DeleteBoard(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GetBoardListByPage
// @Router   /api/v1/board [get]
func (*BoardApi) HandleBoardListByPage(c *gin.Context) {
	var req model.GetBoardListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	boardList, err := service.GroupApp.Board.GetBoardListByPage(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", boardList)
}

// GetBoard
// @Router   /api/v1/board/{id} [get]
func (*BoardApi) HandleBoard(c *gin.Context) {
	id := c.Param("id")
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	board, err := service.GroupApp.Board.GetBoard(id, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", board)
}

// GetBoardListByTenantId
// @Router   /api/v1/board/home [get]
func (*BoardApi) HandleBoardListByTenantId(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	boardList, err := service.GroupApp.Board.GetBoardListByTenantId(userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", boardList)
}

// GetDeviceTotal
// @Router   /api/v1/board/device/total [get]
func (*BoardApi) HandleDeviceTotal(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	board := service.GroupApp.Board
	total, err := board.GetDeviceTotal(c, userClaims.Authority, userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", total)
}

// GetDevice
// @Router   /api/v1/board/device [get]
func (*BoardApi) HandleDevice(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	board := service.GroupApp.Board
	data, err := board.GetDevice(c, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// GetTenant
// @Router   /api/v1/board/tenant [get]
func (*BoardApi) HandleTenant(c *gin.Context) {
	// TODO:: Unsure if user information needs to be revalidated
	users := service.UsersService{}
	data, err := users.GetTenant(c)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// GetTenantUserInfo
// @Router   /api/v1/board/tenant/user/info [get]
func (*BoardApi) HandleTenantUserInfo(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	tenantID := userClaims.TenantID
	// Query tenant information by tenant ID
	tenantInfo, err := service.GroupApp.User.GetTenantInfo(tenantID)
	if err != nil {
		c.Error(err)
		return
	}
	users := service.UsersService{}
	data, err := users.GetTenantUserInfo(c, tenantInfo.Email)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// GetTenantDeviceInfo
// @Router   /api/v1/board/tenant/device/info [get]
func (*BoardApi) HandleTenantDeviceInfo(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	board := service.GroupApp.Board
	total, err := board.GetDeviceByTenantID(c, userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", total)
}

// GetUserInfo
// @Router   /api/v1/board/user/info [get]
func (*BoardApi) HandleUserInfo(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	// 根据租户ID查询租户信息
	users := service.UsersService{}
	data, err := users.GetTenantInfo(c, userClaims.Email)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// UpdateUserInfo
// @Router   /api/v1/board/user/update [post]
func (*BoardApi) UpdateUserInfo(c *gin.Context) {
	var param model.UsersUpdateReq
	if !BindAndValidate(c, &param) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	users := service.UsersService{}
	err := users.UpdateTenantInfo(c, userClaims, &param)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// UpdateUserInfoPassword
// @Router   /api/v1/board/user/update/password [post]
func (*BoardApi) UpdateUserInfoPassword(c *gin.Context) {
	var param model.UsersUpdatePasswordReq
	if !BindAndValidate(c, &param) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	users := service.UsersService{}
	err := users.UpdateTenantInfoPassword(c, userClaims, &param)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GetDeviceTrend
// @Router   /api/v1/board/trend [get]
func (*BoardApi) GetDeviceTrend(c *gin.Context) {
	var deviceTrendReq model.DeviceTrendReq
	if !BindAndValidate(c, &deviceTrendReq) {
		return
	}

	// Get user claims
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	// If tenantID is not specified in the request, use the current user's tenantID
	if deviceTrendReq.TenantID == nil || *deviceTrendReq.TenantID == "" {
		deviceTrendReq.TenantID = &userClaims.TenantID
	}

	// Permission check - Only system administrators can view data from other tenants
	if *deviceTrendReq.TenantID != userClaims.TenantID && userClaims.Authority != "SYS_ADMIN" {
		c.Error(errcode.New(errcode.CodeNoPermission))
		return
	}

	// Call the service layer to fetch trend data
	trend, err := service.GroupApp.Device.GetDeviceTrend(c, *deviceTrendReq.TenantID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", trend)
}
