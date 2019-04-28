package service

import (
	"context"
	"math"
	"sync"

	account "go-common/app/service/main/account/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

func (s *Service) accountInfos(c context.Context, mids []int64) (res map[int64]*account.Info, err error) {
	var (
		g        errgroup.Group
		mu       sync.Mutex
		pagesize = 50
		ids      = make([]int64, 0, len(mids))
		midMap   = make(map[int64]struct{})
	)
	res = make(map[int64]*account.Info)
	if len(mids) == 0 {
		return
	}
	for _, mid := range mids {
		if _, ok := midMap[mid]; !ok {
			midMap[mid] = struct{}{}
			ids = append(ids, mid)
		}
	}
	total := len(ids)
	pageNum := int(math.Ceil(float64(total) / float64(pagesize)))
	for i := 0; i < pageNum; i++ {
		start := i * pagesize
		end := (i + 1) * pagesize
		if end > total {
			end = total
		}
		g.Go(func() (err error) {
			var (
				arg   = &account.MidsReq{Mids: ids[start:end]}
				reply *account.InfosReply
			)
			if reply, err = s.accountRPC.Infos3(c, arg); err != nil {
				log.Error("accRPC.Infos3(%+v) error(%v)", arg, err)
				return
			}
			for mid, info := range reply.GetInfos() {
				mu.Lock()
				res[mid] = info
				mu.Unlock()
			}
			return
		})
	}
	err = g.Wait()
	return
}
