package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/interface/main/push-archive/conf"
	"go-common/app/interface/main/push-archive/dao"
	"go-common/app/interface/main/push-archive/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	pb "go-common/app/service/main/push/api/grpc/v1"
	pushrpc "go-common/app/service/main/push/api/grpc/v1"
	"go-common/library/cache"
	"go-common/library/conf/env"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Service push service.
type Service struct {
	c            *conf.Config
	dao          *dao.Dao
	cache        *cache.Cache
	accRPC       *accrpc.Service3
	wg           sync.WaitGroup
	archiveSub   *databus.Databus
	relationSub  *databus.Databus
	userSettings map[int64]*model.Setting
	arcMo, relMo int64
	CloseCh      chan bool
	ForbidTimes  []conf.ForbidTime
	settingCh    chan *pb.SetSettingRequest
	pushRPC      pushrpc.PushClient
}

// New creates a push service instance.
func New(c *conf.Config) *Service {
	s := &Service{
		c:            c,
		dao:          dao.New(c),
		cache:        cache.New(1, 102400),
		accRPC:       accrpc.New3(c.AccountRPC),
		archiveSub:   databus.New(c.ArchiveSub),
		relationSub:  databus.New(c.RelationSub),
		userSettings: make(map[int64]*model.Setting),
		CloseCh:      make(chan bool),
		ForbidTimes:  c.ArcPush.ForbidTimes,
		settingCh:    make(chan *pb.SetSettingRequest, 3072),
	}
	var err error
	if s.pushRPC, err = pushrpc.NewClient(c.PushRPC); err != nil {
		panic(err)
	}
	s.mappingAbtest()
	go s.loadUserSettingsproc()
	time.Sleep(2 * time.Second) // consumeArchive will notice upper's fans, it depends on user's setting
	s.wg.Add(1)
	go s.loadUserSettingsproc()
	if c.Push.ProdSwitch {
		s.wg.Add(1)
		go s.consumeRelationproc()
		s.wg.Add(1)
		go s.consumeArchiveproc()
		go s.monitorConsume()
	}
	go s.clearStatisticsProc()
	s.wg.Add(1)
	go s.saveStatisticsProc()
	go s.setSettingProc()
	return s
}

// loadUserSettingsproc 若是用户没有设置推送开关，默认会被推送消息；因此，若是新设置load出错/无值，则不替换旧设置
func (s *Service) loadUserSettingsproc() {
	defer s.wg.Done()
	var (
		ps    int64 = 30000
		start int64
		end   int64
		mxID  int64
		err   error
		res   map[int64]*model.Setting
	)

	for {
		select {
		case _, ok := <-s.CloseCh:
			if !ok {
				log.Info("CloseCh is closed, close the loadUserSettingsproc")
				return
			}
		default:
		}

		start = 0
		end = 0
		err = nil
		res = make(map[int64]*model.Setting)
		mxID, err = s.dao.SettingsMaxID(context.TODO())
		if err != nil || mxID == 0 {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		for {
			start = end
			end += ps
			err = s.dao.SettingsAll(context.TODO(), start, end, &res)
			if err != nil {
				break
			}
			if end >= mxID {
				break
			}
		}

		if err != nil || len(res) == 0 {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		s.userSettings = res
		time.Sleep(time.Duration(s.c.Push.LoadSettingsInterval))
	}
}

func (s *Service) monitorConsume() {
	if env.DeployEnv != env.DeployEnvProd {
		return
	}
	var (
		arc int64 // archive result count
		rel int64 // relation count
	)
	for {
		time.Sleep(10 * time.Minute)
		if s.arcMo-arc == 0 {
			msg := "databus: push-archive archiveResult did not consume within ten minute"
			s.dao.WechatMessage(msg)
			log.Warn(msg)
		}
		arc = s.arcMo
		if s.relMo-rel == 0 {
			msg := "databus: push-archive relation did not consume within ten minute"
			s.dao.WechatMessage(msg)
			log.Warn(msg)
		}
		rel = s.relMo
	}
}

// Close closes service.
func (s *Service) Close() {
	s.archiveSub.Close()
	s.relationSub.Close()
	close(s.CloseCh)
	s.wg.Wait()
	s.dao.Close()
}

// Ping checks service.
func (s *Service) Ping(c *bm.Context) (err error) {
	err = s.dao.Ping(c)
	return
}

// getTodayTime 获取当日的某个时间点的时间
func (s *Service) getTodayTime(tm string) (todayTime time.Time, err error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	todayTime, err = time.ParseInLocation("2006-01-02 15:04:05", today+" "+tm, time.Local)
	if err != nil {
		log.Error("clearStatisticsProc time.ParseInLocation error(%v)", err)
	}
	return
}
