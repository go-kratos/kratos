package dao

import (
	"context"
	"time"

	"go-common/app/job/main/member/conf"
	"go-common/app/job/main/member/dao/block"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Dao struct info of Dao.
type Dao struct {
	c           *conf.Config
	block       *block.Dao
	db          *sql.DB
	accCheckDB  *sql.DB
	passLogDB   *sql.DB
	accdb       *sql.DB
	asodb       *sql.DB
	client      *bm.Client
	mc          *memcache.Pool
	mcExpire    int32
	redis       *redis.Pool
	plogDatabus *databus.Databus
	accNotify   *databus.Databus
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		db:          sql.NewMySQL(c.Mysql),
		accCheckDB:  sql.NewMySQL(c.AccCheckMysql),
		passLogDB:   sql.NewMySQL(c.PasslogMysql),
		accdb:       sql.NewMySQL(c.AccMysql),
		asodb:       sql.NewMySQL(c.AsoMysql),
		client:      bm.NewClient(c.HTTPClient),
		mc:          memcache.NewPool(c.Memcache.Config),
		mcExpire:    int32(time.Duration(c.Memcache.Expire) / time.Second),
		redis:       redis.NewPool(c.Redis),
		plogDatabus: databus.New(c.PLogDatabus),
		accNotify:   databus.New(c.AccountNotify),
	}
	d.block = block.New(c,
		memcache.NewPool(c.BlockMemcache),
		sql.NewMySQL(c.BlockDB),
		d.client,
		d.NotifyPurgeCache,
	)
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
	if d.block != nil {
		d.block.Close()
	}
}

// Ping ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// BlockImpl is.
func (d *Dao) BlockImpl() *block.Dao {
	return d.block
}
