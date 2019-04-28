package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/growup/conf"
	article "go-common/app/interface/openplatform/article/rpc/client"
	account "go-common/app/service/main/account/rpc/client"
	vip "go-common/app/service/main/vip/rpc/client"

	"go-common/library/cache/redis"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Dao def dao struct
type Dao struct {
	c    *conf.Config
	db   *sql.DB
	rddb *sql.DB
	// redis
	redis       *redis.Pool
	redisExpire int64
	// search
	httpRead *bm.Client
	// rpc
	acc *account.Service3
	art *article.Service
	vip *vip.Service
	// hbase
	hbase        *hbase.Client
	hbaseTimeOut time.Duration
	// chan
	missch chan func()
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:    c,
		db:   sql.NewMySQL(c.DB.Growup),
		rddb: sql.NewMySQL(c.DB.Allowance),
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int64(time.Duration(c.Redis.Expire) / time.Second),
		// search
		httpRead: bm.NewClient(c.HTTPClient.Read),
		// rpc
		acc: account.New3(c.AccountRPC),
		art: article.New(c.ArticleRPC),
		vip: vip.New(c.VipRPC),
		// hbase
		hbase:        hbase.NewClient(c.HBase.Config),
		hbaseTimeOut: time.Duration(time.Millisecond * 200),
		// chan
		missch: make(chan func(), 1024),
	}
	go d.cacheproc()
	return d
}

// Ping ping db
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Error("d.db.Ping error(%v)", err)
		return
	}
	if err = d.pingRedis(c); err != nil {
		log.Error("d.pingRedis error(%v)", err)
	}
	return
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// Close close db conn
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	if d.hbase != nil {
		d.hbase.Close()
	}
}

// BeginTran begin transcation
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}

// AddCache add to chan for cache
func (d *Dao) AddCache(f func()) {
	select {
	case d.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for execute closure.
func (d *Dao) cacheproc() {
	for {
		f := <-d.missch
		f()
	}
}
