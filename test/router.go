package test

import (
	"github.com/Thingsly/backend/internal/middleware"
	"github.com/Thingsly/backend/internal/middleware/response"
	"github.com/gin-gonic/gin"
)

// SetUpRouter configures and returns a new Gin router with necessary middleware
func SetUpRouter() *gin.Engine {
	router := gin.Default()

	// Add response handler middleware
	responseHandler, _ := response.NewHandler("", "")
	router.Use(responseHandler.Middleware())

	// Add JWT auth middleware
	router.Use(middleware.JWTAuth())

	// Add Casbin RBAC middleware
	router.Use(middleware.CasbinRBAC())

	return router
}
