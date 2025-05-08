package dal

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/query"
	"github.com/Thingsly/backend/pkg/errcode"

	"github.com/sirupsen/logrus"
)

func DeleteServicePlugin(id string) error {

	tx, err := StartTransaction()
	if err != nil {
		Rollback(tx)
		return err
	}
	serviceAccess := tx.ServiceAccess

	_, err = serviceAccess.Where(serviceAccess.ServicePluginID.Eq(id)).Delete()
	if err != nil {
		Rollback(tx)
		return err
	}

	servicePlugin := tx.ServicePlugin

	_, err = servicePlugin.Where(query.ServicePlugin.ID.Eq(id)).Delete()
	if err != nil {
		Rollback(tx)
		return err
	}

	err = Commit(tx)
	return err
}

func UpdateServicePlugin(id string, updates map[string]interface{}) error {
	q := query.ServicePlugin
	queryBuilder := q.WithContext(context.Background())
	_, err := queryBuilder.Where(query.ServicePlugin.ID.Eq(id)).Updates(updates)
	return err
}

func GetServicePluginListByPage(req *model.GetServicePluginByPageReq) (int64, []map[string]interface{}, error) {
	var count int64
	servicePlugins := make([]map[string]interface{}, 0)

	q := query.ServicePlugin
	queryBuilder := q.WithContext(context.Background())
	if req.ServiceType != 0 {
		queryBuilder = queryBuilder.Where(q.ServiceType.Eq(req.ServiceType))
	}

	count, err := queryBuilder.Count()
	if err != nil {
		logrus.Error(err)
		return count, servicePlugins, err
	}
	if req.Page != 0 && req.PageSize != 0 {
		queryBuilder = queryBuilder.Limit(req.PageSize)
		queryBuilder = queryBuilder.Offset((req.Page - 1) * req.PageSize)
	}
	timeNow := time.Now().UTC()
	err = queryBuilder.Select().Order(q.CreateAt.Desc()).Scan(&servicePlugins)
	if err != nil {
		logrus.Error(err)
		return count, servicePlugins, err
	}

	for i := range servicePlugins {
		lastActiveTime, ok := servicePlugins[i]["last_active_time"].(time.Time)
		if !ok {
			logrus.Warn("LastActiveTime is not of type time.Time for plugin ", i)
			servicePlugins[i]["service_heartbeat"] = 2
			continue
		}

		if timeNow.Sub(lastActiveTime) > time.Minute {
			servicePlugins[i]["service_heartbeat"] = 2
		} else {
			servicePlugins[i]["service_heartbeat"] = 1
		}
	}
	return count, servicePlugins, err
}

func GetServicePlugin(id string) (interface{}, error) {
	var servicePlugin *model.ServicePlugin

	q := query.ServicePlugin
	queryBuilder := q.WithContext(context.Background())

	queryBuilder = queryBuilder.Where(q.ID.Eq(id))

	servicePlugin, err := queryBuilder.First()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return servicePlugin, err
}

func GetServicePluginByID(id string) (*model.ServicePlugin, error) {

	q := query.ServicePlugin
	queryBuilder := q.WithContext(context.Background())
	servicePlugin, err := queryBuilder.Where(q.ID.Eq(id)).Select().First()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return servicePlugin, nil
}

func GetServicePluginHttpAddressByID(id string) (*model.ServicePlugin, string, error) {
	servicePlugin, err := GetServicePluginByID(id)
	if err != nil {
		return nil, "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	if servicePlugin.ServiceConfig == nil || *servicePlugin.ServiceConfig == "" {
		return nil, "", errcode.New(200065)
	}

	var serviceAccessConfig model.ServiceAccessConfig
	err = json.Unmarshal([]byte(*servicePlugin.ServiceConfig), &serviceAccessConfig)
	if err != nil {
		return nil, "", errcode.New(200066)
	}

	if serviceAccessConfig.HttpAddress == "" {
		return nil, "", errcode.New(200067)
	}
	return servicePlugin, serviceAccessConfig.HttpAddress, nil
}

func GetServicePluginByServiceIdentifier(serviceIdentifier string) (*model.ServicePlugin, error) {
	if serviceIdentifier == "MQTT" {
		return &model.ServicePlugin{
			Name:              "MQTT",
			ServiceType:       1,
			ServiceConfig:     nil,
			ServiceIdentifier: "MQTT",
		}, nil
	}

	q := query.ServicePlugin
	queryBuilder := q.WithContext(context.Background())
	servicePlugin, err := queryBuilder.Where(q.ServiceIdentifier.Eq(serviceIdentifier)).Select().First()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return servicePlugin, nil
}

func GetServicePluginByDeviceConfigID(deviceConfigID string) (*model.ServicePlugin, error) {

	deviceConfig, err := GetDeviceConfigByID(deviceConfigID)
	if err != nil {
		return nil, err
	}

	return GetServicePluginByServiceIdentifier(*deviceConfig.ProtocolType)
}

func GetServicePluginSubTopicPrefixByDeviceConfigID(deviceConfigID string) (string, error) {
	servicePlugin, err := GetServicePluginByDeviceConfigID(deviceConfigID)
	if err != nil {
		logrus.Error("failed to get service plugin by device config id: ", err)
		return "", err
	}
	var subTopicPrefix string
	if servicePlugin.ServiceType == int32(1) {
		var protocolAccessConfig model.ProtocolAccessConfig
		if servicePlugin.ServiceConfig == nil {
			err = errors.New("service config is empty")
			return "", err
		}
		err = json.Unmarshal([]byte(*servicePlugin.ServiceConfig), &protocolAccessConfig)
		if err != nil {
			logrus.Error("failed to unmarshal service config: ", err)
			return "", err
		}
		if protocolAccessConfig.SubTopicPrefix != "" {
			subTopicPrefix = protocolAccessConfig.SubTopicPrefix
		}
	} else if servicePlugin.ServiceType == int32(2) {
		var serviceAccessConfig model.ServiceAccessConfig
		if servicePlugin.ServiceConfig == nil {
			err = errors.New("service config is empty")
			return "", err
		}
		err = json.Unmarshal([]byte(*servicePlugin.ServiceConfig), &serviceAccessConfig)
		if err != nil {
			logrus.Error("failed to unmarshal service config: ", err)
			return "", err
		}
		if serviceAccessConfig.SubTopicPrefix != "" {
			subTopicPrefix = serviceAccessConfig.SubTopicPrefix
		}

	}
	return subTopicPrefix, nil
}

func UpdateServicePluginHeartbeat(serviceIdentifier string) error {
	q := query.ServicePlugin
	queryBuilder := q.WithContext(context.Background())

	t := time.Now().UTC()
	info, err := queryBuilder.Where(q.ServiceIdentifier.Eq(serviceIdentifier)).Update(q.LastActiveTime, t)
	if err != nil {
		logrus.Error(err)
	}
	if info.RowsAffected == 0 {
		return errors.New("service plugin not found")
	}
	return err
}

// GetServiceSelectList
func GetServiceSelectList() ([]model.ServicePlugin, error) {
	q := query.ServicePlugin
	queryBuilder := q.WithContext(context.Background())
	var servicePlugins []model.ServicePlugin
	err := queryBuilder.Select().Scan(&servicePlugins)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return servicePlugins, nil
}
