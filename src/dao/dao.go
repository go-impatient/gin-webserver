package dao

import(
	"github.com/jinzhu/gorm"
	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/storer"
)

type Dao struct {
	orm 	*gorm.DB
	db *storer.Database
}

func New(cfg *config.Config) *Dao {
	DbInstance := storer.NewDatabase(cfg.Db)
	if err := DbInstance.Open(); err != nil {
		panic(err)
	}

	d := &Dao{
		db: DbInstance,
		orm: DbInstance.Self,
	}

	defer DbInstance.Close()

	return d
}
