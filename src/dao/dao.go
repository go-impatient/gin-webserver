package dao

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/model"
	"github.com/moocss/go-webserver/src/storer"

	"github.com/moocss/go-webserver/src/pkg/log"
)

// Dao 对象
type Dao struct {
	orm *gorm.DB
	DB  *storer.DB
}

// New 实例化
func New(cfg *config.Config) *Dao {
	DbInstance := storer.NewDB(cfg.Db)
	if err := DbInstance.Open(); err != nil {
		panic(err)
	}

	dao := &Dao{
		DB:  DbInstance,
		orm: DbInstance.Self,
	}

	return dao
}

// Create 数据添加
func (d *Dao) Create(tableName string, data interface{}) bool {
	db := d.orm.Table(d.setTableName(tableName)).Create(data)
	if err := db.Error; err != nil {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

// FindMulti 数据复合查询
func (d *Dao) FindMulti(tableName string, data interface{}, query *model.QueryParam) bool {
	db := d.orm.Table(d.setTableName(tableName)).Offset(query.Offset)
	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}
	if query.Fields != "" {
		db = db.Select(query.Fields)
	}
	if query.Order != "" {
		db = db.Order(query.Order)
	}
	db = parseWhereParam(db, query.Where)
	db.Find(data)
	if err := db.Error; err != nil {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

// Count 统计某条件字段的条目数
func (d *Dao) Count(tableName string, count *int, query *model.QueryParam) bool {
	db := d.orm.Table(d.setTableName(tableName))
	db = parseWhereParam(db, query.Where)
	db = db.Count(count)
	if err := db.Error; err != nil {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

// FindOne 根据某一字段查询数据
func (d *Dao) FindOne(tableName string, data interface{}, query *model.QueryParam) bool {
	db := d.orm.Table(d.setTableName(tableName))
	if query.Fields != "" {
		db = db.Select(query.Fields)
	}
	db = parseWhereParam(db, query.Where)
	db = db.First(data)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

// FindById 根据ID查询数据
func (d *Dao) FindById(tableName string, data interface{}, id interface{}) bool {
	db := d.orm.Table(d.setTableName(tableName))
	db.First(data, id)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

// Update 更新数据
func (d *Dao) Update(tableName string, data interface{}, query *model.QueryParam) bool {
	db := d.orm.Table(d.setTableName(tableName))
	db = parseWhereParam(db, query.Where)
	db = db.Updates(data)
	if err := db.Error; err != nil {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

// Delete 删除数据
func (d *Dao) Delete(tableName string, data interface{}, query *model.QueryParam) bool {
	if len(query.Where) == 0 {
		log.Warn("mysql query error: delete failed, where conditions cannot be empty")
		return false
	}
	db := d.orm.Table(d.setTableName(tableName))
	db = parseWhereParam(db, query.Where)
	db = db.Delete(data)
	if err := db.Error; err != nil {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

// DeleteById 更加ID删除数据
func (d *Dao) DeleteById(tableName string, data interface{}) bool {
	db := d.orm.Table(d.setTableName(tableName))
	db.Delete(data)
	if err := db.Error; err != nil {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

// Where 条件拼装
func parseWhereParam(db *gorm.DB, where []model.WhereParam) *gorm.DB {
	if len(where) == 0 {
		return db
	}
	var (
		plain   []string
		prepare []interface{}
	)
	for _, w := range where {
		tag := w.Tag
		if tag == "" {
			tag = "="
		}
		var plainFmt string
		switch tag {
		case "IN":
			plainFmt = fmt.Sprintf("%s IN (?)", w.Field)
		default:
			plainFmt = fmt.Sprintf("%s %s ?", w.Field, tag)
		}
		plain = append(plain, plainFmt)
		prepare = append(prepare, w.Prepare)
	}
	return db.Where(strings.Join(plain, " AND "), prepare...)
}

// setTableName 给数据表拼接前缀
func (d *Dao) setTableName(rawName string) string {
	return strings.Join([]string{
		d.DB.GetTablePrefix(),
		rawName,
	}, "")
}
