package archive

import (
	"context"
	"strconv"

	"go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_prefixArchive = "a3p_"
	_prefixStat    = "stp_"
)

func keyArc(aid int64) string {
	return _prefixArchive + strconv.FormatInt(aid, 10)
}

func keyStat(aid int64) string {
	return _prefixStat + strconv.FormatInt(aid, 10)
}

// arcsCache get archives cache.
func (d *Dao) arcsCache(c context.Context, aids []int64) (cached map[int64]*api.Arc, missed []int64, err error) {
	cached = make(map[int64]*api.Arc, len(aids))
	var rs map[string]*memcache.Item
	conn := d.arcMc.Get(c)
	keys := make([]string, 0, len(aids))
	aidmap := make(map[string]int64, len(aids))
	defer conn.Close()
	for _, aid := range aids {
		k := keyArc(aid)
		if _, ok := aidmap[k]; !ok {
			keys = append(keys, k)
			aidmap[k] = aid
		}
	}
	if rs, err = conn.GetMulti(keys); err != nil {
		err = errors.Wrapf(err, "%v", keys)
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

// statsCache get stat cache by aids
func (d *Dao) statsCache(c context.Context, aids []int64) (cached map[int64]*api.Stat, missed []int64, err error) {
	cached = make(map[int64]*api.Stat, len(aids))
	var rs map[string]*memcache.Item
	conn := d.arcMc.Get(c)
	keys := make([]string, 0, len(aids))
	defer conn.Close()
	for _, aid := range aids {
		keys = append(keys, keyStat(aid))
	}
	if rs, err = conn.GetMulti(keys); err != nil {
		err = errors.Wrapf(err, "%v", keys)
		return
	}
	for _, r := range rs {
		st := &api.Stat{}
		if err = conn.Scan(r, st); err != nil {
			log.Error("conn.Scan(%s) error(%v)", r.Value, err)
			err = nil
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

func (d *Dao) avWithStCaches(c context.Context, aids []int64) (cached map[int64]*api.Arc, avMissed, stMissed []int64, err error) {
	cached = make(map[int64]*api.Arc, len(aids))
	conn := d.arcMc.Get(c)
	defer conn.Close()
	keys := make([]string, 0, len(aids)*2)
	avm := make(map[string]int64, len(aids))
	stm := make(map[string]int64, len(aids))
	for _, aid := range aids {
		ak := keyArc(aid)
		if _, ok := avm[ak]; !ok {
			keys = append(keys, ak)
			avm[ak] = aid
		}
		sk := keyStat(aid)
		if _, ok := stm[sk]; !ok {
			keys = append(keys, sk)
			stm[sk] = aid
		}
	}
	rs, err := conn.GetMulti(keys)
	if err != nil {
		err = errors.Wrapf(err, "%v", keys)
		return
	}
	stCached := make(map[string]*api.Stat, len(aids))
	for k, r := range rs {
		if aid, ok := avm[k]; ok {
			a := &api.Arc{}
			if err = conn.Scan(r, a); err != nil {
				log.Error("conn.Scan(%s) error(%v)", r.Value, err)
				err = nil
				continue
			}
			cached[aid] = a
			// delete hit key
			delete(avm, k)
		}
		if _, ok := stm[k]; ok {
			st := &api.Stat{}
			if err = conn.Scan(r, st); err != nil {
				log.Error("conn.Scan(%s) error(%v)", r.Value, err)
				err = nil
				continue
			}
			stCached[k] = st
		}
	}
	for k, st := range stCached {
		if a, ok := cached[st.Aid]; ok {
			a.Stat = *st
			// delete hit key
			delete(stm, k)
		}
	}
	// missed key
	avMissed = make([]int64, 0, len(avm))
	for _, aid := range avm {
		avMissed = append(avMissed, aid)

	}
	stMissed = make([]int64, 0, len(stm))
	for _, aid := range stm {
		stMissed = append(stMissed, aid)
	}
	return
}
