package service

import (
	"context"
)

// UpdateVipStatus handle user vip staus change
func (s *Service) UpdateVipStatus(c context.Context, mid int64, vs int32) (err error) {
	s.figureDao.UpdateVipStatus(c, mid, vs)
	return
}
