package service

import (
	"context"

	"go-common/app/service/main/vip/model"
)

// AssociatePanel associate panel.
func (s *Service) AssociatePanel(c context.Context, arg *model.ArgAssociatePanel) (res []*model.AssociatePanelInfo, err error) {
	var vps []*model.VipPanelInfo
	res = []*model.AssociatePanelInfo{}
	if vps, err = s.VipUserPanel(c,
		arg.Mid, s.GetPlatID(c, arg.Platform, arg.PanelType, arg.MobiApp, arg.Device, arg.Build),
		arg.SortTP,
		arg.Build); err != nil {
		return
	}
	for _, v := range vps {
		res = append(res, &model.AssociatePanelInfo{
			ID:       v.Id,
			Month:    v.Month,
			PdName:   v.PdName,
			PdID:     v.PdID,
			SubType:  v.SubType,
			SuitType: v.SuitType,
			OPrice:   v.OPrice,
			DPrice:   v.DPrice,
			DRate:    v.DRate,
			Remark:   v.Remark,
			Selected: v.Selected,
		})
	}
	return
}
