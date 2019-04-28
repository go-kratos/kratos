package service

import (
	"context"

	"go-common/app/service/main/passport-game/model"
)

// TokenPBCache get token pb cache.
func (s *Service) TokenPBCache(c context.Context, key string) (res *model.Perm, err error) {
	return s.d.TokenPBCache(c, key)
}

// InfoPBCache get info pb cache.
func (s *Service) InfoPBCache(c context.Context, key string) (res *model.Info, err error) {
	return s.d.InfoPBCache(c, key)
}
