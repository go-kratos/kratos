package service

import (
	"context"

	"go-common/app/service/main/vip/model"
)

func (s *Service) loadAssociateVip() (err error) {
	var list []*model.AssociateVipResp
	if list, err = s.dao.EffectiveAssociateVips(context.Background()); err != nil {
		return
	}
	tmpmap := map[int8][]*model.AssociateVipResp{}
	for _, v := range list {
		if tmpmap[v.AssociatePlatform] == nil {
			tmpmap[v.AssociatePlatform] = []*model.AssociateVipResp{}
		}
		tmpmap[v.AssociatePlatform] = append(tmpmap[v.AssociatePlatform], v)
	}
	s.associateVipMap = tmpmap
	return
}

// AssociateVips get associate vips.
func (s *Service) AssociateVips(c context.Context, arg *model.ArgAssociateVip) (res []*model.AssociateVipResp) {
	return s.associateVipMap[model.AssociatePlatform(arg.Platform, arg.Device, arg.MobiApp)]
}
