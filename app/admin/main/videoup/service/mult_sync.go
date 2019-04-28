package service

import (
	"context"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
)

func (s *Service) multSyncProc() {
	defer s.wg.Done()
	for {
		if s.closed {
			return
		}
		var (
			c    = context.TODO()
			err  error
			sync *archive.MultSyncParam
		)
		if sync, err = s.busCache.PopMultSync(c); err != nil || sync == nil || sync.Action == "" {
			time.Sleep(5 * time.Second)
			continue
		}
		log.Info("sync_action %s %+v", sync.Action, sync)
		switch sync.Action {
		case archive.ActionVideoSubmit:
			if err = s.dealVideo(c, sync.VideoParam); err != nil {
				s.busCache.PushMultSync(c, sync)
				log.Error("s.dealVideo() error(%v)", err)
				continue
			}
		case archive.ActionArchiveSubmit:
			if err = s.dealArchive(c, sync.ArcParam); err != nil {
				s.busCache.PushMultSync(c, sync)
				log.Error("s.dealArchive() error(%v)", err)
				continue
			}
		case archive.ActionArchiveSecondRound:
			if err = s.dealArchiveSecondRound(c, sync.ArcParam); err != nil {
				s.busCache.PushMultSync(c, sync)
				log.Error("s.dealArchive() error(%v)", err)
				continue
			}
		case archive.ActionArchiveAttr:
			if err = s.dealAttrs(c, sync.ArcParam); err != nil {
				s.busCache.PushMultSync(c, sync)
				log.Error("s.dealAttrs() error(%v)", err)
				continue
			}
		case archive.ActionArchiveTypeID:
			if err = s.dealTypeID(c, sync.ArcParam); err != nil {
				s.busCache.PushMultSync(c, sync)
				log.Error("s.dealTypeID() error(%v)", err)
				continue
			}
		case archive.ActionArchiveTag:
			if err = s.dealTag(c, false, sync.ArcParam); err != nil {
				s.busCache.PushMultSync(c, sync)
				log.Error("s.dealTag() error(%v)", err)
				continue
			}
		case archive.ActionArchiveTagRecheck:
			if err = s.dealTag(c, true, sync.ArcParam); err != nil {
				s.busCache.PushMultSync(c, sync)
				log.Error("s.dealTag() error(%v)", err)
				continue
			}
		default:
			log.Info("s.multSyncProc() default action(%s)", sync.Action)
		}
	}
}
