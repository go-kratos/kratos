package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/service/main/member/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_expPrefix   = "exp_%d"
	_moralPrefix = "moral_%d"
	_expExpire   = 86400
)

func expKey(mid int64) string {
	return fmt.Sprintf(_expPrefix, mid)
}

func moralKey(mid int64) string {
	return fmt.Sprintf(_moralPrefix, mid)
}

// pingMC ping memcache.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        "ping",
		Value:      []byte{1},
		Expiration: 86400,
	}); err != nil {
		log.Error("conn.Set(ping, 1) error(%v)", err)
	}
	return
}

//   -------- base --------- //

// BaseInfoCache get base info from mc.
func (d *Dao) BaseInfoCache(c context.Context, mid int64) (info *model.BaseInfo, err error) {
	key := fmt.Sprintf(model.CacheKeyBase, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	info = &model.BaseInfo{}
	if err = conn.Scan(item, info); err != nil {
		log.Error("conn.Scan(%s) error(%v)", string(item.Value), err)
	}
	return
}

// BatchBaseInfoCache get batch base info from mc.
func (d *Dao) BatchBaseInfoCache(c context.Context, mids []int64) (cached map[int64]*model.BaseInfo, missed []int64, err error) {
	cached = make(map[int64]*model.BaseInfo, len(mids))
	if len(mids) == 0 {
		return
	}
	keys := make([]string, 0, len(mids))
	midmap := make(map[string]int64, len(mids))
	for _, mid := range mids {
		k := fmt.Sprintf(model.CacheKeyBase, mid)
		keys = append(keys, k)
		midmap[k] = mid
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	bases, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("conn.Gets(%v) error(%v)", keys, err)
		return
	}
	for _, base := range bases {
		b := &model.BaseInfo{}
		if err = conn.Scan(base, b); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", base.Value, err)
			return
		}
		cached[midmap[base.Key]] = b
		delete(midmap, base.Key)
	}
	missed = make([]int64, 0, len(midmap))
	for _, bid := range midmap {
		missed = append(missed, bid)
	}
	return

}

// SetBatchBaseInfoCache set batch base info to mc.
func (d *Dao) SetBatchBaseInfoCache(c context.Context, bs []*model.BaseInfo) (err error) {
	for _, info := range bs {
		d.SetBaseInfoCache(c, info.Mid, info)
	}
	return
}

// SetBaseInfoCache set base info to mc
func (d *Dao) SetBaseInfoCache(c context.Context, mid int64, info *model.BaseInfo) (err error) {
	key := fmt.Sprintf(model.CacheKeyBase, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        key,
		Object:     info,
		Flags:      memcache.FlagProtobuf,
		Expiration: d.baseTTL,
	}); err != nil {
		log.Error("conn.Set(%s, %v) error(%v)", key, info, err)
	}
	return
}

// DelBaseInfoCache delete baseInfo cache.
func (d *Dao) DelBaseInfoCache(c context.Context, mid int64) (err error) {
	key := fmt.Sprintf(model.CacheKeyBase, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%s) error(%v)", key, err)
	}
	return
}

//   -------- base --------- //

// ----- exp ------ //

// Exp get user exp from cache,if miss get from db.
func (d *Dao) expCache(c context.Context, mid int64) (exp int64, err error) {
	key := expKey(mid)
	conn := d.mc.Get(c)
	res, err := conn.Get(key)
	defer conn.Close()
	if err != nil {
		return
	}
	exp, _ = strconv.ParseInt(string(res.Value), 10, 64)
	return
}

// ExpsCache get users exp cache.
func (d *Dao) expsCache(c context.Context, mids []int64) (exps map[int64]int64, miss []int64, err error) {
	var keys []string
	for _, mid := range mids {
		keys = append(keys, expKey(mid))
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	its, err := conn.GetMulti(keys)
	if err != nil {
		return
	}
	exps = make(map[int64]int64)
	for _, mid := range mids {
		if it, ok := its[expKey(mid)]; ok {
			exp, _ := strconv.ParseInt(string(it.Value), 10, 64)
			exps[mid] = exp
		} else {
			miss = append(miss, mid)
		}
	}
	return
}

// SetExpCache set user exp cache.
func (d *Dao) SetExpCache(c context.Context, mid, exp int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        expKey(mid),
		Value:      []byte(strconv.FormatInt(exp, 10)),
		Expiration: _expExpire,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// moralCache get moral from cache.
func (d *Dao) moralCache(c context.Context, mid int64) (moral *model.Moral, err error) {
	key := moralKey(mid)
	conn := d.mc.Get(c)
	item, err := conn.Get(key)
	defer conn.Close()
	if err != nil {
		return
	}
	moral = &model.Moral{}
	if err = conn.Scan(item, moral); err != nil {
		log.Error("conn.Scan(%s) error(%v)", string(item.Value), err)
	}
	return
}

// SetMoralCache set moral to mc
func (d *Dao) SetMoralCache(c context.Context, mid int64, moral *model.Moral) (err error) {
	key := moralKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if conn.Set(&memcache.Item{
		Key:        key,
		Object:     moral,
		Flags:      memcache.FlagProtobuf,
		Expiration: d.moralTTL,
	}); err != nil {
		log.Error("conn.Set(%s, %v) error(%v)", key, moral, err)
	}
	return
}

// DelMoralCache delete moral cache.
func (d *Dao) DelMoralCache(c context.Context, mid int64) (err error) {
	key := moralKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("DelMoralCache conn.Delete(%s) error(%v)", key, err)
	}
	return
}

// ------realname------
func realnameInfoKey(mid int64) string {
	return fmt.Sprintf("realname_info_%d", mid)
}

func realnameCaptureTimesKey(mid int64) string {
	return fmt.Sprintf("realname_cap_times_%d", mid)
}

func realnameCaptureCodeKey(mid int64) string {
	return fmt.Sprintf("realname_cap_code_%d", mid)
}

func realnameCaptureErrTimesKey(mid int64) string {
	return fmt.Sprintf("realname_cap_err_times%d", mid)
}

// RealnameCaptureTimesCache is
func (d *Dao) RealnameCaptureTimesCache(c context.Context, mid int64) (times int, err error) {
	var (
		key  = realnameCaptureTimesKey(mid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			times = -1
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	if err = conn.Scan(item, &times); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%+v)", item)
		return
	}
	return
}

// IncreaseRealnameCaptureTimes is
func (d *Dao) IncreaseRealnameCaptureTimes(c context.Context, mid int64) (err error) {
	var (
		key  = realnameCaptureTimesKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Increment(key, 1); err != nil {
		err = errors.Wrapf(err, "conn.Increment(%s,1)", key)
		return
	}
	return
}

// SetRealnameCaptureTimes is
func (d *Dao) SetRealnameCaptureTimes(c context.Context, mid int64, times int) (err error) {
	var (
		key  = realnameCaptureTimesKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: times, Flags: memcache.FlagJSON, Expiration: d.captureTimesTTL}); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%+v)", key, times)
		return
	}
	return
}

// RealnameCaptureCodeCache .
// return code : -1 , if code not found
// RealnameCaptureCodeCache is
func (d *Dao) RealnameCaptureCodeCache(c context.Context, mid int64) (code int, err error) {
	var (
		key  = realnameCaptureCodeKey(mid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			code = -1
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	if err = conn.Scan(item, &code); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%+v)", item)
		return
	}
	return
}

// RealnameInfoCache is.
func (d *Dao) RealnameInfoCache(c context.Context, mid int64) (info *model.RealnameCacheInfo, err error) {
	var (
		key  = realnameInfoKey(mid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			info = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	info = &model.RealnameCacheInfo{}
	if err = conn.Scan(item, &info); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%+v)", item)
		return
	}
	return
}

// SetRealnameInfo is.
func (d *Dao) SetRealnameInfo(c context.Context, mid int64, info *model.RealnameCacheInfo) (err error) {
	var (
		key  = realnameInfoKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: info, Flags: memcache.FlagJSON, Expiration: d.applyInfoTTL}); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%+v)", key, info)
		return
	}
	return
}

// SetRealnameCaptureCode is
func (d *Dao) SetRealnameCaptureCode(c context.Context, mid int64, code int) (err error) {
	var (
		key  = realnameCaptureCodeKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: code, Flags: memcache.FlagJSON, Expiration: d.captureCodeTTL}); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%+v)", key, code)
		return
	}
	return
}

// DeleteRealnameCaptureCode is
func (d *Dao) DeleteRealnameCaptureCode(c context.Context, mid int64) (err error) {
	var (
		key  = realnameCaptureCodeKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Delete(%s)", key)
		return
	}
	return
}

// DeleteRealnameInfo is
func (d *Dao) DeleteRealnameInfo(c context.Context, mid int64) (err error) {
	var (
		key  = realnameInfoKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Delete(%s)", key)
		return
	}
	return
}

// RealnameCaptureErrTimesCache is
func (d *Dao) RealnameCaptureErrTimesCache(c context.Context, mid int64) (times int, err error) {
	var (
		key  = realnameCaptureErrTimesKey(mid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			times = -1
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	if err = conn.Scan(item, &times); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%+v)", item)
		return
	}
	return
}

// SetRealnameCaptureErrTimes is
func (d *Dao) SetRealnameCaptureErrTimes(c context.Context, mid int64, times int) (err error) {
	var (
		key  = realnameCaptureErrTimesKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: times, Flags: memcache.FlagJSON, Expiration: d.captureErrTimesTTL}); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%+v)", key, times)
		return
	}
	return
}

// IncreaseRealnameCaptureErrTimes is
func (d *Dao) IncreaseRealnameCaptureErrTimes(c context.Context, mid int64) (err error) {
	var (
		key  = realnameCaptureErrTimesKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Increment(key, 1); err != nil {
		err = errors.Wrapf(err, "conn.Increment(%s,1)", key)
		return
	}
	return
}

// DeleteRealnameCaptureErrTimes is
func (d *Dao) DeleteRealnameCaptureErrTimes(c context.Context, mid int64) (err error) {
	var (
		key  = realnameCaptureErrTimesKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Delete(%s)", key)
		return
	}
	return
}
