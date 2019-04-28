package like

import (
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/log"

	"golang.org/x/net/context"
)

const (
	_prefixAttention = "lg_"
)

func redisKey(key string) string {
	return _prefixAttention + key
}

//RsSet set res
func (d *Dao) RsSet(c context.Context, key string, value string) (err error) {
	var (
		rkey = redisKey(key)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("SET", rkey, value); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", rkey, value, err)
		return
	}
	return
}

// RbSet setRb
func (d *Dao) RbSet(c context.Context, key string, value []byte) (err error) {
	var (
		rkey = redisKey(key)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("SET", rkey, value); err != nil {
		log.Error("conn.Send(SET, %s, %d) error(%v)", rkey, value, err)
		return
	}
	return
}

// RsGet getRs
func (d *Dao) RsGet(c context.Context, key string) (res string, err error) {
	var (
		rkey = redisKey(key)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.String(conn.Do("GET", rkey)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET key(%s)) error(%v)", rkey, err)
		}
		return
	}
	return
}

// RsSetNX NXset get
func (d *Dao) RsSetNX(c context.Context, key string) (res bool, err error) {
	var (
		rkey = redisKey(key)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("SETNX", rkey, 1)); err != nil {
		log.Error("conn.Do(SETNX key(%s)) error(%v)", rkey, err)
		return
	}
	return
}

// Incr incr
func (d *Dao) Incr(c context.Context, key string) (res bool, err error) {
	var (
		rkey = redisKey(key)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("incr", rkey)); err != nil {
		log.Error("conn.Do(incr key(%s)) error(%v)", rkey, err)
		return
	}
	return
}

// CreateSelection Create selection
func (d *Dao) CreateSelection(c context.Context, aid int64, stage int64) (err error) {
	key := strconv.FormatInt(aid, 10) + ":" + strconv.FormatInt(stage, 10)
	var (
		rkeyYes = redisKey(key + ":yes")
		rkeyNo  = redisKey(key + ":no")
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("SET", rkeyYes, 0); err != nil {
		log.Error("conn.Send(SET %s) error(%v)", rkeyYes, err)
		return
	}
	if err = conn.Send("SET", rkeyNo, 0); err != nil {
		log.Error("conn.Send(SET %s) error(%v)", rkeyNo, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// Selection selection
func (d *Dao) Selection(c context.Context, aid int64, stage int64) (yes int64, no int64, err error) {
	key := strconv.FormatInt(aid, 10) + ":" + strconv.FormatInt(stage, 10)
	var (
		rkeyYes = redisKey(key + ":yes")
		rkeyNo  = redisKey(key + ":no")
		conn    = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("GET", rkeyYes); err != nil {
		log.Error("conn.Send(SET %s) error(%v)", rkeyYes, err)
		return
	}
	if err = conn.Send("GET", rkeyNo); err != nil {
		log.Error("conn.Send(SET %s) error(%v)", rkeyNo, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	if yes, err = redis.Int64(conn.Receive()); err != nil {
		log.Error("conn.Receive(yes) error(%v)", err)
		return
	}
	if no, err = redis.Int64(conn.Receive()); err != nil {
		log.Error("conn.Receive(no) error(%v)", err)
		return
	}
	return
}
