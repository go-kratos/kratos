package videoshot

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	// _hashKeyNum the max count of hash key
	_hashKeyNum = int64(100000)
	_keyPrefix  = "vs_"
)

func hashKey(cid int64) string {
	return fmt.Sprintf("%s%d", _keyPrefix, cid%_hashKeyNum)
}

// cache get videoshot's count by id.
func (d *Dao) cache(c context.Context, cid int64) (count, ver int, err error) {
	var (
		key  = hashKey(cid)
		conn = d.rds.Get(c)
		out  int64
	)
	defer conn.Close()
	if out, err = redis.Int64(conn.Do("HGET", key, cid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("HGET(%s, %d) error(%v)", key, cid, err)
		}
		return
	}
	ver = int(out >> 32)
	count = int(int32(out))
	return
}

// addCache set videoshot info into redis.
func (d *Dao) addCache(c context.Context, cid int64, ver, count int) (err error) {
	var (
		key  = hashKey(cid)
		conn = d.rds.Get(c)
		in   int64
	)
	in = int64(ver)<<32 | int64(count)
	defer conn.Close()
	if _, err = conn.Do("HSET", key, cid, in); err != nil {
		log.Error("HSET(%s, %d, %d)", key, cid, count, err)
	}
	return
}
