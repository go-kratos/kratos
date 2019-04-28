package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/history/conf"
	"go-common/app/interface/main/history/dao/history"
	"go-common/app/interface/main/history/dao/toview"
	"go-common/app/interface/main/history/model"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	favrpc "go-common/app/service/main/favorite/api/gorpc"
	hisrpc "go-common/app/service/main/history/api/grpc"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"go-common/library/sync/pipeline/fanout"
)

type playPro struct {
	Type     int8   `json:"type"`
	SubType  int8   `json:"sub_type"`
	Mid      int64  `json:"mid"`
	Sid      int64  `json:"sid"`
	Epid     int64  `json:"epid"`
	Cid      int64  `json:"cid"`
	Progress int64  `json:"progress"`
	IP       string `json:"ip"`
	Ts       int64  `json:"ts"`
	RealTime int64  `json:"realtime"`
}

// Service is history service.
type Service struct {
	conf        *conf.Config
	historyDao  *history.Dao
	toviewDao   *toview.Dao
	delChan     *fanout.Fanout
	mergeChan   chan *model.Merge
	msgs        chan *playPro
	proChan     chan *model.History
	serviceChan chan func()
	favRPC      *favrpc.Service
	arcRPC      *arcrpc.Service2
	hisRPC      hisrpc.HistoryClient
	cache       *fanout.Fanout
	toviewCache *fanout.Fanout
	midMap      map[int64]bool
}

// New new a History service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:        c,
		historyDao:  history.New(c),
		toviewDao:   toview.New(c),
		mergeChan:   make(chan *model.Merge, 1024),
		msgs:        make(chan *playPro, 1024),
		proChan:     make(chan *model.History, 1024),
		serviceChan: make(chan func(), 10240),
		delChan:     fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		arcRPC:      arcrpc.New2(c.RPCClient2.Archive),
		favRPC:      favrpc.New2(c.RPCClient2.Favorite),
		cache:       fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		toviewCache: fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		midMap:      make(map[int64]bool),
	}
	for _, v := range s.conf.History.Mids {
		s.midMap[v] = true
	}
	var err error
	if s.hisRPC, err = hisrpc.NewClient(c.RPCClient2.History); err != nil {
		panic(err)
	}
	go s.playProproc()
	go s.mergeproc()
	go s.proPubproc()
	for i := 0; i < s.conf.History.ConsumeSize; i++ {
		go s.serviceproc()
	}
	return
}

func (s *Service) addMerge(mid, now int64) {
	select {
	case s.mergeChan <- &model.Merge{Mid: mid, Now: now}:
	default:
		log.Warn("mergeChan chan is full")
	}
}

func (s *Service) addPlayPro(p *playPro) {
	select {
	case s.msgs <- p:
	default:
		log.Warn("s.msgs chan is full")
	}
}

func (s *Service) addProPub(p *model.History) {
	select {
	case s.proChan <- p:
	default:
		log.Warn("s.proChan chan is full")
	}
}

func (s *Service) mergeproc() {
	var (
		m        *model.Merge
		ticker   = time.NewTicker(time.Duration(s.conf.History.Ticker))
		mergeMap = make(map[int64]int64)
	)
	for {
		select {
		case m = <-s.mergeChan:
			if m == nil {
				s.merge(mergeMap)
				return
			}
			if _, ok := mergeMap[m.Mid]; !ok {
				mergeMap[m.Mid] = m.Now
			}
			if len(mergeMap) < s.conf.History.Page {
				continue
			}
		case <-ticker.C:
		}
		s.merge(mergeMap)
		mergeMap = make(map[int64]int64)
	}
}

// playProproc send history to databus.
func (s *Service) playProproc() {
	var (
		msg    *playPro
		ms     []*playPro
		ticker = time.NewTicker(time.Second)
	)
	for {
		select {
		case msg = <-s.msgs:
			if msg == nil {
				if len(ms) > 0 {
					s.pushPlayPro(ms)
				}
				return
			}
			ms = append(ms, msg)
			if len(ms) < 100 {
				continue
			}
		case <-ticker.C:
		}
		if len(ms) == 0 {
			continue
		}
		s.pushPlayPro(ms)
		ms = make([]*playPro, 0, 100)
	}
}

func (s *Service) pushPlayPro(ms []*playPro) {
	key := fmt.Sprintf("%d%d", ms[0].Mid, ms[0].Sid)
	for j := 0; j < 3; j++ {
		if err := s.historyDao.PlayPro(context.Background(), key, ms); err == nil {
			return
		}
	}
}

// proPubroc send history to databus.
func (s *Service) proPubproc() {
	for {
		msg := <-s.proChan
		if msg == nil {
			return
		}
		s.proPub(msg)
	}
}

func (s *Service) proPub(msg *model.History) {
	key := fmt.Sprintf("%d%d", msg.Mid, msg.Aid)
	for j := 0; j < 3; j++ {
		if err := s.historyDao.ProPub(context.Background(), key, msg); err == nil {
			break
		}
	}
}

func (s *Service) userActionLog(mid int64, action string) {
	report.User(&report.UserInfo{
		Mid:      mid,
		Business: model.HistoryLog,
		Action:   action,
		Ctime:    time.Now(),
	})
}

func (s *Service) migration(mid int64) bool {
	if !s.conf.History.Migration || mid == 0 {
		return false
	}
	if _, ok := s.midMap[mid]; ok {
		return true
	}
	if s.conf.History.Rate != 0 && mid%s.conf.History.Rate == 0 {
		return true
	}
	return false
}

// Ping ping service.
// +wd:ignore
func (s *Service) Ping(c context.Context) (err error) {
	if s.historyDao != nil {
		err = s.historyDao.Ping(c)
	}
	if s.toviewDao != nil {
		err = s.toviewDao.Ping(c)
	}
	return
}

// Close close resource.
// +wd:ignore
func (s *Service) Close() {
	s.mergeChan <- nil
	s.msgs <- nil
	s.proChan <- nil
	if s.historyDao != nil {
		s.historyDao.Close()
	}
	if s.toviewDao != nil {
		s.toviewDao.Close()
	}
}
