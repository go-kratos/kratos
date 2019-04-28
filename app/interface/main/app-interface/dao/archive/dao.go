package archive

import (
	"context"
	"time"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

// Dao is archive dao.
type Dao struct {
	// http client
	client *bm.Client
	// rpc
	arcRPC *arcrpc.Service2
	// memcache
	arcMc     *memcache.Pool
	expireArc int32
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client: bm.NewClient(c.HTTPWrite),
		arcRPC: arcrpc.New2(c.ArchiveRPC),
		// memcache
		arcMc:     memcache.NewPool(c.Memcache.Archive.Config),
		expireArc: int32(time.Duration(c.Memcache.Archive.ArchiveExpire) / time.Second),
	}
	return
}

// UpArcs3 get upper archives
func (d *Dao) UpArcs3(c context.Context, mid int64, pn, ps int) (as []*api.Arc, err error) {
	arg := &archive.ArgUpArcs2{Mid: mid, Pn: pn, Ps: ps}
	if as, err = d.arcRPC.UpArcs3(c, arg); err != nil {
		if ecode.Cause(err) == ecode.NothingFound {
			err = nil
		}
	}
	return
}

// UpCount2 get upper count.
func (d *Dao) UpCount2(c context.Context, mid int64) (cnt int, err error) {
	arg := &archive.ArgUpCount2{Mid: mid}
	return d.arcRPC.UpCount2(c, arg)
}

// Archives multi get archives.
func (d *Dao) Archives(c context.Context, aids []int64) (as map[int64]*api.Arc, err error) {
	if len(aids) == 0 {
		return
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		var missed []int64
		if as, missed, err = d.arcsCache(ctx, aids); err != nil {
			log.Error("%+v", err)
			missed = aids
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
	var stm map[int64]*api.Stat
	g.Go(func() (err error) {
		var missed []int64
		if stm, missed, err = d.statsCache(ctx, aids); err != nil {
			log.Error("%+v", err)
			missed = aids
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
		return
	}
	for aid, a := range as {
		if st, ok := stm[aid]; ok {
			a.Stat = *st
		}
	}
	return
}

// Archives2 multi get archives.
func (d *Dao) Archives2(c context.Context, aids []int64) (am map[int64]*api.Arc, err error) {
	if len(aids) == 0 {
		return
	}
	var avMissed, stMissed []int64
	if am, avMissed, stMissed, err = d.avWithStCaches(c, aids); err != nil {
		avMissed = aids
		log.Error("%+v", err)
		err = nil
	}
	g, ctx := errgroup.WithContext(c)
	if len(avMissed) != 0 {
		g.Go(func() (err error) {
			arg := &archive.ArgAids2{Aids: avMissed}
			avm, err := d.arcRPC.Archives3(ctx, arg)
			if err != nil {
				err = errors.Wrapf(err, "%v", arg)
				return
			}
			for aid, a := range avm {
				am[aid] = a
			}
			return
		})
	}
	var stm map[int64]*api.Stat
	if len(stMissed) != 0 {
		g.Go(func() error {
			arg := &archive.ArgAids2{Aids: stMissed}
			if stm, err = d.arcRPC.Stats3(ctx, arg); err != nil {
				log.Error("d.arcRPC.Stats3(%v) error(%v)", arg, err)
			}
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		return
	}
	for aid, a := range am {
		if st, ok := stm[aid]; ok {
			a.Stat = *st
		}
	}
	return
}
