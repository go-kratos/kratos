package service

import (
	"context"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_maxPushBatch = 200
)

// PushKeys push a message by keys.
func (s *Service) PushKeys(c context.Context, op int32, keys []string, msg string, contentType int32) (err error) {
	if len(keys) > _maxPushBatch {
		err = ecode.LimitExceed
		return
	}
	servers, err := s.dao.ServersByKeys(c, keys)
	if err != nil {
		return
	}
	pushKeys := make(map[string][]string)
	for i, key := range keys {
		server := servers[i]
		if server != "" && key != "" {
			pushKeys[server] = append(pushKeys[server], key)
		}
	}
	for server := range pushKeys {
		if err = s.dao.PushMsg(c, op, server, msg, pushKeys[server], contentType); err != nil {
			return
		}
	}
	return
}

// PushMids push a message by mid.
func (s *Service) PushMids(c context.Context, op int32, mids []int64, msg string, contentType int32) (err error) {
	if len(mids) > _maxPushBatch {
		err = ecode.LimitExceed
		return
	}
	now := time.Now().Unix()
	s.stats.Info("push_mids", op, xstr.JoinInts(mids), "", "", len(mids), now)
	keyServers, olMids, err := s.dao.KeysByMids(c, mids)
	if err != nil {
		return
	}
	keys := make(map[string][]string)
	for key, server := range keyServers {
		if key != "" && server != "" {
			keys[server] = append(keys[server], key)
		} else {
			log.Warn("push key:%s server:%s is empty", key, server)
		}
	}
	for server, keys := range keys {
		if err = s.dao.PushMsg(c, op, server, msg, keys, contentType); err != nil {
			return
		}
	}
	s.stats.Info("push_mids_ol", op, xstr.JoinInts(olMids), "", "", len(olMids), now)
	return
}

// PushRoom push a message by room.
func (s *Service) PushRoom(c context.Context, op int32, room, msg string, contentType int32) (err error) {
	if err = s.dao.BroadcastRoomMsg(c, op, room, msg, contentType); err != nil {
		return
	}
	s.stats.Info("push_room", op, "", room, "", s.roomCount[room], time.Now().Unix())
	return
}

// PushAll push a message to all.
func (s *Service) PushAll(c context.Context, op, speed int32, msg, platform string, contentType int32) (err error) {
	if err = s.dao.BroadcastMsg(c, op, speed, msg, platform, contentType); err != nil {
		return
	}
	s.stats.Info("push_all", op, "", "", platform, 0, time.Now().Unix())
	return
}
