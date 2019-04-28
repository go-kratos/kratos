package service

import (
	"context"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
)

// VipPriceConfigs get vip price configs.
func (s *Service) VipPriceConfigs(c context.Context, arg *model.ArgVipPrice) (res []*model.VipPriceConfig, err error) {
	var (
		vpcs []*model.VipPriceConfig
		mvd  map[int64]*model.VipDPriceConfig
	)
	if vpcs, err = s.dao.VipPriceConfigs(c); err != nil {
		return
	}
	if mvd, err = s.dao.VipPriceDiscountConfigs(c); err != nil {
		return
	}
	for _, item := range vpcs {
		if arg.Plat != -1 && item.Plat != arg.Plat {
			continue
		}
		if arg.Month != -1 && item.Month != arg.Month {
			continue
		}
		if arg.SubType != -1 && item.SubType != arg.SubType {
			continue
		}
		if arg.SuitType != -1 && item.SuitType != arg.SuitType {
			continue
		}
		if vdc, ok := mvd[item.ID]; ok {
			item.NPrice = vdc.DPrice
			item.PdID = vdc.PdID
		} else {
			item.NPrice = item.OPrice
		}
		res = append(res, item)
	}
	return
}

// VipPriceConfigID get vip price config by id.
func (s *Service) VipPriceConfigID(c context.Context, arg *model.ArgVipPriceID) (res *model.VipPriceConfig, err error) {
	res, err = s.dao.VipPriceConfigID(c, arg)
	return
}

// AddVipPriceConfig add vip price config.
func (s *Service) AddVipPriceConfig(c context.Context, arg *model.ArgAddOrUpVipPrice) (err error) {
	var count int64
	if len(arg.Superscript) > 18 {
		err = ecode.VipPanelSuperscriptTooLongErr
		return
	}
	// 平台校验
	platForm, err := s.dao.PlatformByID(c, int64(arg.Plat))
	if err != nil || platForm == nil {
		err = ecode.VipPanelPlatNotExitErr
		return
	}
	// 价格校验
	if arg.OPrice <= 0 {
		err = ecode.VipPanelValidOPriceErr
		return
	}
	// 面版冲突
	if count, err = s.dao.VipPriceConfigUQCheck(c, arg); err != nil {
		return
	}
	if count > 0 {
		err = ecode.VipPanelConfNotUQErr
		return
	}
	vpc := &model.VipPriceConfig{
		Plat:        arg.Plat,
		PdName:      arg.PdName,
		PdID:        arg.PdID,
		SuitType:    arg.SuitType,
		Month:       arg.Month,
		SubType:     arg.SubType,
		OPrice:      arg.OPrice,
		Remark:      arg.Remark,
		Status:      model.VipPriceConfigStatusON,
		Operator:    arg.Operator,
		OpID:        arg.OpID,
		Superscript: arg.Superscript,
		Selected:    arg.Selected,
		StartBuild:  arg.StartBuild,
		EndBuild:    arg.EndBuild,
	}
	err = s.dao.AddVipPriceConfig(c, vpc)
	return
}

// UpVipPriceConfig update vip price config.
func (s *Service) UpVipPriceConfig(c context.Context, arg *model.ArgAddOrUpVipPrice) (err error) {
	var (
		count int64
		max   float64
		ovpc  *model.VipPriceConfig
	)
	if len(arg.Superscript) > 18 {
		err = ecode.VipPanelSuperscriptTooLongErr
		return
	}
	if ovpc, err = s.dao.VipPriceConfigID(c, &model.ArgVipPriceID{ID: arg.ID}); err != nil {
		return
	}
	if max, err = s.dao.VipMaxPriceDiscount(c, arg); err != nil {
		return
	}
	if ovpc == nil {
		err = ecode.VipPanelConfNotExistErr
		return
	}
	// 价格校验
	if arg.OPrice <= 0 {
		err = ecode.VipPanelValidOPriceErr
		return
	}
	if arg.OPrice < max {
		err = ecode.VipPanelValidDPriceErr
		return
	}
	if arg.Plat != ovpc.Plat || arg.Month != ovpc.Month || arg.SubType != ovpc.SubType || arg.SuitType != ovpc.SuitType ||
		arg.StartBuild != ovpc.StartBuild || arg.EndBuild != ovpc.EndBuild {
		// 面版冲突
		if count, err = s.dao.VipPriceConfigUQCheck(c, arg); err != nil {
			return
		}
		if count > 0 {
			err = ecode.VipPanelConfNotUQErr
			return
		}
	}
	vpc := &model.VipPriceConfig{
		ID:          arg.ID,
		Plat:        arg.Plat,
		PdName:      arg.PdName,
		PdID:        arg.PdID,
		SuitType:    arg.SuitType,
		Month:       arg.Month,
		SubType:     arg.SubType,
		OPrice:      arg.OPrice,
		Remark:      arg.Remark,
		Status:      model.VipPriceConfigStatusON,
		Operator:    arg.Operator,
		OpID:        arg.OpID,
		Superscript: arg.Superscript,
		Selected:    arg.Selected,
		StartBuild:  arg.StartBuild,
		EndBuild:    arg.EndBuild,
	}
	err = s.dao.UpVipPriceConfig(c, vpc)
	return
}

// DelVipPriceConfig delete vip price config.
func (s *Service) DelVipPriceConfig(c context.Context, arg *model.ArgVipPriceID) (err error) {
	err = s.dao.DelVipPriceConfig(c, arg)
	return
}

// VipDPriceConfigs get discount price config list.
func (s *Service) VipDPriceConfigs(c context.Context, arg *model.ArgVipPriceID) (res []*model.VipDPriceConfig, err error) {
	res, err = s.dao.VipDPriceConfigs(c, arg)
	return
}

// VipDPriceConfigID get discount price config by id.
func (s *Service) VipDPriceConfigID(c context.Context, arg *model.ArgVipDPriceID) (res *model.VipDPriceConfig, err error) {
	res, err = s.dao.VipDPriceConfigID(c, arg)
	return
}

// AddVipDPriceConfig add vip discount price config.
func (s *Service) AddVipDPriceConfig(c context.Context, arg *model.ArgAddOrUpVipDPrice) (err error) {
	var (
		vpc *model.VipPriceConfig
		mvd map[int64]*model.VipDPriceConfig
	)
	if vpc, err = s.dao.VipPriceConfigID(c, &model.ArgVipPriceID{ID: arg.ID}); err != nil {
		return
	}
	if vpc == nil {
		err = ecode.VipPanelConfNotExistErr
		return
	}
	// if vpc.CheckProductID(arg) {
	// 	err = ecode.VipPanelProductIDNotNilErr
	// 	return
	// }
	if vpc.OPrice < arg.DPrice {
		err = ecode.VipPanelValidDPriceErr
		return
	}
	if vpc.SubType != model.AutoRenew && arg.FirstPrice > 0 {
		return ecode.VipPanelFirstPriceNotSupportErr
	}
	if arg.ETime != 0 && arg.STime >= arg.ETime {
		err = ecode.VipPanelSTGeqETErr
		return
	}
	// 面版冲突
	if mvd, err = s.dao.VipDPriceConfigUQTime(c, arg); err != nil {
		return
	}
	if len(mvd) > 0 {
		err = ecode.VipPanelValidTimeErr
		return
	}
	nvdc := &model.VipDPriceConfig{
		ID:         arg.ID,
		PdID:       arg.PdID,
		DPrice:     arg.DPrice,
		STime:      arg.STime,
		ETime:      arg.ETime,
		Remark:     arg.Remark,
		Operator:   arg.Operator,
		OpID:       arg.OpID,
		FirstPrice: arg.FirstPrice,
	}
	err = s.dao.AddVipDPriceConfig(c, nvdc)
	return
}

// UpVipDPriceConfig update vip discount price config.
func (s *Service) UpVipDPriceConfig(c context.Context, arg *model.ArgAddOrUpVipDPrice) (err error) {
	var (
		ovpc *model.VipPriceConfig
		ovdc *model.VipDPriceConfig
		mvd  map[int64]*model.VipDPriceConfig
	)
	if ovpc, err = s.dao.VipPriceConfigID(c, &model.ArgVipPriceID{ID: arg.ID}); err != nil {
		return
	}
	if ovpc == nil {
		err = ecode.VipPanelConfNotExistErr
		return
	}
	if ovdc, err = s.dao.VipDPriceConfigID(c, &model.ArgVipDPriceID{DisID: arg.DisID}); err != nil {
		return
	}
	// if ovpc.CheckProductID(arg) {
	// 	err = ecode.VipPanelProductIDNotNilErr
	// 	return
	// }
	if ovpc.SubType != model.AutoRenew && arg.FirstPrice > 0 {
		return ecode.VipPanelFirstPriceNotSupportErr
	}
	if ovpc.OPrice < arg.DPrice {
		err = ecode.VipPanelValidDPriceErr
		return
	}
	if arg.ETime != 0 && arg.STime >= arg.ETime {
		err = ecode.VipPanelSTGeqETErr
		return
	}
	if arg.STime != ovdc.STime || arg.ETime != ovdc.ETime {
		// 面版冲突
		if mvd, err = s.dao.VipDPriceConfigUQTime(c, arg); err != nil {
			return
		}
		if _, ok := mvd[arg.DisID]; ok {
			delete(mvd, arg.DisID)
		}
		if len(mvd) > 0 {
			err = ecode.VipPanelValidTimeErr
			return
		}
	}
	vdc := &model.VipDPriceConfig{
		DisID:      arg.DisID,
		ID:         arg.ID,
		PdID:       arg.PdID,
		DPrice:     arg.DPrice,
		STime:      arg.STime,
		ETime:      arg.ETime,
		Remark:     arg.Remark,
		Operator:   arg.Operator,
		OpID:       arg.OpID,
		FirstPrice: arg.FirstPrice,
	}
	err = s.dao.UpVipDPriceConfig(c, vdc)
	return
}

// DelVipDPriceConfig delete vip discount price config.
func (s *Service) DelVipDPriceConfig(c context.Context, arg *model.ArgVipDPriceID) (err error) {
	err = s.dao.DelVipDPriceConfig(c, arg)
	return
}

// PanelPlatFormTypes .
func (s *Service) PanelPlatFormTypes(c context.Context) (res []*model.TypePlatform, err error) {
	if res, err = s.dao.PlatformTypes(c); err != nil {
		return
	}
	return
}
