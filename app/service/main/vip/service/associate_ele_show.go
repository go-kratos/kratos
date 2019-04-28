package service

import (
	"context"

	"go-common/app/service/main/vip/model"
)

// EleRedPackages ele packages.
func (s *Service) EleRedPackages(c context.Context) (data []*model.EleRedPackagesResp, err error) {
	return s.dao.EleRedPackages(c)
}

// EleSpecailFoods ele specail foods.
func (s *Service) EleSpecailFoods(c context.Context) (data []*model.EleSpecailFoodsResp, err error) {
	return s.dao.EleSpecailFoods(c)
}
