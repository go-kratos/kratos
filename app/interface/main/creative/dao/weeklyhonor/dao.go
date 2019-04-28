package weeklyhonor

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/conf"
	up "go-common/app/service/main/up/api/v1"
	"go-common/library/cache/memcache"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao is data dao.
type Dao struct {
	c *conf.Config
	// hbase
	hbase        *hbase.Client
	hbaseTimeOut time.Duration
	// db
	db *sql.DB
	// mc
	mc            *memcache.Pool
	mcExpire      int32
	mcClickExpire int32
	// grpc
	upClient up.UpClient
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// hbase
		hbase:        hbase.NewClient(c.HBaseOld.Config),
		hbaseTimeOut: time.Duration(time.Millisecond * 200),
		// db
		db: sql.NewMySQL(c.DB.Creative),
		// mc
		mc:            memcache.NewPool(c.Memcache.Honor.Config),
		mcExpire:      int32(time.Duration(c.Memcache.Honor.HonorExpire) / time.Second),
		mcClickExpire: int32(time.Duration(c.Memcache.Honor.ClickExpire) / time.Second),
	}
	var err error
	if d.upClient, err = up.NewClient(c.UpClient); err != nil {
		panic(err)
	}
	return
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMySQL(c); err != nil {
		log.Error("s.pingMySQL.Ping err(%v)", err)
	}
	return
}

// Close hbase close
func (d *Dao) Close() (err error) {
	if d.hbase != nil {
		d.hbase.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
	return
}
