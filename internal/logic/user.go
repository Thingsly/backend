package logic

import (
	"context"

	"github.com/Thingsly/backend/internal/query"
	"github.com/Thingsly/backend/pkg/constant"
)

func UserIsEncrypt(ctx context.Context) bool {
	var (
		sysFunction = query.SysFunction
	)

	info, err := sysFunction.WithContext(ctx).Where(sysFunction.Name.Eq("frontend_res")).First()
	if err != nil {
		return false
	}
	if info.EnableFlag == constant.DisableFlag {
		return false
	}
	return true
}

func UserIsShare(ctx context.Context) bool {
	var (
		sysFunction = query.SysFunction
	)

	info, err := sysFunction.WithContext(ctx).Where(sysFunction.Name.Eq("shared_account")).First()
	if err != nil {
		return false
	}
	if info.EnableFlag == constant.DisableFlag {
		return false
	}
	return true
}
