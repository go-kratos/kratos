package service

import (
	"context"
	"time"

	"go-common/app/interface/openplatform/article/model"
)

// ActInfo .
func (s *Service) ActInfo(c context.Context, plat int8) (res *model.ActInfo, err error) {
	res = &model.ActInfo{Activities: []*model.Activity{}, Banners: []*model.Banner{}}
	if bs, _ := s.actBanners(c, plat, time.Now()); len(bs) > 0 {
		res.Banners = bs
	}
	for _, act := range s.activities {
		res.Activities = append(res.Activities, act)
	}
	return
}
