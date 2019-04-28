package dao

import (
	"context"
	"time"

	"go-common/library/cache/redis"
)

// AddRctFollower is
func (d *Dao) AddRctFollower(c context.Context, mid, fid int64) error {
	key := recentFollower(fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err := conn.Send("ZADD", key, time.Now().Unix(), mid); err != nil {
		return err
	}
	if err := conn.Send("EXPIRE", key, d.UnreadDuration); err != nil {
		return err
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	return nil
}

// DelRctFollower is
func (d *Dao) DelRctFollower(c context.Context, mid, fid int64) error {
	key := recentFollower(fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err := conn.Do("ZREM", key, mid)
	return err
}

// RctFollowerCount is
func (d *Dao) RctFollowerCount(ctx context.Context, fid int64) (int64, error) {
	key := recentFollower(fid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	count, err := redis.Int64(conn.Do("ZCARD", key))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// EmptyRctFollower is
func (d *Dao) EmptyRctFollower(ctx context.Context, fid int64) error {
	key := recentFollower(fid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

// RctFollowerNotify is
func (d *Dao) RctFollowerNotify(c context.Context, fid int64) (bool, error) {
	key := recentFollowerNotify(fid)
	conn := d.redis.Get(c)
	defer conn.Close()
	flagi, err := redis.Int64(conn.Do("HGET", key, fid))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	flag := false
	if flagi > 0 {
		flag = true
	}
	return flag, err
}

// SetRctFollowerNotify is
func (d *Dao) SetRctFollowerNotify(c context.Context, fid int64, flag bool) error {
	key := recentFollowerNotify(fid)
	flagi := 0
	if flag {
		flagi = 1
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if err := conn.Send("HSET", key, fid, flagi); err != nil {
		return err
	}
	if err := conn.Send("EXPIRE", key, d.UnreadDuration); err != nil {
		return err
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	return nil
}
