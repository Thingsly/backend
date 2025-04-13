package initialize

import "github.com/robfig/cron"

func CronInit() {
	c := cron.New()

	c.AddFunc("*/5 * * * * *", func() {

	})

	c.AddFunc("*/5 * * * * *", func() {

	})
}
