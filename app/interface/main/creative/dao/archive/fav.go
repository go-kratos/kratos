package archive

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_upFavTpsPrefix = "up_fav_tps_"
)

func keyUpFavTpsPrefix(mid int64) string {
	return _upFavTpsPrefix + strconv.FormatInt(mid, 10)
}

// FavTypes fn
func (d *Dao) FavTypes(c context.Context, mid int64) (items map[string]int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if items, err = redis.Int64Map(conn.Do("ZRANGE", keyUpFavTpsPrefix(mid), "0", "-1", "WITHSCORES")); err != nil {
		log.Error("redis.Int64Map(conn.Do(ZRANGE, %s, 0, -1)) error(%v)", keyUpFavTpsPrefix(mid), err)
	}
	return
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
