package archive

import (
	"context"
	"strconv"

	"go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixStatPB  = "stp_"
	_prefixClickPB = "clkp_"
)

type batchStat3 struct {
	cached map[int64]*api.Stat
	missed []int64
	err    error
}

func statPBKey(aid int64) string {
	return _prefixStatPB + strconv.FormatInt(aid, 10)
}

func clickPBKey(aid int64) string {
	return _prefixClickPB + strconv.FormatInt(aid, 10)
}

// statCache3 get a archive stat from cache.
func (d *Dao) statCache3(c context.Context, aid int64) (st *api.Stat, err error) {
	var (
		key  = statPBKey(aid)
		rp   *memcache.Item
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			d.errProm.Incr("stat_mc")
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	st = new(api.Stat)
	if err = conn.Scan(rp, st); err != nil {
		log.Error("conn.Scan(%s) error(%v)", rp.Value, err)
		return
	}
	st.DisLike = 0
	return
}

// statCaches3 multi get archives stat, return map[aid]*Stat and missed aids.
func (d *Dao) statCaches3(c context.Context, aids []int64) (cached map[int64]*api.Stat, missed []int64, err error) {
	var (
		lt    = len(aids)
		times int
		reqCh chan batchStat3
	)
	if lt%_multiInterval == 0 {
		times = lt / _multiInterval
	} else {
		times = lt/_multiInterval + 1
	}
	reqCh = make(chan batchStat3, times)
	for i := 0; i < times; i++ {
		var tmps []int64
		if i == times-1 {
			tmps = aids[i*_multiInterval:]
		} else {
			tmps = aids[i*_multiInterval : (i+1)*_multiInterval]
		}
		go d._statCaches3(c, tmps, reqCh)
	}
	if times == 1 {
		req := <-reqCh
		cached = req.cached
		missed = req.missed
		err = req.err
		return
	}
	cached = make(map[int64]*api.Stat, len(aids))
	for i := 0; i < times; i++ {
		req := <-reqCh
		for k, v := range req.cached {
			cached[k] = v
		}
		missed = append(missed, req.missed...)
	}
	return
}

func (d *Dao) _statCaches3(c context.Context, aids []int64, reqCh chan<- batchStat3) {
	var (
		req  batchStat3
		keys = make([]string, 0, len(aids))
		am   = make(map[int64]struct{}, len(aids))
		rps  map[string]*memcache.Item
		err  error
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	for _, aid := range aids {
		keys = append(keys, statPBKey(aid))
		am[aid] = struct{}{}
	}
	if rps, err = conn.GetMulti(keys); err != nil {
		req.missed = aids
		req.err = err
		reqCh <- req
		log.Error("conn.Gets error(%v)", err)
		d.errProm.Incr("stat_mc")
		return
	}
	var (
		cached = make(map[int64]*api.Stat, len(aids))
		missed []int64
	)
	for _, rp := range rps {
		var st = new(api.Stat)
		if e := conn.Scan(rp, st); e != nil {
			log.Error("conn.Scan(%s) error(%v)", rp.Value, e)
			continue
		}
		st.DisLike = 0
		cached[st.Aid] = st
		delete(am, st.Aid)
	}
	for aid := range am {
		missed = append(missed, aid)
	}
	req.cached = cached
	req.missed = missed
	reqCh <- req
}

// addStatCache set archive stat into cache.
func (d *Dao) addStatCache3(c context.Context, st *api.Stat) (err error) {
	key := statPBKey(st.Aid)
	conn := d.mc.Get(c)
	st.DisLike = 0
	if err = conn.Set(&memcache.Item{Key: key, Object: st, Flags: memcache.FlagProtobuf, Expiration: 0}); err != nil {
		d.errProm.Incr("stat_mc")
		log.Error("memcache.Set(%s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// clickCache get a archive click from cache.
func (d *Dao) clickCache3(c context.Context, aid int64) (clk *api.Click, err error) {
	var (
		key  = clickPBKey(aid)
		rp   *memcache.Item
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			d.errProm.Incr("stat_mc")
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	clk = new(api.Click)
	if err = conn.Scan(rp, clk); err != nil {
		log.Error("conn.Scan(%s) error(%v)", rp.Value, err)
	}
	return
}

// addClickCache set archive click into cache.
func (d *Dao) addClickCache3(c context.Context, clk *api.Click) (err error) {
	key := clickPBKey(clk.Aid)
	conn := d.mc.Get(c)
	if err = conn.Set(&memcache.Item{Key: key, Object: clk, Flags: memcache.FlagProtobuf, Expiration: 0}); err != nil {
		d.errProm.Incr("stat_mc")
		log.Error("memcache.Set(%s) error(%v)", key, err)
	}
	conn.Close()
	return
}
