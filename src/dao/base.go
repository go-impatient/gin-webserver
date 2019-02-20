package dao

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/moocss/go-webserver/src/model"
	"github.com/moocss/go-webserver/src/pkg/log"
)

func (d *Dao) Create(tableName string, data interface{}) bool {
	db := d.orm.Table(d.setTableName(tableName)).Create(data)
	if err := db.Error; err != nil {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

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

func (d *Dao) FindById(tableName string, data interface{}, id interface{}) bool {
	db := d.orm.Table(d.setTableName(tableName))
	db.First(data, id)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

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

func (d *Dao) DeleteById(tableName string, data interface{}) bool {
	db := d.orm.Table(d.setTableName(tableName))
	db.Delete(data)
	if err := db.Error; err != nil {
		log.Warn("mysql query error: %v, sql[%v]", err, db.QueryExpr())
		return false
	}
	return true
}

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

func (d *Dao) setTableName(rawName string) string {
	return strings.Join([]string{
		d.db.GetTablePrefix(),
		rawName,
	}, "")
}
