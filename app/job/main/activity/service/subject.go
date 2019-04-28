package service

import (
	"context"
	"encoding/json"

	actmdl "go-common/app/interface/main/activity/model/like"
	l "go-common/app/job/main/activity/model/like"
	"go-common/library/log"
)

// upSubject update act_subject cache .
func (s *Service) upSubject(c context.Context, upMsg json.RawMessage) (err error) {
	var (
		subObj = new(l.ActSubject)
	)
	if err = json.Unmarshal(upMsg, subObj); err != nil {
		log.Error("upSubject json.Unmarshal(%s) error(%+v)", upMsg, err)
		return
	}
	if err = s.actRPC.SubjectUp(c, &actmdl.ArgSubjectUp{Sid: subObj.ID}); err != nil {
		log.Error("s.actRPC.SubjectUp(%d) error(%+v)", subObj.ID, err)
		return
	}
	log.Info("upSubject success s.actRPC.SubjectUp(%d)", subObj.ID)
	return
}

//func (s *Service) addElemeLottery(c context.Context, msg json.RawMessage) {
//	var (
//		vipContent = &l.VipActOrder{}
//	)
//	if err := json.Unmarshal(msg, vipContent); err != nil {
//		log.Error("addElemeLottery json.Unmarshal(%s) error(%v)", msg, err)
//		return
//	}
//	if vipContent.PanelType == "ele" {
//		if err := s.dao.AddLotteryTimes(c, s.c.Rule.EleLotteryID, vipContent.Mid); err != nil {
//			log.Error("s.dao.AddLotteryTimes(%d,%d) error(%v)", s.c.Rule.EleLotteryID, vipContent.Mid, err)
//			return
//		}
//		log.Info("addElemeLottery has AddLotteryTimes %d", vipContent.Mid)
//	}
//	log.Info("addElemeLottery success id:%d,panee:%s", vipContent.ID, vipContent.PanelType)
//}
