package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/service/main/figure/model"
	"go-common/library/cache/redis"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_figureKey = "f:%d"
)

func figureKey(mid int64) string {
	return fmt.Sprintf(_figureKey, mid)
}

// PingRedis check redis connection
func (d *Dao) PingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AddFigureInfoCache put figure to redis
func (d *Dao) AddFigureInfoCache(c context.Context, f *model.Figure) (err error) {
	var (
		key    = figureKey(f.Mid)
		conn   = d.redis.Get(c)
		values []byte
	)
	defer conn.Close()
	if values, err = json.Marshal(f); err != nil {
		err = errors.Wrapf(err, "%+v", f)
		return
	}
	if err = conn.Send("SET", key, values); err != nil {
		err = errors.Wrapf(err, "conn.Send(SET, %s, %d)", key, values)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		err = errors.Wrapf(err, "conn.Send(Expire, %s, %d)", key, d.redisExpire)
		return
	}
	return
}

// FigureInfoCache get user figure info from redis
func (d *Dao) FigureInfoCache(c context.Context, mid int64) (f *model.Figure, err error) {
	key := figureKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	item, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	if err = json.Unmarshal(item, &f); err != nil {
		log.Error("json.Unmarshal(%v) err(%v)", item, err)
	}
	return
}

// FigureBatchInfoCache ...
func (d *Dao) FigureBatchInfoCache(c context.Context, mids []int64) (fs []*model.Figure, missIndex []int, err error) {
	if len(mids) == 0 {
		return
	}
	fs = make([]*model.Figure, len(mids))
	var (
		conn       = d.redis.Get(c)
		valueBytes [][]byte
		keys       []interface{}
	)
	defer conn.Close()
	for _, mid := range mids {
		keys = append(keys, figureKey(mid))
	}
	if valueBytes, err = redis.ByteSlices(conn.Do("MGET", keys...)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	for i, value := range valueBytes {
		if value == nil {
			missIndex = append(missIndex, i)
			continue
		}
		f := &model.Figure{}
		if err = json.Unmarshal(value, &f); err != nil {
			log.Error("%+v", errors.Wrapf(err, "json.Unmarshal(%s)", value))
			err = nil
			missIndex = append(missIndex, i)
			continue
		}
		fs[i] = f
	}
	return
}
