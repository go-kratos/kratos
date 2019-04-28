package view

import (
	"context"
	"strconv"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

const (
	_prefixArc            = "a3p_"
	_prefixView           = "avp_"
	_prefixSt             = "stp_"
	_prefixViewContribute = "avpc_"
)

func keyArc(aid int64) string {
	return _prefixArc + strconv.FormatInt(aid, 10)
}

func keyView(aid int64) string {
	return _prefixView + strconv.FormatInt(aid, 10)
}

func keySt(aid int64) string {
	return _prefixSt + strconv.FormatInt(aid, 10)
}

func keyViewContribute(mid int64) string {
	return _prefixViewContribute + strconv.FormatInt(mid, 10)
}

// UpArcCache update archive cache
func (d *Dao) UpArcCache(c context.Context, a *archive.Archive3) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{Key: keyArc(a.Aid), Object: a, Flags: memcache.FlagProtobuf, Expiration: 0}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%v)", item)
	}
	conn.Close()
	return
}

// DelArcCache delete archive cache
func (d *Dao) DelArcCache(c context.Context, aid int64) (err error) {
	conn := d.mc.Get(c)
	key := keyArc(aid)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			err = errors.Wrapf(err, "conn.Delete(%s)", key)
		}
	}
	conn.Close()
	return
}

// UpViewCache up all app cache .
func (d *Dao) UpViewCache(c context.Context, v *archive.View3) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{Key: keyView(v.Aid), Object: v, Flags: memcache.FlagProtobuf, Expiration: 0}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%v)", item)
	}
	conn.Close()
	return
}

// DelViewCache del view cache
func (d *Dao) DelViewCache(c context.Context, aid int64) (err error) {
	conn := d.mc.Get(c)
	key := keyView(aid)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			err = errors.Wrapf(err, "conn.Delete(%s)", key)
		}
	}
	conn.Close()
	return
}

// StatCache get a archive stat from cache.
func (d *Dao) StatCache(c context.Context, aid int64) (st *api.Stat, err error) {
	conn := d.mc.Get(c)
	key := keySt(aid)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	st = &api.Stat{}
	if err = conn.Scan(r, st); err != nil {
		st = nil
		err = errors.Wrapf(err, "conn.Scan(%s)", r.Value)
	}
	return
}

// UpStatCache up st cache
func (d *Dao) UpStatCache(c context.Context, st *api.Stat) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{Key: keySt(st.Aid), Object: st, Flags: memcache.FlagProtobuf, Expiration: 0}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%v)", item)
	}
	conn.Close()
	return
}

// UpViewContributeCache up app view contribute cache .
func (d *Dao) UpViewContributeCache(c context.Context, mid int64, aids []int64) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{Key: keyViewContribute(mid), Object: aids, Flags: memcache.FlagJSON, Expiration: 0}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%v)", item)
	}
	conn.Close()
	return
}

// ViewCache get view cache from remote memecache .
func (d *Dao) ViewCache(c context.Context, aid int64) (vs *archive.View3, err error) {
	conn := d.mc.Get(c)
	key := keyView(aid)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	vs = &archive.View3{Archive3: &archive.Archive3{}}
	if err = conn.Scan(r, vs); err != nil {
		vs = nil
		err = errors.Wrapf(err, "conn.Scan(%s)", r.Value)
	}
	return
}

// ArcCache get archive cache.
func (d *Dao) ArcCache(c context.Context, aid int64) (a *archive.Archive3, err error) {
	conn := d.mc.Get(c)
	key := keyArc(aid)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	a = &archive.Archive3{}
	if err = conn.Scan(r, a); err != nil {
		a = nil
		err = errors.Wrapf(err, "conn.Scan(%s)", r.Value)
	}
	return
}
