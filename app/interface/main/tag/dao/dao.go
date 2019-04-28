package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/tag/conf"
	accrpc "go-common/app/service/main/account/api"
	tagrpc "go-common/app/service/main/tag/api"
	"go-common/library/cache/redis"
)

// Dao Dao.
type Dao struct {
	c            *conf.Config
	redis        *redis.Pool
	rankRedis    *redis.Pool
	expire       int
	restagExpire int
	rdsExpArc    int
	rdsExpOp     int
	rdsExpAtLike int
	rdsExpAtHate int
	// tagArc
	expireNewArc int
	tagArcMaxNum int

	// tag grpc client
	tagRPC tagrpc.TagRPCClient
	// account grpc client.
	accRPC accrpc.AccountClient
}

// New New.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		redis:        redis.NewPool(c.Redis.Tag.Redis),
		rankRedis:    redis.NewPool(c.Redis.Rank.Redis),
		expire:       int(time.Duration(c.Redis.Tag.Expire.Sub) / time.Second),
		restagExpire: int(time.Duration(c.Redis.Tag.Expire.ArcTag) / time.Second),

		//arctag
		rdsExpArc:    int(time.Duration(c.Redis.Tag.Expire.ArcTag) / time.Second),
		rdsExpOp:     int(time.Duration(c.Redis.Tag.Expire.ArcTagOp) / time.Second),
		rdsExpAtLike: int(time.Duration(c.Redis.Tag.Expire.AtLike) / time.Second),
		rdsExpAtHate: int(time.Duration(c.Redis.Tag.Expire.AtHate) / time.Second),

		expireNewArc: int(time.Duration(c.Redis.Rank.Expire.TagNewArc) / time.Second),
		tagArcMaxNum: c.Tag.MaxArcsLimit,
	}
	var err error
	if d.tagRPC, err = tagrpc.NewClient(c.TagGRPClient); err != nil {
		panic(err)
	}
	if d.accRPC, err = accrpc.NewClient(c.AccGRPClient); err != nil {
		panic(err)
	}
	return
}

// PingRe PingRe.
func (d *Dao) PingRe(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// Close dao close.
func (d *Dao) Close() {
	d.redis.Close()
	d.rankRedis.Close()
}
