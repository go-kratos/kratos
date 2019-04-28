package assist

import (
	"context"

	"go-common/app/interface/main/app-intl/conf"
	"go-common/app/service/main/assist/model/assist"
	assistrpc "go-common/app/service/main/assist/rpc/client"

	"github.com/pkg/errors"
)

// Dao is assist dao
type Dao struct {
	assistRPC *assistrpc.Service
}

// New initial assist dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		assistRPC: assistrpc.New(c.AssistRPC),
	}
	return
}

// Assist get assists data from api.
func (d *Dao) Assist(c context.Context, upMid int64) (asss []int64, err error) {
	arg := &assist.ArgAssists{Mid: upMid}
	if asss, err = d.assistRPC.AssistIDs(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
