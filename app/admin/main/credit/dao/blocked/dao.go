package blocked

import (
	"context"

	"go-common/app/admin/main/credit/conf"
	"go-common/library/database/orm"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao struct info of Dao.
type Dao struct {
	// mysql
	ReadDB *gorm.DB
	DB     *gorm.DB
	// http
	client *bm.Client
	// conf
	c *conf.Config
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// http client
		client: bm.NewClient(c.HTTPClient),
		ReadDB: orm.NewMySQL(c.ORM.Read),
		DB:     orm.NewMySQL(c.ORM.Write),
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	d.ReadDB.LogMode(true)
	d.DB.LogMode(true)
}

// Close close dao.
func (d *Dao) Close() {
	if d.ReadDB != nil {
		d.ReadDB.Close()
	}
	if d.DB != nil {
		d.DB.Close()
	}
}

// Ping check connection of db , mc.
func (d *Dao) Ping(c context.Context) (err error) {
	if d.ReadDB != nil {
		err = d.ReadDB.DB().PingContext(c)
		return
	}
	if d.DB != nil {
		err = d.DB.DB().PingContext(c)
	}
	return
}
