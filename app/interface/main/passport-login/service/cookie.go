package service

import (
	"context"
	"time"

	"go-common/app/interface/main/passport-login/model"
	"go-common/library/ecode"
)

const (
	_defaultCookieType    = 0
	_defaultCookieExpires = 2592000 // 30 days
	_newSessionHexLen     = 32
)

// CookieInfo cookie info by session.
func (s *Service) CookieInfo(c context.Context, sessionData string) (cookie *model.CookieProto, err error) {
	var (
		ck          *model.Cookie
		cacheable   = true
		now         = time.Now()
		sessionByte []byte
	)
	if cookie, err = s.dao.CookieCache(c, sessionData); err != nil {
		cacheable = false
	}
	if cookie != nil && time.Now().Unix() < cookie.Expires {
		// cache hit & not expired
		return
	}
	if sessionByte, err = hexDecode(sessionData); err != nil {
		return
	}
	var delta int
	delta, err = calcMonDelta(sessionData, now)
	if err != nil {
		return
	}
	// check delta
	if delta > _expiresMonth {
		return
	}
	ct := monDiff(now, -delta)
	if ck, err = s.dao.GetCookie(c, sessionByte, ct); err != nil {
		return
	}
	if ck == nil || time.Now().Unix() >= ck.Expires {
		err = ecode.NoLogin
		return
	}
	cookie = ck.ConvertToProto()
	if !cacheable {
		return
	}
	s.addCache(func() {
		s.dao.SetCookieCache(context.TODO(), cookie)
	})
	return
}

// AddCookie add cookie
func (s *Service) AddCookie(c context.Context, mid int64) (res *model.CookieProto, err error) {
	var (
		now  = time.Now()
		sd   []byte
		csrf []byte
	)
	if sd, err = hexDecode(generateSD(mid, int(now.Month()), s.c.DC.Num)); err != nil {
		return nil, ecode.RequestErr
	}
	if csrf, err = hexDecode(generateCSRF(mid)); err != nil {
		return nil, ecode.RequestErr
	}
	cookie := &model.Cookie{
		Mid:     mid,
		Session: sd,
		CSRF:    csrf,
		Type:    _defaultCookieType,
		Expires: now.Unix() + _defaultCookieExpires,
	}

	if _, err = s.dao.AddOldCookie(c, cookie.ConvertToOld()); err != nil {
		return
	}

	if _, err = s.dao.AddCookie(c, cookie, now); err != nil {
		err = nil
	}

	s.addCache(func() {
		s.dao.SetCookieCache(context.TODO(), cookie.ConvertToProto())
	})

	res = cookie.ConvertToProto()
	return
}

// DelCookie delete cookie
func (s *Service) DelCookie(c context.Context, session string, mid int64) (err error) {
	var (
		sessionByte []byte
		now         = time.Now()
	)
	// del old db
	if _, err = s.dao.DelOldCookie(c, session, mid); err != nil {
		return
	}
	// del new db
	if sessionByte, err = decodeSession(session); err != nil {
		err = ecode.RequestErr
		return
	}
	if s.c.Switch.SupportOld {
		if _, err = s.dao.DelCookie(c, sessionByte, now); err != nil {
			err = nil
		}
		// last month
		if _, err = s.dao.DelCookie(c, sessionByte, monDiff(now, -1)); err != nil {
			err = nil
		}
	} else {
		var delta int
		delta, err = calcMonDelta(session, now)
		if err != nil {
			return
		}
		// check delta
		if delta > _expiresMonth {
			return
		}
		ct := monDiff(now, -delta)
		// del new db
		if _, err = s.dao.DelCookie(c, sessionByte, ct); err != nil {
			return
		}
	}
	s.addCache(func() {
		s.dao.DelCookieCache(context.TODO(), session)
	})
	return
}

// DelCookieByMid delete cookie by mid
func (s *Service) DelCookieByMid(c context.Context, mid int64) (err error) {
	now := time.Now()
	// del old db
	if _, err = s.dao.DelOldCookieByMid(c, mid); err != nil {
		return err
	}
	// del new db
	if _, err = s.dao.DelCookieByMid(c, mid, now); err != nil {
		err = nil
	}
	if _, err = s.dao.DelCookieByMid(c, mid, monDiff(now, -1)); err != nil {
		err = nil
	}
	if _, err = s.dao.DelCookieByMid(c, mid, monDiff(now, -2)); err != nil {
		err = nil
	}
	return

}

// AddOldCookie add old cookie
func (s *Service) AddOldCookie(c context.Context, mid int64) (res *model.CookieProto, err error) {
	var (
		now     = time.Now()
		csrf    []byte
		expires = now.Unix() + _defaultCookieExpires
	)
	if csrf, err = hexDecode(generateCSRF(mid)); err != nil {
		return nil, ecode.RequestErr
	}
	cookie := &model.Cookie{
		Mid:     mid,
		Session: []byte(s.oldSession(mid, expires, int(now.Month()))),
		CSRF:    csrf,
		Type:    _defaultCookieType,
		Expires: expires,
	}

	// add old session.
	if _, err = s.dao.AddOldCookie(c, cookie.ConvertToOld()); err != nil {
		return
	}
	// add new session.
	if _, err = s.dao.AddCookie(c, cookie, now); err != nil {
		err = nil
	}

	s.addCache(func() {
		s.dao.SetCookieCache(context.TODO(), cookie.ConvertToProto())
	})

	res = cookie.ConvertToProto()
	return
}
