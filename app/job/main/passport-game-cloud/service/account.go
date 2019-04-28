package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/passport-game-cloud/model"
	"go-common/library/log"
)

const (
	_memberTableUpdateDuration = time.Second * 1
)

func (s *Service) processMemberInfo(bmsg *model.BMsg) {
	n := new(model.Info)
	if err := json.Unmarshal(bmsg.New, n); err != nil {
		log.Error("failed to parse binlog new, json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
		return
	}
	s.memberInterval.Prom(context.TODO(), bmsg.MTS)
	switch bmsg.Action {
	case "insert":
		s.addMemberInfo(context.TODO(), n)
		s.delInfoCache(context.TODO(), n.Mid)
	case "update":
		old := new(model.Info)
		if err := json.Unmarshal(bmsg.Old, old); err != nil {
			log.Error("failed to parse binlog old, json.Unmarshal(%s) error(%v)", string(bmsg.Old), err)
			return
		}
		if n.Equals(old) {
			return
		}
		s.addMemberInfo(context.TODO(), n)
		s.delInfoCache(context.TODO(), n.Mid)
	case "delete":
		s.delMemberInfo(context.TODO(), n.Mid)
		s.delInfoCache(context.TODO(), n.Mid)
	}
}

func (s *Service) addMemberInfo(c context.Context, info *model.Info) (err error) {
	for {
		if _, err = s.d.AddMemberInfo(c, info); err == nil {
			break
		}
		time.Sleep(_memberTableUpdateDuration)
	}
	return
}

func (s *Service) delMemberInfo(c context.Context, mid int64) (err error) {
	for {
		if _, err = s.d.DelMemberInfo(c, mid); err == nil {
			break
		}
		time.Sleep(_memberTableUpdateDuration)
	}
	return
}

func (s *Service) delInfoCache(c context.Context, mid int64) (err error) {
	for i := 0; i < _accountCacheRetryCount; i++ {
		if err = s.d.DelInfoCache(c, mid); err == nil {
			break
		}
		time.Sleep(_accountCacheRetryDuration)
	}
	return
}
