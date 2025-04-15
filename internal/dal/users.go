package dal

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"
	common "github.com/Thingsly/backend/pkg/common"
	global "github.com/Thingsly/backend/pkg/global"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

const (
	SYS_ADMIN    = "SYS_ADMIN"
	TENANT_ADMIN = "TENANT_ADMIN"
	TENANT_USER  = "TENANT_USER"
)

func CreateUsers(user *model.User) error {
	return query.User.Create(user)
}

func GetUsersById(uid string) (*model.User, error) {
	user, err := query.User.Where(query.User.ID.Eq(uid)).First()
	if err != nil {
		return nil, err
	}
	return user, err
}

func GetUsersByEmail(email string) (*model.User, error) {
	q := query.User
	user, err := q.Where(q.Email.Eq(email)).First()
	if err != nil {
		return nil, err
	}
	return user, err
}

func GetUsersByPhoneNumber(phoneNumber string) (*model.User, error) {

	format1 := phoneNumber
	format2 := phoneNumber

	if strings.HasPrefix(phoneNumber, "+") {

		phone := phoneNumber[1:]
		parts := strings.Split(phone, " ")
		if len(parts) > 1 {
			// Case with spaces: +84 18211111111
			format2 = parts[1]
		} else if len(phone) > 4 {
			// Without space: +8418211111111
			// Assume area code does not exceed 4 digits
			format2 = phone[2:] // Remove the first two digits
		}
	} else {
		// If it does not start with a +, add the +84 prefix as the second format
		format2 = "+84 " + phoneNumber
	}

	q := query.User
	user, err := q.Where(q.PhoneNumber.Eq(format1)).
		Or(q.PhoneNumber.Eq(format2)).
		First()

	return user, err
}

func GetUserListByPage(userListReq *model.UserListReq, claims *utils.UserClaims) (int64, interface{}, error) {
	q := query.User
	var count int64
	var userList []map[string]interface{}
	queryBuilder := q.WithContext(context.Background())

	if claims.Authority == TENANT_ADMIN || claims.Authority == TENANT_USER {
		queryBuilder = queryBuilder.Where(q.TenantID.Eq(claims.TenantID))
		queryBuilder = queryBuilder.Where(q.Authority.Eq(TENANT_USER))
	} else if claims.Authority == SYS_ADMIN {
		queryBuilder = queryBuilder.Where(q.Authority.Eq(TENANT_ADMIN))
	} else {
		return count, nil, fmt.Errorf("authority exception")
	}

	if userListReq.Email != nil && *userListReq.Email != "" {
		queryBuilder = queryBuilder.Where(q.Email.Like(fmt.Sprintf("%%%s%%", *userListReq.Email)))
	}
	if userListReq.PhoneNumber != nil && *userListReq.PhoneNumber != "" {
		queryBuilder = queryBuilder.Where(q.PhoneNumber.Eq(*userListReq.PhoneNumber))
	}
	if userListReq.Name != nil && *userListReq.Name != "" {
		queryBuilder = queryBuilder.Where(q.Name.Like(fmt.Sprintf("%%%s%%", *userListReq.Name)))
	}
	if userListReq.Status != nil && *userListReq.Status != "" {
		queryBuilder = queryBuilder.Where(q.Status.Eq(*userListReq.Status))
	}
	count, err := queryBuilder.Count()
	if err != nil {
		return count, nil, err
	}
	if userListReq.Page != 0 && userListReq.PageSize != 0 {
		queryBuilder = queryBuilder.Limit(userListReq.PageSize)
		queryBuilder = queryBuilder.Offset((userListReq.Page - 1) * userListReq.PageSize)
	}

	users, err := queryBuilder.Select(q.ID, q.Name, q.PhoneNumber, q.Email, q.Status, q.Authority, q.TenantID, q.Remark, q.AdditionalInfo, q.CreatedAt, q.UpdatedAt, q.LastVisitTime).Order(q.CreatedAt.Desc()).Find()
	if err != nil {
		return count, users, err
	}
	for _, user := range users {
		roles, _ := GetRolesByUserId(user.ID)
		userMap := map[string]interface{}{
			"id":             user.ID,
			"name":           user.Name,
			"phone_number":   user.PhoneNumber,
			"email":          user.Email,
			"status":         user.Status,
			"authority":      user.Authority,
			"tenant_id":      user.TenantID,
			"remark":         user.Remark,
			"additionalInfo": user.AdditionalInfo,
			"created_at":     user.CreatedAt,
			"updated_at":     user.UpdatedAt,
			"userRoles":      roles,
			"lastVisitTime":  user.LastVisitTime,
		}
		userList = append(userList, userMap)
	}

	return count, userList, err
}

func UpdateUserInfoByIdPersonal(uid string, data *model.UpdateUserInfoReq) (int64, error) {
	q := query.User
	t := time.Now()
	data.UpdatedAt = &t
	r, err := query.User.Where(q.ID.Eq(uid)).Updates(data)
	return r.RowsAffected, err
}

func UpdateUserInfoById(_ string, data *model.User) (int64, error) {
	q := query.User
	r, err := query.User.Where(q.ID.Eq(data.ID)).Updates(data)
	return r.RowsAffected, err
}

func DeleteUsersById(uid string) error {
	_, err := query.User.Where(query.User.ID.Eq(uid)).Delete()
	return err
}

func GetUserIdBYTenantID(tenantID string) (string, error) {
	var (
		userId     string
		cacheKeyId = fmt.Sprintf("GetUserIdBYTenantID:%s", tenantID)
		err        error
	)
	userId, err = global.REDIS.Get(context.Background(), cacheKeyId).Result()
	if err == nil {
		return userId, nil
	}
	err = query.User.Where(query.User.TenantID.Eq(tenantID)).Select(query.User.ID).Scan(&userId)
	if err != nil {
		return userId, err
	}
	global.REDIS.Set(context.Background(), cacheKeyId, userId, time.Hour*6)
	return userId, nil
}

type UserQuery struct {
}

func (UserQuery) Count(ctx context.Context) (count int64, err error) {
	count, err = query.User.Count()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (UserQuery) CountByWhere(ctx context.Context, option ...gen.Condition) (count int64, err error) {
	var users = query.User
	count, err = users.Where(option...).Count()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (UserQuery) GroupByMonthCount(ctx context.Context, email *string) (list []*model.GetBoardUserListMonth) {
	var (
		db = global.DB.WithContext(ctx)
	)
	conn := db.Model(&model.User{}).Select("(EXTRACT(MONTH FROM created_at) ) AS mon,COUNT(1) as num").
		Where("created_at > ? and created_at  IS NOT NULL", common.GetYearStart()).
		Group("EXTRACT(MONTH FROM created_at)").Order("mon")

	if email != nil {
		conn = conn.Where("email = ?", *email)
	}

	conn.Scan(&list)

	return
}

func (UserQuery) First(ctx context.Context, option ...gen.Condition) (info *model.User, err error) {
	var users = query.User

	info, err = users.Where(option...).First()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (UserQuery) Select(ctx context.Context, option ...gen.Condition) (list []*model.User, err error) {
	var users = query.User

	list, err = users.Where(option...).Find()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (UserQuery) UpdateByEmail(ctx context.Context, info *model.User, columns ...field.Expr) (err error) {
	var users = query.User
	//users.Password, users.Name, users.PhoneNumber, users.Remark
	_, err = users.Where(users.Email.Eq(info.Email)).
		Select(columns...).
		UpdateColumns(info)
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (UserQuery) UpdateLastVisitTime(ctx context.Context, uid string) (err error) {
	_, err = query.User.Where(query.User.ID.Eq(uid)).Update(query.User.LastVisitTime, time.Now())
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

type UserVo struct {
}

func (UserVo) PoToVo(userInfo *model.User) (info *model.UsersRes) {
	info = &model.UsersRes{
		ID:       userInfo.ID,
		PhoneNum: userInfo.PhoneNumber,
		Email:    userInfo.Email,
	}
	if userInfo.Name != nil {
		info.Name = *userInfo.Name
	}
	if userInfo.Authority != nil {
		info.Authority = *userInfo.Authority
	}
	if userInfo.TenantID != nil {
		info.TenantID = *userInfo.TenantID
	}
	if userInfo.Remark != nil {
		info.Remark = *userInfo.Remark
	}
	if userInfo.CreatedAt != nil {
		info.CreateTime = common.DateTimeToString(*userInfo.CreatedAt, "")
	}
	if userInfo.AdditionalInfo != nil {
		info.AdditionalInfo = *userInfo.AdditionalInfo
	}
	return
}

func (UserVo) GetTenantAdminList() (list []*model.User, err error) {
	var users = query.User
	userInfoList, err := users.Where(users.Authority.Eq(TENANT_ADMIN)).Find()
	if err != nil {
		logrus.Error(err)
		return
	}
	return userInfoList, nil
}

func GetTenantsById(tenantID string) (info *model.User, err error) {
	var tenants = query.User
	info, err = tenants.Where(tenants.TenantID.Eq(tenantID), tenants.Authority.Eq(TENANT_ADMIN)).First()
	if err != nil {
		logrus.Error(err)
		return
	}
	return info, nil
}
