package service

import (
	"context"
	"time"

	"go-common/app/interface/main/growup/model"
)

// GetBanner get banner for now
func (s *Service) GetBanner(c context.Context) (b *model.Banner, err error) {
	return s.dao.Banner(c, time.Now().Unix())
}
