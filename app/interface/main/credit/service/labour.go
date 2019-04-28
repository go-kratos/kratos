package service

import (
	"context"
	"encoding/json"
	"time"

	model "go-common/app/interface/main/credit/model"
	"go-common/app/service/main/archive/api"
	arcMDL "go-common/app/service/main/archive/model/archive"
	blkmdl "go-common/app/service/main/member/model/block"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

// AddQs add labour question.
func (s *Service) AddQs(c context.Context, qs *model.LabourQs) (err error) {
	return s.dao.AddQs(c, qs)
}

// SetQs set labour question field.
func (s *Service) SetQs(c context.Context, id int64, ans int64, status int64) (err error) {
	return s.dao.SetQs(c, id, ans, status)
}

// DelQs del labour question.
func (s *Service) DelQs(c context.Context, id int64, isDel int64) (err error) {
	return s.dao.DelQs(c, id, isDel)
}

// GetQs get question.
func (s *Service) GetQs(c context.Context, mid int64) (qs []*model.LabourQs, err error) {
	var (
		ok                   bool
		arc                  *api.Arc
		getids, qsids, avids []int64
		marc                 map[int64]*api.Arc
		block                *blkmdl.RPCResInfo
	)
	defer func() {
		if err == nil {
			if len(avids) != 0 {
				marc, _ = s.arcRPC.Archives3(c, &arcMDL.ArgAids2{Aids: avids, RealIP: metadata.String(c, metadata.RemoteIP)})
			}
			for _, q := range qs {
				qsids = append(qsids, q.ID)
				if arc, ok = marc[q.AvID]; !ok {
					log.Warn("aid(%d) is not extists, [mid(%d)-qid(%d)]", q.AvID, mid, q.ID)
					q.AvID = 0
					continue
				}
				q.AvTitle = arc.Title
				if !model.ArcVisible(arc.State) {
					log.Warn("aid(%d) atitle(%s) state(%d) is not visible, [mid(%d)-qid(%d)]", q.AvID, q.AvTitle, arc.State, mid, q.ID)
					q.AvID = 0
				}
			}
			qsCache := &model.QsCache{
				Stime: xtime.Time(time.Now().Unix()),
				QsStr: xstr.JoinInts(qsids),
			}
			s.dao.SetQsCache(c, mid, qsCache)
		}
	}()
	if block, err = s.memRPC.BlockInfo(c, &blkmdl.RPCArgInfo{MID: mid}); err != nil {
		return
	}
	status := int8(block.BlockStatus)
	if status == model.BlockStatusNone {
		err = ecode.CreditNoblock
		return
	}
	if status == model.BlockStatusForever {
		err = ecode.CreditForeverBlock
		return
	}
	var qsIDs *model.AIQsID
	qsIDs, err = s.dao.GetQS(c, mid)
	if err != nil {
		log.Error("s.dao.GetQS(%d,%s) error(%+v)", mid, metadata.String(c, metadata.RemoteIP), err)
		err = nil
		qs = s.question
		avids = s.avIDs
		return
	}
	getids = append(getids, qsIDs.Pend...)
	getids = append(getids, qsIDs.Done...)
	idStr := xstr.JoinInts(getids)
	if _, qs, avids, err = s.dao.QsAllList(c, idStr); err != nil {
		return
	}
	if len(qs) != s.c.Property.QsNum {
		log.Warn("creditQsNumError(mid:%d,idstr:%s,qs:%+v),len:%d", mid, idStr, qs, len(qs))
		qs = s.question
		avids = s.avIDs
	}
	return
}

// CommitQs commit questions.
func (s *Service) CommitQs(c context.Context, mid int64, refer string, ua string, buvid string, ans *model.LabourAns) (commitRs *model.CommitRs, err error) {
	var (
		num   int64
		qs    map[int64]*model.LabourQs
		block *blkmdl.RPCResInfo
	)
	if block, err = s.memRPC.BlockInfo(c, &blkmdl.RPCArgInfo{MID: mid}); err != nil {
		return
	}
	status := int8(block.BlockStatus)
	if status == model.BlockStatusNone {
		err = ecode.CreditNoblock
		return
	}
	if status == model.BlockStatusForever {
		err = ecode.CreditForeverBlock
		return
	}
	if len(ans.ID) != s.c.Property.QsNum || len(ans.Answer) != s.c.Property.QsNum {
		err = ecode.CreditAnsNumError
		log.Error("CreditAnsNumError(mid:%d,id:%+v，ans:%+v)", mid, ans.ID, ans.Answer)
		return
	}
	idStr := xstr.JoinInts(ans.ID)
	qsCache, _ := s.dao.GetQsCache(c, mid)
	if qsCache == nil || qsCache.QsStr != idStr {
		err = ecode.RequestErr
		return
	}
	if qs, _, _, err = s.dao.QsAllList(c, idStr); err != nil {
		return
	}
	if len(qs) != s.c.Property.QsNum {
		log.Error("CreditRightAnsNumError(mid:%d,qs:%+v)", mid, qs)
	}
	for id, qsid := range ans.ID {
		if ans.Answer[id] != 1 && ans.Answer[id] != 2 {
			err = ecode.RequestErr
			log.Error("CreditAnsError(mid:%d,id:%+v，ans:%+v)", mid, ans.ID, ans.Answer)
			return
		}
		if v, ok := qs[qsid]; ok {
			if v.Ans == ans.Answer[id] {
				num++
			}
		}
	}
	commitRs = &model.CommitRs{}
	commitRs.Score = num * s.c.Property.PerScore
	if commitRs.Score >= 100 {
		commitRs.Score = 100
		if status == model.BlockStatusNone {
			commitRs.Day = 0
		} else {
			rts := time.Until(time.Unix(block.EndTime, 0))
			commitRs.Day = int64(rts / (time.Hour * 24))
			if int64(rts%(time.Hour*24)) > 0 {
				commitRs.Day++
			}
		}
	}
	anstr, err := json.Marshal(ans)
	if err != nil {
		log.Error("json.Marshal() error(%v)", err)
		return
	}
	id, err := s.dao.AddAnsLog(c, mid, commitRs.Score, string(anstr), qsCache.Stime)
	if err != nil {
		return
	}
	var msg model.DataBusResult
	msg.Mid = mid
	msg.Buvid = buvid
	msg.IP = metadata.String(c, metadata.RemoteIP)
	msg.Ua = ua
	msg.Refer = refer
	msg.Score = commitRs.Score
	for idx, qsid := range ans.ID {
		var rs = model.Rs{}
		rs.ID = qsid
		rs.Ans = ans.Answer[idx]
		if v, ok := qs[qsid]; ok {
			rs.Question = v.Question
			rs.TrueAns = v.Ans
			rs.AvID = v.AvID
			rs.Status = v.Status
			rs.Source = v.Source
			rs.Ctime = v.Ctime
			rs.Mtime = v.Mtime
		}
		msg.Rs = append(msg.Rs, rs)
	}
	if err = s.dao.PubLabour(c, id, msg); err != nil {
		log.Error("s.dao.PubLabour(%d,%+v) error(%+v)", id, msg, err)
		return
	}
	log.Info("PubLabour id(%d) msg(%+v)", id, msg)
	s.dao.DelQsCache(c, mid)
	return
}

// IsAnswered labour check user is answwered question between the time.
func (s *Service) IsAnswered(c context.Context, mid int64, mtime int64) (state int8, err error) {
	var (
		mc    = true
		found bool
	)
	if state, found, err = s.dao.GetAnswerStateCache(c, mid); err != nil {
		err = nil
		mc = false
	}
	if found {
		return
	}
	var status bool
	if status, err = s.dao.AnswerStatus(c, mid, time.Unix(mtime, 0)); err != nil {
		return
	}
	if status {
		state = model.LabourOkAnswer
	}
	if mc {
		s.addCache(func() {
			s.dao.SetAnswerStateCache(context.TODO(), mid, state)
		})
	}
	return
}
