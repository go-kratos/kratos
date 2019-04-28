package service

import (
	"context"

	v1 "go-common/app/service/main/account/api"
)

// Vip get big member info
func (s *Service) Vip(c context.Context, mid int64) (vip *v1.VipInfo, err error) {
	vip, err = s.dao.Vip(c, mid)
	return
}

// Vips is
func (s *Service) Vips(ctx context.Context, mids []int64) (map[int64]*v1.VipInfo, error) {
	return s.dao.Vips(ctx, mids)
}
