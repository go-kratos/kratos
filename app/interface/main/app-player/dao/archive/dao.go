package archive

import (
	"context"
	"fmt"
	"runtime"

	"go-common/app/interface/main/app-player/conf"
	"go-common/app/interface/main/app-player/model/archive"
	arcrpc "go-common/app/service/main/archive/api"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

// Dao is archive dao.
type Dao struct {
	// memcache
	arcMc *memcache.Pool
	// chan
	mCh chan func()
	// rpc
	arcRPC arcrpc.ArchiveClient
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// memcache
		arcMc: memcache.NewPool(c.Memcache),
		// mc proc
		mCh: make(chan func(), 1024),
	}
	var err error
	d.arcRPC, err = arcrpc.NewClient(c.ArchiveClient)
	if err != nil {
		panic(fmt.Sprintf("archive NewClient error(%v)", err))
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		go d.cacheproc()
	}
	return
}

// Ping ping check memcache connection
func (d *Dao) Ping(c context.Context) (err error) {
	return d.pingMC(c)
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
		f, ok := <-d.mCh
		if !ok {
			return
		}
		f()
	}
}

// ArchiveCache is
func (d *Dao) ArchiveCache(c context.Context, aid int64) (arc *archive.Info, err error) {
	if arc, err = d.archiveCache(c, aid); err != nil {
		log.Error("%+v", err)
		err = nil
	}
	if arc != nil {
		return
	}
	var (
		view *arcrpc.ViewReply
		cids []int64
	)
	if view, err = d.arcRPC.View(c, &arcrpc.ViewRequest{Aid: aid}); err != nil {
		log.Error("d.arcRPC.View3(%d) error(%+v)", aid, err)
		return
	}
	for _, p := range view.Pages {
		cids = append(cids, p.Cid)
	}
	arc = &archive.Info{
		Aid:       view.Arc.Aid,
		State:     view.Arc.State,
		Mid:       view.Arc.Author.Mid,
		Cids:      cids,
		Attribute: view.Arc.Attribute,
	}
	d.addCache(func() {
		d.addArchiveCache(context.Background(), aid, arc)
	})
	return
}

// Views is
func (d *Dao) Views(c context.Context, aids []int64) (arcs map[int64]*archive.Info, err error) {
	var reply *arcrpc.ViewsReply
	if reply, err = d.arcRPC.Views(c, &arcrpc.ViewsRequest{Aids: aids}); err != nil {
		return
	}
	arcs = make(map[int64]*archive.Info)
	for _, v := range reply.Views {
		var (
			info = new(archive.Info)
			cids []int64
		)
		info.Aid = v.Arc.Aid
		info.State = v.Arc.State
		info.Mid = v.Arc.Author.Mid
		info.Attribute = v.Arc.Attribute
		for _, p := range v.Pages {
			cids = append(cids, p.Cid)
		}
		info.Cids = cids
		arcs[info.Aid] = info
	}
	return
}
