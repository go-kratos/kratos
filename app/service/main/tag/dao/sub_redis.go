package dao

import (
	"context"
	"errors"
	"fmt"

	"go-common/app/service/main/tag/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// mid -> tids
const (
	_subkey = "hs_%d" // hashset  key:mid value:tid score:mtime
)

func (d *Dao) subKey(mid int64) string {
	return fmt.Sprintf(_subkey, mid)
}

// ExpireSubCache .
func (d *Dao) ExpireSubCache(c context.Context, mid int64) (ok bool, err error) {
	key := d.subKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.subExpire)); err != nil {
		log.Error("conn.Do(EXPIRE,%s), error(%v)", key, err)
	}
	return
}

// AddSubMapCache .
func (d *Dao) AddSubMapCache(c context.Context, mid int64, subs map[int64]*model.Sub) (err error) {
	key := d.subKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HDEL", key, -1); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	for _, sub := range subs {
		if err = conn.Send("HSET", key, sub.Tid, sub.MTime.Time().Unix()); err != nil {
			log.Error("conn.Send error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.subExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(subs)+2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddSubListCache .
func (d *Dao) AddSubListCache(c context.Context, mid int64, subs []*model.Sub) (err error) {
	key := d.subKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HDEL", key, -1); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	for _, sub := range subs {
		if err = conn.Send("HSET", key, sub.Tid, sub.MTime.Time().Unix()); err != nil {
			log.Error("conn.Send error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.subExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(subs)+2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelSubTidCache .
func (d *Dao) DelSubTidCache(c context.Context, mid, tid int64) (err error) {
	key := d.subKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("HDEL", key, tid); err != nil {
		log.Error("conn.Do(%s,%d) error(%v)", key, tid, err)
	}
	return
}

// DelSubTidsCache .
func (d *Dao) DelSubTidsCache(c context.Context, mid int64, tids []int64) (err error) {
	key := d.subKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, tid := range tids {
		if err = conn.Send("HDEL", key, tid); err != nil {
			log.Error("conn.Send error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.subExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(tids); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelSubCache .
func (d *Dao) DelSubCache(c context.Context, mid int64) (err error) {
	key := d.subKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(%s) error(%v)", key, err)
	}
	return
}

// IsSubCache .
func (d *Dao) IsSubCache(c context.Context, mid int64, tid int64) (ok bool, err error) {
	var unix int64
	conn := d.redis.Get(c)
	defer conn.Close()
	if unix, err = redis.Int64(conn.Do("HGET", d.subKey(mid), tid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("redis.Int64Map error(%v)", err)
		}
		return
	}
	if unix > 0 {
		ok = true
	}
	return
}

// IsSubsCache .
func (d *Dao) IsSubsCache(c context.Context, mid int64, tids []int64) (res map[int64]int32, err error) {
	var (
		values [][]byte
		key    = d.subKey(mid)
		args   = []interface{}{key}
	)
	res = make(map[int64]int32, len(tids))
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, tid := range tids {
		args = append(args, tid)
	}
	values, err = redis.ByteSlices(conn.Do("HMGET", args...))
	if err != nil {
		log.Error("redis.Int64Map error(%v)", err)
		return
	}
	for i, v := range tids {
		if len(values[i]) > 0 {
			res[v] = 1
		}
	}
	return
}

// SubTidsCache .
func (d *Dao) SubTidsCache(c context.Context, mid int64) (res map[int64]int32, err error) {
	var (
		key    = d.subKey(mid)
		conn   = d.redis.Get(c)
		values []interface{}
		ok     bool
		k      []byte
		tid    int64
	)
	defer conn.Close()
	if values, err = redis.Values(conn.Do("HGETALL", key)); err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: Int64Map expects even number of values result")
	}
	res = make(map[int64]int32, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		k, ok = values[i].([]byte)
		if !ok {
			return nil, errors.New("redigo: ScanMap key not a bulk string value")
		}
		tid, err = redis.Int64(k, nil)
		if err != nil {
			return
		}
		res[tid] = 1
	}
	return
}

// SubCache .
func (d *Dao) SubCache(c context.Context, mid int64) (subs []*model.Sub, rem map[int64]*model.Sub, err error) {
	var (
		ok     bool
		key    []byte
		tid, v int64
		values []interface{}
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err = redis.Values(conn.Do("HGETALL", d.subKey(mid)))
	if err != nil {
		return nil, nil, err
	}
	if len(values)%2 != 0 {
		return nil, nil, errors.New("redigo: Int64Map expects even number of values result")
	}
	subs = make([]*model.Sub, 0, len(values)/2)
	rem = make(map[int64]*model.Sub, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok = values[i].([]byte)
		if !ok {
			return nil, nil, errors.New("redigo: ScanMap key not a bulk string value")
		}
		tid, err = redis.Int64(key, nil)
		if err != nil {
			return
		}
		if tid <= 0 {
			continue
		}
		v, err = redis.Int64(values[i+1], nil)
		if err != nil {
			return
		}
		r := &model.Sub{
			Tid:   tid,
			MTime: xtime.Time(v),
		}
		subs = append(subs, r)
		rem[r.Tid] = r
	}
	return
}

// SubListCache .
func (d *Dao) SubListCache(c context.Context, mid int64) (subs []*model.Sub, err error) {
	var (
		ok     bool
		key    []byte
		tid, v int64
		values []interface{}
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err = redis.Values(conn.Do("HGETALL", d.subKey(mid)))
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: Int64Map expects even number of values result")
	}
	subs = make([]*model.Sub, 0, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok = values[i].([]byte)
		if !ok {
			return nil, errors.New("redigo: ScanMap key not a bulk string value")
		}
		tid, err = redis.Int64(key, nil)
		if err != nil {
			return
		}
		if tid <= 0 {
			continue
		}
		v, err = redis.Int64(values[i+1], nil)
		if err != nil {
			return
		}
		subs = append(subs, &model.Sub{Tid: tid, MTime: xtime.Time(v)})
	}
	return
}
