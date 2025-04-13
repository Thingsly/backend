package test

import (
	"fmt"
	"testing"

	"github.com/HustIoTPlatform/backend/initialize"

	"github.com/sirupsen/logrus"
)

func init() {
	initialize.ViperInit("../../configs/conf-localdev.yml")
	initialize.LogInIt()
	initialize.RedisInit()
	alarmCache = initialize.NewAlarmCache()

}

var (
	alarmCache          *initialize.AlarmCache
	group_id            = "group_id_1234"
	scene_automation_id = "scene_automation_id1234"
	device_ids          = []string{"device_id123", "device_id456"}
	contents            = []string{"Temperature greater than 30", "Humidity greater than 27"}
)

func TestSetDevice(t *testing.T) {
	logrus.Debug("Unit test execution started:")
	err := alarmCache.SetDevice(group_id, scene_automation_id, device_ids, contents)
	if err != nil {
		t.Error("Failed to set alarm cache", err)
	}

	res1, err := alarmCache.GetByGroupId(group_id)
	if err != nil {
		t.Error("Failed to query alarm cache by group ID", err)
	}
	fmt.Printf("res:%#v", res1)

	res2, err := alarmCache.GetBySceneAutomationId(scene_automation_id)

	if err != nil {
		t.Error("Failed to query alarm cache by scene automation ID", err)
	}
	fmt.Printf("res:%#v", res2)
}

func TestDeleteByGroupId(t *testing.T) {
	fmt.Println("Testing cache deletion...")

	err := alarmCache.DeleteBygroupId(group_id)
	if err != nil {
		t.Error("Failed to delete alarm cache", err)
	}
}
