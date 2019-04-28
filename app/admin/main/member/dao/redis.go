package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/net/ip"
)

func reviewAuditNotifyLockKey(t time.Time) string {
	prefix := t.Format("2006-01-02-15")
	if t.Minute() < 30 {
		return fmt.Sprintf("review_notify_%s_00", prefix)
	}
	return fmt.Sprintf("review_notify_%s_30", prefix)
}

func (d *Dao) pingRedis(c context.Context) error {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err := conn.Do("SET", "ping", "pong")
	return err
}

// TryLockReviewNotify is
func (d *Dao) TryLockReviewNotify(c context.Context, t time.Time) (bool, error) {
	key := reviewAuditNotifyLockKey(t)
	conn := d.redis.Get(c)
	defer conn.Close()

	locked, err := redis.Bool(conn.Do("SETNX", key, fmt.Sprintf("%s::%s", ip.InternalIP(), t)))
	if err != nil {
		return false, err
	}
	if !locked {
		return false, nil
	}
	if _, err := conn.Do("EXPIRE", key, 60*60); err != nil {
		log.Error("Failed to set expire on key: %s: %+v", key, err)
		// return
	}
	return locked, nil
}
