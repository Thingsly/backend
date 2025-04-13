package dal

import (
	"context"

	"github.com/HustIoTPlatform/backend/internal/model"
	"github.com/HustIoTPlatform/backend/internal/query"

	"github.com/sirupsen/logrus"
	"gorm.io/gen"
)

type DeviceModelCommandsQuery struct {
}

func (DeviceModelCommandsQuery) First(ctx context.Context, option ...gen.Condition) (info *model.DeviceModelCommand, err error) {
	info, err = query.DeviceModelCommand.WithContext(ctx).Where(option...).First()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}

func (DeviceModelCommandsQuery) Find(ctx context.Context, option ...gen.Condition) (list []*model.DeviceModelCommand, err error) {
	list, err = query.DeviceModelCommand.WithContext(ctx).Where(option...).Find()
	if err != nil {
		logrus.Error(ctx, err)
	}
	return
}
