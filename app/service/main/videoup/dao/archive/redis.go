package archive

import (
	"context"
	"fmt"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_videoJamTime = "va_v_jam_time"
	_editlock     = "edit_lock_aid_%d"
)

func lockKey(aid int64) string {
	return fmt.Sprintf(_editlock, aid)
}

// GetVideoJam get video traffic jam time
func (d *Dao) GetVideoJam(c context.Context) (seconds int, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()

	if seconds, err = redis.Int(conn.Do("GET", _videoJamTime)); err != nil {
		log.Error("conn.Do(GET,%s) error(%v)", _videoJamTime, err)
	}
	return
}

//SetNXLock redis lock.
func (d *Dao) SetNXLock(c context.Context, aid int64, times int64) (res bool, err error) {
	var (
		key  = lockKey(aid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("SETNX", key, "1")); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			res = true
			log.Error("conn.Do(SETNX(%s)) error(%v)", key, err)
			return
		}
	}
	if res {
		if _, err = redis.Bool(conn.Do("EXPIRE", key, times)); err != nil {
			log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, times, err)
			return
		}
	}
	return
}

//DelLock del lock.
func (d *Dao) DelLock(c context.Context, aid int64) (err error) {
	var (
		key  = lockKey(aid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(del,%v) err(%v)", key, err)
	}
	return
}
