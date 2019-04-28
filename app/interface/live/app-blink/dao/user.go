package dao

import (
	"context"

	"go-common/app/service/live/user/api/liverpc/v3"
	"go-common/library/log"
)

//GetUserInfo 获取用户信息
func (d *Dao) GetUserInfo(c context.Context, uid []int64) (res map[int64]*v3.UserGetMultipleResp_UserInfo, err error) {
	reply, err := d.UserApi.V3User.GetMultiple(c, &v3.UserGetMultipleReq{
		Uids:       uid,
		Attributes: []string{"exp", "info"},
	})
	if err != nil {
		log.Error("%v", err)
		return
	}
	if reply.Code != 0 {
		log.Error("user_getUserInfo_error:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
		return
	}
	log.Info("user_getUserInfo:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
	res = reply.Data
	return
}
