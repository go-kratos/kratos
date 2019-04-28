package server

import (
	"go-common/app/service/main/vip/model"
	"go-common/library/log"
	"go-common/library/net/rpc/context"
)

// VipPanelInfo rpc user vip panel info.
func (r *RPC) VipPanelInfo(c context.Context, arg *model.ArgPanel, res *[]*model.VipPanelInfo) (err error) {
	var v []*model.VipPanelInfo
	if v, err = r.svc.VipUserPanel(c, arg.Mid, arg.Plat, arg.SortTp, arg.Build); err == nil && res != nil {
		*res = v
	}
	return
}

// VipPanelInfo2 rpc user vip panel info v2.
func (r *RPC) VipPanelInfo2(c context.Context, arg *model.ArgPanel, res *model.VipPirceResp) (err error) {
	var v *model.VipPirceResp
	if v, err = r.svc.VipUserPanelV4(c, arg); err == nil && v != nil {
		*res = *v
	}
	if err != nil {
		log.Error("rpc.VipPanelInfo2(%+v) err(%+v)", arg, err)
	}
	return
}

// VipPanelInfo5 rpc user vip panel info v5.
func (r *RPC) VipPanelInfo5(c context.Context, arg *model.ArgPanel, res *model.VipPirceResp5) (err error) {
	var v *model.VipPirceResp5
	if v, err = r.svc.VipUserPanelV5(c, arg); err == nil && v != nil {
		*res = *v
	}
	return
}
