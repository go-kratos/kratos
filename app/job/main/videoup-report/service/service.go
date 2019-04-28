package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/job/main/videoup-report/conf"
	arcdao "go-common/app/job/main/videoup-report/dao/archive"
	"go-common/app/job/main/videoup-report/dao/data"
	"go-common/app/job/main/videoup-report/dao/email"
	hbasedao "go-common/app/job/main/videoup-report/dao/hbase"
	"go-common/app/job/main/videoup-report/dao/manager"
	"go-common/app/job/main/videoup-report/dao/mission"
	redisdao "go-common/app/job/main/videoup-report/dao/redis"
	"go-common/app/job/main/videoup-report/dao/tag"
	arcmdl "go-common/app/job/main/videoup-report/model/archive"
	taskmdl "go-common/app/job/main/videoup-report/model/task"
	account "go-common/app/service/main/account/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	upsrpc "go-common/app/service/main/up/api/v1"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"strings"
)

const (
	_archiveTable = "archive"
	_videoTable   = "archive_video"
	_upsTable     = "ups"
	_jumpChanSize = int(4000)
	_logChanSize  = int(200000)
)

// Service is service.
type Service struct {
	c              *conf.Config
	arc            *arcdao.Dao
	redis          *redisdao.Dao
	hbase          *hbasedao.Dao
	dataDao        *data.Dao
	email          *email.Dao
	mng            *manager.Dao
	arcUpChs       []chan *arcmdl.UpInfo
	videoUpInfoChs []chan *arcmdl.VideoUpInfo
	// waiter
	waiter sync.WaitGroup
	// databus
	archiveSub   *databus.Databus
	arcResultSub *databus.Databus
	videoupSub   *databus.Databus
	ManagerDBSub *databus.Databus
	// cache
	sfTpsCache        map[int16]*arcmdl.Type
	adtTpsCache       map[int16]struct{}
	taskCache         *arcmdl.TaskCache
	videoAuditCache   *arcmdl.VideoAuditCache
	arcMoveTypeCache  *arcmdl.ArcMoveTypeCache
	arcRoundFlowCache *arcmdl.ArcRoundFlowCache
	xcodeTimeCache    *arcmdl.XcodeTimeCache
	assignCache       map[int64]*taskmdl.AssignConfig
	upperCache        map[int8]map[int64]struct{}
	weightCache       map[int8]map[int64]*taskmdl.ConfigItem
	missTagsCache     map[string]int

	lastjumpMap map[int64]struct{} //上轮插队的这一轮也更新，否则会出现权重只增不减
	jumplist    *taskmdl.JumpList  //插队序列
	jumpchan    chan *taskmdl.WeightLog
	tasklogchan chan *taskmdl.WeightLog
	//rpc
	arcRPCGroup2 *arcrpc.Service2

	//grpc
	accRPC account.AccountClient
	upsRPC upsrpc.UpClient

	// closed
	closed     bool
	tagDao     *tag.Dao
	missionDao *mission.Dao
}

// New is videoup-report-job service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		//dao
		arc:   arcdao.New(c),
		redis: redisdao.New(c),
		hbase: hbasedao.New(c),
		//databus
		archiveSub:   databus.New(c.ArchiveSub),
		arcResultSub: databus.New(c.ArchiveResultSub),
		videoupSub:   databus.New(c.VideoupSub),
		ManagerDBSub: databus.New(c.ManagerDBSub),
		dataDao:      data.New(c),
		email:        email.New(c),
		mng:          manager.New(c),
		// cache
		taskCache: &arcmdl.TaskCache{
			Task: make(map[int64]*arcmdl.Task),
		},
		videoAuditCache: &arcmdl.VideoAuditCache{
			Data: make(map[int16]map[string]int),
		},
		arcMoveTypeCache: &arcmdl.ArcMoveTypeCache{
			Data: make(map[int8]map[int16]map[string]int),
		},
		arcRoundFlowCache: &arcmdl.ArcRoundFlowCache{
			Data: make(map[int8]map[int64]map[string]int),
		},
		xcodeTimeCache: &arcmdl.XcodeTimeCache{
			Data: make(map[int8][]int),
		},
		arcRPCGroup2: arcrpc.New2(c.ArchiveRPCGroup2),
		lastjumpMap:  make(map[int64]struct{}),
		jumplist:     taskmdl.NewJumpList(),
		jumpchan:     make(chan *taskmdl.WeightLog, _jumpChanSize),
		tasklogchan:  make(chan *taskmdl.WeightLog, _logChanSize),
		tagDao:       tag.New(c),
		missionDao:   mission.New(c),
	}
	var err error
	if s.accRPC, err = account.NewClient(conf.Conf.GRPC.AccRPC); err != nil {
		panic(err)
	}

	if s.upsRPC, err = upsrpc.NewClient(conf.Conf.GRPC.UpsRPC); err != nil {
		panic(err)
	}

	for i := 0; i < s.c.ChanSize; i++ {
		log.Info("videoup-report-job chanSize starting(%d)", i)
		s.arcUpChs = append(s.arcUpChs, make(chan *arcmdl.UpInfo, 10240))
		s.waiter.Add(1)
		go s.arcUpdateproc(i)
		s.videoUpInfoChs = append(s.videoUpInfoChs, make(chan *arcmdl.VideoUpInfo, 10240))
		s.waiter.Add(1)
		go s.upVideoproc(i)
	}
	s.loadConf()
	s.loadType()
	// load cache.
	s.loadTask()
	s.loadTaskTookSort()
	s.hdlTraffic()
	s.loadMission()
	go s.cacheproc()
	go s.monitorNotifyProc()
	s.waiter.Add(1)
	go s.hotarchiveproc()
	s.waiter.Add(1)
	go s.arcCanalConsume()
	s.waiter.Add(1)
	go s.arcResultConsume()
	s.waiter.Add(1)
	go s.taskWeightConsumer()
	s.waiter.Add(1)
	go s.taskweightproc()
	s.waiter.Add(1)
	go s.movetaskproc()
	go s.deltaskproc()
	s.waiter.Add(1)
	go s.videoupConsumer()
	s.waiter.Add(1)
	go s.emailProc()
	s.waiter.Add(1)
	go s.emailFastProc()
	s.waiter.Add(1)
	go s.retryProc()
	s.waiter.Add(1)
	go s.managerDBConsume()

	return s
}

func (s *Service) loadType() {
	tpm, err := s.arc.TypeMapping(context.TODO())
	if err != nil {
		log.Error("s.dede.TypeMapping error(%v)", err)
		return
	}
	s.sfTpsCache = tpm
	// audit types
	adt, err := s.arc.AuditTypesConf(context.TODO())
	if err != nil {
		log.Error("s.dede.AuditTypesConf error(%v)", err)
		return
	}
	s.adtTpsCache = adt

	wvc, err := s.arc.WeightValueConf(context.TODO())
	if err != nil {
		log.Error("s.arc.WeightValueConf error(%v)", err)
		return
	}
	taskmdl.WLVConf = wvc
}

func (s *Service) isAuditType(tpID int16) bool {
	_, isAt := s.adtTpsCache[tpID]
	return isAt
}

func (s *Service) topType(tpID int16) (id int16) {
	if tp, ok := s.sfTpsCache[tpID]; ok && tp != nil {
		id = tp.PID
	}
	return
}

func (s *Service) typeName(tpID int16) (name string) {
	if tp, ok := s.sfTpsCache[tpID]; ok && tp != nil {
		name = tp.Name
	}
	return
}

func (s *Service) topTypeName(tpID int16) (name string) {
	pid := s.topType(tpID)
	name = s.typeName(pid)
	return
}

func (s *Service) cacheproc() {
	for {
		time.Sleep(1 * time.Minute)
		// config
		s.loadConf()
		// task
		s.loadTask()
		s.loadTaskTookSort()
		// handle task took
		s.hdlTaskTook()
		s.hdlTaskTookByHourHalf()
		// handle video audit
		s.hdlVideoAuditCount()
		s.hdlMoveTypeCount()
		s.hdlRoundFlowCount()
		//handle calculate video xcode time stats, and save to DB
		s.hdlXcodeStats()
		s.hdlTraffic()
		s.loadMission()
	}
}

func (s *Service) monitorNotifyProc() {
	for {
		s.monitorNotify()
		time.Sleep(30 * time.Minute)
	}
}

// s.missTagsCache: missionName or first tag
func (s *Service) loadMission() {
	mm, err := s.missionDao.Missions(context.TODO())
	if err != nil {
		log.Error("s.missionDao.Mission error(%v)", err)
		return
	}
	s.missTagsCache = make(map[string]int)
	for _, m := range mm {
		if len(m.Tags) > 0 {
			splitedTags := strings.Split(m.Tags, ",")
			s.missTagsCache[splitedTags[0]] = m.ID
		} else {
			s.missTagsCache[m.Name] = m.ID
		}
	}
}

// hotarchiveproc get hot archive which need to recheck
func (s *Service) hotarchiveproc() {
	defer s.waiter.Done()
	for {
		if s.closed {
			return
		}
		s.addHotRecheck()
		time.Sleep(10 * time.Minute)
	}
}

// Close  consumer close.
func (s *Service) Close() {
	s.closed = true
	s.archiveSub.Close()
	s.arcResultSub.Close()
	s.videoupSub.Close()
	time.Sleep(2 * time.Second)
	for i := 0; i < s.c.ChanSize; i++ {
		log.Info("videoup-report-job chanSize closing(%d)", i)
		close(s.arcUpChs[i])
		close(s.videoUpInfoChs[i])
	}
	s.arc.Close()
	s.mng.Close()
	s.redis.Close()
	s.email.Close()
	s.hbase.Close()
	s.waiter.Wait()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.arc.Ping(c); err != nil {
		return
	}
	if err = s.mng.Ping(c); err != nil {
		return
	}
	return
}

func (s *Service) loadConf() {
	var (
		err        error
		assignConf map[int64]*taskmdl.AssignConfig
		weightConf map[int8]map[int64]*taskmdl.ConfigItem
		upperCache map[int8]map[int64]struct{}
	)

	if assignConf, err = s.assignConf(context.TODO()); err != nil {
		log.Error("s.assignConf error(%v)", err)
		return
	}
	s.assignCache = assignConf

	upperCache, err = s.upSpecial(context.TODO())
	if err != nil {
		log.Error("s.upSpecial error(%v)", err)
	} else {
		s.upperCache = upperCache
	}

	if weightConf, err = s.weightConf(context.TODO()); err != nil {
		log.Error(" s.weightConf error(%v)", err)
		return
	}
	s.weightCache = weightConf
}
