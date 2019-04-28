package dao

import (
	"context"
	"fmt"
	"strconv"
	gtime "time"

	"go-common/app/service/main/relation/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	_prefixFollowings         = "at_"        // key of public following with tags datas.
	_prefixMonitor            = "rs_mo_list" // key of monitor
	_prefixRecentFollower     = "rf_"        // recent follower sorted set
	_prefixRecentFollowerTime = "rft_"       // recent follower sorted set
	_prefixDailyNotifyCount   = "dnc_%d_%s"  // daily new-follower notificaiton count
	_notifyCountExpire        = 24 * 3600    // notify count scope is daily
)

func followingsKey(mid int64) string {
	return _prefixFollowings + strconv.FormatInt(mid, 10)
}

func monitorKey() string {
	return _prefixMonitor
}

func recentFollower(mid int64) string {
	return _prefixRecentFollower + strconv.FormatInt(mid, 10)
}

func recentFollowerNotify(mid int64) string {
	return _prefixRecentFollowerTime + strconv.FormatInt(mid%10000, 10)
}

func dailyNotifyCount(mid int64, date gtime.Time) string {
	// _cacheShard  作为sharding
	return fmt.Sprintf(_prefixDailyNotifyCount, mid%_cacheShard, date.Format("2006-01-02"))
}

// pingRedis ping redis.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		log.Error("conn.Do(SET,PING,PONG) error(%v)", err)
	}
	conn.Close()
	return
}

// SetFollowingsCache set followings cache.
func (d *Dao) SetFollowingsCache(c context.Context, mid int64, followings []*model.Following) (err error) {
	key := followingsKey(mid)
	args := redis.Args{}.Add(key)
	expire := d.redisExpire
	if len(followings) == 0 {
		expire = 7200
	}
	ef, _ := d.encode(0, 0, nil, 0)
	args = args.Add(0, ef)
	for i := 0; i < len(followings); i++ {
		var ef []byte
		if ef, err = d.encode(followings[i].Attribute, followings[i].MTime, followings[i].Tag, followings[i].Special); err != nil {
			return
		}
		args = args.Add(followings[i].Mid, ef)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("HMSET", args...); err != nil {
		log.Error("conn.Send(HMSET, %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 3; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() %d error(%v)", i+1, err)
			break
		}
	}
	return
}

// AddFollowingCache add following cache.
func (d *Dao) AddFollowingCache(c context.Context, mid int64, following *model.Following) (err error) {
	var (
		ok  bool
		key = followingsKey(mid)
	)
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpire)); err != nil {
		log.Error("redis.Bool(conn.Do(EXPIRE, %s)) error(%v)", key, err)
	} else if ok {
		var ef []byte
		if ef, err = d.encode(following.Attribute, following.MTime, following.Tag, following.Special); err != nil {
			return
		}
		if _, err = conn.Do("HSET", key, following.Mid, ef); err != nil {
			log.Error("conn.Do(HSET, %s, %d) error(%v)", key, following.Mid, err)
		}
	}
	conn.Close()
	return
}

// DelFollowing del following cache.
func (d *Dao) DelFollowing(c context.Context, mid int64, following *model.Following) (err error) {
	var (
		ok  bool
		key = followingsKey(mid)
	)
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpire)); err != nil {
		log.Error("redis.Bool(conn.Do(EXPIRE, %s)) error(%v)", key, err)
	} else if ok {
		if _, err = conn.Do("HDEL", key, following.Mid); err != nil {
			log.Error("conn.Do(HDEL, %s, %d) error(%v)", key, following.Mid, err)
		}
	}
	conn.Close()
	return
}

// FollowingsCache get followings cache.
func (d *Dao) FollowingsCache(c context.Context, mid int64) (followings []*model.Following, err error) {
	key := followingsKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	tmp, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return
	}
	if err == nil && len(tmp) > 0 {
		for k, v := range tmp {
			if mid, err = strconv.ParseInt(k, 10, 64); err != nil {
				return
			}
			if mid <= 0 {
				continue
			}
			vf := &model.FollowingTags{}
			if err = d.decode([]byte(v), vf); err != nil {
				//todo
				return
			}
			followings = append(followings, &model.Following{
				Mid:       mid,
				Attribute: vf.Attr,
				Tag:       vf.TagIds,
				MTime:     vf.Ts,
				Special:   vf.Special,
			})
		}
	}
	return
}

// DelFollowingsCache delete followings cache.
func (d *Dao) DelFollowingsCache(c context.Context, mid int64) (err error) {
	key := followingsKey(mid)
	conn := d.redis.Get(c)
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// RelationsCache relations cache.
func (d *Dao) RelationsCache(c context.Context, mid int64, fids []int64) (resMap map[int64]*model.Following, err error) {
	var retRedis [][]byte
	key := followingsKey(mid)
	args := redis.Args{}.Add(key)
	for _, fid := range fids {
		args = args.Add(fid)
	}
	args.Add(0)
	conn := d.redis.Get(c)
	defer conn.Close()
	if retRedis, err = redis.ByteSlices(conn.Do("HMGET", args...)); err != nil {
		log.Error("redis.Int64s(conn.DO(HMGET, %v)) error(%v)", args, err)
		return
	}
	resMap = make(map[int64]*model.Following)
	for index, fid := range fids {
		if retRedis[index] == nil {
			continue
		}
		v := &model.FollowingTags{}
		if err = d.decode(retRedis[index], v); err != nil {
			return
		}
		resMap[fid] = &model.Following{
			Mid:       fid,
			Attribute: v.Attr,
			Tag:       v.TagIds,
			MTime:     v.Ts,
			Special:   v.Special,
		}
	}
	return
}

// encode
func (d *Dao) encode(attribute uint32, mtime time.Time, tagids []int64, special int32) (res []byte, err error) {
	ft := &model.FollowingTags{Attr: attribute, Ts: mtime, TagIds: tagids, Special: special}
	return ft.Marshal()
}

// decode
func (d *Dao) decode(src []byte, v *model.FollowingTags) (err error) {
	return v.Unmarshal(src)
}

// MonitorCache monitor cache
func (d *Dao) MonitorCache(c context.Context, mid int64) (exist bool, err error) {
	key := monitorKey()
	conn := d.redis.Get(c)
	if exist, err = redis.Bool(conn.Do("SISMEMBER", key, mid)); err != nil {
		log.Error("redis.Bool(conn.Do(SISMEMBER, %s, %d)) error(%v)", key, mid, err)
	}
	conn.Close()
	return
}

// SetMonitorCache set monitor cache
func (d *Dao) SetMonitorCache(c context.Context, mid int64) (err error) {
	var (
		key  = monitorKey()
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("SADD", key, mid); err != nil {
		log.Error("SADD conn.Do error(%v)", err)
		return
	}
	return
}

// DelMonitorCache del monitor cache
func (d *Dao) DelMonitorCache(c context.Context, mid int64) (err error) {
	var (
		key  = monitorKey()
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = redis.Int64(conn.Do("SREM", key, mid)); err != nil {
		log.Error("SREM conn.Do(%s,%d) err(%v)", key, mid, err)
	}
	return
}

// LoadMonitorCache load monitor cache
func (d *Dao) LoadMonitorCache(c context.Context, mids []int64) (err error) {
	var (
		key  = monitorKey()
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, v := range mids {
		if err = conn.Send("SADD", key, v); err != nil {
			log.Error("SADD conn.Do error(%v)", err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	return
}

// TodayNotifyCountCache get notify count in the current day
func (d *Dao) TodayNotifyCountCache(c context.Context, mid int64) (notifyCount int64, err error) {
	var (
		key  = dailyNotifyCount(mid, gtime.Now())
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if notifyCount, err = redis.Int64(conn.Do("HGET", key, mid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("HGET conn.Do error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	return
}

// IncrTodayNotifyCount increment the today notify count in the current day
func (d *Dao) IncrTodayNotifyCount(c context.Context, mid int64) (err error) {
	var (
		key  = dailyNotifyCount(mid, gtime.Now())
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("HINCRBY", key, mid, 1); err != nil {
		log.Error("HINCRBY conn.Do error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, _notifyCountExpire); err != nil {
		log.Error("EXPIRE conn.Do error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	return
}
