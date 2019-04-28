package dao

import (
	"context"

	"go-common/app/service/live/fans_medal/api/liverpc/v1"
	"go-common/library/log"
)

//OpenFansMealLevel 开启粉丝勋章的主播等级
const OpenFansMealLevel = 10

// GetFansMedalInfo 获取粉丝勋章信息
func (d *Dao) GetFansMedalInfo(c context.Context, uid int64) (resp *v1.MedalQueryResp_Data, err error) {
	reply, err := d.FansMedalApi.V1Medal.Query(c, &v1.MedalQueryReq{Uid: uid, Source: 1})
	if err != nil {
		log.Error("fansMedal_getFansMedalInfo_error:%v", err)
		return
	}
	if reply.Code != 0 {
		log.Error("fansMedal_getFansMedalInfo_error:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
		return
	}
	log.Info("fansMedal_getFansMedalInfo:%d,%s,$v", reply.Code, reply.Msg, reply.Data)
	if reply.Data.IsNull {
		resp = nil
	}
	resp = reply.Data
	return
}
