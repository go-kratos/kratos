package archive

import (
	"context"
	"runtime"
	"time"

	"go-common/app/interface/main/app-intl/conf"
	arcmdl "go-common/app/interface/main/app-intl/model/player/archive"
	"go-common/app/interface/main/app-intl/model/view"
	history "go-common/app/interface/main/history/model"
	hisrpc "go-common/app/interface/main/history/rpc/client"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

// Dao is archive dao.
type Dao struct {
	// http client
	client        *bm.Client
	realteURL     string
	commercialURL string
	relateRecURL  string
	playURL       string
	// rpc
	arcRPC  *arcrpc.Service2
	arcRPC2 *arcrpc.Service2
	hisRPC  *hisrpc.Service
	// mc
	mc        *memcache.Pool
	expireMc  int32
	expireRlt int32
	// chan
	mCh chan func()
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client:        bm.NewClient(c.HTTPWrite),
		realteURL:     c.Host.Data + _realteURL,
		commercialURL: c.Host.APICo + _commercialURL,
		relateRecURL:  c.Host.Data + _relateRecURL,
		playURL:       c.Host.Bvcvod + _playURL,
		// rpc
		arcRPC:  arcrpc.New2(c.ArchiveRPC),
		arcRPC2: arcrpc.New2(c.ArchiveRPC2),
		hisRPC:  hisrpc.New(c.HisRPC),
		// mc
		mc:        memcache.NewPool(c.Memcache.Feed.Config),
		expireMc:  int32(time.Duration(c.Memcache.Feed.Expire) / time.Second),
		expireRlt: int32(time.Duration(c.Memcache.Archive.RelateExpire) / time.Second),
		// mc proc
		mCh: make(chan func(), 10240),
	}
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go d.cacheproc()
	}
	return
}

// Ping ping check memcache connection
func (d *Dao) Ping(c context.Context) (err error) {
	return d.pingMC(c)
}

// Archives multi get archives.
func (d *Dao) Archives(c context.Context, aids []int64) (am map[int64]*api.Arc, err error) {
	if len(aids) == 0 {
		return
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		var missed []int64
		if am, missed, err = d.arcsCache(ctx, aids); err != nil {
			missed = aids
			log.Error("%+v", err)
			err = nil
		}
		if len(missed) == 0 {
			return
		}
		var tmp map[int64]*api.Arc
		arg := &archive.ArgAids2{Aids: missed}
		if tmp, err = d.arcRPC.Archives3(ctx, arg); err != nil {
			err = errors.Wrapf(err, "%v", arg)
			return
		}
		for aid, a := range tmp {
			am[aid] = a
		}
		return
	})
	var stm map[int64]*api.Stat
	g.Go(func() (err error) {
		var missed []int64
		if stm, missed, err = d.statsCache(ctx, aids); err != nil {
			missed = aids
			log.Error("%+v", err)
			err = nil
		}
		if len(missed) == 0 {
			return
		}
		tmp, err := d.arcRPC.Stats3(ctx, &archive.ArgAids2{Aids: missed})
		if err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		for _, st := range tmp {
			stm[st.Aid] = st
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	for aid, arc := range am {
		if st, ok := stm[aid]; ok {
			arc.Stat = *st
		}
	}
	return
}

// ArchivesWithPlayer archives witch player
func (d *Dao) ArchivesWithPlayer(c context.Context, aids []int64, qn int, platform string, fnver, fnval int) (res map[int64]*archive.ArchiveWithPlayer, err error) {
	if len(aids) == 0 {
		return
	}
	// 国际版暂时不秒开
	// ip := metadata.String(c, metadata.RemoteIP)
	ip := ""
	arg := &archive.ArgPlayer{Aids: aids, Qn: qn, Platform: platform, Fnval: fnval, Fnver: fnver, RealIP: ip}
	if res, err = d.arcRPC.ArchivesWithPlayer(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// Archive get archive mc->rpc.
func (d *Dao) Archive(c context.Context, aid int64) (a *api.Arc, err error) {
	if a, err = d.arcCache(c, aid); err != nil {
		log.Error("%+v", err)
	} else if a != nil {
		return
	}
	arg := &archive.ArgAid2{Aid: aid}
	if a, err = d.arcRPC.Archive3(c, arg); err != nil {
		log.Error("d.arcRPC.Archive3(%v) error(%v)", arg, err)
		if a, err = d.arcRPC2.Archive3(c, arg); err != nil {
			err = errors.Wrapf(err, "%v", arg)
		}
	}
	return
}

// Archive3 get archive.
func (d *Dao) Archive3(c context.Context, aid int64) (a *api.Arc, err error) {
	arg := &archive.ArgAid2{Aid: aid}
	if a, err = d.arcRPC.Archive3(c, arg); err != nil {
		log.Error("d.arcRPC.Archive3(%v) error(%+v)", arg, err)
		if a, err = d.arcRPC2.Archive3(c, arg); err != nil {
			err = errors.Wrapf(err, "d.arcRPC2.Archive3(%v)", arg)
			return
		}
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

// UpCount2 get upper count.
func (d *Dao) UpCount2(c context.Context, mid int64) (cnt int, err error) {
	arg := &archive.ArgUpCount2{Mid: mid}
	if cnt, err = d.arcRPC.UpCount2(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// ArchiveCache is
func (d *Dao) ArchiveCache(c context.Context, aid int64) (arc *arcmdl.Info, err error) {
	var (
		vp   *archive.View3
		cids []int64
	)
	if vp, err = d.ViewCache(c, aid); err != nil {
		log.Error("%+v", err)
	}
	if vp == nil || vp.Archive3 == nil || len(vp.Pages) == 0 || vp.AttrVal(archive.AttrBitIsMovie) == archive.AttrYes {
		if vp, err = d.View3(c, aid); err != nil {
			log.Error("%+v", err)
			err = ecode.NothingFound
			return
		}
	}
	if vp == nil || vp.Archive3 == nil || len(vp.Pages) == 0 {
		err = ecode.NothingFound
		return
	}
	for _, p := range vp.Pages {
		cids = append(cids, p.Cid)
	}
	arc = &arcmdl.Info{
		Aid:       vp.Aid,
		State:     vp.State,
		Mid:       vp.Author.Mid,
		Cids:      cids,
		Attribute: vp.Attribute,
	}
	return
}

// addCache add archive to mc or redis
func (d *Dao) addCache(f func()) {
	select {
	case d.mCh <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc write memcache and stat redis use goroutine
func (d *Dao) cacheproc() {
	for {
		f := <-d.mCh
		f()
	}
}
