package model

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// Page request parameters
type PageReq struct {
	Page     int `json:"page" form:"page" validate:"required,gte=1"`                    // Page number
	PageSize int `json:"page_size" form:"page_size" validate:"required,gte=1,lte=1000"` // Items per page
}

// PutMessage
// @DESCRIPTION: Common parameters for sending a message
type PutMessage struct {
	DeviceID string `json:"device_id" form:"device_id" validate:"required,max=36"`
	Value    string `json:"value" form:"value" validate:"required,max=9999"`
}

type PutMessageForCommand struct {
	DeviceID string  `json:"device_id" form:"device_id" validate:"required,max=36"`
	Value    *string `json:"value" form:"value" validate:"omitempty,max=9999"`
	Identify string  `json:"identify" form:"identify" validate:"required,max=255"`
}

type ParamID struct {
	ID string `query:"id" form:"id" json:"id" validate:"required"`
}

const OPEN = "OPEN"
const CLOSE = "CLOSE"

// For parameters from the front-end that cannot have a fixed structure, such as: products.AdditionalInfo
// Use *json.RawMessage to receive them and convert them into a string that can be stored in the database
// Also, remove extra spaces in the JSON string
func JsonRawMessage2Str(in *json.RawMessage) (str string, err error) {
	var data map[string]interface{}
	err = json.Unmarshal([]byte(*in), &data)
	if err != nil {
		return str, err
	}
	compactJson, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err)
		return str, err
	}
	str = string(compactJson)
	return
}
