package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/job/main/figure/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_keyWaitDealUser = "w:u" // b_batch_no wait block
	_figureKey       = "f:%d"
)

func figureKey(mid int64) string {
	return fmt.Sprintf(_figureKey, mid)
}

// keyWaitBlock return block cache key.
func keyWaitBlock(version int64, mid int64) string {
	return fmt.Sprintf("%s%d%d", _keyWaitDealUser, version, mid%10000)
}

// SetWaiteUserCache set waite deal user cache.
func (d *Dao) SetWaiteUserCache(c context.Context, mid int64, ver int64) (err error) {
	if mid <= 0 {
		log.Error("%+v", errors.Errorf("SetWaiteUserCache mid [%d] ver [%d] error", mid, ver))
		return
	}
	var (
		key  = keyWaitBlock(ver, mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("SADD", key, mid); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.waiteMidExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// AddFigureInfoCache put figure to redis
func (d *Dao) AddFigureInfoCache(c context.Context, f *model.Figure) (err error) {
	key := figureKey(f.Mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := json.Marshal(f)
	if err != nil {
		return
	}
	if err = conn.Send("SET", key, values); err != nil {
		log.Error("conn.Send(SET, %s, %d) error(%v)", key, values, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.redisExpire, err)
		return
	}
	return
}

// PingRedis check redis connection
func (d *Dao) PingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
