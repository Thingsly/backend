package sseapi

import (
	"time"

	"github.com/Thingsly/backend/pkg/errcode"
	"github.com/Thingsly/backend/pkg/global"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SSEApi struct{}

// api/v1/events

func (*SSEApi) HandleSystemEvents(c *gin.Context) {
	userClaims, ok := c.MustGet("claims").(*utils.UserClaims)
	if !ok {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"error": "UserClaims not found",
		}))
		return
	}

	logrus.WithFields(logrus.Fields{
		"tenantID":  userClaims.TenantID,
		"userEmail": userClaims.Email,
	}).Info("User connected to SSE")

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	clientID := global.TPSSEManager.AddClient(userClaims.TenantID, userClaims.ID, c.Writer)
	defer global.TPSSEManager.RemoveClient(userClaims.TenantID, clientID)

	// Send initial success message
	c.SSEvent("message", "Connected to system events")
	c.Writer.Flush()

	// Create a timer for sending heartbeat
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	// Create a channel to check if the client is still connected
	done := make(chan bool)
	go func() {
		<-c.Request.Context().Done()
		done <- true
	}()

	for {
		select {
		case <-heartbeatTicker.C:
			// Send heartbeat message
			c.SSEvent("heartbeat", time.Now().Unix())
			c.Writer.Flush()
		case <-done:
			logrus.WithFields(logrus.Fields{
				"tenantID":  userClaims.TenantID,
				"userEmail": userClaims.Email,
			}).Info("User disconnected from SSE")
			return
		}
	}
}
