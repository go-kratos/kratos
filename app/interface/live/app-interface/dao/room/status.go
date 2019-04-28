package room

import (
	"context"
	"github.com/pkg/errors"
	ServiceConf "go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
	"time"
)

// GetRoomInfo ...
// 获取room信息
func (d *Dao) GetRoomInfo(ctx context.Context, req *roomV1.RoomGetStatusInfoByUidsReq, params ServiceConf.ChunkCallInfo) (resp *roomV1.RoomGetStatusInfoByUidsResp, err error) {
	// ret, err := dao.FansMedalApi.V1FansMedal.TargetsWithMedal(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond),
	resp = &roomV1.RoomGetStatusInfoByUidsResp{}
	ret, err := dao.RoomApi.V1Room.GetStatusInfoByUids(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond),
		&roomV1.RoomGetStatusInfoByUidsReq{Uids: req.Uids, FilterOffline: req.FilterOffline, NeedBroadcastType: req.NeedBroadcastType})
	if err != nil {
		return
	}
	erelongInfo, success := ServiceConf.CheckReturn(err, ret.Code, ret.Msg, ServiceConf.GetStatusInfoByUfos, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
	if !success {
		if err != nil {
			err = errors.WithMessage(ecode.RoomGetStatusInfoRPCError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.Uids)
		} else {
			err = errors.WithMessage(ecode.RoomGetStatusInfoRPCError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.Uids)
		}
		return
	}
	if ret.Data == nil || len(ret.Data) <= 0 {
		erelongInfo.ErrType = ServiceConf.EmptyResultEn
		erelongInfo.ErrDesc = ServiceConf.EmptyResult
		// log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
		// 	erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.Uids)
		return
	}
	resp = ret
	return
}
