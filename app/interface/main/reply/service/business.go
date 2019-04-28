package service

import (
	"context"

	"go-common/app/interface/main/reply/model/reply"
)

// ListBusiness return all non-deleted business record.
func (s *Service) ListBusiness(c context.Context) (business []*reply.Business, err error) {
	return s.dao.Business.ListBusiness(c)
}

// loadBusiness load business
func (s *Service) loadBusiness() (err error) {
	var business []*reply.Business
	if business, err = s.ListBusiness(context.Background()); err != nil {
		return
	}
	for _, b := range business {
		s.typeMapping[b.Type] = b.Alias
		s.aliasMapping[b.Alias] = b.Type
	}
	return
}
