package dao

import (
	"context"
	"time"

	"go-common/app/service/main/usersuit/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_beFormal = "/api/member/beFormal"
)

// Dao struct answer history of Dao
type Dao struct {
	db    *sql.DB
	redis *redis.Pool
	//http
	httpClient   *bm.Client
	c            *conf.Config
	beFormalURI  string
	inviteExpire int32
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		db:           sql.NewMySQL(c.MySQL),
		redis:        redis.NewPool(c.Redis.Config),
		inviteExpire: int32(time.Duration(c.Redis.InviteExpire) / time.Second),
		httpClient:   bm.NewClient(c.HTTPClient),
		beFormalURI:  c.AccountIntranetURI + _beFormal,
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	d.db.Close()
	d.redis.Close()
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Error("dao.db.Ping() error(%v)", err)
		return
	}
	if err = d.pingRedis(c); err != nil {
		log.Error("dao.pingRedis() error(%v)", err)
	}
	return
}
