package dao

import (
	"context"

	v12 "go-common/app/service/main/member/api"
	"go-common/library/log"
)

// GetIdentityStatus  获取身份申请信息
func (d *Dao) GetIdentityStatus(c context.Context, mid int64) (resp int8, err error) {
	reply, err := d.memberCli.RealnameApplyStatus(c, &v12.MemberMidReq{Mid: mid})
	if err != nil {
		log.Error("main_member_GetIdentityStatus_error:%v", err)
		return
	}
	log.Info("main_member_GetIdentityStatus:%v", reply)
	resp = reply.Status
	return
}
