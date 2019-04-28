package service

import (
	"go-common/app/service/main/upcredit/common/election"
	"go-common/app/service/main/upcredit/conf"
	"go-common/app/service/main/upcredit/dao/upcrmdao"
	"go-common/app/service/main/upcredit/mathutil"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"sync"
	"time"
)

// Service is service.
type Service struct {
	c                    *conf.Config
	httpClient           *bm.Client
	upcrmdb              *upcrmdao.Dao
	creditLogSub         *databus.Databus
	businessBinLogSub    *databus.Databus
	wg                   sync.WaitGroup
	CreditScoreInputChan chan *upcrmmodel.UpScoreHistory
	CalcSvc              *CalcService
	limit                *mathutil.Limiter
	running              bool
	closeChan            chan struct{}
	businessScoreChan    chan *upcrmdao.UpQualityInfo
	zkElection           *election.ZkElection
}

// New is go-common/app/service/videoup service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                    c,
		httpClient:           bm.NewClient(c.HTTPClient.Normal),
		upcrmdb:              upcrmdao.New(c),
		creditLogSub:         databus.New(c.CreditLogSub),
		businessBinLogSub:    databus.New(c.BusinessBinLogSub),
		CreditScoreInputChan: make(chan *upcrmmodel.UpScoreHistory, 10240),
		running:              true,
		closeChan:            make(chan struct{}),
		businessScoreChan:    make(chan *upcrmdao.UpQualityInfo, 200),
	}

	if conf.Conf.ElectionZooKeeper != nil {
		var zc = conf.Conf.ElectionZooKeeper
		s.zkElection = election.New(zc.Addrs, zc.Root, time.Duration(zc.Timeout))
		var err = s.zkElection.Init()
		if err != nil {
			log.Error("zk elect init fail, err=%s", err)
		} else {
			err = s.zkElection.Elect()
			if err != nil {
				log.Error("zk elect fail, err=%s", err)
			} else {
				go func() {
					for {
						conf.IsMaster = <-s.zkElection.C
						if conf.IsMaster {
							log.Info("this is master, node=%s", s.zkElection.NodePath)
						} else {
							log.Info("this is follower, node=%s, master=%s", s.zkElection.NodePath, s.zkElection.MasterPath)
						}
					}
				}()
			}
		}
	}

	s.CalcSvc = NewCalc(c, s.CreditScoreInputChan, s.upcrmdb)
	if c.MiscConf.BusinessBinLogLimitRate <= 0 {
		c.MiscConf.BusinessBinLogLimitRate = 300
	}
	s.limit = mathutil.NewLimiter(c.MiscConf.BusinessBinLogLimitRate)
	s.CalcSvc.Run()
	// credit log databus
	{
		s.wg.Add(2)
		go s.arcCreditLogConsume()
		go s.arcBusinessBinLogCanalConsume()
	}

	for i := 0; i < c.MiscConf.CreditLogWriteRoutineNum; i++ {
		s.wg.Add(1)
		go s.WriteStatData()
	}
	s.wg.Add(1)
	go s.updateScoreProc()
	s.wg.Add(1)
	go s.batchWriteProc()
	return s
}

//Close service close
func (s *Service) Close() {
	s.creditLogSub.Close()
	s.businessBinLogSub.Close()
	s.CalcSvc.Close()
	s.running = false
	close(s.closeChan)
	s.wg.Wait()
	s.upcrmdb.Close()
}

//
//func (s *Service) ArchiveAuditResult(c *context.Context, ap *archive.ArcParam) {
//	// bussiness type = 1 是稿件
//	// 检查ap的state是不是需要记录的state 对应到 type 中
//	// 检查对应的reason id， 对应到具体的optype 分类中
//	// 检查对应的
//}
