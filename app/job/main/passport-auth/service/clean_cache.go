package service

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/passport-auth/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

type authTokenBMsg struct {
	Action string
	Table  string
	New    *model.AuthToken
}

type authCookieBMsg struct {
	Action string
	Table  string
	New    *model.AuthCookie
}

func (s *Service) authConsumeProc() {
	// fill callbacks
	s.authGroup.New = func(msg *databus.Message) (res interface{}, err error) {
		bmsg := new(model.BMsg)
		if err = json.Unmarshal(msg.Value, &bmsg); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			return
		}
		log.Info("receive auth msg action(%s) table(%s) key(%s) partition(%d) offset(%d) timestamp(%d) New(%s) Old(%s)",
			bmsg.Action, bmsg.Table, msg.Key, msg.Partition, msg.Offset, msg.Timestamp, string(bmsg.New), string(bmsg.Old))
		if strings.HasPrefix(bmsg.Table, "user_token_") {
			tokenBMsg := &authTokenBMsg{
				Action: bmsg.Action,
				Table:  bmsg.Table,
			}
			newToken := new(model.AuthToken)
			if err = json.Unmarshal(bmsg.New, &newToken); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", bmsg.New, err)
				return
			}
			tokenBMsg.New = newToken
			return tokenBMsg, nil
		} else if strings.HasPrefix(bmsg.Table, "user_cookie_") {
			cookieBMsg := &authCookieBMsg{
				Action: bmsg.Action,
				Table:  bmsg.Table,
			}
			newCookie := new(model.AuthCookie)
			if err = json.Unmarshal(bmsg.New, newCookie); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", bmsg.New, err)
				return
			}
			cookieBMsg.New = newCookie
			return cookieBMsg, nil
		}
		return
	}
	s.authGroup.Split = func(msg *databus.Message, data interface{}) int {
		if t, ok := data.(*authTokenBMsg); ok {
			return int(t.New.Mid)
		} else if t, ok := data.(*authCookieBMsg); ok {
			return int(t.New.Mid)
		}
		return 0
	}
	s.authGroup.Do = func(msgs []interface{}) {
		for _, m := range msgs {
			if msg, ok := m.(*authTokenBMsg); ok {
				if msg.Action != "delete" {
					return
				}
				for {
					if err := s.cleanTokenCache(msg.New.Token, msg.New.Mid); err != nil {
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}
			} else if msg, ok := m.(*authCookieBMsg); ok {
				if msg.Action != "delete" {
					return
				}
				for {
					if err := s.cleanCookieCache(msg.New.Session, msg.New.Mid); err != nil {
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}
			}
		}
	}
	// start the group
	s.authGroup.Start()
}

func (s *Service) cleanTokenCache(tokenBase64 string, mid int64) (err error) {
	var bytes []byte
	if bytes, err = base64.StdEncoding.DecodeString(tokenBase64); err != nil {
		log.Error("cleanTokenCache base64 decode err %v", err)
		err = nil
		return
	}
	token := hex.EncodeToString(bytes)
	if err = s.authRPC.DelTokenCache(context.Background(), token); err != nil {
		log.Error("cleanTokenCache err, %v", err)
		return
	}
	if err = s.dao.AsoCleanCache(context.Background(), token, "", mid); err != nil {
		return
	}
	return
}

func (s *Service) cleanCookieCache(cookieBase64 string, mid int64) (err error) {
	var bytes []byte
	if bytes, err = base64.StdEncoding.DecodeString(cookieBase64); err != nil {
		log.Error("cleanCookieCache base64 decode err %v", err)
		err = nil
		return
	}
	session := string(bytes)
	if err = s.authRPC.DelCookieCookie(context.Background(), session); err != nil {
		log.Error("cleanCookieCache err, %v", err)
		return
	}
	if err = s.dao.AsoCleanCache(context.Background(), "", session, mid); err != nil {
		return
	}
	return
}
