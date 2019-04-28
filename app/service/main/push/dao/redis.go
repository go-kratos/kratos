package dao

import (
	"context"

	"go-common/library/log"
)

// pingRedis ping redis.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		PromError("redis: ping remote")
		log.Error("remote redis: conn.Do(SET,PING,PONG) error(%v)", err)
	}
	return
}
