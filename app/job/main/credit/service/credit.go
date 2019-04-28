package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/credit/model"
	"go-common/library/log"
)

const (
	_blockedJuryTable         = "blocked_jury"
	_blockedInfoTable         = "blocked_info"
	_blockedCaseTable         = "blocked_case"
	_voteOpinion              = "blocked_opinion"
	_blockedKpiTable          = "blocked_kpi"
	_blockedPublishTable      = "blocked_publish"
	_blockedVoteCaseTable     = "blocked_case_vote"
	_blockedCaseApplyLogTable = "blocked_case_apply_log"
	_blockedLabourAnswerLog   = "blocked_labour_answer_log"
	_retry                    = 3
	_retrySleep               = time.Second * 1
)

func (s *Service) creditConsumer() {
	var err error
	for res := range s.credbSub.Messages() {
		mu := &model.Message{}
		if err = json.Unmarshal(res.Value, mu); err != nil {
			log.Error("credit-job,json.Unmarshal (%v) error(%v)", string(res.Value), err)
			continue
		}
		for i := 0; ; i++ {
			err = s.dealCredit(mu)
			if err != nil {
				log.Error("s.flush error(%v)", err)
				time.Sleep(_retrySleep)
				if i > _retry && s.c.Env == "prod" {
					s.dao.Sms(context.TODO(), s.c.Sms.Phone, s.c.Sms.Token, "credit-job update cache fail for 5 times")
					i = 0
				}
				continue
			}
			break
		}
		if err = res.Commit(); err != nil {
			log.Error("databus.Commit err(%v)", err)
		}
		log.Info("subproc key:%v,topic: %v, part:%v offset:%v,message %s,", res.Key, res.Topic, res.Partition, res.Offset, res.Value)
	}
}

// dealAction deal databus action
func (s *Service) dealCredit(mu *model.Message) (err error) {
	switch mu.Table {
	case _blockedCaseTable:
		if mu.Action == "insert" {
			s.RegReply(context.TODO(), mu.Table, mu.New, mu.Old)
		}
		err = s.Judge(context.TODO(), mu.New, mu.Old)
		s.GrantCase(context.TODO(), mu.New, mu.Old)
		s.DelGrantCase(context.TODO(), mu.New, mu.Old)
		s.DelCaseInfoCache(context.TODO(), mu.New)
	case _blockedInfoTable:
		if mu.Action == "insert" {
			s.RegReply(context.TODO(), mu.Table, mu.New, mu.Old)
			s.InvalidJury(context.TODO(), mu.New, mu.Old)
		}
		if mu.Action == "update" {
			s.UnBlockAccount(context.TODO(), mu.New, mu.Old)
		}
		s.UpdateCache(context.TODO(), mu.New)
	case _blockedKpiTable:
		if mu.Action == "insert" {
			s.KPIReward(context.TODO(), mu.New, mu.Old)
		}
	case _voteOpinion:
		s.DeleteIdx(context.TODO(), mu.New)
	case _blockedPublishTable:
		s.RegReply(context.TODO(), mu.Table, mu.New, mu.Old)
	case _blockedVoteCaseTable:
		s.DelVoteCaseCache(context.TODO(), mu.New)
	case _blockedLabourAnswerLog:
		if mu.Action == "insert" {
			s.NotifyBlockAnswer(context.TODO(), mu.New)
		}
	case _blockedCaseApplyLogTable:
		if mu.Action == "insert" {
			s.DealCaseApplyReason(context.TODO(), mu.New)
		}
	case _blockedJuryTable:
		s.DelJuryInfoCache(context.TODO(), mu.New)
	}
	return
}
