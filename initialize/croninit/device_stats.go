// initialize/croninit/device_stats.go
package croninit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Thingsly/backend/internal/dal"
	"github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/global"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type DeviceStatsData struct {
	DeviceTotal int64     `json:"device_total"`
	DeviceOn    int64     `json:"device_on"`
	Timestamp   time.Time `json:"timestamp"`
}

const (
	deviceStatsKeyPattern = "device_stats:%s:%s"

	dataRetentionPeriod = 48 * time.Hour
)

// InitDeviceStatsCron initializes the device statistics scheduled task
func InitDeviceStatsCron(c *cron.Cron) {
	// Executes every hour on the hour
	c.AddFunc("0 0 * * * *", func() {
		collectDeviceStats()
	})
}

// collectDeviceStats collects device statistics data
func collectDeviceStats() {
	ctx := context.Background()
	logrus.Info("Starting device status statistics task")

	// Get the list of all tenant IDs
	userList, err := dal.UserVo{}.GetTenantAdminList()
	if err != nil {
		logrus.Errorf("Failed to get tenant list: %v", err)
		return
	}

	currentTime := time.Now()
	dateStr := currentTime.Format("2006-01-02")

	// Iterate through the tenants and collect device statistics
	for _, user := range userList {
		// Get the device statistics for the tenant
		deviceStats, err := service.GroupApp.Board.GetDeviceByTenantID(ctx, *user.TenantID)
		if err != nil {
			logrus.Errorf("Failed to get device statistics for tenant %s: %v", *user.TenantID, err)
			continue
		}

		// Build the statistics data
		statsData := DeviceStatsData{
			DeviceTotal: deviceStats.DeviceTotal,
			DeviceOn:    deviceStats.DeviceOn,
			Timestamp:   currentTime,
		}

		// Serialize the data
		statsJSON, err := json.Marshal(statsData)
		if err != nil {
			logrus.Errorf("Failed to serialize statistics data: %v", err)
			continue
		}

		// Construct the Redis key
		key := fmt.Sprintf(deviceStatsKeyPattern, *user.TenantID, dateStr)

		// Store the data in the Redis List
		err = global.REDIS.RPush(ctx, key, string(statsJSON)).Err()
		if err != nil {
			logrus.Errorf("Failed to store statistics data in Redis: %v", err)
			continue
		}

		// Set the expiration time for the key
		err = global.REDIS.Expire(ctx, key, dataRetentionPeriod).Err()
		if err != nil {
			logrus.Errorf("Failed to set expiration time for Redis key: %v", err)
		}
	}

	logrus.Info("Device status statistics task completed")
}
