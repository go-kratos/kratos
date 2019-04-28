package like

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

const (
	_redDotKeyFmt = "rd_%d"
)

func keyRedDot(mid int64) string {
	return fmt.Sprintf(_redDotKeyFmt, mid)
}

// CacheRedDotTs .
func (d *Dao) CacheRedDotTs(c context.Context, mid int64) (ts int64, err error) {
	key := keyRedDot(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if ts, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			err = errors.Wrapf(err, "conn.Do(GET, %s)", key)
		}
	}
	return
}

// AddCacheRedDotTs .
func (d *Dao) AddCacheRedDotTs(c context.Context, mid, ts int64) (err error) {
	key := keyRedDot(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("SET", key, ts); err != nil {
		err = errors.Wrapf(err, "conn.Send(SET, %s, %d)", key, ts)
		return
	}
	if err = conn.Send("EXPIRE", key, d.hotDotExpire); err != nil {
		err = errors.Wrapf(err, "conn.Send(EXPIRE, %s, %d)", key, ts)
		return
	}
	if err = conn.Flush(); err != nil {
		err = errors.Wrap(err, "AddCacheHotDotTs conn.Flush")
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			err = errors.Wrapf(err, "add conn.Receive(%d)", i+1)
			return
		}
	}
	return
}
