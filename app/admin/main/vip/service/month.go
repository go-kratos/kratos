package service

import (
	"context"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// MonthList .
func (s *Service) MonthList(c context.Context) (res []*model.VipMonth, err error) {
	return s.dao.MonthList(c)
}

// MonthEdit .
func (s *Service) MonthEdit(c context.Context, id int64, status int8, op string) (err error) {
	var (
		m *model.VipMonth
	)
	if m, err = s.dao.GetMonth(c, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if m == nil {
		err = ecode.VipMonthErr
		return
	}
	_, err = s.dao.MonthEdit(c, id, status, op)
	return
}

// PriceList .
func (s *Service) PriceList(c context.Context, mID int64) (res []*model.VipMonthPrice, err error) {
	return s.dao.PriceList(c, mID)
}

// PriceAdd .
func (s *Service) PriceAdd(c context.Context, mp *model.VipMonthPrice) (err error) {
	_, err = s.dao.PriceAdd(c, mp)
	return
}

// PriceEdit .
func (s *Service) PriceEdit(c context.Context, mp *model.VipMonthPrice) (err error) {
	var (
		vmp *model.VipMonthPrice
	)
	if vmp, err = s.dao.GetPrice(c, mp.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if vmp == nil {
		err = ecode.VipMonthPriceErr
		return
	}
	_, err = s.dao.PriceEdit(c, mp)
	return
}
