package service

import (
	"context"
	"math"
	"sync"

	"go-common/app/interface/main/dm2/model"
	"go-common/app/service/main/archive/api"
	arcMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

func (s *Service) archiveInfos(c context.Context, aids []int64) (archiveInfos map[int64]*api.Arc, err error) {
	var (
		pagesize = 100
		wg       errgroup.Group
		mu       sync.Mutex
	)
	archiveInfos = make(map[int64]*api.Arc)
	if len(aids) <= 0 {
		return
	}
	page := int(math.Ceil(float64(len(aids)) / float64(pagesize)))
	for i := 0; i < page; i++ {
		start := i * pagesize
		end := (i + 1) * pagesize
		if end > len(aids) {
			end = len(aids)
		}
		wg.Go(func() (err error) {
			arg := &arcMdl.ArgAids2{Aids: aids[start:end]}
			infos, err := s.arcRPC.Archives3(c, arg)
			if err != nil {
				log.Error("s.arcRPC.Archives3(%v) error(%v)", arg, err)
				return
			}
			for _, info := range infos {
				mu.Lock()
				archiveInfos[info.Aid] = info
				mu.Unlock()
			}
			return
		})
	}
	err = wg.Wait()
	return
}

// videoDuration return video duration cid.
func (s *Service) videoDuration(c context.Context, aid, cid int64) (duration int64, err error) {
	var cache = true
	if duration, err = s.dao.DurationCache(c, cid); err != nil {
		log.Error("dao.Duration(cid:%d) error(%v)", cid, err)
		err = nil
		cache = false
	} else if duration != model.NotFound {
		return
	}
	arg := &arcMdl.ArgVideo2{Aid: aid, Cid: cid, RealIP: metadata.String(c, metadata.RemoteIP)}
	page, err := s.arcRPC.Video3(c, arg)
	if err != nil {
		if ecode.Cause(err).Code() == ecode.NothingFound.Code() {
			duration = 0
			err = nil
			log.Warn("acvSvc.Video3(%v) error(duration not exist)", arg)
		} else {
			log.Error("acvSvc.Video3(%v) error(%v)", arg, err)
		}
	} else {
		duration = page.Duration * 1000
	}
	if cache {
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.SetDurationCache(ctx, cid, duration)
		})
	}
	return
}
