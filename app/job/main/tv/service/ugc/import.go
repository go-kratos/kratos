package ugc

import (
	"fmt"
	"time"

	appDao "go-common/app/job/main/tv/dao/app"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	arccli "go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_arcRetry = 3
	_apiRetry = 5
	_sleep    = 100 * time.Millisecond
)

// upImportproc always runs to init the uppers
func (s *Service) upImportproc() {
	var (
		err    error
		uppers []*ugcmdl.Upper
	)
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("upImportproc DB closed!")
			return
		}
		// if no more data, we scan per 30s
		if uppers, err = s.dao.Import(ctx); err != nil && err != sql.ErrNoRows {
			log.Error("upperImport error %v", err)
			appDao.PromError("ImportMid:Err")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.ImportFre))
			continue
		}
		if len(uppers) == 0 && err == sql.ErrNoRows {
			log.Info("No Import Data")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.ImportFre))
			continue
		}
		if err = s.upImport(uppers); err != nil {
			log.Error("upImport Error %v", err)
			appDao.PromError("ImportMid:Err")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.ImportFre))
			continue
		}
		appDao.PromInfo("ImportMid:Succ")
		time.Sleep(1 * time.Second)
	}
}

// upImport loads 20 uppers to init, and load them one by one
func (s *Service) upImport(uppers []*ugcmdl.Upper) (err error) {
	for _, v := range uppers {
		// import data
		if err = s.InitUpper(v.MID); err != nil {
			log.Error("initUpper MID: %v, Err: %v, Postpone the MID", v.MID, err)
			s.dao.PpUpper(ctx, v.MID)
			continue
		}
		// update the status
		if err = s.dao.FinishUpper(ctx, v.MID); err != nil {
			log.Error("FinishUpper Mid: %d, Err: %v", v.MID, err)
			return
		}
		time.Sleep(time.Duration(s.c.UgcSync.Frequency.UpperPause)) // pause after import each upper
	}
	return
}

// InitUpper takes the upper's archive & videos, load them into our DB
func (s *Service) InitUpper(mid int64) (err error) {
	var (
		arcCount int
		ps       = s.c.UgcSync.Batch.ArcPS // page size to pick archives
		ptn      int                       // total page number
		pMatch   map[int64]*arccli.Arc     // the mapping of aid to archive model of one page
		pAids    []int64                   // the aids of one page
		videoNum int64
		begin    = time.Now()
	)
	// count upper's archive and get the total number of pages to get
	if arcCount, err = s.arcCount(mid); err != nil {
		return
	}
	log.Info("InitUpper mid %d, Count: %d", mid, arcCount)
	if arcCount == 0 {
		log.Error("Upper %d Arc Count is 0", mid)
		return
	}
	if arcCount%ps == 0 {
		ptn = arcCount / ps
	} else {
		ptn = arcCount/ps + 1
	}
	// get the upper's archives page by page
	for i := 1; i <= ptn; i++ {
		if pMatch, pAids, err = s.UpArchives(mid, i, ps); err != nil {
			log.Error("Mid %d, Page %d Error %v", mid, i, err)
			return
		}
		if len(pMatch) == 0 { // which means this page is all existing
			log.Error("Mid %d, Page %d, no need to import Due to Types Hit", mid, i)
			continue
		}
		if err = s.dao.FilterExist(ctx, &pMatch, pAids); err != nil { // filter the existing ones
			log.Error("Mid %d, Page %d Error %v", mid, i, err)
			return
		}
		if len(pMatch) == 0 { // which means this page is all existing
			log.Error("Mid %d, Page %d, no need to impot Due to Existing", mid, i)
			continue
		}
		if err = s.arcsIn(pMatch); err != nil { // insert this page's arc & views data into our DB
			log.Error("Mid %d, Page %d Error %v", mid, i, err)
			return
		}
		videoNum = videoNum + int64(len(pMatch))
		time.Sleep(time.Duration(s.c.UgcSync.Frequency.UpInitFre)) // pause after import each page of upper's archive
	}
	log.Info("ImportUpper Mid %d, Page Number %d, Page Size %d, "+
		"Video Number %d, Time %v", mid, ptn, ps, videoNum, time.Since(begin)) // record init upper time
	return
}

// get map's keys
func mapKeys(myMap map[int64]*arccli.Arc) (keys []int64) {
	for k := range myMap {
		keys = append(keys, k)
	}
	return
}

// UpArchives picks one page of the up's archives
func (s *Service) UpArchives(mid int64, pn int, ps int) (match map[int64]*arccli.Arc, aids []int64, err error) {
	var res []*arccli.Arc
	match = make(map[int64]*arccli.Arc)
	if err = Retry(func() (err error) {
		if res, err = s.arcRPC.UpArcs3(ctx, &arcmdl.ArgUpArcs2{
			Mid: mid,
			Pn:  pn,
			Ps:  ps,
		}); err != nil {
			log.Error("%+v", err)
		}
		return
	}, _arcRetry, _sleep); err != nil {
		log.Error("upArchives Error %+v", err)
		return
	} else if len(res) == 0 {
		err = fmt.Errorf("result empty")
		return
	}
	for _, v := range res {
		arcAllow := &ugcmdl.ArcAllow{}
		arcAllow.FromArcmdl(v)
		if allow := s.arcAllowImport(arcAllow); !allow { // check whether the archive is allowed to import into TV db
			continue
		}
		match[v.Aid] = v
		aids = append(aids, v.Aid)
	}
	return
}

// Retry . retry one function until no error
func Retry(callback func() error, retry int, sleep time.Duration) (err error) {
	for i := 0; i < retry; i++ {
		if err = callback(); err == nil {
			return
		}
		time.Sleep(sleep)
	}
	return
}

// arcsIn picks one page of archive data and their views data, to import them into the DB one by one
func (s *Service) arcsIn(pMatch map[int64]*arccli.Arc) (err error) {
	var (
		tx     *sql.Tx
		pViews map[int64]*arccli.ViewReply
		pAids  []int64
	)
	// get the filtered aids to get the views
	pAids = mapKeys(pMatch)
	if pViews, err = s.arcViews(pAids); err != nil {
		log.Error("arcsIn Error %v", err)
		return
	}
	// import the arc & its video one by one
	for aid, arc := range pMatch {
		// begin the transaction and insert the archive data
		if tx, err = s.dao.BeginTran(ctx); err != nil { // begin transaction
			return
		}
		arc.Pic = s.coverURL(arc.Pic, s.c.UgcSync.Cfg.BFSPrefix)
		if err = s.dao.TxImportArc(tx, arc); err != nil {
			tx.Rollback()
			return
		}
		cViews, ok := pViews[arc.Aid]
		if !ok {
			log.Error("arcIn View Data for %d not found", arc.Aid)
			tx.Rollback()
			return
		}
		if err = s.dao.TxMnlVideos(tx, cViews); err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
		log.Info("Succ Add Arc & View for Aid: %d", aid)
	}
	return
}
