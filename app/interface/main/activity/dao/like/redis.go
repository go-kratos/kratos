package like

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixAttention = "lg_"
)

func redisKey(key string) string {
	return _prefixAttention + key
}

// RsSet Dao
func (dao *Dao) RsSet(c context.Context, key string, value string) (ok bool, err error) {
	var (
		rkey = redisKey(key)
		conn = dao.redis.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("SET", rkey, value)); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", rkey, value, err)
		return
	}
	return
}

// RsGet Dao
func (dao *Dao) RsGet(c context.Context, key string) (res string, err error) {
	var (
		rkey = redisKey(key)
		conn = dao.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.String(conn.Do("GET", rkey)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET %s) error(%v)", rkey, err)
		}
	}
	return
}

// RsSetNX Dao
func (dao *Dao) RsSetNX(c context.Context, key string) (res bool, err error) {
	var (
		rkey = redisKey(key)
		conn = dao.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("SETNX", rkey, "1")); err != nil {
		if err == redis.ErrNil {
			log.Error("conn.Do(GET key(%s)) error(%v)", rkey, err)
			err = nil
		} else {
			log.Error("conn.Do(GET key(%s)) error(%v)", rkey, err)
			return
		}
	}
	fmt.Print(res)
	return
}

// Rb Dao
func (dao *Dao) Rb(c context.Context, key string) (res []byte, err error) {
	var (
		rkey = redisKey(key)
		conn = dao.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Bytes(conn.Do("GET", rkey)); err != nil {
		if err == redis.ErrNil {
			res = nil
			err = nil
		} else {
			log.Error("conn.Do(GET key(%v)) error(%v)", rkey, err)
		}
	}
	return
}

// Incr Dao
func (dao *Dao) Incr(c context.Context, key string) (res bool, err error) {
	var (
		rkey = redisKey(key)
		conn = dao.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("INCR", rkey)); err != nil {
		log.Error("conn.Do(INCR key(%s)) error(%v)", rkey, err)
	}
	return
}

// Incrby Dao
func (dao *Dao) Incrby(c context.Context, key string) (res bool, err error) {
	var (
		rkey = redisKey(key)
		conn = dao.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("INCRBY", rkey, 222)); err != nil {
		log.Error("conn.Do(INCRBY key(%s)) error(%v)", rkey, err)
	}
	return
}
