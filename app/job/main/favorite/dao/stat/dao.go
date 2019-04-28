package stat

import (
	"fmt"
	"time"

	"go-common/app/job/main/favorite/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
)

// Dao favorite dao.
type Dao struct {
	db          *sql.DB
	redis       *redis.Pool
	mc          *memcache.Pool
	redisExpire int
	ipExpire    int
	buvidExpire int
	mcExpire    int32
}

// New new a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.DB.Fav),
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int(time.Duration(c.Redis.Expire) / time.Second),
		ipExpire:    int(time.Duration(c.Redis.IPExpire) / time.Second),
		buvidExpire: int(time.Duration(c.Redis.BuvidExpire) / time.Second),
		// memcache
		mc:       memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
	}
	return
}

// hit .
func hit(id int64) (fid int64, table string) {
	fid = id / _folderStatSharding
	table = fmt.Sprintf("%02d", id%_folderStatSharding)
	return
}
