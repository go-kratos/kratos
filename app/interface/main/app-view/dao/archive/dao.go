package archive

import (
	"context"
	"runtime"
	"time"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/view"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/sync/errgroup"

	history "go-common/app/interface/main/history/model"
	hisrpc "go-common/app/interface/main/history/rpc/client"
	arcrpc "go-common/app/service/main/archive/api/gorpc"

	"go-common/library/cache/memcache"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is archive dao.
type Dao struct {
	// http client
	client        *bm.Client
	httpClient    *bm.Client
	realteURL     string
	commercialURL string
	relateRecURL  string
	playURL       string
	// rpc
	arcRPC *arcrpc.Service2
	hisRPC *hisrpc.Service
	// memcache
	arcMc     *memcache.Pool
	expireArc int32
	expireRlt int32
	// chan
	mCh chan func()
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client:        bm.NewClient(c.HTTPWrite),
		httpClient:    bm.NewClient(c.HTTPClient),
		arcRPC:        arcrpc.New2(c.ArchiveRPC),
		hisRPC:        hisrpc.New(c.HisRPC),
		realteURL:     c.Host.Data + _realteURL,
		commercialURL: c.Host.APICo + _commercialURL,
		relateRecURL:  c.Host.Data + _relateRecURL,
		playURL:       c.Host.Bvcvod + _playURL,
		// memcache
		arcMc:     memcache.NewPool(c.Memcache.Archive.Config),
		expireArc: int32(time.Duration(c.Memcache.Archive.ArchiveExpire) / time.Second),
		expireRlt: int32(time.Duration(c.Memcache.Archive.RelateExpire) / time.Second),
		// mc proc
		mCh: make(chan func(), 10240),
	}
	// video db
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go d.cacheproc()
	}
	return
}

// Ping ping check memcache connection
func (d *Dao) Ping(c context.Context) (err error) {
	return d.pingMC(c)
}

// Archive3 get archive.
func (d *Dao) Archive3(c context.Context, aid int64) (a *api.Arc, err error) {
	arg := &archive.ArgAid2{Aid: aid}
	if a, err = d.arcRPC.Archive3(c, arg); err != nil {
		log.Error("d.arcRPC.Archive3(%v) error(%+v)", arg, err)
		return
	}
	return
}

// Archives multi get archives.
func (d *Dao) Archives(c context.Context, aids []int64) (as map[int64]*api.Arc, err error) {
	if len(aids) == 0 {
		return
	}
	var stm map[int64]*api.Stat
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		var missed []int64
		if as, missed, err = d.arcsCache(ctx, aids); err != nil {
			as = make(map[int64]*api.Arc, len(aids))
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
			log.Error("d.arcRPC.Archives3(%v) error(%v)", arg, err)
			return
		}
		for aid, a := range tmp {
			as[aid] = a
		}
		return
	})
	g.Go(func() (err error) {
		var missed []int64
		if stm, missed, err = d.statsCache(ctx, aids); err != nil {
			stm = make(map[int64]*api.Stat, len(aids))
			missed = aids
			log.Error("%+v", err)
			err = nil
		}
		if len(missed) == 0 {
			return
		}
		var tmp map[int64]*api.Stat
		arg := &archive.ArgAids2{Aids: missed}
		if tmp, err = d.arcRPC.Stats3(ctx, arg); err != nil {
			log.Error("d.arcRPC.Stats3(%v) error(%v)", arg, err)
			err = nil
			return
		}
		for aid, st := range tmp {
			stm[aid] = st
		}
		return
	})
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	for aid, a := range as {
		if st, ok := stm[aid]; ok {
			a.Stat = *st
		}
	}
	return
}

// Shot get video shot.
func (d *Dao) Shot(c context.Context, aid, cid int64) (shot *archive.Videoshot, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &archive.ArgCid2{Aid: aid, Cid: cid, RealIP: ip}
	return d.arcRPC.Videoshot2(c, arg)
}

// UpCount2 get upper count.
func (d *Dao) UpCount2(c context.Context, mid int64) (cnt int, err error) {
	arg := &archive.ArgUpCount2{Mid: mid}
	if cnt, err = d.arcRPC.UpCount2(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// UpArcs3 get upper archives.
func (d *Dao) UpArcs3(c context.Context, mid int64, pn, ps int) (as []*api.Arc, err error) {
	arg := &archive.ArgUpArcs2{Mid: mid, Pn: pn, Ps: ps}
	if as, err = d.arcRPC.UpArcs3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
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

// Archive 用的时候注意了！！！这个方法得到的稿件Stat不是最新的！！！
func (d *Dao) Archive(c context.Context, aid int64) (a *api.Arc, err error) {
	if a, err = d.arcCache(c, aid); err != nil {
		log.Error("%+v", err)
	} else if a != nil {
		return
	}
	if a, err = d.arcRPC.Archive3(c, &archive.ArgAid2{Aid: aid}); err != nil {
		log.Error("d.arcRPC.Archive3(%d) error(%v)", aid, err)
		return
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
