package archive

import (
	"context"

	"go-common/app/interface/main/app-channel/conf"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is archive dao.
type Dao struct {
	// rpc
	arcRPC *arcrpc.Service2
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// rpc
		arcRPC: arcrpc.New2(c.ArchiveRPC),
	}
	return
}

// Archive get archive by aid.
func (d *Dao) Archive(ctx context.Context, aid int64) (a *api.Arc, err error) {
	arg := &archive.ArgAid2{Aid: aid}
	a, err = d.arcRPC.Archive3(ctx, arg)
	return
}

// Archives multi get archives.
func (d *Dao) Archives(ctx context.Context, aids []int64) (as map[int64]*api.Arc, err error) {
	arg := &archive.ArgAids2{Aids: aids}
	as, err = d.arcRPC.Archives3(ctx, arg)
	return
}

// ArchivesWithPlayer archives witch player
func (d *Dao) ArchivesWithPlayer(c context.Context, aids []int64, qn int, platform string, fnver, fnval, build int) (res map[int64]*archive.ArchiveWithPlayer, err error) {
	if len(aids) == 0 {
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &archive.ArgPlayer{Aids: aids, Qn: qn, Platform: platform, Fnval: fnval, Fnver: fnver, RealIP: ip, Build: build}
	if res, err = d.arcRPC.ArchivesWithPlayer(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
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
