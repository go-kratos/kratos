package service

import (
	"context"
	"strings"
	"sync"

	"go-common/app/job/main/aegis/conf"
	"go-common/app/job/main/aegis/dao"
	"go-common/app/job/main/aegis/dao/email"
	"go-common/app/job/main/aegis/dao/monitor"
	"go-common/app/job/main/aegis/model"
	accApi "go-common/app/service/main/account/api"
	upApi "go-common/app/service/main/up/api/v1"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
)

// Service struct
type Service struct {
	c       *conf.Config
	acc     accApi.AccountClient
	up      upApi.UpClient
	dao     *dao.Dao
	moniDao *monitor.Dao
	email   *email.Dao
	// databus
	binLogDataBus    *databus.Databus
	archiveDataBus   *databus.Databus
	aegisRscDataBus  *databus.Databus
	aegisTaskDataBus *databus.Databus
	//channel
	chanReport chan *model.RIR
	// cache
	Cache
	//权重计算器
	wmHash     map[string]*WeightManager
	rschandle  map[string]RscHandler
	taskhandle map[string]TaskHandler

	wg sync.WaitGroup

	//databus group
	resourceGroup *databusutil.Group
	taskGroup     *databusutil.Group
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		dao:     dao.New(c),
		moniDao: monitor.New(c),
		email:   email.New(c),

		binLogDataBus:    databus.New(c.DataBus.BinLogSub),
		chanReport:       make(chan *model.RIR, 1024),
		archiveDataBus:   databus.New(c.DataBus.ArchiveSub),
		aegisRscDataBus:  databus.New(c.DataBus.ResourceSub),
		aegisTaskDataBus: databus.New(c.DataBus.TaskSub),
	}
	if !s.c.Debug {
		var err error
		if s.acc, err = accApi.NewClient(c.GRPC.Acc); err != nil {
			panic(err)
		}
		if s.up, err = upApi.NewClient(c.GRPC.Up); err != nil {
			panic(err)
		}
	}
	initHandler(s)
	s.initCache()
	s.startWeightManager()

	s.resourceGroup = databusutil.NewGroup(c.Databusutil.Resource, s.aegisRscDataBus.Messages())
	s.resourceGroup.New = s.newrsc
	s.resourceGroup.Split = s.splitrsc
	s.resourceGroup.Do = s.dorsc
	s.resourceGroup.Start()

	s.taskGroup = databusutil.NewGroup(c.Databusutil.Task, s.aegisTaskDataBus.Messages())
	s.taskGroup.New = s.newtask
	s.taskGroup.Split = s.splittask
	s.taskGroup.Do = s.dotask
	s.taskGroup.Start()

	go s.cacheProc()
	go s.taskProc()
	go s.monitorNotify()
	s.wg.Add(1)
	go s.taskconsumeproc()
	s.wg.Add(1)
	go s.archiveConsumeProc()
	s.wg.Add(1)
	go s.monitorEmailProc()
	return s
}

// Cache .
type Cache struct {
	upCache          map[int64]map[int64]struct{}
	rangeWeightCfg   map[int64]map[string]*model.RangeWeightConfig
	equalWeightCfg   map[string][]*model.EqualWeightConfig
	assignConfig     map[string][]*model.AssignConfig
	consumerCache    map[string]map[int64]struct{}
	ccMux            sync.RWMutex
	oldactiveBizFlow map[string]struct{}
	newactiveBizFlow map[string]struct{}
}

// DebugCache .
func (s *Service) DebugCache(keys string) map[string]interface{} {
	dc := map[string]interface{}{
		"upCache":          s.upCache,
		"rangeWeightCfg":   s.rangeWeightCfg,
		"equalWeightCfg":   s.equalWeightCfg,
		"assignConfig":     s.assignConfig,
		"consumerCache":    s.consumerCache,
		"oldactiveBizFlow": s.oldactiveBizFlow,
		"newactiveBizFlow": s.newactiveBizFlow,
	}
	res := make(map[string]interface{})
	if len(keys) > 0 {
		for _, key := range strings.Split(keys, ",") {
			res[key] = dc[key]
		}
	}
	return res
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.binLogDataBus.Close()
	s.archiveDataBus.Close()
	s.aegisRscDataBus.Close()
	s.aegisTaskDataBus.Close()
	s.resourceGroup.Close()
	s.taskGroup.Close()
	s.wg.Wait()
	s.dao.Close()
	s.moniDao.Close()
}
