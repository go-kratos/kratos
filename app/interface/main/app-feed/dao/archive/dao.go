package archive

import (
	"context"
	"time"

	"go-common/app/interface/main/app-feed/conf"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

// Dao is archive dao.
type Dao struct {
	// rpc
	arcRPC *arcrpc.Service2
	// mc
	mc       *memcache.Pool
	expireMc int32
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// rpc
		arcRPC: arcrpc.New2(c.ArchiveRPC),
		// mc
		mc:       memcache.NewPool(c.Memcache.Feed.Config),
		expireMc: int32(time.Duration(c.Memcache.Feed.ExpireArchive) / time.Second),
	}
	return
}

func (d *Dao) PingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{Key: "ping", Value: []byte{1}, Flags: memcache.FlagRAW, Expiration: d.expireMc}
	err = conn.Set(item)
	conn.Close()
	return
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
		if tmp, err = d.arcRPC.Archives3(c, arg); err != nil {
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
func (d *Dao) ArchivesWithPlayer(c context.Context, aids []int64, qn int, platform string, fnver, fnval, forceHost, build int) (res map[int64]*archive.ArchiveWithPlayer, err error) {
	if len(aids) == 0 {
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &archive.ArgPlayer{Aids: aids, Qn: qn, Platform: platform, Fnval: fnval, Fnver: fnver, RealIP: ip, ForceHost: forceHost, Build: build}
	if res, err = d.arcRPC.ArchivesWithPlayer(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
