package service

import (
	"context"

	"go-common/app/admin/main/tag/model"
)

// ListBusiness return all business
func (s *Service) ListBusiness(c context.Context, state int32) (business []*model.Business, err error) {
	return s.dao.ListBusiness(c, state)
}

// GetBusiness return a business by type
func (s *Service) GetBusiness(c context.Context, tp int32) (business *model.Business, err error) {
	return s.dao.Business(c, tp)
}

// AddBusiness add a business
func (s *Service) AddBusiness(c context.Context, tp int32, name, appkey, remark, alias string) (id int64, err error) {
	return s.dao.InBusiness(c, tp, name, appkey, remark, alias)
}

// UpBusiness update a business's name appkey and remark
func (s *Service) UpBusiness(c context.Context, name, appkey, remark, alias string, tp int32) (id int64, err error) {
	return s.dao.UpBusiness(c, name, appkey, remark, alias, tp)
}

// UpBusinessState update a business's state
func (s *Service) UpBusinessState(c context.Context, state, tp int32) (id int64, err error) {
	return s.dao.UpBusinessState(c, state, tp)
}
