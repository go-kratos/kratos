package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
)

const (
	_cacheShard = 10000
	_upPrompt   = "rl_up_%d_%d"    // key of upper prompt; hashes(fid-count)
	_buPrompt   = "rl_bu_%d_%d_%d" // key of business type prompt;hashes(mid-count)
)

// key upPrompt : rl_up_mid_ts/period
func (d *Dao) upPrompt(mid, ts int64) string {
	return fmt.Sprintf(_upPrompt, mid, ts/d.period)
}

// key _buPrompt : rl_bu_businesstype_mid/10000_ts
func (d *Dao) buPrompt(btype int8, mid, ts int64) string {
	return fmt.Sprintf(_buPrompt, btype, mid/_cacheShard, ts/d.period)
}

// IncrPromptCount incr up prompt count and business type prompt count.
func (d *Dao) IncrPromptCount(c context.Context, mid, fid, ts int64, btype int8) (ucount, bcount int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	keyUp := d.upPrompt(mid, ts)
	keyBs := d.buPrompt(btype, mid, ts)
	conn.Send("HINCRBY", keyUp, fid, 1)
	conn.Send("EXPIRE", keyUp, d.period)
	conn.Send("HINCRBY", keyBs, mid, 1)
	conn.Send("EXPIRE", keyBs, d.period)
	err = conn.Flush()
	if err != nil {
		return
	}
	ucount, err = redis.Int64(conn.Receive())
	if err != nil {
		return
	}
	conn.Receive()
	bcount, err = redis.Int64(conn.Receive())
	if err != nil {
		return
	}
	conn.Receive()
	return
}

// ClosePrompt set prompt count to max config value.
func (d *Dao) ClosePrompt(c context.Context, mid, fid, ts int64, btype int8) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	keyUp := d.upPrompt(mid, ts)
	keyBs := d.buPrompt(btype, mid, ts)
	conn.Send("HSET", keyUp, fid, d.ucount)
	conn.Send("HSET", keyBs, mid, d.bcount)
	return conn.Flush()
}

// UpCount get upper prompt count.
func (d *Dao) UpCount(c context.Context, mid, fid, ts int64) (count int64, err error) {
	conn := d.redis.Get(c)
	count, err = redis.Int64(conn.Do("HGET", d.upPrompt(mid, ts), fid))
	conn.Close()
	return
}

// BCount get business type prompt count.
func (d *Dao) BCount(c context.Context, mid, ts int64, btype int8) (count int64, err error) {
	conn := d.redis.Get(c)
	count, err = redis.Int64(conn.Do("HGET", d.buPrompt(btype, mid, ts), mid))
	conn.Close()
	return
}
