package black

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

const (
	_prefixBlack = "b_"
)

// keyBlack is.
func keyBlack(mid int64) string {
	return _prefixBlack + strconv.FormatInt(mid, 10)
}

func (d *Dao) blackCache(c context.Context, mid int64) (aidm map[int64]struct{}, err error) {
	var aids []int64
	conn := d.redis.Get(c)
	key := keyBlack(mid)
	defer conn.Close()
	if aids, err = redis.Int64s(conn.Do("ZREVRANGE", key, 0, -1)); err != nil {
		err = errors.Wrapf(err, "conn.Do(ZREVRANGE,%s,0,-1)", key)
		return
	}
	aidm = make(map[int64]struct{}, len(aids))
	for _, aid := range aids {
		aidm[aid] = struct{}{}
	}
	return
}

// addBlackCache is.
func (d *Dao) addBlackCache(c context.Context, mid int64, aids ...int64) (err error) {
	if len(aids) == 0 {
		return
	}
	key := keyBlack(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, aid := range aids {
		if err = conn.Send("ZADD", key, aid, aid); err != nil {
			err = errors.Wrapf(err, "conn.Send(ZADD,%s,%d,%d)", key, aid, aid)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.expireRds); err != nil {
		err = errors.Wrapf(err, "conn.Send(EXPIRE,%s,%d)", key, d.expireRds)
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < len(aids)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return
}

// delBlackCache is.
func (d *Dao) delBlackCache(c context.Context, mid, aid int64) (err error) {
	key := keyBlack(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", key, aid); err != nil {
		err = errors.Wrapf(err, "conn.Do(ZREM,%s,%d)", key, aid)
	}
	return
}

// expireBlackCache is.
func (d *Dao) expireBlackCache(c context.Context, mid int64) (ok bool, err error) {
	key := keyBlack(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expireRds)); err != nil {
		err = errors.Wrapf(err, "conn.Do(EXPIRE,%s,%d)", key, d.expireRds)
	}
	return
}
