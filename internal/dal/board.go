package dal

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"

	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

func CreateBoard(boards *model.Board) error {
	if _, err := query.Board.Where(query.Board.HomeFlag.Eq("Y"), query.Board.TenantID.Eq(boards.TenantID)).First(); err == nil {
		return fmt.Errorf("dashboard already exists on the homepage")
	}
	return query.Board.Create(boards)
}

func UpdateBoard(boards *model.Board) error {
	p := query.Board
	r, err := query.Board.Where(p.ID.Eq(boards.ID)).Updates(boards)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if r.RowsAffected == 0 {
		return fmt.Errorf("no data updated")
	}
	return err
}

func DeleteBoard(id string) error {
	r, err := query.Board.Where(query.Board.ID.Eq(id)).Delete()
	// The interface with an incorrect ID also returns success.
	if r.RowsAffected == 0 {
		return nil
	}
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func GetBoardListByPage(boards *model.GetBoardListByPageReq, tenantId string) (int64, interface{}, error) {
	q := query.Board
	var count int64
	queryBuilder := q.WithContext(context.Background())
	queryBuilder = queryBuilder.Where(q.TenantID.Eq(tenantId))

	if boards.Name != nil && *boards.Name != "" {
		queryBuilder = queryBuilder.Where(q.Name.Like(fmt.Sprintf("%%%s%%", *boards.Name)))
	}

	if boards.HomeFlag != nil && *boards.HomeFlag != "" {
		queryBuilder = queryBuilder.Where(q.HomeFlag.Eq(*boards.HomeFlag))
	}
	count, err := queryBuilder.Count()
	if err != nil {
		logrus.Error(err)
		return count, nil, err
	}
	if boards.Page != 0 && boards.PageSize != 0 {
		queryBuilder = queryBuilder.Limit(boards.PageSize)
		queryBuilder = queryBuilder.Offset((boards.Page - 1) * boards.PageSize)
	}
	queryBuilder = queryBuilder.Order(q.CreatedAt.Desc())
	boardsList, err := queryBuilder.Select(q.ID, q.Name, q.HomeFlag, q.MenuFlag, q.UpdatedAt, q.CreatedAt, q.Description, q.Remark, q.TenantID).Find()
	if err != nil {
		logrus.Error(err)
		return count, boardsList, err
	}

	return count, boardsList, err
}

func GetBoard(id string, tenantId string) (interface{}, error) {
	p := query.Board
	board, err := query.Board.Where(p.ID.Eq(id)).Where(p.TenantID.Eq(tenantId)).Select().First()
	if err != nil {
		logrus.Error(err)
	}
	return board, err
}

func GetBoardListByTenantId(tenantid string) (int64, interface{}, error) {
	q := query.Board
	var count int64
	queryBuilder := q.WithContext(context.Background())
	boardsList, err := queryBuilder.Where(q.TenantID.Eq(tenantid), q.HomeFlag.Eq("Y")).Select().First()
	if err != nil {
		// If there is no homepage dashboard, return empty.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return count, nil, nil
		}
		logrus.Error(err)
		return count, boardsList, err
	}
	count, err = queryBuilder.Count()
	return count, boardsList, err
}

type BoardQuery struct{}

func (BoardQuery) Create(ctx context.Context, info *model.Board) (*model.Board, error) {
	var (
		board = query.Board
		err   error
	)
	if err = board.WithContext(ctx).Create(info); err != nil {
		logrus.Error(ctx, "[BoardQuery]First failed:", err)
	}
	return info, err
}

func (BoardQuery) First(ctx context.Context, option ...gen.Condition) (info *model.Board, err error) {
	board := query.Board
	info, err = board.WithContext(ctx).Where(option...).First()
	if err != nil {
		logrus.Error(ctx, "[BoardQuery]First failed:", err)
	}
	return info, err
}

// Set other tenant homepage dashboards as non-homepage.
func (BoardQuery) UpdateHomeFlagN(ctx context.Context, tenantid string) error {
	var (
		board = query.Board
		err   error
	)
	if _, err := board.WithContext(ctx).Where(query.Board.TenantID.Eq(tenantid), query.Board.HomeFlag.Eq("Y")).Updates(map[string]interface{}{"home_flag": "N"}); err != nil {
		logrus.Error(ctx, "update failed:", err)
	}
	return err
}

// Add a default homepage dashboard for the newly added tenant.
func (BoardQuery) CreateDefaultBoard(ctx context.Context, tenantid string) error {
	var (
		board  = query.Board
		config = `[{"x":3,"y":0,"w":3,"h":2,"minW":2,"minH":2,"i":1745505262993650,"data":{"cardId":"on-num","type":"builtin","title":"Online Device Count","config":{},"layout":{"w":3,"h":2,"minH":2,"minW":2},"basicSettings":{},"dataSource":{"origin":"system","systemSource":[{}],"deviceSource":[{}]}},"moved":false},{"x":6,"y":0,"w":3,"h":2,"minW":2,"minH":2,"i":1745505261029718,"data":{"cardId":"off-num","type":"builtin","title":"Offline Device Count","config":{},"layout":{"w":3,"h":2,"minH":2,"minW":2},"basicSettings":{},"dataSource":{"origin":"system","systemSource":[{}],"deviceSource":[{}]}},"moved":false},{"x":0,"y":0,"w":3,"h":2,"minW":2,"minH":2,"i":1744117893640210,"data":{"cardId":"access-num","type":"builtin","title":"Visit Count","config":{},"layout":{"w":3,"h":2,"minH":2,"minW":2},"basicSettings":{},"dataSource":{"origin":"system","systemSource":[{}],"deviceSource":[{}]}},"moved":false},{"x":6,"y":2,"w":3,"h":2,"minW":2,"minH":2,"i":1745417880007150,"data":{"cardId":"memory-usage","type":"builtin","title":"card.memoryUsage","config":{},"layout":{"w":3,"h":2,"minH":2,"minW":2},"basicSettings":{},"dataSource":{"origin":"system","systemSource":[{}],"deviceSource":[{}]}},"moved":false},{"x":0,"y":2,"w":3,"h":2,"minW":2,"minH":2,"i":1745418459170384,"data":{"cardId":"cpu-usage","type":"builtin","title":"card.cpuUsage","config":{},"layout":{"w":3,"h":2,"minH":2,"minW":2},"basicSettings":{},"dataSource":{"origin":"system","systemSource":[{}],"deviceSource":[{}]}},"moved":false},{"x":3,"y":2,"w":3,"h":2,"minW":2,"minH":2,"i":1745419824863842,"data":{"cardId":"disk-usage","type":"builtin","title":"card.diskUsage","config":{},"layout":{"w":3,"h":2,"minH":2,"minW":2},"basicSettings":{},"dataSource":{"origin":"system","systemSource":[{}],"deviceSource":[{}]}},"moved":false},{"x":3,"y":4,"w":2,"h":6,"minW":2,"minH":2,"i":1745505353977076,"data":{"cardId":"recently-visited","type":"builtin","title":"Recent Access","config":{},"layout":{"w":3,"h":2,"minH":2,"minW":2},"basicSettings":{},"dataSource":{"origin":"system","systemSource":[{}],"deviceSource":[{}]}},"moved":false},{"x":9,"y":0,"w":3,"h":4,"minW":2,"minH":2,"i":1745504866294083,"data":{"cardId":"metrics-history","type":"builtin","title":"System Metrics History","config":{},"layout":{"w":2,"h":2,"minW":2,"minH":2},"basicSettings":{},"dataSource":{"origin":"system","systemSource":[{}],"deviceSource":[{}]}},"moved":false},{"x":0,"y":4,"w":3,"h":6,"minW":2,"minH":2,"i":1745506739499397,"data":{"cardId":"operation-guide","type":"builtin","title":"Operation Guide","config":{"guideList":[{"titleKey":"card.operationGuideCard.guideItems.addDevice.title","descriptionKey":"card.operationGuideCard.guideItems.addDevice.description","link":"/device/manage"},{"titleKey":"card.operationGuideCard.guideItems.configureDevice.title","descriptionKey":"card.operationGuideCard.guideItems.configureDevice.description"},{"titleKey":"card.operationGuideCard.guideItems.createDashboard.title","descriptionKey":"card.operationGuideCard.guideItems.createDashboard.description"}],"guideListAdmin":[{"titleKey":"card.operationGuideAdmin.guideItems.createTenant.title","descriptionKey":"card.operationGuideAdmin.guideItems.createTenant.description","link":"/management/user"},{"titleKey":"card.operationGuideAdmin.guideItems.configureNotification.title","descriptionKey":"card.operationGuideAdmin.guideItems.configureNotification.description"},{"titleKey":"card.operationGuideAdmin.guideItems.configurePlugin.title","descriptionKey":"card.operationGuideAdmin.guideItems.configurePlugin.description"}]},"layout":{"w":3,"h":5,"minW":2,"minH":2},"basicSettings":{},"dataSource":{"origin":"system","isSupportTimeRange":false,"dataTimeRange":"","isSupportAggregate":false,"dataAggregateRange":"","systemSource":[],"deviceSource":[]}},"moved":false},{"x":5,"y":4,"w":5,"h":6,"minW":2,"minH":2,"i":1745508228132635,"data":{"cardId":"tenant-chart","type":"builtin","title":"cards.alarmInfo5","config":{},"layout":{"w":2,"h":2,"minW":2,"minH":2},"basicSettings":{},"dataSource":{"origin":"device","isSupportTimeRange":true,"dataTimeRange":"1h","isSupportAggregate":true,"dataAggregateRange":"1m","systemSource":[],"deviceSource":[]}},"moved":false},{"x":10,"y":7,"w":2,"h":3,"minW":2,"minH":1,"i":1745510975816929,"data":{"cardId":"version-info","type":"builtin","title":"Version Info","config":{},"layout":{"w":3,"h":1,"minW":2,"minH":1},"basicSettings":{},"dataSource":{"origin":"system","systemSource":[{}],"deviceSource":[{}]}},"moved":false}]`
	)
	// Create the default homepage dashboard according to the above SQL statement
	err := board.WithContext(ctx).Create(&model.Board{
		ID:        uuid.New(),
		Name:      "Home",
		Config:    &config,
		TenantID:  tenantid,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		HomeFlag:  "Y",
		Remark:    nil,
	})
	if err != nil {
		logrus.Error(ctx, "[BoardQuery]CreateDefaultBoard failed:", err)
	}
	return err
}
