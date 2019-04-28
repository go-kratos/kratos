package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/identify/model"
	mdl "go-common/app/service/main/identify/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

type actionToken struct {
	model.Token
	Action string
}

type actionCookie struct {
	model.Cookie
	Action string
}

func (s *Service) identifyNew(msg *databus.Message) (interface{}, error) {
	bmsg := new(model.BMsg)
	if err := json.Unmarshal(msg.Value, bmsg); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
		return nil, err
	}
	log.Info("identify service process databus message, table(%s)", bmsg.Table)
	if strings.HasPrefix(bmsg.Table, _tokenTable) {
		t := new(actionToken)
		if err := json.Unmarshal(bmsg.New, t); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
			return nil, err
		}
		log.Info("process databus message, table(%s), action(%s)", bmsg.Table, bmsg.Action)
		t.Action = bmsg.Action
		return t, nil
	}
	if strings.HasPrefix(bmsg.Table, _cookieTable) {
		t := new(actionCookie)
		if err := json.Unmarshal(bmsg.New, t); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
			return nil, err
		}
		log.Info("process databus message, table(%s), action(%s)", bmsg.Table, bmsg.Action)
		t.Action = bmsg.Action
		return t, nil
	}
	return bmsg, nil
}

func (s *Service) identifySplit(msg *databus.Message, data interface{}) int {
	mergeNum := int64(s.c.Databusutil.Num)
	mid := int64(0)
	switch t := data.(type) {
	case *model.Cookie:
		mid = t.Mid
	case *model.Token:
		mid = t.Mid
	}
	return int(mid % mergeNum)
}

func (s *Service) processIdentifyInfo(bmsgs []interface{}) {
	for _, msg := range bmsgs {
		switch t := msg.(type) {
		case *actionToken:
			info := &mdl.IdentifyInfo{
				Mid:     t.Mid,
				Expires: t.Expires,
			}
			if ok := isGameAppID(t.APPID); ok {
				continue
			}
			s.processIdentify(t.Action, t.AccessToken, info)
		case *actionCookie:
			info := &mdl.IdentifyInfo{
				Mid:     t.Mid,
				Csrf:    t.CSRFToken,
				Expires: t.ExpireTime,
			}
			s.processIdentify(t.Action, t.SessionData, info)
		}
	}
}

func (s *Service) processIdentify(action, key string, info *mdl.IdentifyInfo) (err error) {
	if err = s.processMc(action, key, info); err != nil {
		return
	}
	log.Info("identify process action(%s) mid(%d) csrf(%s) key(%s) success", action, info.Mid, info.Csrf, key[:8])
	return
}

// processMc .
func (s *Service) processMc(action, k string, info *mdl.IdentifyInfo) (err error) {
	if len(s.poolm) == 0 {
		return
	}
	if _insertAction == action {
		expire := info.Expires
		if expire < int32(time.Now().Unix()) {
			log.Error("identify expire error(%d,%d)", info.Expires, time.Now().Unix())
			return
		}
		for name, p := range s.poolm {
			mcc, ok := s.c.Memcaches[name]
			if !ok || mcc == nil {
				return
			}
			key := mcc.Prefix + k
			conn := p.Get(context.Background())
			err = conn.Set(&memcache.Item{Key: key, Object: info, Flags: memcache.FlagProtobuf, Expiration: expire})
			conn.Close()
			if err != nil {
				log.Error("old identify set error(%s,%d,%v)", key, info.Expires, err)
				err = nil
			}
		}
		return
	}
	if _delteAction == action {
		for name, p := range s.poolm {
			mcc, ok := s.c.Memcaches[name]
			if !ok || mcc == nil {
				return
			}
			key := mcc.Prefix + k
			// if delete failed, retry until success
			for {
				conn := p.Get(context.Background())
				err := conn.Delete(key)
				if err == nil || err == memcache.ErrNotFound {
					conn.Close()
					break
				}
				log.Error("dao.DelCache(%s) error(%v)", key, err)
				conn.Close()
				time.Sleep(time.Second)
			}
		}
	}
	return
}

func isGameAppID(appid int64) (res bool) {
	for i := 0; i < len(_gameAppID); i++ {
		if _gameAppID[i] == appid {
			return true
		}
	}
	return false
}
