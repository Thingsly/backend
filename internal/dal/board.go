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
		board = query.Board
		err   error
	)

	// Get demo tenant's board configuration
	demoBoard, err := board.WithContext(ctx).Where(query.Board.TenantID.Eq("d616bcbb"), query.Board.HomeFlag.Eq("Y")).First()
	if err != nil {
		logrus.Error(ctx, "[BoardQuery]GetDemoBoard failed:", err)
		return err
	}

	// Create the default homepage dashboard with demo tenant's configuration
	err = board.WithContext(ctx).Create(&model.Board{
		ID:          uuid.New(),
		Name:        demoBoard.Name,
		Config:      demoBoard.Config,
		TenantID:    tenantid,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		HomeFlag:    "Y",
		Description: demoBoard.Description,
		Remark:      demoBoard.Remark,
		MenuFlag:    demoBoard.MenuFlag,
	})
	if err != nil {
		logrus.Error(ctx, "[BoardQuery]CreateDefaultBoard failed:", err)
	}
	return err
}
