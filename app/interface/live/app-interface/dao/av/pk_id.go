package av

import (
	"context"
	"github.com/pkg/errors"
	ServiceConf "go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	avV1 "go-common/app/service/live/av/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
	"time"
)

// GetPkID ...
// 获取PkId信息
func (d *Dao) GetPkID(ctx context.Context, req *avV1.PkGetPkIdsByRoomIdsReq, params ServiceConf.ChunkCallInfo) (avResult *avV1.PkGetPkIdsByRoomIdsResp, err error) {
	avResult = &avV1.PkGetPkIdsByRoomIdsResp{}
	ret, err := dao.AvApi.V1Pk.GetPkIdsByRoomIds(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &avV1.PkGetPkIdsByRoomIdsReq{RoomIds: req.RoomIds, Platform: req.Platform})
	if err != nil {
		return
	}
	erelongInfo, success := ServiceConf.CheckReturn(err, ret.Code, ret.Msg, ServiceConf.GetPkIdsByRoomIds, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
	if !success {
		if err != nil {
			err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.RoomIds)
		} else {
			err = errors.WithMessage(ecode.PkIDLiveRPCCodeError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.RoomIds)
		}
		return
	}
	if ret.Data == nil || len(ret.Data) <= 0 {
		erelongInfo.ErrType = ServiceConf.EmptyResultEn
		erelongInfo.ErrDesc = ServiceConf.EmptyResult
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
			erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.RoomIds)
		return
	}
	avResult = ret
	return
}
