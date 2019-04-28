package v2

import (
	"context"
	"math"
	"strconv"

	"github.com/pkg/errors"

	ServiceConf "go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	relationT "go-common/app/interface/live/app-interface/service/v1/relation"
	avV1 "go-common/app/service/live/av/api/liverpc/v1"
	fansMedalV1 "go-common/app/service/live/fans_medal/api/liverpc/v1"
	liveDataV1 "go-common/app/service/live/live_data/api/liverpc/v1"
	relationV1 "go-common/app/service/live/relation/api/liverpc/v1"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	roomExV1 "go-common/app/service/live/room_ex/api/liverpc/v1"
	userExV1 "go-common/app/service/live/userext/api/liverpc/v1"
	liveUserExpM "go-common/app/service/live/xuser/api/grpc/v1"
	accountM "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
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
	// GoRoutingErr ...
	// 协程wait错误
	GoRoutingErr = "协程等待数据错误"
)

// UIDs2roomIDs ...
// uid转换roomID,每批最大400
func (s *IndexService) UIDs2roomIDs(ctx context.Context, UIDs []int64) (roomIDs map[int64]int64, err error) {
	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(ServiceConf.GetRoomID)
	params := ServiceConf.ChunkCallInfo{ParamsName: "uids", URLName: ServiceConf.GetRoomID, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	roomIDs = make(map[int64]int64)
	lens := len(UIDs)
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
				chunkUfosIds = UIDs[(x-1)*params.ChunkSize:]
			} else {
				chunkUfosIds = UIDs[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := s.roomDao.UIDs2roomIDs(ctx, &roomV2.RoomRoomIdByUidMultiReq{Uids: chunkUfosIds}, params)
			if err != nil {
				return err
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ServiceConf.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = ServiceConf.GetRoomID
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
				roomIDs[Index] = itemInt
			}
		}
	}
	return
}

// GetRoomInfo ...
// 获取room信息
func (s *IndexService) GetRoomInfo(ctx context.Context, input *roomV1.RoomGetStatusInfoByUidsReq) (roomResult map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo, err error) {
	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(ServiceConf.GetStatusInfoByUfos)
	params := ServiceConf.ChunkCallInfo{ParamsName: "uids", URLName: ServiceConf.GetStatusInfoByUfos, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
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
			ret, err := s.roomDao.GetRoomInfo(ctx, &roomV1.RoomGetStatusInfoByUidsReq{Uids: chunkUfosIds, FilterOffline: input.FilterOffline, NeedBroadcastType: input.NeedBroadcastType}, params)
			if err != nil {
				return err
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ServiceConf.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = ServiceConf.GetStatusInfoByUfos
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
func (s *IndexService) GetLastLiveTime(ctx context.Context, rolaids []int64) (resp map[string]string, err error) {
	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(ServiceConf.Record)
	params := ServiceConf.ChunkCallInfo{ParamsName: "uids", URLName: ServiceConf.Record, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	resp = make(map[string]string)
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
			ret, err := s.livedataDao.GetLastLiveTime(ctx, &liveDataV1.RecordGetReq{Roomids: chunkRoomIds}, params)
			if err != nil {
				return err
			}
			chunkResult[x-1] = ret
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ServiceConf.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = ServiceConf.Record
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
				resp[k] = item.RecentEndTime
			}
		}
	}
	return
}

// GetRoomNewsInfo ...
// 获取公告信息
func (s *IndexService) GetRoomNewsInfo(ctx context.Context, rolaids *roomExV1.RoomNewsMultiGetReq) (roomNewsResult map[int64]*roomExV1.RoomNewsMultiGetResp_Data, err error) {
	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(ServiceConf.RoomNews)
	params := ServiceConf.ChunkCallInfo{ParamsName: "roomids", URLName: ServiceConf.RoomNews, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
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
			ret, err := s.roomexDao.GetRoomNewsInfo(ctx, &roomExV1.RoomNewsMultiGetReq{RoomIds: chunkRoomIds, IsDecoded: rolaids.IsDecoded}, params)
			if err != nil {
				return err
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ServiceConf.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = ServiceConf.RoomNews
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
func (s *IndexService) GetRoomPendantInfo(ctx context.Context, req *roomV1.RoomPendantGetPendantByIdsReq) (roomNewsResult map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result, err error) {
	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(ServiceConf.RoomPendent)
	params := ServiceConf.ChunkCallInfo{ParamsName: "ids", URLName: ServiceConf.RoomPendent, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
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
			ret, err := s.roomDao.GetRoomPendantInfo(ctx, &roomV1.RoomPendantGetPendantByIdsReq{Ids: chunkRoomIds, Type: req.Type, Position: req.Position}, params)
			if err != nil || ret == nil {
				return err
			}
			chunkResult[x-1] = ret.Data.Result
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ServiceConf.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = ServiceConf.RoomPendent
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
func (s *IndexService) GetPkID(ctx context.Context, req *avV1.PkGetPkIdsByRoomIdsReq) (avResult map[string]int64, err error) {
	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(ServiceConf.GetPkIdsByRoomIds)
	params := ServiceConf.ChunkCallInfo{ParamsName: "ids", URLName: ServiceConf.GetPkIdsByRoomIds, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
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
			ret, err := s.avDao.GetPkID(ctx, &avV1.PkGetPkIdsByRoomIdsReq{RoomIds: chunkRoomIds, Platform: req.Platform}, params)
			if err != nil {
				return err
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ServiceConf.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = ServiceConf.GetPkIdsByRoomIds
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
func (s *IndexService) GetFansMedal(ctx context.Context, req *fansMedalV1.FansMedalTargetsWithMedalReq) (fansResult map[int64]bool, err error) {
	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(ServiceConf.TargetsWithMedal)
	params := ServiceConf.ChunkCallInfo{ParamsName: "ids", URLName: ServiceConf.TargetsWithMedal, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
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
			ret, err := s.fansMedalDao.GetFansMedal(ctx, &fansMedalV1.FansMedalTargetsWithMedalReq{Uid: req.Uid, TargetIds: chunkRoomIds}, params)
			if err != nil {
				return err
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := ServiceConf.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = ServiceConf.TargetsWithMedal
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
func (s *IndexService) GetGrayRule(ctx context.Context, req *userExV1.GrayRuleGetByMarkReq) (extResult *userExV1.GrayRuleGetByMarkResp_Data, err error) {
	extResult = &userExV1.GrayRuleGetByMarkResp_Data{}
	if req == nil {
		return nil, nil
	}
	ret, err := s.userextDao.GetGrayRule(ctx, req)
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
	relationParams := &relationV1.BaseInfoGetGiftInfoReq{}
	giftInfo = make(map[int64]int64)
	ret, err := dao.RelationApi.V1BaseInfo.GetGiftInfo(ctx, relationParams)

	for _, v := range ret.Data {
		giftInfo[v.Mid] = v.Gold
	}
	return
}

// GetUserInfo 获取用户信息
func (s *IndexService) GetUserInfo(ctx context.Context, UIDs []int64) (userResult map[int64]*accountM.Card, err error) {
	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(ServiceConf.AccountGRPC)
	params := ServiceConf.ChunkCallInfo{ParamsName: "uids", URLName: ServiceConf.AccountGRPC, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	userResult = make(map[int64]*accountM.Card)
	lens := len(UIDs)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([]map[int64]*accountM.Card, params.ChunkNum)
	wg, _ := errgroup.WithContext(ctx)
	for i := int64(1); i <= params.ChunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkUfosIds := make([]int64, 20)
			if x == params.ChunkNum {
				chunkUfosIds = UIDs[(x-1)*params.ChunkSize:]
			} else {
				chunkUfosIds = UIDs[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := s.accountDao.GetUserInfoData(ctx, chunkUfosIds)
			if err != nil {
				err = errors.WithMessage(ecode.AccountGRPCError, "GET SEA PATROL FAIL")
				log.Error("Call main.Account.Cards Error.Infos(%+v) error(%+v)", chunkUfosIds, err)
				return nil
			}
			chunkResult[x-1] = ret
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := relationT.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = relationT.AccountGRPC
		erelongInfo.ErrDesc = relationT.GoRoutingErr
		erelongInfo.Code = 1003001
		erelongInfo.RPCTimeout = params.RPCTimeout
		erelongInfo.ErrorPtr = &err
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
			*erelongInfo.ErrorPtr, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		err = errors.WithMessage(ecode.AccountGRPCFrameError, "GET SEA PATROL FAIL")
		return nil, err
	}
	for _, chunkItemList := range chunkResult {
		for _, item := range chunkItemList {
			if item != nil {
				userResult[item.Mid] = item
			}
		}
	}
	return
}

// GetLiveUserExp 获取用户经验信息
func (s *IndexService) GetLiveUserExp(ctx context.Context, UIDs []int64) (userResult map[int64]*liveUserExpM.LevelInfo, err error) {
	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(ServiceConf.LiveUserExpGRPC)
	params := ServiceConf.ChunkCallInfo{ParamsName: "uids", URLName: ServiceConf.LiveUserExpGRPC, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	userResult = make(map[int64]*liveUserExpM.LevelInfo)
	lens := len(UIDs)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([]map[int64]*liveUserExpM.LevelInfo, params.ChunkNum)
	wg, _ := errgroup.WithContext(ctx)
	for i := int64(1); i <= params.ChunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkUfosIds := make([]int64, 20)
			if x == params.ChunkNum {
				chunkUfosIds = UIDs[(x-1)*params.ChunkSize:]
			} else {
				chunkUfosIds = UIDs[(x-1)*params.ChunkSize : x*params.ChunkSize]
			}
			ret, err := s.xuserDao.GetUserExpData(ctx, chunkUfosIds)
			if err != nil {
				err = errors.WithMessage(ecode.AccountGRPCError, "GET SEA PATROL FAIL")
				log.Error("Call main.Account.Cards Error.Infos(%+v) error(%+v)", chunkUfosIds, err)
				return nil
			}
			chunkResult[x-1] = ret
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := relationT.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = relationT.LiveUserExpGRPC
		erelongInfo.ErrDesc = relationT.GoRoutingErr
		erelongInfo.Code = 1003001
		erelongInfo.RPCTimeout = params.RPCTimeout
		erelongInfo.ErrorPtr = &err
		log.Error(erelongInfo.ErrType+"|"+erelongInfo.URLName+"|error:%+v"+"|Code:%d"+"|Msg:%s"+"|RPCTimeout:%d"+"|ChunkSize:%d"+"|ChunkNum:%d"+"|ParamsName:%s",
			*erelongInfo.ErrorPtr, erelongInfo.Code, erelongInfo.Msg, erelongInfo.RPCTimeout, erelongInfo.ChunkSize, erelongInfo.ChunkNum, params.ParamsName)
		err = errors.WithMessage(ecode.AccountGRPCFrameError, "GET SEA PATROL FAIL")
		return nil, err
	}
	for _, chunkItemList := range chunkResult {
		for _, item := range chunkItemList {
			if item != nil {
				userResult[item.Uid] = item
			}
		}
	}
	return
}
