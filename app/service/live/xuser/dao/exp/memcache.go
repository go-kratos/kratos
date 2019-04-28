package exp

import (
	"context"
	expModel "go-common/app/service/live/xuser/model/exp"
	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"
	"go-common/library/sync/errgroup"
	"math/rand"
	"strconv"
	"time"
)

// redis cache
const (
	_prefixExp        = "json_e_" // 用户经验cache key,json协议
	_emptyExpire      = 20 * 24 * 3600
	_errorMcLogPrefix = "xuser.exp.dao.memcache"
	_promGetSuccess   = "xuser_exp_mc:获取用户经验cache成功"
	_promDelSuccess   = "xuser_exp_mc:成功删除用户经验cache"
	_promDelErr       = "xuser_exp_mc:删除用户经验cache失败"
	_promGetErr       = "xuser_exp_mc:批量获取用户经验key失败"
	_promScanErr      = "xuser_exp_mc:解析用户经验key失败"
	// _recExpire   = 5 * 24 * 3600
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

func expKey(mid int64) string {
	return _prefixExp + strconv.FormatInt(mid, 10)
}

// SetExpListCache 批量设置用户经验cache
func (d *Dao) SetExpListCache(c context.Context, expList map[int64]*expModel.LevelInfo) (err error) {
	return d.setExpListCache(c, expList)
}

// DelExpFromMemCache 删除获取用户经验cache,不支持批量
func (d *Dao) DelExpFromMemCache(c context.Context, mid int64) (err error) {
	return d.delExpFromMemCache(c, expKey(mid))
}

// GetExpFromMemCache 对外接口
func (d *Dao) GetExpFromMemCache(ctx context.Context, mids []int64) (expList map[int64]*expModel.LevelInfo, missedUids []int64, err error) {
	return d.getExpFromMemCache(ctx, mids)
}

// getExpFromMemCache 批量获取用户经验cache
func (d *Dao) getExpFromMemCache(ctx context.Context, mids []int64) (expList map[int64]*expModel.LevelInfo, arrayMissedUids []int64, err error) {
	var expKeys []string
	expList = make(map[int64]*expModel.LevelInfo)
	mapMissedUids := make(map[int64]bool)
	arrayMissedUids = make([]int64, 0)
	for _, uid := range mids {
		expKeys = append(expKeys, expKey(uid))
		mapMissedUids[uid] = true
	}
	group := errgroup.Group{}
	group.Go(func() error {
		conn := d.memcache.Get(context.TODO())
		defer conn.Close()
		resp, err := conn.GetMulti(expKeys)
		if err != nil {
			PromError(_promGetErr)
			log.Error(_errorMcLogPrefix+"|conn.Gets(%v) error(%v)", expKeys, err)
			return nil
		}
		// miss uids
		for _, item := range resp {
			element := &expModel.LevelInfo{}
			if err = conn.Scan(item, &element); err != nil {
				PromError(_promScanErr)
				log.Error(_errorMcLogPrefix+"|item.Scan(%s) error(%v)", item.Value, err)
				return err
			}
			expList[element.UID] = element
			mapMissedUids[element.UID] = false
		}
		return nil
	})
	group.Wait()
	for k, v := range mapMissedUids {
		if v {
			arrayMissedUids = append(arrayMissedUids, k)
		}

	}
	PromInfo(_promGetSuccess)
	return
}

func (d *Dao) delExpFromMemCache(ctx context.Context, key string) (err error) {
	conn := d.memcache.Get(ctx)
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			log.Error(_errorMcLogPrefix+"|Delete(%s) error(%v)", key, err)
			PromError(_promDelErr)
		}
	}
	PromInfo(_promDelSuccess)
	conn.Close()
	return
}

func (d *Dao) setExpListCache(ctx context.Context, expList map[int64]*expModel.LevelInfo) (err error) {
	expire := d.getExpire()
	rand.Seed(time.Now().UnixNano())
	expire = expire + rand.Int31n(3600)
	conn := d.memcache.Get(context.TODO())
	defer conn.Close()
	for _, v := range expList {
		item := &gmc.Item{Key: expKey(v.UID),
			Object:     &expModel.LevelInfo{UID: v.UID, UserLevel: v.UserLevel, AnchorLevel: v.AnchorLevel, CTime: v.CTime, MTime: v.MTime},
			Expiration: expire, Flags: gmc.FlagJSON}
		err := conn.Set(item)
		if err != nil {
			PromError("mc:设置上报")
			log.Error(_errorMcLogPrefix+"|conn.Set(%v) error(%v)", expKey(v.UID), err)
			continue
		}
	}
	return

}
