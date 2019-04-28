package spam

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

// Cache Cache
type Cache struct {
	redisPool *redis.Pool
	expireRp  int
	expireAct int
}

// NewCache NewCache
func NewCache(c *redis.Config) *Cache {
	return &Cache{
		redisPool: redis.NewPool(c),
		expireRp:  60, // 60s
		expireAct: 20, // 20s
	}
}

func (c *Cache) keyRcntCnt(mid int64) string {
	return "rc_" + strconv.FormatInt(mid, 10)
}

func (c *Cache) keyUpRcntCnt(mid int64) string {
	return "urc_" + strconv.FormatInt(mid, 10)
}

func (c *Cache) keyDailyCnt(mid int64) string {
	return "rd_" + strconv.FormatInt(mid, 10)
}

func (c *Cache) keyActRec(mid int64) string {
	return "ra_" + strconv.FormatInt(mid, 10)
}

func (c *Cache) keySpamRpRec(mid int64) string {
	return "sr_" + strconv.FormatInt(mid, 10)
}

func (c *Cache) keySpamRpDaily(mid int64) string {
	return "sd_" + strconv.FormatInt(mid, 10)
}

func (c *Cache) keySpamActRec(mid int64) string {
	return "sa_" + strconv.FormatInt(mid, 10)
}

// IncrReply incr user reply count.
func (c *Cache) IncrReply(ctx context.Context, mid int64, isUp bool) (count int, err error) {
	key := c.keyRcntCnt(mid)
	if isUp {
		key = c.keyUpRcntCnt(mid)
	}
	conn := c.redisPool.Get(ctx)
	defer conn.Close()
	conn.Send("INCR", key)
	conn.Send("EXPIRE", key, c.expireRp)
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if count, err = redis.Int(conn.Receive()); err != nil {
		log.Error("conn.Receive error(%v)", key, err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", key, err)
		return
	}
	return
}

// IncrAct incr user action count.
func (c *Cache) IncrAct(ctx context.Context, mid int64) (count int, err error) {
	key := c.keyActRec(mid)
	conn := c.redisPool.Get(ctx)
	defer conn.Close()
	conn.Send("INCR", key)
	conn.Send("EXPIRE", key, c.expireAct)
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if count, err = redis.Int(conn.Receive()); err != nil {
		log.Error("conn.Receive error(%v)", key, err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", key, err)
		return
	}
	return
}

// IncrDailyReply IncrDailyReply
func (c *Cache) IncrDailyReply(ctx context.Context, mid int64) (count int, err error) {
	key := c.keyDailyCnt(mid)
	conn := c.redisPool.Get(ctx)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("INCR", key)); err != nil {
		log.Error("conn.Do(INCRBY, %s), error(%v)", key, err)
	}
	return
}

// TTLDailyReply TTLDailyReply
func (c *Cache) TTLDailyReply(ctx context.Context, mid int64) (ttl int, err error) {
	key := c.keyDailyCnt(mid)
	conn := c.redisPool.Get(ctx)
	defer conn.Close()
	if ttl, err = redis.Int(conn.Do("TTL", key)); err != nil {
		log.Error("conn.Do(TTL, %s), error(%v)", key, err)
	}
	return
}

// ExpireDailyReply ExpireDailyReply
func (c *Cache) ExpireDailyReply(ctx context.Context, mid int64, exp int) (err error) {
	key := c.keyDailyCnt(mid)
	conn := c.redisPool.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("EXPIRE", key, exp); err != nil {
		log.Error("conn.Do(EXPIRE, %s), error(%v)", key, err)
	}
	return
}

// SetReplyRecSpam SetReplyRecSpam
func (c *Cache) SetReplyRecSpam(ctx context.Context, mid int64, code, exp int) (err error) {
	key := c.keySpamRpRec(mid)
	conn := c.redisPool.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SETEX", key, exp, code); err != nil {
		log.Error("conn.Do error(%v)", err)
	}
	return
}

// SetReplyDailySpam SetReplyDailySpam
func (c *Cache) SetReplyDailySpam(ctx context.Context, mid int64, code, exp int) (err error) {
	key := c.keySpamRpDaily(mid)
	conn := c.redisPool.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SETEX", key, exp, code); err != nil {
		log.Error("conn.Do error(%v)", err)
	}
	return
}

// SetActionRecSpam SetActionRecSpam
func (c *Cache) SetActionRecSpam(ctx context.Context, mid int64, code, exp int) (err error) {
	key := c.keySpamActRec(mid)
	conn := c.redisPool.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SETEX", key, exp, code); err != nil {
		log.Error("conn.Do error(%v)", err)
	}
	return
}
