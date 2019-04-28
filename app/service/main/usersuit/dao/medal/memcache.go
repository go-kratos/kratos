package medal

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"go-common/app/service/main/usersuit/model"
	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixActivatedNid = "usma:" // key of activated medal nid
	_prefixOwners       = "umos:" // key of owners info
	_prefixRedPoint     = "usrp:" // key of red point
	_prefixPopup        = "uspp:" // key of new medal popup
)

// medalactivated medal nid key.
func activatedNidKey(mid int64) string {
	return _prefixActivatedNid + strconv.FormatInt(mid, 10)
}

// ownersKey medal_owner key.
func ownersKey(mid int64) string {
	return _prefixOwners + strconv.FormatInt(mid, 10)
}

//RedPointKey new medal RedPoint key.
func RedPointKey(mid int64) string {
	return _prefixRedPoint + strconv.FormatInt(mid, 10)
}

// PopupKey new medal popup key.
func PopupKey(mid int64) string {
	return _prefixPopup + strconv.FormatInt(mid, 10)
}

func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&gmc.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// MedalOwnersCache get medal_owner cache.
func (d *Dao) MedalOwnersCache(c context.Context, mid int64) (res []*model.MedalOwner, notFound bool, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(ownersKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			res = nil
			err = nil
			notFound = true
			return
		}
		err = errors.WithStack(err)
		return
	}
	res = make([]*model.MedalOwner, 0)
	if err = conn.Scan(item, &res); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// SetMedalOwnersache set medal_owner cache.
func (d *Dao) SetMedalOwnersache(c context.Context, mid int64, nos []*model.MedalOwner) (err error) {
	key := ownersKey(mid)
	item := &gmc.Item{Key: key, Object: nos, Expiration: d.mcExpire, Flags: gmc.FlagJSON}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DelMedalOwnersCache delete medal_owner cache.
func (d *Dao) DelMedalOwnersCache(c context.Context, mid int64) (err error) {
	key := ownersKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// MedalsActivatedCache multi get user activated medal nid from memcache.
func (d *Dao) medalsActivatedCache(c context.Context, mids []int64) (nids map[int64]int64, missed []int64, err error) {
	nids = make(map[int64]int64, len(mids))
	keys := make([]string, len(mids))
	mm := make(map[string]int64, len(mids))
	for i, mid := range mids {
		var key = activatedNidKey(mid)
		keys[i] = key
		mm[key] = mid
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	items, err := conn.GetMulti(keys)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	for _, item := range items {
		var nid int64
		if err = conn.Scan(item, &nid); err != nil {
			log.Error("conn.Scan(%s) error(%v)", item.Value, err)
			continue
		}
		nids[mm[item.Key]] = nid
		delete(mm, item.Key)
	}
	missed = make([]int64, 0, len(mm))
	for _, m := range mm {
		missed = append(missed, m)
	}
	return
}

// MedalActivatedCache get user activated medal nid.
func (d *Dao) medalActivatedCache(c context.Context, mid int64) (nid int64, notFound bool, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(activatedNidKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			nid = 0
			err = nil
			notFound = true
			return
		}
		err = errors.WithStack(err)
		return
	}
	if err = conn.Scan(item, &nid); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// SetMedalActivatedCache set activated medal  cache.
func (d *Dao) SetMedalActivatedCache(c context.Context, mid, nid int64) (err error) {
	key := activatedNidKey(mid)
	item := &gmc.Item{Key: key, Object: nid, Expiration: d.mcExpire, Flags: gmc.FlagJSON}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DelMedalActivatedCache delete activated medal cache.
func (d *Dao) DelMedalActivatedCache(c context.Context, mid int64) (err error) {
	key := activatedNidKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// PopupCache get new medal info popup cache.
func (d *Dao) PopupCache(c context.Context, mid int64) (nid int64, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(PopupKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			nid = 0
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	if err = conn.Scan(item, &nid); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// SetPopupCache set popup cache.
func (d *Dao) SetPopupCache(c context.Context, mid, nid int64) (err error) {
	key := PopupKey(mid)
	item := &gmc.Item{Key: key, Object: nid, Expiration: d.pointExpire, Flags: gmc.FlagJSON}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DelPopupCache delete new medal info popup cache.
func (d *Dao) DelPopupCache(c context.Context, mid int64) (err error) {
	key := PopupKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// RedPointCache get new medal info red point cache.
func (d *Dao) RedPointCache(c context.Context, mid int64) (nid int64, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(RedPointKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	if err = conn.Scan(item, &nid); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// SetRedPointCache set red point cache.
func (d *Dao) SetRedPointCache(c context.Context, mid, nid int64) (err error) {
	key := RedPointKey(mid)
	item := &gmc.Item{Key: key, Object: nid, Expiration: d.pointExpire, Flags: gmc.FlagJSON}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DelRedPointCache delete new medal info red point cache.
func (d *Dao) DelRedPointCache(c context.Context, mid int64) (err error) {
	key := RedPointKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}
