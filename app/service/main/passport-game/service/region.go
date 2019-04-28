package service

import (
	"context"
	"strings"

	"go-common/app/service/main/passport-game/model"
)

const (
	_segmentation = "_"

	_origin = "origin"
)

// Regions get region list.
func (s *Service) Regions(c context.Context) (res []*model.RegionInfo) {
	return s.regionItems
}

// region returns region for token, use it in ok pattern:
// if token end with "_", returns "", false,
// if token end with "_foo", then returns "foo", true,
// otherwise returns "origin".
func region(token string) (string, bool) {
	index := strings.Index(token, _segmentation)
	if index < 0 {
		return _origin, true
	}
	suffix := token[index+1:]
	if suffix == "" {
		return "", false
	}
	return suffix, true
}
