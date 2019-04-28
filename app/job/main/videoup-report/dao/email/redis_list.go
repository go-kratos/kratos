package email

import (
	"context"
	"encoding/json"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

// PushRedis rpush fail item to redis
func (d *Dao) PushRedis(c context.Context, a interface{}, key string) (err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(a); err != nil {
		log.Error("json.Marshal(%v) error(%v) key(%s)", a, err, key)
		return
	}
	if _, err = conn.Do("RPUSH", key, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s, %s) error(%v)", key, bs, err)
	}
	return
}

// PopRedis lpop fail item from redis
func (d *Dao) PopRedis(c context.Context, key string) (bs []byte, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()

	if bs, err = redis.Bytes(conn.Do("LPOP", key)); err != nil && err != redis.ErrNil {
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%v)", key, err)
	}
	return
}

//RemoveRedis lrem an element from redis list
func (d *Dao) RemoveRedis(c context.Context, key string, member ...interface{}) (reply int, err error) {
	var (
		lmem = len(member)
		conn redis.Conn
	)

	if lmem < 1 {
		return
	}

	conn = d.redis.Get(c)
	defer conn.Close()

	if lmem == 1 {
		reply, err = redis.Int(conn.Do("LREM", key, 0, member[0]))
	} else {
		lua := "local a=0;for k in pairs(ARGV) do a=a+redis.call('LREM',KEYS[1],0,ARGV[k]) end;return a;"
		args := []interface{}{lua, 1, key}
		args = append(args, member...)
		reply, err = redis.Int(conn.Do("EVAL", args...))
	}

	if err != nil {
		log.Error("RemoveRedis conn.Do(%s) member(%+v) error(%v)", key, member, err)
	}
	return
}
