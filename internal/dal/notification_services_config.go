package dal

import (
	"errors"

	model "github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"

	"gorm.io/gorm"
)

func GetNotificationServicesConfigByType(noticeType string) (*model.NotificationServicesConfig, error) {
	data, err := query.NotificationServicesConfig.Where(query.NotificationServicesConfig.NoticeType.Eq(noticeType)).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, nil
		}
		return nil, err
	}
	return data, nil
}

func SaveNotificationServicesConfig(data *model.NotificationServicesConfig) (*model.NotificationServicesConfig, error) {
	err := query.NotificationServicesConfig.Save(data)
	if err != nil {

		return nil, err
	}
	return data, nil
}
