package live_data

import (
	"context"
	"time"

	"github.com/pkg/errors"

	ServiceConf "go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	liveDataV1 "go-common/app/service/live/live_data/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
)

// GetLastLiveTime ...
// 获取Record信息
func (d *Dao) GetLastLiveTime(ctx context.Context, req *liveDataV1.RecordGetReq, params ServiceConf.ChunkCallInfo) (resp map[string]*liveDataV1.RecordGetResp_TimeInfo, err error) {
	resp = make(map[string]*liveDataV1.RecordGetResp_TimeInfo)
	ret, err := dao.LiveDataApi.V1Record.Get(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond),
		&liveDataV1.RecordGetReq{Roomids: req.Roomids})
	if err != nil {
		return
	}
	erelongInfo, success := ServiceConf.CheckReturn(err, ret.Code, ret.Msg, ServiceConf.Record, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
	if !success {
		if err != nil {
			err = errors.WithMessage(ecode.RecordRecordFrameWorkCallError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.Roomids)
		} else {
			err = errors.WithMessage(ecode.RecordLiveRPCCodeError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.Roomids)
		}
		return
	}
	if ret.Data == nil || len(ret.Data) <= 0 {
		erelongInfo.ErrType = ServiceConf.EmptyResultEn
		erelongInfo.ErrDesc = ServiceConf.EmptyResult
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
			erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.Roomids)
		return
	}
	resp = ret.Data
	return
}
