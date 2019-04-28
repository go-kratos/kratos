package service

import (
	"context"

	"go-common/app/interface/main/passport-login/model"
)

// ProxyAddCookie proxy add cookie
func (s *Service) ProxyAddCookie(c context.Context, mid int64) (res *model.CookieProto, err error) {
	return s.AddOldCookie(c, mid)
}

// ProxyDeleteCookie proxy delete cookie
func (s *Service) ProxyDeleteCookie(c context.Context, mid int64, session string) (err error) {
	return s.DelCookie(c, session, mid)
}

// ProxyDeleteCookies proxy delete cookies
func (s *Service) ProxyDeleteCookies(c context.Context, mid int64) (err error) {
	return s.DelCookieByMid(c, mid)
}
