package history

import (
	"context"
	"time"

	"go-common/app/interface/main/history/conf"
	eprpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	"go-common/library/cache"
	"go-common/library/cache/redis"
	"go-common/library/queue/databus"

	"go-common/library/database/hbase.v2"
)

// Dao dao.
type Dao struct {
	conf       *conf.Config
	info       *hbase.Client
	redis      *redis.Pool
	playPro    *databus.Databus
	merge      *databus.Databus
	experience *databus.Databus
	proPub     *databus.Databus
	delChan    *cache.Cache
	expire     int

	epidGRPC eprpc.EpisodeClient
}

// New new history dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:       c,
		info:       hbase.NewClient(c.Info.Config),
		redis:      redis.NewPool(c.Redis.Config),
		playPro:    databus.New(c.DataBus.PlayPro),
		merge:      databus.New(c.DataBus.Merge),
		proPub:     databus.New(c.DataBus.Pub),
		experience: databus.New(c.DataBus.Experience),
		delChan:    cache.New(1, 1024),
		expire:     int(time.Duration(c.Redis.Expire) / time.Second),
	}
	var err error
	if d.epidGRPC, err = eprpc.NewClient(nil); err != nil {
		panic(err)
	}
	return
}

// Ping check connection success.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.PingRedis(c)
}

// Close close the redis and kafka resource.
func (d *Dao) Close() {
	if d.redis != nil {
		d.redis.Close()
	}
	if d.playPro != nil {
		d.playPro.Close()
	}
	if d.merge != nil {
		d.merge.Close()
	}
	if d.experience != nil {
		d.experience.Close()
	}
	if d.proPub != nil {
		d.proPub.Close()
	}
}
