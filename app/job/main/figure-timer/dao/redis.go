package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/job/main/figure-timer/model"
	"go-common/library/cache/redis"
	"go-common/library/log"

	"github.com/pkg/errors"
)

func keyFigure(mid int64) string {
	return fmt.Sprintf("f:%d", mid)
}

func keyPendingMids(ver int64, shard int64) string {
	return fmt.Sprintf("w:u%d%d", ver, shard)
}

// PingRedis check redis connection
func (d *Dao) PingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("SET", "PING", "PONG")
	return
}

// FigureCache get FigureUser from cache
func (d *Dao) FigureCache(c context.Context, mid int64) (figure *model.Figure, err error) {
	key := keyFigure(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	item, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	figure = &model.Figure{}
	if err = json.Unmarshal(item, &figure); err != nil {
		log.Error("json.Unmarshal(%v) err(%v)", item, err)
	}
	return
}

// SetFigureCache set FigureUser to cache
func (d *Dao) SetFigureCache(c context.Context, figure *model.Figure) (err error) {
	key := keyFigure(figure.Mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := json.Marshal(figure)
	if err != nil {
		return
	}
	if err = conn.Send("SET", key, values); err != nil {
		log.Error("conn.Send(SET, %s, %v) error(%v)", key, values, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.redisExpire, err)
		return
	}
	return
}

// PendingMidsCache get PendingUser set from cache
func (d *Dao) PendingMidsCache(c context.Context, version int64, shard int64) (mids []int64, err error) {
	var (
		conn = d.redis.Get(c)
		key  = keyPendingMids(version, shard)
	)
	defer conn.Close()
	if mids, err = redis.Int64s(conn.Do("SMEMBERS", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		err = errors.Wrapf(err, "redis.Int64s(conn.Do(SMEMEBERS,%s))", key)
		return
	}
	return
}

// RemoveCache remove figure cache
func (d *Dao) RemoveCache(c context.Context, mid int64) (err error) {
	key := keyFigure(mid)
	conn := d.redis.Get(c)
	defer conn.Close()

	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	return
}
