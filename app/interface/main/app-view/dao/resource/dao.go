package resource

import (
	"context"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/service/main/resource/model"
	rscrpc "go-common/app/service/main/resource/rpc/client"
	"go-common/library/ecode"

	"github.com/pkg/errors"
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

func (d *Dao) Paster(c context.Context, plat, adType int8, aid, typeID, buvid string) (res *model.Paster, err error) {
	arg := &model.ArgPaster{Platform: plat, AdType: adType, Aid: aid, TypeId: typeID, Buvid: buvid}
	if res, err = d.rscRPC.PasterAPP(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

func (d *Dao) PlayerIcon(c context.Context) (res *model.PlayerIcon, err error) {
	if res, err = d.rscRPC.PlayerIcon(c); err != nil {
		if ecode.Cause(err) == ecode.NothingFound {
			res, err = nil, nil
		}
	}
	return
}
