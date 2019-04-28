package dao

import (
	"context"
	"go-common/app/service/bbq/topic/api"
	"go-common/library/sync/pipeline/fanout"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/conf/paladin"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -sync=true -batch=10 -max_group=10 -batch_err=break -nullcache=&api.VideoExtension{Svid:-1} -check_null_code=$==nil||$.Svid==-1
	VideoExtension(c context.Context, ids []int64) (map[int64]*api.VideoExtension, error)
	// cache: -sync=true -batch=10 -max_group=10 -batch_err=break -nullcache=&api.TopicInfo{TopicId:-1} -check_null_code=$==nil||$.TopicId==-1
	TopicInfo(c context.Context, ids []int64) (map[int64]*api.TopicInfo, error)
}

// Dao dao.
type Dao struct {
	cache       *fanout.Fanout
	db          *sql.DB
	redis       *redis.Pool
	topicExpire int32
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
			Topic *sql.Config
		}
		rc struct {
			Topic       *redis.Config
			TopicExpire xtime.Duration
		}
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))
	dao = &Dao{
		cache: fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		// mysql
		db: sql.NewMySQL(dc.Topic),
		// redis
		redis:       redis.NewPool(rc.Topic),
		topicExpire: int32(time.Duration(rc.TopicExpire) / time.Second),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	if err = d.pingRedis(ctx); err != nil {
		return
	}
	return d.db.Ping(ctx)
}

func (d *Dao) pingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}
