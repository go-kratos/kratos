package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/vip/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/orm"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"

	"github.com/jinzhu/gorm"
)

// Dao dao conf
type Dao struct {
	c        *conf.Config
	db       *xsql.DB
	vip      *gorm.DB
	client   *bm.Client
	mc       *memcache.Pool
	mcExpire int32
	errProm  *prom.Prom
}

// New init mysql db
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:        c,
		db:       xsql.NewMySQL(c.MySQL),
		vip:      orm.NewMySQL(c.ORM.Vip),
		client:   bm.NewClient(c.HTTPClient),
		mc:       memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
		errProm:  prom.BusinessErrCount,
	}
	d.initORM()
	return d
}

// Close close the resource.
func (d *Dao) Close() (err error) {
	return d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return d.db.Ping(c)
}

// BeginTran start tx .
func (d *Dao) BeginTran(c context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(c)
}

// BeginGormTran start gorm tx .
func (d *Dao) BeginGormTran(c context.Context) (tx *gorm.DB) {
	return d.vip.Begin()
}

func (d *Dao) initORM() {
	d.vip.LogMode(true)
}
