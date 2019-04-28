package dao

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

const (
	_prefixBlackList = "zl_"
)

func keyZlimit(aid int64) (key string) {
	key = _prefixBlackList + strconv.FormatInt(aid, 10)
	return
}

// ExistsAuth if existes ruls in redis.
func (d *Dao) ExistsAuth(c context.Context, aid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXISTS", keyZlimit(aid))); err != nil {
		err = errors.Wrapf(err, "EXISTS %s", keyZlimit(aid))
	}
	return
}

// Auth get zone rule from redis
func (d *Dao) Auth(c context.Context, aid int64, zoneids []int64) (res []int64, err error) {
	var playauth int64
	key := keyZlimit(aid)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, v := range zoneids {
		if err = conn.Send("HGET", key, v); err != nil {
			err = errors.Wrapf(err, "HGET %s %d", key, v)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		err = errors.WithStack(err)
		return
	}
	for range zoneids {
		if playauth, err = redis.Int64(conn.Receive()); err != nil {
			if err != redis.ErrNil {
				err = errors.WithStack(err)
				return
			}
			err = nil
		}
		res = append(res, playauth)
	}
	return
}

// AddAuth add zone rule from redis
func (d *Dao) AddAuth(c context.Context, zoneids map[int64]map[int64]int64) (err error) {
	var key string
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	for aid, zids := range zoneids {
		if key == "" {
			key = keyZlimit(aid)
		}
		for zid, auth := range zids {
			if err = conn.Send("HSET", key, zid, auth); err != nil {
				err = errors.Wrapf(err, "HGET %s %d", key, zid)
				return
			}
			count++
		}
	}
	if err = conn.Send("EXPIRE", key, d.expire); err != nil {
		err = errors.Wrapf(err, "EXPIRE %s %d", key, d.expire)
		return
	}
	if err = conn.Flush(); err != nil {
		err = errors.WithStack(err)
		return
	}
	for i := 0; i <= count; i++ {
		if _, err = conn.Receive(); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}
