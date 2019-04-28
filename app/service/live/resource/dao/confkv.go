package dao

import (
	"context"
	"encoding/json"

	v1pb "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_findConfkvSQL   = "SELECT `value` FROM `confkv` WHERE `key` = ?"
	_addConfkvSQL    = "INSERT INTO `confkv` (`key`,`value`) VALUES (?,?)"
	_updateConfkvSQL = "UPDATE `confkv` SET `value`=? WHERE `key`=?"
)

const (
	_confLiveCheck = "live_check"

	_platformAndroid = "android"
	_platformIos     = "ios"
)

// GetLiveCheck live.app-interface call
// cache -> db
func (d *Dao) GetLiveCheck(c context.Context, platform, system, mobile string) (isLive int64) {
	isLive = int64(1)
	inst := 0
	res, ok := d.sCache[inst].Get(cacheLiveCheckKey(platform, system, mobile))
	if !ok {
		value, err := d.ConfKv(c, _confLiveCheck)
		if err != nil {
			log.Error("[LiveCheck] get live check error by from source")
			return
		}
		if value == "" {
			log.Error("[LiveCheck] get live check error by source data empty")
			return
		}
		list := &v1pb.GetLiveCheckListResp{}
		err = json.Unmarshal([]byte(value), list)
		if err != nil {
			log.Error("[LiveCheck] get live check error by source data wrong format")
			return
		}
		log.Info("[LiveCheck] live_check list is %v", list)
		switch platform {
		case _platformAndroid:
			for _, v := range list.Android {
				if v.System == system {
					for _, m := range v.Mobile {
						if m == mobile {
							isLive = int64(0)
						}
					}
				}
			}
		case _platformIos:
			for _, v := range list.Ios {
				if v.System == system {
					for _, m := range v.Mobile {
						log.Info("[LiveCheck] range m %v mobile %v", m, mobile)
						if m == mobile {
							isLive = int64(0)
						}
					}
				}
			}
		}
		d.sCache[inst].Put(cacheLiveCheckKey(platform, system, mobile), isLive)
		return
	}
	isLive = res.(int64)
	return
}

// ConfKv get data from cache if miss will call source method, then add to cache.
func (d *Dao) ConfKv(c context.Context, key string) (value string, err error) {
	inst := 0
	res, ok := d.sCache[inst].Get(cacheConfKey(key))
	if !ok {
		log.Info("[LiveCheck] conf cache miss")
		value, err = d.RawConfKv(c, key)
		if err != nil {
			return
		}
		d.sCache[inst].Put(cacheConfKey(key), value)
		return
	}
	log.Info("[LiveCheck] conf cache hit")
	value = res.(string)
	return
}

// SetLiveCheck set live_check conf
func (d *Dao) SetLiveCheck(c context.Context, value string) (err error) {
	err = d.AddOrUpdateConfKv(c, _confLiveCheck, value)
	return
}

// GetLiveCheckList get live_check conf
func (d *Dao) GetLiveCheckList(c context.Context) (value string, err error) {
	value, err = d.RawConfKv(c, _confLiveCheck)
	return
}

// RawConfKv get conf
func (d *Dao) RawConfKv(c context.Context, key string) (value string, err error) {
	row := d.db.QueryRow(c, _findConfkvSQL, key)
	if err = row.Scan(&value); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("[SelectConfKv] row.Scan() error(%v)", err)
	}
	return
}

// AddOrUpdateConfKv add or update conf
func (d *Dao) AddOrUpdateConfKv(c context.Context, key string, value string) (err error) {
	oldValue, err := d.RawConfKv(c, key)
	if err != nil {
		return
	}
	if oldValue != "" {
		//update
		log.Info("[LiveCheck] update db value %v", value)
		if _, err = d.db.Exec(c, _updateConfkvSQL, value, key); err != nil {
			log.Error("[AddOrUpdateConfKv] UpdateConfKv:db.Exec(%v,$v) error(%v)", key, value, err)
		}
		return
	}
	//add
	log.Info("[LiveCheck] add db value %v", value)
	if _, err = d.db.Exec(c, _addConfkvSQL, key, value); err != nil {
		log.Error("[AddOrUpdateConfKv] AddConfKv:db.Exec(%v,$v) error(%v)", key, value, err)
	}
	return
}
