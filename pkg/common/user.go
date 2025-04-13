package common

import (
	constant "github.com/HustIoTPlatform/backend/pkg/constant"
)

func CheckUserIsAdmin(authority string) bool {
	return authority == constant.SYS_ADMIN
}
