// internal/middleware/response/response.go
package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Thingsly/backend/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Handler
type Handler struct {
	ErrManager *errcode.ErrorManager
}

// NewHandler
func NewHandler(configPath string, strConfigPath string) (*Handler, error) {
	errManager := errcode.NewErrorManager(configPath, strConfigPath)
	if err := errManager.LoadMessages(); err != nil {
		return nil, err
	}
	return &Handler{ErrManager: errManager}, nil
}

// Middleware
func (h *Handler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Handle panic
		defer func() {
			if err := recover(); err != nil {
				// For panic, create a system error
				sysErr := errcode.NewWithMessage(errcode.CodeSystemError, fmt.Sprint(err))
				h.handleError(c, sysErr)
				c.Abort()
			}
		}()

		c.Next()

		// 2. If the response has already been written, return
		if c.Writer.Written() {
			return
		}

		// 3. Handle error response
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			switch e := err.(type) {
			case *errcode.Error:
				h.handleError(c, e)
			default:
				// For other errors, create a system error
				sysErr := errcode.NewWithMessage(errcode.CodeSystemError, err.Error())
				h.handleError(c, sysErr)
			}
			return
		}

		// 4. Handle success response
		if data, exists := c.Get("data"); exists {
			h.responseSuccess(c, data)
		}
	}
}

// responseSuccess
func (h *Handler) responseSuccess(c *gin.Context, data interface{}) {
	lang := c.GetHeader("Accept-Language")
	c.JSON(http.StatusOK, &Response{
		Code:    errcode.CodeSuccess,
		Message: h.ErrManager.GetMessage(errcode.CodeSuccess, lang),
		Data:    data,
	})
}

// pkg/response/response.go
func (h *Handler) handleError(c *gin.Context, err *errcode.Error) {
	var msg string
	if err.UseCustomMsg {
		msg = err.CustomMsg
	} else {
		lang := c.GetHeader("Accept-Language")
		msg = h.ErrManager.GetMessage(err.Code, lang)

		if err.Args != nil {
			msg = fmt.Sprintf(msg, err.Args...)
		}

		if err.Variables != nil {
			for k, v := range err.Variables {
				msg = strings.ReplaceAll(msg, "${"+k+"}", fmt.Sprint(v))
			}
		}
	}

	resp := &Response{
		Code:    err.Code,
		Message: msg,
	}
	if err.Data != nil {
		resp.Data = err.Data
	}

	c.JSON(http.StatusOK, resp)
}
