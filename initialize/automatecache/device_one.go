package automatecache

import "github.com/Thingsly/backend/internal/model"

type OneDeviceCache struct{}

func NewOneDeviceCache() *OneDeviceCache {
	return &OneDeviceCache{}
}

func (*OneDeviceCache) GetAutomateCacheKeyPrefix() string {
	return "one"
}

func (*OneDeviceCache) GetDeviceTriggerConditionType() string {
	return model.DEVICE_TRIGGER_CONDITION_TYPE_ONE
}
