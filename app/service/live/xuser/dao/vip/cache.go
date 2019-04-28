package vip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go-common/app/service/live/xuser/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"time"
)

// redis cache
const (
	_userInfoRedisKey = "us:infoo_v2:%d" // 用户缓存key prefix
	_vipFieldName     = "vip"            // v3 vip attr field
	_levelFieldName   = "level"          // v2 level attr field, Todo: remove level attr
	_userExpired      = 86400            // user cache expire time
)

type vipCache struct {
	Vip      interface{} `json:"vip"`
	VipTime  string      `json:"vip_time"`
	Svip     interface{} `json:"svip"`
	SvipTime string      `json:"svip_time"`
}

// GetVipFromCache get user vip info from cache
func (d *Dao) GetVipFromCache(ctx context.Context, uid int64) (info *model.VipInfo, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	reply, err := redis.String(conn.Do("HGET", getUserCacheKey(uid), _vipFieldName))
	if err != nil {
		if err == redis.ErrNil {
			// key or field not exists, return nil, nil, back to db
			log.Info("[dao.vip.cache|GetVipFromCache] cache key or field not exists err(%v), uid(%d)", err, uid)
			return nil, nil
		}
		log.Error("[dao.vip.cache|GetVipFromCache] hget error(%v), uid(%d)", err, uid)
		return
	}
	if reply == "" {
		return nil, nil
	}

	// ===== begin eat others' dog food =====
	// 1.兼容缓存中vip/svip可能是int or string的问题
	rawInfo := &vipCache{}
	if err = json.Unmarshal([]byte(reply), rawInfo); err != nil {
		log.Error("[dao.vip.cache|GetVipFromCache] json.Unmarshal rawInfo error(%v), uid(%d), reply(%s)",
			err, uid, reply)
		// parse cache json error, return nil, nil, back to db and restore cache
		return nil, nil
	}
	if info, err = d.formatVipCache(rawInfo); err != nil {
		log.Error("[dao.vip.cache|GetVipFromCache] format rawInfo error(%v), uid(%d), reply(%s)", err, uid, reply)
		return nil, nil
	}

	// 2.注意!!! cache里的vip_time/svip_time不一定正确，可能含有已经过期的time
	currentTime := time.Now().Unix()
	// vip time
	if info.Vip, err = d.checkVipTime(info.VipTime, info.Vip, currentTime); err != nil {
		log.Error("[dao.vip.cache|GetVipFromCache] check vip time error(%v), uid(%d), info(%v), reply(%s)",
			err, uid, info, reply)
		return nil, nil
	}
	if info.Svip, err = d.checkVipTime(info.SvipTime, info.Svip, currentTime); err != nil {
		log.Error("[dao.vip.cache|GetVipFromCache] check svip time error(%v), uid(%d), info(%v), reply(%s)",
			err, uid, info, reply)
		return nil, nil
	}
	// ===== end =====

	return
}

// formatVipCache 转换vip/svip的格式
func (d *Dao) formatVipCache(info *vipCache) (v *model.VipInfo, err error) {
	v = &model.VipInfo{
		VipTime:  info.VipTime,
		SvipTime: info.SvipTime,
	}
	if v.Vip, err = toInt(info.Vip); err != nil {
		return
	}
	if v.Svip, err = toInt(info.Svip); err != nil {
		return
	}

	// format info struct
	v = d.initInfo(v)

	return
}

// checkVipTime 检查缓存中vip_time/svip_time是否过期
func (d *Dao) checkVipTime(t string, f int, compare int64) (int, error) {
	if t == model.TimeEmpty {
		if f != 0 {
			return 0, errors.New("empty time with not zero flag.")
		}
	} else {
		vt, err := time.Parse(model.TimeNano, t)
		if err != nil {
			return 0, errors.New("time parse error.")
		}
		if vt.Unix() <= compare {
			return 0, nil
		}
	}
	return f, nil
}

// SetVipCache set vip to cache
func (d *Dao) SetVipCache(ctx context.Context, uid int64, info *model.VipInfo) (err error) {
	var vipJson []byte
	conn := d.redis.Get(ctx)
	key := getUserCacheKey(uid)
	defer conn.Close()
	// format info struct
	info = d.initInfo(info)
	// format info json string
	vipJson, err = json.Marshal(info)
	if err != nil {
		log.Error("[dao.vip.cache|SetVipCache] json.Marshal error(%v), uid(%d), info(%v)", err, uid, info)
		// if marshal error, clear cache
		goto CLEAR
	}
	_, err = conn.Do("HSET", key, _vipFieldName, string(vipJson))
	if err != nil {
		log.Error("[dao.vip.cache|SetVipCache] HSET error(%v), uid(%d), info(%v)", err, uid, info)
		// if hset error, clear cache
		goto CLEAR
	}
	_, err = conn.Do("EXPIRE", key, _userExpired)
	if err != nil {
		log.Error("[dao.vip.cache|SetVipCache] EXPIRE error(%v), uid(%d), info(%v)", err, uid, info)
		// if set expire error, clear cache
		goto CLEAR
	}
	return

CLEAR:
	log.Error("[dao.vip.cache|SetVipCache] set error, aysnc clear, uid(%d), info(%v)", uid, info)
	go d.ClearCache(metadata.WithContext(ctx), uid)
	return
}

// ClearCache clear user's vip and level field cache
// Todo: remove level attr
func (d *Dao) ClearCache(ctx context.Context, uid int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	key := getUserCacheKey(uid)
	_, err = conn.Do("HDEL", key, _vipFieldName, _levelFieldName)
	if err != nil {
		err = errors.Wrapf(err, "conn.Do(HDEL, %s, %s, %s)", key, _vipFieldName, _levelFieldName)
		log.Error("[dao.vip.cache|ClearCache] hdel uid(%d) vip and level attr err(%v)", uid, err)
	}
	return
}

func getUserCacheKey(uid int64) string {
	return fmt.Sprintf(_userInfoRedisKey, uid)
}
