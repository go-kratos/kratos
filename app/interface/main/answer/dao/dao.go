package dao

import (
	"time"

	"go-common/app/interface/main/answer/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Dao struct answer history of Dao
type Dao struct {
	c                 *conf.Config
	db                *sql.DB
	redis             *redis.Pool
	mc                *memcache.Pool
	mcExpire          int32
	redisExpire       int32
	ansCountExpire    int32
	ansAddFlagExpire  int32
	answerBlockExpire int32
	dbExtraAnswerRet  *databus.Databus
	captcha           *httpx.Client
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                 c,
		db:                sql.NewMySQL(c.Mysql),
		redis:             redis.NewPool(c.Redis.Config),
		mc:                memcache.NewPool(c.Memcache.Config),
		redisExpire:       int32(time.Duration(c.Redis.Expire) / time.Second),
		ansCountExpire:    int32(time.Duration(c.Redis.AnsCountExpire) / time.Second),
		ansAddFlagExpire:  int32(time.Duration(c.Redis.AnsCountExpire) / time.Second),
		mcExpire:          int32(time.Duration(c.Memcache.Expire) / time.Second),
		answerBlockExpire: int32(time.Duration(c.Memcache.AnswerBolckExpire) / time.Second),
		dbExtraAnswerRet:  databus.New(c.DataBus.ExtraAnswer),
		captcha:           httpx.NewClient(c.Captcha),
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.redis != nil {
		d.redis.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
}
