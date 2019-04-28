package service

import (
	"context"

	"go-common/app/interface/main/esports/model"
)

// Search  search video list.
func (s *Service) Search(c context.Context, mid int64, p *model.ParamSearch, buvid string) (rs *model.SearchEsp, err error) {
	rs, err = s.dao.Search(c, mid, p, buvid)
	return
}
