package service

import (
	"context"
	"sync"

	artmdl "go-common/app/interface/openplatform/article/model"
	accmdl "go-common/app/service/main/account/model"
	"go-common/app/service/main/feed/dao"
	"go-common/library/log"

	"go-common/library/sync/errgroup"
)

const _upsArtBulkSize = 50

// attenUpArticles get new articles of attention uppers.
func (s *Service) attenUpArticles(c context.Context, minTotalCount int, mid int64, ip string) (res map[int64][]*artmdl.Meta, err error) {
	var mids []int64
	arg := &accmdl.ArgMid{Mid: mid}
	if mids, err = s.accRPC.Attentions3(c, arg); err != nil {
		dao.PromError("关注rpc接口:Attentions", "s.accRPC.Attentions(%d) error(%v)", mid, err)
		return
	}
	if len(mids) == 0 {
		return
	}
	count := minTotalCount/len(mids) + s.c.Feed.MinUpCnt
	return s.upsArticle(c, count, ip, mids...)
}

func (s *Service) upsArticle(c context.Context, count int, ip string, mids ...int64) (res map[int64][]*artmdl.Meta, err error) {
	dao.MissedCount.Add("upArt", int64(len(mids)))
	var (
		group      *errgroup.Group
		errCtx     context.Context
		midsLen, i int
		mutex      = sync.Mutex{}
	)
	res = make(map[int64][]*artmdl.Meta)
	group, errCtx = errgroup.WithContext(c)
	midsLen = len(mids)
	for ; i < midsLen; i += _upsArtBulkSize {
		var partMids []int64
		if i+_upsArcBulkSize > midsLen {
			partMids = mids[i:]
		} else {
			partMids = mids[i : i+_upsArtBulkSize]
		}
		group.Go(func() (err error) {
			var tmpRes map[int64][]*artmdl.Meta
			arg := &artmdl.ArgUpsArts{Mids: partMids, Pn: 1, Ps: count, RealIP: ip}
			if tmpRes, err = s.artRPC.UpsArtMetas(errCtx, arg); err != nil {
				log.Error("s.artRPC.UpsArtMetas(%+v) error(%v)", arg, err)
				err = nil
				return
			}
			mutex.Lock()
			for mid, arcs := range tmpRes {
				for _, arc := range arcs {
					if arc.AttrVal(artmdl.AttrBitNoDistribute) {
						continue
					}
					res[mid] = append(res[mid], arc)
				}
			}
			mutex.Unlock()
			return
		})
	}
	group.Wait()
	return
}

func (s *Service) articles(c context.Context, ip string, aids ...int64) (res map[int64]*artmdl.Meta, err error) {
	var (
		mutex    = sync.Mutex{}
		bulkSize = s.c.Feed.BulkSize
	)
	res = make(map[int64]*artmdl.Meta, len(aids))
	group, errCtx := errgroup.WithContext(c)
	aidsLen := len(aids)
	for i := 0; i < aidsLen; i += bulkSize {
		var partAids []int64
		if i+bulkSize < aidsLen {
			partAids = aids[i : i+bulkSize]
		} else {
			partAids = aids[i:aidsLen]
		}
		group.Go(func() error {
			var (
				tmpRes map[int64]*artmdl.Meta
				artErr error
				arg    *artmdl.ArgAids
			)
			arg = &artmdl.ArgAids{Aids: partAids, RealIP: ip}
			if tmpRes, artErr = s.artRPC.ArticleMetas(errCtx, arg); artErr != nil {
				log.Error("s.artRPC.ArticleMetas() error(%v)", artErr)
				return nil
			}
			mutex.Lock()
			for aid, arc := range tmpRes {
				if arc.AttrVal(artmdl.AttrBitNoDistribute) {
					continue
				}
				res[aid] = arc
			}
			mutex.Unlock()
			return nil
		})
	}
	group.Wait()
	return
}
