package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/interface/main/kvo/model"

	mc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_userConfKeyPrefix = "u"
)

func cacheUserConfKey(mid int64, moduleKey int) string {
	return fmt.Sprintf("%v_%v_%v", _userConfKeyPrefix, mid, moduleKey)
}

// UserConfCache user config cache
func (d *Dao) UserConfCache(ctx context.Context, mid int64, moduleKey int) (uc *model.UserConf, err error) {
	var r *mc.Item
	conn := d.cache.Get(ctx)
	defer conn.Close()
	if r, err = conn.Get(cacheUserConfKey(mid, moduleKey)); err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", cacheUserConfKey(mid, moduleKey), err)
		return
	}
	if err = json.Unmarshal(r.Value, &uc); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		uc = nil
	}
	return
}

// SetUserConfCache set user config cache
func (d *Dao) SetUserConfCache(ctx context.Context, uc *model.UserConf) (err error) {
	bs, err := json.Marshal(uc)
	if err != nil {
		log.Error("SetUserConfCache.Marshal err:%v", err)
		return
	}
	conn := d.cache.Get(ctx)
	defer conn.Close()
	if err = conn.Set(&mc.Item{
		Key:        cacheUserConfKey(uc.Mid, uc.ModuleKey),
		Value:      bs,
		Expiration: d.mcExpire,
	}); err != nil {
		log.Error("dao.SetUserConfCache(%v,%v) err:%v", cacheUserConfKey(uc.Mid, uc.ModuleKey), bs, err)
		return
	}
	return
}

// DelUserConfCache del user config cache
func (d *Dao) DelUserConfCache(ctx context.Context, mid int64, moduleKey int) (err error) {
	conn := d.cache.Get(ctx)
	defer conn.Close()
	if err = conn.Delete(cacheUserConfKey(mid, moduleKey)); err == mc.ErrNotFound {
		err = nil
	}
	return
}
