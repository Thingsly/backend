package dal

import (
	"context"
	"errors"
	"time"

	model "github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

const (
	DEVICE_TYPE_1 = "DIRECT_ATTACHED_PROTOCOL"
	DEVICE_TYPE_2 = "GATEWAY_PROTOCOL"
)

func CreateProtocolPluginWithDict(p *model.CreateProtocolPluginReq) (*model.ProtocolPlugin, error) {

	var dictCode string
	switch p.DeviceType {
	case 1:
		dictCode = DEVICE_TYPE_1
	case 2:
		dictCode = DEVICE_TYPE_2
	default:
		return nil, errors.New("deviceType is invalid")
	}
	logrus.Info("dictCode:", dictCode)

	tx, err := StartTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	t := time.Now().UTC()

	var dict = model.SysDict{}
	dictId := uuid.New()
	dict.ID = dictId
	dict.DictCode = dictCode
	dict.DictValue = p.ProtocolType
	dict.CreatedAt = t
	if err := CreateDict(&dict, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	var dictLanguage = model.SysDictLanguage{}
	dictLanguage.ID = uuid.New()
	dictLanguage.DictID = dictId
	dictLanguage.LanguageCode = p.LanguageCode
	dictLanguage.Translation = p.Name

	if err := CreateDictLanguage(&dictLanguage, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	var protocolPlugin = model.ProtocolPlugin{}
	protocolPlugin.ID = uuid.New()
	protocolPlugin.Name = p.Name
	protocolPlugin.DeviceType = p.DeviceType
	protocolPlugin.ProtocolType = p.ProtocolType
	protocolPlugin.AccessAddress = p.AccessAddress
	protocolPlugin.HTTPAddress = p.HTTPAddress
	protocolPlugin.SubTopicPrefix = p.SubTopicPrefix
	protocolPlugin.Description = p.Description
	protocolPlugin.AdditionalInfo = p.AdditionalInfo
	protocolPlugin.CreatedAt = t
	protocolPlugin.UpdateAt = t
	protocolPlugin.Remark = p.Remark
	if err := tx.ProtocolPlugin.Create(&protocolPlugin); err != nil {
		tx.Rollback()
		return nil, err
	} else {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
	}
	return &protocolPlugin, nil
}

func DeleteProtocolPluginWithDict(id string) error {

	tx, err := StartTransaction()
	if err != nil {
		return err
	}

	dictLanguage, err := tx.ProtocolPlugin.Where(tx.ProtocolPlugin.ID.Eq(id)).First()
	if err != nil {
		return err
	}

	var dictCode string
	if dictLanguage.DeviceType == 1 {
		dictCode = DEVICE_TYPE_1
	} else {
		dictCode = DEVICE_TYPE_2
	}

	dict, err := tx.SysDict.Where(tx.SysDict.DictCode.Eq(dictCode), tx.SysDict.DictValue.Eq(dictLanguage.ProtocolType)).First()
	if err != nil {
		return err
	}

	_, err = tx.SysDict.Where(tx.SysDict.DictCode.Eq(dictCode), tx.SysDict.DictValue.Eq(dictLanguage.ProtocolType)).Delete()
	if err != nil {
		return err
	}

	_, err = tx.SysDictLanguage.Where(tx.SysDictLanguage.DictID.Eq(dict.ID)).Delete()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ProtocolPlugin.Where(tx.ProtocolPlugin.ID.Eq(id)).Delete()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func UpdateProtocolPluginWithDict(p *model.UpdateProtocolPluginReq) error {

	tx, err := StartTransaction()
	if err != nil {
		return err
	}

	oldProtocolPlugin, err := tx.ProtocolPlugin.Where(tx.ProtocolPlugin.ID.Eq(p.Id)).First()
	if err != nil {
		return err
	}

	var pp = model.ProtocolPlugin{}
	pp.ID = p.Id
	pp.Name = p.Name
	pp.DeviceType = p.DeviceType
	pp.ProtocolType = p.ProtocolType
	pp.AccessAddress = p.AccessAddress
	pp.HTTPAddress = p.HTTPAddress
	pp.SubTopicPrefix = p.SubTopicPrefix
	pp.Description = p.Description
	pp.AdditionalInfo = p.AdditionalInfo
	pp.UpdateAt = time.Now().UTC()
	pp.Remark = p.Remark

	err = tx.ProtocolPlugin.Save(&pp)
	if err != nil {
		return err
	}

	// if oldProtocolPlugin.DeviceType != p.DeviceType || oldProtocolPlugin.ProtocolType != p.ProtocolType {
	var dictCode string
	if oldProtocolPlugin.DeviceType == int16(1) {
		dictCode = DEVICE_TYPE_1
	} else {
		dictCode = DEVICE_TYPE_2
	}

	dict, err := tx.SysDict.Where(tx.SysDict.DictCode.Eq(dictCode), tx.SysDict.DictValue.Eq(oldProtocolPlugin.ProtocolType)).First()

	if err != nil {
		tx.Rollback()
		return err
	}

	var newDict = model.SysDict{}
	var newDictCode string
	if p.DeviceType == int16(1) {
		newDictCode = DEVICE_TYPE_1
	} else {
		newDictCode = DEVICE_TYPE_2
	}
	newDict.ID = dict.ID
	newDict.DictCode = newDictCode
	newDict.DictValue = p.ProtocolType

	err = tx.SysDict.Save(&newDict)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldDictLanguage, err := tx.SysDictLanguage.Where(tx.SysDictLanguage.DictID.Eq(dict.ID)).First()
	if err != nil {
		tx.Rollback()
		return err
	}
	if oldProtocolPlugin.Name != p.Name || oldDictLanguage.LanguageCode != p.LanguageCode {
		var newDictLanguage = model.SysDictLanguage{}
		newDictLanguage.LanguageCode = p.LanguageCode
		newDictLanguage.Translation = p.Name
		_, err = tx.SysDictLanguage.Where(tx.SysDictLanguage.DictID.Eq(dict.ID)).Updates(newDictLanguage)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// }

	tx.Commit()

	return nil
}

func GetProtocolPluginListByPage(p *model.GetProtocolPluginListByPageReq) (int64, interface{}, error) {
	q := query.ProtocolPlugin
	var count int64
	queryBuilder := q.WithContext(context.Background())
	count, err := queryBuilder.Count()
	if err != nil {
		logrus.Error(err)
		return count, nil, err
	}

	if p.Page != 0 && p.PageSize != 0 {
		queryBuilder = queryBuilder.Limit(p.PageSize)
		queryBuilder = queryBuilder.Offset((p.Page - 1) * p.PageSize)
	}

	protocolPluginList, err := queryBuilder.Select().Order(q.CreatedAt.Desc()).Find()
	if err != nil {
		logrus.Error(err)
		return count, protocolPluginList, err
	}
	return count, protocolPluginList, err
}

func GetProtocolPluginByDeviceTypeAndProtocolType(deviceType int16, protocolType string) (*model.ProtocolPlugin, error) {
	q := query.ProtocolPlugin
	protocolPlugin, err := q.Where(q.DeviceType.Eq(deviceType), q.ProtocolType.Eq(protocolType)).First()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return protocolPlugin, nil
}

func GetProtocolPluginByDeviceConfigID(deviceConfigID string) (*model.ProtocolPlugin, error) {
	q := query.DeviceConfig
	deviceConfig, err := q.Where(q.ID.Eq(deviceConfigID)).First()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if deviceConfig.ProtocolType == nil {
		return nil, errors.New("protocolType is nil")
	}

	var deviceType int16
	switch deviceConfig.DeviceType {
	case "1":
		deviceType = 1
	case "2", "3":
		deviceType = 2
	default:
		return nil, errors.New("deviceType is invalid")
	}

	if *deviceConfig.ProtocolType == "MQTT" {
		return nil, nil
	}
	protocolPlugin, err := GetProtocolPluginByDeviceTypeAndProtocolType(deviceType, *deviceConfig.ProtocolType)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return protocolPlugin, nil
}
