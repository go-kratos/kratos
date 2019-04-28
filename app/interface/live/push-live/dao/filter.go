package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/live/push-live/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_intervalUserkey   = "i:%d"     // 用户推送间隔缓存
	_limitUserDailyKey = "daily:%d" //用户每日推送额度缓存
	_defaultPushLimit  = 4          // 用户每日默认最大推送额度
)

// FilterConfig FilterConfig
type FilterConfig struct {
	Business        int
	IntervalExpired int32
	IntervalValue   string
	DailyExpired    float64
	Task            *model.ApPushTask
}

// Filter Filter
type Filter struct {
	conf *FilterConfig
	conn redis.Conn
}

// FilterChain FilterChain
type FilterChain map[string]func(ctx context.Context, mid int64) (bool, error)

// NewFilter NewFilter
func (d *Dao) NewFilter(conf *FilterConfig) (f *Filter, err error) {
	var conn redis.Conn
	// redis conn
	conn, err = redis.Dial(d.c.Redis.PushInterval.Proto, d.c.Redis.PushInterval.Addr, d.RedisOption()...)
	if err != nil {
		log.Error("[dao.filter|NewFilter] redis.Dial error(%v), conf(%v)", err, conf)
		return
	}
	f = &Filter{
		conf: conf,
		conn: conn,
	}
	return
}

// NewFilterChain NewFilterChain
func (d *Dao) NewFilterChain(f *Filter) FilterChain {
	funcs := make(FilterChain)
	if d.needLimit(f.conf.Business) {
		funcs["limit"] = f.dailyLimitFilter
	}
	if d.needSmooth(f.conf.Business) && f.conf.IntervalExpired > 0 {
		if f.conf.Business == model.ActivityBusiness {
			funcs["smooth"] = f.appointSmoothFilter
		} else {
			funcs["smooth"] = f.intervalSmoothFilter
		}
	}
	return funcs
}

// needSmooth
func (d *Dao) needSmooth(business int) bool {
	return !ignoreFilter(business, d.c.Push.PushFilterIgnores.Smooth)
}

// needLimit
func (d *Dao) needLimit(business int) bool {
	return !ignoreFilter(business, d.c.Push.PushFilterIgnores.Limit)
}

// Done do some close work
func (f *Filter) Done() {
	if f.conn != nil {
		f.conn.Close()
	}
}

// dailyLimitFilter 判断是否到达每日推送上限
func (f *Filter) dailyLimitFilter(ctx context.Context, mid int64) (b bool, err error) {
	var (
		left int
		key  = fmt.Sprintf(_limitUserDailyKey, mid)
	)

	// fetch daily push count
	left, err = redis.Int(f.conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			// key not exists, should return false & nil, first push today
			err = nil
			return
		}
		// actually error occurs
		return
	}

	// daily push limit
	if left <= 0 {
		b = true
		return
	}
	return
}

// intervalSmoothFilter 判断是否被平滑推送逻辑过滤
func (f *Filter) intervalSmoothFilter(ctx context.Context, mid int64) (b bool, err error) {
	var reply interface{}
	key := fmt.Sprintf(_intervalUserkey, mid)
	reply, err = f.conn.Do("SET", key, f.conf.IntervalValue, "EX", f.conf.IntervalExpired, "NX")
	if err != nil {
		return
	}
	// key exists, nil returned
	// key not exists, will return OK
	if reply == nil {
		b = true
		return
	}
	return
}

// appointSmoothFilter 预约逻辑的平滑
func (f *Filter) appointSmoothFilter(ctx context.Context, mid int64) (b bool, err error) {
	var reply interface{}
	key := fmt.Sprintf(_intervalUserkey, mid)
	reply, err = f.conn.Do("SET", key, f.conf.IntervalValue, "EX", f.conf.IntervalExpired, "NX")
	if err != nil {
		return
	}
	// key exists, nil returned
	// key not exists, will return OK
	if reply == nil {
		// 活动预约有特殊判断逻辑
		reply, err = redis.String(f.conn.Do("GET", key))
		if err != nil {
			return
		}
		// 相同房间会被过滤
		if reply == f.conf.IntervalValue {
			b = true
		}
		return
	}
	return
}

// BatchFilter 对输入mid序列执行所有过滤方法，返回过滤结果
func (f *Filter) BatchFilter(ctx context.Context, filterChain FilterChain, mids []int64) (resMids []int64) {
	if len(mids) == 0 || len(filterChain) == 0 {
		return
	}
	var (
		isFiltered bool
		err        error
		filterMids = make(map[string][]int64)
		errMids    = make(map[string][]int64)
	)
	defer func() {
		filterMids = nil
		errMids = nil
	}()
	resMids = make([]int64, 0, len(mids))

	// 记录被过滤掉的mid，过滤发生错误的mid
	for name := range filterChain {
		filterMids[name] = make([]int64, 0, len(mids))
		errMids[name] = make([]int64, 0, len(mids))
	}

MidLoop:
	for _, mid := range mids {
		for name, fc := range filterChain {
			isFiltered, err = fc(ctx, mid)
			// error occurs, next mid
			if err != nil {
				errMids[name] = append(errMids[name], mid)
				continue MidLoop
			}
			// filtered by any filterChain func, next mid
			if isFiltered {
				filterMids[name] = append(filterMids[name], mid)
				continue MidLoop
			}
		}
		// mid here is filter result, should push
		resMids = append(resMids, mid)
	}

	// log
	for name, ids := range filterMids {
		if len(ids) == 0 {
			continue
		}
		log.Info("[dao.filter|BatchFilter] BatchFilter filterMids, task(%v), len(%d), name(%s), mids(%d)",
			f.conf.Task, len(ids), name, len(mids))
	}
	for name, ids := range errMids {
		if len(ids) == 0 {
			continue
		}
		log.Error("[dao.filter|BatchFilter] BatchFilter errMids, task(%v), len(%d), name(%s), err(%v)",
			f.conf.Task, len(ids), name, err)
	}
	return
}

// BatchDecreaseLimit 批量减少配额
func (f *Filter) BatchDecreaseLimit(ctx context.Context, mids []int64) (total int, err error) {
	defer func() {
		if f != nil {
			f.Done()
		}
		log.Info("[dao.filter|BatchDecreaseLimit] business(%d), input(%d), exec(%d), err(%v)",
			f.conf.Business, len(mids), total, err)
	}()
	if len(mids) == 0 {
		return
	}
	initLeft := _defaultPushLimit - 1
	for _, mid := range mids {
		key := fmt.Sprintf(_limitUserDailyKey, mid)
		left, err := redis.Int(f.conn.Do("GET", key))
		if err != nil {
			if err == redis.ErrNil {
				// key not exists
				f.conn.Do("SET", key, initLeft, "EX", f.conf.DailyExpired)
				total++
			}
			continue
		}
		f.conn.Do("SET", key, left-1, "EX", f.conf.DailyExpired)
		total++
	}
	return
}

// ignoreFilter 判断business是否能够不需要过滤
func ignoreFilter(business int, ignores []int) bool {
	var f = false
	for _, ignore := range ignores {
		if business == ignore {
			f = true
		}
	}
	return f
}

// GetIntervalKey return interval smooth redis key
func GetIntervalKey(mid int64) string {
	return fmt.Sprintf(_intervalUserkey, mid)
}

// GetDailyLimitKey return daily limit redis key
func GetDailyLimitKey(mid int64) string {
	return fmt.Sprintf(_limitUserDailyKey, mid)
}
