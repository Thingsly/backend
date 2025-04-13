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

func newSysDictLanguage(db *gorm.DB, opts ...gen.DOOption) sysDictLanguage {
	_sysDictLanguage := sysDictLanguage{}

	_sysDictLanguage.sysDictLanguageDo.UseDB(db, opts...)
	_sysDictLanguage.sysDictLanguageDo.UseModel(&model.SysDictLanguage{})

	tableName := _sysDictLanguage.sysDictLanguageDo.TableName()
	_sysDictLanguage.ALL = field.NewAsterisk(tableName)
	_sysDictLanguage.ID = field.NewString(tableName, "id")
	_sysDictLanguage.DictID = field.NewString(tableName, "dict_id")
	_sysDictLanguage.LanguageCode = field.NewString(tableName, "language_code")
	_sysDictLanguage.Translation = field.NewString(tableName, "translation")

	_sysDictLanguage.fillFieldMap()

	return _sysDictLanguage
}

type sysDictLanguage struct {
	sysDictLanguageDo

	ALL          field.Asterisk
	ID           field.String 
	DictID       field.String 
	LanguageCode field.String 
	Translation  field.String 

	fieldMap map[string]field.Expr
}

func (s sysDictLanguage) Table(newTableName string) *sysDictLanguage {
	s.sysDictLanguageDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s sysDictLanguage) As(alias string) *sysDictLanguage {
	s.sysDictLanguageDo.DO = *(s.sysDictLanguageDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *sysDictLanguage) updateTableName(table string) *sysDictLanguage {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewString(table, "id")
	s.DictID = field.NewString(table, "dict_id")
	s.LanguageCode = field.NewString(table, "language_code")
	s.Translation = field.NewString(table, "translation")

	s.fillFieldMap()

	return s
}

func (s *sysDictLanguage) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *sysDictLanguage) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 4)
	s.fieldMap["id"] = s.ID
	s.fieldMap["dict_id"] = s.DictID
	s.fieldMap["language_code"] = s.LanguageCode
	s.fieldMap["translation"] = s.Translation
}

func (s sysDictLanguage) clone(db *gorm.DB) sysDictLanguage {
	s.sysDictLanguageDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s sysDictLanguage) replaceDB(db *gorm.DB) sysDictLanguage {
	s.sysDictLanguageDo.ReplaceDB(db)
	return s
}

type sysDictLanguageDo struct{ gen.DO }

type ISysDictLanguageDo interface {
	gen.SubQuery
	Debug() ISysDictLanguageDo
	WithContext(ctx context.Context) ISysDictLanguageDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ISysDictLanguageDo
	WriteDB() ISysDictLanguageDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ISysDictLanguageDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ISysDictLanguageDo
	Not(conds ...gen.Condition) ISysDictLanguageDo
	Or(conds ...gen.Condition) ISysDictLanguageDo
	Select(conds ...field.Expr) ISysDictLanguageDo
	Where(conds ...gen.Condition) ISysDictLanguageDo
	Order(conds ...field.Expr) ISysDictLanguageDo
	Distinct(cols ...field.Expr) ISysDictLanguageDo
	Omit(cols ...field.Expr) ISysDictLanguageDo
	Join(table schema.Tabler, on ...field.Expr) ISysDictLanguageDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ISysDictLanguageDo
	RightJoin(table schema.Tabler, on ...field.Expr) ISysDictLanguageDo
	Group(cols ...field.Expr) ISysDictLanguageDo
	Having(conds ...gen.Condition) ISysDictLanguageDo
	Limit(limit int) ISysDictLanguageDo
	Offset(offset int) ISysDictLanguageDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ISysDictLanguageDo
	Unscoped() ISysDictLanguageDo
	Create(values ...*model.SysDictLanguage) error
	CreateInBatches(values []*model.SysDictLanguage, batchSize int) error
	Save(values ...*model.SysDictLanguage) error
	First() (*model.SysDictLanguage, error)
	Take() (*model.SysDictLanguage, error)
	Last() (*model.SysDictLanguage, error)
	Find() ([]*model.SysDictLanguage, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.SysDictLanguage, err error)
	FindInBatches(result *[]*model.SysDictLanguage, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.SysDictLanguage) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ISysDictLanguageDo
	Assign(attrs ...field.AssignExpr) ISysDictLanguageDo
	Joins(fields ...field.RelationField) ISysDictLanguageDo
	Preload(fields ...field.RelationField) ISysDictLanguageDo
	FirstOrInit() (*model.SysDictLanguage, error)
	FirstOrCreate() (*model.SysDictLanguage, error)
	FindByPage(offset int, limit int) (result []*model.SysDictLanguage, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ISysDictLanguageDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (s sysDictLanguageDo) Debug() ISysDictLanguageDo {
	return s.withDO(s.DO.Debug())
}

func (s sysDictLanguageDo) WithContext(ctx context.Context) ISysDictLanguageDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s sysDictLanguageDo) ReadDB() ISysDictLanguageDo {
	return s.Clauses(dbresolver.Read)
}

func (s sysDictLanguageDo) WriteDB() ISysDictLanguageDo {
	return s.Clauses(dbresolver.Write)
}

func (s sysDictLanguageDo) Session(config *gorm.Session) ISysDictLanguageDo {
	return s.withDO(s.DO.Session(config))
}

func (s sysDictLanguageDo) Clauses(conds ...clause.Expression) ISysDictLanguageDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s sysDictLanguageDo) Returning(value interface{}, columns ...string) ISysDictLanguageDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s sysDictLanguageDo) Not(conds ...gen.Condition) ISysDictLanguageDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s sysDictLanguageDo) Or(conds ...gen.Condition) ISysDictLanguageDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s sysDictLanguageDo) Select(conds ...field.Expr) ISysDictLanguageDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s sysDictLanguageDo) Where(conds ...gen.Condition) ISysDictLanguageDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s sysDictLanguageDo) Order(conds ...field.Expr) ISysDictLanguageDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s sysDictLanguageDo) Distinct(cols ...field.Expr) ISysDictLanguageDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s sysDictLanguageDo) Omit(cols ...field.Expr) ISysDictLanguageDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s sysDictLanguageDo) Join(table schema.Tabler, on ...field.Expr) ISysDictLanguageDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s sysDictLanguageDo) LeftJoin(table schema.Tabler, on ...field.Expr) ISysDictLanguageDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s sysDictLanguageDo) RightJoin(table schema.Tabler, on ...field.Expr) ISysDictLanguageDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s sysDictLanguageDo) Group(cols ...field.Expr) ISysDictLanguageDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s sysDictLanguageDo) Having(conds ...gen.Condition) ISysDictLanguageDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s sysDictLanguageDo) Limit(limit int) ISysDictLanguageDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s sysDictLanguageDo) Offset(offset int) ISysDictLanguageDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s sysDictLanguageDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ISysDictLanguageDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s sysDictLanguageDo) Unscoped() ISysDictLanguageDo {
	return s.withDO(s.DO.Unscoped())
}

func (s sysDictLanguageDo) Create(values ...*model.SysDictLanguage) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s sysDictLanguageDo) CreateInBatches(values []*model.SysDictLanguage, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s sysDictLanguageDo) Save(values ...*model.SysDictLanguage) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s sysDictLanguageDo) First() (*model.SysDictLanguage, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.SysDictLanguage), nil
	}
}

func (s sysDictLanguageDo) Take() (*model.SysDictLanguage, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.SysDictLanguage), nil
	}
}

func (s sysDictLanguageDo) Last() (*model.SysDictLanguage, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.SysDictLanguage), nil
	}
}

func (s sysDictLanguageDo) Find() ([]*model.SysDictLanguage, error) {
	result, err := s.DO.Find()
	return result.([]*model.SysDictLanguage), err
}

func (s sysDictLanguageDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.SysDictLanguage, err error) {
	buf := make([]*model.SysDictLanguage, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s sysDictLanguageDo) FindInBatches(result *[]*model.SysDictLanguage, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s sysDictLanguageDo) Attrs(attrs ...field.AssignExpr) ISysDictLanguageDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s sysDictLanguageDo) Assign(attrs ...field.AssignExpr) ISysDictLanguageDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s sysDictLanguageDo) Joins(fields ...field.RelationField) ISysDictLanguageDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s sysDictLanguageDo) Preload(fields ...field.RelationField) ISysDictLanguageDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s sysDictLanguageDo) FirstOrInit() (*model.SysDictLanguage, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.SysDictLanguage), nil
	}
}

func (s sysDictLanguageDo) FirstOrCreate() (*model.SysDictLanguage, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.SysDictLanguage), nil
	}
}

func (s sysDictLanguageDo) FindByPage(offset int, limit int) (result []*model.SysDictLanguage, count int64, err error) {
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

func (s sysDictLanguageDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s sysDictLanguageDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s sysDictLanguageDo) Delete(models ...*model.SysDictLanguage) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *sysDictLanguageDo) withDO(do gen.Dao) *sysDictLanguageDo {
	s.DO = *do.(*gen.DO)
	return s
}
