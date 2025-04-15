// /initialize/croninit/cron.go
package croninit

import (
	"time"

	"github.com/Thingsly/backend/internal/service"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var (
	c = cron.New()
)

// Define task initialization
func CronInit() {

	// Initialize device statistics scheduled task
	InitDeviceStatsCron(c)

	// Define one-time task - executed every 5 seconds
	c.AddFunc("*/5 * * * * *", func() {
		logrus.Debug("Automation one-time task starts:")
		service.GroupApp.OnceTaskExecute()
	})

	// Define periodic task - executed every 5 seconds
	c.AddFunc("*/5 * * * * *", func() {
		logrus.Debug("Automation periodic task starts:")
		service.GroupApp.PeriodicTaskExecute()
	})

	// Data cleanup task, executed at 2 AM every day
	c.AddFunc("0 2 * * *", func() {
		logrus.Debug("System data cleanup task starts:")
		service.GroupApp.CleanSystemDataByCron()
	})

	// Run script task, executed at 1 AM every day
	c.AddFunc("0 1 * * *", func() {
		service.GroupApp.RunScript()
	})

	// Message push management cleanup task, executed at 2 AM every day
	err := c.AddFunc("2 0 * * * *", func() {
		logrus.Debug("Task starts:", time.Now())
		service.GroupApp.MessagePush.MessagePushMangeClear()
	})
	if err != nil {
		logrus.Error("Message cleanup scheduled task failed to start")
	}
	c.Start()
}
