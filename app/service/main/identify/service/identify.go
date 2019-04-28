package service

import (
	"context"
	"net/url"
	"time"

	"go-common/app/service/main/identify/api/grpc"
	"go-common/app/service/main/identify/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	// cache -101
	_noLoginMid         = 0
	_CookieSessionField = "SESSDATA"
	_CookieBuvidField   = "buvid3"
)

var (
	_noLoginIdentify = &model.IdentifyInfo{
		Mid:     _noLoginMid,
		Expires: 86400,
	}
)

// GetCookieInfo verify user info by cookie.
func (s *Service) GetCookieInfo(c context.Context, cookie string) (res *model.IdentifyInfo, err error) {
	var cache = true
	ssda := readCookiesVal(cookie, _CookieSessionField)[0]
	if ssda == "" {
		err = ecode.NoLogin
		return
	}
	ssda, err = url.QueryUnescape(ssda)
	if err != nil {
		log.Error("cookie (%s , %s) QueryUnescape error(%v)", cookie, ssda, err)
		err = ecode.RequestErr
		return
	}

	// from cache
	if res, err = s.d.AccessCache(c, ssda); err != nil {
		cache = false
	}
	if res != nil {
		if res.Mid == _noLoginMid {
			// check if is no login mid cache
			err = ecode.NoLogin
			return
		}
		//get buvid3 from cookie
		buvid3 := readCookiesVal(cookie, _CookieBuvidField)[0]
		s.loginLog(res.Mid, metadata.String(c, metadata.RemoteIP), metadata.String(c, metadata.RemotePort), buvid3)
		return
	}
	noLoginCache := false
	// from passport
	if res, err = s.d.AccessCookie(c, cookie); err != nil {
		if err != ecode.NoLogin && err != ecode.AccessKeyErr {
			return
		}
		noLoginCache = true
	} else if res.Expires < int32(time.Now().Unix()) {
		noLoginCache = true
	}
	if noLoginCache {
		res = _noLoginIdentify
	}
	if cache && res != nil {
		s.cache.Save(func() {
			s.d.SetAccessCache(context.Background(), ssda, res)
		})
		// if cache err or res nil, don't call addLoginLog
		//get buvid3 from cookie
		buvid3 := readCookiesVal(cookie, _CookieBuvidField)[0]
		s.loginLog(res.Mid, metadata.String(c, metadata.RemoteIP), metadata.String(c, metadata.RemotePort), buvid3)
	}
	if res.Mid == _noLoginMid {
		err = ecode.NoLogin
		return
	}
	return
}

// GetTokenInfo verify user info by accesskey.
func (s *Service) GetTokenInfo(c context.Context, token *v1.GetTokenInfoReq) (res *model.IdentifyInfo, err error) {
	var cache = true
	if res, err = s.d.AccessCache(c, token.Token); err != nil {
		cache = false
	}
	if res != nil {
		if res.Mid == _noLoginMid {
			err = ecode.NoLogin
			return
		}
		s.loginLog(res.Mid, metadata.String(c, metadata.RemoteIP), metadata.String(c, metadata.RemotePort), token.Buvid)
		return
	}
	if res, err = s.d.AccessToken(c, token.Token); err != nil {
		if err != ecode.NoLogin && err != ecode.AccessKeyErr {
			return
		}
		// no login and need cache 30s
		res = _noLoginIdentify
	}
	if cache && res != nil {
		s.cache.Save(func() {
			s.d.SetAccessCache(context.Background(), token.Token, res)
		})
		// if cache err or res nil, don't call addLoginLog
		s.loginLog(res.Mid, metadata.String(c, metadata.RemoteIP), metadata.String(c, metadata.RemotePort), token.Buvid)
	}
	if res.Mid == _noLoginMid {
		err = ecode.NoLogin
		return
	}
	return
}

// DelCache delete access cache when user change pwd or logout.
func (s *Service) DelCache(c context.Context, key string) (err error) {
	return s.d.DelCache(c, key)
}

// SetCache delete access cache when user change pwd or logout.
func (s *Service) SetCache(c context.Context, key string, res *model.IdentifyInfo) (err error) {
	s.d.SetAccessCache(c, key, res)
	return
}

func (s *Service) loginLog(mid int64, ip string, ipport string, buvid string) {
	if mid <= 0 || ip == "" {
		return
	}
	s.addLoginLog(func() {
		if ok, err := s.d.ExistMIDAndIP(context.Background(), mid, ip); err != nil {
			return
		} else if ok {
			return
		}
		if s.isIntranetIP(ip) {
			s.d.SetLoginCache(context.Background(), mid, ip, s.loginCacheExpires)
			log.Warn("user ip error %s", ip)
			return
		}
		if err := s.loginLogDataBus.Send(context.Background(), ip, model.NewLoginLog(mid, ip, ipport, buvid)); err == nil {
			s.d.SetLoginCache(context.Background(), mid, ip, s.loginCacheExpires)
		}
	})
}
