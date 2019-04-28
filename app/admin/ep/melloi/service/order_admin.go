package service

import (
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
)

// QueryOrderAdmin get administrator for order by current username
func (s *Service) QueryOrderAdmin(userName string) (*model.OrderAdmin, error) {
	return s.dao.QueryOrderAdmin(userName)
}

// AddOrderAdmin add administrator for order
func (s *Service) AddOrderAdmin(admin *model.OrderAdmin) (err error) {
	var oa *model.OrderAdmin
	oa, _ = s.dao.QueryOrderAdmin(admin.UserName)
	if oa.UserName == admin.UserName {
		return ecode.MelloiAdminExist
	}
	return s.dao.AddOrderAdmin(admin)
}
