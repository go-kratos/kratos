package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/passport-auth/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

var (
	local, _ = time.LoadLocation("Local")
)

type tokenBMsg struct {
	Action string
	Table  string
	New    *model.OldToken
}

type cookieBMsg struct {
	Action string
	Table  string
	New    *model.OldCookie
}

func (s *Service) consumeproc() {
	// fill callbacks
	s.g.New = func(msg *databus.Message) (res interface{}, err error) {
		bmsg := new(model.BMsg)
		if err = json.Unmarshal(msg.Value, &bmsg); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			return
		}
		log.Info("receive aso msg action(%s) table(%s) key(%s) partition(%d) offset(%d) timestamp(%d) New(%s) Old(%s)",
			bmsg.Action, bmsg.Table, msg.Key, msg.Partition, msg.Offset, msg.Timestamp, string(bmsg.New), string(bmsg.Old))
		if strings.HasPrefix(bmsg.Table, "aso_app_perm") {
			tokenBMsg := &tokenBMsg{
				Action: bmsg.Action,
				Table:  bmsg.Table,
			}
			newToken := new(model.OldToken)
			if err = json.Unmarshal(bmsg.New, &newToken); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", bmsg.New, err)
				return
			}
			tokenBMsg.New = newToken
			return tokenBMsg, nil
		} else if strings.HasPrefix(bmsg.Table, "aso_cookie_token") {
			cookieBMsg := &cookieBMsg{
				Action: bmsg.Action,
				Table:  bmsg.Table,
			}
			newCookie := new(model.OldCookie)
			if err = json.Unmarshal(bmsg.New, newCookie); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", bmsg.New, err)
				return
			}
			cookieBMsg.New = newCookie
			return cookieBMsg, nil
		}
		return
	}
	s.g.Split = func(msg *databus.Message, data interface{}) int {
		if t, ok := data.(*tokenBMsg); ok {
			return int(t.New.Mid)
		} else if t, ok := data.(*cookieBMsg); ok {
			return int(t.New.Mid)
		}
		return 0
	}
	s.g.Do = func(msgs []interface{}) {
		for _, m := range msgs {
			if msg, ok := m.(*tokenBMsg); ok {
				for {
					if err := s.handleToken(msg); err != nil {
						log.Error("do handleToken error", err)
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}
			} else if msg, ok := m.(*cookieBMsg); ok {
				for {
					if err := s.handleCookie(msg); err != nil {
						log.Error("do handleCookie error", err)
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}
			}
		}
	}
	// start the group
	s.g.Start()
}

func (s *Service) handleCookie(cookie *cookieBMsg) (err error) {
	newCookie := &model.Cookie{
		Mid:     cookie.New.Mid,
		Session: cookie.New.Session,
		CSRF:    cookie.New.CSRFToken,
		Expires: cookie.New.Expires,
	}
	csrfByte, _ := hex.DecodeString(cookie.New.CSRFToken)
	if strings.ToLower(cookie.Action) == "insert" {
		if _, err = s.dao.AddCookie(context.Background(), newCookie, []byte(newCookie.Session), csrfByte, time.Now()); err != nil {
			return
		}
		return
	} else if strings.ToLower(cookie.Action) == "delete" {
		now := time.Now()
		if _, err = s.dao.DelCookie(context.Background(), []byte(cookie.New.Session), now); err != nil {
			log.Error("del db cookie(%s) error(%+v)", cookie.New.Session, err)
			return
		}
		if _, err = s.dao.DelCookie(context.Background(), []byte(cookie.New.Session), previousMonth(now, -1)); err != nil {
			log.Error("del db cookie(%s) error(%+v)", cookie.New.Session, err)
			return
		}
		if newCookie.Expires > now.Unix() {
			var t time.Time
			if cookie.New.ModifyTime == "0000-00-00 00:00:00" {
				t = time.Now()
			} else {
				t, err = time.ParseInLocation("2006-01-02 15:04:05", cookie.New.ModifyTime, local)
				if err != nil {
					log.Error("handleCookie error: ctime parse error(%+v) cookie(%+v)", err, cookie.New)
					return
				}
			}
			timestamp := xtime.Time(t.Unix())
			newCookie.Ctime = timestamp
			if _, err = s.dao.AddCookieDeleted(context.Background(), newCookie, []byte(newCookie.Session), csrfByte, t); err != nil {
				return
			}
		}
		if err = s.authRPC.DelCookieCookie(context.Background(), cookie.New.Session); err != nil {
			log.Error("del cache cookie(%s) error(%+v)", cookie.New.Session, err)
		}
		if err = s.dao.AsoCleanCache(context.Background(), "", cookie.New.Session, cookie.New.Mid); err != nil {
			return
		}
	}
	return
}

func (s *Service) handleToken(token *tokenBMsg) (err error) {
	var t time.Time
	t, err = time.ParseInLocation("2006-01-02 15:04:05", token.New.CTime, local)
	if err != nil {
		log.Error("handleToken error: ctime parse err. (%+v)", token.New)
		return
	}
	timestamp := xtime.Time(t.Unix())
	newToken := &model.Token{
		Mid:     token.New.Mid,
		Token:   token.New.AccessToken,
		AppID:   token.New.AppID,
		Expires: token.New.Expires,
		Type:    token.New.Type,
		Ctime:   timestamp,
	}
	tokenByte, _ := hex.DecodeString(newToken.Token)
	if strings.ToLower(token.Action) == "insert" {
		if _, err = s.dao.AddToken(context.Background(), newToken, tokenByte, t); err != nil {
			return
		}
		if token.New.RefreshToken == "" {
			return
		}
		newRefresh := &model.Refresh{
			Mid:     token.New.Mid,
			Refresh: token.New.RefreshToken,
			Token:   token.New.AccessToken,
			AppID:   token.New.AppID,
			Expires: token.New.Expires, // 需要考虑+30天
		}
		tokenByteRefresh, _ := hex.DecodeString(newRefresh.Token)
		refreshByte, _ := hex.DecodeString(newRefresh.Refresh)
		if _, err = s.dao.AddRefresh(context.Background(), newRefresh, refreshByte, tokenByteRefresh, t); err != nil {
			return
		}
	} else if strings.ToLower(token.Action) == "delete" {
		now := time.Now()
		if _, err = s.dao.DelToken(context.Background(), tokenByte, now); err != nil {
			log.Error("del db token(%s) error(%+v)", token.New.AccessToken, err)
			return
		}
		if _, err = s.dao.DelToken(context.Background(), tokenByte, previousMonth(now, -1)); err != nil {
			log.Error("del db token(%s) error(%+v)", token.New.AccessToken, err)
			return
		}
		if _, err = s.dao.DelToken(context.Background(), tokenByte, previousMonth(now, -2)); err != nil {
			log.Error("del db token(%s) error(%+v)", token.New.AccessToken, err)
			return
		}
		if _, err = s.dao.AddTokenDeleted(context.Background(), newToken, tokenByte, t); err != nil {
			return
		}
		if err = s.authRPC.DelTokenCache(context.Background(), token.New.AccessToken); err != nil {
			log.Error("del cache token(%s) error(%+v)", token.New.AccessToken, err)
		}
		if err = s.dao.AsoCleanCache(context.Background(), token.New.AccessToken, "", token.New.Mid); err != nil {
			return
		}
		if token.New.RefreshToken == "" {
			return
		}
		refreshByte, _ := hex.DecodeString(token.New.RefreshToken)
		if _, err = s.dao.DelRefresh(context.Background(), refreshByte, t); err != nil {
			log.Error("del db refresh(%s) error(%+v)", token.New.RefreshToken, err)
			return
		}
		if _, err = s.dao.DelRefresh(context.Background(), refreshByte, previousMonth(t, -2)); err != nil {
			log.Error("del db refresh(%s) error(%+v)", token.New.RefreshToken, err)
			return
		}
	}
	return
}

func previousMonth(t time.Time, delta int) time.Time {
	if delta == 0 {
		return t
	}
	year, month, _ := t.Date()
	thisMonthFirstDay := time.Date(year, month, 1, 1, 1, 1, 1, t.Location())
	return thisMonthFirstDay.AddDate(0, delta, 0)
}
