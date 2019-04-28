package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/dm/model"
	dm2Mdl "go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/model/oplog"
	account "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_recallLt  = 5                 //普通用户每天可以撤回几个
	_recallTO  = 120               //多久以前的弹幕可以撤回
	_recallCap = 20000             // 字幕君
	_recallOK  = "撤回成功，你还有%s次撤回机会" // 弹幕撤回成功提示
)

// Recall 撤回弹幕
func (s *Service) Recall(c context.Context, mid, cid int64, id int64) (msg string, err error) {
	var (
		card    *account.CardReply
		cnt     int
		ts      string
		isSuper bool
	)
	if card, err = s.accountSvc.Card3(c, &account.MidReq{Mid: mid}); err != nil {
		log.Error("s.actSvc.Card3(%d) error(%v)", mid, err)
		return
	}
	ts = "无限"
	isSuper = card.GetCard().GetRank() >= _recallCap
	if !isSuper {
		if cnt, err = s.dao.RecallCnt(c, mid); err != nil {
			log.Error("s.dao.RecallCnt(%d) error(%v)", mid, err)
			return
		}
		if cnt >= _recallLt {
			err = ecode.DMRecallLimit
			return
		}
		ts = strconv.Itoa(_recallLt - cnt - 1)
	}
	dm, err := s.dao.Index(c, model.SubTypeVideo, cid, id)
	if err != nil {
		return
	}
	if dm == nil || !dm.NeedDisplay() || dm.Mid != mid {
		err = ecode.DMRecallDeleted
		return
	}
	if (time.Now().Unix()-int64(dm.Ctime)) > _recallTO && !(isSuper && (dm.Pool == 1 || dm.Pool == 2)) {
		err = ecode.DMRecallTimeout
		return
	}
	if err = s.EditDMState(c, 1, cid, mid, dm2Mdl.StateUserDelete, oplog.SourcePlayer, oplog.OperatorMember, id); err != nil {
		log.Error("s.EditDMStat(%d,%d) error(%v)", cid, id, err)
		err = ecode.DMRecallError
		return
	}
	if err = s.dao.UptRecallCnt(c, mid); err != nil {
		log.Error("s.dao.Item(%d,%d) error(%v)", cid, mid, err)
		err = nil
	}
	msg = fmt.Sprintf(_recallOK, ts)
	return
}

// EditDMState edit dm state used rpc method in dm2.
func (s *Service) EditDMState(c context.Context, tp int32, oid, mid int64, state int32, source oplog.Source, operatorType oplog.OperatorType, dmids ...int64) (err error) {
	arg := &dm2Mdl.ArgEditDMState{
		Type:         tp,
		Oid:          oid,
		Mid:          mid,
		State:        state,
		Dmids:        dmids,
		Source:       source,
		OperatorType: operatorType,
	}
	if err = s.dmRPC.EditDMState(c, arg); err != nil {
		log.Error("dmRPC.EditDMState(%v) error(%v)", arg, err)
	}
	return
}

// MidHash 弹幕发送者mid hash.
func (s *Service) MidHash(c context.Context, mid int64) (hash string, err error) {
	hash = model.Hash(mid, 0)
	return
}

// TransferJob set task to db
func (s *Service) TransferJob(c context.Context, mid, fromCid, toCid int64, offset float64) (err error) {
	job, err := s.dao.CheckTransferJob(c, fromCid, toCid)
	if err != nil {
		log.Error("dao.CheckTransferJob(from:%d,to:%d) err(%v)", fromCid, toCid, err)
		return
	}
	if job != nil && job.FromCID == fromCid && job.ToCID == toCid && job.State != model.TransferJobStatFailed {
		err = ecode.DMTransferRepet
		return
	}
	_, err = s.dao.AddTransferJob(c, fromCid, toCid, mid, offset, model.TransferJobStatInit)
	if err != nil {
		log.Error("dao.AddTransferJob(from:%d,to:%d) err(%v)", fromCid, toCid, err)
	}
	return
}

// TransferList service
func (s *Service) TransferList(c context.Context, cid int64) (hiss []*model.TransferHistory, err error) {
	hiss, err = s.dao.TransferList(c, cid)
	if err != nil || len(hiss) == 0 {
		return
	}
	for _, his := range hiss {
		cidInfo, err := s.dao.CidInfo(c, his.CID)
		if err != nil {
			log.Error("dao.CidInfo(%d) err(%v)", cid, err)
			continue
		}
		his.Title = cidInfo.Title
		his.PartID = int32(cidInfo.Index)
	}
	return
}

// TransferRetry change transferjob state
func (s *Service) TransferRetry(c context.Context, id, mid int64) (err error) {
	job, err := s.dao.CheckTransferID(c, id)
	if err != nil {
		log.Error("dao.CheckTransferID(%d) err(%v)", id, err)
		return
	}
	if job.State != model.TransferJobStatFailed || job.MID != mid {
		err = ecode.RequestErr
		return
	}
	_, err = s.dao.SetTransferState(c, id, model.TransferJobStatInit)
	if err != nil {
		log.Error("dao.TransferList(%d %d %d) err(%v)", id, err)
	}
	return
}

// CheckExist check exit of up id.
func (s *Service) CheckExist(c context.Context, mid, cid int64) (err error) {
	sub, err := s.subject(c, 1, cid)
	if err != nil {
		if ecode.Cause(err) == ecode.NothingFound {
			err = ecode.DMTransferNotFound
		}
		return
	}
	if sub.Mid == 0 {
		err = ecode.DMTransferNotFound
		return
	}
	if sub.Mid != mid {
		err = ecode.DMTransferNotBelong
		return
	}
	return
}

// dms get dm list by dmid from database
func (s *Service) dms(c context.Context, tp int32, oid int64, ids []int64) (dms []*model.DM, err error) {
	var (
		idxMap     = make(map[int64]*model.DM)
		contentSpe = make(map[int64]*model.ContentSpecial)
		special    []int64
		contents   []*model.Content
	)
	if idxMap, special, err = s.dao.IndexsByID(c, tp, oid, ids); err != nil || len(idxMap) == 0 {
		return
	}
	if contents, err = s.dao.Contents(c, oid, ids); err != nil {
		return
	}
	if len(special) > 0 {
		if contentSpe, err = s.dao.ContentsSpecial(c, special); err != nil {
			return
		}
	}
	for _, content := range contents {
		if dm, ok := idxMap[content.ID]; ok {
			dm.Content = content
			if dm.Pool == model.PoolSpecial {
				if _, ok = contentSpe[dm.ID]; ok {
					dm.ContentSpe = contentSpe[dm.ID]
				}
			}
			dms = append(dms, dm)
		}
	}
	return
}
