package common

import (
	"context"
	"time"

	"github.com/Thingsly/backend/pkg/global"

	"github.com/sirupsen/logrus"
)

func AcquireLock(lockKey string, expiration time.Duration) bool {

	ok, err := global.REDIS.SetNX(context.Background(), lockKey, true, expiration).Result()
	if err != nil {
		return false
	}
	return ok
}

func ReleaseLock(lockKey string) {

	err := global.REDIS.Del(context.Background(), lockKey).Err()
	if err != nil {
		logrus.Error("Error releasing lock:", err)
	}
}
