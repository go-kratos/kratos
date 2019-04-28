package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go-common/app/job/main/stat/conf"
	"go-common/app/job/main/stat/dao"
	"go-common/app/job/main/stat/model"
	arcmdl "go-common/app/service/main/archive/api"
	archive "go-common/app/service/main/archive/api/gorpc"
	"go-common/library/cache/memcache"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_sharding = 100
)

type lastTmStat struct {
	last int64
	stat *arcmdl.Stat
}

// Service is stat job service.
type Service struct {
	c *conf.Config
	// dao
	dao *dao.Dao
	// wait
	waiter sync.WaitGroup
	closed bool
	// databus
	subMap     map[string]*databus.Databus
	subMonitor map[string]*model.Monitor
	subStatCh  []chan *model.StatMsg
	mu         sync.Mutex
	// stat map
	statSM []map[int64]*lastTmStat
	// rpc
	arcRPC  *archive.Service2
	arcRPC2 *archive.Service2
	// max aid
	maxAid    int64
	memcaches []*memcache.Pool
}

// New is stat-job service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		// dao
		dao: dao.New(c),
		// rpc
		arcRPC:     archive.New2(c.ArchiveRPC),
		arcRPC2:    archive.New2(c.ArchiveRPC2),
		subMap:     make(map[string]*databus.Databus),
		subMonitor: make(map[string]*model.Monitor),
	}
	for _, mc := range s.c.Memcaches {
		s.memcaches = append(s.memcaches, memcache.NewPool(mc))
	}
	// view
	s.subMap[model.TypeForView] = databus.New(c.ViewSub)
	s.subMonitor[model.TypeForView] = &model.Monitor{Topic: c.ViewSub.Topic, Count: 0}
	// dm
	s.subMap[model.TypeForDm] = databus.New(c.DmSub)
	s.subMonitor[model.TypeForDm] = &model.Monitor{Topic: c.DmSub.Topic, Count: 0}
	// reply
	s.subMap[model.TypeForReply] = databus.New(c.ReplySub)
	s.subMonitor[model.TypeForReply] = &model.Monitor{Topic: c.ReplySub.Topic, Count: 0}
	// fav
	s.subMap[model.TypeForFav] = databus.New(c.FavSub)
	s.subMonitor[model.TypeForFav] = &model.Monitor{Topic: c.FavSub.Topic, Count: 0}
	// coin
	s.subMap[model.TypeForCoin] = databus.New(c.CoinSub)
	s.subMonitor[model.TypeForCoin] = &model.Monitor{Topic: c.CoinSub.Topic, Count: 0}
	// share
	s.subMap[model.TypeForShare] = databus.New(c.ShareSub)
	s.subMonitor[model.TypeForShare] = &model.Monitor{Topic: c.ShareSub.Topic, Count: 0}
	// rank
	s.subMap[model.TypeForRank] = databus.New(c.RankSub)
	// like
	s.subMap[model.TypeForLike] = databus.New(c.LikeSub)
	s.subMonitor[model.TypeForLike] = &model.Monitor{Topic: c.LikeSub.Topic, Count: 0}
	for i := int64(0); i < _sharding; i++ {
		s.subStatCh = append(s.subStatCh, make(chan *model.StatMsg, 10240))
		s.statSM = append(s.statSM, map[int64]*lastTmStat{})
		s.waiter.Add(1)
		go s.statDealproc(i)
	}
	go s.loadproc()
	if env.DeployEnv == env.DeployEnvProd {
		go s.monitorproc()
	}
	for k, d := range s.subMap {
		s.waiter.Add(1)
		go s.consumerproc(k, d)
	}
	return
}

func (s *Service) loadproc() {
	for {
		time.Sleep(1 * time.Minute)
		id, err := s.dao.MaxAID(context.TODO())
		if err != nil {
			s.maxAid = 0
			log.Error("s.dao.MaxAid error(%+v)", err)
			continue
		}
		s.maxAid = id
	}
}

func (s *Service) monitorproc() {
	for {
		time.Sleep(90 * time.Second)
		s.mu.Lock()
		for _, mo := range s.subMonitor {
			if mo.Count == 0 {
				s.dao.SendQiyeWX(fmt.Sprintf("日志报警:stat-job topic(%s) 没消费！！！！", mo.Topic))
			}
			mo.Count = 0
		}
		s.mu.Unlock()
	}
}

// consumerproc consumer all topic
func (s *Service) consumerproc(k string, d *databus.Databus) {
	defer s.waiter.Done()
	var msgs = d.Messages()
	for {
		var (
			err error
			ok  bool
			msg *databus.Message
			now = time.Now().Unix()
		)
		msg, ok = <-msgs
		if !ok || s.closed {
			log.Info("databus(%s) consumer exit", k)
			return
		}
		msg.Commit()
		var ms = &model.StatCount{}
		if err = json.Unmarshal(msg.Value, ms); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}
		if ms.Aid <= 0 || (ms.Type != "archive" && ms.Type != "archive_his") {
			log.Warn("message(%s) error", msg.Value)
			continue
		}
		if now-ms.TimeStamp > 8*60*60 {
			log.Warn("topic(%s) message(%s) too early", msg.Topic, msg.Value)
			continue
		}
		stat := &model.StatMsg{Aid: ms.Aid, Type: k, Ts: ms.TimeStamp}
		switch k {
		case model.TypeForView:
			stat.Click = ms.Count
		case model.TypeForDm:
			stat.DM = ms.Count
		case model.TypeForReply:
			stat.Reply = ms.Count
		case model.TypeForFav:
			stat.Fav = ms.Count
		case model.TypeForCoin:
			stat.Coin = ms.Count
		case model.TypeForShare:
			stat.Share = ms.Count
		case model.TypeForRank:
			stat.HisRank = ms.Count
		case model.TypeForLike:
			stat.Like = ms.Count
			stat.DisLike = ms.DisLike
		default:
			log.Error("unknow type(%s) message(%s)", k, msg.Value)
			continue
		}
		s.mu.Lock()
		if _, ok := s.subMonitor[k]; ok {
			s.subMonitor[k].Count++
		}
		s.mu.Unlock()
		s.subStatCh[stat.Aid%_sharding] <- stat
		log.Info("got message(%+v)", stat)
	}
}

// Close Databus consumer close.
func (s *Service) Close() (err error) {
	s.closed = true
	time.Sleep(2 * time.Second)
	log.Info("start close job")
	for k, d := range s.subMap {
		d.Close()
		log.Info("databus(%s) cloesed", k)
	}
	for i := int64(0); i < _sharding; i++ {
		close(s.subStatCh[i])
	}
	log.Info("end close job")
	s.waiter.Wait()
	return
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}
