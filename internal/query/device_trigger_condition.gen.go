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

func newDeviceTriggerCondition(db *gorm.DB, opts ...gen.DOOption) deviceTriggerCondition {
	_deviceTriggerCondition := deviceTriggerCondition{}

	_deviceTriggerCondition.deviceTriggerConditionDo.UseDB(db, opts...)
	_deviceTriggerCondition.deviceTriggerConditionDo.UseModel(&model.DeviceTriggerCondition{})

	tableName := _deviceTriggerCondition.deviceTriggerConditionDo.TableName()
	_deviceTriggerCondition.ALL = field.NewAsterisk(tableName)
	_deviceTriggerCondition.ID = field.NewString(tableName, "id")
	_deviceTriggerCondition.SceneAutomationID = field.NewString(tableName, "scene_automation_id")
	_deviceTriggerCondition.Enabled = field.NewString(tableName, "enabled")
	_deviceTriggerCondition.GroupID = field.NewString(tableName, "group_id")
	_deviceTriggerCondition.TriggerConditionType = field.NewString(tableName, "trigger_condition_type")
	_deviceTriggerCondition.TriggerSource = field.NewString(tableName, "trigger_source")
	_deviceTriggerCondition.TriggerParamType = field.NewString(tableName, "trigger_param_type")
	_deviceTriggerCondition.TriggerParam = field.NewString(tableName, "trigger_param")
	_deviceTriggerCondition.TriggerOperator = field.NewString(tableName, "trigger_operator")
	_deviceTriggerCondition.TriggerValue = field.NewString(tableName, "trigger_value")
	_deviceTriggerCondition.Remark = field.NewString(tableName, "remark")
	_deviceTriggerCondition.TenantID = field.NewString(tableName, "tenant_id")

	_deviceTriggerCondition.fillFieldMap()

	return _deviceTriggerCondition
}

type deviceTriggerCondition struct {
	deviceTriggerConditionDo

	ALL                  field.Asterisk
	ID                   field.String 
	SceneAutomationID    field.String 
	Enabled              field.String 
	GroupID              field.String 
	TriggerConditionType field.String 
	TriggerSource        field.String 
	TriggerParamType     field.String 
	TriggerParam         field.String 
	TriggerOperator      field.String 
	TriggerValue         field.String 
	Remark               field.String
	TenantID             field.String 

	fieldMap map[string]field.Expr
}

func (d deviceTriggerCondition) Table(newTableName string) *deviceTriggerCondition {
	d.deviceTriggerConditionDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d deviceTriggerCondition) As(alias string) *deviceTriggerCondition {
	d.deviceTriggerConditionDo.DO = *(d.deviceTriggerConditionDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *deviceTriggerCondition) updateTableName(table string) *deviceTriggerCondition {
	d.ALL = field.NewAsterisk(table)
	d.ID = field.NewString(table, "id")
	d.SceneAutomationID = field.NewString(table, "scene_automation_id")
	d.Enabled = field.NewString(table, "enabled")
	d.GroupID = field.NewString(table, "group_id")
	d.TriggerConditionType = field.NewString(table, "trigger_condition_type")
	d.TriggerSource = field.NewString(table, "trigger_source")
	d.TriggerParamType = field.NewString(table, "trigger_param_type")
	d.TriggerParam = field.NewString(table, "trigger_param")
	d.TriggerOperator = field.NewString(table, "trigger_operator")
	d.TriggerValue = field.NewString(table, "trigger_value")
	d.Remark = field.NewString(table, "remark")
	d.TenantID = field.NewString(table, "tenant_id")

	d.fillFieldMap()

	return d
}

func (d *deviceTriggerCondition) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *deviceTriggerCondition) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 12)
	d.fieldMap["id"] = d.ID
	d.fieldMap["scene_automation_id"] = d.SceneAutomationID
	d.fieldMap["enabled"] = d.Enabled
	d.fieldMap["group_id"] = d.GroupID
	d.fieldMap["trigger_condition_type"] = d.TriggerConditionType
	d.fieldMap["trigger_source"] = d.TriggerSource
	d.fieldMap["trigger_param_type"] = d.TriggerParamType
	d.fieldMap["trigger_param"] = d.TriggerParam
	d.fieldMap["trigger_operator"] = d.TriggerOperator
	d.fieldMap["trigger_value"] = d.TriggerValue
	d.fieldMap["remark"] = d.Remark
	d.fieldMap["tenant_id"] = d.TenantID
}

func (d deviceTriggerCondition) clone(db *gorm.DB) deviceTriggerCondition {
	d.deviceTriggerConditionDo.ReplaceConnPool(db.Statement.ConnPool)
	return d
}

func (d deviceTriggerCondition) replaceDB(db *gorm.DB) deviceTriggerCondition {
	d.deviceTriggerConditionDo.ReplaceDB(db)
	return d
}

type deviceTriggerConditionDo struct{ gen.DO }

type IDeviceTriggerConditionDo interface {
	gen.SubQuery
	Debug() IDeviceTriggerConditionDo
	WithContext(ctx context.Context) IDeviceTriggerConditionDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IDeviceTriggerConditionDo
	WriteDB() IDeviceTriggerConditionDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IDeviceTriggerConditionDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IDeviceTriggerConditionDo
	Not(conds ...gen.Condition) IDeviceTriggerConditionDo
	Or(conds ...gen.Condition) IDeviceTriggerConditionDo
	Select(conds ...field.Expr) IDeviceTriggerConditionDo
	Where(conds ...gen.Condition) IDeviceTriggerConditionDo
	Order(conds ...field.Expr) IDeviceTriggerConditionDo
	Distinct(cols ...field.Expr) IDeviceTriggerConditionDo
	Omit(cols ...field.Expr) IDeviceTriggerConditionDo
	Join(table schema.Tabler, on ...field.Expr) IDeviceTriggerConditionDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IDeviceTriggerConditionDo
	RightJoin(table schema.Tabler, on ...field.Expr) IDeviceTriggerConditionDo
	Group(cols ...field.Expr) IDeviceTriggerConditionDo
	Having(conds ...gen.Condition) IDeviceTriggerConditionDo
	Limit(limit int) IDeviceTriggerConditionDo
	Offset(offset int) IDeviceTriggerConditionDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IDeviceTriggerConditionDo
	Unscoped() IDeviceTriggerConditionDo
	Create(values ...*model.DeviceTriggerCondition) error
	CreateInBatches(values []*model.DeviceTriggerCondition, batchSize int) error
	Save(values ...*model.DeviceTriggerCondition) error
	First() (*model.DeviceTriggerCondition, error)
	Take() (*model.DeviceTriggerCondition, error)
	Last() (*model.DeviceTriggerCondition, error)
	Find() ([]*model.DeviceTriggerCondition, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.DeviceTriggerCondition, err error)
	FindInBatches(result *[]*model.DeviceTriggerCondition, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.DeviceTriggerCondition) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IDeviceTriggerConditionDo
	Assign(attrs ...field.AssignExpr) IDeviceTriggerConditionDo
	Joins(fields ...field.RelationField) IDeviceTriggerConditionDo
	Preload(fields ...field.RelationField) IDeviceTriggerConditionDo
	FirstOrInit() (*model.DeviceTriggerCondition, error)
	FirstOrCreate() (*model.DeviceTriggerCondition, error)
	FindByPage(offset int, limit int) (result []*model.DeviceTriggerCondition, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IDeviceTriggerConditionDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (d deviceTriggerConditionDo) Debug() IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Debug())
}

func (d deviceTriggerConditionDo) WithContext(ctx context.Context) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d deviceTriggerConditionDo) ReadDB() IDeviceTriggerConditionDo {
	return d.Clauses(dbresolver.Read)
}

func (d deviceTriggerConditionDo) WriteDB() IDeviceTriggerConditionDo {
	return d.Clauses(dbresolver.Write)
}

func (d deviceTriggerConditionDo) Session(config *gorm.Session) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Session(config))
}

func (d deviceTriggerConditionDo) Clauses(conds ...clause.Expression) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d deviceTriggerConditionDo) Returning(value interface{}, columns ...string) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d deviceTriggerConditionDo) Not(conds ...gen.Condition) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d deviceTriggerConditionDo) Or(conds ...gen.Condition) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d deviceTriggerConditionDo) Select(conds ...field.Expr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d deviceTriggerConditionDo) Where(conds ...gen.Condition) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d deviceTriggerConditionDo) Order(conds ...field.Expr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d deviceTriggerConditionDo) Distinct(cols ...field.Expr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d deviceTriggerConditionDo) Omit(cols ...field.Expr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d deviceTriggerConditionDo) Join(table schema.Tabler, on ...field.Expr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d deviceTriggerConditionDo) LeftJoin(table schema.Tabler, on ...field.Expr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d deviceTriggerConditionDo) RightJoin(table schema.Tabler, on ...field.Expr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d deviceTriggerConditionDo) Group(cols ...field.Expr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d deviceTriggerConditionDo) Having(conds ...gen.Condition) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d deviceTriggerConditionDo) Limit(limit int) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d deviceTriggerConditionDo) Offset(offset int) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d deviceTriggerConditionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d deviceTriggerConditionDo) Unscoped() IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Unscoped())
}

func (d deviceTriggerConditionDo) Create(values ...*model.DeviceTriggerCondition) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d deviceTriggerConditionDo) CreateInBatches(values []*model.DeviceTriggerCondition, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d deviceTriggerConditionDo) Save(values ...*model.DeviceTriggerCondition) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d deviceTriggerConditionDo) First() (*model.DeviceTriggerCondition, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTriggerCondition), nil
	}
}

func (d deviceTriggerConditionDo) Take() (*model.DeviceTriggerCondition, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTriggerCondition), nil
	}
}

func (d deviceTriggerConditionDo) Last() (*model.DeviceTriggerCondition, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTriggerCondition), nil
	}
}

func (d deviceTriggerConditionDo) Find() ([]*model.DeviceTriggerCondition, error) {
	result, err := d.DO.Find()
	return result.([]*model.DeviceTriggerCondition), err
}

func (d deviceTriggerConditionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.DeviceTriggerCondition, err error) {
	buf := make([]*model.DeviceTriggerCondition, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d deviceTriggerConditionDo) FindInBatches(result *[]*model.DeviceTriggerCondition, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d deviceTriggerConditionDo) Attrs(attrs ...field.AssignExpr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d deviceTriggerConditionDo) Assign(attrs ...field.AssignExpr) IDeviceTriggerConditionDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d deviceTriggerConditionDo) Joins(fields ...field.RelationField) IDeviceTriggerConditionDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d deviceTriggerConditionDo) Preload(fields ...field.RelationField) IDeviceTriggerConditionDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d deviceTriggerConditionDo) FirstOrInit() (*model.DeviceTriggerCondition, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTriggerCondition), nil
	}
}

func (d deviceTriggerConditionDo) FirstOrCreate() (*model.DeviceTriggerCondition, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTriggerCondition), nil
	}
}

func (d deviceTriggerConditionDo) FindByPage(offset int, limit int) (result []*model.DeviceTriggerCondition, count int64, err error) {
	result, err = d.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = d.Offset(-1).Limit(-1).Count()
	return
}

func (d deviceTriggerConditionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d deviceTriggerConditionDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d deviceTriggerConditionDo) Delete(models ...*model.DeviceTriggerCondition) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *deviceTriggerConditionDo) withDO(do gen.Dao) *deviceTriggerConditionDo {
	d.DO = *do.(*gen.DO)
	return d
}
