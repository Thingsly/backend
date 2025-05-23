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

	"github.com/Thingsly/backend/internal/model"
)

func newDeviceTemplate(db *gorm.DB, opts ...gen.DOOption) deviceTemplate {
	_deviceTemplate := deviceTemplate{}

	_deviceTemplate.deviceTemplateDo.UseDB(db, opts...)
	_deviceTemplate.deviceTemplateDo.UseModel(&model.DeviceTemplate{})

	tableName := _deviceTemplate.deviceTemplateDo.TableName()
	_deviceTemplate.ALL = field.NewAsterisk(tableName)
	_deviceTemplate.ID = field.NewString(tableName, "id")
	_deviceTemplate.Name = field.NewString(tableName, "name")
	_deviceTemplate.Author = field.NewString(tableName, "author")
	_deviceTemplate.Version = field.NewString(tableName, "version")
	_deviceTemplate.Description = field.NewString(tableName, "description")
	_deviceTemplate.TenantID = field.NewString(tableName, "tenant_id")
	_deviceTemplate.CreatedAt = field.NewTime(tableName, "created_at")
	_deviceTemplate.UpdatedAt = field.NewTime(tableName, "updated_at")
	_deviceTemplate.Flag = field.NewInt16(tableName, "flag")
	_deviceTemplate.Label = field.NewString(tableName, "label")
	_deviceTemplate.WebChartConfig = field.NewString(tableName, "web_chart_config")
	_deviceTemplate.AppChartConfig = field.NewString(tableName, "app_chart_config")
	_deviceTemplate.Remark = field.NewString(tableName, "remark")
	_deviceTemplate.Path = field.NewString(tableName, "path")

	_deviceTemplate.fillFieldMap()

	return _deviceTemplate
}

type deviceTemplate struct {
	deviceTemplateDo

	ALL            field.Asterisk
	ID             field.String 
	Name           field.String 
	Author         field.String 
	Version        field.String 
	Description    field.String 
	TenantID       field.String
	CreatedAt      field.Time
	UpdatedAt      field.Time
	Flag           field.Int16  
	Label          field.String 
	WebChartConfig field.String 
	AppChartConfig field.String 
	Remark         field.String 
	Path           field.String 

	fieldMap map[string]field.Expr
}

func (d deviceTemplate) Table(newTableName string) *deviceTemplate {
	d.deviceTemplateDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d deviceTemplate) As(alias string) *deviceTemplate {
	d.deviceTemplateDo.DO = *(d.deviceTemplateDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *deviceTemplate) updateTableName(table string) *deviceTemplate {
	d.ALL = field.NewAsterisk(table)
	d.ID = field.NewString(table, "id")
	d.Name = field.NewString(table, "name")
	d.Author = field.NewString(table, "author")
	d.Version = field.NewString(table, "version")
	d.Description = field.NewString(table, "description")
	d.TenantID = field.NewString(table, "tenant_id")
	d.CreatedAt = field.NewTime(table, "created_at")
	d.UpdatedAt = field.NewTime(table, "updated_at")
	d.Flag = field.NewInt16(table, "flag")
	d.Label = field.NewString(table, "label")
	d.WebChartConfig = field.NewString(table, "web_chart_config")
	d.AppChartConfig = field.NewString(table, "app_chart_config")
	d.Remark = field.NewString(table, "remark")
	d.Path = field.NewString(table, "path")

	d.fillFieldMap()

	return d
}

func (d *deviceTemplate) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *deviceTemplate) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 14)
	d.fieldMap["id"] = d.ID
	d.fieldMap["name"] = d.Name
	d.fieldMap["author"] = d.Author
	d.fieldMap["version"] = d.Version
	d.fieldMap["description"] = d.Description
	d.fieldMap["tenant_id"] = d.TenantID
	d.fieldMap["created_at"] = d.CreatedAt
	d.fieldMap["updated_at"] = d.UpdatedAt
	d.fieldMap["flag"] = d.Flag
	d.fieldMap["label"] = d.Label
	d.fieldMap["web_chart_config"] = d.WebChartConfig
	d.fieldMap["app_chart_config"] = d.AppChartConfig
	d.fieldMap["remark"] = d.Remark
	d.fieldMap["path"] = d.Path
}

func (d deviceTemplate) clone(db *gorm.DB) deviceTemplate {
	d.deviceTemplateDo.ReplaceConnPool(db.Statement.ConnPool)
	return d
}

func (d deviceTemplate) replaceDB(db *gorm.DB) deviceTemplate {
	d.deviceTemplateDo.ReplaceDB(db)
	return d
}

type deviceTemplateDo struct{ gen.DO }

type IDeviceTemplateDo interface {
	gen.SubQuery
	Debug() IDeviceTemplateDo
	WithContext(ctx context.Context) IDeviceTemplateDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IDeviceTemplateDo
	WriteDB() IDeviceTemplateDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IDeviceTemplateDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IDeviceTemplateDo
	Not(conds ...gen.Condition) IDeviceTemplateDo
	Or(conds ...gen.Condition) IDeviceTemplateDo
	Select(conds ...field.Expr) IDeviceTemplateDo
	Where(conds ...gen.Condition) IDeviceTemplateDo
	Order(conds ...field.Expr) IDeviceTemplateDo
	Distinct(cols ...field.Expr) IDeviceTemplateDo
	Omit(cols ...field.Expr) IDeviceTemplateDo
	Join(table schema.Tabler, on ...field.Expr) IDeviceTemplateDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IDeviceTemplateDo
	RightJoin(table schema.Tabler, on ...field.Expr) IDeviceTemplateDo
	Group(cols ...field.Expr) IDeviceTemplateDo
	Having(conds ...gen.Condition) IDeviceTemplateDo
	Limit(limit int) IDeviceTemplateDo
	Offset(offset int) IDeviceTemplateDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IDeviceTemplateDo
	Unscoped() IDeviceTemplateDo
	Create(values ...*model.DeviceTemplate) error
	CreateInBatches(values []*model.DeviceTemplate, batchSize int) error
	Save(values ...*model.DeviceTemplate) error
	First() (*model.DeviceTemplate, error)
	Take() (*model.DeviceTemplate, error)
	Last() (*model.DeviceTemplate, error)
	Find() ([]*model.DeviceTemplate, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.DeviceTemplate, err error)
	FindInBatches(result *[]*model.DeviceTemplate, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.DeviceTemplate) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IDeviceTemplateDo
	Assign(attrs ...field.AssignExpr) IDeviceTemplateDo
	Joins(fields ...field.RelationField) IDeviceTemplateDo
	Preload(fields ...field.RelationField) IDeviceTemplateDo
	FirstOrInit() (*model.DeviceTemplate, error)
	FirstOrCreate() (*model.DeviceTemplate, error)
	FindByPage(offset int, limit int) (result []*model.DeviceTemplate, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IDeviceTemplateDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (d deviceTemplateDo) Debug() IDeviceTemplateDo {
	return d.withDO(d.DO.Debug())
}

func (d deviceTemplateDo) WithContext(ctx context.Context) IDeviceTemplateDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d deviceTemplateDo) ReadDB() IDeviceTemplateDo {
	return d.Clauses(dbresolver.Read)
}

func (d deviceTemplateDo) WriteDB() IDeviceTemplateDo {
	return d.Clauses(dbresolver.Write)
}

func (d deviceTemplateDo) Session(config *gorm.Session) IDeviceTemplateDo {
	return d.withDO(d.DO.Session(config))
}

func (d deviceTemplateDo) Clauses(conds ...clause.Expression) IDeviceTemplateDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d deviceTemplateDo) Returning(value interface{}, columns ...string) IDeviceTemplateDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d deviceTemplateDo) Not(conds ...gen.Condition) IDeviceTemplateDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d deviceTemplateDo) Or(conds ...gen.Condition) IDeviceTemplateDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d deviceTemplateDo) Select(conds ...field.Expr) IDeviceTemplateDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d deviceTemplateDo) Where(conds ...gen.Condition) IDeviceTemplateDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d deviceTemplateDo) Order(conds ...field.Expr) IDeviceTemplateDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d deviceTemplateDo) Distinct(cols ...field.Expr) IDeviceTemplateDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d deviceTemplateDo) Omit(cols ...field.Expr) IDeviceTemplateDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d deviceTemplateDo) Join(table schema.Tabler, on ...field.Expr) IDeviceTemplateDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d deviceTemplateDo) LeftJoin(table schema.Tabler, on ...field.Expr) IDeviceTemplateDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d deviceTemplateDo) RightJoin(table schema.Tabler, on ...field.Expr) IDeviceTemplateDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d deviceTemplateDo) Group(cols ...field.Expr) IDeviceTemplateDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d deviceTemplateDo) Having(conds ...gen.Condition) IDeviceTemplateDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d deviceTemplateDo) Limit(limit int) IDeviceTemplateDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d deviceTemplateDo) Offset(offset int) IDeviceTemplateDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d deviceTemplateDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IDeviceTemplateDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d deviceTemplateDo) Unscoped() IDeviceTemplateDo {
	return d.withDO(d.DO.Unscoped())
}

func (d deviceTemplateDo) Create(values ...*model.DeviceTemplate) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d deviceTemplateDo) CreateInBatches(values []*model.DeviceTemplate, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d deviceTemplateDo) Save(values ...*model.DeviceTemplate) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d deviceTemplateDo) First() (*model.DeviceTemplate, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTemplate), nil
	}
}

func (d deviceTemplateDo) Take() (*model.DeviceTemplate, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTemplate), nil
	}
}

func (d deviceTemplateDo) Last() (*model.DeviceTemplate, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTemplate), nil
	}
}

func (d deviceTemplateDo) Find() ([]*model.DeviceTemplate, error) {
	result, err := d.DO.Find()
	return result.([]*model.DeviceTemplate), err
}

func (d deviceTemplateDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.DeviceTemplate, err error) {
	buf := make([]*model.DeviceTemplate, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d deviceTemplateDo) FindInBatches(result *[]*model.DeviceTemplate, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d deviceTemplateDo) Attrs(attrs ...field.AssignExpr) IDeviceTemplateDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d deviceTemplateDo) Assign(attrs ...field.AssignExpr) IDeviceTemplateDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d deviceTemplateDo) Joins(fields ...field.RelationField) IDeviceTemplateDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d deviceTemplateDo) Preload(fields ...field.RelationField) IDeviceTemplateDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d deviceTemplateDo) FirstOrInit() (*model.DeviceTemplate, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTemplate), nil
	}
}

func (d deviceTemplateDo) FirstOrCreate() (*model.DeviceTemplate, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeviceTemplate), nil
	}
}

func (d deviceTemplateDo) FindByPage(offset int, limit int) (result []*model.DeviceTemplate, count int64, err error) {
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

func (d deviceTemplateDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d deviceTemplateDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d deviceTemplateDo) Delete(models ...*model.DeviceTemplate) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *deviceTemplateDo) withDO(do gen.Dao) *deviceTemplateDo {
	d.DO = *do.(*gen.DO)
	return d
}
