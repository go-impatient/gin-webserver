package storer

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // enable the mysql dialect
	_ "github.com/jinzhu/gorm/dialects/postgres" // enable the postgres dialect
	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/log"
	"time"
)

type Database struct {
	Self 	*gorm.DB
	cfg	 	*config.ConfigDb
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
	db.LogMode(true)
	db.DB().SetMaxOpenConns(d.cfg.MaxOpenConns) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(d.cfg.MaxIdleConns) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.DB().SetConnMaxLifetime(time.Second * time.Duration(d.cfg.ConnMaxLifeTime))

	return &Database{Self: db}, nil
}

func (d *Database) Close() {
	err := d.Self.Close()
	if err != nil {
		log.Error("Disconnect from database failed: ", err)
	}
}

func (d *Database) GetTablePrefix() string {
	return d.cfg.TablePrefix
}

func realDSN(dbname, username, password, addr, charset string) string {
	conf := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", username, password, addr, dbname, charset)
	return conf
}

func (d *Database) parseConnConfig() string {
	var connHost string
	if d.cfg.Unix != "" {
		connHost = fmt.Sprintf("unix(%s)", d.cfg.Unix)
	} else {
		connHost = fmt.Sprintf("tcp(%s:%s)", d.cfg.Host, d.cfg.Port)
	}
	s := realDSN(d.cfg.Username, d.cfg.Password, connHost, d.cfg.DbName, d.cfg.Charset);
	return s
}