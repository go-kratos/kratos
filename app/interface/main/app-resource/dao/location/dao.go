package location

import (
	"context"
	"strconv"

	"go-common/app/interface/main/app-resource/conf"
	locmdl "go-common/app/service/main/location/model"
	locrpc "go-common/app/service/main/location/rpc/client"
	"go-common/library/log"
)

// Dao is location dao.
type Dao struct {
	// rpc
	locRPC *locrpc.Service
}

// New new a location dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// rpc
		locRPC: locrpc.New(c.LocationRPC),
	}
	return
}

func (d *Dao) Info(c context.Context, ipaddr string) (info *locmdl.Info, err error) {
	if info, err = d.locRPC.Info(c, &locmdl.ArgIP{IP: ipaddr}); err != nil {
		log.Error("%v", err)
	}
	return
}

func (d *Dao) AuthPIDs(c context.Context, pids, ipaddr string) (res map[string]*locmdl.Auth, err error) {
	var auths map[int64]*locmdl.Auth
	if auths, err = d.locRPC.AuthPIDs(c, &locmdl.ArgPids{Pids: pids, IP: ipaddr}); err != nil {
		log.Error("%v", err)
		return
	}
	res = make(map[string]*locmdl.Auth)
	for pid, auth := range auths {
		p := strconv.FormatInt(pid, 10)
		res[p] = auth
	}
	return
}
