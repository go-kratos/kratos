package dao

import (
	"context"

	"go-common/app/service/live/av/api/liverpc/v1"
	"go-common/library/log"
)

// GetFansMedalInfo 获取粉丝勋章信息
func (d *Dao) GetPkStatus(c context.Context, roomId int64) (resp *v1.PkGetPkStatusResp_Data, err error) {
	reply, err := d.AvApi.V1Pk.GetPkStatus(c, &v1.PkGetPkStatusReq{RoomId: roomId})
	if err != nil {
		log.Error("av_GetPkStatus_error:%v", err)
		return
	}
	if reply.Code != 0 {
		log.Error("av_GetPkStatus_error:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
		return
	}
	log.Info("av_GetPkStatus:%d,%s,$v", reply.Code, reply.Msg, reply.Data)
	resp = reply.Data
	return
}
