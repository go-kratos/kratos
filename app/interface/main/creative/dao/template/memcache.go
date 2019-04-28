package template

import (
	"context"
	"go-common/app/interface/main/creative/model/template"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"strconv"
)

const (
	_prefix = "tpl_"
)

func keyTpl(mid int64) string {
	return _prefix + strconv.FormatInt(mid, 10)
}

// tplCache get tpl cache.
func (d *Dao) tplCache(c context.Context, mid int64) (tps []*template.Template, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	// get cache
	r, err = conn.Get(keyTpl(mid))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache conn.Get2(%d) error(%v)", mid, err)
		}
		return
	}
	tps = []*template.Template{}
	if err = conn.Scan(r, &tps); err != nil {
		log.Error("tplCache json.Unmarshal(%s) error(%v)", r.Value, err)
		tps = nil
	}
	return
}

// addTplCache add tpl cache.
func (d *Dao) addTplCache(c context.Context, mid int64, tps []*template.Template) (err error) {
	var (
		key = keyTpl(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: tps, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("tplCache memcache.Set(%v) error(%v)", key, err)
	}

	return
}

// delTplCache del tpl cache.
func (d *Dao) delTplCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	// del cache
	if err = conn.Delete(keyTpl(mid)); err == memcache.ErrNotFound {
		err = nil
	}
	return
}
