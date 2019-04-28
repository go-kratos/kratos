package service

import (
	"context"
	"encoding/json"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"go-common/app/job/main/tag/conf"
	"go-common/app/job/main/tag/dao"
	"go-common/app/job/main/tag/model"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/cache"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_updateArc = "update"
	_insertArc = "insert"
	_archive   = "archive"
)

// Service service .
type Service struct {
	conf          *conf.Config
	waiter        *sync.WaitGroup
	dao           *dao.Dao
	conChan       []chan *model.Message
	arcSub        *databus.Databus
	tagSub        *databus.Databus
	rpMap         map[int64]int64
	coreNumber    int
	arcRPC        *arcrpc.Service2
	businessCache atomic.Value
	cacheCh       *cache.Cache
}

// New new .
func New(c *conf.Config) (s *Service) {
	coreNumber := runtime.NumCPU()
	s = &Service{
		conf:       c,
		waiter:     new(sync.WaitGroup),
		dao:        dao.New(c),
		conChan:    make([]chan *model.Message, coreNumber),
		arcSub:     databus.New(c.Databus.Archive),
		tagSub:     databus.New(c.Databus.Tag),
		rpMap:      make(map[int64]int64),
		coreNumber: coreNumber,
		arcRPC:     arcrpc.New2(c.ArchiveRPC),
		cacheCh:    cache.New(1, 1024),
	}
	if err := s.businessCaches(); err != nil {
		panic(err)
	}
	s.waiter.Add(1)
	go s.subProc()
	for i := 0; i < s.coreNumber; i++ {
		s.waiter.Add(1)
		ch := make(chan *model.Message, 1024)
		s.conChan[i] = ch
		go s.jobProc(ch)
	}
	go s.ramProc()
	go s.hotTagNewArcProc()
	go s.writeTagInfoproc()
	go s.businessCacheproc()
	go s.resTagActionConsumeproc()
	return
}

func (s *Service) subProc() {
	defer s.waiter.Done()
	for {
		msg, ok := <-s.arcSub.Messages()
		if !ok {
			log.Info("subProc exit")
			for i := 0; i < s.coreNumber; i++ {
				log.Info("close channel(%d)", i)
				close(s.conChan[i])
			}
			return
		}
		msg.Commit()
		m := &model.Message{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal() error(%v)", err)
			continue
		}
		switch m.Table {
		case _archive:
			newArc, err := s.parseNewMsg(m)
			if err != nil {
				continue
			}
			s.conChan[newArc.Aid%int64(s.coreNumber)] <- m
		}
	}
}

func (s *Service) jobProc(ch chan *model.Message) {
	defer s.waiter.Done()
	for {
		m, ok := <-ch
		ctx := context.TODO()
		if !ok {
			log.Info("jobProc exit")
			return
		}
		switch m.Action {
		case _insertArc:
			newArc, err := s.parseNewMsg(m)
			if err != nil {
				continue
			}
			s.insertArcCache(ctx, newArc)
		case _updateArc:
			newArc, oldArc, err := s.parseAllMsg(m)
			if err != nil {
				continue
			}
			s.upTagArcCache(ctx, newArc, oldArc)
		default:
		}
	}

}

// Close close .
func (s *Service) Close() (err error) {
	s.dao.Close()
	if err = s.arcSub.Close(); err != nil {
		log.Error("s.arcSub.Close() error(%v)", err)
		return
	}
	if err = s.tagSub.Close(); err != nil {
		log.Error("s.tagSub.Close() error(%v)", err)
	}
	return
}

// Wait wait chan chan .
func (s *Service) Wait() {
	s.waiter.Wait()
}

// Ping ping db .
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

func (s *Service) ramProc() {
	for {
		if rpMap, err := s.dao.Rids(context.Background()); err == nil && len(rpMap) > 0 {
			s.rpMap = rpMap
		} else {
			log.Error("d.rids() error(%v) len rpMap(%d)", err, len(rpMap))
		}
		time.Sleep(time.Duration(s.conf.Tag.Tick))
	}
}

func (s *Service) hotTagNewArcProc() {
	for {
		ridTagMap, err := s.dao.HotMap(context.Background())
		if err != nil || len(ridTagMap) == 0 {
			continue
		}
		for rid, tids := range ridTagMap {
			for _, tid := range tids {
				oids, err := s.dao.TagResources(context.Background(), tid)
				if err != nil {
					continue
				}
				s.batchArchivesDelay(context.Background(), tid, oids)
			}
			log.Info("hotTagNewArcProc() finsh rid:%d", rid)
		}
		time.Sleep(time.Minute * 10)
	}
}

const _archiveInterval = 50

// batchArchivesDelay time delay  .
func (s *Service) batchArchivesDelay(c context.Context, tid int64, aids []int64) (err error) {
	var (
		tmpRes map[int64]*api.Arc
		n      = _archiveInterval
	)
	for len(aids) > 0 {
		if n > len(aids) {
			n = len(aids)
		}
		arg := &archive.ArgAids2{Aids: aids[:n]}
		aids = aids[n:]
		if tmpRes, err = s.arcRPC.Archives3(c, arg); err != nil {
			log.Error("s.arcRPC.Archives3(%v) error(%v)", arg.Aids, err)
			err = nil
			continue
		}
		normalArc := make(map[int64]*api.Arc)
		for k, v := range tmpRes {
			if v.IsNormal() {
				normalArc[k] = v
			}
		}
		if len(normalArc) > 0 {
			s.dao.UpdateTagNewArcCache(context.Background(), tid, normalArc)
		}
		time.Sleep(time.Millisecond * 50)
	}
	return
}
