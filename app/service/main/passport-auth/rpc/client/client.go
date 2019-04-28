package client

import (
	"context"

	"go-common/library/net/rpc"
)

var (
	_noRes = &struct{}{}
)

const (
	// token
	_delTokenCache = "RPC.DelTokenCache"

	// cookie
	_delCookieCache = "RPC.DelCookieCache"
)

// Service is a question service.
type Service struct {
	client *rpc.Client2
}

// New new a question service.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{
		client: rpc.NewDiscoveryCli("passport.service.auth", c),
	}
	return
}

// DelTokenCache query token.
func (s *Service) DelTokenCache(c context.Context, token string) (err error) {
	err = s.client.Call(c, _delTokenCache, token, &_noRes)
	return
}

// DelCookieCookie del cookie.
func (s *Service) DelCookieCookie(c context.Context, session string) (err error) {
	err = s.client.Call(c, _delCookieCache, session, &_noRes)
	return
}
