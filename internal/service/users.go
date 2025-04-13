package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/HustIoTPlatform/backend/initialize"
	dal "github.com/HustIoTPlatform/backend/internal/dal"
	"github.com/HustIoTPlatform/backend/internal/logic"
	model "github.com/HustIoTPlatform/backend/internal/model"
	query "github.com/HustIoTPlatform/backend/internal/query"
	common "github.com/HustIoTPlatform/backend/pkg/common"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gen/field"
)

type UsersService struct{}

// GetTenant
func (*UsersService) GetTenant(ctx context.Context) (model.GetTenantRes, error) {
	var (
		list []*model.GetBoardUserListMonth
		data model.GetTenantRes

		user = query.User
		db   = dal.UserQuery{}
	)

	total, err := db.Count(ctx)
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	yesterday, err := db.CountByWhere(ctx, user.CreatedAt.Gte(common.GetYesterdayBegin()))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	month, err := db.CountByWhere(ctx, user.CreatedAt.Gte(common.GetMonthStart()))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	list = db.GroupByMonthCount(ctx, nil)

	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		return data, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	data = model.GetTenantRes{
		UserTotal:          total,
		UserAddedYesterday: yesterday,
		UserAddedMonth:     month,
		UserListMonth:      list,
	}
	return data, err
}

// GetTenantUserInfo
func (*UsersService) GetTenantUserInfo(ctx context.Context, email string) (model.GetTenantRes, error) {
	var (
		err                     error
		total, yesterday, month int64
		list                    []*model.GetBoardUserListMonth
		data                    model.GetTenantRes

		user = query.User
		db   = dal.UserQuery{}
	)

	total, err = db.CountByWhere(ctx, user.Email.Eq(email))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	yesterday, err = db.CountByWhere(ctx, user.CreatedAt.Gte(common.GetYesterdayBegin()), user.Email.Eq(email))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	month, err = db.CountByWhere(ctx, user.CreatedAt.Gte(common.GetMonthStart()), user.Email.Eq(email))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	list = db.GroupByMonthCount(ctx, &email)

	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		return data, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	data = model.GetTenantRes{
		UserTotal:          total,
		UserAddedYesterday: yesterday,
		UserAddedMonth:     month,
		UserListMonth:      list,
	}
	return data, err
}

// GetTenantInfo
func (*UsersService) GetTenantInfo(ctx context.Context, email string) (*model.UsersRes, error) {
	var (
		info *model.UsersRes

		db   = dal.UserQuery{}
		user = query.User
	)

	UserInfo, err := db.First(ctx, user.Email.Eq(email))
	if err != nil {
		logrus.Error(ctx, "[GetTenantInfo]Users info failed:", err)
		return info, errcode.WithData(101001, map[string]interface{}{
			"error": err.Error(),
		})
	}
	info = dal.UserVo{}.PoToVo(UserInfo)

	return info, err
}

// UpdateTenantInfo
func (*UsersService) UpdateTenantInfo(ctx context.Context, userInfo *utils.UserClaims, param *model.UsersUpdateReq) error {
	var (
		db   = dal.UserQuery{}
		user = query.User
	)
	info, err := db.First(ctx, user.Email.Eq(userInfo.Email))
	if err != nil {
		logrus.Error(ctx, "[UpdateTenantInfo]Get Users info failed:", err)
		return errcode.WithData(101001, map[string]interface{}{
			"error": err.Error(),
		})
	}
	var columns []field.Expr
	columns = append(columns, user.Name)
	if param.Name != "" {
		info.Name = &param.Name
	}
	if param.AdditionalInfo != nil {
		info.AdditionalInfo = param.AdditionalInfo
		columns = append(columns, user.AdditionalInfo)
	}
	if param.PhoneNumber != nil {
		var phonePrefix string
		if param.PhonePrefix != nil {
			phonePrefix = *param.PhonePrefix
		}
		info.PhoneNumber = fmt.Sprintf("%s %s", phonePrefix, *param.PhoneNumber)
		columns = append(columns, user.PhoneNumber)
	}
	if err = db.UpdateByEmail(ctx, info, columns...); err != nil {
		logrus.Error(ctx, "[UpdateTenantInfo]Update Users info failed:", err)
		return errcode.WithData(101001, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return err
}

// UpdateTenantInfoPassword
func (*UsersService) UpdateTenantInfoPassword(ctx context.Context, userInfo *utils.UserClaims, param *model.UsersUpdatePasswordReq) error {

	if userInfo.Email == "test@hust.edu.vn" {
		return errcode.New(200044)
	}

	err := utils.ValidatePassword(param.Password)
	if err != nil {
		return err
	}

	var (
		db   = dal.UserQuery{}
		user = query.User
	)

	info, err := db.First(ctx, user.Email.Eq(userInfo.Email))
	if err != nil {
		logrus.Error(ctx, "[UpdateTenantInfoPassword]Get Users info failed:", err)
		return errcode.WithData(101001, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if logic.UserIsEncrypt(ctx) {
		password, err := initialize.DecryptPassword(param.Password)
		if err != nil {
			return errcode.New(200043)
		}
		passwords := strings.TrimSuffix(string(password), param.Salt)
		param.Password = passwords
	}

	if !utils.BcryptCheck(param.OldPassword, info.Password) {
		return errcode.New(200045)
	}

	t := time.Now().UTC()
	info.UpdatedAt = &t
	info.PasswordLastUpdated = &t

	info.Password = utils.BcryptHash(param.Password)
	if err = db.UpdateByEmail(ctx, info, user.Password, user.UpdatedAt, user.PasswordLastUpdated); err != nil {
		logrus.Error(ctx, "[UpdateTenantInfoPassword]Update Users info failed:", err)
		return errcode.WithData(101001, map[string]interface{}{
			"error": err.Error(),
			"email": userInfo.Email,
		})
	}

	return nil
}
