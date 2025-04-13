// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/HustIoTPlatform/backend/internal/model"
)

func newSceneAutomationLog(db *gorm.DB, opts ...gen.DOOption) sceneAutomationLog {
	_sceneAutomationLog := sceneAutomationLog{}

	_sceneAutomationLog.sceneAutomationLogDo.UseDB(db, opts...)
	_sceneAutomationLog.sceneAutomationLogDo.UseModel(&model.SceneAutomationLog{})

	tableName := _sceneAutomationLog.sceneAutomationLogDo.TableName()
	_sceneAutomationLog.ALL = field.NewAsterisk(tableName)
	_sceneAutomationLog.SceneAutomationID = field.NewString(tableName, "scene_automation_id")
	_sceneAutomationLog.ExecutedAt = field.NewTime(tableName, "executed_at")
	_sceneAutomationLog.Detail = field.NewString(tableName, "detail")
	_sceneAutomationLog.ExecutionResult = field.NewString(tableName, "execution_result")
	_sceneAutomationLog.TenantID = field.NewString(tableName, "tenant_id")
	_sceneAutomationLog.Remark = field.NewString(tableName, "remark")

	_sceneAutomationLog.fillFieldMap()

	return _sceneAutomationLog
}

type sceneAutomationLog struct {
	sceneAutomationLogDo

	ALL               field.Asterisk
	SceneAutomationID field.String 
	ExecutedAt        field.Time   
	Detail            field.String 
	ExecutionResult   field.String 
	TenantID          field.String
	Remark            field.String

	fieldMap map[string]field.Expr
}

func (s sceneAutomationLog) Table(newTableName string) *sceneAutomationLog {
	s.sceneAutomationLogDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s sceneAutomationLog) As(alias string) *sceneAutomationLog {
	s.sceneAutomationLogDo.DO = *(s.sceneAutomationLogDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *sceneAutomationLog) updateTableName(table string) *sceneAutomationLog {
	s.ALL = field.NewAsterisk(table)
	s.SceneAutomationID = field.NewString(table, "scene_automation_id")
	s.ExecutedAt = field.NewTime(table, "executed_at")
	s.Detail = field.NewString(table, "detail")
	s.ExecutionResult = field.NewString(table, "execution_result")
	s.TenantID = field.NewString(table, "tenant_id")
	s.Remark = field.NewString(table, "remark")

	s.fillFieldMap()

	return s
}

func (s *sceneAutomationLog) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *sceneAutomationLog) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 6)
	s.fieldMap["scene_automation_id"] = s.SceneAutomationID
	s.fieldMap["executed_at"] = s.ExecutedAt
	s.fieldMap["detail"] = s.Detail
	s.fieldMap["execution_result"] = s.ExecutionResult
	s.fieldMap["tenant_id"] = s.TenantID
	s.fieldMap["remark"] = s.Remark
}

func (s sceneAutomationLog) clone(db *gorm.DB) sceneAutomationLog {
	s.sceneAutomationLogDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s sceneAutomationLog) replaceDB(db *gorm.DB) sceneAutomationLog {
	s.sceneAutomationLogDo.ReplaceDB(db)
	return s
}

type sceneAutomationLogDo struct{ gen.DO }

type ISceneAutomationLogDo interface {
	gen.SubQuery
	Debug() ISceneAutomationLogDo
	WithContext(ctx context.Context) ISceneAutomationLogDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ISceneAutomationLogDo
	WriteDB() ISceneAutomationLogDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ISceneAutomationLogDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ISceneAutomationLogDo
	Not(conds ...gen.Condition) ISceneAutomationLogDo
	Or(conds ...gen.Condition) ISceneAutomationLogDo
	Select(conds ...field.Expr) ISceneAutomationLogDo
	Where(conds ...gen.Condition) ISceneAutomationLogDo
	Order(conds ...field.Expr) ISceneAutomationLogDo
	Distinct(cols ...field.Expr) ISceneAutomationLogDo
	Omit(cols ...field.Expr) ISceneAutomationLogDo
	Join(table schema.Tabler, on ...field.Expr) ISceneAutomationLogDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ISceneAutomationLogDo
	RightJoin(table schema.Tabler, on ...field.Expr) ISceneAutomationLogDo
	Group(cols ...field.Expr) ISceneAutomationLogDo
	Having(conds ...gen.Condition) ISceneAutomationLogDo
	Limit(limit int) ISceneAutomationLogDo
	Offset(offset int) ISceneAutomationLogDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ISceneAutomationLogDo
	Unscoped() ISceneAutomationLogDo
	Create(values ...*model.SceneAutomationLog) error
	CreateInBatches(values []*model.SceneAutomationLog, batchSize int) error
	Save(values ...*model.SceneAutomationLog) error
	First() (*model.SceneAutomationLog, error)
	Take() (*model.SceneAutomationLog, error)
	Last() (*model.SceneAutomationLog, error)
	Find() ([]*model.SceneAutomationLog, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.SceneAutomationLog, err error)
	FindInBatches(result *[]*model.SceneAutomationLog, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.SceneAutomationLog) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ISceneAutomationLogDo
	Assign(attrs ...field.AssignExpr) ISceneAutomationLogDo
	Joins(fields ...field.RelationField) ISceneAutomationLogDo
	Preload(fields ...field.RelationField) ISceneAutomationLogDo
	FirstOrInit() (*model.SceneAutomationLog, error)
	FirstOrCreate() (*model.SceneAutomationLog, error)
	FindByPage(offset int, limit int) (result []*model.SceneAutomationLog, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ISceneAutomationLogDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (s sceneAutomationLogDo) Debug() ISceneAutomationLogDo {
	return s.withDO(s.DO.Debug())
}

func (s sceneAutomationLogDo) WithContext(ctx context.Context) ISceneAutomationLogDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s sceneAutomationLogDo) ReadDB() ISceneAutomationLogDo {
	return s.Clauses(dbresolver.Read)
}

func (s sceneAutomationLogDo) WriteDB() ISceneAutomationLogDo {
	return s.Clauses(dbresolver.Write)
}

func (s sceneAutomationLogDo) Session(config *gorm.Session) ISceneAutomationLogDo {
	return s.withDO(s.DO.Session(config))
}

func (s sceneAutomationLogDo) Clauses(conds ...clause.Expression) ISceneAutomationLogDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s sceneAutomationLogDo) Returning(value interface{}, columns ...string) ISceneAutomationLogDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s sceneAutomationLogDo) Not(conds ...gen.Condition) ISceneAutomationLogDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s sceneAutomationLogDo) Or(conds ...gen.Condition) ISceneAutomationLogDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s sceneAutomationLogDo) Select(conds ...field.Expr) ISceneAutomationLogDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s sceneAutomationLogDo) Where(conds ...gen.Condition) ISceneAutomationLogDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s sceneAutomationLogDo) Order(conds ...field.Expr) ISceneAutomationLogDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s sceneAutomationLogDo) Distinct(cols ...field.Expr) ISceneAutomationLogDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s sceneAutomationLogDo) Omit(cols ...field.Expr) ISceneAutomationLogDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s sceneAutomationLogDo) Join(table schema.Tabler, on ...field.Expr) ISceneAutomationLogDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s sceneAutomationLogDo) LeftJoin(table schema.Tabler, on ...field.Expr) ISceneAutomationLogDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s sceneAutomationLogDo) RightJoin(table schema.Tabler, on ...field.Expr) ISceneAutomationLogDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s sceneAutomationLogDo) Group(cols ...field.Expr) ISceneAutomationLogDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s sceneAutomationLogDo) Having(conds ...gen.Condition) ISceneAutomationLogDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s sceneAutomationLogDo) Limit(limit int) ISceneAutomationLogDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s sceneAutomationLogDo) Offset(offset int) ISceneAutomationLogDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s sceneAutomationLogDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ISceneAutomationLogDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s sceneAutomationLogDo) Unscoped() ISceneAutomationLogDo {
	return s.withDO(s.DO.Unscoped())
}

func (s sceneAutomationLogDo) Create(values ...*model.SceneAutomationLog) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s sceneAutomationLogDo) CreateInBatches(values []*model.SceneAutomationLog, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s sceneAutomationLogDo) Save(values ...*model.SceneAutomationLog) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s sceneAutomationLogDo) First() (*model.SceneAutomationLog, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.SceneAutomationLog), nil
	}
}

func (s sceneAutomationLogDo) Take() (*model.SceneAutomationLog, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.SceneAutomationLog), nil
	}
}

func (s sceneAutomationLogDo) Last() (*model.SceneAutomationLog, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.SceneAutomationLog), nil
	}
}

func (s sceneAutomationLogDo) Find() ([]*model.SceneAutomationLog, error) {
	result, err := s.DO.Find()
	return result.([]*model.SceneAutomationLog), err
}

func (s sceneAutomationLogDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.SceneAutomationLog, err error) {
	buf := make([]*model.SceneAutomationLog, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s sceneAutomationLogDo) FindInBatches(result *[]*model.SceneAutomationLog, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s sceneAutomationLogDo) Attrs(attrs ...field.AssignExpr) ISceneAutomationLogDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s sceneAutomationLogDo) Assign(attrs ...field.AssignExpr) ISceneAutomationLogDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s sceneAutomationLogDo) Joins(fields ...field.RelationField) ISceneAutomationLogDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s sceneAutomationLogDo) Preload(fields ...field.RelationField) ISceneAutomationLogDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s sceneAutomationLogDo) FirstOrInit() (*model.SceneAutomationLog, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.SceneAutomationLog), nil
	}
}

func (s sceneAutomationLogDo) FirstOrCreate() (*model.SceneAutomationLog, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.SceneAutomationLog), nil
	}
}

func (s sceneAutomationLogDo) FindByPage(offset int, limit int) (result []*model.SceneAutomationLog, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s sceneAutomationLogDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s sceneAutomationLogDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s sceneAutomationLogDo) Delete(models ...*model.SceneAutomationLog) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *sceneAutomationLogDo) withDO(do gen.Dao) *sceneAutomationLogDo {
	s.DO = *do.(*gen.DO)
	return s
}
