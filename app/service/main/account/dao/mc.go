package dao

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	v1 "go-common/app/service/main/account/api"
	mc "go-common/library/cache/memcache"
)

const (
	_prefixInfo    = "i3_"
	_prefixCard    = "c3_"
	_prefixVip     = "v3_"
	_prefixProfile = "p3_"
)

func keyInfo(mid int64) string {
	return _prefixInfo + strconv.FormatInt(mid, 10)
}

func keyCard(mid int64) string {
	return _prefixCard + strconv.FormatInt(mid, 10)
}

func keyVip(mid int64) string {
	return _prefixVip + strconv.FormatInt(mid, 10)
}

func keyProfile(mid int64) string {
	return _prefixProfile + strconv.FormatInt(mid, 10)
}

// CacheInfo get account info from cache.
func (d *Dao) CacheInfo(c context.Context, mid int64) (v *v1.Info, err error) {
	key := keyInfo(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao cache info")
		return
	}
	v = &v1.Info{}
	if err = conn.Scan(r, v); err != nil {
		err = errors.Wrap(err, "dao cache scan info")
	}
	return
}

// AddCacheInfo set account info into cache.
func (d *Dao) AddCacheInfo(c context.Context, mid int64, v *v1.Info) (err error) {
	item := &mc.Item{
		Key:        keyInfo(mid),
		Object:     v,
		Flags:      mc.FlagProtobuf,
		Expiration: d.mcExpire,
	}
	conn := d.mc.Get(c)
	err = conn.Set(item)
	conn.Close()
	if err != nil {
		err = errors.Wrap(err, "dao add info cache")
	}
	return
}

// CacheInfos multi get account info from cache.
func (d *Dao) CacheInfos(c context.Context, mids []int64) (res map[int64]*v1.Info, err error) {
	keys := make([]string, 0, len(mids))
	keyMidMap := make(map[string]int64, len(mids))
	for _, mid := range mids {
		key := keyInfo(mid)
		if _, ok := keyMidMap[key]; !ok {
			// duplicate mid
			keyMidMap[key] = mid
			keys = append(keys, key)
		}
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	rs, err := conn.GetMulti(keys)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao infos cache")
		return
	}
	res = make(map[int64]*v1.Info, len(mids))
	for _, r := range rs {
		ai := &v1.Info{}
		conn.Scan(r, ai)
		res[ai.Mid] = ai
	}
	return
}

// AddCacheInfos set account infos cache.
func (d *Dao) AddCacheInfos(c context.Context, im map[int64]*v1.Info) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, i := range im {
		item := &mc.Item{
			Key:        keyInfo(i.Mid),
			Object:     i,
			Flags:      mc.FlagProtobuf,
			Expiration: d.mcExpire,
		}
		err = conn.Set(item)
		if err != nil {
			err = errors.Wrap(err, "dao add infos cache")
		}
	}
	return
}

// CacheCard get account card from cache.
func (d *Dao) CacheCard(c context.Context, mid int64) (v *v1.Card, err error) {
	key := keyCard(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao cache card")
		return
	}
	v = &v1.Card{}
	if err = conn.Scan(r, v); err != nil {
		err = errors.Wrap(err, "dao cache scan card")
	}
	return
}

// AddCacheCard set account card into cache.
func (d *Dao) AddCacheCard(c context.Context, mid int64, v *v1.Card) (err error) {
	item := &mc.Item{
		Key:        keyCard(mid),
		Object:     v,
		Flags:      mc.FlagProtobuf,
		Expiration: d.mcExpire,
	}
	conn := d.mc.Get(c)
	err = conn.Set(item)
	conn.Close()
	if err != nil {
		err = errors.Wrap(err, "dao add card cache")
	}
	return
}

// CacheCards multi get account cards from cache.
func (d *Dao) CacheCards(c context.Context, mids []int64) (res map[int64]*v1.Card, err error) {
	keys := make([]string, 0, len(mids))
	keyMidMap := make(map[string]int64, len(mids))
	for _, mid := range mids {
		key := keyCard(mid)
		if _, ok := keyMidMap[key]; !ok {
			// duplicate mid
			keyMidMap[key] = mid
			keys = append(keys, key)
		}
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	rs, err := conn.GetMulti(keys)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao cards cache")
		return
	}
	res = make(map[int64]*v1.Card, len(mids))
	for _, r := range rs {
		ai := &v1.Card{}
		conn.Scan(r, ai)
		res[ai.Mid] = ai
	}
	return
}

// AddCacheCards set account cards cache.
func (d *Dao) AddCacheCards(c context.Context, cm map[int64]*v1.Card) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, card := range cm {
		item := &mc.Item{
			Key:        keyCard(card.Mid),
			Object:     card,
			Flags:      mc.FlagProtobuf,
			Expiration: d.mcExpire,
		}
		err = conn.Set(item)
		if err != nil {
			err = errors.Wrap(err, "dao add cards cache")
		}
	}
	return
}

// CacheVip get vip cache.
func (d *Dao) CacheVip(c context.Context, mid int64) (v *v1.VipInfo, err error) {
	key := keyVip(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao vip cache")
		return
	}
	v = new(v1.VipInfo)
	if err = conn.Scan(r, v); err != nil {
		err = errors.Wrap(err, "dao vip cache scan")
	}
	return
}

// AddCacheVip set vip cache.
func (d *Dao) AddCacheVip(c context.Context, mid int64, v *v1.VipInfo) (err error) {
	conn := d.mc.Get(c)
	conn.Set(&mc.Item{
		Key:        keyVip(mid),
		Object:     v,
		Flags:      mc.FlagProtobuf,
		Expiration: d.mcExpire,
	})
	conn.Close()
	if err != nil {
		err = errors.Wrap(err, "dao vip add cache")
	}
	return
}

// CacheVips multi get account cards from cache.
func (d *Dao) CacheVips(c context.Context, mids []int64) (res map[int64]*v1.VipInfo, err error) {
	keys := make([]string, 0, len(mids))
	keyMidMap := make(map[string]int64, len(mids))
	for _, mid := range mids {
		key := keyVip(mid)
		if _, ok := keyMidMap[key]; !ok {
			// duplicate mid
			keyMidMap[key] = mid
			keys = append(keys, key)
		}
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	rs, err := conn.GetMulti(keys)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao vips cache")
		return
	}
	res = make(map[int64]*v1.VipInfo, len(mids))
	for _, r := range rs {
		ai := &v1.VipInfo{}
		conn.Scan(r, ai)
		res[keyMidMap[r.Key]] = ai
	}
	return
}

// AddCacheVips set account vips cache.
func (d *Dao) AddCacheVips(c context.Context, vm map[int64]*v1.VipInfo) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for mid, v := range vm {
		item := &mc.Item{
			Key:        keyVip(mid),
			Object:     v,
			Flags:      mc.FlagProtobuf,
			Expiration: d.mcExpire,
		}
		err = conn.Set(item)
		if err != nil {
			err = errors.Wrap(err, "dao add vips cache")
		}
	}
	return
}

// CacheProfile get profile cache.
func (d *Dao) CacheProfile(c context.Context, mid int64) (v *v1.Profile, err error) {
	key := keyProfile(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao profile cache")
		return
	}
	v = new(v1.Profile)
	if err = conn.Scan(r, v); err != nil {
		err = errors.Wrap(err, "dao profile cache scan")
	}
	return
}

// AddCacheProfile set profile cache.
func (d *Dao) AddCacheProfile(c context.Context, mid int64, v *v1.Profile) (err error) {
	conn := d.mc.Get(c)
	conn.Set(&mc.Item{
		Key:        keyProfile(mid),
		Object:     v,
		Flags:      mc.FlagProtobuf,
		Expiration: d.mcExpire,
	})
	conn.Close()
	if err != nil {
		err = errors.Wrap(err, "dao profile add cache")
	}
	return
}

// DelCache delete cache.
func (d *Dao) DelCache(c context.Context, mid int64) []error {
	conn := d.mc.Get(c)
	errs := make([]error, 0, 5)
	if err := conn.Delete(keyInfo(mid)); err != nil {
		errs = append(errs, errors.Wrap(err, keyInfo(mid)))
	}
	if err := conn.Delete(keyCard(mid)); err != nil {
		errs = append(errs, errors.Wrap(err, keyCard(mid)))
	}
	if err := conn.Delete(keyVip(mid)); err != nil {
		errs = append(errs, errors.Wrap(err, keyVip(mid)))
	}
	if err := conn.Delete(keyProfile(mid)); err != nil {
		errs = append(errs, errors.Wrap(err, keyProfile(mid)))
	}
	if err := conn.Close(); err != nil {
		errs = append(errs, errors.Wrap(err, "conn close"))
	}
	d.cache.Do(c, func(ctx context.Context) {
		d.Info(ctx, mid)
		d.Card(ctx, mid)
		d.Vip(ctx, mid)
		d.Profile(ctx, mid)
	})
	return errs
}
