package ugc

import (
	"time"

	appDao "go-common/app/job/main/tv/dao/app"
	"go-common/app/job/main/tv/dao/lic"
	"go-common/library/database/sql"
	"go-common/library/log"
)

func (s *Service) delArcproc() {
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("delArcproc DB closed!")
			return
		}
		// build the skeleton, arc + video data
		cAid, err := s.dao.DeletedArc(ctx)
		if err != nil && err != sql.ErrNoRows {
			log.Error("DeletedArc Error %v", err)
			appDao.PromError("SyncDelAid:Err")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		if err == sql.ErrNoRows || cAid == 0 {
			log.Info("SyncDelAid No Data to Sync")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		if err = s.delLic(cAid); err != nil {
			appDao.PromError("SyncDelAid:Err")
			log.Error("delLic error %v, aid %d", err, cAid)
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		appDao.PromInfo("SyncDelAid:Succ")
	}
}

// delArcErr: it logs the error and postpone the videos for the next submit
func (s *Service) delArcErr(aid int64, fmt string, err error) {
	s.dao.PpDelArc(ctx, aid)
	log.Error(fmt, aid, err)
}

// delLic: sync our arc data to License owner
func (s *Service) delLic(cAid int64) (err error) {
	var (
		xmlBody string
		sign    = s.c.Sync.Sign
		prefix  = s.c.Sync.UGCPrefix
	)
	xmlBody = lic.PrepareXML(lic.DelLic(sign, prefix, cAid))
	// call api
	if _, err = s.licDao.CallRetry(ctx, s.c.Sync.API.DelSeasonURL, xmlBody); err != nil {
		s.delArcErr(cAid, "xml call %d error %v", err)
		return
	}
	// update the arc & videos' submit status to finish
	if err = s.dao.FinishDelArc(ctx, cAid); err != nil {
		s.delArcErr(cAid, "FinishDelArc %d error %v", err)
	}
	return
}

func (s *Service) delArc(aid int64) (err error) {
	var tx *sql.Tx
	// check whether the arc exist in our DB
	if !s.arcExist(aid) {
		log.Warn("Del Arc %d, it doesn't exist", aid)
		return
	}
	// delete the arc, put submit to 1
	if tx, err = s.dao.BeginTran(ctx); err != nil { // begin transaction
		return
	}
	if err = s.dao.TxDelArc(tx, aid); err != nil {
		appDao.PromError("DelArc:Err")
		tx.Rollback()
		return
	}
	// delete the videos put submit to 1
	if err = s.dao.TxDelVideos(tx, aid); err != nil {
		appDao.PromError("DelArc:Err")
		tx.Rollback()
		return
	}
	appDao.PromInfo("DelArc:Succ")
	tx.Commit()
	return
}
