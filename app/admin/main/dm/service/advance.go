package service

import (
	"context"

	"go-common/app/admin/main/dm/model"
	accountApi "go-common/app/service/main/account/api"
	"go-common/library/log"
)

// Advances 高级弹幕列表
func (s *Service) Advances(c context.Context, dmInid int64, typ, mode string, pn, ps int64) (res []*model.Advance, total int64, err error) {
	var mids = make([]int64, 0)
	if res, total, err = s.dao.Advances(c, dmInid, typ, mode, pn, ps); err != nil {
		log.Error("dao.Advances(cid:%d, typ:%s, mode:%s, pn:%d,ps:%d) error(%v)", dmInid, typ, mode, pn, ps, err)
		return
	}
	for _, r := range res {
		mids = append(mids, r.Mid)
	}
	arg := &accountApi.MidsReq{Mids: mids}
	uInfos, err := s.accountRPC.Infos3(c, arg)
	if err != nil {
		log.Error("s.accRPC.Infos3(%v) error(%v)", mids, err)
		return
	}
	for _, r := range res {
		if v, ok := uInfos.GetInfos()[r.Mid]; ok {
			r.Name = v.Name
		}
	}
	return
}
