package service

import (
	"context"
	"time"

	"go-common/app/service/main/passport-auth/model"
	"go-common/library/ecode"
)

const (
	_expiresMonth = 1

	_refreshExpiresMonth = 2
)

// RefreshInfo refreshToken by refreshToken
func (s *Service) RefreshInfo(c context.Context, rk string) (res *model.Refresh, err error) {
	return s.refreshInfo(c, rk)
}

// TokenInfo verify user info by accesskey.
func (s *Service) TokenInfo(c context.Context, ak string) (res *model.AuthReply, err error) {
	var r *model.Token
	if r, err = s.tokenInfo(c, ak); err != nil {
		return
	}
	if r == nil || time.Now().Unix() > r.Expires {
		res = _noLogin
		return
	}
	res = &model.AuthReply{
		Login:   true,
		Mid:     r.Mid,
		Expires: r.Expires,
	}
	return
}

// OldTokenInfo get old token info by accesskey.
func (s *Service) OldTokenInfo(c context.Context, ak string) (res *model.Token, err error) {
	return s.tokenInfo(c, ak)
}

// DelTokenCache del mc cache
func (s *Service) DelTokenCache(c context.Context, ak string) (err error) {
	return s.d.DelTokenCache(c, ak)
}

// DelCookieCache del mc cache
func (s *Service) DelCookieCache(c context.Context, session string) (err error) {
	return s.d.DelCookieCache(c, session)
}

func (s *Service) refreshInfo(c context.Context, rk string) (res *model.Refresh, err error) {
	// load from db
	if s.c.ServiceConf.SupportOld {
		res, err = s.refreshInfoFromDBOld(c, rk)
	} else {
		res, err = s.refreshInfoFromDBNew(c, rk)
	}
	return
}

func (s *Service) refreshInfoFromDBNew(c context.Context, rk string) (res *model.Refresh, err error) {
	now := time.Now()
	delta, err := calcMonDelta(rk, now)
	if err != nil {
		return
	}
	// check delta
	if delta > _refreshExpiresMonth {
		return
	}
	ct := monDiff(now, -delta)

	rkBytes, err := hexDecode(rk)
	if err != nil {
		return
	}

	res, err = s.d.Refresh(c, rkBytes, ct)
	return
}

func (s *Service) refreshInfoFromDBOld(c context.Context, rk string) (res *model.Refresh, err error) {
	// cur month
	now := time.Now()
	rkBytes, err := hexDecode(rk)
	if err != nil {
		return
	}
	if res, err = s.d.Refresh(c, rkBytes, now); err != nil || res != nil {
		return
	}
	// last one
	if res, err = s.d.Refresh(c, rkBytes, monDiff(now, -1)); err != nil || res != nil {
		return
	}
	// last two
	res, err = s.d.Refresh(c, rkBytes, monDiff(now, -2))
	return
}

func (s *Service) tokenInfo(c context.Context, ak string) (res *model.Token, err error) {
	cache := true
	if res, err = s.d.TokenCache(c, ak); err != nil {
		cache = false
	} else if res != nil {
		return
	}
	if s.c.ServiceConf.SupportOld {
		res, err = s.tokenInfoFromDBOld(c, ak)
	} else {
		res, err = s.tokenInfoFromDBNew(c, ak)
	}
	if err != nil || !cache {
		return
	}
	if res != nil {
		s.addCache(func() {
			s.d.SetTokenCache(context.TODO(), ak, res)
		})
		return
	}
	return
}

func (s *Service) tokenInfoFromDBNew(c context.Context, ak string) (res *model.Token, err error) {
	now := time.Now()
	delta, err := calcMonDelta(ak, now)
	if err != nil {
		return
	}
	// check delta
	if delta > _expiresMonth {
		return
	}
	ct := monDiff(now, -delta)

	tokenBytes, err := hexDecode(ak)
	if err != nil {
		err = ecode.RequestErr
		return
	}

	res, err = s.d.Token(c, tokenBytes, ct)
	return
}

func (s *Service) tokenInfoFromDBOld(c context.Context, ak string) (res *model.Token, err error) {
	// cur month
	now := time.Now()
	tokenBytes, err := hexDecode(ak)
	if err != nil {
		return
	}

	if res, err = s.d.Token(c, tokenBytes, now); err != nil || res != nil {
		return
	}
	// last one
	if res, err = s.d.Token(c, tokenBytes, monDiff(now, -1)); err != nil || res != nil {
		return
	}
	// last two
	res, err = s.d.Token(c, tokenBytes, monDiff(now, -2))
	return
}
