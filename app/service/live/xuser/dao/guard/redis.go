package guard

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	dahanghaiModel "go-common/app/service/live/xuser/model/dhh"
	gmc "go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/stat/prom"
	"math/rand"
	"strconv"
	"time"
)

// redis cache
const (
	_prefixUID           = "live_user:guard:uid:v1:"       // 用户侧key
	_prefixTopList       = "GOVERNOR_SHOW_TID:"            // 最近购买总督key
	_prefixAnchorID      = "live_user:guard:target_id:v1:" // 主播侧key
	_emptyExpire         = 3600
	_errorRedisLogPrefix = "xuser.dahanghai.dao.redis"
	_promGetSuccess      = "xuser_dahanghai_redis:获取用户大航海cache成功"
	_promDelSuccess      = "xuser_dahanghai_redis:成功删除用户大航海cache"
	_promDelErr          = "xuser_dahanghai_redis:删除用户大航海cache失败"
	_promGetErr          = "xuser_dahanghai_redis:批量获取用户大航海key失败"
	// _promScanErr         = "xuser_dahanghai_redis:解析用户大航海key失败"
)

var (
	errorsCount    = prom.BusinessErrCount
	infosCount     = prom.BusinessInfoCount
	cacheHitCount  = prom.CacheHit
	cacheMissCount = prom.CacheMiss
)

// PromError prometheus error count.
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo prometheus info count.
func PromInfo(name string) {
	infosCount.Incr(name)
}

// PromCacheHit prometheus cache hit count.
func PromCacheHit(name string) {
	cacheHitCount.Incr(name)
}

// PromCacheMiss prometheus cache hit count.
func PromCacheMiss(name string) {
	cacheMissCount.Incr(name)
}

func dahanghaiUIDKey(mid int64) string {
	return _prefixUID + strconv.FormatInt(mid, 10)
}

func guardAnchorUIDKey(mid int64) string {
	return _prefixAnchorID + strconv.FormatInt(mid, 10)
}

func recentGuardTopKey(mid int64) string {
	return _prefixTopList + strconv.FormatInt(mid, 10)
}

// SetDHHListCache ... 批量设置用户守护cache
func (d *GuardDao) SetDHHListCache(c context.Context, dhhList []dahanghaiModel.DaHangHaiRedis2, uid int64) (err error) {
	return d.setDHHListCache(c, dhhList, uid)
}

// SetAnchorGuardListCache ... 批量设置主播维度守护信息cache
func (d *GuardDao) SetAnchorGuardListCache(c context.Context, dhhList []dahanghaiModel.DaHangHaiRedis2, uid int64) (err error) {
	return d.setAnchorGuardListCache(c, dhhList, uid)
}

// DelDHHFromRedis 删除获取用户守护cache,不支持批量
func (d *GuardDao) DelDHHFromRedis(c context.Context, mid int64) (err error) {
	return d.delDHHFromRedis(c, dahanghaiUIDKey(mid))
}

// GetUIDAllGuardFromRedis 获取单个用户的全量守护信息
func (d *GuardDao) GetUIDAllGuardFromRedis(ctx context.Context, mids []int64) (dhhList []*dahanghaiModel.DaHangHaiRedis2, err error) {
	return d.getUIDAllGuardFromRedis(ctx, mids)
}

// GetUIDAllGuardFromRedisBatch 获取批量用户的全量守护信息
func (d *GuardDao) GetUIDAllGuardFromRedisBatch(ctx context.Context, mids []int64) (dhhList []*dahanghaiModel.DaHangHaiRedis2, err error) {
	return d.getUIDAllGuardFromRedisBatch(ctx, mids)
}

// GetAnchorAllGuardFromRedis 获取单个主播的全量被守护信息(同一个主播仅获取最高级别)
func (d *GuardDao) GetAnchorAllGuardFromRedis(ctx context.Context, anchorUIDs []int64) (dhhList []*dahanghaiModel.DaHangHaiRedis2, err error) {
	return d.getAnchorAllGuardFromRedis(ctx, anchorUIDs)
}

// GetGuardTopListCache 获取单个用户的全量守护信息
func (d *GuardDao) GetGuardTopListCache(ctx context.Context, uid int64) (dhhList []*dahanghaiModel.DaHangHaiRedis2, err error) {
	return d.getGuardTopListCache(ctx, uid)
}

// GetAnchorRecentTopGuardCache 获取单个主播最近的总督信息
func (d *GuardDao) GetAnchorRecentTopGuardCache(ctx context.Context, uid int64) (resp map[int64]int64, err error) {
	return d.getAnchorRecentTopGuardCache(ctx, uid)
}

func (d *GuardDao) getGuardTopListCache(ctx context.Context, uid int64) (dhhList []*dahanghaiModel.DaHangHaiRedis2, err error) {
	return
}

// getUIDAllGuardFromRedis 批量获取用户cache
func (d *GuardDao) getUIDAllGuardFromRedis(ctx context.Context, mids []int64) (dhhList []*dahanghaiModel.DaHangHaiRedis2, err error) {
	var (
		conn        = d.redis.Get(ctx)
		args        = redis.Args{}
		cacheResult [][]byte
	)
	defer conn.Close()
	for _, uid := range mids {
		args = args.Add(dahanghaiUIDKey(uid))
	}
	if cacheResult, err = redis.ByteSlices(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			PromError(_promGetErr)
			log.Error(_errorRedisLogPrefix+"|conn.MGET(%v) error(%v)", args, err)
			err = errors.Wrapf(err, "redis.StringMap(conn.Do(MGET,%v)", args)
		}
		return
	}
	dhhList = make([]*dahanghaiModel.DaHangHaiRedis2, 0)
	dhhListSingle := &dahanghaiModel.DaHangHaiRedis2{}
	if len(cacheResult) > 0 {
		for k, v := range cacheResult {
			if v == nil {
				return nil, nil
			}
			if len(v) > 0 {
				if err = json.Unmarshal([]byte(v), &dhhList); err != nil {
					if err = json.Unmarshal([]byte(v), &dhhListSingle); err != nil {
						log.Error("[dao.dahanghai.cache|GetDHHFromRedis] json.Unmarshal rawInfo error(%v), uid(%d), reply(%s)",
							err, k, v)
						return nil, nil
					}
					dhhList = append(dhhList, dhhListSingle)
				}
			}
		}
	}
	PromInfo(_promGetSuccess)
	return
}

// getUIDAllGuardFromRedis 批量获取用户cache
func (d *GuardDao) getUIDAllGuardFromRedisBatch(ctx context.Context, mids []int64) (dhhList []*dahanghaiModel.DaHangHaiRedis2, err error) {
	var (
		conn        = d.redis.Get(ctx)
		args        = redis.Args{}
		cacheResult [][]byte
	)
	defer conn.Close()
	for _, uid := range mids {
		args = args.Add(dahanghaiUIDKey(uid))
	}
	if cacheResult, err = redis.ByteSlices(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			PromError(_promGetErr)
			log.Error(_errorRedisLogPrefix+"|conn.MGET(%v) error(%v)", args, err)
			err = errors.Wrapf(err, "redis.StringMap(conn.Do(MGET,%v)", args)
		}
		return
	}
	dhhList = make([]*dahanghaiModel.DaHangHaiRedis2, 0)
	if len(cacheResult) > 0 {
		for k, v := range cacheResult {
			if v == nil {
				continue
			}
			dhhListLoop := make([]*dahanghaiModel.DaHangHaiRedis2, 0)
			dhhListSingle := &dahanghaiModel.DaHangHaiRedis2{}
			if len(v) > 0 {
				if err = json.Unmarshal([]byte(v), &dhhListLoop); err != nil {
					if err = json.Unmarshal([]byte(v), &dhhListSingle); err != nil {
						log.Error("[dao.dahanghai.cache|GetDHHFromRedis] json.Unmarshal rawInfo error(%v), uid(%d), reply(%s)",
							err, k, v)
						return nil, nil
					}
					dhhList = append(dhhList, dhhListSingle)
				} else {
					dhhList = append(dhhList, dhhListLoop...)
				}
			}
		}

	}
	PromInfo(_promGetSuccess)
	return
}

// getAnchorAllGuardFromRedis 批量获取用户cache
func (d *GuardDao) getAnchorAllGuardFromRedis(ctx context.Context, mids []int64) (dhhList []*dahanghaiModel.DaHangHaiRedis2, err error) {
	var (
		conn        = d.redis.Get(ctx)
		args        = redis.Args{}
		cacheResult [][]byte
	)
	defer conn.Close()
	for _, uid := range mids {
		args = args.Add(guardAnchorUIDKey(uid))
	}
	if cacheResult, err = redis.ByteSlices(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			PromError(_promGetErr)
			log.Error(_errorRedisLogPrefix+"|conn.MGET(%v) error(%v)", args, err)
			err = errors.Wrapf(err, "redis.StringMap(conn.Do(MGET,%v)", args)
		}
		return
	}
	dhhList = make([]*dahanghaiModel.DaHangHaiRedis2, 0)
	dhhListSingle := &dahanghaiModel.DaHangHaiRedis2{}
	if len(cacheResult) > 0 {
		for k, v := range cacheResult {
			if v == nil {
				return nil, nil
			}
			if len(v) > 0 {
				if err = json.Unmarshal([]byte(v), &dhhList); err != nil {
					if err = json.Unmarshal([]byte(v), &dhhListSingle); err != nil {
						log.Error("[dao.dahanghai.cache|GetDHHFromRedis] json.Unmarshal rawInfo error(%v), uid(%d), reply(%s)",
							err, k, v)
						return nil, nil
					}
					dhhList = append(dhhList, dhhListSingle)
				}
			}
		}
	}
	PromInfo(_promGetSuccess)
	return
}

func (d *GuardDao) delDHHFromRedis(ctx context.Context, key string) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	if err == gmc.ErrNotFound {
		err = nil
	} else {
		log.Error(_errorRedisLogPrefix+"|Delete(%s) error(%v)", key, err)
		PromError(_promDelErr)
	}
	PromInfo(_promDelSuccess)
	conn.Close()
	return
}

func (d *GuardDao) setDHHListCache(ctx context.Context, dhhList []dahanghaiModel.DaHangHaiRedis2, uid int64) (err error) {
	expire := d.getExpire()
	rand.Seed(time.Now().UnixNano())
	expire = expire + rand.Int31n(60)
	var (
		argsMid = redis.Args{}
		conn    = d.redis.Get(ctx)
		dhhJSON []byte
	)
	defer conn.Close()

	key := dahanghaiUIDKey(uid)
	if len(dhhList) == 0 {
		argsMid = argsMid.Add(key).Add("")
	} else {
		dhhJSON, err = json.Marshal(dhhList)
		if err != nil {
			log.Error("[dao.dahanghai.cache|GetDHHFromRedis] json.Marshal rawInfo error(%v), uid(%d)", err, uid)
			return
		}
		argsMid = argsMid.Add(key).Add(string(dhhJSON))
	}

	if err = conn.Send("SET", argsMid...); err != nil {
		err = errors.Wrap(err, "conn.Send(SET) error")
		return
	}
	rand.Seed(time.Now().UnixNano())
	expire = expire + rand.Int31n(60)
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("setDHHListCache conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	return
}

func (d *GuardDao) setAnchorGuardListCache(ctx context.Context, dhhList []dahanghaiModel.DaHangHaiRedis2, uid int64) (err error) {
	expire := d.getExpire()
	rand.Seed(time.Now().UnixNano())
	expire = expire + rand.Int31n(60)
	var (
		argsMid = redis.Args{}
		conn    = d.redis.Get(ctx)
		dhhJSON []byte
	)
	defer conn.Close()

	key := guardAnchorUIDKey(uid)
	if len(dhhList) == 0 {
		argsMid = argsMid.Add(key).Add("")
	} else {
		dhhJSON, err = json.Marshal(dhhList)
		if err != nil {
			log.Error("[dao.dahanghai.cache|setAnchorGuardListCache] json.Marshal rawInfo error(%v), uid(%d)", err, uid)
			return
		}
		argsMid = argsMid.Add(key).Add(string(dhhJSON))
	}

	if err = conn.Send("SET", argsMid...); err != nil {
		err = errors.Wrap(err, "conn.Send(SET) error")
		return
	}
	rand.Seed(time.Now().UnixNano())
	expire = expire + rand.Int31n(60)
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("setAnchorGuardListCache conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	return
}

func (d *GuardDao) getAnchorRecentTopGuardCache(ctx context.Context, uid int64) (resp map[int64]int64, err error) {
	resp = make(map[int64]int64)
	nowTime := time.Now().Unix()
	cacheKey := recentGuardTopKey(uid)
	var (
		conn = d.redis.Get(ctx)
	)
	values, err := redis.Values(conn.Do("ZRANGEBYSCORE", cacheKey, nowTime, "INF", "WITHSCORES"))
	if err != nil {
		log.Error("getAnchorRecentTopGuardCache.conn.Do(ZRANGEBYSCORE %v) error(%v)", cacheKey, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var aid, unix int64
	for len(values) > 0 {
		if values, err = redis.Scan(values, &aid, &unix); err != nil {
			log.Error("getAnchorRecentTopGuardCache.redis.Scan(%v) error(%v)", values, err)
			return
		}
		resp[aid] = unix
	}
	return
}
