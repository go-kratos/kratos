package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/passport-game-cloud/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) binlogconsumeproc() {
	mergeRoutineNum := int64(s.c.Group.BinLog.Num)
	for {
		msg, ok := <-s.binLogDataBus.Messages()
		if !ok {
			log.Error("binlogconsumeproc closed")
			return
		}
		// marked head to first commit
		m := &message{data: msg}
		s.mu.Lock()
		if s.head == nil {
			s.head = m
			s.last = m
		} else {
			s.last.next = m
			s.last = m
		}
		s.mu.Unlock()
		bmsg := new(model.BMsg)
		if err := json.Unmarshal(msg.Value, bmsg); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}
		mid := int64(0)
		if bmsg.Table == _asoAccountTable {
			t := new(model.OriginAsoAccount)
			if err := json.Unmarshal(bmsg.New, t); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
				continue
			}
			mid = t.Mid
			bmsg.MTS = s.asoAccountInterval.MTS(context.TODO(), t.Mtime)
		} else if strings.HasPrefix(bmsg.Table, _tokenTablePrefix) {
			t := new(model.OriginPerm)
			if err := json.Unmarshal(bmsg.New, t); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
				continue
			}
			mid = t.Mid
			bmsg.MTS = s.tokenInterval.MTS(context.TODO(), t.Mtime)
		} else if strings.HasPrefix(bmsg.Table, _memberTablePrefix) {
			t := new(model.OriginMember)
			if err := json.Unmarshal(bmsg.New, t); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
				continue
			}
			mid = t.Mid
			bmsg.MTS = s.memberInterval.MTS(context.TODO(), t.Mtime)
		}
		m.object = bmsg
		// use specify goroutine to merge messages
		s.binLogMergeChans[mid%mergeRoutineNum] <- m
		log.Info("binlogconsumeproc table:%s key:%s partition:%d offset:%d", bmsg.Table, msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) binlogcommitproc() {
	commits := make(map[int32]*databus.Message, s.c.Group.BinLog.Size)
	for {
		done := <-s.binLogDoneChan
		// merge partitions to commit offset
		for _, d := range done {
			d.done = true
		}
		s.mu.Lock()
		for ; s.head != nil && s.head.done; s.head = s.head.next {
			commits[s.head.data.Partition] = s.head.data
		}
		s.mu.Unlock()
		for k, m := range commits {
			log.Info("binlogcommitproc committed, key:%s partition:%d offset:%d", m.Key, m.Partition, m.Offset)
			m.Commit()
			delete(commits, k)
		}
	}
}

func (s *Service) binlogmergeproc(c chan *message) {
	var (
		max    = s.c.Group.BinLog.Size
		merges = make([]*model.BMsg, 0, max)
		marked = make([]*message, 0, max)
		ticker = time.NewTicker(time.Duration(s.c.Group.BinLog.Ticker))
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("binlogmergeproc closed")
				return
			}
			p, assertOk := msg.object.(*model.BMsg)
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
			s.process(merges)
			merges = make([]*model.BMsg, 0, max)
		}
		if len(marked) > 0 {
			s.binLogDoneChan <- marked
			marked = make([]*message, 0, max)
		}
	}
}

func (s *Service) process(bmsgs []*model.BMsg) {
	for _, bmsg := range bmsgs {
		if bmsg.Table == _asoAccountTable {
			s.processUserInfo(bmsg)
		} else if strings.HasPrefix(bmsg.Table, _tokenTablePrefix) {
			s.processToken(bmsg)
		} else if strings.HasPrefix(bmsg.Table, _memberTablePrefix) {
			s.processMemberInfo(bmsg)
		}
	}
}
