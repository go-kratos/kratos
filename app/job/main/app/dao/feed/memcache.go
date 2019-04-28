package feed

import (
	"context"
	"strconv"

	"go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_prefixArchivePB = "a3p_"
)

func keyArcPB(aid int64) string {
	return _prefixArchivePB + strconv.FormatInt(aid, 10)
}

// arcsCache get archives cache.
func (d *Dao) arcsCache(c context.Context, aids []int64) (cached map[int64]*api.Arc, missed []int64, err error) {
	var (
		keys   = make([]string, 0, len(aids))
		conn   = d.mc.Get(c)
		aidmap = make(map[string]int64, len(aids))
		rs     map[string]*memcache.Item
		a      *api.Arc
	)
	cached = make(map[int64]*api.Arc, len(aids))
	defer conn.Close()
	for _, aid := range aids {
		k := keyArcPB(aid)
		if _, ok := aidmap[k]; !ok {
			keys = append(keys, k)
			aidmap[k] = aid
		}
	}
	if rs, err = conn.GetMulti(keys); err != nil {
		err = errors.Wrapf(err, "conn.GetMulti(%v)", keys)
		return
	}
	for k, r := range rs {
		a = &api.Arc{}
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
