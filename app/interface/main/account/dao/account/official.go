package account

import (
	"context"
	"fmt"
	"time"

	"go-common/library/cache/redis"
)

const (
	_oneMonth = 31 * 24 * 60 * 60
)

func keyMonthlyOfficialSubmittedTimes(t time.Time, mid int64) string {
	return fmt.Sprintf("ot_%d_%d", t.Month(), mid)
}

// IncreaseMonthlyOfficialSubmittedTimes is
func (d *Dao) IncreaseMonthlyOfficialSubmittedTimes(ctx context.Context, mid int64) (int64, error) {
	key := keyMonthlyOfficialSubmittedTimes(time.Now(), mid)
	conn := d.redis.Get(ctx)
	defer conn.Close()

	conn.Send("INCR", key)
	conn.Send("EXPIRE", key, _oneMonth)
	if err := conn.Flush(); err != nil {
		return 0, err
	}

	new, err := redis.Int64(conn.Receive())
	if err != nil {
		return 0, err
	}
	conn.Receive() // drain the pipe line

	return new, nil
}

// GetMonthlyOfficialSubmittedTimes is
func (d *Dao) GetMonthlyOfficialSubmittedTimes(ctx context.Context, mid int64) (int64, error) {
	key := keyMonthlyOfficialSubmittedTimes(time.Now(), mid)
	conn := d.redis.Get(ctx)
	defer conn.Close()

	v, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return 0, nil
		}
		return 0, err
	}
	return v, nil
}
