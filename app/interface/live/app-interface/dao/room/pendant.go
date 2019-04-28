package room

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"go-common/app/interface/live/app-interface/conf"
	ServiceConf "go-common/app/interface/live/app-interface/conf"
	cDao "go-common/app/interface/live/app-interface/dao"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
)

// GetRoomPendant
func (d *Dao) GetRoomPendant(ctx context.Context, roomIds []int64, pendantType string, position int64) (pendantRoomListResp map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result, err error) {
	pendantRoomListResp = make(map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result)
	getPendantByIdsTimeout := time.Duration(conf.GetTimeout("getPendantByIds", 50)) * time.Millisecond
	pendantRoomList, err := cDao.RoomApi.V1RoomPendant.GetPendantByIds(rpcCtx.WithTimeout(ctx, getPendantByIdsTimeout), &roomV1.RoomPendantGetPendantByIdsReq{
		Ids:      roomIds,
		Type:     pendantType,
		Position: position, // 历史原因，取右上，但客户端展示在左上
	})

	if err != nil {
		log.Error("[GetRoomPendant]room.v1.getPendantByIds rpc error:%+v", err)
		err = errors.WithMessage(ecode.RoomPendantError, fmt.Sprintf("room.v1.getPendantByIds rpc error:%+v", err))
		return
	}
	if pendantRoomList.Code != 0 {
		log.Error("[GetRoomPendant]room.v1.getPendantByIds response code:%d,msg:%s", pendantRoomList.Code, pendantRoomList.Msg)
		err = errors.WithMessage(ecode.RoomPendantReturnError, fmt.Sprintf("room.v1.getPendantByIds response error, code:%d, msg:%s", pendantRoomList.Code, pendantRoomList.Msg))
		return
	}
	if pendantRoomList.Data == nil || pendantRoomList.Data.Result == nil {
		log.Error("[GetRoomPendant]room.v1.getPendantByIds empty error")
		err = errors.WithMessage(ecode.RoomPendantReturnError, "")
		return
	}
	pendantRoomListResp = pendantRoomList.Data.Result
	return
}

// GetRoomPendantInfo ...
// 获取角标信息
func (d *Dao) GetRoomPendantInfo(ctx context.Context, req *roomV1.RoomPendantGetPendantByIdsReq, params ServiceConf.ChunkCallInfo) (roomNewsResult *roomV1.RoomPendantGetPendantByIdsResp, err error) {
	ret, err := cDao.RoomApi.V1RoomPendant.GetPendantByIds(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &roomV1.RoomPendantGetPendantByIdsReq{Ids: req.Ids, Type: req.Type, Position: req.Position})
	if err != nil {
		return
	}
	erelongInfo, success := ServiceConf.CheckReturn(err, ret.Code, ret.Msg, ServiceConf.RoomPendent, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
	if !success {
		if err != nil {
			err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.Ids)
		} else {
			err = errors.WithMessage(ecode.PkIDLiveRPCCodeError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.Ids)
		}
		return nil, err
	}
	if ret.Data == nil {
		erelongInfo.ErrType = ServiceConf.EmptyResultEn
		erelongInfo.ErrDesc = ServiceConf.EmptyResult
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
			erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.Ids)
		return nil, err
	}
	roomNewsResult = ret
	return
}
