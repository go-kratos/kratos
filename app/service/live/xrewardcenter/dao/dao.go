package dao

import (
	"context"

	"go-common/app/service/live/xrewardcenter/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	gift_api "go-common/app/service/live/gift/api/liverpc"
	room_api "go-common/app/service/live/room/api/liverpc"
	"go-common/library/net/rpc/liverpc"
)

// Dao dao
type Dao struct {
	c     *conf.Config
	mc    *memcache.Pool
	redis *redis.Pool
	db    *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		mc:    memcache.NewPool(c.Memcache),
		redis: redis.NewPool(c.Redis),
		db:    xsql.NewMySQL(c.MySQL),
	}
	return
}

// RoomAPI .
var RoomAPI *room_api.Client

// GiftAPI .
var GiftAPI *gift_api.Client

// InitAPI init all service APIs
func InitAPI() {
	RoomAPI = room_api.New(getConf("room"))
	GiftAPI = gift_api.New(getConf("gift"))
}

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc
	if c != nil {
		return c[appName]
	}
	return nil
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	if err := d.db.Ping(c); err != nil {
		log.Error("ping db error(%v)", err)
		return err
	}

	return d.pingMemcache(c)
}

func (d *Dao) pingMemcache(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	err = conn.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0})

	if err != nil {
		log.Error("mc.ping.Store error(%v)", err)
		return err
	}
	return err
}
