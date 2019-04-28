package service

import (
	"context"
	"go-common/app/interface/main/passport-login/model"
)

// ProxyAddToken proxy add token
func (s *Service) ProxyAddToken(c context.Context, appID, mid int64) (res *model.RefreshTokenResp, err error) {
	return s.AddToken(c, appID, mid)
}

// ProxyDeleteToken proxy delete cookie
func (s *Service) ProxyDeleteToken(c context.Context, token string) (err error) {
	return s.DelToken(c, token, "")
}

// ProxyDeleteTokens proxy delete tokens
func (s *Service) ProxyDeleteTokens(c context.Context, mid int64) (err error) {
	return s.DelTokenByMid(c, mid)
}

// ProxyDeleteGameTokens proxy delete tokens
func (s *Service) ProxyDeleteGameTokens(c context.Context, mid, appID int64) (err error) {
	return s.DelTokenByMid(c, mid)
}

// ProxyRenewToken renew game token
func (s *Service) ProxyRenewToken(c context.Context, ak string) (res *model.Token, err error) {
	return s.RenewGameToken(c, ak)
}
