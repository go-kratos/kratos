package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/job/openplatform/article/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	// view
	_viewPrefix = "v_"

	// retry stat
	_retryStatKey = "retry_stat"
	// RetryUpdateStatCache .
	RetryUpdateStatCache = 1
	// RetryUpdateStatDB .
	RetryUpdateStatDB = 2
	// RetryStatCount is retry upper limit.
	RetryStatCount = 10

	// retry cache
	_retryArtCacheKey     = "retry_art_cache"
	_retryGameCacheKey    = "artj_retry_game_cache"
	_retryFlowCacheKey    = "artj_retry_flow_cache"
	_retryDynamicCacheKey = "artj_retry_dynamic_cache"
	// RetryAddArtCache .
	RetryAddArtCache = 1
	// RetryUpdateArtCache .
	RetryUpdateArtCache = 2
	// RetryDeleteArtCache .
	RetryDeleteArtCache = 3
	// RetryDeleteArtRecCache .
	RetryDeleteArtRecCache = 4

	// retry reply
	_retryReplyKey = "retry_reply"
	// retry purge cdn
	_retryCDNKey   = "retry_cdn"
	_recheckArtKey = "recheck_lock_%d"

	// reading start set
	_readPingSet = "art:readping"
	// reading during on some device for some article
	_prefixReadPing = "art:readping:%s:%d"
)

// StatRetry .
type StatRetry struct {
	Action int             `json:"action"`
	Count  int             `json:"count"`
	Data   *artmdl.StatMsg `json:"data"`
}

func viewKey(aid, mid int64, ip string) (key string) {
	if ip == "" {
		// let it pass if ip is empty.
		return
	}
	key = _viewPrefix + strconv.FormatInt(aid, 10) + ip
	if mid != 0 {
		key += strconv.FormatInt(mid, 10)
	}
	return
}

func dupViewKey(aid, mid int64) (key string) {
	return fmt.Sprintf("dv_%v_%v", aid, mid)
}

func recheckKey(aid int64) (key string) {
	return fmt.Sprintf(_recheckArtKey, aid)
}

// pingRedis checks redis healthy.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		PromError("redis:Ping")
		log.Error("redis: conn.Do(SET,PING,PONG) error(%+v)", err)
	}
	return
}

// Intercept intercepts illegal views.
func (d *Dao) Intercept(c context.Context, aid, mid int64, ip string) (ban bool) {
	var (
		err   error
		exist bool
		key   = viewKey(aid, mid, ip)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if key == "" {
		return
	}
	if exist, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		log.Error("conn.Do(EXISTS, %s) error(%+v)", key, err)
		PromError("redis:EXISTS阅读数")
		return
	}
	if exist {
		ban = true
		return
	}
	if err = conn.Send("SET", key, "1"); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", key, err)
		PromError("redis:SET阅读数")
		return
	}
	if err = conn.Send("EXPIRE", key, d.viewCacheTTL); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", key, err)
		PromError("redis:EXPIRE阅读数")
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%+v)", err)
		PromError("redis:阅读数缓存Flush")
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%+v)", err)
			PromError("redis:阅读数缓存Receive")
			return
		}
	}
	return
}

// DupViewIntercept intercepts illegal views.
func (d *Dao) DupViewIntercept(c context.Context, aid, mid int64) (ban bool) {
	if mid == 0 {
		return
	}
	var (
		err   error
		exist bool
		key   = dupViewKey(aid, mid)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if exist, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		log.Error("conn.Do(EXISTS, %s) error(%+v)", key, err)
		PromError("redis:EXISTS连续阅读数")
		return
	}
	if exist {
		ban = true
		return
	}
	if err = conn.Send("SET", key, "1"); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", key, err)
		PromError("redis:SET连续阅读数")
		return
	}
	if err = conn.Send("EXPIRE", key, d.dupViewCacheTTL); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%+v)", key, err)
		PromError("redis:EXPIRE连续阅读数")
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%+v)", err)
		PromError("redis:连续阅读数缓存Flush")
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%+v)", err)
			PromError("redis:连续阅读数缓存Receive")
			return
		}
	}
	return
}

// PushStat pushs failed item to redis.
func (d *Dao) PushStat(c context.Context, retry *StatRetry) (err error) {
	var (
		length int64
		bs     []byte
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if length, err = redis.Int64(conn.Do("LLEN", _retryStatKey)); err != nil {
		log.Error("conn.Do(%s) error(%+v)", _retryStatKey, err)
		PromError("redis:计数重试LLEN")
		return
	}
	cacheLen.State("redis:retry_stat_length", length)
	if bs, err = json.Marshal(retry); err != nil {
		log.Error("json.Marshal(%v) error(%+v)", retry, err)
		PromError("redis:计数重试消息Marshal")
		return
	}
	if _, err = conn.Do("RPUSH", _retryStatKey, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%+v)", bs, err)
		PromError("redis:计数重试RPUSH")
	}
	return
}

// PopStat pops failed item from redis.
func (d *Dao) PopStat(c context.Context) (bs []byte, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", _retryStatKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%+v)", _retryStatKey, err)
		PromError("redis:计数重试LPOP")
	}
	return
}

// PushReply opens article's reply.
func (d *Dao) PushReply(c context.Context, aid, mid int64) (err error) {
	var (
		length int64
		conn   = d.redis.Get(c)
		bs     = []byte(strconv.FormatInt(aid, 10) + "_" + strconv.FormatInt(mid, 10))
	)
	defer conn.Close()
	if length, err = redis.Int64(conn.Do("LLEN", _retryReplyKey)); err != nil {
		log.Error("conn.Do(%s) error(%+v)", _retryReplyKey, err)
		PromError("redis:打开评论重试LLEN")
		return
	}
	cacheLen.State("redis:retry_reply_length", length)
	if _, err = conn.Do("RPUSH", _retryReplyKey, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%+v)", bs, err)
		PromError("redis:打开评论重试RPUSH")
	}
	return
}

// PopReply consume reply's job.
func (d *Dao) PopReply(c context.Context) (aid, mid int64, err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
		v    string
		arr  []string
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", _retryReplyKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%+v)", _retryReplyKey, err)
		PromError("redis:打开评论重试LPOP")
		return
	}
	if v = string(bs); v == "" {
		return
	}
	if arr = strings.Split(v, "_"); len(arr) < 2 {
		log.Error("reply retry param error (%s)", v)
		PromError("redis:打开评论重试消息内容错误")
		return
	}
	aid, _ = strconv.ParseInt(arr[0], 10, 64)
	mid, _ = strconv.ParseInt(arr[1], 10, 64)
	return
}

// PushCDN .
func (d *Dao) PushCDN(c context.Context, file string) (err error) {
	var (
		length int64
		conn   = d.redis.Get(c)
		bs     = []byte(file)
	)
	defer conn.Close()
	if length, err = redis.Int64(conn.Do("LLEN", _retryCDNKey)); err != nil {
		log.Error("conn.Do(%s) error(%+v)", _retryCDNKey, err)
		PromError("redis:重试刷新CDN LLEN")
		return
	}
	cacheLen.State("redis:retry_cdn_length", length)
	if _, err = conn.Do("RPUSH", _retryCDNKey, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%+v)", bs, err)
		PromError("redis:重试刷新CDN RPUSH")
	}
	return
}

// PopCDN .
func (d *Dao) PopCDN(c context.Context) (file string, err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", _retryCDNKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%+v)", _retryCDNKey, err)
		PromError("redis:重试刷新CDN LPOP")
		return
	}
	file = string(bs)
	return
}

// CacheRetry struct of retry cache info.
type CacheRetry struct {
	Action int   `json:"action"`
	Aid    int64 `json:"aid"`
	Mid    int64 `json:"mid"`
	Cid    int64 `json:"cid"`
}

// PushArtCache .
func (d *Dao) PushArtCache(c context.Context, info *CacheRetry) (err error) {
	var (
		length int64
		bs     []byte
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if length, err = redis.Int64(conn.Do("LLEN", _retryArtCacheKey)); err != nil {
		log.Error("conn.Do(%s) error(%+v)", _retryArtCacheKey, err)
		PromError("redis:重试文章缓存LLEN")
		return
	}
	cacheLen.State("redis:retry_art_cache_length", length)
	if bs, err = json.Marshal(info); err != nil {
		log.Error("json.Marshal(%v) error(%+v)", info, err)
		PromError("redis:重试文章缓存Marshal")
		return
	}
	if _, err = conn.Do("RPUSH", _retryArtCacheKey, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%+v)", bs, err)
		PromError("redis:重试文章缓存RPUSH")
	}
	return
}

// PopArtCache .
func (d *Dao) PopArtCache(c context.Context) (bs []byte, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", _retryArtCacheKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%+v)", _retryArtCacheKey, err)
		PromError("redis:重试文章缓存LPOP")
	}
	return
}

// PushGameCache .
func (d *Dao) PushGameCache(c context.Context, info *model.GameCacheRetry) (err error) {
	var (
		length int64
		bs     []byte
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if length, err = redis.Int64(conn.Do("LLEN", _retryGameCacheKey)); err != nil {
		log.Error("conn.Do(%s) error(%+v)", _retryGameCacheKey, err)
		PromError("redis:重试游戏缓存LLEN")
		return
	}
	cacheLen.State("redis:retry_game_cache_length", length)
	if bs, err = json.Marshal(info); err != nil {
		log.Error("json.Marshal(%v) error(%+v)", info, err)
		PromError("redis:重试游戏缓存Marshal")
		return
	}
	if _, err = conn.Do("RPUSH", _retryGameCacheKey, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%+v)", bs, err)
		PromError("redis:重试游戏缓存RPUSH")
	}
	return
}

// PopGameCache .
func (d *Dao) PopGameCache(c context.Context) (res *model.GameCacheRetry, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	var bs []byte
	if bs, err = redis.Bytes(conn.Do("LPOP", _retryGameCacheKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%+v)", _retryGameCacheKey, err)
		PromError("redis:重试游戏缓存LPOP")
		return
	}
	res = new(model.GameCacheRetry)
	if err = json.Unmarshal(bs, res); err != nil {
		log.Error("redis.Unmarshal(%s) error(%+v)", bs, err)
		PromError("redis:解析游戏缓存")
	}
	return
}

// PushFlowCache .
func (d *Dao) PushFlowCache(c context.Context, info *model.FlowCacheRetry) (err error) {
	var (
		length int64
		bs     []byte
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if length, err = redis.Int64(conn.Do("LLEN", _retryFlowCacheKey)); err != nil {
		log.Error("conn.Do(%s) error(%+v)", _retryFlowCacheKey, err)
		PromError("redis:重试flow缓存LLEN")
		return
	}
	cacheLen.State("redis:retry_flow_cache_length", length)
	if bs, err = json.Marshal(info); err != nil {
		log.Error("json.Marshal(%v) error(%+v)", info, err)
		PromError("redis:重试flow缓存Marshal")
		return
	}
	if _, err = conn.Do("RPUSH", _retryFlowCacheKey, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%+v)", bs, err)
		PromError("redis:重试flow缓存RPUSH")
	}
	return
}

// PopFlowCache .
func (d *Dao) PopFlowCache(c context.Context) (res *model.FlowCacheRetry, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	var bs []byte
	if bs, err = redis.Bytes(conn.Do("LPOP", _retryFlowCacheKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%+v)", _retryFlowCacheKey, err)
		PromError("redis:重试flow缓存LPOP")
		return
	}
	res = new(model.FlowCacheRetry)
	if err = json.Unmarshal(bs, res); err != nil {
		log.Error("redis.Unmarshal(%s) error(%+v)", bs, err)
		PromError("redis:解析flow缓存")
	}
	return
}

// PushDynamicCache put dynamic to redis
func (d *Dao) PushDynamicCache(c context.Context, info *model.DynamicCacheRetry) (err error) {
	var (
		length int64
		bs     []byte
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if length, err = redis.Int64(conn.Do("LLEN", _retryDynamicCacheKey)); err != nil {
		log.Error("conn.Do(%s) error(%+v)", _retryDynamicCacheKey, err)
		PromError("redis:重试dynamic缓存LLEN")
		return
	}
	cacheLen.State("redis:retry_dynamic_cache_length", length)
	if bs, err = json.Marshal(info); err != nil {
		log.Error("json.Marshal(%v) error(%+v)", info, err)
		PromError("redis:重试dynamic缓存Marshal")
		return
	}
	if _, err = conn.Do("RPUSH", _retryDynamicCacheKey, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%+v)", bs, err)
		PromError("redis:重试dynamic缓存RPUSH")
	}
	return
}

// PopDynamicCache .
func (d *Dao) PopDynamicCache(c context.Context) (res *model.DynamicCacheRetry, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	var bs []byte
	if bs, err = redis.Bytes(conn.Do("LPOP", _retryDynamicCacheKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%+v)", _retryDynamicCacheKey, err)
		PromError("redis:重试dynamic缓存LPOP")
		return
	}
	res = new(model.DynamicCacheRetry)
	if err = json.Unmarshal(bs, res); err != nil {
		log.Error("redis.Unmarshal(%s) error(%+v)", bs, err)
		PromError("redis:解析dynamic缓存")
	}
	return
}

// GetRecheckCache get recheck info from redis
func (d *Dao) GetRecheckCache(c context.Context, aid int64) (isRecheck bool, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	key := recheckKey(aid)
	if isRecheck, err = redis.Bool(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("redis.BOOL(conn.Do(GET, %s)) error(%+v)", key, err)
		PromError("redis:获取回查缓存")
		return
	}
	return
}

// SetRecheckCache set recheck info to redis
func (d *Dao) SetRecheckCache(c context.Context, aid int64) (isRecheck bool, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	key := recheckKey(aid)
	if _, err = conn.Do("SETEX", key, 604800, true); err != nil {
		log.Error("redis.BOOL(conn.Do(SETEX, %s)) error(%+v)", key, err)
		PromError("redis:设置回查缓存")
		return
	}
	return
}

func readPingSetKey() string {
	return _readPingSet
}

func readPingKey(buvid string, aid int64) string {
	return fmt.Sprintf(_prefixReadPing, buvid, aid)
}

// ReadPingSet 获取所有阅读记录（不删除）
func (d *Dao) ReadPingSet(c context.Context) (res []*model.Read, err error) {
	var (
		key    = readPingSetKey()
		conn   = d.artRedis.Get(c)
		tmpRes []string
		tmpArr []string
	)
	defer conn.Close()

	if err = conn.Send("SMEMBERS", key); err != nil {
		log.Error("conn.Send(SMEMBERS, %s) error(%+v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	if tmpRes, err = redis.Strings(conn.Receive()); err != nil {
		log.Error("conn.Receive error(%+v)", err)
		return
	}
	for _, tmp := range tmpRes {
		tmpArr = strings.Split(tmp, "|")
		if len(tmpArr) != 6 {
			log.Error("redis key(%s)存在脏数据(%s)", key, tmp)
			if _, err = conn.Do("SREM", key, tmp); err != nil {
				log.Error("d.Redis.SREM error(%+v), set(%s), key(%s)", err, key, tmp)
			}
			continue
		}
		read := &model.Read{
			Buvid:   tmpArr[0],
			EndTime: 0,
		}
		read.Aid, _ = strconv.ParseInt(tmpArr[1], 10, 64)
		read.Mid, _ = strconv.ParseInt(tmpArr[2], 10, 64)
		read.IP = tmpArr[3]
		read.StartTime, _ = strconv.ParseInt(tmpArr[4], 10, 64)
		read.From = tmpArr[5]
		res = append(res, read)
	}
	return
}

// ReadPing 获取上次阅读心跳时间，不存在则返回0
func (d *Dao) ReadPing(c context.Context, buvid string, aid int64) (last int64, err error) {
	var (
		key  = readPingKey(buvid, aid)
		conn = d.artRedis.Get(c)
	)
	defer conn.Close()
	if last, err = redis.Int64(conn.Do("GET", key)); err != nil && err != redis.ErrNil {
		log.Error("conn.Do(GET, %s) error(%+v)", key, err)
		return
	}
	err = nil
	return
}

// DelReadPingSet 删除阅读记录缓存
func (d *Dao) DelReadPingSet(c context.Context, read *model.Read) (err error) {
	if read == nil {
		return
	}
	var (
		elemKey = readPingKey(read.Buvid, read.Aid)
		setKey  = readPingSetKey()
		value   = fmt.Sprintf("%s|%d|%d|%s|%d|%s", read.Buvid, read.Aid, read.Mid, read.IP, read.StartTime, read.From)
		conn    = d.artRedis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", elemKey); err != nil {
		log.Error("conn.Do(DEL, %s) error(%+v)", elemKey, err)
		return
	}
	if _, err = conn.Do("SREM", setKey, value); err != nil {
		log.Error("conn.Do(SREM, %s, %s) error(%+v)", setKey, value, err)
		return
	}
	return
}
