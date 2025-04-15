package dal

import (
	"context"
	"errors"
	"fmt"
	"strings"

	model "github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"gorm.io/gen"

	"github.com/sirupsen/logrus"
)

func CreateDevice(device *model.Device) error {
	return query.Device.Create(device)
}

// Batch create devices.
func CreateDeviceBatch(devices []*model.Device) error {
	return query.Device.Create(devices...)
}

func CreateDeviceBath(devices []*model.Device) error {
	return query.Device.Create(devices...)
}

func UpdateDevice(device *model.Device) (*model.Device, error) {
	info, err := query.Device.Where(query.Device.ID.Eq(device.ID)).Updates(device)
	if err != nil {
		logrus.Error(err)
		return nil, err
	} else if info.RowsAffected == 0 {
		return nil, fmt.Errorf("update device failed, no rows affected")
	}
	return device, err
}

func UpdateDeviceByMap(deviceID string, deviceMap map[string]interface{}) (*model.Device, error) {
	info, err := query.Device.Where(query.Device.ID.Eq(deviceID)).Updates(deviceMap)
	if err != nil {
		logrus.Error(err)
		return nil, err
	} else if info.RowsAffected == 0 {
		return nil, fmt.Errorf("update device failed, no rows affected")
	}
	// Query the updated data.
	device, err := query.Device.Where(query.Device.ID.Eq(deviceID)).First()
	if err != nil {
		logrus.Error(err)
	}
	return device, err
}

func UpdateDeviceStatus(deviceId string, status int16) error {
	info, err := query.Device.Where(query.Device.ID.Eq(deviceId)).Update(query.Device.IsOnline, status)
	if err != nil {
		logrus.Error(err)
	}
	if info.RowsAffected == 0 {
		return fmt.Errorf("update device status failed, no rows affected")
	}
	return err
}

func DeleteDevice(id string, tenantID string) error {
	_, err := query.Device.Where(query.Device.ID.Eq(id), query.Device.TenantID.Eq(tenantID)).Delete()
	if err != nil {
		logrus.Error(err)
	}
	return err
}

// DeleteDeviceWithTx deletes a device with the given ID and tenant ID using the provided transaction.
func DeleteDeviceWithTx(id string, tenantID string, tx *query.QueryTx) error {
	_, err := tx.Device.Where(query.Device.ID.Eq(id), query.Device.TenantID.Eq(tenantID)).Delete()
	if err != nil {
		logrus.Error(err)
	}
	return err
}

// GetParentDeviceBySubDeviceID retrieves the parent device information by sub-device ID.
func GetParentDeviceBySubDeviceID(subDeviceID string) (info *model.Device, err error) {
	device := query.Device
	info, err = device.Where(device.ID.Eq(subDeviceID)).First()
	if err != nil {
		logrus.Error(err)
	}
	return
}

func GetDeviceByID(id string) (*model.Device, error) {
	device, err := query.Device.Where(query.Device.ID.Eq(id)).First()
	if err != nil {
		logrus.Error(err)
	} else if device == nil {
		return nil, fmt.Errorf("device not found")
	}
	return device, err
}

// GetDeviceDetail retrieves device details by device ID.
func GetDeviceDetail(id string) (map[string]interface{}, error) {
	device := query.Device
	deviceConfig := query.DeviceConfig
	t := query.TelemetryCurrentData
	t2 := query.TelemetryCurrentData.As("t2")
	data := make(map[string]interface{})
	// Get device details by ID
	err := device.LeftJoin(deviceConfig, deviceConfig.ID.EqCol(device.DeviceConfigID)).
		LeftJoin(t.Select(t.T.Max().As("ts"), t.DeviceID).Group(t.DeviceID).As("t2"), t2.DeviceID.EqCol(device.ID)).
		Where(device.ID.Eq(id)).
		Select(device.ALL, deviceConfig.Name.As("device_config_name"), t2.T).Scan(&data)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if data["parent_id"] != nil {
		parentDevice, err := GetDeviceByID(data["parent_id"].(string))
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		data["gateway_device_name"] = parentDevice.Name
	}
	return data, err
}

// GetDeviceByVoucher retrieves device information by voucher.
func GetDeviceByVoucher(voucher string) (*model.Device, error) {
	device, err := query.Device.Where(query.Device.Voucher.Eq(voucher)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get device by voucher: %s failed: %v", voucher, err)
		}
		return nil, err
	}
	return device, err
}

// UpdateDeviceOnlineStatus updates the online status of a device.
func UpdateDeviceOnlineStatus(deviceId string, status int16) error {
	_, err := query.Device.Where(query.Device.ID.Eq(deviceId)).Update(query.Device.IsOnline, status)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

// GetDeviceByDeviceNumber retrieves device information by device number.
func GetDeviceByDeviceNumber(deviceNumber string) (*model.Device, error) {
	device, err := query.Device.Where(query.Device.DeviceNumber.Eq(deviceNumber)).First()
	if err != nil {
		logrus.Error(err)
	}
	return device, err
}

func GetDeviceBySubDeviceAddress(deviceAddress []string, parentId string) (map[string]*model.Device, error) {
	devices, err := query.Device.Where(query.Device.SubDeviceAddr.In(deviceAddress...)).
		Where(query.Device.ParentID.Eq(parentId)).Find()
	if err != nil {
		return nil, err
	}
	result := make(map[string]*model.Device)
	for _, d := range devices {
		result[*d.SubDeviceAddr] = d
	}
	return result, err
}

// RemoveSubDevice removes the sub-device information from the device.
func RemoveSubDevice(deviceId string, tenant_id string) error {
	info, err := query.Device.Where(query.Device.ID.Eq(deviceId), query.Device.TenantID.Eq(tenant_id)).UpdateSimple(query.Device.ParentID.Null(), query.Device.SubDeviceAddr.Null())
	if err != nil {
		logrus.Error(err)
	} else if info.RowsAffected == 0 {
		return fmt.Errorf("remove sub device failed, device not found")
	}
	return err
}

// GetDeviceListByPage retrieves a paginated list of devices based on the provided request parameters.
func GetDeviceListByPage(req *model.GetDeviceListByPageReq, tenant_id string) (int64, []model.GetDeviceListByPageRsp, error) {
	q := query.Device
	c := query.DeviceConfig
	lda := query.LatestDeviceAlarm
	var count int64
	deviceList := []model.GetDeviceListByPageRsp{}
	queryBuilder := q.WithContext(context.Background())

	queryBuilder = queryBuilder.Where(q.TenantID.Eq(tenant_id))

	if req.GroupId != nil && *req.GroupId != "" {

		groupIds, err := GetGroupChildrenIds(*req.GroupId)
		if err != nil {
			logrus.Error(err)
			return count, deviceList, err
		}
		ids, err2 := GetDeviceIdsByGroupIds(groupIds)
		if err2 != nil {
			logrus.Error(err2)
			return count, deviceList, err2
		}
		ids = append(ids, *req.GroupId)
		queryBuilder = queryBuilder.Where(q.ID.In(ids...))
	}

	queryBuilder = queryBuilder.Where(q.ActivateFlag.Eq("active"))

	if req.IsEnabled != nil && *req.IsEnabled != "" {
		queryBuilder = queryBuilder.Where(q.IsEnabled.Eq(*req.IsEnabled))
	}

	if req.ProductID != nil && *req.ProductID != "" {
		queryBuilder = queryBuilder.Where(q.ProductID.Eq(*req.ProductID))
	}

	if req.ServiceIdentifier != nil && *req.ServiceIdentifier != "" {
		if *req.ServiceIdentifier == "mqtt" {
			queryBuilder = queryBuilder.Where(query.Device.Where(c.ProtocolType.Eq(*req.ServiceIdentifier)).Or(q.DeviceConfigID.IsNull()))
		} else {
			queryBuilder = queryBuilder.Where(c.ProtocolType.Eq(*req.ServiceIdentifier))
		}
	}
	if req.ServiceAccessID != nil && *req.ServiceAccessID != "" {
		queryBuilder = queryBuilder.Where(q.ServiceAccessID.Eq(*req.ServiceAccessID))
	}
	if req.DeviceType != nil && *req.DeviceType != "" {
		if *req.DeviceType == "1" {
			queryBuilder = queryBuilder.Where(query.Device.Where(q.DeviceConfigID.IsNull()).Or(c.DeviceType.Eq(*req.DeviceType)))
		} else {
			queryBuilder = queryBuilder.Where(c.DeviceType.Eq(*req.DeviceType))
		}
	}

	if req.AccessWay != nil && *req.AccessWay != "" {
		queryBuilder = queryBuilder.Where(q.AccessWay.Eq(*req.AccessWay))
	}

	if req.Label != nil && *req.Label != "" {
		queryBuilder = queryBuilder.Where(q.Label.Like(fmt.Sprintf("%%%s%%", *req.Label)))
	}

	if req.Search != nil && *req.Search != "" {
		queryBuilder = queryBuilder.
			Where(query.Device.Where(
				// q.Name.Like(fmt.Sprintf("%%%s%%", *req.Search)),
				q.Name.Lower().Like(fmt.Sprintf("%%%s%%", strings.ToLower(*req.Search))),
			).Or(
				// q.DeviceNumber.Like(fmt.Sprintf("%%%s%%", *req.Search)),
				q.DeviceNumber.Lower().Like(fmt.Sprintf("%%%s%%", strings.ToLower(*req.Search))),
			),
			)
	}

	if req.Name != nil && *req.Name != "" {
		queryBuilder = queryBuilder.Where(q.Name.Like(fmt.Sprintf("%%%s%%", *req.Name)))
	}

	if req.CurrentVersion != nil && *req.CurrentVersion != "" {
		queryBuilder = queryBuilder.Where(q.CurrentVersion.Eq(*req.CurrentVersion))
	}

	if req.BatchNumber != nil && *req.BatchNumber != "" {
		queryBuilder = queryBuilder.Where(q.BatchNumber.Like(fmt.Sprintf("%%%s%%", *req.BatchNumber)))
	}

	if req.DeviceNumber != nil && *req.DeviceNumber != "" {
		queryBuilder = queryBuilder.Where(q.DeviceNumber.Like(fmt.Sprintf("%%%s%%", *req.DeviceNumber)))
	}
	if req.DeviceConfigId != nil && *req.DeviceConfigId != "" {
		queryBuilder = queryBuilder.Where(q.DeviceConfigID.Eq(*req.DeviceConfigId))
	}

	if req.IsOnline != nil {
		if *req.IsOnline == int(1) {
			queryBuilder = queryBuilder.Where(q.IsOnline.Eq(1))
		} else if *req.IsOnline == int(0) {
			queryBuilder = queryBuilder.Where(q.IsOnline.Neq(1))
		} else {
			return count, deviceList, fmt.Errorf("is_online param error")
		}
	}
	queryBuilder = queryBuilder.LeftJoin(c, c.ID.EqCol(q.DeviceConfigID))
	queryBuilder = queryBuilder.LeftJoin(lda, lda.DeviceID.EqCol(q.ID))

	if req.WarnStatus != nil && *req.WarnStatus != "" {

		if *req.WarnStatus == "N" {
			queryBuilder = queryBuilder.Where(lda.AlarmStatus.Eq("N")).Or(lda.AlarmStatus.IsNull())
		} else {
			queryBuilder = queryBuilder.Where(lda.AlarmStatus.Neq("N"))
		}
	}

	count, err := queryBuilder.Count()
	if err != nil {
		logrus.Error(err)
		return count, deviceList, err
	}

	if req.Page != 0 && req.PageSize != 0 {
		queryBuilder = queryBuilder.Limit(req.PageSize)
		queryBuilder = queryBuilder.Offset((req.Page - 1) * req.PageSize)
	}

	t := query.TelemetryCurrentData
	t2 := query.TelemetryCurrentData.As("t2")
	// q.ID, q.DeviceNumber, q.Name, q.DeviceConfigID, q.ActivateFlag, q.ActivateAt, q.BatchNumber
	err = queryBuilder.Select(lda.AlarmStatus.As("warn_status"), q.ID, q.DeviceNumber, q.Name, q.DeviceConfigID, q.ActivateFlag, q.ActivateAt, q.BatchNumber, q.Location, q.CurrentVersion, q.CreatedAt, q.IsOnline, q.AccessWay, c.ProtocolType, c.DeviceType, c.Name.As("DeviceConfigName"), t2.T).
		LeftJoin(t.Select(t.T.Max().As("ts"), t.DeviceID).Group(t.DeviceID).As("t2"), t2.DeviceID.EqCol(q.ID)).
		Order(q.CreatedAt.Desc()).
		Scan(&deviceList)
	if err != nil {
		logrus.Error(err)
		return count, deviceList, err
	}
	return count, deviceList, err
}

func GetDevicePreRegisterListByPage(req *model.GetDevicePreRegisterListByPageReq, tenant_id string) (int64, []model.GetDevicePreRegisterListByPageRsp, error) {
	q := query.Device
	var count int64
	deviceList := []model.GetDevicePreRegisterListByPageRsp{}
	queryBuilder := q.WithContext(context.Background())

	queryBuilder = queryBuilder.Where(q.TenantID.Eq(tenant_id))

	if req.ActivateFlag != nil && *req.ActivateFlag != "" {
		queryBuilder = queryBuilder.Where(q.ActivateFlag.Eq(*req.ActivateFlag))
	}

	if req.IsEnabled != nil && *req.IsEnabled != "" {
		queryBuilder = queryBuilder.Where(q.IsEnabled.Eq(*req.IsEnabled))
	}

	if req.ProductID != "" {
		queryBuilder = queryBuilder.Where(q.ProductID.Eq(req.ProductID))
	}

	if req.DeviceConfigID != nil && *req.DeviceConfigID != "" {
		queryBuilder = queryBuilder.Where(q.DeviceConfigID.Eq(*req.DeviceConfigID))
	}

	if req.BatchNumber != nil && *req.BatchNumber != "" {
		queryBuilder = queryBuilder.Where(q.BatchNumber.Like(fmt.Sprintf("%%%s%%", *req.BatchNumber)))
	}

	if req.DeviceNumber != nil && *req.DeviceNumber != "" {
		queryBuilder = queryBuilder.Where(q.DeviceNumber.Like(fmt.Sprintf("%%%s%%", *req.DeviceNumber)))
	}

	if req.Name != nil && *req.Name != "" {
		queryBuilder = queryBuilder.Where(q.Name.Like(fmt.Sprintf("%%%s%%", *req.Name)))
	}

	count, err := queryBuilder.Count()
	if err != nil {
		logrus.Error(err)
		return count, deviceList, err
	}

	if req.Page != 0 && req.PageSize != 0 {
		queryBuilder = queryBuilder.Limit(req.PageSize)
		queryBuilder = queryBuilder.Offset((req.Page - 1) * req.PageSize)
	}

	err = queryBuilder.Select(
		q.ID, q.Name, q.DeviceNumber, q.ActivateFlag, q.ActivateAt, q.BatchNumber, q.CurrentVersion).
		Order(q.CreatedAt.Desc()).
		Scan(&deviceList)
	if err != nil {
		logrus.Error(err)
		return count, deviceList, err
	}
	return count, deviceList, err
}

func GetDevicesCount() int64 {
	count, _ := query.Device.Count()
	return count
}

// GetDeviceCacheById retrieves device information by device ID.
func GetDeviceCacheById(deviceId string) (*model.Device, error) {
	device, err := query.Device.Where(query.Device.ID.Eq(deviceId)).First()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return device, nil
}

type DeviceQuery struct{}

func (DeviceQuery) Count(ctx context.Context) (count int64, err error) {
	count, err = query.Device.Count()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (DeviceQuery) CountByTenantID(ctx context.Context, TenantID string) (count int64, err error) {
	device := query.Device
	count, err = device.Where(device.TenantID.Eq(TenantID)).Count()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

// GetGatewayUnrelatedDeviceList retrieves a list of devices that are not related to any gateway.
func (DeviceQuery) GetGatewayUnrelatedDeviceList(ctx context.Context, tenantId string) (list []map[string]interface{}, err error) {
	device := query.Device
	deviceConfig := query.DeviceConfig
	// Condition: device - parent device is empty, device configuration is not empty
	// Condition: device_config_id - device type is 3 - child device
	err = device.
		WithContext(ctx).
		Select(device.ID, device.Name, device.DeviceConfigID.As("device_config_id"), deviceConfig.Name.As("device_config_name")).
		Where(device.TenantID.Eq(tenantId)).
		Where(device.DeviceConfigID.IsNotNull()).
		Where(device.ParentID.IsNull()). // Parent device is empty
		LeftJoin(deviceConfig, deviceConfig.ID.EqCol(device.DeviceConfigID)).
		Where(deviceConfig.DeviceType.Eq("3")). // Device type is 3 - child device
		Scan(&list)
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (DeviceQuery) CountByWhere(ctx context.Context, option ...gen.Condition) (count int64, err error) {
	device := query.Device
	count, err = device.Where(option...).Count()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (DeviceQuery) First(ctx context.Context, option ...gen.Condition) (info *model.Device, err error) {
	info, err = query.Device.WithContext(ctx).Where(option...).First()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (DeviceQuery) Find(ctx context.Context, option ...gen.Condition) (list []*model.Device, err error) {
	list, err = query.Device.WithContext(ctx).Where(option...).Find()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

// Get device dropdown list
// Return device ID, device name, device configuration ID, device configuration name
func (DeviceQuery) GetDeviceSelect(tenantId string, deviceName string, bindConfig int) (list []map[string]interface{}, err error) {
	device := query.Device
	deviceConfig := query.DeviceConfig
	query := device.
		WithContext(context.Background()).
		Select(device.ID, device.Name, device.DeviceConfigID.As("device_config_id"), deviceConfig.Name.As("device_config_name")).
		Where(device.TenantID.Eq(tenantId)).
		Where(device.ActivateFlag.Eq("active")). // Device is active
		Where(device.Name.Like(fmt.Sprintf("%%%s%%", deviceName))).
		LeftJoin(deviceConfig, deviceConfig.ID.EqCol(device.DeviceConfigID)).Order(device.CreatedAt.Desc())
	switch bindConfig {
	case 1:
		query = query.Where(device.DeviceConfigID.IsNotNull())
	case 2:
		query = query.Where(device.DeviceConfigID.IsNull())
	}
	err = query.Scan(&list)
	if err != nil {
		logrus.Error(err)
	}
	return
}

// Update specified field.
func (DeviceQuery) Update(ctx context.Context, info *model.Device, option ...field.Expr) error {
	device := query.Device
	_, err := query.Device.WithContext(ctx).Where(device.ID.Eq(info.ID)).Select(option...).UpdateColumns(info)
	if err != nil {
		logrus.Error(ctx, err)
	}
	return err
}

// Update device configuration ID
func (DeviceQuery) ChangeDeviceConfig(deviceID string, deviceConfigID *string) error {
	device := query.Device
	info, err := device.Where(device.ID.Eq(deviceID)).Update(device.DeviceConfigID, deviceConfigID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if info.RowsAffected == 0 {
		return fmt.Errorf("update device config failed, no rows affected")
	}
	return err
}

func (DeviceQuery) GetSubList(ctx context.Context, parent_id string, pageSize, page int64, tenantID string) ([]model.GetSubListResp, int64, error) {
	var (
		q     = query.Device
		count int64
		resp  []model.GetSubListResp
	)
	query := q.WithContext(ctx).Where(q.ParentID.Eq(parent_id), q.TenantID.Eq(tenantID))
	count, err := query.Count()
	if err != nil {
		return resp, count, err
	}
	err = query.Offset(int(page-1) * int(pageSize)).Limit(int(pageSize)).Order(q.CreatedAt.Desc()).Scan(&resp)
	if err != nil {
		return resp, count, err
	}
	return resp, count, nil
}

// GetSubDeviceListByParentID retrieves a list of sub-devices by parent device ID.
func GetSubDeviceListByParentID(parentId string) ([]*model.Device, error) {
	device := query.Device
	list, err := device.Where(device.ParentID.Eq(parentId)).Find()
	if err != nil {
		logrus.Error(err)
	}
	return list, err
}

func GetDeviceTemplateChartSelect(tenantId string) (any, error) {
	data := []map[string]interface{}{}
	d := query.Device
	dc := query.DeviceConfig
	dm := query.DeviceTemplate
	err := d.LeftJoin(dc, dc.ID.EqCol(d.DeviceConfigID)).
		LeftJoin(dm, dm.ID.EqCol(dc.DeviceTemplateID)).
		Where(d.TenantID.Eq(tenantId)).
		Where(d.ActivateFlag.Eq("active")).
		Where(d.DeviceConfigID.IsNotNull()).
		Where(dc.DeviceTemplateID.IsNotNull()).
		Where(dm.WebChartConfig.IsNotNull()).
		Select(d.ID.As("device_id"), d.Name.As("device_name"), dm.WebChartConfig).Scan(&data)
	if err != nil {
		logrus.Error(err)
	}
	return data, nil
}

func GetDeviceCurrentStatus(deviceId string) (string, error) {
	data, err := query.Device.Where(query.Device.ID.Eq(deviceId)).First()
	var result string = "OFF-LINE"
	if err != nil {
		return result, err
	} else if err == gorm.ErrRecordNotFound {
		return result, nil
	}
	if data.IsOnline == 1 {
		result = "ON-LINE"
	}
	return result, nil
}

func GetDeviceTemplateIdByDeviceId(deviceId string) (string, error) {
	var result model.DeviceConfig
	query.Device.LeftJoin(query.DeviceConfig, query.Device.DeviceConfigID.EqCol(query.DeviceConfig.ID)).
		Where(query.Device.ID.Eq(deviceId)).Select(query.DeviceConfig.DeviceTemplateID).Scan(&result)
	if result.DeviceTemplateID != nil {
		return *result.DeviceTemplateID, nil
	}
	return "", nil
}

// GetDeviceConfigIdByDeviceId retrieves the device configuration ID by device ID.
func GetDevicesByDeviceConfigID(deviceConfigID string) ([]*model.Device, error) {
	device := query.Device
	list, err := device.Where(device.DeviceConfigID.Eq(deviceConfigID)).Find()
	if err != nil {
		logrus.Error(err)
	}
	return list, err
}

// func GetDevicesBySubDeviceConfigID(deviceConfigID string) ([]*model.Device, error) {
// 	var device = query.Device
// 	list, err := device.Where(device.DeviceConfigID.Eq(deviceConfigID), device.ParentID.IsNotNull()).Find()
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	return list, err
// }

// GetGatewayDevicesBySubDeviceConfigID retrieves the gateway devices by sub-device configuration ID.
func GetGatewayDevicesBySubDeviceConfigID(deviceConfigID string) ([]string, error) {
	device := query.Device
	var deviceIDList []string
	err := device.Where(device.DeviceConfigID.Eq(deviceConfigID), device.ParentID.IsNotNull()).Select(device.ParentID.Distinct()).Scan(&deviceIDList)
	if err != nil {
		logrus.Error(err)
	}
	return deviceIDList, err
}

// GetServiceDeviceList retrieves the list of devices associated with a service access ID.
func GetServiceDeviceList(serviceAccessId string) ([]model.Device, error) {
	var devices []model.Device
	err := query.Device.Where(query.Device.ServiceAccessID.Eq(serviceAccessId)).Scan(&devices)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return devices, err
}

// GetSubDeviceExists
// @description GetSubDeviceExists checks if a sub-device exists under a parent device
func GetSubDeviceExists(deviceId, subAddr string) bool {
	num, err := query.Device.Where(query.Device.ParentID.Eq(deviceId), query.Device.SubDeviceAddr.Eq(subAddr)).Count()
	if num > 0 || err != nil {
		return true
	}
	return false
}

// CheckDeviceNumberExists
// CheckDeviceNumberExists checks if a device number already exists in the database
func CheckDeviceNumberExists(deviceNumber string) (bool, error) {
	count, err := query.Device.Where(query.Device.DeviceNumber.Eq(deviceNumber)).Count()
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	return count > 0, nil
}

// GetDeviceSelector retrieves a list of devices based on the provided request parameters.
func GetDeviceSelector(req model.DeviceSelectorReq, tenantId string) (*model.DeviceSelectorRes, error) {
	device := query.Device

	query := device.WithContext(context.Background())

	if req.HasDeviceConfig != nil {
		if *req.HasDeviceConfig {
			query = query.Where(device.DeviceConfigID.IsNotNull())
		} else {
			query = query.Where(device.DeviceConfigID.IsNull())
		}
	}

	query = query.Where(device.TenantID.Eq(tenantId))

	query = query.Select(device.ID.As("device_id"), device.Name.As("device_name"))

	query = query.Order(device.CreatedAt.Desc())

	count, err := query.Count()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	query = query.Limit(req.PageSize)
	query = query.Offset((req.Page - 1) * req.PageSize)

	var list []*model.DeviceSelectorData
	err = query.Scan(&list)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &model.DeviceSelectorRes{
		Total: count,
		List:  list,
	}, nil
}
