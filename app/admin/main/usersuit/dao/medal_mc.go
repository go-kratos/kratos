package dao

import (
	"context"
	"strconv"

	"go-common/app/admin/main/usersuit/model"
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
		log.Error("conn.Store(set, ping, 1) error(%v)", err)
	}
	return
}

// MedalOwnersCache get medal_owner cache.
func (d *Dao) MedalOwnersCache(c context.Context, mid int64) (res []*model.MedalOwner, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(ownersKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			res = nil
			err = nil
			return
		}
		log.Error("d.MedalOwnersCache err(%v)", err)
		return
	}
	res = make([]*model.MedalOwner, 0)
	if err = conn.Scan(item, &res); err != nil {
		log.Error("d.MedalOwnersCache err(%v)", err)
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
		log.Error("SetMedalOwnersache err(%v)", err)
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
			log.Error("d.DelMedalOwnersCache(%s) error(%v)", key, err)
		}
	}
	return
}

// MedalActivatedCache get user activated medal nid.
func (d *Dao) MedalActivatedCache(c context.Context, mid int64) (nid int64, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(activatedNidKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			nid = 0
			err = nil
			return
		}
		log.Error("d.MedalActivatedCache(mid:%d) err(%v)", mid, err)
		return
	}
	if err = conn.Scan(item, &nid); err != nil {
		log.Error("d.MedalActivatedCache(mid:%d) err(%v)", mid, err)
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
		log.Error("SetMedalActivatedCache err(%v)", err)
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
			log.Error("d.DelMedalActivatedCache(%s) error(%v)", key, err)
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
		log.Error("d.PopupCache(mid:%d) err(%v)", mid, err)
		return
	}
	if err = conn.Scan(item, &nid); err != nil {
		log.Error("d.PopupCache(mid:%d) err(%v)", mid, err)
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
		log.Error("SetMedalOwnersache err(%v)", err)
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
			log.Error("d.DelPopupCache(%s) error(%v)", key, err)
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
		log.Error("d.RedPointCache(mid:%d) err(%v)", mid, err)
		return
	}
	if err = conn.Scan(item, &nid); err != nil {
		log.Error("d.RedPointCache(mid:%d) err(%v)", mid, err)
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
		log.Error("SetRedPointCache(%d %d) err(%v)", mid, nid, err)
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
			log.Error("d.DelRedPointCache(%d) error(%v)", mid, err)
		}
	}
	return
}
