package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/job/main/passport-game-local/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) asobinlogconsumeproc() {
	mergeNum := int64(s.c.Group.AsoBinLog.Num)
	for {
		msg, ok := <-s.dsAsoBinLogSub.Messages()
		if !ok {
			log.Error("asobinlogconsumeproc closed")
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
			t := new(model.AsoAccount)
			if err := json.Unmarshal(bmsg.New, t); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
			}
			mid = t.Mid
			m.object = bmsg
			log.Info("asobinlogconsumeproc table:%s key:%s partition:%d offset:%d", bmsg.Table, msg.Key, msg.Partition, msg.Offset)
		} else {
			continue
		}
		s.merges[mid%mergeNum] <- m
	}
}

func (s *Service) asobinlogcommitproc() {
	for {
		done := <-s.done
		commits := make(map[int32]*databus.Message)
		for _, d := range done {
			d.done = true
		}
		s.mu.Lock()
		for ; s.head != nil && s.head.done; s.head = s.head.next {
			commits[s.head.data.Partition] = s.head.data
		}
		s.mu.Unlock()
		for _, m := range commits {
			m.Commit()
		}
	}
}

func (s *Service) asobinlogmergeproc(c chan *message) {
	var (
		max    = s.c.Group.AsoBinLog.Size
		merges = make([]*model.BMsg, 0, max)
		marked = make([]*message, 0, max)
		ticker = time.NewTicker(time.Duration(s.c.Group.AsoBinLog.Ticker))
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("asobinlogmergeproc closed")
				return
			}
			p, assertOk := msg.object.(*model.BMsg)
			if assertOk && p.Action != "" && (p.Table == _asoAccountTable) {
				merges = append(merges, p)
			}
			marked = append(marked, msg)
			if len(marked) < max && len(merges) < max {
				continue
			}
		case <-ticker.C:
		}
		if len(merges) > 0 {
			s.processAsoAccLogInfo(merges)
			merges = make([]*model.BMsg, 0, max)
		}
		if len(marked) > 0 {
			s.done <- marked
			marked = make([]*message, 0, max)
		}
	}
}

func (s *Service) processAsoAccLogInfo(bmsgs []*model.BMsg) {
	for _, msg := range bmsgs {
		s.processAsoAccLog(msg)
	}
}

func (s *Service) processAsoAccLog(msg *model.BMsg) {
	aso := new(model.OriginAsoAccount)
	if err := json.Unmarshal(msg.New, aso); err != nil {
		log.Error("failed to parse binlog new, json.Unmarshal(%s) error(%v)", string(msg.New), err)
		return
	}
	pmsg := new(model.PMsg)
	if "update" == msg.Action {
		old := new(model.AsoAccount)
		if err := json.Unmarshal(msg.Old, old); err != nil {
			log.Error("failed to parse binlog new, json.Unmarshal(%s) error(%v)", string(msg.New), err)
			return
		}
		if old.Pwd != aso.Pwd {
			pmsg.Flag = 1
		}
	}
	pmsg.Action = msg.Action
	pmsg.Table = msg.Table
	pmsg.Data = model.Default(aso)
	key := strconv.FormatInt(aso.Mid, 10)
	for {
		if err := s.dsAsoEncryptTransPub.Send(context.TODO(), key, pmsg); err == nil {
			return
		}
		time.Sleep(time.Second)
	}
}
