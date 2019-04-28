package service

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/passport/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) authBinLogconsumeproc() {
	mergeNum := s.c.Group.AuthBinLog.Num
	var (
		err  error
		n    int
		msgs = s.authBinLog.Messages()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.authBinLogconsumeproc closed")
			return
		}
		// marked head to first commit
		m := &message{data: msg}
		if n, err = strconv.Atoi(msg.Key); err != nil {
			log.Error("strconv.Atoi(%s) error(%v)", msg.Key, err)
			continue
		}
		s.authBinLogMu.Lock()
		if s.authBinLogHead == nil {
			s.authBinLogHead = m
			s.authBinLogLast = m
		} else {
			s.authBinLogLast.next = m
			s.authBinLogLast = m
		}
		s.authBinLogMu.Unlock()
		// use specify goroutine to merge messages
		s.authBinLogMergeChans[n%mergeNum] <- m
		log.Info("authBinLogconsumeproc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) authBinLogcommitproc() {
	commits := make(map[int32]*databus.Message, s.c.Group.AuthBinLog.Size)
	for {
		done := <-s.authBinLogDoneChan
		// merge partitions to commit offset
		for _, d := range done {
			d.done = true
		}
		s.mu.Lock()
		for ; s.authBinLogHead != nil && s.authBinLogHead.done; s.authBinLogHead = s.authBinLogHead.next {
			commits[s.authBinLogHead.data.Partition] = s.authBinLogHead.data
		}
		s.mu.Unlock()
		for k, m := range commits {
			log.Info("authBinLogcommitproc committed, key:%s partition:%d offset:%d", m.Key, m.Partition, m.Offset)
			m.Commit()
			delete(commits, k)
		}
	}
}

func (s *Service) authBinLogmergeproc(c chan *message) {
	var (
		err    error
		max    = s.c.Group.AuthBinLog.Size
		merges = make([]*model.AuthToken, 0, max)
		marked = make([]*message, 0, max)
		ticker = time.NewTicker(time.Duration(s.c.Group.AuthBinLog.Ticker))
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("s.authBinLogmergeproc closed")
				return
			}
			bmsg := &model.BMsg{}
			if err = json.Unmarshal(msg.data.Value, bmsg); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(msg.data.Value), err)
				continue
			}
			if bmsg.Action == "delete" && strings.HasPrefix(bmsg.Table, "user_token_") {
				t := &model.AuthToken{}
				if err = json.Unmarshal(bmsg.New, t); err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
					continue
				}
				merges = append(merges, t)
			}
			marked = append(marked, msg)
			if len(marked) < max && len(merges) < max {
				continue
			}
		case <-ticker.C:
		}
		if len(merges) > 0 {
			s.cleanAuthTokens(merges)
			merges = make([]*model.AuthToken, 0, max)
		}
		if len(marked) > 0 {
			s.authBinLogDoneChan <- marked
			marked = make([]*message, 0, max)
		}
	}
}

// cleanTokens by auth .
func (s *Service) cleanAuthTokens(authTokens []*model.AuthToken) {
	for _, authToken := range authTokens {
		var (
			bytes []byte
			err   error
		)
		if bytes, err = base64.StdEncoding.DecodeString(authToken.Token); err != nil {
			log.Error("cleanAuthTokens base64 decode err %v", err)
			continue
		}
		token := hex.EncodeToString(bytes)
		log.Info("auth binlog clear cleanAuthTokens,msg is (%+v)", authToken)
		t := &model.AccessInfo{
			Mid:     authToken.Mid,
			AppID:   int32(authToken.AppID),
			Token:   token,
			Expires: authToken.Expires,
		}
		s.cleanToken(t)
	}
}
