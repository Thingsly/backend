package model

type SaveNotificationServicesConfigReq struct {
	EMailConfig *EmailConfig `json:"email_config" validate:"omitempty"`                         // Email configuration
	SMEConfig   *SMEConfig   `json:"sme_config" validate:"omitempty"`                           // SMS configuration
	NoticeType  string       `json:"notice_type" form:"notice_type" validate:"required,max=36"` // Notification type: EMAIL / SME
	Status      string       `json:"status" form:"status" validate:"required,max=36"`           // Status: ON/OFF
	Remark      *string      `json:"remark" form:"remark" validate:"omitempty,max=36"`          // Remarks
}

const NoticeType_Email = "EMAIL"
const NoticeType_SME = "SME"
const NoticeType_Member = "MEMBER"
const NoticeType_Voice = "VOICE"
const NoticeType_Webhook = "WEBHOOK"

type EmailConfig struct {
	// Email        string `json:"email" validate:"required"`
	Host         string `json:"host" form:"host" validate:"required"`
	Port         int    `json:"port" form:"port" validate:"required"`
	FromPassword string `json:"from_password" form:"from_password" validate:"required"`
	FromEmail    string `json:"from_email" form:"from_email" validate:"required"`
	SSL          *bool  `json:"ssl" form:"ssl" validate:"omitempty"`
}

type SMEConfig struct {
}

type SendTestEmailReq struct {
	Email string `json:"email" validate:"required"`
	Body  string `json:"body" form:"body" validate:"required"`
}
