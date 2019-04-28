package dao

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_addPrefix2 = "nca_"
)

func hashField(aid, tp int64) int64 {
	return aid*1000 + tp
}

func addKey2(mid int64) (key string) {
	key = _addPrefix2 + strconv.FormatInt(mid, 10)
	return
}

// CoinsAddedCache get coin added of archive.
func (d *Dao) CoinsAddedCache(c context.Context, mid, aid, tp int64) (added int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := addKey2(mid)
	if added, err = redis.Int64(conn.Do("HGET", key, hashField(aid, tp))); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		PromError("redis:CoinsAddedCache")
		log.Errorv(c, log.KV("log", "redis.Do(HGET)"), log.KV("err", err), log.KV("mid", mid))
	}
	return
}

// SetCoinAddedCache set coin added of archive
func (d *Dao) SetCoinAddedCache(c context.Context, mid, aid, tp, count int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := addKey2(mid)
	defer func() {
		if err != nil {
			PromError("redis:SetCoinAddedCache")
			log.Errorv(c,
				log.KV("log", "s.coinDao.SetCoinAdded()"),
				log.KV("mid", mid),
				log.KV("aid", aid),
				log.KV("err", err),
			)
		}
	}()
	if err = conn.Send("HSETNX", key, hashField(aid, tp), count); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", key, d.expireAdded); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = redis.Bool(conn.Receive()); err != nil {
		return
	}
	conn.Receive()
	return
}

// SetCoinAddedsCache multiset added cache
func (d *Dao) SetCoinAddedsCache(c context.Context, mid int64, counts map[int64]int64) (err error) {
	if len(counts) == 0 {
		// 空缓存
		counts = map[int64]int64{-1: -1}
	}
	key := addKey2(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	for field, count := range counts {
		if err = conn.Send("HSETNX", key, field, count); err != nil {
			log.Errorv(c,
				log.KV("log", "conn.Send(HSETNX)"),
				log.KV("err", err),
				log.KV("mid", mid),
			)
			PromError("redis:SetCoinAddedsCache")
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.expireAdded); err != nil {
		log.Errorv(c,
			log.KV("log", "conn.Send(EXPIRE)"),
			log.KV("err", err),
			log.KV("mid", mid),
		)
		PromError("redis:SetCoinAddedsCache")
		return
	}
	if err = conn.Flush(); err != nil {
		log.Errorv(c, log.KV("log", "conn.Flush()"), log.KV("err", err), log.KV("mid", mid))
		PromError("redis:SetCoinAddedsCache")
		return
	}
	for i := 0; i < len(counts)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Errorv(c, log.KV("log", "conn.Recive()"), log.KV("err", err), log.KV("mid", mid))
			PromError("redis:SetCoinAddedsCache")
			return
		}
	}
	return
}

// IncrCoinAddedCache Incr coin added
func (d *Dao) IncrCoinAddedCache(c context.Context, mid, aid, tp, count int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := addKey2(mid)
	if _, err = conn.Do("HINCRBY", key, hashField(aid, tp), count); err != nil {
		PromError("redis:IncrCoinAdded")
		log.Errorv(c, log.KV("log", "conn.Do(HINCRBY) error"), log.KV("err", err), log.KV("mid", mid), log.KV("aid", aid))
	}
	return
}

// ExpireCoinAdded set expire time for coinadded
func (d *Dao) ExpireCoinAdded(c context.Context, mid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := addKey2(mid)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expireAdded)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		PromError("redis:ExpireCoinAdded")
		log.Errorv(c, log.KV("log", "conn.Do(EXPIRE)"), log.KV("err", err))
	}
	return
}
