package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"go-common/app/service/main/broadcast/model"
	identify "go-common/app/service/main/identify/api/grpc"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/zhenjl/cityhash"
)

const (
	roomsShard = 32
)

// Connect connect a conn.
func (s *Service) Connect(c context.Context, server, serverKey, cookie string, token []byte) (mid int64, key, roomID string, paltform string, accepts []int32, err error) {
	var auth model.AuthToken
	if err = json.Unmarshal(token, &auth); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", token, err)
		return
	}
	accepts = auth.Accepts
	paltform = auth.Platform
	key = serverKey
	// 兼容goim-chat
	if auth.Aid != 0 && auth.Cid != 0 {
		roomID = fmt.Sprintf("video://%d_%d", auth.Aid, auth.Cid)
		return
	}
	if auth.DeviceID != "" {
		key = auth.DeviceID
	}
	// new auth
	if auth.RoomID != "" {
		var u *url.URL
		if u, err = url.Parse(auth.RoomID); err != nil || u.Scheme == "" {
			err = ecode.RequestErr
			log.Error("url.Parse(%s) error(%v)", auth.RoomID, err)
			return
		}
		roomID = auth.RoomID
	}
	if auth.AccessKey != "" {
		info, err1 := s.identifyCli.GetTokenInfo(c, &identify.GetTokenInfoReq{Token: auth.AccessKey})
		if err1 == nil {
			mid = info.Mid
		} else {
			log.Error("s.identifyCli.GetTokenInfo(%s) failed!err:=%v", auth.AccessKey, err)
		}
	} else if cookie != "" {
		info, err1 := s.identifyCli.GetCookieInfo(c, &identify.GetCookieInfoReq{Cookie: cookie})
		if err1 == nil {
			mid = info.Mid
		} else {
			log.Error("s.identifyCli.GetCookieInfo(%s) failed!err:=%v", cookie, err)
		}
	}
	if mid > 0 {
		s.cache.Save(func() {
			if err1 := s.dao.AddMapping(context.Background(), mid, key, server); err1 != nil {
				log.Error("s.dao.AddMapping(%d,%s,%s) error(%v)", mid, key, server, err1)
			}
		})
	}
	log.Info("conn connected key:%s server:%s mid:%d token:%s", key, server, mid, token)
	return
}

// Disconnect disconnect a conn.
func (s *Service) Disconnect(c context.Context, mid int64, key, server string) (has bool, err error) {
	if has, err = s.dao.DelMapping(c, mid, key, server); err != nil {
		log.Error("s.dao.DelMapping(%d,%s,%s) error(%v)", mid, key, server, err)
		return
	}
	log.Info("conn disconnected key:%s server:%s mid:%d", key, server, mid)
	return
}

// Heartbeat heartbeat a conn.
func (s *Service) Heartbeat(c context.Context, mid int64, key, server string) (err error) {
	has, err := s.dao.ExpireMapping(c, mid, key)
	if err != nil {
		log.Error("s.dao.ExpireMapping(%d,%s,%s) error(%v)", mid, key, server, err)
		return
	}
	if !has {
		s.cache.Save(func() {
			if err1 := s.dao.AddMapping(context.Background(), mid, key, server); err1 != nil {
				log.Error("s.dao.AddMapping(%d,%s,%s) error(%v)", mid, key, server, err1)
				return
			}
		})
	}
	log.Info("conn heartbeat key:%s server:%s mid:%d", key, server, mid)
	return
}

// RenewOnline renew a server online.
func (s *Service) RenewOnline(c context.Context, server string, shard int32, roomCount map[string]int32) (mergedRoomCount map[string]int32, err error) {
	online := &model.Online{
		Server:    server,
		RoomCount: roomCount,
		Updated:   time.Now().Unix(),
	}
	mergedRoomCount = make(map[string]int32)
	for roomID, count := range s.roomCount {
		hash := cityhash.CityHash32([]byte(roomID), uint32(len(roomID))) % roomsShard
		if hash == uint32(shard) {
			mergedRoomCount[roomID] = count
		}
	}
	if err := s.dao.AddServerOnline(c, server, shard, online); err != nil {
		log.Error("s.dao.AddServerOnline(%s %d %d) merged:%d error(%v)", server, shard, len(roomCount), len(mergedRoomCount), err)
	}
	return
}

// Receive receive a message.
func (s *Service) Receive(c context.Context, mid int64, proto *model.Proto) (err error) {
	// TODO upstream message
	log.Info("conn receive a message mid:%d proto:%+v", mid, proto)
	return
}
