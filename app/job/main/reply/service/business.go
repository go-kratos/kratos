package service

import (
	"context"

	"go-common/app/job/main/reply/model/reply"
)

// ListBusiness return all non-deleted business record.
func (s *Service) ListBusiness(c context.Context) (business []*reply.Business, err error) {
	return s.dao.Business.ListBusiness(c)
}
