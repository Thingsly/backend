package dal

import (
	"context"

	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/query"

	"github.com/sirupsen/logrus"
)

func DeleteServiceAccess(id string) error {
	q := query.ServiceAccess
	queryBuilder := q.WithContext(context.Background())
	_, err := queryBuilder.Where(q.ID.Eq(id)).Delete()
	return err
}

func UpdateServiceAccess(id string, updates map[string]interface{}) error {
	q := query.ServiceAccess
	queryBuilder := q.WithContext(context.Background())
	_, err := queryBuilder.Where(q.ID.Eq(id)).Updates(updates)
	return err
}

func GetServiceAccessListByPage(req *model.GetServiceAccessByPageReq, tenantID string) (int64, interface{}, error) {
	var count int64
	var serviceAccess = []model.ServiceAccess{}

	q := query.ServiceAccess
	queryBuilder := q.WithContext(context.Background())
	queryBuilder = queryBuilder.Where(q.ServicePluginID.Eq(req.ServicePluginID))
	queryBuilder = queryBuilder.Where(q.TenantID.Eq(tenantID))

	count, err := queryBuilder.Count()
	if err != nil {
		logrus.Error(err)
		return count, serviceAccess, err
	}
	if req.Page != 0 && req.PageSize != 0 {
		queryBuilder = queryBuilder.Limit(req.PageSize)
		queryBuilder = queryBuilder.Offset((req.Page - 1) * req.PageSize)
	}

	err = queryBuilder.Select().Order(q.CreateAt.Desc()).Scan(&serviceAccess)
	if err != nil {
		logrus.Error(err)
		return count, serviceAccess, err
	}
	return count, serviceAccess, err
}

func GetServiceAccessByVoucher(voucher string, tenantID string) (*model.ServiceAccess, error) {

	q := query.ServiceAccess
	queryBuilder := q.WithContext(context.Background())
	serviceAccess, err := queryBuilder.Where(q.Voucher.Eq(voucher)).Where(q.TenantID.Eq(tenantID)).First()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return serviceAccess, nil
}

func GetServiceAccessListByServicePluginID(servicePluginID string) ([]model.ServiceAccess, error) {
	q := query.ServiceAccess
	queryBuilder := q.WithContext(context.Background())
	var serviceAccess = []model.ServiceAccess{}
	err := queryBuilder.Where(q.ServicePluginID.Eq(servicePluginID)).Select().Scan(&serviceAccess)
	if err != nil {
		logrus.Error(err)
		return serviceAccess, err
	}
	return serviceAccess, nil
}

func GetServiceAccessByID(id string) (*model.ServiceAccess, error) {

	q := query.ServiceAccess
	queryBuilder := q.WithContext(context.Background())
	serviceAccess, err := queryBuilder.Where(q.ID.Eq(id)).First()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return serviceAccess, nil
}
