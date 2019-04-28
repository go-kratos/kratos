package account

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"go-common/app/admin/main/credit/conf"
	creditMDL "go-common/app/admin/main/credit/model"
	blocked "go-common/app/admin/main/credit/model/blocked"
	accgrpc "go-common/app/service/main/account/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// Dao is account dao.
type Dao struct {
	// grpc
	accountClient accgrpc.AccountClient
}

// New is initial for account .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{}
	var err error
	if d.accountClient, err = accgrpc.NewClient(c.AccClient); err != nil {
		panic(errors.WithMessage(err, "Failed to dial account service"))
	}
	return
}

// RPCInfo rpc info get by  muti mid .
func (d *Dao) RPCInfo(c context.Context, mid int64) (res *accgrpc.InfoReply, err error) {
	arg := &accgrpc.MidReq{Mid: mid}
	if res, err = d.accountClient.Info3(c, arg); err != nil {
		log.Error("d.accountClient.Info3 error(%v)", err)
	}
	return
}

// RPCInfos rpc info get by  muti mid .
func (d *Dao) RPCInfos(c context.Context, mids []int64) (res map[int64]*accgrpc.Info, err error) {
	var (
		g      errgroup.Group
		l      sync.RWMutex
		args   *accgrpc.MidsReq
		accRes *accgrpc.InfosReply
	)
	mids = creditMDL.ArrayUnique(mids)
	total := len(mids)
	pageNum := total / blocked.AccMaxPageSize
	if total%blocked.AccMaxPageSize != 0 {
		pageNum++
	}
	res = make(map[int64]*accgrpc.Info, total)
	for i := 0; i < pageNum; i++ {
		start := i * blocked.AccMaxPageSize
		end := (i + 1) * blocked.AccMaxPageSize
		if end > total {
			end = total
		}
		g.Go(func() (err error) {
			args = &accgrpc.MidsReq{Mids: mids[start:end]}
			if accRes, err = d.accountClient.Infos3(c, args); err != nil {
				log.Error("d.accountClient.Infos3(%+v) error(%v)", mids[start:end], err)
				err = nil
				return
			}
			for mid, info := range accRes.Infos {
				l.Lock()
				res[mid] = info
				l.Unlock()
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("g.Wait error(%v)", err)
	}
	return
}
