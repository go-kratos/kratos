package archive

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_upFavTpsPrefix = "up_fav_tps_"
)

func keyUpFavTpsPrefix(mid int64) string {
	return _upFavTpsPrefix + strconv.FormatInt(mid, 10)
}

// FilenameExpires get filename expire time.
func (d *Dao) FilenameExpires(c context.Context, vs []*archive.VideoParam) (ves []*archive.VideoExpire, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	for _, v := range vs {
		conn.Send("GET", v.Filename)
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v) | vs(%#v)", err, vs)
		return
	}
	for _, v := range vs {
		var exp int64
		if exp, err = redis.Int64(conn.Receive()); err != nil && err != redis.ErrNil {
			log.Error("conn.Receive error(%+v) | filename(%s)", err, v.Filename)
			return
		}
		err = nil // NOTE: maybe err==redis.ErrNil
		ves = append(ves, &archive.VideoExpire{
			Filename: v.Filename,
			Expire:   exp,
		})
	}
	return
}

// FreshFavTypes fn
func (d *Dao) FreshFavTypes(c context.Context, mid int64, tp int) (err error) {
	var (
		conn  = d.redis.Get(c)
		score = time.Now().Unix()
	)
	defer conn.Close()
	if err = conn.Send("ZADD", keyUpFavTpsPrefix(mid), score, strconv.Itoa(tp)); err != nil {
		log.Error("conn.Send(ZADD, %s, %d) error(%v)", _upFavTpsPrefix, tp, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
