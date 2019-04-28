package relation

import (
	"context"
	"math"
	"strconv"
	"time"

	"go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	avV1 "go-common/app/service/live/av/api/liverpc/v1"
	fansMedalV1 "go-common/app/service/live/fans_medal/api/liverpc/v1"
	liveDataV1 "go-common/app/service/live/live_data/api/liverpc/v1"
	relationV1 "go-common/app/service/live/relation/api/liverpc/v1"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	roomExV1 "go-common/app/service/live/room_ex/api/liverpc/v1"
	userExV1 "go-common/app/service/live/userext/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

// ChunkCallInfo ...
// 日志结构体
type ChunkCallInfo struct {
	ParamsName string
	URLName    string
	ChunkSize  int64
	ChunkNum   int64
	RPCTimeout int64
}

const (
	getStatusInfoByUfos = "room/v1/Room/get_status_info_by_uids"
	targetsWithMedal    = "fans_medal/v1/FansMedal/targetsWithMedal"
	uuid2roomed         = "room/v2/Room/room_id_by_uid_multi"
	record              = "live_data/v1/Record/get"
	getPkIdsByRoomIds   = "av/v1/Pk/getPkIdsByRoomIds"
	roomPendent         = "room/v1/RoomPendant/getPendantByIds"
	roomNews            = "/room_ex/v1/RoomNews/multiGet"
	relationGiftInfo    = "/relation/v1/BaseInfo/getGiftInfo"
	// AccountGRPC ...
	// 主站grpc用户信息
	AccountGRPC = "Cards3"
	// LiveUserExpGRPC ...
	// 直播用户经验grpc
	LiveUserExpGRPC = "xuserExp"
	// FansNum ...
	// 直播粉丝
	FansNum = "GetUserFcBatch"
	// LiveDomain implementation
	// 域名
	LiveDomain = "http://live.bilibili.com/"
	// BoastURL implementation
	// 秒开url
	BoastURL      = "?broadcast_type="
	emptyResult   = "调用直播服务返回data为空"
	emptyResultEn = "got_empty_result"
	// GoRoutingErr ...
	// 协程wait错误
	GoRoutingErr = "协程等待数据错误"
	// App533CardType implementation
	// 大卡类型
	App533CardType = 1
	// PendentMobileBadge implementation
	// 角标类型
	PendentMobileBadge = "mobile_index_badge"
	// PendentPosition implementation
	// 角标位置
	PendentPosition = 2
	// App531GrayRule implementation
	// 灰度策略
	App531GrayRule = "r_big_card"
	// App536GrayRule implementation
	// 灰度策略
	App536GrayRule = "r_homepage_card536"
	// SelfUID implementation
	// 调试UID
	SelfUID = 22973824
	// DummyUIDEnable implementation
	// 调试开
	DummyUIDEnable = 1
)

// UIDs2roomIDs ...
// uid转换roomID,每批最大400
func UIDs2roomIDs(ctx context.Context, ufos []int64) (rolaids map[int64]int64, err error) {
	rpcChunkSize, RPCTimeout, err := GetChunkInfo(uuid2roomed)
	params := ChunkCallInfo{ParamsName: "ufos", URLName: uuid2roomed, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	rolaids = make(map[int64]int64)
	lens := len(ufos)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([]map[string]string, params.ChunkNum)
	wg, _ := errgroup.WithContext(ctx)

	for i := int64(1); i <= params.ChunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkUfosIds := make([]int64, 20)
			if x == params.ChunkNum {
				chunkUfosIds = ufos[(x-1)*params.ChunkSize:]
			} else {
				chunkUfosIds = ufos[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := dao.RoomApi.V2Room.RoomIdByUidMulti(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &roomV2.RoomRoomIdByUidMultiReq{Uids: chunkUfosIds})
			if err != nil {
				ret = &roomV2.RoomRoomIdByUidMultiResp{}
				ret.Code = -1
				ret.Msg = "liveprc_error"
			}
			erelongInfo, success := CheckReturn(err, ret.Code, ret.Msg, uuid2roomed, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
			if !success {
				if err != nil {
					err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkUfosIds)
				}
				return nil
			}
			if ret.Data == nil || len(ret.Data) <= 0 {
				erelongInfo.ErrType = emptyResultEn
				erelongInfo.ErrDesc = emptyResult
				// log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				// 	erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkUfosIds)
				return nil
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = uuid2roomed
		erelongInfo.ErrDesc = GoRoutingErr
		erelongInfo.Code = 1003001
		erelongInfo.RPCTimeout = params.RPCTimeout
		erelongInfo.ErrorPtr = &err
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
			*erelongInfo.ErrorPtr, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		err = errors.WithMessage(ecode.RelationFrameWorkGoRoutingError, "GET SEA PATROL FAIL")
		return nil, err
	}
	// 整理数据
	for _, chunkItemList := range chunkResult {
		for k, item := range chunkItemList {
			if item != "" {
				Index := RParseInt(k, 1)
				itemInt := RParseInt(item, 1)
				rolaids[Index] = itemInt
			}
		}
	}
	return
}

// GetRoomInfo ...
// 获取room信息
func GetRoomInfo(ctx context.Context, input *roomV1.RoomGetStatusInfoByUidsReq) (roomResult map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo, err error) {
	rpcChunkSize, RPCTimeout, err := GetChunkInfo(getStatusInfoByUfos)
	params := ChunkCallInfo{ParamsName: "uids", URLName: getStatusInfoByUfos, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	roomResult = make(map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo)
	lens := len(input.Uids)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([]map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo, params.ChunkNum)
	wg, _ := errgroup.WithContext(ctx)

	for i := int64(1); i <= params.ChunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkUfosIds := make([]int64, 20)
			if x == params.ChunkNum {
				chunkUfosIds = input.Uids[(x-1)*params.ChunkSize:]
			} else {
				chunkUfosIds = input.Uids[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := dao.RoomApi.V1Room.GetStatusInfoByUids(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &roomV1.RoomGetStatusInfoByUidsReq{Uids: chunkUfosIds, FilterOffline: input.FilterOffline, NeedBroadcastType: input.NeedBroadcastType})
			if err != nil {
				if err != nil {
					ret = &roomV1.RoomGetStatusInfoByUidsResp{}
					ret.Code = -1
					ret.Msg = "liveprc_error"
				}
			}
			erelongInfo, success := CheckReturn(err, ret.Code, ret.Msg, getStatusInfoByUfos, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
			if !success {
				if err != nil {
					err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkUfosIds)
				} else {
					err = errors.WithMessage(ecode.PkIDLiveRPCCodeError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkUfosIds)
				}
				return nil
			}
			if ret.Data == nil || len(ret.Data) <= 0 {
				erelongInfo.ErrType = emptyResultEn
				erelongInfo.ErrDesc = emptyResult
				// log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				// 	erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkUfosIds)
				return nil
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = getStatusInfoByUfos
		erelongInfo.ErrDesc = GoRoutingErr
		erelongInfo.Code = 1003001
		erelongInfo.RPCTimeout = params.RPCTimeout
		erelongInfo.ErrorPtr = &err
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
			*erelongInfo.ErrorPtr, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		err = errors.WithMessage(ecode.RoomFrameWorkGoRoutingError, "GET SEA PATROL FAIL")
		return nil, err
	}
	// 整理数据
	for _, chunkItemList := range chunkResult {
		for _, item := range chunkItemList {
			if item != nil {
				roomResult[item.Uid] = item
			}
		}
	}
	return
}

// GetLastLiveTime ...
// 获取Record信息
func GetLastLiveTime(ctx context.Context, rolaids []int64) (literature map[string]string, err error) {
	rpcChunkSize, RPCTimeout, err := GetChunkInfo(record)
	params := ChunkCallInfo{ParamsName: "rolaids", URLName: record, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	literature = make(map[string]string)
	lens := len(rolaids)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([]map[string]*liveDataV1.RecordGetResp_TimeInfo, params.ChunkNum)
	wg, _ := errgroup.WithContext(ctx)

	for i := int64(1); i <= params.ChunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkRoomIds := make([]int64, 20)
			if x == params.ChunkNum {
				chunkRoomIds = rolaids[(x-1)*params.ChunkSize:]
			} else {
				chunkRoomIds = rolaids[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := dao.LiveDataApi.V1Record.Get(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &liveDataV1.RecordGetReq{Roomids: chunkRoomIds})
			if err != nil {
				if err != nil {
					ret = &liveDataV1.RecordGetResp{}
					ret.Code = -1
					ret.Msg = "liveprc_error"
				}
			}
			erelongInfo, success := CheckReturn(err, ret.Code, ret.Msg, record, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
			if !success {
				if err != nil {
					err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				} else {
					err = errors.WithMessage(ecode.PkIDLiveRPCCodeError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				}
				return nil
			}
			if ret.Data == nil || len(ret.Data) <= 0 {
				erelongInfo.ErrType = emptyResultEn
				erelongInfo.ErrDesc = emptyResult
				// log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				// 	erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				return nil
			}
			// chunkResult = append(chunkResult, ret.Data)
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = record
		erelongInfo.ErrDesc = GoRoutingErr
		erelongInfo.Code = 1003001
		erelongInfo.RPCTimeout = params.RPCTimeout
		erelongInfo.ErrorPtr = &err
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
			*erelongInfo.ErrorPtr, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		err = errors.WithMessage(ecode.RecordFrameWorkGoRoutingError, "GET SEA PATROL FAIL")
		return nil, err
	}
	// 整理数据
	for _, chunkItemList := range chunkResult {
		for k, item := range chunkItemList {
			if item != nil {
				literature[k] = item.RecentEndTime
			}
		}
	}
	return
}

// GetRoomNewsInfo ...
// 获取公告信息
func GetRoomNewsInfo(ctx context.Context, rolaids *roomExV1.RoomNewsMultiGetReq) (roomNewsResult map[int64]*roomExV1.RoomNewsMultiGetResp_Data, err error) {
	rpcChunkSize, RPCTimeout, err := GetChunkInfo(roomNews)
	params := ChunkCallInfo{ParamsName: "rolaids", URLName: roomNews, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	roomNewsResult = make(map[int64]*roomExV1.RoomNewsMultiGetResp_Data)
	lens := len(rolaids.RoomIds)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([][]*roomExV1.RoomNewsMultiGetResp_Data, params.ChunkNum)
	wg, _ := errgroup.WithContext(ctx)
	for i := int64(1); i <= params.ChunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkRoomIds := make([]int64, 20)
			if x == params.ChunkNum {
				chunkRoomIds = rolaids.RoomIds[(x-1)*params.ChunkSize:]
			} else {
				chunkRoomIds = rolaids.RoomIds[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := dao.RoomExtApi.V1RoomNews.MultiGet(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &roomExV1.RoomNewsMultiGetReq{RoomIds: chunkRoomIds, IsDecoded: rolaids.IsDecoded})
			if err != nil {
				if err != nil {
					ret = &roomExV1.RoomNewsMultiGetResp{}
					ret.Code = -1
					ret.Msg = "liveprc_error"
				}
			}
			erelongInfo, success := CheckReturn(err, ret.Code, ret.Msg, roomNews, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
			if !success {
				if err != nil {
					err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				} else {
					err = errors.WithMessage(ecode.PkIDLiveRPCCodeError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				}
				return nil
			}
			if ret.Data == nil {
				erelongInfo.ErrType = emptyResultEn
				erelongInfo.ErrDesc = emptyResult
				// log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				// 	erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				return nil
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = roomNews
		erelongInfo.ErrDesc = GoRoutingErr
		erelongInfo.Code = 1003001
		erelongInfo.RPCTimeout = params.RPCTimeout
		erelongInfo.ErrorPtr = &err
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
			*erelongInfo.ErrorPtr, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		err = errors.WithMessage(ecode.RoomNewsFrameWorkGoRoutingError, "GET SEA PATROL FAIL")
		return nil, err
	}
	// 整理数据
	for _, chunkItemList := range chunkResult {
		for _, item := range chunkItemList {
			if item != nil {
				if mid, err := strconv.ParseInt(item.Roomid, 10, 64); err == nil {
					roomNewsResult[mid] = item
				}

			}
		}
	}
	return
}

// GetRoomPendantInfo ...
// 获取角标信息
func GetRoomPendantInfo(ctx context.Context, req *roomV1.RoomPendantGetPendantByIdsReq) (roomNewsResult map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result, err error) {
	rpcChunkSize, RPCTimeout, err := GetChunkInfo(roomPendent)
	params := ChunkCallInfo{ParamsName: "ids", URLName: roomPendent, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	roomNewsResult = make(map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result)
	lens := len(req.Ids)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([]map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result, params.ChunkNum)
	wg, _ := errgroup.WithContext(ctx)
	for i := int64(1); i <= params.ChunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkRoomIds := make([]int64, 20)
			if x == params.ChunkNum {
				chunkRoomIds = req.Ids[(x-1)*params.ChunkSize:]
			} else {
				chunkRoomIds = req.Ids[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := dao.RoomApi.V1RoomPendant.GetPendantByIds(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &roomV1.RoomPendantGetPendantByIdsReq{Ids: chunkRoomIds, Type: req.Type, Position: req.Position})
			if err != nil {
				if err != nil {
					ret = &roomV1.RoomPendantGetPendantByIdsResp{}
					ret.Code = -1
					ret.Msg = "liveprc_error"
				}
			}
			erelongInfo, success := CheckReturn(err, ret.Code, ret.Msg, roomPendent, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
			if !success {
				if err != nil {
					err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				} else {
					err = errors.WithMessage(ecode.PkIDLiveRPCCodeError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				}
				return nil
			}
			if ret.Data == nil {
				erelongInfo.ErrType = emptyResultEn
				erelongInfo.ErrDesc = emptyResult
				return nil
			}
			chunkResult[x-1] = ret.Data.Result
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = roomPendent
		erelongInfo.ErrDesc = GoRoutingErr
		erelongInfo.Code = 1003001
		erelongInfo.RPCTimeout = params.RPCTimeout
		erelongInfo.ErrorPtr = &err
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
			*erelongInfo.ErrorPtr, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		err = errors.WithMessage(ecode.RoomPendentFrameWorkGoRoutingError, "GET SEA PATROL FAIL")
		return nil, err
	}
	// 整理数据
	for _, chunkItemList := range chunkResult {
		for k, item := range chunkItemList {
			if item != nil {
				roomNewsResult[k] = item

			}
		}
	}
	return
}

// GetPkID ...
// 获取PkId信息
func GetPkID(ctx context.Context, req *avV1.PkGetPkIdsByRoomIdsReq) (avResult map[string]int64, err error) {
	rpcChunkSize, RPCTimeout, err := GetChunkInfo(getPkIdsByRoomIds)
	params := ChunkCallInfo{ParamsName: "roomids", URLName: getPkIdsByRoomIds, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	avResult = make(map[string]int64)
	lens := len(req.RoomIds)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([]map[string]int64, params.ChunkNum)
	wg, _ := errgroup.WithContext(ctx)
	for i := int64(1); i <= params.ChunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkRoomIds := make([]int64, 20)
			if x == params.ChunkNum {
				chunkRoomIds = req.RoomIds[(x-1)*params.ChunkSize:]
			} else {
				chunkRoomIds = req.RoomIds[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := dao.AvApi.V1Pk.GetPkIdsByRoomIds(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &avV1.PkGetPkIdsByRoomIdsReq{RoomIds: chunkRoomIds, Platform: req.Platform})
			if err != nil {
				if err != nil {
					ret = &avV1.PkGetPkIdsByRoomIdsResp{}
					ret.Code = -1
					ret.Msg = "liveprc_error"
				}
			}
			erelongInfo, success := CheckReturn(err, ret.Code, ret.Msg, getPkIdsByRoomIds, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
			if !success {
				if err != nil {
					err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				} else {
					err = errors.WithMessage(ecode.PkIDLiveRPCCodeError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				}
				return nil
			}
			if ret.Data == nil || len(ret.Data) <= 0 {
				erelongInfo.ErrType = emptyResultEn
				erelongInfo.ErrDesc = emptyResult
				// log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				// 	erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				return nil
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = getPkIdsByRoomIds
		erelongInfo.ErrDesc = GoRoutingErr
		erelongInfo.Code = 1003001
		erelongInfo.RPCTimeout = params.RPCTimeout
		erelongInfo.ErrorPtr = &err
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
			*erelongInfo.ErrorPtr, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		err = errors.WithMessage(ecode.PkIDFrameWorkGoRoutingError, "GET SEA PATROL FAIL")
		return nil, err
	}
	// 整理数据
	for _, chunkItemList := range chunkResult {
		for k, item := range chunkItemList {
			avResult[k] = item
		}
	}
	return
}

// GetFansMedal ...
// 获取粉丝勋章佩戴信息
func GetFansMedal(ctx context.Context, req *fansMedalV1.FansMedalTargetsWithMedalReq) (fansResult map[int64]bool, err error) {
	rpcChunkSize, RPCTimeout, err := GetChunkInfo(targetsWithMedal)
	params := ChunkCallInfo{ParamsName: "target_ids", URLName: targetsWithMedal, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	fansResult = make(map[int64]bool)
	lens := len(req.TargetIds)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([][]int64, params.ChunkNum)
	wg, _ := errgroup.WithContext(ctx)
	for i := int64(1); i <= params.ChunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkRoomIds := make([]int64, 20)
			if x == params.ChunkNum {
				chunkRoomIds = req.TargetIds[(x-1)*params.ChunkSize:]
			} else {
				chunkRoomIds = req.TargetIds[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := dao.FansMedalApi.V1FansMedal.TargetsWithMedal(rpcCtx.WithTimeout(ctx, time.Duration(params.RPCTimeout)*time.Millisecond), &fansMedalV1.FansMedalTargetsWithMedalReq{Uid: req.Uid, TargetIds: chunkRoomIds})
			if err != nil {
				if err != nil {
					ret = &fansMedalV1.FansMedalTargetsWithMedalResp{}
					ret.Code = -1
					ret.Msg = "liveprc_error"
				}
			}
			erelongInfo, success := CheckReturn(err, ret.Code, ret.Msg, targetsWithMedal, params.RPCTimeout, params.ChunkSize, params.ChunkNum)
			if !success {
				if err != nil {
					err = errors.WithMessage(ecode.PkIDRecordFrameWorkCallError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				} else {
					err = errors.WithMessage(ecode.PkIDLiveRPCCodeError, "GET SEA PATROL FAIL")
					log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
						err, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				}
				return nil
			}
			if ret.Data == nil || len(ret.Data) <= 0 {
				erelongInfo.ErrType = emptyResultEn
				erelongInfo.ErrDesc = emptyResult
				// log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s"+"|Params:%v",
				// 	erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName, chunkRoomIds)
				return nil
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = targetsWithMedal
		erelongInfo.ErrDesc = GoRoutingErr
		erelongInfo.Code = 1003001
		erelongInfo.RPCTimeout = params.RPCTimeout
		erelongInfo.ErrorPtr = &err
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
			*erelongInfo.ErrorPtr, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		err = errors.WithMessage(ecode.FansMedalFrameWorkGoRoutingError, "GET SEA PATROL FAIL")
		return nil, err
	}
	// 整理数据
	for _, chunkItemList := range chunkResult {
		for _, item := range chunkItemList {
			fansResult[item] = true
		}
	}
	return
}

// GetGrayRule ...
// 获取灰度规则信息
func GetGrayRule(ctx context.Context, req *userExV1.GrayRuleGetByMarkReq) (extResult *userExV1.GrayRuleGetByMarkResp_Data, err error) {
	extResult = &userExV1.GrayRuleGetByMarkResp_Data{}
	if req == nil {
		return nil, nil
	}
	ret, err := dao.UserExtApi.V1GrayRule.GetByMark(ctx, req)
	if err != nil {
		log.Error("call_userExt_grayRule error,err:%v", err)
		err = errors.WithMessage(ecode.GetGrayRuleError, "GET SEA PATROL FAIL")
		return
	}
	extResult = ret.Data
	return
}

// GetGiftInfo ...
// 获取送礼信息
func GetGiftInfo(ctx context.Context) (giftInfo map[int64]int64, err error) {
	_, RPCTimeout, _ := GetChunkInfo(relationGiftInfo)
	relationParams := &relationV1.BaseInfoGetGiftInfoReq{}
	giftInfo = make(map[int64]int64)
	ret, err := dao.RelationApi.V1BaseInfo.GetGiftInfo(ctx, relationParams)
	if err != nil {
		if err != nil {
			ret = &relationV1.BaseInfoGetGiftInfoResp{}
			ret.Code = -1
			ret.Msg = "liveprc_error"
		}
	}
	params := ChunkCallInfo{ParamsName: "", URLName: relationGiftInfo, ChunkSize: 1, RPCTimeout: RPCTimeout}
	erelongInfo, success := CheckReturn(err, ret.Code, ret.Msg, "gift", params.RPCTimeout, params.ChunkSize, params.ChunkNum)
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
		erelongInfo.ErrType = emptyResultEn
		erelongInfo.ErrDesc = emptyResult
		// log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d",
		// 	erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		return giftInfo, nil
	}
	for _, v := range ret.Data {
		giftInfo[v.Mid] = v.Gold
	}
	return
}

// GetChunkInfo ...
// 获取分块信息
func GetChunkInfo(rpcName string) (rpcChunkSize int64, RPCTimeout int64, err error) {
	rpcChunkSize = conf.GetChunkSize(rpcName, 20)
	RPCTimeout = conf.GetTimeout(rpcName, 100)
	return
}
