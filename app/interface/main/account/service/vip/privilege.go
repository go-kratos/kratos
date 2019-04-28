package vip

import (
	"context"

	"go-common/app/service/main/vip/model"
)

// PrivilegebySid privilege by sid.
func (s *Service) PrivilegebySid(c context.Context, arg *model.ArgPrivilegeBySid) (res *model.PrivilegesResp, err error) {
	res, err = s.vipRPC.PrivilegeBySid(c, arg)
	return
}

// PrivilegebyType privilege by type.
func (s *Service) PrivilegebyType(c context.Context, arg *model.ArgPrivilegeDetail) (res []*model.PrivilegeDetailResp, err error) {
	res, err = s.vipRPC.PrivilegeByType(c, arg)
	return
}
