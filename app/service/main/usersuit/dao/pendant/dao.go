package pendant

import (
	"context"
	"time"

	"go-common/app/service/main/usersuit/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	"github.com/bluele/gcache"
)

const (
	_info = "/internal/v1/user/"
)

// Dao struct info of Dao.
type Dao struct {
	db *sql.DB

	c      *conf.Config
	client *bm.Client
	// redis
	redis         *redis.Pool
	pendantExpire int32
	// memcache
	mc          *memcache.Pool
	pointExpire int32
	vipInfoURL  string
	payURL      string
	notifyURL   string
	// equipStore
	equipStore gcache.Cache
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		db:     sql.NewMySQL(c.MySQL),
		client: bm.NewClient(c.HTTPClient),
		// redis
		redis:         redis.NewPool(c.Redis.Config),
		pendantExpire: int32(time.Duration(c.Redis.PendantExpire) / time.Second),
		// memcache
		mc:          memcache.NewPool(c.Memcache.Config),
		pointExpire: int32(time.Duration(c.Memcache.PointExpire) / time.Second),
		vipInfoURL:  c.VipURI + _info,
		payURL:      c.PayURL,
		notifyURL:   c.NotifyURL,
		equipStore:  gcache.New(c.EquipCache.Size).LFU().Build(),
	}
	return
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.pingRedis(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.redis != nil {
		d.redis.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
