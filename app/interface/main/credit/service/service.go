package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/interface/main/credit/conf"
	dao "go-common/app/interface/main/credit/dao"
	model "go-common/app/interface/main/credit/model"
	accgrpc "go-common/app/service/main/account/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	fligrpc "go-common/app/service/main/filter/api/grpc/v1"
	memrpc "go-common/app/service/main/member/api/gorpc"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// Service struct of service.
type Service struct {
	dao *dao.Dao
	// rpc
	arcRPC *arcrpc.Service2
	memRPC *memrpc.Service
	// grpc
	accountClient accgrpc.AccountClient
	fliClient     fligrpc.FilterClient
	// conf
	c        *conf.Config
	question []*model.LabourQs
	avIDs    []int64
	missch   chan func()
	// announcement
	announcement *announcement
	managers     map[string]int64
	tagMap       map[int8]int64
}

type announcement struct {
	def   []*model.BlockedAnnouncement
	alist map[int8][]*model.BlockedAnnouncement
	amap  map[int64]*model.BlockedAnnouncement
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		missch: make(chan func(), 1024000),
		arcRPC: arcrpc.New2(c.RPCClient2.Archive),
		memRPC: memrpc.New(c.RPCClient2.Member),
		tagMap: make(map[int8]int64),
		announcement: &announcement{
			def:   make([]*model.BlockedAnnouncement, 0, 4),
			alist: make(map[int8][]*model.BlockedAnnouncement),
			amap:  make(map[int64]*model.BlockedAnnouncement),
		},
	}
	var err error
	if s.fliClient, err = fligrpc.NewClient(c.GRPCClient.Filter); err != nil {
		panic(errors.WithMessage(err, "Failed to dial filter service"))
	}
	if s.accountClient, err = accgrpc.NewClient(c.GRPCClient.Account); err != nil {
		panic(errors.WithMessage(err, "Failed to dial account service"))
	}
	s.initTag()
	s.loadConf()
	s.loadQuestion()
	s.loadManager()
	s.LoadAnnouncement(context.TODO())
	go s.loadConfproc()
	go s.loadQuestionproc()
	go s.loadManagerproc()
	go s.loadAnnouncementproc()
	go s.cacheproc()
	return
}

func (s *Service) loadConfproc() {
	for {
		time.Sleep(time.Duration(s.c.Judge.ConfTimer))
		s.loadConf()
	}
}

func (s *Service) loadQuestionproc() {
	for {
		time.Sleep(time.Duration(s.c.Judge.ConfTimer))
		s.loadQuestion()
	}
}

func (s *Service) loadManagerproc() {
	for {
		time.Sleep(time.Duration(s.c.Judge.LoadManagerTime))
		s.loadManager()
	}
}

func (s *Service) loadAnnouncementproc() {
	for {
		time.Sleep(time.Duration(s.c.Judge.ConfTimer))
		s.LoadAnnouncement(context.TODO())
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

func (s *Service) initTag() {
	s.tagMap[model.OriginReply] = s.c.TagID.Reply
	s.tagMap[model.OriginDM] = s.c.TagID.DM
	s.tagMap[model.OriginMsg] = s.c.TagID.Msg
	s.tagMap[model.OriginTag] = s.c.TagID.Tag
	s.tagMap[model.OriginMember] = s.c.TagID.Member
	s.tagMap[model.OriginArchive] = s.c.TagID.Archive
	s.tagMap[model.OriginMusic] = s.c.TagID.Music
	s.tagMap[model.OriginArticle] = s.c.TagID.Article
	s.tagMap[model.OriginSpaceTop] = s.c.TagID.SpaceTop
}

func (s *Service) loadManager() {
	managers, err := s.dao.Managers(context.TODO())
	if err != nil {
		log.Error("s.dao.Managers error(%v)", err)
		return
	}
	s.managers = managers
}

func (s *Service) loadQuestion() {
	audit, avIDs, err := s.dao.LastAuditQuestion(context.TODO())
	if err != nil {
		log.Error("s.dao.LastAuditQuestion error(%v)", err)
		return
	}
	noAudit, noAvIDs, err := s.dao.LastNoAuditQuestion(context.TODO())
	if err != nil {
		log.Error("s.dao.LastNoAuditQuestion error(%v)", err)
		return
	}
	audit = append(audit, noAudit...)
	avIDs = append(avIDs, noAvIDs...)
	s.question = audit
	s.avIDs = avIDs
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.Ping(c); err != nil {
		return
	}
	return s.dao.Ping(c)
}

// Close dao.
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for executing closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}
