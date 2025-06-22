package global

import (
	"github.com/Thingsly/backend/internal/middleware/response"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	VERSION         = "0.0.7"
	VERSION_NUMBER  = 7
	SYSTEM_VERSION  = "v1.1.7"
	DB              *gorm.DB
	REDIS           *redis.Client
	STATUS_REDIS    *redis.Client
	CasbinEnforcer  *casbin.Enforcer
	OtaAddress      string
	TLSSEManager    *SSEManager
	ResponseHandler *response.Handler
)

type EventData struct {
	Name    string
	Message string
}

var EventChan chan EventData
