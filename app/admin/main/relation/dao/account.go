package dao

import (
	"context"
	"sync"

	accountApi "go-common/app/service/main/account/api"
	account "go-common/app/service/main/account/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_perCall = 50
)

func chain(ids ...[]int64) []int64 {
	res := make([]int64, 0, len(ids))
	for _, l := range ids {
		res = append(res, l...)
	}
	return res
}

func uniq(ids ...[]int64) []int64 {
	hm := make(map[int64]struct{})
	for _, i := range chain(ids...) {
		hm[i] = struct{}{}
	}
	res := make([]int64, 0, len(ids))
	for i := range hm {
		res = append(res, i)
	}
	return res
}

// RPCInfos rpc info get by  muti mid .
func (d *Dao) RPCInfos(c context.Context, mids []int64) (res map[int64]*account.Info, err error) {
	var (
		g errgroup.Group
		l sync.RWMutex
	)
	mids = uniq(mids)
	total := len(mids)
	pageNum := total / _perCall
	if total%_perCall != 0 {
		pageNum++
	}
	res = make(map[int64]*account.Info, total)
	for i := 0; i < pageNum; i++ {
		start := i * _perCall
		end := (i + 1) * _perCall
		if end > total {
			end = total
		}
		g.Go(func() (err error) {
			midsReq := &accountApi.MidsReq{Mids: mids[start:end]}
			infosReply, err := d.accountClient.Infos3(c, midsReq)
			if err != nil {
				log.Error("d.accountClient.Infos3(%+v) error(%v)", midsReq, err)
				err = nil
				return
			}
			for mid, info := range infosReply.Infos {
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
