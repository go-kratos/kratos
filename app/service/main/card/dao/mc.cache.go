package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/service/main/card/model"
	mc "go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"

	"github.com/pkg/errors"
)

const (
	_prequip = "e_"
)

func equipKey(mid int64) string {
	return _prequip + strconv.FormatInt(mid, 10)
}

// pingMC ping memcache.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	err = conn.Set(&mc.Item{
		Key:   "ping",
		Value: []byte("pong"),
	})
	return
}

// CacheEquips get data from mc
func (d *Dao) CacheEquips(c context.Context, mids []int64) (res map[int64]*model.UserEquip, err error) {
	keys := make([]string, 0, len(mids))
	keyMidMap := make(map[string]int64, len(mids))
	for _, mid := range mids {
		key := equipKey(mid)
		if _, ok := keyMidMap[key]; !ok {
			// duplicate mid
			keyMidMap[key] = mid
			keys = append(keys, key)
		}
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	rs, err := conn.GetMulti(keys)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao equips cache")
		return
	}
	res = make(map[int64]*model.UserEquip, len(mids))
	for k, r := range rs {
		e := &model.UserEquip{}
		conn.Scan(r, e)
		res[keyMidMap[k]] = e
	}
	return
}

// CacheEquip get user card equip from cache.
func (d *Dao) CacheEquip(c context.Context, mid int64) (v *model.UserEquip, err error) {
	key := equipKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao cache equip")
		return
	}
	v = &model.UserEquip{}
	if err = conn.Scan(r, v); err != nil {
		err = errors.Wrap(err, "dao cache scan equip")
	}
	return
}

// AddCacheEquips Set data to mc
func (d *Dao) AddCacheEquips(c context.Context, values map[int64]*model.UserEquip) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for k, v := range values {
		item := &mc.Item{
			Key:        equipKey(k),
			Object:     v,
			Flags:      mc.FlagProtobuf,
			Expiration: d.mcExpire,
		}
		if err = conn.Set(item); err != nil {
			err = errors.Wrap(err, "dao add equips cache")
		}
	}
	return
}

// AddCacheEquip set user card equip info into cache.
func (d *Dao) AddCacheEquip(c context.Context, mid int64, v *model.UserEquip) (err error) {
	item := &mc.Item{
		Key:        equipKey(mid),
		Object:     v,
		Flags:      mc.FlagProtobuf,
		Expiration: d.mcExpire,
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		err = errors.Wrap(err, "dao add equip cache")
	}
	return
}

// DelCacheEquips delete data from mc
func (d *Dao) DelCacheEquips(c context.Context, ids []int64) (err error) {
	if len(ids) == 0 {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, id := range ids {
		key := equipKey(id)
		if err = conn.Delete(key); err != nil {
			if err == mc.ErrNotFound {
				err = nil
				continue
			}
			prom.BusinessErrCount.Incr("mc:DelCacheEquips")
			log.Errorv(c, log.KV("DelCacheEquips", fmt.Sprintf("%+v", err)), log.KV("key", key))
			return
		}
	}
	return
}

// DelCacheEquip delete data from mc
func (d *Dao) DelCacheEquip(c context.Context, id int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := equipKey(id)
	if err = conn.Delete(key); err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:DelCacheEquip")
		log.Errorv(c, log.KV("DelCacheEquip", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
