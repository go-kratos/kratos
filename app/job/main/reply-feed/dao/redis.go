package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/reply-feed/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	// r_<实验组名>_<oid>_<type>
	// 用redis ZSet存储热门评论列表，score为热评分数，member为rpID
	_replyZSetFormat = "r_%s_%d_%d"

	// c_<oid>_<type>
	_refreshCheckerFormat = "c_%d_%d"

	// h_<oid>_<type>
	// 用一个set来存某一个评论区下的应该存在于热门评论列表的评论ID
	// 上热评的最低门槛，点赞数大于3，且该评论区根评论数目大于20
	_replyListFormat = "h_%d_%d"

	// 行为, 小时，slot, 种类(全量或者针对热评)
	_uvFormat = "uv_%s_%d_%d_%s"

	_uvExp = 3600
)

func keyRefreshChecker(oid int64, tp int) string {
	return fmt.Sprintf(_refreshCheckerFormat, oid, tp)
}

func keyReplyZSet(name string, oid int64, tp int) string {
	return fmt.Sprintf(_replyZSetFormat, name, oid, tp)
}

func keyReplySet(oid int64, tp int) string {
	return fmt.Sprintf(_replyListFormat, oid, tp)
}

// KeyUV ...
func (d *Dao) KeyUV(action string, hour, slot int, kind string) string {
	return keyUV(action, hour, slot, kind)
}

func keyUV(action string, hour, slot int, kind string) string {
	return fmt.Sprintf(_uvFormat, action, hour, slot, kind)
}

// PingRedis redis health check.
func (d *Dao) PingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SET", "ping", "pong")
	return
}

// AddUV ...
func (d *Dao) AddUV(ctx context.Context, action string, hour, slot int, mid int64, kind string) (err error) {
	key := keyUV(action, hour, slot, kind)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if err = conn.Send("SADD", key, mid); err != nil {
		log.Error("redis SADD(%s, %d) error(%v)", key, mid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, _uvExp); err != nil {
		log.Error("redis EXPIRE(%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("redis Flush() error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("redis Receive() error(%v)", err)
			return
		}
	}
	return
}

// CountUV ...
func (d *Dao) CountUV(ctx context.Context, keys []string) (counts []int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	j := 0
	for _, key := range keys {
		if err = conn.Send("SCARD", key); err != nil {
			log.Error("redis SCARD(%s) error(%v)", key, err)
			return
		}
		j++
	}
	if err = conn.Flush(); err != nil {
		log.Error("redis Flush() error(%v)", err)
		return
	}
	for i := 0; i < j; i++ {
		var count int64
		if count, err = redis.Int64(conn.Receive()); err != nil && err != redis.ErrNil {
			log.Error("redis Receive() error(%v)", err)
			return
		}
		counts = append(counts, count)
	}
	return
}

// ExpireCheckerRds expire checker.
func (d *Dao) ExpireCheckerRds(ctx context.Context, oid int64, tp int) (ok bool, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyRefreshChecker(oid, tp)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisRefreshExpire)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis EXPIRE key(%s) error(%v)", key, err)
	}
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

// ExpireReplySetRds expire reply set.
func (d *Dao) ExpireReplySetRds(ctx context.Context, oid int64, tp int) (ok bool, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplySet(oid, tp)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisReplySetExpire)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis EXPIRE key(%s) error(%v)", key, err)
	}
	return
}

// ReplySetRds get reply rpIDs from redis set.
func (d *Dao) ReplySetRds(ctx context.Context, oid int64, tp int) (rpIDs []int64, err error) {
	key := keyReplySet(oid, tp)
	var startCursor, endCursor int64
	endCursor = 1
	for endCursor != 0 {
		var (
			conn         = d.redis.Get(ctx)
			chunkedRpIDs []int64
			values       []interface{}
		)
		if values, err = redis.Values(conn.Do("SSCAN", key, startCursor)); err != nil {
			log.Error("redis SSCAN(%s) error(%v)", key, err)
			conn.Close()
			return
		}
		if _, err = redis.Scan(values, &endCursor, &chunkedRpIDs); err != nil {
			log.Error("redis Scan(%v) error(%v)", values, err)
			conn.Close()
			return
		}
		startCursor = endCursor
		rpIDs = append(rpIDs, chunkedRpIDs...)
		conn.Close()
	}
	return
}

// SetReplySetRds set reply list batch, call it when back to source.
func (d *Dao) SetReplySetRds(ctx context.Context, oid int64, tp int, rpIDs []int64) (err error) {
	if len(rpIDs) < 1 {
		return
	}
	for _, chunkedRpIDs := range split(rpIDs, 5000) {
		var (
			key  = keyReplySet(oid, tp)
			args = make([]interface{}, 0, len(chunkedRpIDs)+1)
			conn = d.redis.Get(ctx)
		)
		args = append(args, key)
		for _, rpID := range chunkedRpIDs {
			args = append(args, rpID)
		}
		if err = conn.Send("SADD", args...); err != nil {
			log.Error("redis SADD(%v) error(%v)", args, err)
			conn.Close()
			return
		}
		if err = conn.Send("EXPIRE", key, d.redisReplySetExpire); err != nil {
			log.Error("redis EXPIRE(%s) error(%v)", key, err)
			conn.Close()
			return
		}
		if err = conn.Flush(); err != nil {
			log.Error("redis Flush() error(%v)", err)
			conn.Close()
			return
		}
		for i := 0; i < 2; i++ {
			if _, err = conn.Receive(); err != nil {
				log.Error("redis Receive() error(%v)", err)
				conn.Close()
				return
			}
		}
		conn.Close()
	}
	return
}

// RemReplySetRds remove one rp from set.
func (d *Dao) RemReplySetRds(ctx context.Context, oid, rpID int64, tp int) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplySet(oid, tp)
	if _, err = redis.Int64(conn.Do("SREM", key, rpID)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("SREM conn.Do(%s,%d) err(%v)", key, rpID, err)
	}
	return
}

// DelReplySetRds delete a set key.
func (d *Dao) DelReplySetRds(ctx context.Context, oid int64, tp int) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplySet(oid, tp)
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
}

// AddReplySetRds add a reply into redis set, make sure expire the key first.
func (d *Dao) AddReplySetRds(ctx context.Context, oid int64, tp int, rpID int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplySet(oid, tp)
	if _, err = conn.Do("SADD", key, rpID); err != nil {
		log.Error("redis SADD(%s, %d) error(%v)", key, rpID, err)
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

// SetReplyZSetRds set reply list batch, call it when back to source.
func (d *Dao) SetReplyZSetRds(ctx context.Context, name string, oid int64, tp int, rs []*model.ReplyScore) (err error) {
	if len(rs) < 1 {
		return
	}
	for _, chunkedReplyStats := range splitReplyScore(rs, 5000) {
		var (
			count = 0
			key   = keyReplyZSet(name, oid, tp)
			conn  = d.redis.Get(ctx)
			args  = make([]interface{}, 0, len(chunkedReplyStats)*2+1)
		)
		args = append(args, key)
		for _, s := range chunkedReplyStats {
			args = append(args, s.Score)
			args = append(args, s.RpID)
		}
		if err = conn.Send("ZADD", args...); err != nil {
			log.Error("redis ZADD(%s, %v) error(%v)", key, args, err)
			conn.Close()
			return
		}
		count++
		if err = conn.Send("EXPIRE", key, d.redisReplyZSetExpire); err != nil {
			log.Error("redis EXPIRE(%s) error(%v)", key, err)
			conn.Close()
			return
		}
		count++
		if err = conn.Flush(); err != nil {
			log.Error("redis Flush error(%v)", err)
			conn.Close()
			return
		}
		for i := 0; i < count; i++ {
			if _, err = conn.Receive(); err != nil {
				log.Error("redis Receive (key: %s, %f, %d) error(%v)", key, rs[i].Score, rs[i].RpID, err)
				conn.Close()
				return
			}
		}
		conn.Close()
	}
	return
}

// RangeReplyZSetRds ...
func (d *Dao) RangeReplyZSetRds(ctx context.Context, name string, oid int64, tp, start, end int) (rpIDs []int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplyZSet(name, oid, tp)
	values, err := redis.Values(conn.Do("ZREVRANGE", key, start, end))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	if err = redis.ScanSlice(values, &rpIDs); err != nil {
		log.Error("redis.ScanSlice(%v) error(%v)", values, err)
		return
	}
	return
}

// AddReplyZSetRds add a reply into redis sorted set, make sure expire the key first.
func (d *Dao) AddReplyZSetRds(ctx context.Context, name string, oid int64, tp int, rs *model.ReplyScore) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplyZSet(name, oid, tp)
	if _, err = conn.Do("ZADD", key, rs.Score, rs.RpID); err != nil {
		log.Error("redis ZADD(%s, %f, %d) error(%v)", key, rs.Score, rs.RpID, err)
	}
	return
}

// RemReplyZSetRds remove one rpID from reply ZSet.
func (d *Dao) RemReplyZSetRds(ctx context.Context, name string, oid int64, tp int, rpID int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyReplyZSet(name, oid, tp)
	if _, err = conn.Do("ZREM", key, rpID); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis ZREM(%s, %d) error(%v)", key, rpID, err)
	}
	return
}

// DelReplyZSetRds del a key from reply ZSet.
func (d *Dao) DelReplyZSetRds(ctx context.Context, names []string, oid int64, tp int) (err error) {
	if len(names) < 1 {
		return
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	count := 0
	for _, name := range names {
		key := keyReplyZSet(name, oid, tp)
		if err = conn.Send("DEL", key); err != nil {
			log.Error("conn.Do(DEL %s) error(%v)", key, err)
			return
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		log.Error("redis Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("redis Receive error(%v)", err)
			return
		}
	}
	return
}

// CheckerTsRds get refresh checker timestamp from redis, if time.Now()-ts > strategy.time, then refresh reply list.
func (d *Dao) CheckerTsRds(ctx context.Context, oid int64, tp int) (ts int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyRefreshChecker(oid, tp)
	if ts, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis GET(%s) error(%v)", key, err)
	}
	return
}

// SetCheckerTsRds set refresh checker's timestamp as time.Now(), call it when refresh reply list.
func (d *Dao) SetCheckerTsRds(ctx context.Context, oid int64, tp int) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := keyRefreshChecker(oid, tp)
	if err = conn.Send("SETEX", key, d.redisRefreshExpire, time.Now().Unix()); err != nil {
		log.Error("redis SETEX(%s) error(%v)", key, err)
	}
	return
}
