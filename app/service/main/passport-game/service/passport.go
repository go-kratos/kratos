package service

import (
	"context"
	"errors"

	"go-common/app/service/main/passport-game/model"
)

const (
	_expireSeconds               = 2592000 // 30 days
	_gameAdditionalExpireSeconds = 1296000 // 15 days
)

var (
	// ErrDispatcherError dispatcher error
	ErrDispatcherError = errors.New("dispatcher route map is error")
)

func (s *Service) tokenInfo(c context.Context, accessKey string) (token *model.Perm, err error) {
	cache := true
	if token, err = s.d.TokenCache(c, accessKey); err != nil {
		err = nil
		cache = false
	} else if token != nil {
		return
	}
	if token, err = s.d.Token(c, accessKey); err != nil {
		return
	}
	if cache && token != nil {
		s.addCache(func() {
			s.d.SetTokenCache(context.TODO(), token)
		})
	}
	return
}

// Proxy if login and getKey via model api.
func (s *Service) Proxy(c context.Context) (res bool) {
	return s.proxy
}
