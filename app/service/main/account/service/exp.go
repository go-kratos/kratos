package service

import (
	"context"

	mmodel "go-common/app/service/main/member/model"
	"go-common/library/net/metadata"
)

// AddExp add user exp.
func (s *Service) AddExp(c context.Context, mid int64, money float64, operater, operate, reason string) error {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mmodel.ArgAddExp{Mid: mid, Count: money, Operate: operate, Reason: reason, IP: ip}
	return s.dao.UpdateExp(c, arg)
}

// AddMoral add user moral.
func (s *Service) AddMoral(c context.Context, mid int64, moral float64, oper, reason, remark string) error {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mmodel.ArgUpdateMoral{Mid: mid, Reason: reason, Remark: remark, IP: ip, IsNotify: true}
	delta := int64(moral * 100)
	arg.Delta = delta
	arg.Origin = mmodel.ReportRewardType
	if delta < 0 {
		arg.Origin = mmodel.PunishmentType
	}
	arg.Operator = oper
	if len(oper) == 0 {
		arg.Operator = "系统"
	}
	// 目前只支持评论，reason_type 后续由业务方传递
	arg.ReasonType = mmodel.ReplyReasonType
	return s.dao.AddMoral(c, arg)
}
