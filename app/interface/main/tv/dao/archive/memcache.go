package archive

import (
	"context"
	"strconv"

	"go-common/app/interface/main/tv/model/view"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_prefixRelate     = "al_"
	_prefixViewStatic = "avp_"
	_prefixArchive    = "a3p_"
)

func keyRl(aid int64) string {
	return _prefixRelate + strconv.FormatInt(aid, 10)
}

func keyView(aid int64) string {
	return _prefixViewStatic + strconv.FormatInt(aid, 10)
}

func keyArc(aid int64) string {
	return _prefixArchive + strconv.FormatInt(aid, 10)
}

// AddArcCache add arc cache
func (d *Dao) AddArcCache(aid int64, arc *arcwar.Arc) {
	d.addCache(func() {
		d.addArcCache(context.TODO(), aid, arc)
	})
}

// AddRelatesCache add relates
func (d *Dao) AddRelatesCache(aid int64, rls []*view.Relate) {
	d.addCache(func() {
		d.addRelatesCache(context.TODO(), aid, rls)
	})
}

// AddViewCache add view relates
func (d *Dao) AddViewCache(aid int64, vp *arcwar.ViewReply) {
	d.addCache(func() {
		d.addViewCache(context.TODO(), aid, vp)
	})
}

// addViewCache add relates cache.
func (d *Dao) addViewCache(c context.Context, aid int64, vp *arcwar.ViewReply) (err error) {
	conn := d.arcMc.Get(c)
	key := keyView(aid)
	item := &memcache.Item{Key: key, Object: vp, Flags: memcache.FlagJSON, Expiration: d.expireView}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%v,%d)", key, vp, d.expireView)
	}
	conn.Close()
	return
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

// addRelatesCache add relates cache.
func (d *Dao) addArcCache(c context.Context, aid int64, cached *arcwar.Arc) (err error) {
	conn := d.arcMc.Get(c)
	key := keyArc(aid)
	item := &memcache.Item{Key: key, Object: cached, Flags: memcache.FlagJSON, Expiration: d.expireArc}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%v,%d)", key, cached, d.expireArc)
	}
	conn.Close()
	return
}

// arcsCache get archives cache.
func (d *Dao) arcsCache(c context.Context, aids []int64) (cached map[int64]*arcwar.Arc, missed []int64, err error) {
	cached = make(map[int64]*arcwar.Arc, len(aids))
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
		a := &arcwar.Arc{}
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
