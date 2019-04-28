package service

import (
	"context"
	"time"

	"go-common/app/interface/main/answer/conf"
	"go-common/app/interface/main/answer/dao"
	accountDao "go-common/app/interface/main/answer/dao/account"
	geetestDao "go-common/app/interface/main/answer/dao/geetest"
	"go-common/app/interface/main/answer/model"
	accoutCli "go-common/app/service/main/account/api"
	memrpc "go-common/app/service/main/member/api/gorpc"
	"go-common/library/log"
	"go-common/library/log/anticheat"
	"go-common/library/queue/databus/report"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

// Service struct of service.
type Service struct {
	c                 *conf.Config
	answerDao         *dao.Dao
	geetestDao        *geetestDao.Dao
	accountDao        *accountDao.Dao
	accountSvc        accoutCli.AccountClient
	memRPC            *memrpc.Service
	missch            *fanout.Fanout
	beformalch        chan *model.Formal
	questionTypeCache map[int64]*model.TypeInfo
	rankCache         []*model.RankInfo
	mRankCache        []*model.RankInfo
	tcQestTick        time.Duration
	rankQuestTick     time.Duration
	promBeFormal      *prom.Prom

	infoc2 *anticheat.AntiCheat
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                 c,
		answerDao:         dao.New(c),
		geetestDao:        geetestDao.New(c),
		accountDao:        accountDao.New(c),
		memRPC:            memrpc.New(c.RPCClient.Member),
		missch:            fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		beformalch:        make(chan *model.Formal, 1024),
		questionTypeCache: map[int64]*model.TypeInfo{},
		rankCache:         []*model.RankInfo{},
		mRankCache:        []*model.RankInfo{},
		tcQestTick:        time.Duration(c.Question.TcQestTick),
		rankQuestTick:     time.Duration(c.Question.RankQestTick),
		promBeFormal:      prom.New().WithCounter("answer_beformal_count", []string{"name"}),
	}
	accountSvc, err := accoutCli.NewClient(c.AccountRPC)
	if err != nil {
		panic(err)
	}
	s.accountSvc = accountSvc
	s.loadQidsCache()
	s.loadExtraQidsCache()
	s.loadtypes()
	go s.rankcacheproc()
	go s.beformalproc()
	if c.Infoc2 != nil {
		s.infoc2 = anticheat.New(c.Infoc2)
	}
	return
}

func (s *Service) addRetryBeFormal(msg *model.Formal) {
	select {
	case s.beformalch <- msg:
	default:
		log.Warn("beformalch chan full")
	}
}

func (s *Service) beformalproc() {
	var (
		err error
		c   = context.Background()
		msg *model.Formal
	)
	for {
		msg = <-s.beformalch
		for retries := 0; retries < s.c.Answer.MaxRetries; retries++ {
			if err = s.accountDao.BeFormal(c, msg.Mid, msg.IP); err != nil {
				sleep := s.c.Backoff.Backoff(retries)
				log.Error("beFormal fail(%d) sleep(%d) err(%+v)", msg.Mid, sleep, err)
				time.Sleep(sleep * time.Second)
				continue
			}
			break
		}
	}
}

// Close dao.
func (s *Service) Close() {
	s.answerDao.Close()
}

func (s *Service) rankcacheproc() {
	for {
		time.Sleep(s.tcQestTick)
		s.loadtypes()
		s.loadQidsCache()
		s.loadExtraQidsCache()
	}
}

func (s *Service) userActionLog(mid int64, typ string, ah *model.AnswerHistory) {
	report.User(&report.UserInfo{
		Mid:      mid,
		Business: model.AnswerLogID,
		Action:   model.AnswerUpdate,
		Ctime:    time.Now(),
		Content: map[string]interface{}{
			typ: ah,
		},
	})
}
