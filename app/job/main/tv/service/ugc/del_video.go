package ugc

import (
	appDao "go-common/app/job/main/tv/dao/app"
	"go-common/app/job/main/tv/dao/lic"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

func (s *Service) delVideoproc() {
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("delVideoproc DB closed!")
			return
		}
		// pick deleted videos
		videoIDs, err := s.dao.DeletedVideos(ctx)
		if err != nil && err != sql.ErrNoRows {
			log.Error("videoIDs Error %v", err)
			appDao.PromError("SyncDelVid:Err")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		if err == sql.ErrNoRows || len(videoIDs) == 0 {
			log.Info("No SyncDelVid Data to Sync")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		if err = s.delVideoLic(videoIDs); err != nil {
			appDao.PromError("SyncDelVid:Err")
			log.Error("delLic error %v, cids %s", err, videoIDs)
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		appDao.PromInfo("SyncDelVid:Succ")
	}
}

// delVideoErr: it logs the error and postpone the videos for the next submit
func (s *Service) delVideoErr(cids []int, fmt string, err error) {
	s.dao.PpDelVideos(ctx, cids)
	log.Error(fmt, cids, err)
}

// delVideoLic: sync our deleted video data to License owner
func (s *Service) delVideoLic(videoIDs []int) (err error) {
	var (
		xmlBody string
		sign    = s.c.Sync.Sign
		prefix  = s.c.Sync.UGCPrefix
	)
	xmlBody = lic.DelEpLic(prefix, sign, videoIDs)
	// call api
	if _, err = s.licDao.CallRetry(ctx, s.c.Sync.API.DelEPURL, xmlBody); err != nil {
		s.delVideoErr(videoIDs, "xml call %v error %v", err)
		return
	}
	// update the videos' submit status to finish
	if err = s.dao.FinishDelVideos(ctx, videoIDs); err != nil {
		log.Info("Del Video Finish, Sync For Vids: %v", videoIDs)
		s.delVideoErr(videoIDs, "FinishDelVideos %v error %v", err)
	}
	return
}
