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
// @Summary Create board
// @Description Create a new board
// @Tags Board
// @Accept json
// @Produce json
// @Param create_board_req body model.CreateBoardReq true "Board information details"
// @Success 200 {object} model.CreateBoardReq "Board created successfully"
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
// @Summary Update board
// @Description Update the information of a board
// @Tags Board
// @Accept json
// @Produce json
// @Param update_board_req body model.UpdateBoardReq true "Board information details"
// @Success 200 {object} model.UpdateBoardReq "Board information updated successfully"
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
// @Summary Delete board
// @Description Delete a board by its ID
// @Tags Board
// @Accept json
// @Produce json
// @Param id path string true "Board ID"
// @Success 200 {object} model.UsersUpdateReq "Board deleted successfully"
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
// @Summary Get board list by page
// @Description Get the list of boards by page
// @Tags Board
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {object} model.UsersUpdateReq "Board list retrieved successfully"
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
// @Summary Get board by ID
// @Description Get the information of a board by its ID
// @Tags Board
// @Accept json
// @Produce json
// @Param id path string true "Board ID"
// @Success 200 {object} model.UsersUpdateReq "Board information retrieved successfully"
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
// @Summary Get board list by tenant ID
// @Description Get the list of boards for a specific tenant
// @Tags Board
// @Accept json
// @Produce json
// @Success 200 {object} model.UsersUpdateReq "Board list retrieved successfully"
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
// @Summary Get device total
// @Description Get the total number of devices
// @Tags Board
// @Accept json
// @Produce json
// @Success 200 {object} model.UsersUpdateReq "Device total retrieved successfully"
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
// @Summary Get device information
// @Description Get the information of a device
// @Tags Board
// @Accept json
// @Produce json
// @Success 200 {object} model.UsersUpdateReq "Device information retrieved successfully"
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
// @Summary Get tenant information
// @Description Get the information of a tenant
// @Tags Board
// @Accept json
// @Produce json
// @Success 200 {object} model.UsersUpdateReq "Tenant information retrieved successfully"
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
// @Summary Get user information
// @Description Get the information of a user
// @Tags Board
// @Accept json
// @Produce json
// @Success 200 {object} model.UsersUpdateReq "User information retrieved successfully"
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
// @Summary Get device information
// @Description Get the information of a device
// @Tags Board
// @Accept json
// @Produce json
// @Success 200 {object} model.UsersUpdateReq "Device information retrieved successfully"
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
// @Summary Get user information
// @Description Get the information of a user
// @Tags Board
// @Accept json
// @Produce json
// @Success 200 {object} model.UsersUpdateReq "User information retrieved successfully"
func (*BoardApi) HandleUserInfo(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)

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
// @Summary Update user information
// @Description Update the information of a user
// @Tags Board
// @Accept json
// @Produce json
// @Param users_update_req body model.UsersUpdateReq true "User information details"
// @Success 200 {object} model.UsersUpdateReq "User information updated successfully"
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
// @Summary Update user password
// @Description Update the password of a user
// @Tags Board
// @Accept json
// @Produce json
// @Param users_update_password_req body model.UsersUpdatePasswordReq true "User password details"
// @Success 200 {object} model.UsersUpdatePasswordReq "User password updated successfully"
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
// @Summary Get device trend
// @Description Get device trend data
// @Tags Board
// @Accept json
// @Produce json
// @Param device_trend_req body model.DeviceTrendReq true "Device trend request"
// @Success 200 {object} model.DeviceTrend "Device trend data"
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
