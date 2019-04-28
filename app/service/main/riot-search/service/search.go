package service

import (
	"context"

	"go-common/app/service/main/riot-search/model"
)

// SearchIDOnly return ID Only
func (s *Service) SearchIDOnly(c context.Context, arg *model.RiotSearchReq) (res *model.IDsResp) {
	res = s.dao.SearchIDOnly(arg)
	return
}

// Search return both id and content
func (s *Service) Search(c context.Context, arg *model.RiotSearchReq) (res *model.DocumentsResp) {
	res = s.dao.Search(arg)
	return
}

// Has return DocId exist
func (s *Service) Has(c context.Context, id uint64) bool {
	return s.dao.Has(id)
}
