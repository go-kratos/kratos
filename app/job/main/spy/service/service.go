package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"go-common/app/job/main/spy/conf"
	"go-common/app/job/main/spy/dao"
	"go-common/app/job/main/spy/model"
	cmmdl "go-common/app/service/main/spy/model"
	spyrpc "go-common/app/service/main/spy/rpc/client"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
	xtime "go-common/library/time"

	"github.com/robfig/cron"
)

const (
	_bigdataEvent = "bt"
	_secretEvent  = "st"
)

// Service service def.
type Service struct {
	c               *conf.Config
	waiter          sync.WaitGroup
	dao             *dao.Dao
	eventDatabus    *databus.Databus
	spystatDatabus  *databus.Databus
	secLoginDatabus *databus.Databus
	spyRPC          *spyrpc.Service
	quit            chan struct{}
	cachech         chan func()
	retrych         chan func()
	spyConfig       map[string]interface{}
	configLoadTick  time.Duration
	promBlockInfo   *prom.Prom
	blockTick       time.Duration
	blockWaitTick   time.Duration
	allEventName    map[string]int64
	loadEventTick   time.Duration
	// activity events
	activityEvents map[string]struct{}
}

// New create a instance of Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:               c,
		dao:             dao.New(c),
		eventDatabus:    databus.New(c.Databus.EventData),
		spystatDatabus:  databus.New(c.Databus.SpyStatData),
		secLoginDatabus: databus.New(c.Databus.SecLogin),
		spyRPC:          spyrpc.New(c.SpyRPC),
		quit:            make(chan struct{}),
		cachech:         make(chan func(), 1024),
		retrych:         make(chan func(), 128),
		spyConfig:       make(map[string]interface{}),
		configLoadTick:  time.Duration(c.Property.ConfigLoadTick),
		promBlockInfo:   prom.New().WithCounter("spy_block_info", []string{"name"}),
		blockTick:       time.Duration(c.Property.BlockTick),
		blockWaitTick:   time.Duration(c.Property.BlockWaitTick),
		allEventName:    make(map[string]int64),
		loadEventTick:   time.Duration(c.Property.LoadEventTick),
	}

	if err := s.loadSystemConfig(); err != nil {
		panic(err)
	}
	if err := s.loadeventname(); err != nil {
		panic(err)
	}
	s.initActivityEvents()

	s.waiter.Add(1)
	go s.consumeproc()

	s.waiter.Add(1)
	go s.secloginproc()

	s.waiter.Add(1)
	go s.spystatproc()

	go s.retryproc()
	go s.cacheproc()

	go s.loadconfig()
	go s.blockcacheuser()

	go s.loadeventproc()

	t := cron.New()
	if err := t.AddFunc(s.c.Property.Block.CycleCron, s.cycleblock); err != nil {
		panic(err)
	}
	if err := t.AddFunc(s.c.Property.ReportCron, s.reportjob); err != nil {
		panic(err)
	}
	t.Start()

	return s
}

func (s *Service) consumeproc() {
	defer s.waiter.Done()
	defer func() {
		if x := recover(); x != nil {
			log.Error("eventproc unknown panic(%v)", x)
		}
	}()
	var (
		msg       *databus.Message
		eventMsg  *model.EventMessage
		ok        bool
		err       error
		msgChan   = s.eventDatabus.Messages()
		c         = context.TODO()
		preOffset int64
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("eventConsumeProc msgChan closed")
				return
			}
		case <-s.quit:
			log.Info("quit eventConsumeProc")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit error(%v)", err)
		}
		if preOffset, err = s.dao.OffsetCache(c, _bigdataEvent, msg.Partition); err != nil {
			log.Error("s.dao.OffsetCache(%d) error(%v)", msg.Partition, err)
			preOffset = 0
		} else {
			if msg.Offset > preOffset {
				s.setOffset(_bigdataEvent, msg.Partition, msg.Offset)
			}
		}
		eventMsg = &model.EventMessage{}
		if err = json.Unmarshal([]byte(msg.Value), eventMsg); err != nil {
			log.Error("json.Unmarshall(%s) error(%v)", msg.Value, err)
			s.setOffset(_bigdataEvent, msg.Partition, msg.Offset)
			continue
		}
		if msg.Offset <= preOffset && s.isMsgExpiration(eventMsg.Time) {
			log.Error("drop expired msg (%+v) (now_offset : %d)", preOffset)
			continue
		}
		if err = s.handleEvent(c, eventMsg); err != nil {
			log.Error("s.HandleEvent(%v) error(%v)", eventMsg, err)
			continue
		}
		log.Info("s.handleEvent(%v) eventMsg", eventMsg)
	}
}

func (s *Service) setOffset(event string, partition int32, offset int64) {
	s.cachemiss(func() {
		var err error
		if err = s.dao.SetOffsetCache(context.TODO(), event, partition, offset); err != nil {
			log.Error("s.dao.SetOffsetCache(%d,%d) error(%v)", partition, offset, err)
		}
	})
}

func (s *Service) secloginproc() {
	defer s.waiter.Done()
	defer func() {
		if x := recover(); x != nil {
			log.Error("eventproc unknown panic(%v)", x)
		}
	}()
	var (
		msg       *databus.Message
		eventMsg  *model.EventMessage
		ok        bool
		err       error
		msgChan   = s.secLoginDatabus.Messages()
		c         = context.TODO()
		preOffset int64
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("secloginproc msgChan closed")
				return
			}
		case <-s.quit:
			log.Info("quit secloginproc")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit error(%v)", err)
		}
		if preOffset, err = s.dao.OffsetCache(c, _secretEvent, msg.Partition); err != nil {
			log.Error("s.dao.OffsetCache(%d) error(%v)", msg.Partition, err)
			preOffset = 0
		} else {
			if msg.Offset > preOffset {
				s.setOffset(_secretEvent, msg.Partition, msg.Offset)
			}
		}
		eventMsg = &model.EventMessage{}
		if err = json.Unmarshal([]byte(msg.Value), eventMsg); err != nil {
			log.Error("json.Unmarshall(%s) error(%v)", msg.Value, err)
			s.setOffset(_secretEvent, msg.Partition, msg.Offset)
			continue
		}
		if msg.Offset <= preOffset && s.isMsgExpiration(eventMsg.Time) {
			log.Error("drop expired msg (%+v) (now_offset : %d)", preOffset)
			continue
		}
		if err = s.handleEvent(c, eventMsg); err != nil {
			log.Error("s.HandleEvent(%v) error(%v)", eventMsg, err)
			continue
		}
		log.Info("s.handleEvent(%v) eventMsg", eventMsg)
	}
}

func (s *Service) isMsgExpiration(timeStr string) bool {
	var (
		eventTime time.Time
		now       = time.Now()
		err       error
	)
	if eventTime, err = time.Parse("2006-01-02 15:04:05", timeStr); err != nil {
		return true
	}
	return eventTime.AddDate(0, 1, 1).Before(now)
}

func (s *Service) handleEvent(c context.Context, event *model.EventMessage) (err error) {
	var (
		eventTime time.Time
		argBytes  []byte
	)
	if eventTime, err = time.Parse("2006-01-02 15:04:05", event.Time); err != nil {
		log.Error("time.Parse(%s) errore(%v)", event.Time, err)
		return
	}
	if argBytes, err = json.Marshal(event.Args); err != nil {
		log.Error("json.Marshal(%v) error(%v), so empty it", event.Args, err)
		argBytes = []byte("{}")
	}
	var argHandleEvent = &cmmdl.ArgHandleEvent{
		Time:      xtime.Time(eventTime.Unix()),
		IP:        event.IP,
		Service:   event.Service,
		Event:     event.Event,
		ActiveMid: event.ActiveMid,
		TargetMid: event.TargetMid,
		TargetID:  event.TargetID,
		Args:      string(argBytes),
		Result:    event.Result,
		Effect:    event.Effect,
		RiskLevel: event.RiskLevel,
	}
	if err = s.spyRPC.HandleEvent(c, argHandleEvent); err != nil {
		log.Error("s.spyRPC.HandleEvent(%v) error(%v)", argHandleEvent, err)
		// s.retrymiss(func() {
		// 	log.Info("Start retry rpc error(%v)", err)
		// 	for {
		// 		if err = s.spyRPC.HandleEvent(context.TODO(), argHandleEvent); err != nil {
		// 			log.Error("s.spyRPC.HandleEvent(%v) error(%v)", argHandleEvent, err)
		// 		} else {
		// 			break
		// 		}
		// 	}
		// 	log.Info("End retry error(%v)", err)
		// })
		return
	}
	return
}

func (s *Service) dataproc() {
	for {
		select {
		case <-s.quit:
			log.Info("quit handleData")
			return
		default:
		}
		if !s.c.Property.Debug {
			d := time.Now().AddDate(0, 0, 1)
			ts := time.Date(d.Year(), d.Month(), d.Day(), 5, 0, 0, 0, time.Local).Sub(time.Now())
			time.Sleep(ts)
		} else {
			time.Sleep(s.c.Property.TaskTimer * time.Second)
		}
		s.reBuild()
	}
}

// Ping check service health.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.Ping(c); err != nil {
		log.Error("s.db.Ping() error(%v)", err)
		return
	}
	return
}

// Close all resource.
func (s *Service) Close() (err error) {
	close(s.quit)
	s.dao.Close()
	if err = s.eventDatabus.Close(); err != nil {
		log.Error("s.db.Close() error(%v)", err)
		return
	}
	return
}

// Wait wait all closed.
func (s *Service) Wait() {
	s.waiter.Wait()
}

func (s *Service) cacheproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.cacheproc panic(%v)", x)
			go s.cacheproc()
			log.Info("service.cacheproc recover")
		}
	}()
	for {
		f := <-s.cachech
		f()
	}
}

func (s *Service) retryproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.retryproc panic(%v)", x)
			go s.retryproc()
			log.Info("service.retryproc recover")
		}
	}()
	for {
		f := <-s.retrych
		go f()
	}
}

func (s *Service) cachemiss(f func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.cachemiss panic(%v)", x)
		}
	}()
	select {
	case s.cachech <- f:
	default:
		log.Error("service.cachech full")
	}
}

func (s *Service) retrymiss(f func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.retrymiss panic(%v)", x)
		}
	}()
	select {
	case s.retrych <- f:
	default:
		log.Error("service.retrych full")
	}
}

func (s *Service) loadconfig() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.cycleblock panic(%v)", x)
		}
	}()
	for {
		time.Sleep(s.configLoadTick)
		s.loadSystemConfig()
	}
}

func (s *Service) loadSystemConfig() (err error) {
	var (
		cdb map[string]string
	)
	cdb, err = s.dao.Configs(context.TODO())
	if err != nil {
		log.Error("sys config db get data err(%v)", err)
		return
	}
	if len(cdb) == 0 {
		err = errors.New("sys config no data")
		return
	}
	cs := make(map[string]interface{}, len(cdb))
	for k, v := range cdb {
		switch k {
		case model.LimitBlockCount:
			t, err1 := strconv.ParseInt(v, 10, 64)
			if err1 != nil {
				log.Error("sys config err(%s,%v,%v)", model.LimitBlockCount, t, err1)
				err = err1
				return
			}
			cs[k] = t
		case model.LessBlockScore:
			tmp, err1 := strconv.ParseInt(v, 10, 8)
			if err1 != nil {
				log.Error("sys config err(%s,%v,%v)", model.LessBlockScore, tmp, err1)
				err = err1
				return
			}
			cs[k] = int8(tmp)
		case model.AutoBlock:
			tmp, err1 := strconv.ParseInt(v, 10, 8)
			if err1 != nil {
				log.Error("sys config err(%s,%v,%v)", model.AutoBlock, tmp, err1)
				err = err1
				return
			}
			cs[k] = int8(tmp)
		default:
			cs[k] = v
		}
	}
	s.spyConfig = cs
	log.Info("loadSystemConfig success(%v)", cs)
	return
}

//Config get config.
func (s *Service) Config(key string) (interface{}, bool) {
	if s.spyConfig == nil {
		return nil, false
	}
	v, ok := s.spyConfig[key]
	return v, ok
}

func (s *Service) cycleblock() {
	var (
		c = context.TODO()
	)
	log.Info("cycleblock start (%v)", time.Now())
	if b, err := s.dao.SetNXLockCache(c, model.BlockLockKey, model.DefLockTime); !b || err != nil {
		log.Error("cycleblock had run (%v,%v)", b, err)
		return
	}
	s.BlockTask(c)
	s.dao.DelLockCache(c, model.BlockLockKey)
	log.Info("cycleblock end (%v)", time.Now())
}

func (s *Service) blockcacheuser() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.blockcacheuser panic(%v)", x)
			go s.blockcacheuser()
			log.Info("service.blockcacheuser recover")
		}
	}()
	for {
		mid, err := s.dao.SPOPBlockCache(context.TODO())
		if err != nil {
			log.Error("blockcacheuser err (%v,%v)", mid, err)
			continue
		}
		if mid != 0 {
			s.blockByMid(context.TODO(), mid)
			time.Sleep(s.blockTick)
		} else {
			// when no user should be block
			time.Sleep(s.blockWaitTick)
		}
	}
}

func (s *Service) reportjob() {
	var (
		c = context.TODO()
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.reportjob panic(%v)", x)
			go s.blockcacheuser()
			log.Info("service.reportjob recover")
		}
	}()
	log.Info("reportjob start (%v)", time.Now())
	if b, err := s.dao.SetNXLockCache(c, model.ReportJobKey, model.DefLockTime); !b || err != nil {
		log.Error("reportjob had run (%v,%v)", b, err)
		return
	}
	s.AddReport(c)
	log.Info("reportjob end (%v)", time.Now())
}

func (s *Service) spystatproc() {
	defer s.waiter.Done()
	defer func() {
		if x := recover(); x != nil {
			log.Error("sinstatproc unknown panic(%v)", x)
		}
	}()
	var (
		msg     *databus.Message
		statMsg *model.SpyStatMessage
		ok      bool
		err     error
		msgChan = s.spystatDatabus.Messages()
		c       = context.TODO()
	)
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("spystatproc msgChan closed")
				return
			}
		case <-s.quit:
			log.Info("quit spystatproc")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit error(%v)", err)
		}
		statMsg = &model.SpyStatMessage{}
		if err = json.Unmarshal([]byte(msg.Value), statMsg); err != nil {
			log.Error("json.Unmarshall(%s) error(%v)", msg.Value, err)
			continue
		}
		log.Info(" spystatproc (%v) start", statMsg)
		// check uuid
		unique, _ := s.dao.PfaddCache(c, statMsg.UUID)
		if !unique {
			log.Error("stat duplicate msg (%s) error", statMsg)
			continue
		}
		s.UpdateStatData(c, statMsg)
		log.Info(" spystatproc (%v) handle", statMsg)
	}
}

func (s *Service) loadeventname() (err error) {
	var (
		c  = context.Background()
		es []*model.Event
	)
	es, err = s.dao.AllEvent(c)
	if err != nil {
		log.Error("loadeventname allevent error(%v)", err)
		return
	}
	tmp := make(map[string]int64, len(es))
	for _, e := range es {
		tmp[e.Name] = e.ID
	}
	s.allEventName = tmp
	log.Info("loadeventname (%v) load success", tmp)
	return
}

func (s *Service) loadeventproc() {
	for {
		time.Sleep(s.loadEventTick)
		s.loadeventname()
	}
}

func (s *Service) initActivityEvents() {
	tmp := make(map[string]struct{}, len(s.c.Property.ActivityEvents))
	for _, name := range s.c.Property.ActivityEvents {
		tmp[name] = struct{}{}
	}
	s.activityEvents = tmp
}
