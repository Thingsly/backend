package dal

import (
	model "github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"

	"github.com/go-basic/uuid"
	"gorm.io/gorm"
)

func GetAttributeDataList(deviceId string) ([]*model.AttributeData, error) {
	data, err := query.AttributeData.
		Where(query.AttributeData.DeviceID.Eq(deviceId)).Find()
	if err != nil {
		return nil, err
	}
	return data, nil
}

/*
select ad.*,dma.data_name from attribute_datas ad
left join devices on ad.device_id = devices.id  left join  device_configs dc on devices.device_config_id = dc.id
left join device_templates dt on dt.id = dc.device_template_id
left join device_model_attributes dma on dt.id = dma.device_template_id and ad.key = dma.data_identifier
where devices.id = 'ca33926c-5ee5-3e9f-147e-94e188fde65b'
*/

// Retrieve the list of device attribute data based on device ID and join to get data names as in the SQL above
func GetAttributeDataListWithDeviceName(deviceId string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	err := query.AttributeData.
		Select(query.AttributeData.ALL, query.DeviceModelAttribute.DataName, query.DeviceModelAttribute.Unit, query.DeviceModelAttribute.ReadWriteFlag, query.DeviceModelAttribute.DataType, query.DeviceModelAttribute.AdditionalInfo.As("enum")).
		LeftJoin(query.Device, query.AttributeData.DeviceID.EqCol(query.Device.ID)).
		LeftJoin(query.DeviceConfig, query.Device.DeviceConfigID.EqCol(query.DeviceConfig.ID)).
		LeftJoin(query.DeviceTemplate, query.DeviceConfig.DeviceTemplateID.EqCol(query.DeviceTemplate.ID)).
		LeftJoin(query.DeviceModelAttribute, query.DeviceTemplate.ID.EqCol(query.DeviceModelAttribute.DeviceTemplateID), query.AttributeData.Key.EqCol(query.DeviceModelAttribute.DataIdentifier)).
		Where(query.Device.ID.Eq(deviceId)).Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DeleteAttributeData(id string) error {
	_, err := query.AttributeData.
		Where(query.AttributeData.ID.Eq(id)).
		Delete()
	return err
}

func CreateAttributeData(data *model.AttributeData) error {
	return query.AttributeData.Create(data)
}

// Update device attribute data. If the data does not exist, generate a UUID to create a new entry
func UpdateAttributeData(data *model.AttributeData) (*model.AttributeData, error) {
	// Based on the data type of the new data, directly set other type fields to null
	if data.StringV != nil {
		data.NumberV = nil
		data.BoolV = nil
	} else if data.NumberV != nil {
		data.StringV = nil
		data.BoolV = nil
	} else if data.BoolV != nil {
		data.StringV = nil
		data.NumberV = nil
	}

	// Create an update map that includes null values to ensure null fields are also updated in the database
	updateMap := map[string]interface{}{
		"bool_v":   data.BoolV,
		"number_v": data.NumberV,
		"string_v": data.StringV,
		"ts":       data.T,
	}

	// Try to update the existing record
	result, err := query.AttributeData.Where(
		query.AttributeData.DeviceID.Eq(data.DeviceID),
		query.AttributeData.TenantID.Eq(*data.TenantID),
		query.AttributeData.Key.Eq(data.Key),
	).Updates(updateMap)
	if err != nil {
		return nil, err
	} else if result.RowsAffected == 0 {
		data.ID = uuid.New()
		err = query.AttributeData.Create(data)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return data, nil
}

// Get the latest value of a single device metric. If the data does not exist, return nil
func GetAttributeOneKeys(deviceId string, keys string) (interface{}, error) {
	data, err := query.AttributeData.Where(query.AttributeData.DeviceID.Eq(deviceId), query.AttributeData.Key.Eq(keys)).Order(query.AttributeData.T.Desc()).First()
	var result interface{}
	if err != nil {
		return result, err
	} else if err == gorm.ErrRecordNotFound {
		return result, nil
	}
	if data.BoolV != nil {
		//result = fmt.Sprintf("%t", *data.BoolV)
		result = *data.BoolV
	}
	if data.NumberV != nil {
		//result = fmt.Sprintf("%d", data.NumberV)
		result = *data.NumberV
	}
	if data.StringV != nil {
		result = *data.StringV
	}
	return result, nil
}

// Get the latest value of a single device metric. If the data does not exist, return nil.
func GetAttributeOneKeysByDeviceId(deviceId string, keys string) (*model.AttributeData, error) {
	data, err := query.AttributeData.Where(query.AttributeData.DeviceID.Eq(deviceId), query.AttributeData.Key.Eq(keys)).Order(query.AttributeData.T.Desc()).First()
	if err != nil {
		return &model.AttributeData{}, err
	}
	return data, nil
}

// Delete all attribute data for a specific device ID
func DeleteAttributeDataByDeviceId(deviceId string, tx *query.QueryTx) error {
	_, err := tx.AttributeData.Where(query.AttributeData.DeviceID.Eq(deviceId)).Delete()
	return err
}
