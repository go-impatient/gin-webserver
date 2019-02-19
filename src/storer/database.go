package storer

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"    // enable the mssql dialect
	_ "github.com/jinzhu/gorm/dialects/mysql"    // enable the mysql dialect
	_ "github.com/jinzhu/gorm/dialects/postgres" // enable the postgres dialect
	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/pkg/log"
	"github.com/pkg/errors"

	"github.com/moocss/go-webserver/src/model"
)

var Models = []interface{}{
	&model.User{},
}

type Database struct {
	Self *gorm.DB
	cfg  *config.ConfigDb
}

func NewDatabase(cfg *config.ConfigDb) *Database {
	return &Database{
		cfg: cfg,
	}
}

func (d *Database) Open() (*Database, error) {
	db, err := gorm.Open(d.cfg.Dialect, d.parseConnConfig())
	if err != nil {
		log.Errorf("Database connection failed.", err)
		return nil, err
	}

	// set for db connection
	db.DB().SetMaxOpenConns(d.cfg.MaxOpenConns) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(d.cfg.MaxIdleConns) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.DB().SetConnMaxLifetime(time.Second * time.Duration(d.cfg.ConnMaxLifeTime))
	db.LogMode(true)

	return &Database{Self: db}, nil
}

func (d *Database) Close() {
	err := d.Self.Close()
	if err != nil {
		log.Errorf("Disconnect from database failed: ", err)
	}
}

func (d *Database) GetTablePrefix() string {
	return d.cfg.TablePrefix
}

// migrate migrates database schemas ...
func (d *Database) Migrate() error {
	err := d.Self.AutoMigrate(Models...).Error
	if err != nil {
		return errors.Wrap(err, "auto migrate tables failed")
	}

	return nil
}

// creates necessary database tables
func (d *Database) CreateTables() error {
	for _, model := range Models {
		if !d.Self.HasTable(model) {
			err := d.Self.CreateTable(model).Error
			if err != nil {
				return errors.Wrap(err, "create table failed")
			}
		}
	}

	return nil
}

func realDSN(driver, dbname, username, password, addr, charset string) string {
	connStr := ""
	switch driver {
	case "mysql":
		connStr = fmt.Sprintf("%s:%s@%s/%s?charset=%s&parseTime=True&loc=Local", username, password, addr, dbname, charset)
	case "postgres":
		connStr = fmt.Sprintf("%s dbname=%s user=%s password=%s", addr, dbname, username, password)
	case "mssql":
		connStr = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", username, password, addr, dbname)
	}
	return connStr
}

func (d *Database) parseConnConfig() string {
	connHost := ""
	switch d.cfg.Dialect {
	case "mysql":
		if d.cfg.Unix != "" {
			connHost = fmt.Sprintf("unix(%s)", d.cfg.Unix)
		} else {
			connHost = fmt.Sprintf("tcp(%s:%s)", d.cfg.Host, d.cfg.Port)
		}
	case "postgres":
		connHost = fmt.Sprintf("host=%s port=%s", d.cfg.Host, d.cfg.Port)
	case "mssql":
		connHost = fmt.Sprintf("%s:%s", d.cfg.Host, d.cfg.Port)
	}
	s := realDSN(d.cfg.Dialect, d.cfg.Username, d.cfg.Password, connHost, d.cfg.DbName, d.cfg.Charset)
	return s
}
