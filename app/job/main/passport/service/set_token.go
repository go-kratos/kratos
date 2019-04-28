package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/passport/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) userconsumeproc() {
	mergeRoutineNum := int64(s.c.Group.User.Num)
	msgs := s.dsUser.Messages()
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.userconsumeproc closed")
			return
		}
		// marked head to first commit
		m := &message{data: msg}
		p := new(model.PMsg)
		if err := json.Unmarshal(msg.Value, p); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}
		s.userMu.Lock()
		if s.userHead == nil {
			s.userHead = m
			s.userLast = m
		} else {
			s.userLast.next = m
			s.userLast = m
		}
		s.userMu.Unlock()
		m.object = p
		// use specify goroutine to merge messages
		s.userMergeChans[p.Data.Mid%mergeRoutineNum] <- m
		log.Info("userconsumeproc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) usercommitproc() {
	commits := make(map[int32]*databus.Message, s.c.Group.User.Size)
	for {
		done := <-s.userDoneChan
		// merge partitions to commit offset
		for _, d := range done {
			d.done = true
		}
		s.userMu.Lock()
		for ; s.userHead != nil && s.userHead.done; s.userHead = s.userHead.next {
			commits[s.userHead.data.Partition] = s.userHead.data
		}
		s.userMu.Unlock()
		for k, m := range commits {
			log.Info("usercommitproc committed, key:%s partition:%d offset:%d", m.Key, m.Partition, m.Offset)
			m.Commit()
			delete(commits, k)
		}
	}
}

func (s *Service) usermergeproc(c chan *message) {
	var (
		max    = s.c.Group.User.Size
		merges = make([]*model.PMsg, 0, max)
		marked = make([]*message, 0, max)
		ticker = time.NewTicker(time.Duration(s.c.Group.User.Ticker))
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("s.usermergeproc closed")
				return
			}
			p, assertOk := msg.object.(*model.PMsg)
			if assertOk && strings.HasPrefix(p.Table, "aso_app_perm") && p.Action != "" {
				merges = append(merges, p)
			}
			marked = append(marked, msg)
			if len(marked) < max && len(merges) < max {
				continue
			}
		case <-ticker.C:
		}
		if len(merges) > 0 {
			s.setTokens(merges)
			merges = make([]*model.PMsg, 0, max)
		}
		if len(marked) > 0 {
			s.userDoneChan <- marked
			marked = make([]*message, 0, max)
		}
	}
}

// setTokens for set tokens.
func (s *Service) setTokens(msgs []*model.PMsg) {
	for _, msg := range msgs {
		s.setToken(msg.Action, msg.Data)
	}
}

// setToken set single token.
func (s *Service) setToken(action string, t *model.Token) {
	if action == "" || t == nil || t.Token == "" {
		return
	}
	switch action {
	case "insert":
		for {
			if err := s.d.SetToken(context.TODO(), t); err == nil {
				return
			}
			time.Sleep(time.Second)
		}
	}
}
