package relation

import (
	"context"
	"github.com/pkg/errors"
	ServiceConf "go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	relationV1 "go-common/app/service/live/relation/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
	"time"
)

// GetGiftInfo 获取给我关注的人送礼(金瓜子)信息
func (d *Dao) GetGiftInfo(ctx context.Context, params ServiceConf.ChunkCallInfo) (giftInfo map[int64]int64, err error) {
	relationParams := &relationV1.BaseInfoGetGiftInfoReq{}
	ret, err := dao.RelationApi.V1BaseInfo.GetGiftInfo(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), relationParams)
	if err != nil {
		return
	}
	erelongInfo, success := ServiceConf.CheckReturn(err, ret.Code, ret.Msg, "gift", params.RPCTimeout, params.ChunkSize, params.ChunkNum)
	if !success {
		if err != nil {
			err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		} else {
			err = errors.WithMessage(ecode.PkIDLiveRPCCodeError, "GET SEA PATROL FAIL")
			log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
				err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		}
		return giftInfo, nil
	}
	if ret.Data == nil || len(ret.Data) < 0 {
		erelongInfo.ErrType = ServiceConf.EmptyResultEn
		erelongInfo.ErrDesc = ServiceConf.EmptyResult
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d",
			erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		return giftInfo, nil
	}
	for _, v := range ret.Data {
		giftInfo[v.Mid] = v.Gold
	}
	return
}

// GetAttentionListGroupBy 获取分组的关注列表
func (d *Dao) GetAttentionListGroupBy(ctx context.Context, params ServiceConf.ChunkCallInfo) (relationInfo map[int64]*relationV1.BaseInfoGetFollowTypeResp_UidInfo, attentionErr error) {
	attentionData, attentionErr := dao.RelationApi.V1BaseInfo.GetFollowType(
		rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond),
		&relationV1.BaseInfoGetFollowTypeReq{})
	if attentionErr != nil || attentionData == nil {
		attentionErr = ecode.AttentionListRPCError
		return
	}
	relationInfo = attentionData.Data
	return
}
