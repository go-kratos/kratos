package service

import (
	"context"

	"go-common/app/job/main/archive/model/archive"
	"go-common/app/job/main/archive/model/result"
	"go-common/app/job/main/archive/model/retry"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

func (s *Service) updateSubjectMid(nw *result.Archive, old *result.Archive) {
	if nw.Mid == old.Mid {
		return
	}
	if err := s.replyDao.ChangeSubjectMid(nw.AID, nw.Mid); err == nil {
		log.Info("/x/v2/reply/admin/subject/mid called")
	}
}

// sendMail send e-mail after state change
func (s *Service) sendMail(nw *result.Archive, old *result.Archive) {
	// nw.Attribute == old.Attribute tmp remove
	if nw.Mid == old.Mid && nw.TypeID == old.TypeID && nw.Duration == old.Duration && nw.Title == old.Title &&
		nw.Cover == old.Cover && nw.Content == old.Content && nw.Copyright == old.Copyright &&
		nw.State == old.State && nw.Access == old.Access && nw.Forward == old.Forward && nw.PubTime == old.PubTime && nw.CTime == old.CTime { // all field
		log.Info("archive(%d) change only mtime no email", nw.ID)
		return
	}
	a, err := s.arcServices[0].Archive3(context.TODO(), &arcmdl.ArgAid2{Aid: nw.AID})
	if err != nil || a == nil {
		log.Error("s.arcRPC.Archive3(%d) error(%v) or arc is nil", nw.AID, err)
		return
	}
	// not bangumi or movie
	if a.AttrVal(arcmdl.AttrBitIsBangumi) == arcmdl.AttrNo && a.AttrVal(arcmdl.AttrBitIsMovie) == arcmdl.AttrNo {
		return
	}
	s.emailDao.PGCNotifyMail(a, nw, old)
	log.Info("pgc notify(%d) mail send success", nw.AID)
}

func (s *Service) updateResultCache(nw *result.Archive, old *result.Archive) (err error) {
	var (
		c      = context.TODO()
		oldMid int64
	)
	if old != nil && old.Mid != nw.Mid {
		oldMid = old.Mid
	}
	for k, rpc := range s.arcServices {
		if err = rpc.ArcCache2(c, &arcmdl.ArgCache2{Aid: nw.AID, Tp: arcmdl.CacheUpdate}); err != nil {
			log.Error("s.arcRPC(%d).ArcCache2(%d, %s) error(%v)", k, nw.AID, arcmdl.CacheUpdate, err)
		}
		if nw.State >= 0 || nw.State == -6 {
			if err = rpc.ArcCache2(c, &arcmdl.ArgCache2{Aid: nw.AID, OldMid: oldMid, Tp: arcmdl.CacheAdd}); err != nil {
				log.Error("s.arcRpc(%d).ArcCache2(%d, %s) error(%v)", k, nw.AID, arcmdl.CacheAdd, err)
			}
		} else {
			if err = rpc.ArcCache2(c, &arcmdl.ArgCache2{Aid: nw.AID, OldMid: oldMid, Tp: arcmdl.CacheDelete}); err != nil {
				log.Error("s.arcRpc(%d).ArcCache2(%d, %s) error(%v)", k, nw.AID, arcmdl.CacheDelete, err)
			}
		}
		if err != nil {
			rt := &retry.Info{Action: retry.FailUpCache}
			rt.Data.Aid = nw.AID
			rt.Data.State = nw.State
			s.PushFail(c, rt)
			log.Error("updateResultCache error(%v)", err)
		}
	}
	return
}

func (s *Service) updateResultField(nw *result.Archive, old *result.Archive) {
	var (
		c   = context.TODO()
		err error
	)
	if nw.TypeID != old.TypeID {
		for k, rpc := range s.arcServices {
			if err = rpc.ArcFieldCache2(c, &arcmdl.ArgFieldCache2{Aid: nw.AID, TypeID: nw.TypeID, OldTypeID: old.TypeID}); err != nil {
				log.Error("s.arcRPC(%d).ArcFieldCache2(%d, %d, %d) error(%v)", k, nw.AID, nw.TypeID, old.TypeID, err)
			}
		}
	}
}

func (s *Service) tranResult(c context.Context, aid int64) (changed bool, upCids []int64, delCids []int64, err error) {
	var (
		tx        *sql.Tx
		rows      int64
		a         *archive.Archive
		aResult   *result.Archive
		ad        *archive.Addit
		vs        []*archive.Video
		videosCnt int
		staff     []*archive.Staff
	)
	defer func() {
		if err != nil {
			rt := &retry.Info{Action: retry.FailResultAdd}
			rt.Data.Aid = aid
			s.PushFail(c, rt)
			log.Error("tranResult error(%v)", err)
		}
	}()
	if a, err = s.archiveDao.Archive(c, aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v)", aid, err)
		return
	}
	if !a.IsSyncState() {
		log.Info("archive(%d) state(%d) cant change", aid, a.State)
		// FIXME: eeeee
		if s.isPGC(aid) && !s.hadPassed(c, aid) {
			log.Info("archive(%d) is PGC first change", aid)
		} else {
			return
		}
	}
	if ad, err = s.archiveDao.Addit(c, aid); err != nil {
		log.Error("s.arc.Addit(%d) error(%v)", aid, err)
		return
	}
	if ad == nil {
		ad = &archive.Addit{Aid: aid}
	}
	// if aid%10 == 0 || aid%10 == 1 || aid%10 == 2 {
	if vs, err = s.archiveDao.Videos2(c, aid); err != nil {
		log.Error("s.arc.Videos2(%d) error(%v)", aid, err)
		return
	}
	// } else {
	// if vs, err = s.archiveDao.Videos(c, aid); err != nil {
	// 	log.Error("s.arc.Videos(%d) error(%v)", aid, err)
	// 	return
	// }
	// }
	for _, v := range vs {
		if (v.Status == archive.VideoStatusAccess || v.Status == archive.VideoStatusOpen) && v.State >= 0 {
			videosCnt++
		}
	}
	// 辣鸡dede, check cid
	for _, v := range vs {
		if v.Cid == 0 && v.Status == archive.VideoStatusSubmit {
			// NOTE: 刚上传，没必要同步去
			log.Error("aid(%d) vid(%d) cid(%d) videoStatus(%d) return", v.Aid, v.ID, v.Cid, v.Status)
			return
		}
	}
	if aResult, err = s.resultDao.Archive(c, aid); err != nil {
		log.Error("s.resultDao.Archive error(%+v)", err)
		return
	}
	if tx, err = s.resultDao.BeginTran(c); err != nil {
		log.Error("s.result.BeginTran error(%v)", err)
		return
	}
	var (
		duration   int
		firstCid   int64
		dimensions string
	)
	for _, v := range vs {
		if (v.Status == archive.VideoStatusAccess || v.Status == archive.VideoStatusOpen) && v.State == archive.VideoRelationBind {
			if _, err = s.resultDao.TxAddVideo(c, tx, v); err != nil {
				tx.Rollback()
				log.Error("s.result.TxAddVideo error(%v)", err)
				break
			}
			duration += int(v.Duration)
			upCids = append(upCids, v.Cid)
			if v.Index == 1 && v.SrcType == "vupload" {
				firstCid = v.Cid
				dimensions = v.Dimensions
			}
		} else {
			if _, err = s.resultDao.TxDelVideoByCid(c, tx, aid, v.Cid); err != nil {
				tx.Rollback()
				log.Error("s.result.TxDelVideoByCid error(%v)")
				break
			}
			delCids = append(delCids, v.Cid)
		}
	}
	a.Duration = duration
	if rows, err = s.resultDao.TxAddArchive(c, tx, a, ad, videosCnt, firstCid, dimensions); err != nil {
		tx.Rollback()
		log.Error("s.result.TxAddArchive error(%v)", err)
		return
	}
	if rows == 0 {
		if _, err = s.resultDao.TxUpArchive(c, tx, a, ad, videosCnt, firstCid, dimensions); err != nil {
			tx.Rollback()
			log.Error("s.result.TxUpArchive error(%v)")
			return
		}
	}
	// 更新联合投稿人
	if a.AttrVal(archive.AttrBitIsCooperation) == archive.AttrYes {
		if staff, err = s.archiveDao.Staff(c, aid); err != nil {
			tx.Rollback()
			log.Error("s.archiveDao.Staff aid(%d) error(%v)", aid, err)
			return
		}
		if err = s.resultDao.TxDelStaff(c, tx, aid); err != nil {
			tx.Rollback()
			log.Error("s.result.TxDelStaff aid(%d) error(%v)", aid, err)
			return
		}
		if staff != nil {
			if err = s.resultDao.TxAddStaff(c, tx, aid, staff); err != nil {
				tx.Rollback()
				log.Error("s.result.TxAddStaff aid(%d) error(%v)", aid, err)
				return
			}
		}
	} else { //从联合投稿改为非联合投稿的 删除staff数据
		if aResult != nil && aResult.AttrVal(archive.AttrBitIsCooperation) == archive.AttrYes {
			if err = s.resultDao.TxDelStaff(c, tx, aid); err != nil {
				tx.Rollback()
				log.Error("s.result.TxDelStaff aid(%d) error(%v)", aid, err)
				return
			}
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)")
		return
	}
	log.Info("aid(%d) upCids(%d) delCids(%d) db updated", aid, len(upCids), len(delCids))
	changed = true
	return
}
