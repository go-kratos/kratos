package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/conf/paladin"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_recInfoExpireTtl = 86400
	_recInfoKey       = "rec_info_%d"
)

// Dao dao.
type Dao struct {
	liveAppDb   *sql.DB
	redis       *redis.Pool
	redisExpire int32
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// New new a dao and return.
func New() (dao *Dao) {
	var (
		dc struct {
			LiveApp *sql.Config
		}
		rc struct {
			Rec        *redis.Config
			DemoExpire xtime.Duration
		}
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))
	fmt.Printf("%+v", dc)
	dao = &Dao{
		// mysql
		liveAppDb: sql.NewMySQL(dc.LiveApp),
		// redis
		redis:       redis.NewPool(rc.Rec),
		redisExpire: int32(time.Duration(rc.DemoExpire) / time.Second),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.liveAppDb.Close()
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	if err = d.pingRedis(ctx); err != nil {
		return
	}
	return d.liveAppDb.Ping(ctx)
}

func (d *Dao) pingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}
