package archive

import (
	"context"
	"strconv"

	"go-common/app/interface/main/app-view/model/view"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_prefixRelate         = "al_"
	_prefixViewStatic     = "avp_"
	_prefixStat           = "stp_"
	_prefixArchive        = "a3p_"
	_prefixViewContribute = "avpc_"
)

func keyRl(aid int64) string {
	return _prefixRelate + strconv.FormatInt(aid, 10)
}

func keyView(aid int64) string {
	return _prefixViewStatic + strconv.FormatInt(aid, 10)
}

func keyStat(aid int64) string {
	return _prefixStat + strconv.FormatInt(aid, 10)
}

func keyArc(aid int64) string {
	return _prefixArchive + strconv.FormatInt(aid, 10)
}

func keyViewContribute(mid int64) string {
	return _prefixViewContribute + strconv.FormatInt(mid, 10)
}

// AddRelatesCache add relates
func (d *Dao) AddRelatesCache(aid int64, rls []*view.Relate) {
	d.addCache(func() {
		d.addRelatesCache(context.TODO(), aid, rls)
	})
}

// RelatesCache get relates.
func (d *Dao) RelatesCache(c context.Context, aid int64) (rls []*view.Relate, err error) {
	conn := d.arcMc.Get(c)
	key := keyRl(aid)
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
	if err = conn.Scan(r, &rls); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%s)", r.Value)
	}
	return
}

// addRelatesCache add relates cache.
func (d *Dao) addRelatesCache(c context.Context, aid int64, rls []*view.Relate) (err error) {
	conn := d.arcMc.Get(c)
	key := keyRl(aid)
	item := &memcache.Item{Key: key, Object: rls, Flags: memcache.FlagJSON, Expiration: d.expireRlt}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%v,%d)", key, rls, d.expireRlt)
	}
	conn.Close()
	return
}

// viewCache get view cache from remote memecache .
func (d *Dao) viewCache(c context.Context, aid int64) (vs *archive.View3, err error) {
	conn := d.arcMc.Get(c)
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

// statCache get a archive stat from cache.
func (d *Dao) statCache(c context.Context, aid int64) (st *api.Stat, err error) {
	conn := d.arcMc.Get(c)
	key := keyStat(aid)
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

// statsCache get stat cache by aids
func (d *Dao) statsCache(c context.Context, aids []int64) (cached map[int64]*api.Stat, missed []int64, err error) {
	cached = make(map[int64]*api.Stat, len(aids))
	conn := d.arcMc.Get(c)
	defer conn.Close()
	keys := make([]string, 0, len(aids))
	for _, aid := range aids {
		keys = append(keys, keyStat(aid))
	}
	rs, err := conn.GetMulti(keys)
	if err != nil {
		err = errors.Wrapf(err, "conn.GetMulti(%v)", keys)
		return
	}
	for _, item := range rs {
		var st = &api.Stat{}
		if err = conn.Scan(item, st); err != nil {
			err = nil
			log.Error("conn.Scan(%v) error(%v)", item, err)
			continue
		}
		cached[st.Aid] = st
	}
	if len(cached) == len(aids) {
		return
	}
	for _, aid := range aids {
		if _, ok := cached[aid]; !ok {
			missed = append(missed, aid)
		}
	}
	return
}

// arcsCache get archives cache.
func (d *Dao) arcsCache(c context.Context, aids []int64) (cached map[int64]*api.Arc, missed []int64, err error) {
	cached = make(map[int64]*api.Arc, len(aids))
	conn := d.arcMc.Get(c)
	defer conn.Close()
	keys := make([]string, 0, len(aids))
	aidmap := make(map[string]int64, len(aids))
	for _, aid := range aids {
		k := keyArc(aid)
		if _, ok := aidmap[k]; !ok {
			keys = append(keys, k)
			aidmap[k] = aid
		}
	}
	rs, err := conn.GetMulti(keys)
	if err != nil {
		err = errors.Wrapf(err, "conn.GetMulti(%v)", keys)
		return
	}
	for k, r := range rs {
		a := &api.Arc{}
		if err = conn.Scan(r, a); err != nil {
			log.Error("conn.Scan(%s) error(%v)", r.Value, err)
			err = nil
			continue
		}
		cached[aidmap[k]] = a
		// delete hit key
		delete(aidmap, k)
	}
	// missed key
	missed = make([]int64, 0, len(aidmap))
	for _, aid := range aidmap {
		missed = append(missed, aid)
	}
	return
}

// arcCache get archive cache.
func (d *Dao) arcCache(c context.Context, aid int64) (a *api.Arc, err error) {
	conn := d.arcMc.Get(c)
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
	a = &api.Arc{}
	if err = conn.Scan(r, a); err != nil {
		a = nil
		err = errors.Wrapf(err, "conn.Scan(%s)", r.Value)
	}
	return
}

// ViewContributeCache get archive cache.
func (d *Dao) ViewContributeCache(c context.Context, mid int64) (aids []int64, err error) {
	conn := d.arcMc.Get(c)
	key := keyViewContribute(mid)
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
	if err = conn.Scan(r, &aids); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%s)", r.Value)
	}
	return
}

func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.arcMc.Get(c)
	err = conn.Set(&memcache.Item{Key: "ping", Object: []byte{1}, Flags: memcache.FlagJSON, Expiration: 0})
	conn.Close()
	return
}
