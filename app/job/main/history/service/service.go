package service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"go-common/app/job/main/history/conf"
	"go-common/app/job/main/history/dao"
	"go-common/app/job/main/history/model"
	hmdl "go-common/app/service/main/history/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline"
	"go-common/library/sync/pipeline/fanout"
	"go-common/library/xstr"

	"golang.org/x/time/rate"
)

const (
	_chanSize    = 1024
	_runtineSzie = 32
	_retryCnt    = 3
)

type message struct {
	next *message
	data *databus.Message
	done bool
}

// Service struct of service.
type Service struct {
	c             *conf.Config
	waiter        *sync.WaitGroup
	dao           *dao.Dao
	hisSub        *databus.Databus
	serviceHisSub *databus.Databus
	sub           *databus.Databus
	mergeChan     []chan *message
	doneChan      chan []*message
	merge         *pipeline.Pipeline
	businesses    map[int64]*hmdl.Business
	cache         *fanout.Fanout
	limit         *rate.Limiter
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:             c,
		dao:           dao.New(c),
		waiter:        new(sync.WaitGroup),
		hisSub:        databus.New(c.HisSub),
		serviceHisSub: databus.New(c.ServiceHisSub),
		sub:           databus.New(c.Sub),
		mergeChan:     make([]chan *message, _chanSize),
		doneChan:      make(chan []*message, _chanSize),
		cache:         fanout.New("cache"),
		limit:         rate.NewLimiter(rate.Limit(c.Job.QPSLimit), c.Job.ServiceBatch*2),
	}
	s.businesses = s.dao.BusinessesMap
	go s.subproc()
	go s.consumeproc()
	go s.serviceConsumeproc()
	go s.deleteproc()
	s.initMerge()
	for i := 0; i < _runtineSzie; i++ {
		c := make(chan *message, _chanSize)
		s.mergeChan[i] = c
		go s.mergeproc(c)
	}
	return
}

func (s *Service) consumeproc() {
	var (
		err        error
		n          int
		head, last *message
		msgs       = s.hisSub.Messages()
	)
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Error("s.consumeproc closed")
				return
			}
			// marked head to first commit
			m := &message{data: msg}
			if head == nil {
				head = m
				last = m
			} else {
				last.next = m
				last = m
			}
			if n, err = strconv.Atoi(msg.Key); err != nil {
				log.Error("strconv.Atoi(%s) error(%v)", msg.Key, err)
			}
			// use specify goruntine to flush
			s.mergeChan[n%_runtineSzie] <- m
			msg.Commit()
			log.Info("consumeproc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
		case done := <-s.doneChan:
			// merge partitions to commit offset
			commits := make(map[int32]*databus.Message)
			for _, d := range done {
				d.done = true
			}
			for ; head != nil && head.done; head = head.next {
				commits[head.data.Partition] = head.data
			}
			// for _, m := range commits {
			// m.Commit()
			// }
		}
	}
}

func (s *Service) mergeproc(c chan *message) {
	var (
		err    error
		max    = s.c.Job.Max
		merges = make(map[int64]int64, 10240)
		marked = make([]*message, 0, 10240)
		ticker = time.NewTicker(time.Duration(s.c.Job.Expire))
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("s.mergeproc closed")
				return
			}
			ms := make([]*model.Merge, 0, 32)
			if err = json.Unmarshal(msg.data.Value, &ms); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", msg.data.Value, err)
				continue
			}
			for _, m := range ms {
				if now, ok := merges[m.Mid]; !ok || now > m.Now {
					merges[m.Mid] = m.Now
				}
			}
			marked = append(marked, msg)
			if len(merges) < max {
				continue
			}
		case <-ticker.C:
		}
		if len(merges) > 0 {
			s.flush(merges)
			s.doneChan <- marked
			merges = make(map[int64]int64, 10240)
			marked = make([]*message, 0, 10240)
		}
	}
}

func (s *Service) flush(res map[int64]int64) {
	var (
		err   error
		ts    int64
		mids  []int64
		batch = s.c.Job.Batch
	)
	for mid, now := range res {
		if now < ts || ts == 0 {
			ts = now
		}
		mids = append(mids, mid)
	}
	for len(mids) > 0 {
		if len(mids) < batch {
			batch = len(mids)
		}
		for i := 0; i < _retryCnt; i++ {
			if err = s.dao.Flush(context.Background(), xstr.JoinInts(mids[:batch]), ts); err == nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		mids = mids[batch:]
	}
}

// Ping ping .
func (s *Service) Ping() {}

// Close .
func (s *Service) Close() {
	if s.sub != nil {
		s.sub.Close()
	}
	if s.serviceHisSub != nil {
		s.serviceHisSub.Close()
	}
	s.merge.Close()
	s.waiter.Wait()
}
