package service

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"go-common/library/log"
)

const (
	ALL      int8 = 0
	MAIN_VIP int8 = 10
)

func (s *Service) GetPpcsMap() map[int8][]*model.PanelPriceConfig {
	return s.ppcsMap.Load().(map[int8][]*model.PanelPriceConfig)
}

func (s *Service) PanelPriceConfigsBySuitType(c context.Context, st int8) (ppcs []*model.PanelPriceConfig, err error) {
	ppcsMap := s.GetPpcsMap()
	ppcs, ok := ppcsMap[st]
	if !ok {
		return make([]*model.PanelPriceConfig, 0), nil
	}
	return ppcs, nil
}

func (s *Service) PanelPriceConfigByProductId(c context.Context, productId string) (ppc *model.PanelPriceConfig, err error) {
	ppcsMap := s.GetPpcsMap()
	for _, ppcs := range ppcsMap {
		for _, ppc = range ppcs {
			if ppc.ProductId == productId {
				return ppc, nil
			}
		}
	}
	return nil, nil
}

func (s *Service) PanelPriceConfigByPid(c context.Context, pid int32) (ppc *model.PanelPriceConfig, err error) {
	ppcsMap := s.GetPpcsMap()
	for _, ppcs := range ppcsMap {
		for _, ppc = range ppcs {
			if ppc.ID == pid {
				return ppc, nil
			}
		}
	}
	return nil, nil
}

func (s *Service) GuestPanelPriceConfigs(c context.Context) (ppcs []*model.PanelPriceConfig, err error) {
	return s.PanelPriceConfigsBySuitType(c, ALL)
}

func (s *Service) MVipPanelPriceConfigs(c context.Context, mid int64) (ppcs []*model.PanelPriceConfig, err error) {
	mv, err := s.dao.MainVip(c, mid)
	if err != nil {
		return nil, err
	}
	ppcs, err = s.PanelPriceConfigsBySuitType(c, MAIN_VIP)
	if err != nil {
		log.Error("s.PanelPriceConfigsBySuitType(%d), err(%v)", MAIN_VIP, err)
		return nil, err
	}
	if !mv.IsVip() || mv.Months() < 1 {
		return make([]*model.PanelPriceConfig, 0), nil
	}
	for _, ppc := range ppcs {
		ppc.MaxNum = mv.Months()
	}
	return ppcs, nil
}

func (s *Service) GuestPanelInfo(c context.Context) (pi map[int8][]*model.PanelPriceConfig, err error) {
	gpps, err := s.GuestPanelPriceConfigs(c)
	if err != nil {
		log.Error("s.GuestPanelPriceConfigs err(%v)", err)
		return
	}
	pi = make(map[int8][]*model.PanelPriceConfig)
	pi[ALL] = gpps
	return pi, nil
}

func (s *Service) PanelInfo(c context.Context, mid int64) (pi map[int8][]*model.PanelPriceConfig, err error) {
	gpps, err := s.GuestPanelPriceConfigs(c)
	if err != nil {
		log.Error("s.GuestPanelPriceConfigs err(%v)", err)
		return
	}
	mvpps, err := s.MVipPanelPriceConfigs(c, mid)
	if err != nil {
		log.Error("s.MVipPanelPriceConfigs(%d) err(%v)", mid, err)
		return
	}
	// NOTE: do we need to use uc replace ui?
	ui, err := s.dao.RawUserInfoByMid(c, mid)
	if err != nil {
		return
	}
	// NOTE: 包月用户无法看到包月套餐，即使用户vip过期
	if ui != nil && ui.PayType == model.VipPayTypeSub {
		for i := 0; i < len(gpps); i++ {
			if gpps[i].SubType == model.SubTypeContract {
				gpps = append(gpps[:i], gpps[i+1:]...)
				i--
			}
		}
		for i := 0; i < len(mvpps); i++ {
			if mvpps[i].SubType == model.SubTypeContract {
				mvpps = append(mvpps[:i], mvpps[i+1:]...)
				i--
			}
		}
	}
	// NOTE: tv 会员无法通过主站会员升级方式购买
	if ui != nil && ui.IsVip() {
		mvpps = make([]*model.PanelPriceConfig, 0)
	}
	pi = make(map[int8][]*model.PanelPriceConfig)
	pi[ALL] = gpps
	pi[MAIN_VIP] = mvpps
	return pi, nil
}
