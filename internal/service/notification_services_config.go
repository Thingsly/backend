package service

import (
	"encoding/json"
	"strings"
	"time"

	dal "github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/pkg/errcode"
	utils "github.com/Thingsly/backend/pkg/utils"
	"github.com/Thingsly/backend/third_party/others/http_client"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type NotificationServicesConfig struct{}

func (*NotificationServicesConfig) SaveNotificationServicesConfig(req *model.SaveNotificationServicesConfigReq) (*model.NotificationServicesConfig, error) {

	c, err := dal.GetNotificationServicesConfigByType(req.NoticeType)
	if err != nil {
		return nil, err
	}

	config := model.NotificationServicesConfig{}

	var strconf []byte
	switch req.NoticeType {
	case model.NoticeType_Email:
		strconf, err = json.Marshal(req.EMailConfig)
		if err != nil {
			return nil, err
		}
	case model.NoticeType_SME_CODE:
		strconf, err = json.Marshal(req.SMEConfig)
		if err != nil {
			return nil, err
		}
	}

	if c == nil {
		config.ID = uuid.New()
	} else {
		config.ID = c.ID
	}

	configStr := string(strconf)
	config.NoticeType = req.NoticeType
	config.Remark = req.Remark
	config.Status = req.Status
	config.Config = &configStr

	data, err := dal.SaveNotificationServicesConfig(&config)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (*NotificationServicesConfig) GetNotificationServicesConfig(noticeType string) (*model.NotificationServicesConfig, error) {
	c, err := dal.GetNotificationServicesConfigByType(noticeType)
	return c, err
}

// SendTestEmail sends a test email based on the provided request
// @params req model.SendTestEmailReq
func (*NotificationServicesConfig) SendTestEmail(req *model.SendTestEmailReq) error {
	if !utils.ValidateEmail(req.Email) {
		return errcode.New(200014)
	}
	c, err := dal.GetNotificationServicesConfigByType(model.NoticeType_Email)
	if err != nil {
		// Return an error if notification configuration is not found
		return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"notice_type": err.Error(),
		})
	}

	// Parse email configuration from the retrieved data
	var emailConf model.EmailConfig
	err = json.Unmarshal([]byte(*c.Config), &emailConf)
	if err != nil {
		// Return an error if parsing the email configuration fails
		return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Create a new email message
	m := gomail.NewMessage()
	// Set the sender's email address
	m.SetHeader("From", emailConf.FromEmail)
	// Set the recipient's email address (can be multiple recipients)
	m.SetHeader("To", req.Email)
	// Set the email subject
	m.SetHeader("Subject", "Iot Platform - Verification Code Notification")
	// Set the email body (can be plain text or HTML)
	m.SetBody("text/html", req.Body)

	// Set up SMTP server credentials (using Gmail as an example)
	d := gomail.NewDialer(emailConf.Host, emailConf.Port, emailConf.FromEmail, emailConf.FromPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		// Return an error if sending the email fails
		return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return nil // Return nil if no error occurs
}

// Send email message
func sendEmailMessage(message string, subject string, tenantId string, to ...string) (err error) {
	c, err := dal.GetNotificationServicesConfigByType(model.NoticeType_Email)
	if err != nil {
		return err
	}
	var emailConf model.EmailConfig
	err = json.Unmarshal([]byte(*c.Config), &emailConf)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(emailConf.Host, emailConf.Port, emailConf.FromEmail, emailConf.FromPassword)

	// if emailConf.SSL != nil {
	// 	if *emailConf.SSL {
	// 		d.TLSConfig = &tls.Config{
	// 			MinVersion:         tls.VersionTLS12,
	// 			MaxVersion:         tls.VersionTLS13,
	// 			InsecureSkipVerify: false,
	// 			CipherSuites: []uint16{
	// 				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	// 				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// 				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	// 				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	// 			},
	// 		}
	// 	}
	// }

	m := gomail.NewMessage()
	m.SetHeader("From", emailConf.FromEmail)
	m.SetHeader("To", to...)
	m.SetBody("text/html", message)
	m.SetHeader("Subject", subject)

	if err := d.DialAndSend(m); err != nil {
		logrus.Error(err)
		result := "FAILURE"
		remark := err.Error()
		GroupApp.NotificationHisory.SaveNotificationHistory(&model.NotificationHistory{
			ID:               uuid.New(),
			SendTime:         time.Now().UTC(),
			SendContent:      &message,
			SendTarget:       to[0],
			SendResult:       &result,
			NotificationType: model.NoticeType_Email,
			TenantID:         tenantId,
			Remark:           &remark,
		})
		return err
	} else {
		result := "SUCCESS"
		GroupApp.NotificationHisory.SaveNotificationHistory(&model.NotificationHistory{
			ID:               uuid.New(),
			SendTime:         time.Now().UTC(),
			SendContent:      &message,
			SendTarget:       to[0],
			SendResult:       &result,
			NotificationType: model.NoticeType_Email,
			TenantID:         tenantId,
			Remark:           nil,
		})
	}
	return nil
}

// Send notification
func (*NotificationServicesConfig) ExecuteNotification(notificationGroupId, title, content string) {

	notificationGroup, err := dal.GetNotificationGroupById(notificationGroupId)
	if err != nil {
		return
	}

	if notificationGroup.Status != "OPEN" {
		return
	}

	switch notificationGroup.NotificationType {
	case model.NoticeType_Member:
		// TODO: SEND TO MEMBER

	case model.NoticeType_Email:
		nConfig := make(map[string]string)
		err := json.Unmarshal([]byte(*notificationGroup.NotificationConfig), &nConfig)
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Debug("Notification configuration:", nConfig)
		emailList := strings.Split(nConfig["EMAIL"], ",")
		for _, ev := range emailList {
			logrus.Debug("Sending email address:", ev)
			sendEmailMessage(title, content, notificationGroup.TenantID, ev)
		}
	case model.NoticeType_Webhook:
		type WebhookConfig struct {
			PayloadURL string `json:"payload_url"`
			Secret     string `json:"secret"`
		}
		nConfig := make(map[string]WebhookConfig)
		err := json.Unmarshal([]byte(*notificationGroup.NotificationConfig), &nConfig)
		if err != nil {
			logrus.Error(err)
		}
		info := make(map[string]string)
		info["alert_title"] = title
		info["alert_details"] = content
		infoByte, _ := json.Marshal(info)
		err = http_client.SendSignedRequest(nConfig["webhook"].PayloadURL, string(infoByte), nConfig["webhook"].Secret)
		if err != nil {
			logrus.Error(err)
		}
	default:

		return
	}
}
