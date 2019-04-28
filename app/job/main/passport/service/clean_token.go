package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/passport/model"
	igmdl "go-common/app/service/main/identify-game/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_changePwd = "changePwd"

	_retryCount    = 3
	_retryDuration = time.Second
)

func (s *Service) tokenconsumeproc() {
	mergeNum := s.c.Group.AsoBinLog.Num
	var (
		err  error
		n    int
		msgs = s.dsToken.Messages()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.tokenconsumeproc closed")
			return
		}
		// marked head to first commit
		m := &message{data: msg}
		if n, err = strconv.Atoi(msg.Key); err != nil {
			log.Error("strconv.Atoi(%s) error(%v)", msg.Key, err)
			continue
		}
		s.mu.Lock()
		if s.head == nil {
			s.head = m
			s.last = m
		} else {
			s.last.next = m
			s.last = m
		}
		s.mu.Unlock()
		// use specify goroutine to merge messages
		s.tokenMergeChans[n%mergeNum] <- m
		log.Info("tokenconsumeproc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) tokencommitproc() {
	commits := make(map[int32]*databus.Message, s.c.Group.AsoBinLog.Size)
	for {
		done := <-s.tokenDoneChan
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
			log.Info("tokencommitproc committed, key:%s partition:%d offset:%d", m.Key, m.Partition, m.Offset)
			m.Commit()
			delete(commits, k)
		}
	}
}

func (s *Service) tokenmergeproc(c chan *message) {
	var (
		err    error
		max    = s.c.Group.AsoBinLog.Size
		merges = make([]*model.AccessInfo, 0, max)
		marked = make([]*message, 0, max)
		ticker = time.NewTicker(time.Duration(s.c.Group.AsoBinLog.Ticker))
	)
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				log.Error("s.tokenmergeproc closed")
				return
			}
			bmsg := &model.BMsg{}
			if err = json.Unmarshal(msg.data.Value, bmsg); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(msg.data.Value), err)
				continue
			}
			if bmsg.Action == "delete" && strings.HasPrefix(bmsg.Table, "aso_app_perm") {
				t := &model.AccessInfo{}
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
			s.cleanTokens(merges)
			merges = make([]*model.AccessInfo, 0, max)
		}
		if len(marked) > 0 {
			s.tokenDoneChan <- marked
			marked = make([]*message, 0, max)
		}
	}
}

// cleanTokens clean tokens.
func (s *Service) cleanTokens(tokens []*model.AccessInfo) {
	for _, token := range tokens {
		s.cleanToken(token)
	}
}

// cleanToken to notify other clean access token.
func (s *Service) cleanToken(token *model.AccessInfo) (err error) {
	if token == nil || token.Expires < time.Now().Unix() {
		return
	}
	isGame := false
	for _, id := range s.gameAppIDs {
		if id == token.AppID {
			isGame = true
			break
		}
	}
	if !isGame {
		return
	}
	for {
		if err = s.d.DelCache(context.TODO(), token.Token); err == nil {
			break
		}
		time.Sleep(_retryDuration)
	}
	for i := 0; i < _retryCount; i++ {
		arg := &igmdl.CleanCacheArgs{
			Token: token.Token,
			Mid:   token.Mid,
		}
		if err = s.igRPC.DelCache(context.TODO(), arg); err == nil {
			break
		}
		log.Error("service.identifyGameRPC.DelCache(%+v) error(%v)", arg, err)
		time.Sleep(_retryDuration)
	}
	for i := 0; i < _retryCount; i++ {
		if err = s.d.NotifyGame(token, _changePwd); err == nil {
			return
		}
		time.Sleep(_retryDuration)
	}
	log.Error("notify err, token(%+v)", token)
	return
}
