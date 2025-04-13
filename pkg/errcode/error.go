// pkg/errcode/error.go
package errcode

import "fmt"


type Error struct {
	Code         int                    `json:"code"`
	Data         interface{}            `json:"data,omitempty"`
	Variables    map[string]interface{} `json:"-"`                 
	Args         []interface{}          `json:"-"`                 
	CustomMsg    string                 `json:"message,omitempty"` 
	UseCustomMsg bool                   `json:"-"`                 
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error Code: %d", e.Code)
}


func New(code int) *Error {
	return &Error{
		Code: code,
	}
}


func NewWithMessage(code int, message string) *Error {
	return &Error{
		Code:         code,
		CustomMsg:    message,
		UseCustomMsg: true,
	}
}


func WithData(code int, data interface{}) *Error {
	return &Error{
		Code: code,
		Data: data,
	}
}


func Newf(code int, args ...interface{}) *Error {
	return &Error{
		Code: code,
		Args: args,
	}
}


func WithVars(code int, vars map[string]interface{}) *Error {
	return &Error{
		Code:      code,
		Variables: vars,
	}
}
