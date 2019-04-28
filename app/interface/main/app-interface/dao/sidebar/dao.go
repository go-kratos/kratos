package sidebar

import (
	"context"

	"go-common/app/interface/main/app-interface/conf"
	resmodel "go-common/app/service/main/resource/model"
	resrpc "go-common/app/service/main/resource/rpc/client"

	"github.com/pkg/errors"
)

// Dao is sidebar dao
type Dao struct {
	resRPC *resrpc.Service
}

// New initial sidebar dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		resRPC: resrpc.New(c.ResourceRPC),
	}
	return
}

// Sidebars from resource service
func (d *Dao) Sidebars(c context.Context) (res *resmodel.SideBars, err error) {
	if res, err = d.resRPC.SideBars(c); err != nil {
		err = errors.Wrapf(err, "d.resRPC.SideBars(%+v)")
	}
	return
}
