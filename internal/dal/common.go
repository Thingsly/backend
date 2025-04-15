package dal

import (
	"github.com/Thingsly/backend/internal/query"
	"github.com/Thingsly/backend/pkg/global"
)

func StartTransaction() (*query.QueryTx, error) {
	tx := query.Use(global.DB).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func Rollback(tx *query.QueryTx) error {
	if err := tx.Rollback(); err != nil {
		return err
	}
	return nil
}

func Commit(tx *query.QueryTx) error {
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
