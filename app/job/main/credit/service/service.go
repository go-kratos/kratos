package service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"go-common/app/job/main/credit/conf"
	"go-common/app/job/main/credit/dao"
	archive "go-common/app/service/main/archive/api/gorpc"
	memrpc "go-common/app/service/main/member/api/gorpc"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service struct of service.
type Service struct {
	c           *conf.Config
	dao         *dao.Dao
	credbSub    *databus.Databus
	replyAllSub *databus.Databus
	labourSub   *databus.Databus
	arcRPC      *archive.Service2
	memRPC      *memrpc.Service
	// wait group
	wg sync.WaitGroup
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		dao:         dao.New(c),
		credbSub:    databus.New(c.DataBus.CreditDBSub),
		replyAllSub: databus.New(c.DataBus.ReplyAllSub),
		labourSub:   databus.New(c.DataBus.LabourSub),
		arcRPC:      archive.New2(c.RPCClient.Archive),
		memRPC:      memrpc.New(c.RPCClient.Member),
	}
	s.loadConf()
	s.loadCase()
	s.loadDealWrongCase()
	s.wg.Add(1)
	go s.replyAllConsumer()
	go s.creditConsumer()
	go s.loadConfproc()
	return
}

// Ping check service health.
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}

func (s *Service) loadConfproc() {
	for {
		time.Sleep(time.Duration(s.c.Judge.ConfTimer))
		s.loadConf()
		s.loadCase()
		s.loadDealWrongCase()
	}
}

func (s *Service) loadConf() {
	m, err := s.dao.LoadConf(context.TODO())
	if err != nil {
		log.Error("loadConf error(%v)", err)
		return
	}
	if s.c.Judge.CaseGiveHours, err = strconv.ParseInt(m["case_give_hours"], 10, 64); err != nil {
		log.Error("loadConf CaseGiveHours error(%v)", err)
	}
	if s.c.Judge.CaseCheckTime, err = strconv.ParseInt(m["case_check_hours"], 10, 64); err != nil {
		log.Error("loadConf CaseCheckTime error(%v)", err)
	}
	if s.c.Judge.JuryRatio, err = strconv.ParseInt(m["jury_vote_radio"], 10, 64); err != nil {
		log.Error("loadConf JuryRatio error(%v)", err)
	}
	if s.c.Judge.JudgeRadio, err = strconv.ParseInt(m["case_judge_radio"], 10, 64); err != nil {
		log.Error("loadConf JudgeRadio error(%v)", err)
	}
	if s.c.Judge.CaseVoteMin, err = strconv.ParseInt(m["case_vote_min"], 10, 64); err != nil {
		log.Error("loadConf CaseVoteMin error(%v)", err)
	}
	if s.c.Judge.CaseObtainMax, err = strconv.ParseInt(m["case_obtain_max"], 10, 64); err != nil {
		log.Error("loadConf CaseObtainMax error(%v)", err)
	}
	if s.c.Judge.CaseVoteMax, err = strconv.ParseInt(m["case_vote_max"], 10, 64); err != nil {
		log.Error("loadConf CaseVoteMax error(%v)", err)
	}
	if s.c.Judge.JuryApplyMax, err = strconv.ParseInt(m["jury_apply_max"], 10, 64); err != nil {
		log.Error("loadConf JuryApplyMax error(%v)", err)
	}
	if s.c.Judge.CaseLoadMax, err = strconv.Atoi(m["case_load_max"]); err != nil {
		log.Error("loadConf CaseLoadMax error(%v)", err)
	}
	var caseLoadSwitch int64
	if caseLoadSwitch, err = strconv.ParseInt(m["case_load_switch"], 10, 64); err != nil {
		log.Error("loadConf CaseLoadSwitch error(%v)", err)
	}
	s.c.Judge.CaseLoadSwitch = int8(caseLoadSwitch)
	if s.c.Judge.CaseVoteMaxPercent, err = strconv.Atoi(m["case_vote_max_percent"]); err != nil {
		log.Error("loadConf CaseVoteMaxPercent error(%v)", err)
	}
	if _, ok := m["vote_num"]; !ok {
		s.c.Judge.VoteNum.RateS = 1
		s.c.Judge.VoteNum.RateA = 1
		s.c.Judge.VoteNum.RateB = 1
		s.c.Judge.VoteNum.RateC = 1
		s.c.Judge.VoteNum.RateD = 1
		return
	}
	if err = json.Unmarshal([]byte(m["vote_num"]), &s.c.Judge.VoteNum); err != nil {
		log.Error("loadConf vote_num error(%v)", err)
	}
}

// Close kafka consumer close.
func (s *Service) Close() {
	if s == nil {
		return
	}
	if s.dao != nil {
		s.dao.Close()
	}
	if s.credbSub != nil {
		s.credbSub.Close()
	}
	if s.replyAllSub != nil {
		s.replyAllSub.Close()
	}
	s.wg.Wait()
}
