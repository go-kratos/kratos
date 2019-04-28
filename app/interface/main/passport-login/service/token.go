package service

import (
	"context"
	"time"

	"go-common/app/interface/main/passport-login/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_defaultTokenType     = 0
	_tokenExpireSeconds   = 2592000     // 30 days
	_refreshExpireSeconds = 2592000 * 2 // 60 days

	_expiresMonth        = 1
	_refreshExpiresMonth = 2
)

// AddToken add token
func (s *Service) AddToken(c context.Context, appID, mid int64) (res *model.RefreshTokenResp, err error) {
	var (
		accessToken  []byte
		refreshToken []byte
		now          = time.Now()
	)
	if accessToken, err = hexDecode(generateToken(mid, now, s.c.DC.Num)); err != nil {
		return
	}
	if refreshToken, err = hexDecode(generateRefresh(mid, now, s.c.DC.Num)); err != nil {
		return
	}
	token := &model.Token{
		Mid:     mid,
		AppID:   appID,
		Token:   accessToken,
		Type:    _defaultTokenType,
		Expires: now.Unix() + _tokenExpireSeconds,
	}
	refresh := &model.Refresh{
		Mid:     mid,
		AppID:   appID,
		Refresh: refreshToken,
		Token:   accessToken,
		Expires: now.Unix() + _refreshExpireSeconds,
	}
	// write to old db
	if _, err = s.dao.AddOldToken(c, token.ConvertToOld(refreshToken, 0), now); err != nil {
		return
	}
	// write to new db
	if _, err = s.dao.AddToken(c, token, now); err != nil {
		// ignore
		err = nil
	}
	if _, err = s.dao.AddRefresh(c, refresh, now); err != nil {
		// ignore
		err = nil
	}

	res = &model.RefreshTokenResp{
		Mid:     mid,
		Token:   hexEncode(token.Token),
		Refresh: hexEncode(refresh.Refresh),
		Expires: token.Expires,
	}
	s.addCache(func() {
		s.dao.SetTokenCache(context.TODO(), token.ConvertToProto())
	})
	return
}

// DelToken delete token
func (s *Service) DelToken(c context.Context, token, refresh string) (err error) {
	now := time.Now()
	// del old db
	if err = s.delOldDBToken(c, token, now); err != nil {
		return
	}
	// delete new db
	if err = s.delNewDBToken(c, token, refresh, now); err != nil {
		log.Error("delete new db token and refresh error,token is %s,refresh is %s", token, refresh)
		err = nil
	}
	s.addCache(func() {
		s.dao.DelTokenCache(context.TODO(), token)
	})
	return
}

// DelTokenByMid delete token by mid
func (s *Service) DelTokenByMid(c context.Context, mid int64) (err error) {
	now := time.Now()
	if err = s.delOldDBTokens(c, mid, now); err != nil {
		return
	}
	// del new db
	if err = s.delNewDBTokens(c, mid, now); err != nil {
		log.Error("delete  db tokens error. mid %d", mid)
		err = nil
	}
	return
}

// TokenInfo get token
func (s *Service) TokenInfo(c context.Context, token string) (res *model.Token, err error) {
	var (
		tokenBytes []byte
		cacheable  = true
		now        = time.Now()
		tokenPb    *model.TokenProto
	)
	if tokenPb, err = s.dao.TokenCache(c, token); err != nil {
		cacheable = false
	}
	if tokenPb != nil && time.Now().Unix() < tokenPb.Expires {
		// cache hit & not expired
		if tokenBytes, err = hexDecode(tokenPb.Token); err != nil {
			err = ecode.RequestErr
			return
		}
		res = &model.Token{
			Mid:     tokenPb.Mid,
			AppID:   tokenPb.AppID,
			Token:   tokenBytes,
			Type:    tokenPb.Type,
			Expires: tokenPb.Expires,
		}
		return
	}
	if tokenBytes, err = hexDecode(token); err != nil {
		err = ecode.RequestErr
		return
	}
	if s.c.Switch.SupportOld {
		if res, err = s.getOldToken(c, tokenBytes, now); err != nil {
			return
		}
	} else {
		var delta int
		delta, err = calcMonDelta(token, now)
		if err != nil {
			return
		}
		// check delta
		if delta > _expiresMonth {
			return
		}
		ct := monDiff(now, -delta)
		res, err = s.dao.GetToken(c, tokenBytes, ct)
	}
	if !cacheable || res == nil {
		log.Info("get token(%s) not found", hexEncode(tokenBytes))
		return
	}
	s.addCache(func() {
		s.dao.SetTokenCache(context.TODO(), res.ConvertToProto())
	})
	return
}

func (s *Service) getOldToken(c context.Context, tokenBytes []byte, now time.Time) (res *model.Token, err error) {
	if res, err = s.dao.GetToken(c, tokenBytes, now); err != nil || res != nil {
		return
	}
	if res, err = s.dao.GetToken(c, tokenBytes, monDiff(now, -1)); err != nil || res != nil {
		return
	}
	if res, err = s.dao.GetToken(c, tokenBytes, monDiff(now, -2)); err != nil || res != nil {
		return
	}
	return
}

// RenewGameToken renew game token
func (s *Service) RenewGameToken(c context.Context, ak string) (res *model.Token, err error) {
	var token *model.Token
	if token, err = s.TokenInfo(c, ak); err != nil || token == nil {
		res = nil
		return
	}
	now := time.Now()
	res = &model.Token{
		Mid:     token.Mid,
		AppID:   token.AppID,
		Token:   token.Token,
		Type:    _defaultTokenType,
		Expires: now.Unix() + _tokenExpireSeconds,
	}
	if err = s.DelToken(c, ak, ""); err != nil {
		return
	}
	// write to old db
	if _, err = s.dao.AddOldToken(c, token.ConvertToOld([]byte{}, 0), now); err != nil {
		return
	}
	// write to new db
	if _, err = s.dao.AddToken(c, token, now); err != nil {
		// ignore
		err = nil
	}
	s.addCache(func() {
		s.dao.SetTokenCache(context.TODO(), res.ConvertToProto())
	})
	return
}

// delOldDBTokens delete token by mid.
func (s *Service) delOldDBTokens(c context.Context, mid int64, now time.Time) (err error) {
	if _, err = s.dao.DelOldTokenByMid(c, mid, now); err != nil {
		return
	}
	if _, err = s.dao.DelOldTokenByMid(c, mid, monDiff(now, -1)); err != nil {
		return
	}
	if _, err = s.dao.DelOldTokenByMid(c, mid, monDiff(now, -2)); err != nil {
		return
	}
	return
}

// delNewDBTokens delete token by mid.
func (s *Service) delNewDBTokens(c context.Context, mid int64, now time.Time) (err error) {
	// del new db
	if _, err = s.dao.DelTokenByMid(c, mid, now); err != nil {
		err = nil
	}
	if _, err = s.dao.DelTokenByMid(c, mid, monDiff(now, -1)); err != nil {
		err = nil
	}
	if _, err = s.dao.DelTokenByMid(c, mid, monDiff(now, -2)); err != nil {
		err = nil
	}
	if _, err = s.dao.DelRefreshByMid(c, mid, now); err != nil {
		err = nil
	}
	if _, err = s.dao.DelRefreshByMid(c, mid, monDiff(now, -1)); err != nil {
		err = nil
	}
	if _, err = s.dao.DelRefreshByMid(c, mid, monDiff(now, -2)); err != nil {
		err = nil
	}
	return
}

// delOldDBToken delete token by token.
func (s *Service) delOldDBToken(c context.Context, token string, now time.Time) (err error) {
	if _, err = s.dao.DelOldToken(c, token, now); err != nil {
		return
	}
	if _, err = s.dao.DelOldToken(c, token, monDiff(now, -1)); err != nil {
		return
	}
	if _, err = s.dao.DelOldToken(c, token, monDiff(now, -2)); err != nil {
		return
	}
	return
}

// delNewDBToken delete token by token.
func (s *Service) delNewDBToken(c context.Context, token, refresh string, now time.Time) (err error) {
	var (
		tokenBytes   []byte
		refreshBytes []byte
	)
	if tokenBytes, err = hexDecode(token); err != nil {
		err = ecode.RequestErr
		return
	}
	if s.c.Switch.SupportOld {
		// del token
		if _, err = s.dao.DelToken(c, tokenBytes, now); err != nil {
			err = nil
		}
		if _, err = s.dao.DelToken(c, tokenBytes, monDiff(now, -1)); err != nil {
			err = nil
		}
	} else {
		var delta int
		delta, err = calcMonDelta(token, now)
		if err != nil {
			return
		}
		// check delta
		if delta > _expiresMonth {
			return
		}
		ct := monDiff(now, -delta)
		if _, err = s.dao.DelToken(c, tokenBytes, ct); err != nil {
			return
		}
	}
	if refresh != "" {
		if refreshBytes, err = hexDecode(refresh); err != nil {
			err = ecode.RequestErr
			return
		}
		if s.c.Switch.SupportOld {
			if _, err = s.dao.DelRefresh(c, refreshBytes, now); err != nil {
				err = nil
			}
			if _, err = s.dao.DelRefresh(c, refreshBytes, monDiff(now, -1)); err != nil {
				err = nil
			}
			if _, err = s.dao.DelRefresh(c, refreshBytes, monDiff(now, -2)); err != nil {
				err = nil
			}
		} else {
			var delta int
			delta, err = calcMonDelta(token, now)
			if err != nil {
				return
			}
			// check delta
			if delta > _refreshExpiresMonth {
				return
			}
			ct := monDiff(now, -delta)
			if _, err = s.dao.DelRefresh(c, refreshBytes, ct); err != nil {
				err = nil
			}
		}
	}
	return
}
