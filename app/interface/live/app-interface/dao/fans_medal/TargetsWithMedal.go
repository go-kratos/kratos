package fans_medal

import (
	"context"
	"time"

	"github.com/pkg/errors"

	ServiceConf "go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	fansMedalV1 "go-common/app/service/live/fans_medal/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
)

// GetFansMedal ...
// 获取粉丝勋章佩戴信息
func (d *Dao) GetFansMedal(ctx context.Context, req *fansMedalV1.FansMedalTargetsWithMedalReq, params ServiceConf.ChunkCallInfo) (fansResult *fansMedalV1.FansMedalTargetsWithMedalResp, err error) {
	fansResult = &fansMedalV1.FansMedalTargetsWithMedalResp{}
	ret, err := dao.FansMedalApi.V1FansMedal.TargetsWithMedal(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond),
		&fansMedalV1.FansMedalTargetsWithMedalReq{Uid: req.Uid, TargetIds: req.TargetIds})
	if err != nil {
		return
	}
	erelongInfo, success := ServiceConf.CheckReturn(err, ret.Code, ret.Msg, ServiceConf.TargetsWithMedal, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
	if !success {
		if err != nil {
			err = errors.WithMessage(ecode.FansMedalFrameWorkCallError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.TargetIds)
		} else {
			err = errors.WithMessage(ecode.FansMedalLiveRPCCodeError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.TargetIds)
		}
		return
	}
	if ret.Data == nil || len(ret.Data) <= 0 {
		erelongInfo.ErrType = ServiceConf.EmptyResultEn
		erelongInfo.ErrDesc = ServiceConf.EmptyResult
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
			erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, req.TargetIds)
		return
	}
	fansResult = ret
	return
}
