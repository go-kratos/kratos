package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/credit/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_blockIdx = "bl_%d_%d"
)

func blockIndexKey(otype, btype int8) string {
	return fmt.Sprintf(_blockIdx, otype, btype)
}

// BlockedIdxCache get block list idx.
func (d *Dao) BlockedIdxCache(c context.Context, otype, btype int8, start, end int) (ids []int64, err error) {
	key := blockIndexKey(otype, btype)
	conn := d.redis.Get(c)
	if ids, err = redis.Int64s(conn.Do("ZREVRANGE", key, start, end)); err != nil {
		log.Info("Redis.ZREVRANGE err(%v)", err)
	}
	conn.Close()
	return
}

// ExpireBlockedIdx expire case index cache.
func (d *Dao) ExpireBlockedIdx(c context.Context, otype, btype int8) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", blockIndexKey(otype, btype), d.redisExpire)); err != nil {
		log.Error("redis.bool err(%v)", err)
	}
	return
}

// LoadBlockedIdx laod blocked info index.
func (d *Dao) LoadBlockedIdx(c context.Context, otype, btype int8, infos []*model.BlockedInfo) (err error) {
	key := blockIndexKey(otype, btype)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, info := range infos {
		if err = conn.Send("ZADD", key, info.PublishTime, info.ID); err != nil {
			log.Error("ZADD err(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("EXPIRE err(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < len(infos)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return
}
