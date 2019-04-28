package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/app/job/main/videoup-report/model/email"
	"go-common/library/log"
)

func (s *Service) addRetry(c context.Context, aid int64, action string, flag int64, flagA int64) (err error) {
	member := email.Retry{
		AID:        aid,
		Action:     action,
		Flag:       flag,
		FlagA:      flagA,
		CreateTime: time.Now().Unix(),
	}

	err = s.email.PushRedis(c, member, email.RetryListKey)
	return
}

//popRetry spop a retry element from redis set
func (s *Service) popRetry(c context.Context, key string) (member email.Retry, err error) {
	var bs []byte
	if bs, err = s.email.PopRedis(c, key); err != nil {
		return
	}

	if err = json.Unmarshal(bs, &member); err != nil {
		log.Error("PopRetry(%s) json.Unmarshal(%s) error(%v)", key, string(bs), err)
	}
	return
}

func (s *Service) removeRetry(c context.Context, aid int64, action string) (err error) {
	var (
		bs    []byte
		list  []interface{}
		reply int
	)

	for _, snew := range archive.ReplyState {
		for _, sold := range archive.ReplyState {
			member := email.Retry{
				AID:    aid,
				Action: action,
				Flag:   snew,
				FlagA:  sold,
			}
			if bs, err = json.Marshal(member); err != nil {
				log.Error("removeRetry json.Marshal error(%v) member(%+v)", err, member)
				err = nil
				continue
			}
			list = append(list, string(bs))
		}
	}

	if reply, err = s.email.RemoveRedis(c, email.RetryListKey, list...); err != nil {
		log.Error("removeRetry s.email.RemoveRedis error(%v) aid(%d) action(%s)", err, aid, action)
	} else {
		log.Info("removeRetry s.email.RemoveRedis success reply(%d) aid(%d) action(%s)", reply, aid, action)
	}
	return
}

func (s *Service) retryProc() {
	defer s.waiter.Done()
	for {
		if s.closed {
			return
		}

		c := context.TODO()
		member, err := s.popRetry(c, email.RetryListKey)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		log.Info("retry member(%+v)", member)

		if member.CreateTime > 0 && (time.Now().Unix()-member.CreateTime >= 180) {
			log.Error("retry list too long for 3min")
			time.Sleep(5 * time.Millisecond)
			continue
		}

		switch member.Action {
		case email.RetryActionReply:
			a, err := s.arc.ArchiveByAid(c, member.AID)
			if err != nil {
				log.Error("retryProc s.arc.ArchiveByAid(%d) error(%v) member(%+v)", member.AID, err, member)
				s.addRetry(c, member.AID, member.Action, member.Flag, member.FlagA)
				continue
			}
			isOpen := isOpenReplyState(a.State) > 0
			if member.Flag == archive.ReplyOn && isOpen {
				err = s.openReply(c, a, member.FlagA)
			}
			if member.Flag == archive.ReplyOff && !isOpen {
				err = s.closeReply(c, a, member.FlagA)
			}
			//if err != nil {
			//s.addRetry(c, member.AID, member.Action, member.Flag, member.FlagA)
			//}
		default:
			log.Warn("retryProc unknown action(%s) member(%+v)", member.Action, member)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
