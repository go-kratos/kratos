package service

import (
	"context"
	"time"

	"go-common/app/job/main/dm2/model"
	"go-common/library/log"
)

func (s *Service) transferProc() {
	var (
		c        = context.TODO()
		interval = time.Duration(time.Second * 60)
	)
	for {
		time.Sleep(interval)
		if !s.dao.AddTransferLock(c) {
			continue
		}
		trans, err := s.dao.Transfers(c, model.StatInit)
		if err != nil || len(trans) == 0 {
			continue
		}
		for _, t := range trans {
			log.Info("dm transfer(%+v) start", t)
			s.transfer(c, t)
		}
	}
}

// transfer transfer dm.
func (s *Service) transfer(c context.Context, t *model.Transfer) {
	var (
		err     error
		limit   int64 = 500
		startID       = t.Dmid
		tp            = model.SubTypeVideo
	)
	t.State = model.StatTransfing
	if _, err = s.dao.UpdateTransfer(c, t); err != nil {
		log.Error("s.dao.UpdateTransfer(%+v) error(%v)", t, err)
		return
	}
	if err = s.dao.DelTransferLock(c); err != nil {
		log.Error("s.dao.DelTransferLock() error")
	}
	targetSub, err := s.dao.Subject(c, tp, t.ToCid)
	if err != nil || targetSub == nil {
		log.Error("s.dao.Subject(cid:%d) error(%v)", t.ToCid, err)
		s.transerFailNow(c, t)
		return
	}
	originSub, err := s.dao.Subject(c, tp, t.FromCid)
	if err != nil || originSub == nil {
		log.Error("s.dao.Subject(cid:%d) error(%v)", t.ToCid, err)
		s.transerFailNow(c, t)
		return
	}
	for {
		// get transfer dm per page
		var dms []*model.DM
		if dms, err = s.transferDMS(c, tp, originSub.Oid, startID, limit); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		if len(dms) == 0 {
			break
		}
		for _, dm := range dms {
			if dm.ID <= startID {
				continue
			} else {
				startID = dm.ID
			}
			var id int64
			if id, err = s.seqRPC.ID(c, s.seqArg); err != nil {
				log.Error("seqRPC.ID() error(%v)", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if dm.Pool == model.PoolSpecial {
				dm.ContentSpe.ID = id
			}
			dm.Oid = targetSub.Oid // 修改这个dm 的主键id和oid
			dm.ID = id
			dm.Content.ID = id
			if t.Offset != 0 {
				dm.Progress = dm.Progress + int32(t.Offset*1000)
			}
			if err = s.actionAddDM(c, targetSub, dm); err != nil {
				continue
			}
			t.Dmid = startID //记录转移到的dmid
		}
		s.dao.UpdateTransfer(c, t)
		time.Sleep(1 * time.Second)
	}
	t.State = model.StatFinished
	if _, err = s.dao.UpdateTransfer(c, t); err != nil {
		log.Error("s.dao.UpdateTransfer(%+v) error(%v)", t, err)
	}
	// 刷新弹幕缓存
	s.flushDmCache(c, &model.Flush{Oid: t.ToCid, Type: tp, Force: true})
	s.flushAllDmSegCache(c, t.ToCid, tp)
}

func (s *Service) transerFailNow(c context.Context, t *model.Transfer) {
	t.State = model.StatFailed
	if _, err := s.dao.UpdateTransfer(c, t); err != nil {
		log.Error("s.dao.UpdateTransfer(%+v) error(%v)", t, err)
	}
}

// NewCommentList get dm list from new db
func (s *Service) transferDMS(c context.Context, tp int32, oid, minID, limit int64) (dms []*model.DM, err error) {
	contentSpec := make(map[int64]*model.ContentSpecial)
	idxMap, dmids, special, err := s.dao.DMIndexs(c, tp, oid, minID, limit)
	if err != nil {
		log.Error("s.dao.DMIndexs(oid:%d mindID:%d) error(%v)", oid, minID, err)
		return
	}
	if len(dmids) == 0 {
		return
	}
	contents, err := s.dao.Contents(c, oid, dmids)
	if err != nil {
		log.Error("s.dao.Contents(oid:%d dmids:%v) error(%v)", oid, dmids, err)
		return
	}
	if len(special) > 0 {
		if contentSpec, err = s.dao.ContentsSpecial(c, special); err != nil {
			log.Error("s.dao.ContentSpecials(oid:%d special:%v) error(%v)", oid, special, err)
			return
		}
	}
	for _, dmid := range dmids {
		dm, ok := idxMap[dmid]
		if !ok {
			continue
		}
		content, ok := contents[dmid]
		if !ok {
			continue
		}
		dm.Content = content
		if dm.Pool == model.PoolSpecial {
			contentspe, ok := contentSpec[dm.ID]
			if !ok {
				continue
			}
			dm.ContentSpe = contentspe
		}
		dms = append(dms, dm)
	}
	return
}
