package archive

import (
	"context"
	"strconv"

	"go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixPBPage = "psb_"
)

func pagePBKey(aid int64) string {
	return _prefixPBPage + strconv.FormatInt(aid, 10)
}

func videoPBKey2(aid, cid int64) string {
	return _prefixPBPage + strconv.FormatInt(aid, 10) + strconv.FormatInt(cid, 10)
}

type batchVdo3 struct {
	cached map[int64][]*api.Page
	missed []int64
	err    error
}

// addPageCache3 add pages cache by aid.
func (d *Dao) addPageCache3(c context.Context, aid int64, ps []*api.Page) (err error) {
	var vs = &api.AidVideos{Aid: aid, Pages: ps}
	key := pagePBKey(aid)
	conn := d.mc.Get(c)
	if err = conn.Set(&memcache.Item{Key: key, Object: vs, Expiration: 0, Flags: memcache.FlagProtobuf}); err != nil {
		d.errProm.Incr("video_mc")
		log.Error("memcache.Set(%s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// addVideoCache3 add video2 cache by aid & cid.
func (d *Dao) addVideoCache3(c context.Context, aid, cid int64, p *api.Page) (err error) {
	key := videoPBKey2(aid, cid)
	conn := d.mc.Get(c)
	if err = conn.Set(&memcache.Item{Key: key, Object: p, Expiration: 0, Flags: memcache.FlagProtobuf}); err != nil {
		d.errProm.Incr("video3_mc")
		log.Error("conn.Set(%s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// DelVideoCache3 del video2 cache by aid & cid.
func (d *Dao) DelVideoCache3(c context.Context, aid, cid int64) (err error) {
	key := videoPBKey2(aid, cid)
	conn := d.mc.Get(c)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			d.errProm.Incr("video3_mc")
			log.Error("conn.Set(%s) error(%v)", key, err)
		}
	}
	conn.Close()
	return
}

// pageCache3 get page cache by aid.
func (d *Dao) pageCache3(c context.Context, aid int64) (ps []*api.Page, err error) {
	var (
		key  = pagePBKey(aid)
		rp   *memcache.Item
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			ps = nil
		} else {
			d.errProm.Incr("video_mc")
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	var vs = new(api.AidVideos)
	if err = conn.Scan(rp, vs); err != nil {
		log.Error("conn.Scan(%s) error(%v)", rp.Value, err)
		ps = nil
		return
	}
	ps = vs.Pages
	return
}

// pagesCache3 get pages cache by aids
func (d *Dao) pagesCache3(c context.Context, aids []int64) (cached map[int64][]*api.Page, missed []int64, err error) {
	var (
		lt    = len(aids)
		times int
		reqCh chan batchVdo3
	)
	if lt%_multiInterval == 0 {
		times = lt / _multiInterval
	} else {
		times = lt/_multiInterval + 1
	}
	reqCh = make(chan batchVdo3, times)
	for i := 0; i < times; i++ {
		var tmps []int64
		if i == times-1 {
			tmps = aids[i*_multiInterval:]
		} else {
			tmps = aids[i*_multiInterval : (i+1)*_multiInterval]
		}
		go d._videoCaches3(c, tmps, reqCh)
	}
	if times == 1 {
		req := <-reqCh
		cached = req.cached
		missed = req.missed
		err = req.err
		return
	}
	cached = make(map[int64][]*api.Page, len(aids))
	for i := 0; i < times; i++ {
		req := <-reqCh
		for k, v := range req.cached {
			cached[k] = v
		}
		missed = append(missed, req.missed...)
	}
	return
}

// videoCache3 get video cache by aid & cid.
func (d *Dao) videoCache3(c context.Context, aid, cid int64) (p *api.Page, err error) {
	var (
		key  = videoPBKey2(aid, cid)
		rp   *memcache.Item
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			p = nil
		} else {
			d.errProm.Incr("video2_mc")
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	p = new(api.Page)
	if err = conn.Scan(rp, p); err != nil {
		log.Error("conn.Scan(%s) error(%v)", rp.Value, err)
		p = nil
		return
	}
	return
}

func (d *Dao) _videoCaches3(c context.Context, aids []int64, reqCh chan<- batchVdo3) {
	var (
		req  batchVdo3
		keys = make([]string, 0, len(aids))
		am   = make(map[int64]struct{}, len(aids))
		rps  map[string]*memcache.Item
		err  error
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	for _, aid := range aids {
		keys = append(keys, pagePBKey(aid))
		am[aid] = struct{}{}
	}
	if rps, err = conn.GetMulti(keys); err != nil {
		req.missed = aids
		req.cached = make(map[int64][]*api.Page)
		req.err = err
		reqCh <- req
		log.Error("conn.Gets error(%v)", err)
		d.errProm.Incr("video_mc")
		return
	}
	var (
		cached = make(map[int64][]*api.Page, len(aids))
		missed []int64
	)
	for _, rp := range rps {
		var v = new(api.AidVideos)
		if e := conn.Scan(rp, v); e != nil {
			log.Error("conn.Scan(%s) error(%v)", rp.Value, e)
			d.errProm.Incr("video_mc")
			continue
		}
		cached[v.Aid] = v.Pages
		delete(am, v.Aid)
	}
	for aid := range am {
		missed = append(missed, aid)
	}
	req.cached = cached
	req.missed = missed
	reqCh <- req
}
