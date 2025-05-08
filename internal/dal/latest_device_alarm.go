package dal

import (
	"context"

	"github.com/Thingsly/backend/internal/query"
)

type LatestDeviceAlarmQuery struct{}

func (q *LatestDeviceAlarmQuery) CountDevicesByTenantAndStatus(ctx context.Context, tenantID string) (int64, error) {
	lda := query.LatestDeviceAlarm

	count, err := lda.WithContext(ctx).
		Where(lda.TenantID.Eq(tenantID)).
		Where(lda.AlarmStatus.Neq("N")).
		Distinct(lda.DeviceID).
		Count()

	return count, err
}