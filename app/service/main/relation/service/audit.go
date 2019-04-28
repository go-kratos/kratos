package service

import (
	"context"
	mmodel "go-common/app/service/main/member/model"
	bmodel "go-common/app/service/main/member/model/block"
	"go-common/app/service/main/relation/model"
	"go-common/library/log"
)

// Audit get member audit info
func (s *Service) Audit(c context.Context, mid int64, realIP string) (rs *model.Audit, err error) {
	rs = &model.Audit{
		Mid: mid,
	}
	var (
		detail    *model.PassportDetail
		blockInfo *bmodel.RPCResInfo
		baseInfo  *mmodel.BaseInfo
	)

	// get bindMainStatus and bindTelStatus from passport-service
	if detail, err = s.dao.PassportDetail(c, mid, realIP); err != nil {
		log.Error("s.accRPC.PassportDetail() error(%v) return(%v)", err, detail)
		return
	}
	rs.BindMail = bindEmailStatus(detail.Email, detail.Spacesta)
	rs.BindTel = bindPhoneStatus(detail.Phone)

	// get block status from block-service
	blockArg := &bmodel.RPCArgInfo{
		MID: mid,
	}
	if blockInfo, err = s.memberRPC.BlockInfo(c, blockArg); err != nil {
		log.Error("s.memberRPC.BlockInfo() error(%v) return(%v)", err, blockInfo)
		return
	}
	if blockInfo.BlockStatus != bmodel.BlockStatusFalse {
		rs.Blocked = true
	}

	// get rank from member-service
	memberArg := &mmodel.ArgMemberMid{
		Mid:      mid,
		RemoteIP: realIP,
	}
	if baseInfo, err = s.memberRPC.Base(c, memberArg); err != nil {
		log.Error("s.memberRPC.Base() error(%v) return(%v)", err, baseInfo)
		return
	}
	rs.Rank = baseInfo.Rank

	return
}

func bindEmailStatus(email string, spacesta int8) bool {
	return spacesta > -10 && len(email) > 0
}

func bindPhoneStatus(phone string) bool {
	return len(phone) > 0
}
