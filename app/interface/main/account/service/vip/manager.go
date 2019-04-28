package vip

import (
	"context"

	"go-common/app/interface/main/account/model"
)

// ManagerInfo manager info.
func (s *Service) ManagerInfo(c context.Context) (res *model.ManagerResp, err error) {
	res = new(model.ManagerResp)
	res.JointlyInfo, err = s.vipRPC.Jointly(c)
	return
}
