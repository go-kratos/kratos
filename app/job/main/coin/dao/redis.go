package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_loginShard    = 10000
	_loginedPrefix = "logined_%d_%d"
)

func loginKey(mid, day int64) string {
	return fmt.Sprintf(_loginedPrefix, day, mid/_loginShard)
}

// SetLogin set user logined,
func (d *Dao) SetLogin(c context.Context, mid, day int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("SETBIT", loginKey(mid, day), mid%_loginShard, 1); err != nil {
		PromError("redis:SetLogin")
		log.Error("d.SetLogin(%v,%v) redis: err: %v", mid, day, err)
		return
	}
	if err = conn.Send("EXPIRE", loginKey(mid, day), d.loginExpire); err != nil {
		PromError("redis:SetLogin")
		log.Error("d.SetLogin(%v,%v) redis: err: %v", mid, day, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("d.SetLogin(%v,%v) redis: err: %v", mid, day, err)
		PromError("redis:SetLogin")
		return
	}
	if _, err = redis.Bool(conn.Receive()); err != nil {
		log.Error("d.SetLogin(%v,%v) redis: err: %v", mid, day, err)
		PromError("redis:SetLogin")
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("d.SetLogin(%v,%v) redis: err: %v", mid, day, err)
		PromError("redis:SetLogin")
	}
	return
}

// Logined check if user logined.
func (d *Dao) Logined(c context.Context, mid, day int64) (b bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if b, err = redis.Bool(conn.Do("GETBIT", loginKey(mid, day), mid%_loginShard)); err == redis.ErrNil {
		err = nil
	}
	if err != nil {
		log.Error("d.Logined(%v,%v) redis err: %+v", mid, day, err)
	}
	return
}
