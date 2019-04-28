package resource

import (
	"context"

	"go-common/app/interface/main/app-intl/conf"
	"go-common/app/service/main/resource/model"
	rscrpc "go-common/app/service/main/resource/rpc/client"
	"go-common/library/ecode"
)

// Dao is archive dao.
type Dao struct {
	// rpc
	rscRPC *rscrpc.Service
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// rpc
		rscRPC: rscrpc.New(c.ResourceRPC),
	}
	return
}

// PlayerIcon is.
func (d *Dao) PlayerIcon(c context.Context) (res *model.PlayerIcon, err error) {
	if res, err = d.rscRPC.PlayerIcon(c); err != nil {
		if ecode.Cause(err) == ecode.NothingFound {
			res, err = nil, nil
		}
	}
	return
}

// PasterCID get all paster cid.
func (d *Dao) PasterCID(c context.Context) (cids map[int64]int64, err error) {
	return d.rscRPC.PasterCID(c)
}
