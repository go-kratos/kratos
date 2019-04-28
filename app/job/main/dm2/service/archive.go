package service

import (
	"context"

	"go-common/app/job/main/dm2/model"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

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
	arg := &archive.ArgVideo2{Aid: aid, Cid: cid, RealIP: metadata.String(c, metadata.RemoteIP)}
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
