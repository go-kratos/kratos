package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/ugcpay/conf"
	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

func assetRelationKey(mid int64) string {
	return fmt.Sprintf("up_ar_%d", mid)
}

func assetRelationField(oid int64, otype string) string {
	return fmt.Sprintf("%s_%d", otype, oid)
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// CacheAssetRelationState get asset relation state.
func (d *Dao) CacheAssetRelationState(c context.Context, oid int64, otype string, mid int64) (state string, err error) {
	var (
		key   = assetRelationKey(mid)
		field = assetRelationField(oid, otype)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if state, err = redis.String(conn.Do("HGET", key, field)); err != nil {
		if err == redis.ErrNil {
			err = nil
			state = "miss"
			return
		}
		err = errors.Wrapf(err, "conn.Do(HGET, %s ,%s)", key, field)
		return
	}
	return
}

// AddCacheAssetRelationState set asset relation state.
func (d *Dao) AddCacheAssetRelationState(c context.Context, oid int64, otype string, mid int64, state string) (err error) {
	var (
		key   = assetRelationKey(mid)
		field = assetRelationField(oid, otype)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("HSET", key, field, state); err != nil {
		err = errors.Wrapf(err, "conn.Do(HSET, %s, %s, %s)", key, field, state)
		return
	}
	if _, err = conn.Do("EXPIRE", key, conf.Conf.CacheTTL.AssetRelationTTL); err != nil {
		err = errors.Wrapf(err, "conn.Do(EXPIRE, %s, %d)", key, conf.Conf.CacheTTL.AssetRelationTTL)
		return
	}
	return
}

// DelCacheAssetRelationState delete asset relation state.
func (d *Dao) DelCacheAssetRelationState(c context.Context, oid int64, otype string, mid int64) (err error) {
	var (
		key   = assetRelationKey(mid)
		field = assetRelationField(oid, otype)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("HDEL", key, field); err != nil {
		err = errors.Wrapf(err, "conn.Do(HDEL, %s, %s)", key, field)
		return
	}
	return
}

// DelCacheAssetRelation delete assetrelation.
func (d *Dao) DelCacheAssetRelation(c context.Context, mid int64) (err error) {
	var (
		key  = assetRelationKey(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		err = errors.Wrapf(err, "conn.Do(DEL, %s)", key)
		return
	}
	return
}
