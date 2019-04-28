package dao

import (
	"context"

	"go-common/app/admin/main/push/conf"
	"go-common/library/database/orm"
	sqlx "go-common/library/database/sql"
	xhttp "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao struct user of Dao.
type Dao struct {
	c          *conf.Config
	db         *sqlx.DB
	DB         *gorm.DB
	httpClient *xhttp.Client
}

// New create a instance of Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		db:         sqlx.NewMySQL(c.MySQL),
		DB:         orm.NewMySQL(c.ORM),
		httpClient: xhttp.NewClient(c.HTTPClient),
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if defaultTableName == "push_business" {
			return defaultTableName
		}
		return "push_" + defaultTableName
	}
	d.DB.LogMode(true)
}

// Ping check connection of db , mc.
func (d *Dao) Ping(c context.Context) (err error) {
	if d.DB != nil {
		err = d.DB.DB().PingContext(c)
	}
	return
}

// Close close connection of db , mc.
func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}
