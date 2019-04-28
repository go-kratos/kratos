package service

import (
	"context"
	"time"

	"go-common/app/service/main/passport-auth/model"
)

const (
	_newTokenHexLen     = 32
	_newTokenBinByteLen = 16
)

// CookieInfo verify user info by cookie.
func (s *Service) CookieInfo(c context.Context, sd string) (res *model.AuthReply, err error) {
	var r *model.Cookie
	if r, err = s.cookieInfo(c, sd); err != nil {
		return
	}
	if r == nil || time.Now().Unix() > r.Expires {
		res = _noLogin
		return
	}
	res = &model.AuthReply{
		Login:   true,
		Mid:     r.Mid,
		CSRF:    r.CSRF,
		Expires: r.Expires,
	}
	return
}

func (s *Service) cookieInfo(c context.Context, sd string) (res *model.Cookie, err error) {
	cache := true
	if res, err = s.d.CookieCache(c, sd); err != nil {
		cache = false
	} else if res != nil {
		return
	}
	if s.c.ServiceConf.SupportOld {
		res, err = s.cookieInfoFromDBOld(c, sd)
	} else {
		res, err = s.cookieInfoFromDBNew(c, sd)
	}
	if err != nil || !cache {
		return
	}
	if res != nil {
		s.addCache(func() {
			s.d.SetCookieCache(context.TODO(), sd, res)
		})
		return
	}
	return
}

func (s *Service) cookieInfoFromDBNew(c context.Context, sd string) (res *model.Cookie, err error) {
	now := time.Now()
	delta, err := calcMonDelta(sd, now)
	if err != nil {
		return
	}
	// check delta
	if delta > _expiresMonth {
		return
	}
	ct := monDiff(now, -delta)

	sdBytes, err := decodeSD(sd)
	if err != nil {
		return
	}
	res, sdb, err := s.d.Cookie(c, sdBytes, ct)
	if err != nil {
		return
	}
	res.Session = encodeSD(sdb)
	return
}

func (s *Service) cookieInfoFromDBOld(c context.Context, sd string) (res *model.Cookie, err error) {
	// cur month
	sdBytes, err := decodeSD(sd)
	if err != nil {
		return
	}
	var sdb []byte
	now := time.Now()
	if res, sdb, err = s.d.Cookie(c, sdBytes, now); err != nil || res != nil {
		res.Session = encodeSD(sdb)
		return
	}
	// last month
	if res, sdb, err = s.d.Cookie(c, sdBytes, monDiff(now, -1)); err != nil || res != nil {
		res.Session = encodeSD(sdb)
		return
	}
	return
}

func decodeSD(sd string) (res []byte, err error) {
	// if is new sd
	if len(sd) == _newTokenHexLen {
		return hexDecode(sd)
	}
	// else if is old sd
	return []byte(sd), nil
}

func encodeSD(b []byte) (s string) {
	// format new
	if len(b) == _newTokenBinByteLen {
		return hexEncode(b)
	}
	// or format old
	return string(b)
}
