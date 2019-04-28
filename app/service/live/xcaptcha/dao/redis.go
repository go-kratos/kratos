package dao

import (
	"context"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

// RedisIncr incr a key
func (d *Dao) RedisIncr(ctx context.Context, key string) (num int64, err error) {
	num = 0
	conn := d.redis.Get(ctx)
	defer conn.Close()
	err = conn.Send("INCR", key)
	if err != nil {
		log.Error("[XCaptcha][Redis][error] conn.Send error(%v)", err)
		return
	}
	err = conn.Send("EXPIRE", key, 3)
	if err != nil {
		log.Error("[XCaptcha][Redis][error] conn.Send error(%v)", err)
		return
	}
	err = conn.Flush()
	if err != nil {
		log.Error("[XCaptcha][Redis][error] conn.Flush error(%v)", err)
		return
	}
	if num, err = redis.Int64(conn.Receive()); err != nil {
		log.Error("[XCaptcha][Redis][error] INCR conn.Receive error(%v)", key, err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("[XCaptcha][Redis][error] EXPIRE conn.Receive error(%v)", key, err)
		return
	}
	return
}

// RedisGet get a string
func (d *Dao) RedisGet(ctx context.Context, key string) (value int64, err error) {
	value = 0
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if value, err = redis.Int64(conn.Do("GET", key)); err != nil {
		log.Error("[XCaptcha][Redis][error] GET conn.do error(%v)", key, err)
		return
	}
	return
}

// RedisSet Set a string and expire
func (d *Dao) RedisSet(ctx context.Context, key string, value int64, timeout int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", key, value, "EX", timeout); err != nil {
		log.Error("[XCaptcha][Redis][error] SET conn.do error(%v)", key, err)
		return
	}
	return
}

// RedisDel delete a key
func (d *Dao) RedisDel(ctx context.Context, key string) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("[XCaptcha][Redis][error] Delete conn.do error(%v)", key, err)
		return
	}
	return
}
