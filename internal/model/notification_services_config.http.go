package model

type SaveNotificationServicesConfigReq struct {
	EMailConfig *EmailConfig `json:"email_config" validate:"omitempty"`                                              
	SMEConfig   *SMEConfig   `json:"sme_config" validate:"omitempty"`                                                
	NoticeType  string       `json:"notice_type" form:"notice_type" validate:"required,max=36,oneof=EMAIL SME_CODE"` 
	Status      string       `json:"status" form:"status" validate:"required,max=36,oneof=OPEN CLOSE"`               
	Remark      *string      `json:"remark" form:"remark" validate:"omitempty,max=36"`                               
}

const (
	NoticeType_Email    = "EMAIL"
	NoticeType_SME_CODE = "SME_CODE"
	NoticeType_Member   = "MEMBER"
	NoticeType_Voice    = "VOICE"
	NoticeType_Webhook  = "WEBHOOK"
)

type EmailConfig struct {
	// Email        string `json:"email" validate:"required"`
	Host         string `json:"host" form:"host" validate:"required"`
	Port         int    `json:"port" form:"port" validate:"required"`
	FromPassword string `json:"from_password" form:"from_password" validate:"required"`
	FromEmail    string `json:"from_email" form:"from_email" validate:"required"`
	SSL          *bool  `json:"ssl" form:"ssl" validate:"omitempty"`
}

type SMEConfig struct {
	Provider string `json:"provider" form:"provider" validate:"required,max=36,oneof=TWILIO"`
	TwilioSMSConfig *TwilioSMSConfig `json:"twilio_sms_config" form:"twilio_sms_config" validate:"omitempty"`
}

type TwilioSMSConfig struct {
	AccessKeyID     string `json:"access_key_id" form:"access_key_id" validate:"required,max=100"`         
	AccessKeySecret string `json:"access_key_secret" form:"access_key_secret" validate:"required,max=100"` 
	Endpoint        string `json:"endpoint" form:"endpoint" validate:"required,max=100"`                   
	SignName        string `json:"sign_name" form:"sign_name" validate:"required,max=100"`                 
	TemplateCode    string `json:"template_code" form:"template_code" validate:"required,max=36"`         
}

type SendTestEmailReq struct {
	Email string `json:"email" validate:"required"`
	Body  string `json:"body" form:"body" validate:"required"`
}
