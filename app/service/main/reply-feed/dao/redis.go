package dao

import (
	"context"
	"fmt"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	// r_<实验组名>_<oid>_<type>
	// 用redis ZSet存储热门评论列表，score为热评分数，member为rpID
	_replyZSetFormat = "r_%s_%d_%d"
)

func keyReplyZSet(name string, oid int64, tp int) string {
	return fmt.Sprintf(_replyZSetFormat, name, oid, tp)
}

// PingRedis redis health check.
func (d *Dao) PingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SET", "ping", "pong")
	return
}

// ExpireReplyZSetRds expire reply list.
func (d *Dao) ExpireReplyZSetRds(ctx context.Context, name string, oid int64, tp int) (ok bool, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplyZSet(name, oid, tp)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisReplyZSetExpire)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis EXPIRE key(%s) error(%v)", key, err)
	}
	return
}

// ReplyZSetRds get reply list from redis sorted set.
func (d *Dao) ReplyZSetRds(ctx context.Context, name string, oid int64, tp, start, end int) (rpIDs []int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplyZSet(name, oid, tp)
	values, err := redis.Values(conn.Do("ZREVRANGE", key, start, end))
	if err != nil {
		log.Error("redis ZREVRANGE(%s, %d, %d) error(%v)", key, start, end, err)
		return
	}
	if err = redis.ScanSlice(values, &rpIDs); err != nil {
		log.Error("redis ScanSlice(%v) error(%v)", values, err)
	}
	return
}

// CountReplyZSetRds count reply num.
func (d *Dao) CountReplyZSetRds(ctx context.Context, name string, oid int64, tp int) (count int, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplyZSet(name, oid, tp)
	if count, err = redis.Int(conn.Do("ZCARD", key)); err != nil {
		log.Error("conn.Do(ZCARD, %s) error(%v)", key, err)
	}
	return
}
