package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/share/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	xip "go-common/library/net/ip"

	farm "github.com/dgryski/go-farm"
	"github.com/pkg/errors"
)

func redisKey(oid int64, tp int) string {
	return fmt.Sprintf("%d_%d", oid, tp)
}

func redisValue(p *model.ShareParams) int64 {
	return int64(farm.Hash64([]byte(fmt.Sprintf("%d_%d_%d_%s", p.MID, p.OID, p.TP, p.IP))))
}

func shareKey(oid int64, tp int) string {
	return fmt.Sprintf("c_%d_%d", oid, tp)
}

// AddShareMember add share
func (d *Dao) AddShareMember(ctx context.Context, p *model.ShareParams) (ok bool, err error) {
	var (
		conn  = d.rds.Get(ctx)
		key   = redisKey(p.OID, p.TP)
		value = (p.MID << 32) | int64(xip.InetAtoN(p.IP))
	)
	log.Info("oid-%d mid-%d ip-%s tp-%d key-%s value-%d", p.OID, p.MID, p.IP, p.TP, key, value)
	defer conn.Close()
	if err = conn.Send("SADD", key, value); err != nil {
		err = errors.Wrapf(err, "conn.Do(SADD, %s, %d)", key, value)
		return
	}
	if err = conn.Send("EXPIRE", key, d.c.RedisExpire); err != nil {
		err = errors.Wrapf(err, "conn.Do(SADD, %s, %d)", key, value)
		return
	}
	if err = conn.Flush(); err != nil {
		err = errors.Wrap(err, "conn.Flush")
		return
	}
	if ok, err = redis.Bool(conn.Receive()); err != nil {
		log.Error("sadd failed mid(%d) oid(%d) type(%d) ip(%s) key(%s) value(%d)",
			p.MID, p.OID, p.TP, p.IP, key, value)
		err = errors.Wrap(err, "redis.Bool(conn.Receive)")
		return
	}
	if _, err = conn.Receive(); err != nil {
		err = errors.Wrap(err, "conn.Receive")
		return
	}
	return
}

// SetShareCache set share cache
func (d *Dao) SetShareCache(c context.Context, oid int64, tp int, shared int64) (err error) {
	var (
		conn = d.rds.Get(c)
		key  = shareKey(oid, tp)
	)
	defer conn.Close()
	if _, err = conn.Do("SET", key, shared); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// ShareCache return oid share count
func (d *Dao) ShareCache(c context.Context, oid int64, tp int) (shared int64, err error) {
	var (
		conn = d.rds.Get(c)
		key  = shareKey(oid, tp)
	)
	defer conn.Close()
	if shared, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			shared = -1
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// SharesCache return oids share
func (d *Dao) SharesCache(c context.Context, oids []int64, tp int) (shares map[int64]int64, err error) {
	conn := d.rds.Get(c)
	defer conn.Close()
	for _, oid := range oids {
		if err = conn.Send("GET", shareKey(oid, tp)); err != nil {
			log.Error("conn.Send(GET, %s) error(%v)", shareKey(oid, tp), err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	shares = make(map[int64]int64, len(oids))
	for _, oid := range oids {
		var cnt int64
		if cnt, err = redis.Int64(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
				continue
			}
			return
		}
		shares[oid] = cnt
	}
	return
}
