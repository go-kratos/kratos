package dao

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"go-common/app/service/main/filter/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

func mcFilterKey(area string, tpid int64, keys []string, content string) string {
	return fmt.Sprintf("f_%s_%d_%s_%s", area, tpid, base64.StdEncoding.EncodeToString([]byte(strings.Join(keys, "|"))), base64.StdEncoding.EncodeToString([]byte(content)))
}

// FilterCache .
func (d *Dao) FilterCache(c context.Context, area string, tpid int64, keys []string, content string) (res *model.FilterCacheRes, err error) {
	var (
		key   = mcFilterKey(area, tpid, keys, content)
		conn  = d.mc.Get(c)
		reply *memcache.Item
	)
	defer conn.Close()

	if reply, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(get,%s)", key)
		return
	}
	res = &model.FilterCacheRes{}
	if err = conn.Scan(reply, &res); err != nil {
		err = errors.Wrapf(err, "reply.Scan(%s)", string(reply.Value))
	}
	log.Info("Filter hit cache key (%s) level (%d)", key, res.Level)
	return
}

// SetFilterCache .
func (d *Dao) SetFilterCache(c context.Context, area string, tpid int64, keys []string, content string, res *model.FilterCacheRes) (err error) {
	var (
		key  = mcFilterKey(area, tpid, keys, content)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcFilterExp}); err != nil {
		err = errors.WithStack(err)
		return
	}
	log.Info("Filter set cache key (%s) level (%d)", key, res.Level)
	return
}

func mcKey(key, area string) string {
	return fmt.Sprintf("%s_%s", key, area)
}

// KeyAreaCache .
func (d *Dao) KeyAreaCache(c context.Context, key string, areas []string) (rs []*model.KeyAreaInfo, miss []string, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	var (
		mcKeys  []string
		rpMap   map[string]*memcache.Item
		areaMap = make(map[string]struct{})
	)
	for _, area := range areas {
		tempKey := mcKey(key, area)
		mcKeys = append(mcKeys, tempKey)
		areaMap[area] = struct{}{}
	}
	if rpMap, err = conn.GetMulti(mcKeys); err != nil {
		if err == memcache.ErrNotFound {
			miss = areas
			err = nil
		}
		return
	}
	for _, r := range rpMap {
		val := []*model.KeyAreaInfo{}
		if err = conn.Scan(r, &val); err != nil {
			log.Error("r.Scan() error(%v)", err)
			return
		}
		if len(val) > 0 {
			rs = append(rs, val...)
			mm := strings.Split(r.Key, "_")
			delete(areaMap, mm[1])
		}
	}
	for area := range areaMap {
		miss = append(miss, area)
	}
	return
}

// SetKeyAreaCache .
func (d *Dao) SetKeyAreaCache(c context.Context, key, area string, rs []*model.KeyAreaInfo) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	mcKey := mcKey(key, area)
	if err = conn.Set(&memcache.Item{Key: mcKey, Object: rs, Flags: memcache.FlagJSON, Expiration: d.mcKeyFilterExp}); err != nil {
		log.Error("conn.Set(%v) error(%v)", rs, err)
	}
	return
}

// DelKeyAreaCache .
func (d *Dao) DelKeyAreaCache(c context.Context, key, area string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	mcKey := mcKey(key, area)
	if err = conn.Delete(mcKey); err != nil {
		log.Error("conn.Delete(%s) error(%v)", mcKey, err)
		return
	}
	return
}
