package service

import (
	"context"

	"go-common/app/service/main/identify-game/model"
)

// Regions get region list.
func (s *Service) Regions(c context.Context) (res []*model.RegionInfo) {
	return s.regionInfos
}
