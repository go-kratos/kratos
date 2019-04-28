package ugc

import (
	"context"
	"time"

	appDao "go-common/app/job/main/tv/dao/app"
	ugcMdl "go-common/app/job/main/tv/model/ugc"
	account "go-common/app/service/main/account/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_arcPiece = 20
	_upName   = 1
	_upFace   = 2
)

// refreshUpproc refreshes the upper info regularly
func (s *Service) refreshUpproc() {
	var c = context.Background()
	for {
		time.Sleep(time.Duration(s.c.UgcSync.Frequency.UpperRefresh))
		s.refreshUp(c, true)
		log.Info("RefreshUpper Succ")
	}
}

// refreshUp picks all the upper in our DB and compare them with RPC result to refresh if different
func (s *Service) refreshUp(ctx context.Context, refresh bool) {
	var (
		count     int64
		nbPiece   int
		err       error
		maxID     int64
		activeUps = make(map[int64]int)
	)
	if count, err = s.upDao.CountUP(ctx); err != nil {
		log.Error("[refreshUpper] CountUP error [%v]", err)
		return
	}
	nbPiece = appDao.NumPce(int(count), _arcPiece)
	log.Info("[refreshUpper] Numero Piece %d", nbPiece)
	for i := 0; i < nbPiece; i++ {
		res, newMaxID, err := s.upDao.PickUppers(ctx, maxID, _arcPiece)
		if err != nil {
			log.Error("[refreshUpper] Pick Piece %d Error, Ignore it", i)
			continue
		}
		if newMaxID <= maxID {
			log.Error("MaxID is not increasing! [%d,%d]", newMaxID, maxID)
			return
		}
		maxID = newMaxID
		for _, v := range res { // travel the 50 mids
			activeUps[v] = 1
			if refresh {
				s.upCheck(ctx, v)
			}
		}
	}
	if len(activeUps) > 0 {
		s.activeUps = activeUps
		log.Info("[refreshUpper] ActiveUps Len %d", len(activeUps))
	}
}

// upCheck checks the upper remote and local, modify if different
func (s *Service) upCheck(ctx context.Context, v int64) (err error) {
	var (
		upCMS *ugcMdl.Upper
		upRPC *account.Card
	)
	if upCMS, err = s.upDao.LoadUpMeta(ctx, v); err != nil { // load local upper data
		log.Warn("[refreshUpper] LoadUpMeta Mid %d, Err %d", v, err)
		return
	}
	if upRPC, err = s.upDao.Card3(ctx, v); err != nil { // load remote upper data
		log.Warn("[refreshUpper] Card3 Mid %d, Err %d", v, err)
		return
	}
	if err = s.upModify(ctx, upRPC, upCMS); err != nil {
		log.Warn("[refreshUpper] upModify Mid %d, Err %d", v, err)
	}
	return
}

// upModify checks whether the upper's info has been modified or not, if yes, refresh it in DB and MC
func (s *Service) upModify(ctx context.Context, upRPC *account.Card, upCMS *ugcMdl.Upper) (err error) {
	fs, ns := upCMS.IsSame(upRPC.Name, upRPC.Face)
	if fs && ns { // compare, if same, just jump to the next upper
		return
	}
	if !fs {
		log.Info("Mid %d Face Modified, Old: %s, New %s", upRPC.Mid, upCMS.OriFace, upRPC.Face)
		req := &ugcMdl.ReqSetUp{
			Value:  upRPC.Face,
			MID:    upRPC.Mid,
			UpType: _upFace,
		}
		if err = s.upDao.RefreshUp(ctx, req); err != nil {
			log.Error("RefreshUp Req %v, Err %v", req, err)
			return
		}
	}
	if !ns {
		log.Info("Mid %d Name Modified, Old: %s, New %s", upRPC.Mid, upCMS.OriName, upRPC.Name)
		req := &ugcMdl.ReqSetUp{
			Value:  upRPC.Name,
			MID:    upRPC.Mid,
			UpType: _upName,
		}
		if err = s.upDao.RefreshUp(ctx, req); err != nil {
			log.Error("RefreshUp Req %v, Err %v", req, err)
			return
		}
	}
	err = s.upDao.SendUpper(ctx, upRPC.Mid) // only update once if both face and name are modified
	return
}

// delete Upper treatment proc
func (s *Service) delUpproc() {
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("delUpproc DB closed!")
			return
		}
		// build the skeleton, arc + video data
		cMid, err := s.dao.DeletedUp(ctx)
		if err != nil && err != sql.ErrNoRows {
			log.Error("DeletedUp Error %v", err)
			appDao.PromError("DelUp:Err")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		if err == sql.ErrNoRows || cMid == 0 {
			log.Info("SyncDelAid No Data to Sync")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		if err = s.delUp(cMid); err != nil {
			appDao.PromError("DelUp:Err")
			log.Error("delUp error %v, mid %d", err, cMid)
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		appDao.PromInfo("DelUp:Succ")
	}
}

// delete Upper treatment
func (s *Service) delUp(mid int64) (err error) {
	var (
		aids  []int64
		count int64
		page  int
	)
	if count, err = s.dao.CountUpArcs(ctx, mid); err != nil {
		s.delUpErr(mid, "CountUpArcs Mid %d, Err %v", err)
		return
	}
	page = appDao.NumPce(int(count), _arcPiece)
	log.Info("Ready to Remove Mid %d Arcs, Count: %d, Page: %d", mid, count, page)
	for i := 0; i < page; i++ {
		log.Info("Operating DelUp Mid %d Page %d", mid, i+1)
		if aids, err = s.dao.UpArcs(ctx, mid); err != nil && err != sql.ErrNoRows {
			return
		}
		if err == sql.ErrNoRows || len(aids) == 0 { // means end of the loup
			break
		}
		for _, v := range aids {
			if err = s.delArc(v); err != nil {
				s.delUpErr(mid, "delArc Mid %d, Err %v", err)
				return
			}
		}
	}
	// change the upper's status to tell it's finish
	if err = s.dao.FinishDelUp(ctx, mid); err != nil {
		s.delUpErr(mid, "delArc Mid %d, Err %v", err)
		return
	}
	log.Info("DelUp Mid %d Succ, Count %d, Page %d", mid, count, page)
	return
}

// delUpErr: it logs the error and postpone the videos for the next submit
func (s *Service) delUpErr(mid int64, fmt string, err error) {
	s.dao.PpDelUp(ctx, mid)
	log.Error(fmt, mid, err)
}

// syncUpproc picks all the submit=1 uppers and updates their archives' submit to 1
func (s *Service) syncUpproc() {
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("syncUpproc DB closed!")
			return
		}
		mids, err := s.upDao.TosyncUps(ctx) // pick to sync mids
		if err != nil {
			log.Error("syncUpproc Error %v", err)
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		if len(mids) == 0 {
			log.Info("No Upper to Sync")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.SyncFre))
			continue
		}
		log.Info("syncUpproc treats %d uppers", len(mids))
		for _, v := range mids { // updates these uppers' archives' submit to 1
			if err = s.submitUps(ctx, v); err != nil {
				log.Error("syncUpproc updates Archive error %v, mid %d", err, v)
				continue
			}
			if err = s.upDao.FinsyncUps(ctx, v); err != nil {
				log.Error("syncUpproc finish upper error %v, mid %d", err, v)
			}
		}
	}
}

func (s *Service) submitUps(ctx context.Context, mid int64) (err error) {
	var (
		count int64
	)
	if count, err = s.dao.CountUpArcs(ctx, mid); err != nil {
		s.delUpErr(mid, "CountUpArcs Mid %d, Err %v", err)
		return
	}
	page := appDao.NumPce(int(count), _arcPiece)
	log.Info("Ready to Submit Mid %d Arcs, Count: %d, Page: %d", mid, count, page)
	for i := 0; i < page; i++ {
		log.Info("Operating DelUp Mid %d Page %d", mid, i+1)
		var aids []int64
		if aids, err = s.dao.UpArcs(ctx, mid); err != nil && err != sql.ErrNoRows {
			return
		}
		if err == sql.ErrNoRows || len(aids) == 0 { // means end of the loup
			break
		}
		s.modArcCh <- aids // put them into the channel
		time.Sleep(50 * time.Millisecond)
	}
	return
}
