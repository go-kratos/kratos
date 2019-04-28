package room_ex

import (
	"context"
	"github.com/pkg/errors"
	ServiceConf "go-common/app/interface/live/app-interface/conf"
	cDao "go-common/app/interface/live/app-interface/dao"
	roomExV1 "go-common/app/service/live/room_ex/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
	"time"
)

func (d *Dao) GetRoomNewsInfo(ctx context.Context, req *roomExV1.RoomNewsMultiGetReq, params ServiceConf.ChunkCallInfo) (roomNewsResult *roomExV1.RoomNewsMultiGetResp, err error) {
	roomNewsResult = &roomExV1.RoomNewsMultiGetResp{}
	ret, err := cDao.RoomExtApi.V1RoomNews.MultiGet(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond),
		&roomExV1.RoomNewsMultiGetReq{RoomIds: req.RoomIds, IsDecoded: req.IsDecoded})
	if err != nil {
		return
	}
	erelongInfo, success := ServiceConf.CheckReturn(err, ret.Code, ret.Msg, ServiceConf.RoomNews, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
	if !success {
		if err != nil {
			err = errors.WithMessage(ecode.RoomNewsRecordFrameWorkCallError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.RoomIds)
		} else {
			err = errors.WithMessage(ecode.RoomNewsLiveRPCCodeError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.RoomIds)
		}
		return nil, err
	}
	if ret.Data == nil {
		erelongInfo.ErrType = ServiceConf.EmptyResultEn
		erelongInfo.ErrDesc = ServiceConf.EmptyResult
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
			erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.RoomIds)
		return nil, err
	}
	roomNewsResult = ret
	return
}
