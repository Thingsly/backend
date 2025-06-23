package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserApi struct{}

// Login
// @Summary      User Login
// @Description  An authentication token (Token) will be generated and returned upon successful login. The client must include this token in the 'x-token' field of the HTTP request header for all subsequent API requests that require authentication. The server will validate this token to confirm the user's identity and authorize access to protected resources.
// @Tags         User Authentication
// @Accept       json
// @Produce      json
// @Param        request body model.LoginReq true "Login credentials" example({"email":"user@example.com","password":"yourpassword"})
// @Success      200 {object} model.LoginRsp "Success"
// @Failure      400 {object} errcode.Error "Error response"
// @Router       /api/v1/login [post]
// @example request - "Request example" {"email":"test@thingsly.vn","password":"123456"}
func (*UserApi) Login(c *gin.Context) {
	var loginReq model.LoginReq
	if !BindAndValidate(c, &loginReq) {
		return
	}

	result := utils.ValidateInput(loginReq.Email)
	if !result.IsValid {
		c.Error(errcode.WithData(200013, map[string]interface{}{
			"error": result.Message,
		}))
		return
	}

	if result.Type == utils.Phone {
		// Retrieve user email by phone number
		email, err := service.GroupApp.User.GetUserEmailByPhoneNumber(loginReq.Email)
		if err != nil {
			c.Error(err)
			return
		}
		loginReq.Email = email
	}

	loginLock := service.NewLoginLock()

	// Check if the user is allowed to log in
	if loginLock.MaxFailedAttempts > 0 {
		if err := loginLock.GetAllowLogin(c, loginReq.Email); err != nil {
			c.Error(err)
			return
		}
	}

	loginRsp, err := service.GroupApp.User.Login(c, &loginReq)
	if err != nil {
		_ = loginLock.LoginFail(c, loginReq.Email)
		c.Error(err)
		return
	}
	_ = loginLock.LoginSuccess(c, loginReq.Email)
	c.Set("data", loginRsp)
}

// Logout
// @Summary Logout
// @Description Logout
// @Tags user
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.LogoutRes
// @Router   /api/v1/user/logout [get]
func (*UserApi) Logout(c *gin.Context) {
	token := c.GetHeader("x-token")
	err := service.GroupApp.User.Logout(token)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// Refresh token
// @Summary Refresh token
// @Description Refresh token
// @Tags user
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.RefreshTokenRes
// @Router   /api/v1/user/refresh [get]
func (*UserApi) RefreshToken(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	loginRsp, err := service.GroupApp.User.RefreshToken(userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", loginRsp)
}

// Get verification code
// @Summary Get verification code
// @Description Get verification code
// @Tags user
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Param is_register query string true "Is register"
// @Success 200 {object} model.GetVerificationCodeRes
// @Router   /api/v1/verification/code [get]
func (*UserApi) HandleVerificationCode(c *gin.Context) {
	email := c.Query("email")
	isRegister := c.Query("is_register")
	err := service.GroupApp.User.GetVerificationCode(email, isRegister)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// Reset password
// @Summary Reset password
// @Description Reset password
// @Tags user
// @Accept json
// @Produce json
// @Param reset_password_req body model.ResetPasswordReq true "Reset password request"
// @Success 200 {object} model.ResetPasswordRes
// @Router   /api/v1/reset/password [post]
func (*UserApi) ResetPassword(c *gin.Context) {
	var resetPasswordReq model.ResetPasswordReq
	if !BindAndValidate(c, &resetPasswordReq) {
		return
	}

	err := service.GroupApp.User.ResetPassword(c, &resetPasswordReq)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// CreateUser
// @Summary Create user
// @Description Create user
// @Tags user
// @Accept json
// @Produce json
// @Param create_user_req body model.CreateUserReq true "Create user request"
// @Success 200 {object} model.CreateUserRes
// @Router   /api/v1/user [post]
func (*UserApi) CreateUser(c *gin.Context) {
	var createUserReq model.CreateUserReq

	if !BindAndValidate(c, &createUserReq) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.User.CreateUser(&createUserReq, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// GetUserListByPage
// @Summary Get user list by page
// @Description Get user list by page
// @Tags user
// @Accept json
// @Produce json
// @Param user_list_req body model.UserListReq true "User list request"
// @Success 200 {object} model.UserListRes
// @Router   /api/v1/user [get]
func (*UserApi) HandleUserListByPage(c *gin.Context) {
	var userListReq model.UserListReq

	if !BindAndValidate(c, &userListReq) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	userList, err := service.GroupApp.User.GetUserListByPage(&userListReq, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", userList)
}

// UpdateUser
// @Summary Update user
// @Description Update user
// @Tags user
// @Accept json
// @Produce json
// @Param update_user_req body model.UpdateUserReq true "Update user request"
// @Success 200 {object} model.UpdateUserRes
// @Router   /api/v1/user [put]
func (*UserApi) UpdateUser(c *gin.Context) {
	var updateUserReq model.UpdateUserReq

	if !BindAndValidate(c, &updateUserReq) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.User.UpdateUser(&updateUserReq, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// DeleteUser
// @Summary Delete user
// @Description Delete user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User id"
// @Success 200 {object} model.DeleteUserRes
// @Router   /api/v1/user/{id} [delete]
func (*UserApi) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.User.DeleteUser(id, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// GetUser
// @Summary Get user
// @Description Get user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User id"
// @Success 200 {object} model.GetUserRes
// @Router   /api/v1/user/{id} [get]
func (*UserApi) HandleUser(c *gin.Context) {
	id := c.Param("id")

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	user, err := service.GroupApp.User.GetUser(id, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	user.Password = ""

	c.Set("data", user)
}

// GetUserDetail
// @Summary Get user detail
// @Description Get user detail
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} model.GetUserDetailRes
// @Router   /api/v1/user/detail [get]
func (*UserApi) HandleUserDetail(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	user, err := service.GroupApp.User.GetUserDetail(userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", user)
}

// UpdateUsers
// @Summary Update users
// @Description Update users
// @Tags user
// @Accept json
// @Produce json
// @Param update_user_info_req body model.UpdateUserInfoReq true "Update user info request"
// @Success 200 {object} model.UpdateUserInfoRes
// @Router   /api/v1/user/update [put]
func (*UserApi) UpdateUsers(c *gin.Context) {
	var updateUserInfoReq model.UpdateUserInfoReq

	if !BindAndValidate(c, &updateUserInfoReq) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.User.UpdateUserInfo(c, &updateUserInfoReq, userClaims)
	if err != nil {
		c.Error(err)
	}

	c.Set("data", nil)
}

// Transform user
// @Summary Transform user
// @Description Transform user
// @Tags user
// @Accept json
// @Produce json
// @Param transform_user_req body model.TransformUserReq true "Transform user request"
// @Success 200 {object} model.TransformUserRes
// @Router   /api/v1/user/transform [post]
func (*UserApi) TransformUser(c *gin.Context) {
	var transformUserReq model.TransformUserReq

	if !BindAndValidate(c, &transformUserReq) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	loginRsp, err := service.GroupApp.User.TransformUser(&transformUserReq, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", loginRsp)
}

// EmailRegister
// @Summary Email register
// @Description Email register
// @Tags user
// @Accept json
// @Produce json
// @Param email_register_req body model.EmailRegisterReq true "Email register request"
// @Success 200 {object} model.EmailRegisterRes
// @Router   /api/v1/tenant/email/register [post]
func (*UserApi) EmailRegister(c *gin.Context) {
	var req model.EmailRegisterReq
	if !BindAndValidate(c, &req) {
		return
	}
	loginRsp, err := service.GroupApp.EmailRegister(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", loginRsp)
}

// GetTenantID
// @Summary Get tenant id
// @Description Get tenant id
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} model.GetTenantIDRes
// @Router   /api/v1/user/tenant/id [get]
func (*UserApi) GetTenantID(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	tenantID := userClaims.TenantID

	c.Set("data", tenantID)
}
