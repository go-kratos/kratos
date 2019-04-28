package dao

import (
	"context"
	"go-common/app/service/live/gift/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/net/rpc/liverpc"

	activity "go-common/app/service/live/activity/api/liverpc"
	fans_medal "go-common/app/service/live/fans_medal/api/liverpc"
	live_user "go-common/app/service/live/live_user/api/liverpc"
	room "go-common/app/service/live/room/api/liverpc"
	user "go-common/app/service/live/user/api/liverpc"
	xuser "go-common/app/service/live/xuser/api/grpc/v1"
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

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return d.db.Ping(ctx)
}

var (
	// RoomApi RoomApi
	RoomApi *room.Client
	// LiveUserApi LiveUserApi
	LiveUserApi *live_user.Client
	// UserApi UserApi
	UserApi *user.Client
	// FansMedalApi FansMedalApi
	FansMedalApi *fans_medal.Client
	// ActivityApi ActivityApi
	ActivityApi *activity.Client
	// XuserClient XuserClient
	XuserClient *xuser.Client
)

//InitApi InitApi
func InitApi() {
	RoomApi = room.New(getConf("room"))
	LiveUserApi = live_user.New(getConf("live_user"))
	UserApi = user.New(getConf("user"))
	FansMedalApi = fans_medal.New(getConf("fans_medal"))
	ActivityApi = activity.New(getConf("activity"))
	var err error
	XuserClient, err = xuser.NewClient(nil)
	if err != nil {
		panic(err)
	}
}

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc
	if c != nil {
		return c[appName]
	}
	return nil
}
