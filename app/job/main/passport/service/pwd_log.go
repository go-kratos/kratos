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

type pwdLogBMsg struct {
	Action string
	Table  string
	New    *model.PwdLog
}

func (s *Service) pwdlogconsumeproc() {
	mergeRoutineNum := int64(s.c.Group.PwdLog.Num)
	for {
		msg, ok := <-s.dsPwdLog.Messages()
		if !ok {
			log.Error("s.pwdlogconsumeproc closed")
			return
		}
		// marked head to first commit
		m := &message{data: msg}
		p := &pwdLogBMsg{}
		if err := json.Unmarshal(msg.Value, p); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}
		// 只处理 aso_pwd_log insert binlog
		if p.Table != "aso_pwd_log" {
			continue
		}
		if p.Action != "insert" {
			continue
		}
		m.object = p
		s.pwdLogMu.Lock()
		if s.pwdLogHead == nil {
			s.pwdLogHead = m
			s.pwdLogLast = m
		} else {
			s.pwdLogLast.next = m
			s.pwdLogLast = m
		}
		s.pwdLogMu.Unlock()
		// use specify goroutine to merge messages
		s.pwdLogMergeChans[p.New.Mid%mergeRoutineNum] <- m
		log.Info("pwdlogconsumeproc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) pwdlogcommitproc() {
	commits := make(map[int32]*databus.Message, s.c.Group.PwdLog.Size)
	for {
		done := <-s.pwdLogDoneChan
		// merge partitions to commit offset
		for _, d := range done {
			d.done = true
		}
		s.pwdLogMu.Lock()
		for ; s.pwdLogHead != nil && s.pwdLogHead.done; s.pwdLogHead = s.pwdLogHead.next {
			commits[s.pwdLogHead.data.Partition] = s.pwdLogHead.data
		}
		s.pwdLogMu.Unlock()
		for k, m := range commits {
			log.Info("pwdlogcommitproc committed, key:%s partition:%d offset:%d", m.Key, m.Partition, m.Offset)
			m.Commit()
			delete(commits, k)
		}
	}
}

func (s *Service) pwdlogmergeproc(c chan *message) {
	var (
		max    = s.c.Group.PwdLog.Size
		merges = make([]*model.PwdLog, 0, max)
		marked = make([]*message, 0, max)
		ticker = time.NewTicker(time.Duration(s.c.Group.PwdLog.Ticker))
		err    error
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("s.pwdlogmergeproc closed")
				return
			}

			bmsg := &model.BMsg{}
			if err = json.Unmarshal(msg.data.Value, bmsg); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(msg.data.Value), err)
				continue
			}
			if bmsg.Action == "insert" && strings.HasPrefix(bmsg.Table, "aso_pwd_log") {
				p := &model.PwdLog{}
				if err = json.Unmarshal(bmsg.New, p); err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
					continue
				}
				merges = append(merges, p)
			}

			marked = append(marked, msg)
			if len(marked) < max && len(merges) < max {
				continue
			}
		case <-ticker.C:
		}
		if len(merges) > 0 {
			s.pwdlogprocessMerges(merges)
			merges = make([]*model.PwdLog, 0, max)
		}
		if len(marked) > 0 {
			s.logDoneChan <- marked
			marked = make([]*message, 0, max)
		}
	}
}

func (s *Service) pwdlogprocessMerges(merges []*model.PwdLog) {
	for _, v := range merges {
		for {
			res, err := s.d.GetPwdLog(context.Background(), v.ID)
			if err != nil {
				log.Error("fail to get pwd log, id(%d) err(%v)", v.ID, err)
				time.Sleep(_addHBaseRetryDuration)
				continue
			}
			if err := s.addPwdLog(context.Background(), res); err != nil {
				time.Sleep(_addHBaseRetryDuration)
				continue
			}
			break
		}

	}
}

func (s *Service) addPwdLog(c context.Context, v *model.PwdLog) (err error) {
	for i := 0; i < _addHBaseRetryCount; i++ {
		if err = s.d.AddPwdLogHBase(c, v); err == nil {
			return
		}
		log.Error("failed to add pwd log to hbase, service.dao.AddPwdLogHBase(%+v) error(%v)", v, err)
		time.Sleep(_addHBaseRetryDuration)
	}
	return
}
