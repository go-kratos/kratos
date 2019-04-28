package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/passport-game-cloud/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) encrypttransconsumeproc() {
	var (
		mergeRoutineNum = int64(s.c.Group.EncryptTrans.Num)
		msgs            = s.encryptTransDataBus.Messages()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("encrypttransconsumeproc closed")
			return
		}
		// marked head to first commit
		m := &message{data: msg}
		s.asoMu.Lock()
		if s.asoHead == nil {
			s.asoHead = m
			s.asoLast = m
		} else {
			s.asoLast.next = m
			s.asoLast = m
		}
		s.asoMu.Unlock()
		p := new(model.PMsg)
		if err := json.Unmarshal(msg.Value, p); err != nil {
			log.Error("encrypttransconsumeproc unmarshal failed, json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if p.Table == _asoAccountTable {
			p.MTS = s.transInterval.MTS(context.TODO(), p.Data.Mtime)
		}
		m.object = p
		s.encryptTransMergeChans[p.Data.Mid%mergeRoutineNum] <- m
		log.Info("encrypttransconsumeproc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) encrypttranscommitproc() {
	commits := make(map[int32]*databus.Message, s.c.Group.EncryptTrans.Size)
	for {
		done := <-s.encryptTransDoneChan
		for _, d := range done {
			d.done = true
		}
		s.asoMu.Lock()
		for ; s.asoHead != nil && s.asoHead.done; s.asoHead = s.asoHead.next {
			commits[s.asoHead.data.Partition] = s.asoHead.data
		}
		s.asoMu.Unlock()
		for k, m := range commits {
			log.Info("encrypttranscommitproc committed, key:%s partition:%d offset:%d", m.Key, m.Partition, m.Offset)
			m.Commit()
			delete(commits, k)
		}
	}
}

func (s *Service) encrypttransmergeproc(c chan *message) {
	var (
		max    = s.c.Group.EncryptTrans.Size
		merges = make([]*model.PMsg, 0, max)
		marked = make([]*message, 0, max)
		ticker = time.NewTicker(time.Duration(s.c.Group.EncryptTrans.Ticker))
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("encrypttransmergeproc closed")
				return
			}
			p, assertOK := msg.object.(*model.PMsg)
			if assertOK {
				merges = append(merges, p)
			}
			marked = append(marked, msg)
			if len(marked) < max && len(merges) < max {
				continue
			}
		case <-ticker.C:
		}
		if len(merges) > 0 {
			s.processAsoAcc(merges)
			merges = make([]*model.PMsg, 0, max)
		}
		if len(marked) > 0 {
			s.encryptTransDoneChan <- marked
			marked = make([]*message, 0, max)
		}
	}
}

func (s *Service) processAsoAcc(pmsgs []*model.PMsg) {
	for _, p := range pmsgs {
		if p.Action != "" && p.Table == _asoAccountTable {
			s.processAsoAccSub(p)
		}
	}
}
