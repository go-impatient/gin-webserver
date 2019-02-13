package storer

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

func (db *Database) Open() *gorm.DB {
	g, err := gorm.Open("mysql", parseConnConfig(db))
	if err != nil {
		log.Errorf("Database connection failed.", err)
	} else {
		log.Infof("Database connection succeed.")
	}

	// set for db connection
	g.LogMode(true)
	g.DB().SetMaxOpenConns(db.cfg.MaxOpenConns) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	g.DB().SetMaxIdleConns(db.cfg.MaxIdleConns) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	g.DB().SetConnMaxLifetime(time.Second * time.Duration(db.cfg.ConnMaxLifeTime))

	db.Self = g

	return g
}

func (db *Database) Close() {
	err := db.Self.Close()
	if err != nil {
		log.Error("Disconnect from database failed: ", err)
	}
}

func realDSN(dbname, username, password, addr, charset string) string {
	conf := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", username, password, addr, dbname, charset)
	return conf
}

func parseConnConfig(db *Database) string {
	var connHost string
	if db.cfg.Unix != "" {
		connHost = fmt.Sprintf("unix(%s)", db.cfg.Unix)
	} else {
		connHost = fmt.Sprintf("tcp(%s:%s)", db.cfg.Host, db.cfg.Port)
	}
	s := realDSN(db.cfg.Username, db.cfg.Password, connHost, db.cfg.DbName, db.cfg.Charset);
	return s
}

func (db *Database) GetTablePrefix() string {
	return db.cfg.TablePrefix
}
