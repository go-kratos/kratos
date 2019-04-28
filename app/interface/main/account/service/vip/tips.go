package vip

import (
	"context"

	"go-common/app/interface/main/account/model"
	vipmol "go-common/app/service/main/vip/model"

	"github.com/pkg/errors"
)

// Tips vip tips info.
func (s *Service) Tips(c context.Context, req *model.TipsReq) (res *vipmol.TipsResp, err error) {
	var rs []*vipmol.TipsResp
	arg := &vipmol.ArgTips{
		Platform: req.Platform,
		Version:  req.Version,
		Position: req.Position,
	}
	if rs, err = s.vipRPC.Tips(c, arg); err != nil {
		err = errors.WithStack(err)
	}
	if len(rs) > 0 {
		res = rs[0]
	}
	return
}

// TipsV2 vip tips info v2.
func (s *Service) TipsV2(c context.Context, req *model.TipsReq) (res []*vipmol.TipsResp, err error) {
	arg := &vipmol.ArgTips{
		Platform: req.Platform,
		Version:  req.Version,
		Position: req.Position,
	}
	if res, err = s.vipRPC.Tips(c, arg); err != nil {
		err = errors.WithStack(err)
	}
	return
}
