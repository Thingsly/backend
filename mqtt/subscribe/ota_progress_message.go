package subscribe

import (
	"encoding/json"
	"strconv"

	initialize "github.com/Thingsly/backend/initialize"
	"github.com/Thingsly/backend/internal/query"

	"github.com/sirupsen/logrus"
)

type DeviceProgressMsg struct {
	UpgradeProgress interface{} `json:"step,omitempty" alias:"Progress"` // Progress of the upgrade step (e.g., current step in the process)
	StatusDetail    string      `json:"desc" alias:"Description"`        // Detailed description of the current status
	Module          string      `json:"module,omitempty" alias:"Module"` // The module related to the upgrade (optional)
	//UpgradeStatus    string      `json:"upgrade_status,omitempty"`         // Upgrade status (optional, commented out)
	//StatusUpdateTime string      `json:"status_update_time" alias:"Upgrade Update Time"`  // Time when the upgrade status was last updated (optional, commented out)
}

func OtaUpgrade(payload []byte, _ string) {

	logrus.Debug("ota progress message:", string(payload))
	progressMsgPayload, err := verifyPayload(payload)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	device, err := initialize.GetDeviceCacheById(progressMsgPayload.DeviceId)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	var progressMsg DeviceProgressMsg
	err = json.Unmarshal(progressMsgPayload.Values, &progressMsg)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	switch progressMsg.UpgradeProgress.(type) {
	case float64:

		progressMsg.UpgradeProgress = strconv.FormatInt(int64(progressMsg.UpgradeProgress.(float64)), 10)

	case string:

	default:
		logrus.Error("Unsupported data type")
		return
	}

	otaTaskDetail, err := query.OtaUpgradeTaskDetail.
		Where(query.OtaUpgradeTaskDetail.DeviceID.Eq(device.ID),
			query.OtaUpgradeTaskDetail.Status.In(2, 3),
		).First()
	if err != nil && otaTaskDetail != nil {
		logrus.Errorf("Upgrade task not found")
		return
	}

	intProgress, err := strconv.Atoi(progressMsg.UpgradeProgress.(string))
	if err != nil {
		desc := progressMsg.UpgradeProgress.(string) + " " + progressMsg.StatusDetail
		otaTaskDetail.StatusDescription = &desc
	}

	switch {
	case intProgress == -1:
		desc := "Error code -1, upgrade failed " + progressMsg.StatusDetail
		otaTaskDetail.Status = 5
		otaTaskDetail.StatusDescription = &desc
	case intProgress == -2:
		desc := "Error code -2, download failed " + progressMsg.StatusDetail
		otaTaskDetail.Status = 5
		otaTaskDetail.StatusDescription = &desc
	case intProgress == -3:
		desc := "Error code -3, verification failed " + progressMsg.StatusDetail
		otaTaskDetail.Status = 5
		otaTaskDetail.StatusDescription = &desc
	case intProgress == -4:
		desc := "Error code -4, write failed " + progressMsg.StatusDetail
		otaTaskDetail.Status = 5
		otaTaskDetail.StatusDescription = &desc
	case intProgress >= 1 && intProgress < 100:
		otaTaskDetail.Status = 3
		otaTaskDetail.StatusDescription = &progressMsg.StatusDetail
	case intProgress == 100:
		otaTaskDetail.Status = 4
		otaTaskDetail.StatusDescription = &progressMsg.StatusDetail
	default:
		logrus.Error("Invalid data format")
		return
	}

	_, err = query.OtaUpgradeTaskDetail.Where(query.OtaUpgradeTaskDetail.ID.Eq(otaTaskDetail.ID)).Updates(otaTaskDetail)
	if err != nil {
		logrus.Error(err)
		return
	}
}
