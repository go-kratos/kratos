package archive

import (
	"context"
	"strconv"
	"sync"

	"go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_prefixArc3 = "a3p_"
	_prefixDesc = "desc_"
)

func descKey(aid int64) string {
	return _prefixDesc + strconv.FormatInt(aid, 10)
}

func arcPBKey(aid int64) string {
	return _prefixArc3 + strconv.FormatInt(aid, 10)
}

// archive3Cache get a archive info from cache.
func (d *Dao) archive3Cache(c context.Context, aid int64) (a *api.Arc, err error) {
	var (
		key  = arcPBKey(aid)
		rp   *memcache.Item
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			d.errProm.Incr("archive_mc")
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	a = &api.Arc{}
	if err = conn.Scan(rp, a); err != nil {
		log.Error("conn.Scan(%s) error(%v)", rp.Value, err)
	}
	return
}

// archive3Caches multi get archives, return map[aid]*Archive and missed aids.
func (d *Dao) archive3Caches(c context.Context, aids []int64) (cached map[int64]*api.Arc, err error) {
	var (
		keys   = make([]string, 0)
		keyMap = make(map[int64]struct{}, len(aids))
		eg     = errgroup.Group{}
		mutex  = sync.Mutex{}
	)
	cached = make(map[int64]*api.Arc)
	for _, aid := range aids {
		if _, ok := keyMap[aid]; ok {
			continue
		}
		keyMap[aid] = struct{}{}
		keys = append(keys, arcPBKey(aid))
		if len(keys) == 50 {
			egks := keys
			eg.Go(func() (err error) {
				as, _ := d._archiveCaches(c, egks)
				mutex.Lock()
				for _, a := range as {
					cached[a.Aid] = a
				}
				mutex.Unlock()
				return
			})
			keys = make([]string, 0)
		}
	}
	if len(keys) > 0 {
		eg.Go(func() (err error) {
			as, _ := d._archiveCaches(c, keys)
			mutex.Lock()
			for _, a := range as {
				cached[a.Aid] = a
			}
			mutex.Unlock()
			return
		})
	}
	eg.Wait()
	return
}

// _archiveCaches is
func (d *Dao) _archiveCaches(c context.Context, keys []string) (as map[int64]*api.Arc, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	rs, err := conn.GetMulti(keys)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		return
	}
	as = make(map[int64]*api.Arc)
	for _, r := range rs {
		a := &api.Arc{}
		if err = conn.Scan(r, a); err != nil {
			log.Error("conn.Scan error(%v)", err)
			continue
		}
		as[a.Aid] = a
	}
	return
}

// addArchivePBCache set archive into cache.
func (d *Dao) addArchive3Cache(c context.Context, a *api.Arc) (err error) {
	key := arcPBKey(a.Aid)
	conn := d.mc.Get(c)
	if err = conn.Set(&memcache.Item{Key: key, Object: a, Flags: memcache.FlagProtobuf, Expiration: 0}); err != nil {
		log.Error("memcache.Set(%s) error(%v)", key, err)
		d.errProm.Incr("archive_mc")
	}
	conn.Close()
	return
}

func (d *Dao) descCache(c context.Context, aid int64) (desc string, err error) {
	key := descKey(aid)
	conn := d.mc.Get(c)
	defer conn.Close()
	var item *memcache.Item
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	if err = conn.Scan(item, &desc); err != nil {
		log.Error("conn.Scan error(%v)", err)
		return
	}
	return
}

func (d *Dao) addDescCache(c context.Context, aid int64, desc string) (err error) {
	key := descKey(aid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Value: []byte(desc), Flags: memcache.FlagRAW, Expiration: 0}); err != nil {
		log.Error("conn.Set(%s, %s) error(%v)", key, desc, err)
		return
	}
	return
}
