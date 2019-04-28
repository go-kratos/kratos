package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/passport/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_addHBaseRetryCount    = 3
	_addHBaseRetryDuration = time.Second
)

func (s *Service) logconsumeproc() {
	mergeRoutineNum := int64(s.c.Group.Log.Num)
	for {
		msg, ok := <-s.dsLog.Messages()
		if !ok {
			log.Error("s.logconsumeproc closed")
			return
		}
		// marked head to first commit
		m := &message{data: msg}
		p := &model.LoginLog{}
		if err := json.Unmarshal(msg.Value, p); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}
		s.logMu.Lock()
		if s.logHead == nil {
			s.logHead = m
			s.logLast = m
		} else {
			s.logLast.next = m
			s.logLast = m
		}
		s.logMu.Unlock()
		m.object = p
		// use specify goroutine to merge messages
		s.logMergeChans[p.Mid%mergeRoutineNum] <- m
		log.Info("logconsumeproc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) logcommitproc() {
	commits := make(map[int32]*databus.Message, s.c.Group.Log.Size)
	for {
		done := <-s.logDoneChan
		// merge partitions to commit offset
		for _, d := range done {
			d.done = true
		}
		s.logMu.Lock()
		for ; s.logHead != nil && s.logHead.done; s.logHead = s.logHead.next {
			commits[s.logHead.data.Partition] = s.logHead.data
		}
		s.logMu.Unlock()
		for k, m := range commits {
			log.Info("logcommitproc committed, key:%s partition:%d offset:%d", m.Key, m.Partition, m.Offset)
			m.Commit()
			delete(commits, k)
		}
	}
}

func (s *Service) logmergeproc(c chan *message) {
	var (
		max    = s.c.Group.Log.Size
		merges = make([]*model.LoginLog, 0, max)
		marked = make([]*message, 0, max)
		ticker = time.NewTicker(time.Duration(s.c.Group.Log.Ticker))
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("s.logmergeproc closed")
				return
			}
			p, assertOk := msg.object.(*model.LoginLog)
			if assertOk {
				merges = append(merges, p)
			}
			marked = append(marked, msg)
			if len(marked) < max && len(merges) < max {
				continue
			}
		case <-ticker.C:
		}
		if len(merges) > 0 {
			s.processMerges(merges)
			merges = make([]*model.LoginLog, 0, max)
		}
		if len(marked) > 0 {
			s.logDoneChan <- marked
			marked = make([]*message, 0, max)
		}
	}
}

func (s *Service) processMerges(merges []*model.LoginLog) {
	s.d.AddLoginLog(merges)
	for _, v := range merges {
		s.addLoginLog(context.TODO(), v)
	}
}

func (s *Service) addLoginLog(c context.Context, v *model.LoginLog) (err error) {
	for i := 0; i < _addHBaseRetryCount; i++ {
		if err = s.d.AddLoginLogHBase(c, v); err == nil {
			return
		}
		log.Error("failed to add login log to hbase, service.dao.AddLoginLogHBase(%+v) error(%v)", v, err)
		time.Sleep(_addHBaseRetryDuration)
	}
	return
}
