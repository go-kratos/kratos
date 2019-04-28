package dao

import (
	"context"
	"go-common/app/service/openplatform/anti-fraud/conf"
	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

// Dao init dao
type Dao struct {
	c           *conf.Config
	db          *sql.DB     // db
	payShieldDb *sql.DB     // payShieldDb
	redis       *redis.Pool // redis
	client      *httpx.Client
	payData     chan *model.ShieldData
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		db:          sql.NewMySQL(c.DB.AntiFraud),
		payShieldDb: sql.NewMySQL(c.DB.PayShield),
		redis:       redis.NewPool(c.Redis.Config),
		client:      httpx.NewClient(conf.Conf.HTTPClient.Read),
		payData:     make(chan *model.ShieldData, 2000),
	}
	go d.SyncPayShield(context.Background())
	go d.PopAnswer(context.TODO())
	return
}

// Close close connections.
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
	if err = d.db.Ping(c); err != nil {
		log.Error("PingDb error(%v)", err)
		return
	}
	if err = d.PingRedis(c); err != nil {
		log.Error("PingRedis error(%v)", err)
		return
	}
	return
}

// BeginTran begin mysql transaction
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}
