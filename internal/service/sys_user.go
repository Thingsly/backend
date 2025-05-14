package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Thingsly/backend/pkg/common"
	"github.com/Thingsly/backend/pkg/errcode"

	"gorm.io/gorm"

	"github.com/Thingsly/backend/initialize"
	dal "github.com/Thingsly/backend/internal/dal"
	"github.com/Thingsly/backend/internal/logic"
	model "github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"
	global "github.com/Thingsly/backend/pkg/global"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type User struct{}

func (u *User) CreateUser(createUserReq *model.CreateUserReq, claims *utils.UserClaims) error {
	user := model.User{}

	user.ID = uuid.New()
	user.Name = createUserReq.Name
	user.PhoneNumber = createUserReq.PhoneNumber
	user.Email = createUserReq.Email
	user.Status = StringPtr("N")
	user.Remark = createUserReq.Remark

	if createUserReq.AdditionalInfo == nil {
		user.AdditionalInfo = StringPtr("{}")
	} else {
		var js map[string]interface{}
		if err := json.Unmarshal(*createUserReq.AdditionalInfo, &js); err != nil {
			return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
				"error": fmt.Sprintf("Failed to unmarshal AdditionalInfo: %v", err),
			})
		}
		user.AdditionalInfo = StringPtr(string(*createUserReq.AdditionalInfo))
	}

	switch claims.Authority {
	case "SYS_ADMIN": // System admin creates tenant admin
		user.Authority = StringPtr("TENANT_ADMIN")
		user.TenantID = StringPtr(strings.Split(uuid.New(), "-")[0])
	case "TENANT_ADMIN": // Tenant admin creates tenant user
		user.Authority = StringPtr("TENANT_USER")
		a, err := u.GetUserById(claims.ID)
		if err != nil {
			logrus.Error(err)
			return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"error":    err.Error(),
				"admin_id": claims.ID,
			})
		}
		user.TenantID = a.TenantID
	default:
		// If the current user is not a system admin or tenant admin, return an error
		return errcode.WithVars(errcode.CodeNoPermission, map[string]interface{}{
			"required_role": "SYS_ADMIN or TENANT_ADMIN",
			"current_role":  claims.Authority,
		})
	}
	t := time.Now().UTC()
	user.CreatedAt = &t
	user.UpdatedAt = &t
	user.PasswordLastUpdated = &t

	if len(createUserReq.Password) < 6 {
		return errcode.New(200040)
	}

	hashedPassword := utils.BcryptHash(createUserReq.Password)
	if hashedPassword == "" {
		return errcode.WithData(errcode.CodeDecryptError, map[string]interface{}{
			"error": "Failed to hash password",
		})
	}
	user.Password = hashedPassword

	err := dal.CreateUsers(&user)
	if err != nil {
		logrus.Error(err)
		if strings.Contains(err.Error(), "users_un") {
			return errcode.New(200008)
		}
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":      err.Error(),
			"user_email": user.Email,
		})
	}

	// If the user is a tenant admin, create a default home dashboard for the tenant
	if claims.Authority == "SYS_ADMIN" {
		err = dal.BoardQuery{}.CreateDefaultBoard(context.Background(), *user.TenantID)
		if err != nil {
			logrus.Error(err)
		}
	}

	if len(createUserReq.RoleIDs) > 0 {
		ok := GroupApp.Casbin.AddRolesToUser(user.ID, createUserReq.RoleIDs)
		if !ok {
			logrus.Error("Failed to add roles to user")
			return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
				"error":    "Failed to add roles to user",
				"user_id":  user.ID,
				"role_ids": createUserReq.RoleIDs,
			})
		}
	}
	return err
}

func (u *User) Login(ctx context.Context, loginReq *model.LoginReq) (*model.LoginRsp, error) {

	user, err := dal.GetUsersByEmail(loginReq.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, errcode.New(errcode.CodeInvalidAuth)
		}

		return nil, errcode.New(errcode.CodeDBError)
	}

	if logic.UserIsEncrypt(ctx) {
		password, err := initialize.DecryptPassword(loginReq.Password)
		if err != nil {
			return nil, errcode.New(errcode.CodeDecryptError)
		}
		passwords := strings.TrimSuffix(string(password), loginReq.Salt)
		loginReq.Password = passwords
	}

	if !utils.BcryptCheck(loginReq.Password, user.Password) {
		return nil, errcode.New(errcode.CodeInvalidAuth)
	}

	if *user.Status != "N" {
		return nil, errcode.New(errcode.CodeUserDisabled)
	}

	logrsp, err := u.UserLoginAfter(user)
	if err != nil {
		return nil, err
	}

	err = dal.UserQuery{}.UpdateLastVisitTime(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return logrsp, nil
}

func (*User) UserLoginAfter(user *model.User) (*model.LoginRsp, error) {
	key := viper.GetString("jwt.key")

	jwt := utils.NewJWT([]byte(key))
	claims := utils.UserClaims{
		ID:         user.ID,
		Email:      user.Email,
		Authority:  *user.Authority,
		CreateTime: time.Now().UTC(),
		TenantID:   *user.TenantID,
	}
	token, err := jwt.GenerateToken(claims)
	if err != nil {
		return nil, errcode.New(errcode.CodeTokenGenerateError)
	}
	timeout := viper.GetInt("session.timeout")
	reset_on_request := viper.GetBool("session.reset_on_request")
	if reset_on_request {
		if timeout == 0 {

			timeout = 60
		}
	}

	global.REDIS.Set(context.Background(), token, "1", time.Duration(timeout)*time.Minute)

	if !logic.UserIsShare(context.Background()) {
		oldToken, err := global.REDIS.Get(context.Background(), user.Email+"_token").Result()
		if err != nil {
			logrus.Error(err)
		} else {
			global.REDIS.Del(context.Background(), oldToken)
		}
		global.REDIS.Set(context.Background(), user.Email+"_token", token, 0)
	}

	loginRsp := &model.LoginRsp{
		Token:     &token,
		ExpiresIn: int64(timeout * 60),
	}
	return loginRsp, nil
}

func (*User) Logout(token string) error {
	if err := global.REDIS.Del(context.Background(), token).Err(); err != nil {
		return errcode.New(errcode.CodeTokenDeleteError)
	}
	return nil
}

func (*User) RefreshToken(userClaims *utils.UserClaims) (*model.LoginRsp, error) {

	user, err := dal.GetUsersByEmail(userClaims.Email)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"operation": "query_user",
			"email":     userClaims.Email,
			"error":     err.Error(),
		})
	}

	if *user.Status != "N" {
		return nil, errcode.New(errcode.CodeUserDisabled)
	}

	key := viper.GetString("jwt.key")

	jwt := utils.NewJWT([]byte(key))
	claims := utils.UserClaims{
		ID:         user.ID,
		Email:      user.Email,
		Authority:  *user.Authority,
		CreateTime: time.Now().UTC(),
		TenantID:   *user.TenantID,
	}
	token, err := jwt.GenerateToken(claims)
	if err != nil {
		return nil, errcode.New(errcode.CodeTokenGenerateError)
	}

	global.REDIS.Set(context.Background(), token, "1", 24*7*time.Hour)

	loginRsp := &model.LoginRsp{
		Token:     &token,
		ExpiresIn: int64(24 * 7 * time.Hour.Seconds()),
	}
	return loginRsp, nil
}

func (*User) GetVerificationCode(email, isRegister string) error {
	user, err := dal.GetUsersByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"operation": "query_user",
			"email":     email,
			"error":     err.Error(),
		})
	}

	switch {
	case user == nil && isRegister != "1":
		return errcode.New(200007)
	case user != nil && isRegister == "1":
		return errcode.New(200008)
	}

	verificationCode, err := common.GenerateNumericCode(6)
	if err != nil {
		return errcode.WithData(200009, map[string]interface{}{
			"email": email,
		})
	}

	err = global.REDIS.Set(context.Background(), email+"_code", verificationCode, 5*time.Minute).Err()
	if err != nil {
		return errcode.WithData(errcode.CodeCacheError, map[string]interface{}{
			"operation": "save_verification_code",
			"email":     email,
			"error":     err.Error(),
		})
	}

	logrus.Warningf("Verification code:%s", verificationCode)
	err = GroupApp.NotificationServicesConfig.SendTestEmail(&model.SendTestEmailReq{
		Email: email,
		Body: fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #f9f9f9;
            border-radius: 8px;
            padding: 30px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            max-width: 150px;
            margin-bottom: 20px;
        }
        .verification-code {
            background-color: #f0f0f0;
            padding: 15px;
            border-radius: 4px;
            text-align: center;
            font-size: 24px;
            font-weight: bold;
            color: #2c3e50;
            margin: 20px 0;
            letter-spacing: 5px;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #666;
        }
        .important {
            color: #e74c3c;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verification Code</h1>
        </div>
        <p>Hello,</p>
        <p>Thank you for using Thingsly Platform. To complete your verification, please use the following code:</p>
        
        <div class="verification-code">%s</div>
        
        <p class="important">This code will expire in 5 minutes.</p>
        
        <p>If you didn't request this verification code, please ignore this email or contact our support team if you have concerns.</p>
        
        <div class="footer">
            <p>This is an automated message, please do not reply to this email.</p>
            <p>&copy; 2025 Thingsly Platform. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, verificationCode),
	})
	if err != nil {
		return errcode.WithData(200010, map[string]interface{}{
			"email": email,
			"error": err.Error(),
		})
	}
	return nil
}

func (*User) ResetPassword(ctx context.Context, resetPasswordReq *model.ResetPasswordReq) error {
	if err := utils.ValidatePassword(resetPasswordReq.Password); err != nil {
		return err
	}

	verificationCode, err := global.REDIS.Get(context.Background(), resetPasswordReq.Email+"_code").Result()
	if err != nil {
		return errcode.New(200011)
	}
	if verificationCode != resetPasswordReq.VerifyCode {
		return errcode.New(200012)
	}

	var (
		db   = dal.UserQuery{}
		user = query.User
	)
	info, err := db.First(ctx, user.Email.Eq(resetPasswordReq.Email))
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"operation": "query_user",
			"email":     resetPasswordReq.Email,
			"error":     err.Error(),
		})
	}
	t := time.Now().UTC()
	info.PasswordLastUpdated = &t
	info.Password = utils.BcryptHash(resetPasswordReq.Password)
	if err = db.UpdateByEmail(ctx, info, user.Password, user.PasswordLastUpdated); err != nil {
		logrus.Error(ctx, "[ResetPasswordByCode]Update Users info failed:", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"operation": "update_password",
			"email":     resetPasswordReq.Email,
			"error":     err.Error(),
		})
	}
	return nil
}

func (*User) GetUserById(id string) (*model.User, error) {
	user, err := dal.GetUsersById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (*User) GetUserListByPage(userListReq *model.UserListReq, claims *utils.UserClaims) (map[string]interface{}, error) {
	total, list, err := dal.GetUserListByPage(userListReq, claims)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"operation": "query_user",
			"error":     err.Error(),
		})
	}
	userListRspMap := make(map[string]interface{})
	userListRspMap["total"] = total
	userListRspMap["list"] = list
	return userListRspMap, nil
}

func (*User) UpdateUser(updateUserReq *model.UpdateUserReq, claims *utils.UserClaims) error {

	if updateUserReq.Password != nil {
		if len(*updateUserReq.Password) == 0 {
			updateUserReq.Password = nil
		} else if len(*updateUserReq.Password) < 6 {
			return errcode.New(200040)
		}
	}

	user, err := dal.GetUsersById(updateUserReq.ID)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": updateUserReq.ID,
		})
	}

	if claims.Authority == "TENANT_ADMIN" || claims.Authority == "TENANT_USER" {
		if *user.TenantID != claims.TenantID {
			return errcode.New(errcode.CodeNoPermission)
		}

		if claims.Authority == "TENANT_ADMIN" && *user.Authority == "TENANT_ADMIN" && *user.Status != *updateUserReq.Status {
			if updateUserReq.Status != nil {
				if updateUserReq.Status != nil {
					return errcode.New(errcode.CodeOpDenied)
				}
			}
		}
	}

	t := time.Now().UTC()

	if updateUserReq.Password != nil {
		hashedPassword := utils.BcryptHash(*updateUserReq.Password)
		if hashedPassword == "" {
			return errcode.New(errcode.CodeDecryptError)
		}
		user.Password = hashedPassword
		user.PasswordLastUpdated = &t
	}

	user.UpdatedAt = &t
	user.Name = updateUserReq.Name
	user.PhoneNumber = *updateUserReq.PhoneNumber
	user.AdditionalInfo = updateUserReq.AdditionalInfo
	user.Status = updateUserReq.Status
	user.Remark = updateUserReq.Remark

	_, err = dal.UpdateUserInfoById(claims.ID, user)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": claims.ID,
		})
	}

	if updateUserReq.RoleIDs != nil {

		GroupApp.Casbin.RemoveUserAndRole(updateUserReq.ID)

		if len(updateUserReq.RoleIDs) > 0 {
			ok := GroupApp.Casbin.AddRolesToUser(updateUserReq.ID, updateUserReq.RoleIDs)
			if !ok {
				return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
					"error":    "Failed to update user roles",
					"user_id":  updateUserReq.ID,
					"role_ids": updateUserReq.RoleIDs,
				})
			}
		}
	}

	return nil
}

func (*User) DeleteUser(id string, claims *utils.UserClaims) error {

	user, err := dal.GetUsersById(id)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		})
	}

	if claims.Authority == "TENANT_ADMIN" || claims.Authority == "TENANT_USER" {

		if *user.TenantID != claims.TenantID {
			return errcode.WithVars(errcode.CodeNoPermission, map[string]interface{}{
				"required_tenant": *user.TenantID,
				"current_tenant":  claims.TenantID,
				"operation":       "delete_user",
			})
		}

		// if claims.Authority == "TENANT_ADMIN" && *user.Authority == "TENANT_ADMIN" {
		// 	return errcode.WithVars(errcode.CodeOpDenied, map[string]interface{}{
		// 		"reason":  "cannot_delete_self",
		// 		"user_id": id,
		// 	})
		// }
	}

	if *user.Authority == "SYS_ADMIN" {
		return errcode.WithVars(errcode.CodeOpDenied, map[string]interface{}{
			"reason":  "cannot_delete_sys_admin",
			"user_id": id,
		})
	}

	err = dal.DeleteUsersById(id)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":     err.Error(),
			"user_id":   id,
			"operation": "delete_user",
		})
	}

	return nil
}

func (*User) GetUser(id string, claims *utils.UserClaims) (*model.User, error) {

	user, err := dal.GetUsersById(id)
	if err != nil {

		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		})
	}

	if claims.Authority == "TENANT_ADMIN" || claims.Authority == "TENANT_USER" {
		if *user.TenantID != claims.TenantID {
			return nil, errcode.WithVars(errcode.CodeNoPermission, map[string]interface{}{
				"required_tenant": *user.TenantID,
				"current_tenant":  claims.TenantID,
				"user_authority":  claims.Authority,
			})
		}
	}

	return user, nil
}

func (*User) GetUserDetail(claims *utils.UserClaims) (*model.User, error) {
	user, err := dal.GetUsersById(claims.ID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": claims.ID,
		})
	}
	return user, nil
}

func (*User) UpdateUserInfo(ctx context.Context, updateUserReq *model.UpdateUserInfoReq, claims *utils.UserClaims) error {
	user, err := dal.GetUsersById(claims.ID)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": claims.ID,
		})
	}

	if user.ID != claims.ID {
		return errcode.WithVars(errcode.CodeNoPermission, map[string]interface{}{
			"reason":  "cannot_update_other_user_info",
			"user_id": claims.ID,
		})
	}

	if logic.UserIsEncrypt(ctx) {
		password, err := initialize.DecryptPassword(*updateUserReq.Password)
		if err != nil {
			return errcode.WithData(errcode.CodeDecryptError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		passwords := strings.TrimSuffix(string(password), updateUserReq.Salt)
		*updateUserReq.Password = passwords
	}

	if updateUserReq.Password != nil {
		updateUserReq.Password = StringPtr(utils.BcryptHash(*updateUserReq.Password))
	}

	r, err := dal.UpdateUserInfoByIdPersonal(user.ID, updateUserReq)
	if r == 0 {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": claims.ID,
		})
	}
	return err
}

func (*User) TransformUser(transformUserReq *model.TransformUserReq, claims *utils.UserClaims) (*model.LoginRsp, error) {

	if claims.Authority != "SYS_ADMIN" && claims.Authority != "TENANT_ADMIN" {
		return nil, errcode.WithVars(errcode.CodeNoPermission, map[string]interface{}{
			"required_authority": "SYS_ADMIN or TENANT_ADMIN",
			"current_authority":  claims.Authority,
		})
	}

	becomeUser, err := dal.GetUsersById(transformUserReq.BecomeUserID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": transformUserReq.BecomeUserID,
		})
	}

	if *becomeUser.Status != "N" {
		return nil, errcode.WithVars(errcode.CodeUserDisabled, map[string]interface{}{
			"user_id":         becomeUser.ID,
			"current_status":  *becomeUser.Status,
			"required_status": "N",
		})
	}

	key := viper.GetString("jwt.key")
	if key == "" {
		return nil, errcode.New(errcode.CodeSystemError)
	}

	becomeUserClaims := utils.UserClaims{
		ID:         becomeUser.ID,
		Email:      becomeUser.Email,
		Authority:  *becomeUser.Authority,
		CreateTime: time.Now().UTC(),
		TenantID:   *becomeUser.TenantID,
	}

	jwt := utils.NewJWT([]byte(key))
	token, err := jwt.GenerateToken(becomeUserClaims)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeTokenGenerateError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": becomeUser.ID,
		})
	}

	err = global.REDIS.Set(context.Background(), token, "1", 24*7*time.Hour).Err()
	if err != nil {
		return nil, errcode.WithData(errcode.CodeTokenSaveError, map[string]interface{}{
			"error":   err.Error(),
			"user_id": becomeUser.ID,
		})
	}

	loginRsp := &model.LoginRsp{
		Token:     &token,
		ExpiresIn: int64(24 * 7 * time.Hour.Seconds()),
	}

	return loginRsp, nil
}

func (u *User) EmailRegister(ctx context.Context, req *model.EmailRegisterReq) (*model.LoginRsp, error) {

	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	verificationCode, err := global.REDIS.Get(context.Background(), req.Email+"_code").Result()
	if err != nil {
		return nil, errcode.New(200011)
	}
	if verificationCode != req.VerifyCode {
		return nil, errcode.New(200012)
	}

	if req.ConfirmPassword != nil && *req.ConfirmPassword != req.Password {
		return nil, errcode.New(200041)
	}

	user, err := dal.GetUsersByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"operation": "query_user",
			"email":     req.Email,
			"error":     err.Error(),
		})
	}
	if user != nil {
		return nil, errcode.New(200008)
	}

	if logic.UserIsEncrypt(ctx) {
		if req.Salt == nil {
			return nil, errcode.New(200042)
		}
		password, err := initialize.DecryptPassword(req.Password)
		if err != nil {
			return nil, errcode.New(200043)
		}
		passwords := strings.TrimSuffix(string(password), *req.Salt)
		req.Password = passwords
	}

	req.Password = utils.BcryptHash(req.Password)

	now := time.Now().UTC()
	tenantID, err := common.GenerateRandomString(8)
	if err != nil {
		logrus.Error("Failed to generate tenant ID ", err)
		return nil, errcode.New(errcode.CodeSystemError)
	}

	userInfo := &model.User{
		ID:                  uuid.New(),
		Name:                &req.Email,
		PhoneNumber:         fmt.Sprintf("%s %s", req.PhonePrefix, req.PhoneNumber),
		Email:               req.Email,
		Status:              StringPtr("N"),
		Authority:           StringPtr("TENANT_ADMIN"),
		Password:            req.Password,
		TenantID:            StringPtr(tenantID),
		Remark:              StringPtr(now.Add(365 * 24 * time.Hour).String()),
		CreatedAt:           &now,
		UpdatedAt:           &now,
		PasswordLastUpdated: &now,
	}

	if err = dal.CreateUsers(userInfo); err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"operation": "create_user",
			"email":     req.Email,
			"error":     err.Error(),
		})
	}

	err = dal.BoardQuery{}.CreateDefaultBoard(ctx, tenantID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"operation": "create_default_board",
			"tenant_id": tenantID,
			"error":     err.Error(),
		})
	}

	return u.UserLoginAfter(userInfo)
}

func (u *User) GetUserEmailByPhoneNumber(phoneNumber string) (string, error) {

	user, err := dal.GetUsersByPhoneNumber(phoneNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errcode.New(200007)
		}
		return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"message": "get_user_by_phone_number",
			"error":   err.Error(),
		})
	}
	return user.Email, nil
}

func (u *User) GetTenantInfo(tenantID string) (*model.User, error) {
	tenant, err := dal.GetTenantsById(tenantID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error":     err.Error(),
			"tenant_id": tenantID,
		})
	}
	return tenant, nil
}
