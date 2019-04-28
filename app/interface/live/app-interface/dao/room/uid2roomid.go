package room

import (
	"context"
	"github.com/pkg/errors"
	ServiceConf "go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
	"time"
)

// UIDs2roomIDs ...
// uid转换roomID,每批最大400
func (d *Dao) UIDs2roomIDs(ctx context.Context, req *roomV2.RoomRoomIdByUidMultiReq, params ServiceConf.ChunkCallInfo) (ret *roomV2.RoomRoomIdByUidMultiResp, err error) {
	ret, err = dao.RoomApi.V2Room.RoomIdByUidMulti(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &roomV2.RoomRoomIdByUidMultiReq{Uids: req.Uids})
	if err != nil {
		return
	}
	erelongInfo, success := ServiceConf.CheckReturn(err, ret.Code, ret.Msg, ServiceConf.GetRoomID, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
	if !success {
		if err != nil {
			err = errors.WithMessage(ecode.FansMedalFrameWorkCallError, "GET SEA PATROL FAIL")
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
	return
}
