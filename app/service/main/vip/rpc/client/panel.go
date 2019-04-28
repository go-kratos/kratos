package client

import (
	"context"

	"go-common/app/service/main/vip/model"
)

const (
	_VipPanelInfo  = "RPC.VipPanelInfo"
	_vipPanelInfo2 = "RPC.VipPanelInfo2"
	_vipPanelInfo5 = "RPC.VipPanelInfo5"
)

// VipPanelInfo rpc user vip panel info.
func (s *Service) VipPanelInfo(c context.Context, arg *model.ArgPanel) (res []*model.VipPanelInfo, err error) {
	err = s.client.Call(c, _VipPanelInfo, arg, &res)
	return
}

// VipPanelInfo2 vip panel v2.
func (s *Service) VipPanelInfo2(c context.Context, arg *model.ArgPanel) (res *model.VipPirceResp, err error) {
	res = new(model.VipPirceResp)
	err = s.client.Call(c, _vipPanelInfo2, arg, &res)
	return
}

// VipPanelInfo5 vip panel v5.
func (s *Service) VipPanelInfo5(c context.Context, arg *model.ArgPanel) (res *model.VipPirceResp5, err error) {
	res = new(model.VipPirceResp5)
	err = s.client.Call(c, _vipPanelInfo5, arg, &res)
	return
}
