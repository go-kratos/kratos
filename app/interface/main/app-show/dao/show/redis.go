package show

import (
	"context"

	"go-common/library/log"

	"go-common/library/cache/redis"
)

const (
	_prefix = "s_"
)

func keyRcmmd(mid string) string {
	return _prefix + mid
}

func keyCnt(mid string) string {
	return _prefix + mid + "_c"
}

// ExistRcmmndCache check recommend cache exists.
func (d *Dao) ExistRcmmndCache(c context.Context, mid string) (exist bool, err error) {
	conn := d.rcmmndRds.Get(c)
	defer conn.Close()
	key := keyCnt(mid)
	exist, err = redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Error("conn.Do(EXISTS, %s) error(%v)", key, err)
	}
	return
}

// AddRcmmndCache add recommend cache.
func (d *Dao) AddRcmmndCache(c context.Context, mid string, aids ...int64) (err error) {
	conn := d.rcmmndRds.Get(c)
	defer conn.Close()
	key := keyRcmmd(mid)
	cntk := keyCnt(mid)
	args := redis.Args{}.Add(key).AddFlat(aids)
	conn.Send("RPUSH", args...)
	conn.Send("EXPIRE", key, d.rcmmndExp)
	conn.Send("INCRBY", cntk, len(aids))
	conn.Send("EXPIRE", cntk, d.rcmmndExp)
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush err(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}

// PopRcmmndCache pop recommend cache.
func (d *Dao) PopRcmmndCache(c context.Context, mid string, cnt int) (aids []int64, err error) {
	conn := d.rcmmndRds.Get(c)
	defer conn.Close()
	key := keyRcmmd(mid)
	for i := 0; i < cnt; i++ {
		conn.Send("LPOP", key)
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	var aid int64
	for i := 0; i < cnt; i++ {
		aid, err = redis.Int64(conn.Receive())
		if err != nil {
			if err == redis.ErrNil {
				err = nil
				continue
			} else {
				log.Error("conn.Do(ZREVRANGE, %v)", err)
			}
			return
		}
		aids = append(aids, aid)
	}
	return
}
