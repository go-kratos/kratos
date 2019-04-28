package ugc

import (
	"context"
	"strings"
	"time"

	appDao "go-common/app/job/main/tv/dao/app"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	arccli "go-common/app/service/main/archive/api"
	"go-common/library/database/sql"
	"go-common/library/log"
)

func (s *Service) manualproc() {
	var (
		err  error
		arcs []*ugcmdl.Archive
	)
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("manualproc DB closed!")
			return
		}
		if arcs, err = s.dao.Manual(ctx); err != nil && err != sql.ErrNoRows {
			log.Error("manualproc Error %v", err)
			appDao.PromError("Manual:Err")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.ManualFre))
			continue
		}
		if err == sql.ErrNoRows || len(arcs) == 0 {
			log.Info("No Manual Data")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.ManualFre))
			continue
		}
		if err = s.manual(arcs); err != nil {
			log.Error("manualproc Error %v", err)
			appDao.PromError("Manual:Err")
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.ManualFre))
			continue
		}
		appDao.PromInfo("Manual:Succ")
		time.Sleep(1 * time.Second)
	}
}

func (s *Service) manual(arcs []*ugcmdl.Archive) (err error) {
	for _, v := range arcs { // locate each archive
		if err = s.importArc(context.Background(), v.AID, true); err != nil {
			log.Error("importArc Error %v", err)
			s.dao.Ppmnl(ctx, v.AID) // postpone the next retry
			continue
		}
	}
	return
}

//coverURL is used for completing url start with bfs,example: /bfs/archive/diuren.png
func (s *Service) coverURL(uri string, prefix string) string {
	if strings.HasPrefix(uri, "/bfs") {
		return prefix + uri
	}
	return uri
}

// importArc imports an brand new archive data, isManual means whether we need update the arc's status
func (s *Service) importArc(ctx context.Context, aid int64, isManual bool) (err error) {
	var (
		tx       *sql.Tx
		arcGrpc  *arccli.Arc
		view     *arccli.ViewReply
		arcAllow = &ugcmdl.ArcAllow{}
		arc      = &ugcmdl.Archive{}
	)
	if arcGrpc, err = s.arcPick(ctx, aid); err != nil { // pick archive api data
		return
	}
	arcAllow.FromArcReply(arcGrpc)
	if allow := s.arcAllowImport(arcAllow); !allow { // check whether the archive is allowed to import into TV db
		if isManual {
			err = s.delArc(aid)
		}
		return
	}
	arc.FromArcReply(arcGrpc)
	arc.Cover = s.coverURL(arc.Cover, s.c.UgcSync.Cfg.BFSPrefix)
	if view, err = s.videoPick(ctx, aid); err != nil { // pick video api data
		return
	}
	if tx, err = s.dao.BeginTran(ctx); err != nil { // begin transaction
		return
	}
	if isManual {
		if err = s.dao.TxMnlArc(tx, arc); err != nil { // manual import archive data, update
			tx.Rollback()
			return
		}
		if arc.MID != 0 {
			s.manualUp(ctx, arc.MID)
		}
	} else {
		if err = s.dao.TxAutoArc(tx, arc); err != nil { // databus import archive data, insert
			tx.Rollback()
			return
		}
	}
	if err = s.dao.TxMnlVideos(tx, view); err != nil { // import video data
		tx.Rollback()
		return
	}
	if isManual {
		if err = s.dao.TxMnlStatus(tx, aid); err != nil { // update the manual to 0, finish the operation
			tx.Rollback()
			return
		}
	}
	log.Info("ImportArc Aid %d Succ", aid)
	tx.Commit()
	return
}

// manualUp imports the manual submit archive's upper
func (s *Service) manualUp(ctx context.Context, mid int64) {
	if cmid, _ := s.dao.UpInList(ctx, mid); cmid != 0 { // mid already in list, no need to import
		return
	}
	upRPC, err := s.upDao.Card3(ctx, mid)
	if err != nil { // load remote upper data
		log.Warn("[manualUp] Card3 Mid %d, Err %d", mid, err)
		return
	}
	s.upDao.ImportUp(ctx, &ugcmdl.EasyUp{ // import this upper into DB and cache
		MID:  mid,
		Face: upRPC.Face,
		Name: upRPC.Name,
	})
}
