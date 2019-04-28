package service

import (
	"context"

	"go-common/app/service/openplatform/ticket-item/model"
)

// VersionSearch 项目版本查询
func (s *ItemService) VersionSearch(c context.Context, in *model.VersionSearchParam) (versions *model.VersionSearchList, err error) {
	versions, err = s.dao.VersionSearch(c, in)
	return
}
