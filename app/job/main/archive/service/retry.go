package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/archive/model/result"
	"go-common/app/job/main/archive/model/retry"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

func (s *Service) retryproc() {
	defer s.waiter.Done()
	for {
		if s.closeRetry {
			return
		}
		var (
			c   = context.TODO()
			bs  []byte
			err error
		)
		bs, err = s.PopFail(c)
		if err != nil || bs == nil {
			time.Sleep(5 * time.Second)
			continue
		}
		msg := &retry.Info{}
		if err = json.Unmarshal(bs, msg); err != nil {
			log.Error("json.Unretry dedeSyncmarshal(%s) error(%v)", bs, err)
			continue
		}
		log.Info("retry %s %s", retry.FailUpCache, bs)
		switch msg.Action {
		case retry.FailUpCache:
			s.updateResultCache(&result.Archive{AID: msg.Data.Aid, State: msg.Data.State}, nil)
		case retry.FailDatabus:
			var upInfo = &result.ArchiveUpInfo{Table: msg.Data.DatabusMsg.Table, Action: msg.Data.DatabusMsg.Table, Nw: msg.Data.DatabusMsg.Nw, Old: msg.Data.DatabusMsg.Old}
			s.sendNotify(upInfo)
		case retry.FailUpVideoCache:
			s.upVideoCache(msg.Data.Aid, msg.Data.Cids)
		case retry.FailDelVideoCache:
			s.delVideoCache(msg.Data.Aid, msg.Data.Cids)
		case retry.FailResultAdd:
			s.arcUpdate(msg.Data.Aid)
		default:
			continue
		}
	}
}

// PushFail rpush fail item to redis
func (s *Service) PushFail(c context.Context, a interface{}) (err error) {
	var (
		conn = s.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(a); err != nil {
		log.Error("json.Marshal(%v) error(%v)", a, err)
		return
	}
	if _, err = conn.Do("RPUSH", retry.FailList, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s, %s) error(%v)")
	}
	return
}

// PopFail lpop fail item from redis
func (s *Service) PopFail(c context.Context) (bs []byte, err error) {
	var conn = s.redis.Get(c)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", retry.FailList)); err != nil && err != redis.ErrNil {
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%v)", retry.FailList, err)
		return
	}
	return
}
