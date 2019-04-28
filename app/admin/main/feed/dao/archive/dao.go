package archive

import (
	"context"

	"go-common/app/admin/main/feed/conf"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

// Dao is archive dao.
type Dao struct {
	// rpc
	arcRPC *arcrpc.Service2
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		arcRPC: arcrpc.New2(c.ArchiveRPC),
	}
	return
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
