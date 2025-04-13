package middleware

import (
	"net/http"
	"strings"

	service "github.com/HustIoTPlatform/backend/internal/service"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Use Casbin; if the resource is in the table, validation is required; if not in the table, no validation is performed
// RBAC: User-Role-Function-Resource-Action

func CasbinRBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		// URLs that require validation: *api*
		if strings.Contains(c.Request.URL.Path, "api") {
			url := strings.TrimLeft(c.Request.URL.Path, "/")
			// Determine if the interface requires validation
			isVerify := service.GroupApp.Casbin.GetUrl(url)
			if isVerify {
				userClaims := c.MustGet("claims").(*utils.UserClaims)
				isSuccess := service.GroupApp.Casbin.Verify(userClaims.ID, url)
				if !isSuccess {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized access"})
					c.Abort()
					return
				}
			}
		}
	}
}
