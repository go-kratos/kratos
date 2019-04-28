package dao

import (
	"context"
	"fmt"

	"go-common/library/log"
)

// RedisDecr 指定 key 减去 num
func (d *Dao) RedisDecr(c context.Context, key string, num int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	if num == 1 {
		_, err = conn.Do("DECR", key)
	} else {
		_, err = conn.Do("DERCBY", key, num)
	}
	if err != nil {
		log.Error("d.RedisDecr(%s, %d) error(%v)", key, num, err)
	}
	return
}

// RedisDecrExist 当 key 存在时 给 key 减去指定数值 key 不存在时不做操作
func (d *Dao) RedisDecrExist(c context.Context, key string, num int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	lua := `if redis.call("EXISTS",KEYS[1])==1 then return redis.call("%s",KEYS[1]%s);else return nil;end`
	if num == 1 {
		log.Info(fmt.Sprintf(fmt.Sprintf(lua, "DECR", "")+"%d %s", 1, key))
		_, err = conn.Do("EVAL", fmt.Sprintf(lua, "DECR", ""), 1, key)
	} else {
		log.Info(fmt.Sprintf(fmt.Sprintf(lua, "DECRBY", ",ARGV[1]")+"%d %s %d"), 1, key, num)
		_, err = conn.Do("EVAL", fmt.Sprintf(lua, "DECRBY", ",ARGV[1]"), 1, key, num)
	}
	if err != nil {
		log.Error("d.RedisDecrExist(%s, %d) error(%v)", key, num, err)
	}
	return
}

// RedisDel del keys
func (d *Dao) RedisDel(c context.Context, key ...interface{}) (err error) {
	if len(key) == 0 {
		return
	}

	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key...); err != nil {
		log.Error("d.RedisDel(%v) error(%v)", key, err)
	}
	return
}

// RedisSetnx setnx
func (d *Dao) RedisSetnx(c context.Context, key string, val interface{}, ttl int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	if ttl > 0 {
		lua := `if redis.call("SETNX",KEYS[1],ARGV[1])==1 then return redis.call("EXPIRE",KEYS[1],ARGV[2]);else return 0;end'`
		_, err = conn.Do("EVAL", lua, 1, key, val, ttl)
	} else {
		_, err = conn.Do("SETNX", key, val)
	}
	if err != nil {
		log.Error("d.RedisSetnx(%s, %v, %d) error(%v)", key, val, ttl, err)
	}
	return
}
