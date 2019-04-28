package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/member/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"strconv"

	errors "github.com/pkg/errors"
)

const (
	_expPrefix = "exp_%d"
	_expExpire = 86400
)

func expKey(mid int64) string {
	return fmt.Sprintf(_expPrefix, mid)
}

func (d *Dao) mcBaseKey(mid int64) (key string) {
	return fmt.Sprintf(model.CacheKeyBase, mid)
}

func (d *Dao) moralKey(mid int64) (key string) {
	return fmt.Sprintf(model.CacheKeyMoral, mid)
}

// SetStartCache set data import start
func (d *Dao) SetStartCache(c context.Context, mid int64, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if conn.Set(&memcache.Item{
		Key:        key,
		Object:     mid,
		Flags:      memcache.FlagJSON,
		Expiration: 0,
	}); err != nil {
		log.Error("conn.Set(%s, %d) error(%v)", key, mid, err)
	}
	return
}

// DelMoralCache delete moral cache.
func (d *Dao) DelMoralCache(c context.Context, mid int64) (err error) {
	key := d.moralKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%s) error(%v)", key, err)
	}
	return
}

// DelBaseInfoCache delete baseInfo cache.
func (d *Dao) DelBaseInfoCache(c context.Context, mid int64) (err error) {
	key := d.mcBaseKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%s) error(%v)", key, err)
	}
	return
}

// SetExpCache set user exp cache.
func (d *Dao) SetExpCache(c context.Context, mid, exp int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        expKey(mid),
		Value:      []byte(strconv.FormatInt(exp, 10)),
		Expiration: _expExpire,
	}); err != nil {
		log.Error("setexpcache mid %d err %v ", mid, err)
	}
	return
}

func realnameInfoKey(mid int64) string {
	return fmt.Sprintf("realname_info_%d", mid)
}

func realnameApplyStatusKey(mid int64) string {
	return fmt.Sprintf("realname_apply_%d", mid)
}

// DeleteRealnameCache delete all realname cache
func (d *Dao) DeleteRealnameCache(c context.Context, mid int64) (err error) {
	var (
		key1 = realnameInfoKey(mid)
		key2 = realnameApplyStatusKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key1); err != nil {
		if err != memcache.ErrNotFound {
			err = errors.Wrapf(err, "conn.Delete(%s)", key1)
			return
		}
		err = nil
	}
	if err = conn.Delete(key2); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			err = errors.Wrapf(err, "conn.Delete(%s)", key2)
		}
		return
	}
	return
}
