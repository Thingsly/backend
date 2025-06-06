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

func newOtaUpgradePackage(db *gorm.DB, opts ...gen.DOOption) otaUpgradePackage {
	_otaUpgradePackage := otaUpgradePackage{}

	_otaUpgradePackage.otaUpgradePackageDo.UseDB(db, opts...)
	_otaUpgradePackage.otaUpgradePackageDo.UseModel(&model.OtaUpgradePackage{})

	tableName := _otaUpgradePackage.otaUpgradePackageDo.TableName()
	_otaUpgradePackage.ALL = field.NewAsterisk(tableName)
	_otaUpgradePackage.ID = field.NewString(tableName, "id")
	_otaUpgradePackage.Name = field.NewString(tableName, "name")
	_otaUpgradePackage.Version = field.NewString(tableName, "version")
	_otaUpgradePackage.TargetVersion = field.NewString(tableName, "target_version")
	_otaUpgradePackage.DeviceConfigID = field.NewString(tableName, "device_config_id")
	_otaUpgradePackage.Module = field.NewString(tableName, "module")
	_otaUpgradePackage.PackageType = field.NewInt16(tableName, "package_type")
	_otaUpgradePackage.SignatureType = field.NewString(tableName, "signature_type")
	_otaUpgradePackage.AdditionalInfo = field.NewString(tableName, "additional_info")
	_otaUpgradePackage.Description = field.NewString(tableName, "description")
	_otaUpgradePackage.PackageURL = field.NewString(tableName, "package_url")
	_otaUpgradePackage.CreatedAt = field.NewTime(tableName, "created_at")
	_otaUpgradePackage.UpdatedAt = field.NewTime(tableName, "updated_at")
	_otaUpgradePackage.Remark = field.NewString(tableName, "remark")
	_otaUpgradePackage.Signature = field.NewString(tableName, "signature")
	_otaUpgradePackage.TenantID = field.NewString(tableName, "tenant_id")

	_otaUpgradePackage.fillFieldMap()

	return _otaUpgradePackage
}

type otaUpgradePackage struct {
	otaUpgradePackageDo

	ALL            field.Asterisk
	ID             field.String 
	Name           field.String 
	Version        field.String 
	TargetVersion  field.String 
	DeviceConfigID field.String 
	Module         field.String 
	PackageType    field.Int16  
	SignatureType  field.String 
	AdditionalInfo field.String 
	Description    field.String 
	PackageURL     field.String 
	CreatedAt      field.Time   
	UpdatedAt      field.Time   
	Remark         field.String 
	Signature      field.String 
	TenantID       field.String

	fieldMap map[string]field.Expr
}

func (o otaUpgradePackage) Table(newTableName string) *otaUpgradePackage {
	o.otaUpgradePackageDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o otaUpgradePackage) As(alias string) *otaUpgradePackage {
	o.otaUpgradePackageDo.DO = *(o.otaUpgradePackageDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *otaUpgradePackage) updateTableName(table string) *otaUpgradePackage {
	o.ALL = field.NewAsterisk(table)
	o.ID = field.NewString(table, "id")
	o.Name = field.NewString(table, "name")
	o.Version = field.NewString(table, "version")
	o.TargetVersion = field.NewString(table, "target_version")
	o.DeviceConfigID = field.NewString(table, "device_config_id")
	o.Module = field.NewString(table, "module")
	o.PackageType = field.NewInt16(table, "package_type")
	o.SignatureType = field.NewString(table, "signature_type")
	o.AdditionalInfo = field.NewString(table, "additional_info")
	o.Description = field.NewString(table, "description")
	o.PackageURL = field.NewString(table, "package_url")
	o.CreatedAt = field.NewTime(table, "created_at")
	o.UpdatedAt = field.NewTime(table, "updated_at")
	o.Remark = field.NewString(table, "remark")
	o.Signature = field.NewString(table, "signature")
	o.TenantID = field.NewString(table, "tenant_id")

	o.fillFieldMap()

	return o
}

func (o *otaUpgradePackage) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *otaUpgradePackage) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 16)
	o.fieldMap["id"] = o.ID
	o.fieldMap["name"] = o.Name
	o.fieldMap["version"] = o.Version
	o.fieldMap["target_version"] = o.TargetVersion
	o.fieldMap["device_config_id"] = o.DeviceConfigID
	o.fieldMap["module"] = o.Module
	o.fieldMap["package_type"] = o.PackageType
	o.fieldMap["signature_type"] = o.SignatureType
	o.fieldMap["additional_info"] = o.AdditionalInfo
	o.fieldMap["description"] = o.Description
	o.fieldMap["package_url"] = o.PackageURL
	o.fieldMap["created_at"] = o.CreatedAt
	o.fieldMap["updated_at"] = o.UpdatedAt
	o.fieldMap["remark"] = o.Remark
	o.fieldMap["signature"] = o.Signature
	o.fieldMap["tenant_id"] = o.TenantID
}

func (o otaUpgradePackage) clone(db *gorm.DB) otaUpgradePackage {
	o.otaUpgradePackageDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o otaUpgradePackage) replaceDB(db *gorm.DB) otaUpgradePackage {
	o.otaUpgradePackageDo.ReplaceDB(db)
	return o
}

type otaUpgradePackageDo struct{ gen.DO }

type IOtaUpgradePackageDo interface {
	gen.SubQuery
	Debug() IOtaUpgradePackageDo
	WithContext(ctx context.Context) IOtaUpgradePackageDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOtaUpgradePackageDo
	WriteDB() IOtaUpgradePackageDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOtaUpgradePackageDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOtaUpgradePackageDo
	Not(conds ...gen.Condition) IOtaUpgradePackageDo
	Or(conds ...gen.Condition) IOtaUpgradePackageDo
	Select(conds ...field.Expr) IOtaUpgradePackageDo
	Where(conds ...gen.Condition) IOtaUpgradePackageDo
	Order(conds ...field.Expr) IOtaUpgradePackageDo
	Distinct(cols ...field.Expr) IOtaUpgradePackageDo
	Omit(cols ...field.Expr) IOtaUpgradePackageDo
	Join(table schema.Tabler, on ...field.Expr) IOtaUpgradePackageDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IOtaUpgradePackageDo
	RightJoin(table schema.Tabler, on ...field.Expr) IOtaUpgradePackageDo
	Group(cols ...field.Expr) IOtaUpgradePackageDo
	Having(conds ...gen.Condition) IOtaUpgradePackageDo
	Limit(limit int) IOtaUpgradePackageDo
	Offset(offset int) IOtaUpgradePackageDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOtaUpgradePackageDo
	Unscoped() IOtaUpgradePackageDo
	Create(values ...*model.OtaUpgradePackage) error
	CreateInBatches(values []*model.OtaUpgradePackage, batchSize int) error
	Save(values ...*model.OtaUpgradePackage) error
	First() (*model.OtaUpgradePackage, error)
	Take() (*model.OtaUpgradePackage, error)
	Last() (*model.OtaUpgradePackage, error)
	Find() ([]*model.OtaUpgradePackage, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OtaUpgradePackage, err error)
	FindInBatches(result *[]*model.OtaUpgradePackage, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.OtaUpgradePackage) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOtaUpgradePackageDo
	Assign(attrs ...field.AssignExpr) IOtaUpgradePackageDo
	Joins(fields ...field.RelationField) IOtaUpgradePackageDo
	Preload(fields ...field.RelationField) IOtaUpgradePackageDo
	FirstOrInit() (*model.OtaUpgradePackage, error)
	FirstOrCreate() (*model.OtaUpgradePackage, error)
	FindByPage(offset int, limit int) (result []*model.OtaUpgradePackage, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOtaUpgradePackageDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (o otaUpgradePackageDo) Debug() IOtaUpgradePackageDo {
	return o.withDO(o.DO.Debug())
}

func (o otaUpgradePackageDo) WithContext(ctx context.Context) IOtaUpgradePackageDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o otaUpgradePackageDo) ReadDB() IOtaUpgradePackageDo {
	return o.Clauses(dbresolver.Read)
}

func (o otaUpgradePackageDo) WriteDB() IOtaUpgradePackageDo {
	return o.Clauses(dbresolver.Write)
}

func (o otaUpgradePackageDo) Session(config *gorm.Session) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Session(config))
}

func (o otaUpgradePackageDo) Clauses(conds ...clause.Expression) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o otaUpgradePackageDo) Returning(value interface{}, columns ...string) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o otaUpgradePackageDo) Not(conds ...gen.Condition) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o otaUpgradePackageDo) Or(conds ...gen.Condition) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o otaUpgradePackageDo) Select(conds ...field.Expr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o otaUpgradePackageDo) Where(conds ...gen.Condition) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o otaUpgradePackageDo) Order(conds ...field.Expr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o otaUpgradePackageDo) Distinct(cols ...field.Expr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o otaUpgradePackageDo) Omit(cols ...field.Expr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o otaUpgradePackageDo) Join(table schema.Tabler, on ...field.Expr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o otaUpgradePackageDo) LeftJoin(table schema.Tabler, on ...field.Expr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o otaUpgradePackageDo) RightJoin(table schema.Tabler, on ...field.Expr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o otaUpgradePackageDo) Group(cols ...field.Expr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o otaUpgradePackageDo) Having(conds ...gen.Condition) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o otaUpgradePackageDo) Limit(limit int) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o otaUpgradePackageDo) Offset(offset int) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o otaUpgradePackageDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o otaUpgradePackageDo) Unscoped() IOtaUpgradePackageDo {
	return o.withDO(o.DO.Unscoped())
}

func (o otaUpgradePackageDo) Create(values ...*model.OtaUpgradePackage) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o otaUpgradePackageDo) CreateInBatches(values []*model.OtaUpgradePackage, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o otaUpgradePackageDo) Save(values ...*model.OtaUpgradePackage) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o otaUpgradePackageDo) First() (*model.OtaUpgradePackage, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.OtaUpgradePackage), nil
	}
}

func (o otaUpgradePackageDo) Take() (*model.OtaUpgradePackage, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.OtaUpgradePackage), nil
	}
}

func (o otaUpgradePackageDo) Last() (*model.OtaUpgradePackage, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.OtaUpgradePackage), nil
	}
}

func (o otaUpgradePackageDo) Find() ([]*model.OtaUpgradePackage, error) {
	result, err := o.DO.Find()
	return result.([]*model.OtaUpgradePackage), err
}

func (o otaUpgradePackageDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OtaUpgradePackage, err error) {
	buf := make([]*model.OtaUpgradePackage, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o otaUpgradePackageDo) FindInBatches(result *[]*model.OtaUpgradePackage, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o otaUpgradePackageDo) Attrs(attrs ...field.AssignExpr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o otaUpgradePackageDo) Assign(attrs ...field.AssignExpr) IOtaUpgradePackageDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o otaUpgradePackageDo) Joins(fields ...field.RelationField) IOtaUpgradePackageDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o otaUpgradePackageDo) Preload(fields ...field.RelationField) IOtaUpgradePackageDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o otaUpgradePackageDo) FirstOrInit() (*model.OtaUpgradePackage, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.OtaUpgradePackage), nil
	}
}

func (o otaUpgradePackageDo) FirstOrCreate() (*model.OtaUpgradePackage, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.OtaUpgradePackage), nil
	}
}

func (o otaUpgradePackageDo) FindByPage(offset int, limit int) (result []*model.OtaUpgradePackage, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o otaUpgradePackageDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o otaUpgradePackageDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o otaUpgradePackageDo) Delete(models ...*model.OtaUpgradePackage) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *otaUpgradePackageDo) withDO(do gen.Dao) *otaUpgradePackageDo {
	o.DO = *do.(*gen.DO)
	return o
}
