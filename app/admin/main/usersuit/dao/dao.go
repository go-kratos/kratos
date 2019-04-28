package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/usersuit/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

const (
	_msgURL          = "/api/notify/send.user.notify.do"
	_managersURI     = "/x/admin/manager/users"
	_managerTotalURI = "/x/admin/manager/users/total"
)

// Dao struct answer history of Dao
type Dao struct {
	db *sql.DB
	c  *conf.Config
	// redis
	redis *redis.Pool
	// http
	client          *bm.Client
	mc              *memcache.Pool
	mcExpire        int32
	pointExpire     int32
	msgURL          string
	managersURL     string
	managerTotalURL string
	bucket          string
	key             string
	secret          string
	bfs             string
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.DB.Usersuit),
		// http client
		client: bm.NewClient(c.HTTPClient),
		// redis
		redis:       redis.NewPool(c.Redis),
		mc:          memcache.NewPool(c.Memcache.Config),
		mcExpire:    int32(time.Duration(c.Memcache.Expire) / time.Second),
		pointExpire: int32(time.Duration(c.Memcache.PointExpire) / time.Second),
	}
	d.msgURL = c.Host.Msg + _msgURL
	d.managersURL = c.Host.Manager + _managersURI
	d.managerTotalURL = c.Host.Manager + _managerTotalURI
	d.bucket = c.BFS.Bucket
	d.key = c.BFS.Key
	d.secret = c.BFS.Secret
	d.bfs = c.Host.Bfs
	return
}

// BeginTran begin tran.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	tx, err = d.db.Begin(c)
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}
