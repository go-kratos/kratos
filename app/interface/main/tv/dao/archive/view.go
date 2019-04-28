package archive

import (
	"context"

	history "go-common/app/interface/main/history/model"
	"go-common/app/interface/main/tv/model/view"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Archive3 picks one archive data
func (d *Dao) Archive3(c context.Context, aid int64) (arc *arcwar.Arc, err error) {
	var (
		arcReply *arcwar.ArcReply
		arg      = &arcwar.ArcRequest{Aid: aid}
	)
	arcReply, err = d.arcClient.Arc(c, arg)
	if err != nil {
		log.Error("d.arcRPC.Archive(%v) error(%+v)", arg, err)
		if ecode.Cause(err) == ecode.NothingFound {
			err = nil
			return
		}
	}
	arc = arcReply.Arc
	return
}

// LoadViews picks view meta information from Cache & RPC
func (d *Dao) LoadViews(ctx context.Context, aids []int64) (resMetas map[int64]*arcwar.ViewReply) {
	var (
		missedMetas = make(map[int64]*arcwar.ViewReply)
		addCache    = true
	)
	resMetas = make(map[int64]*arcwar.ViewReply, len(aids))
	cachedMetas, missed, err := d.ViewsCache(ctx, aids)
	if err != nil {
		log.Error("LoadViews ViewsCache Sids:%v, Error:%v", aids, err)
		err = nil
		addCache = false // mc error, we don't add
	}
	if len(missed) > 0 {
		for _, vv := range missed {
			if vp, _ := d.view3(ctx, vv); vp != nil {
				missedMetas[vv] = vp
				resMetas[vv] = vp
			}
		}
	}
	// merge info from DB and the info from MC
	for sid, v := range cachedMetas {
		resMetas[sid] = v
	}
	log.Info("LoadViews Info Hit %d, Missed %d", len(cachedMetas), len(missed))
	if addCache && len(missedMetas) > 0 { // async Reset the DB data in MC for next time
		log.Info("Set MissedMetas %d Data in MC", missedMetas)
		for aid, vp := range missedMetas {
			d.AddViewCache(aid, vp)
		}
	}
	return
}

// GetView gets the aid's View info from Cache or RPC
func (d *Dao) GetView(c context.Context, aid int64) (vp *arcwar.ViewReply, err error) {
	var (
		addCache = false
	)
	if vp, err = d.viewCache(c, aid); err != nil {
		log.Error("ViewPage viewCache AID %d, Err %v", aid, err)
	}
	if !validView(vp, true) { // back source
		if vp, err = d.view3(c, aid); err != nil {
			log.Error("%+v", err)
			err = ecode.NothingFound
			return
		}
		addCache = true
	}
	if !validView(vp, false) {
		err = ecode.NothingFound
		return
	}
	if addCache {
		d.AddViewCache(aid, vp)
		return
	}
	return
}

// view3 view archive with pages pb.
func (d *Dao) view3(c context.Context, aid int64) (reply *arcwar.ViewReply, err error) {
	var arg = &arcwar.ViewRequest{Aid: aid}
	if reply, err = d.arcClient.View(c, arg); err != nil {
		log.Error("d.arcRPC.view3(%v) error(%+v)", arg, err)
		if ecode.Cause(err) == ecode.NothingFound {
			err = nil
			return
		}
	}
	return
}

// viewCache get view cache from remote memecache .
func (d *Dao) viewCache(c context.Context, aid int64) (vs *arcwar.ViewReply, err error) {
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
	vs = &arcwar.ViewReply{Arc: &arcwar.Arc{}}
	if err = conn.Scan(r, vs); err != nil {
		vs = nil
		err = errors.Wrapf(err, "conn.Scan(%s)", r.Value)
	}
	return
}

// Progress is  archive plays progress .
func (d *Dao) Progress(c context.Context, aid, mid int64) (h *view.History, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &history.ArgPro{Mid: mid, Aids: []int64{aid}, RealIP: ip}
	his, err := d.hisRPC.Progress(c, arg)
	if err != nil {
		log.Error("d.hisRPC.Progress(%v) error(%v)", arg, err)
		return
	}
	if his[aid] != nil {
		h = &view.History{Cid: his[aid].Cid, Progress: his[aid].Pro}
	}
	return
}

// ViewsCache pick view from cache
func (d *Dao) ViewsCache(c context.Context, aids []int64) (cached map[int64]*arcwar.ViewReply, missed []int64, err error) {
	if len(aids) == 0 {
		return
	}
	var allKeys []string
	cached = make(map[int64]*arcwar.ViewReply, len(aids))
	idmap := make(map[string]int64, len(aids))
	for _, id := range aids {
		k := keyView(id)
		allKeys = append(allKeys, k)
		idmap[k] = id
	}
	conn := d.arcMc.Get(c)
	defer conn.Close()
	replys, err := conn.GetMulti(allKeys)
	if err != nil {
		log.Error("conn.Gets(%v) error(%v)", allKeys, err)
		err = nil
		return
	}
	for key, item := range replys {
		vp := &arcwar.ViewReply{}
		if err = conn.Scan(item, vp); err != nil {
			log.Error("item.Scan(%s) error(%v)", item.Value, err)
			err = nil
			continue
		}
		if !validView(vp, true) {
			continue
		}
		cached[idmap[key]] = vp
		delete(idmap, key)
	}
	missed = make([]int64, 0, len(idmap))
	for _, id := range idmap {
		missed = append(missed, id)
	}
	return
}

// Archives multi get archives.
func (d *Dao) Archives(c context.Context, aids []int64) (as map[int64]*arcwar.Arc, err error) {
	if len(aids) == 0 {
		return
	}
	var (
		missed []int64
		tmp    map[int64]*arcwar.Arc
		reply  *arcwar.ArcsReply
	)
	if as, missed, err = d.arcsCache(c, aids); err != nil {
		as = make(map[int64]*arcwar.Arc, len(aids))
		missed = aids
		log.Error("%+v", err)
		err = nil
	}
	if len(missed) == 0 {
		return
	}
	arg := &arcwar.ArcsRequest{Aids: missed}
	if reply, err = d.arcClient.Arcs(c, arg); err != nil {
		log.Error("d.arcRPC.Archives3(%v) error(%v)", arg, err)
		if reply, err = d.arcClient.Arcs(c, arg); err != nil {
			err = errors.Wrapf(err, "%v", arg)
			return
		}
	}
	tmp = reply.Arcs
	for aid, a := range tmp {
		as[aid] = a
		d.AddArcCache(aid, a) // re-fill the cache
	}
	return
}
