package guard

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

// redis cache
const (
	_lockKey          = "saveGuard:%s"
	_cacheKeyUid      = "live_user:guard:uid:v1:%d"
	_cacheKeyTargetId = "live_user:guard:target_id:v1:%d"
)

// LockOrder lock for same order
func (d *GuardDao) LockOrder(ctx context.Context, orderID string) (ok bool, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = redis.String(conn.Do("SET", fmt.Sprintf(_lockKey, orderID), 1, "EX", 3*86400, "NX"))
	if err == redis.ErrNil {
		log.Info("LockOrder(%s) is ErrNil!", orderID)
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// UnlockOrder release lock for same order
func (d *GuardDao) UnlockOrder(ctx context.Context, orderID string) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = redis.String(conn.Do("DEL", fmt.Sprintf(_lockKey, orderID)))
	return
}

// ClearCache delete cache for guard
func (d *GuardDao) ClearCache(ctx context.Context, uid int64, ruid int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = redis.String(conn.Do("DEL", fmt.Sprintf(_cacheKeyUid, uid)))
	_, err = redis.String(conn.Do("DEL", fmt.Sprintf(_cacheKeyTargetId, ruid)))
	return
}
