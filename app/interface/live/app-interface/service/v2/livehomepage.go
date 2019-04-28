package v2

import (
	"context"
	"math"
	"strconv"

	"go-common/app/service/live/third_api/bvc"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
	"go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	relationV1 "go-common/app/interface/live/app-interface/service/v1"
	relationT "go-common/app/interface/live/app-interface/service/v1/relation"
	avV1 "go-common/app/service/live/av/api/liverpc/v1"
	relationRpcV1 "go-common/app/service/live/relation/api/liverpc/v1"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	accountM "go-common/app/service/main/account/model"
	actmdl "go-common/app/service/main/account/model"
	account "go-common/app/service/main/account/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
)

// RelationService struct
type RelationService struct {
	conf       *conf.Config
	accountRPC *account.Service3
}

const (
	relationPageSize       = 4
	app536relationPageSize = 2
)

// CheckLiveAnchorParams ... implementation
// [app端关注首页]入参校验
func CheckLiveAnchorParams(ctx context.Context, page int64) (uid int64, relationPage int64, err error) {
	mid := relationT.GetUIDFromHeader(ctx)
	relationPage = page
	err = nil
	if mid == 0 {
		err = errors.WithMessage(ecode.NeedLogIn, "GET SEA PATROL FAIL")
		return
	}
	if page <= 0 {
		err = errors.WithMessage(ecode.ResourceParamErr, "GET SEA PATROL FAIL")
		return
	}
	return mid, relationPage, err
}

// LiveAnchorHomePage ... implementation
// [app端关注首页]正在直播接口
func (s *IndexService) LiveAnchorHomePage(ctx context.Context, relationPage int64, build int64, platform string, quality int64) (Resp []*v2pb.MMyIdol) {
	List := make([]*v2pb.MyIdolItem, 0)
	ExtraInfo := &v2pb.MyIdolExtra{CardType: relationT.App533CardType}
	Resp = make([]*v2pb.MMyIdol, 0)
	s.MakeLiveAnchorDefaultResult(Resp, ExtraInfo)
	uid, relationPage, err := CheckLiveAnchorParams(ctx, relationPage)
	if err != nil && uid != 0 {
		log.Error("[LiveAnchorHomePage]CheckParamsError,uid:%d,relationPage:%d", uid, relationPage)
		return
	}
	wg, _ := errgroup.WithContext(ctx)
	relationInfo, groupList, mapUfos2Rolaids, mapRoomID2UID, AllRoomID, err := relationV1.GetAttentionListAndGroup(ctx)
	if err != nil {
		log.Error("[LiveAnchorHomePage]get_attentionList_rpc_error")
		return
	}
	// 获取全量room信息,不过滤
	roomParams := &roomV1.RoomGetStatusInfoByUidsReq{Uids: groupList["all"], FilterOffline: 0, NeedBroadcastType: 1}
	// room
	roomResp := make(map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo)
	userResp := make(map[int64]*accountM.Card)
	roomCornerResp := make(map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result)
	pkResp := make(map[string]int64)
	wg.Go(func() error {
		roomResp, err = s.GetRoomInfo(ctx, roomParams)
		return err
	})
	if err = wg.Wait(); nil != err {
		log.Error("[LiveAnchorHomePage][first_step]get_room_rpc_error")
		return
	}
	livingUfos := make([]int64, 0)
	livingRolaids := make([]int64, 0)
	livingRoomInfo := GetLivingRooms(roomResp)
	// 没有人直播
	if len(livingRoomInfo) == 0 {
		GetLastLiveAnchorInfo(ctx, roomResp, AllRoomID, mapRoomID2UID, ExtraInfo)
		moduleInfoMap := s.GetAllModuleInfoMapFromCache(ctx)
		for _, m := range moduleInfoMap[_feedType] {
			Resp = append(Resp, &v2pb.MMyIdol{ModuleInfo: m, List: List, ExtraInfo: ExtraInfo})
		}
		return
	}
	wgHasLive := &errgroup.Group{}
	wgHasLive, _ = errgroup.WithContext(ctx)
	for k, v := range livingRoomInfo {
		livingUfos = append(livingUfos, k)
		livingRolaids = append(livingRolaids, v.RoomId)
	}
	// user信息
	wgHasLive.Go(func() error {
		userResp, err = s.GetUserInfo(ctx, livingUfos)
		return err
	})
	// room
	roomPendentParams := &roomV1.RoomPendantGetPendantByIdsReq{Ids: livingRolaids, Type: relationT.PendentMobileBadge, Position: relationT.PendentPosition}
	wgHasLive.Go(func() error {
		roomCornerResp, err = s.GetRoomPendantInfo(ctx, roomPendentParams)
		return err
	})
	// pk_id
	pkParams := &avV1.PkGetPkIdsByRoomIdsReq{RoomIds: livingRolaids, Platform: platform}
	wgHasLive.Go(func() error {
		pkResp, err = s.GetPkID(ctx, pkParams)
		return err
	})
	if err = wgHasLive.Wait(); nil != err {
		log.Error("[LiveAnchorHomePage][second_step]room/main.account/pkID/rpc_error")
		return
	}
	attentionRoomListPlayURLMap := dao.BvcApi.GetPlayUrlMulti(ctx, livingRolaids, 0, quality, build, platform)
	// 下游数据收集完成
	mapSp := make([]int64, 0)
	normalSp := make([]int64, 0)
	mapSp = append(mapSp, groupList["special"]...)
	normalSp = append(normalSp, groupList["normal"]...)
	List = AdaptLivingField(livingRoomInfo, roomCornerResp, userResp, relationInfo, pkResp, attentionRoomListPlayURLMap, mapSp, normalSp, mapUfos2Rolaids)
	ExtraInfo.TotalCount = int64(len(List))

	// 注释原因：app536灰度策略,需要恢复2卡样式(之前在535已全量4卡,但是需求变了)产品：古月
	// https://www.tapd.cn/20082211/prong/stories/view/1120082211001104459

	// userExtParams := &userExV1.GrayRuleGetByMarkReq{Mark: relationT.App536GrayRule}
	// grayRule, err := relationT.GetGrayRule(ctx, userExtParams)
	// var UserExApp536Rule string
	// if err != nil {
	// 	log.Error("[LiveAnchorHomePage]get_GrayRule_rpc_error")
	// 	UserExApp536Rule = ""
	// } else if grayRule != nil {
	// 	UserExApp536Rule = grayRule.Content
	// }
	SliceList, page := s.SliceForHomePage(List, relationPage, uid, platform)
	ExtraInfo.RelationPage = page
	var result v2pb.MMyIdol
	result.ExtraInfo = ExtraInfo
	result.List = SliceList
	moduleInfoMap := s.GetAllModuleInfoMapFromCache(ctx)
	for _, m := range moduleInfoMap[_feedType] {
		Resp = append(Resp, &v2pb.MMyIdol{ModuleInfo: m, List: result.List, ExtraInfo: result.ExtraInfo})
	}

	return
}

// AdaptLivingField ... implementation
// [app端关注首页]填充数据
func AdaptLivingField(roomInfo map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo,
	roomPendentInfo map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result,
	userResult map[int64]*accountM.Card,
	relationInfo map[int64]*relationRpcV1.BaseInfoGetFollowTypeResp_UidInfo,
	pkIDInfo map[string]int64, playURLInfo map[int64]*bvc.PlayUrlItem, specialUID []int64, normalUID []int64,
	mapUfos2Rolaids map[int64]int64) (resp []*v2pb.MyIdolItem) {

	resp = make([]*v2pb.MyIdolItem, 0)
	normalResp := make([]*v2pb.MyIdolItem, 0)
	resp = make([]*v2pb.MyIdolItem, 0)
	if len(specialUID) > 0 {
		item := LiveFireField(roomInfo, roomPendentInfo, userResult, pkIDInfo, playURLInfo, relationInfo, specialUID, mapUfos2Rolaids)
		resp = AppSortRuleOnline(item)
	}
	if len(normalUID) > 0 {
		tempResp := LiveFireField(roomInfo, roomPendentInfo, userResult, pkIDInfo, playURLInfo, relationInfo, normalUID, mapUfos2Rolaids)
		normalResp = AppSortRuleOnline(tempResp)
	}
	if len(normalResp) > 0 {
		resp = append(resp, normalResp...)
	}
	return
}

// SliceForHomePage ... implementation
// app534规则 [app端关注首页]首页slice逻辑,客户端只显示偶数个数,为兼容推荐去重,当个数为3时返回2
//            https://www.tapd.cn/20082211/prong/stories/view/1120082211001067961
//            https://www.tapd.cn/20082211/prong/stories/view/1120082211001085685
//
// app536规则 https://www.tapd.cn/20082211/prong/stories/view/1120082211001104459
func (s *IndexService) SliceForHomePage(input []*v2pb.MyIdolItem, page int64, uid int64, platform string) (resp []*v2pb.MyIdolItem, relationPage int64) {
	resp = make([]*v2pb.MyIdolItem, 0)
	grayRule := s.App536ABTest(uid, platform)
	relationPage = page
	if len(input) <= 0 {
		return
	}
	count := int64(len(input))

	// 536规则
	if grayRule == 1 {

		switch count {
		case 1:
			{
				resp = input[:]
				relationPage = 1
				return
			}
		case 2:
			{
				resp = input[:]
				relationPage = 1
				return
			}
		}
		var pageSize int64
		pageSize = page
		if page < 1 {
			pageSize = 1
		}
		start := (pageSize - 1) * app536relationPageSize
		end := int64(start + app536relationPageSize)
		// 正常slice
		if end <= count {
			resp = input[start:end]
		} else {
			// 回环逻辑,最后一页不足pagesize时返回第一页
			relationPage = 1
			startIndex := 0
			var startCount int64
			if count > app536relationPageSize {
				startCount = app536relationPageSize
			} else {
				startCount = count
			}
			resp = input[startIndex:startCount]
		}
		return
	}

	// 536之前4卡逻辑
	switch count {
	case 1:
		{
			resp = input[:]
			relationPage = 1
			return
		}
	case 2:
		{
			resp = input[:]
			relationPage = 1
			return
		}
	case 3:
		{
			resp = input[0:2]
			relationPage = 1
			return
		}
	}
	var pageSize int64
	pageSize = page
	if page < 1 {
		pageSize = 1
	}
	start := (pageSize - 1) * relationPageSize
	end := int64(start + relationPageSize)
	if end <= count {
		resp = input[start:end]
	} else {
		relationPage = 1
		startIndex := 0
		var startCount int64
		if count > relationPageSize {
			startCount = relationPageSize
		} else {
			startCount = count
		}
		resp = input[startIndex:startCount]
	}
	return
}

// LiveFireField ... implementation
// [app端关注首页]填充数据
func LiveFireField(roomInfo map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo,
	roomPendentInfo map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result,
	userResult map[int64]*accountM.Card,
	pkIDInfo map[string]int64, playURLInfo map[int64]*bvc.PlayUrlItem,
	relationInfo map[int64]*relationRpcV1.BaseInfoGetFollowTypeResp_UidInfo,
	ufos []int64, mapUfos2Rolaids map[int64]int64) (resp []*v2pb.MyIdolItem) {
	for _, v := range ufos {
		item := v2pb.MyIdolItem{}
		roomID, roomIDExist := mapUfos2Rolaids[v]
		if !roomIDExist {
			continue
		}
		roomItem := roomInfo[v]
		roomPendentItem := roomPendentInfo[roomID]
		userItem := userResult[v]
		relationItem := relationInfo[v]
		pkItem := pkIDInfo[strconv.Itoa(int(roomID))]
		playURLItem := playURLInfo[roomID]
		if roomItem == nil || userItem == nil || relationItem == nil {
			continue
		}
		PlayURL := ""
		PlayURL265 := ""
		PlayURLAcc := make([]int64, 0)
		PlayURLCur := 0
		PendentRu := ""
		PendentRuColor := ""
		PendentRuPic := ""
		if playURLItem != nil {
			PlayURL = playURLItem.Url["h264"]
			PlayURL265 = playURLItem.Url["h265"]
			PlayURLAcc = playURLItem.AcceptQuality
			PlayURLCur = int(playURLItem.CurrentQuality)
		}
		if roomPendentItem != nil {
			PendentRu = roomPendentItem.Value
			PendentRuColor = roomPendentItem.BgColor
			PendentRuPic = roomPendentItem.BgPic
		}
		item.Roomid = roomItem.RoomId
		item.Uid = roomItem.Uid
		item.Uname = userItem.Name
		item.Face = userItem.Face
		item.Title = roomItem.Title
		item.LiveTagName = roomItem.AreaV2Name
		item.LiveTime = roomItem.LiveTime
		item.Online = roomItem.Online
		item.PlayUrl = PlayURL
		item.PlayUrlH265 = PlayURL265
		item.AcceptQuality = PlayURLAcc
		item.CurrentQuality = int64(PlayURLCur)
		item.PkId = pkItem
		item.Area = roomItem.Area
		item.AreaName = roomItem.AreaName
		item.AreaV2Id = roomItem.AreaV2Id
		item.AreaV2Name = roomItem.AreaV2Name
		item.AreaV2ParentId = roomItem.AreaV2ParentId
		item.AreaV2ParentName = roomItem.AreaV2ParentName
		item.BroadcastType = roomItem.BroadcastType
		item.Link = relationT.LiveDomain + strconv.Itoa(int(roomID)) + relationT.BoastURL + strconv.Itoa(int(item.BroadcastType))
		item.OfficialVerify = int64(RoleMap(userItem.Official.Role))
		item.SpecialAttention = relationItem.Special
		item.PendentRu = PendentRu
		item.PendentRuColor = PendentRuColor
		item.PendentRuPic = PendentRuPic
		if len(roomItem.CoverFromUser) == 0 {
			item.Cover = roomItem.Keyframe
		} else {
			item.Cover = roomItem.CoverFromUser
		}
		resp = append(resp, &item)
	}
	return
}

// MakeLiveAnchorDefaultResult ...
// 正在直播默认返回
func (s *IndexService) MakeLiveAnchorDefaultResult(Resp []*v2pb.MMyIdol, ExtraInfo *v2pb.MyIdolExtra) {
	if ExtraInfo != nil {
		ExtraInfo.TotalCount = 0
		ExtraInfo.TagsDesc = ""
		ExtraInfo.UnameDesc = ""
		ExtraInfo.TimeDesc = ""
		// [历史原因]cardType只能为1,否则客户端报错,见 https://www.tapd.cn/20082211/prong/stories/view/1120082211001086997
		ExtraInfo.CardType = 1
		ExtraInfo.RelationPage = 1
	}
	moduleInfoMap := s.GetAllModuleInfoMapFromCache(context.TODO())
	for _, m := range moduleInfoMap[_feedType] {
		Resp = append(Resp, &v2pb.MMyIdol{ModuleInfo: m, List: []*v2pb.MyIdolItem{}, ExtraInfo: ExtraInfo})
	}
}

// GetLivingRooms ... implementation
// [app端关注首页]获取正在直播房间
func GetLivingRooms(roomResult map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo) (liveRoom map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo) {
	liveRoom = make(map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo)
	if len(roomResult) == 0 {
		return
	}
	for k, v := range roomResult {
		if v.LiveStatus == relationV1.RoomStatusLive {
			liveRoom[k] = v
		}
	}
	return
}

// CheckLiveAnchorParams implementation
// 入参校验
func (s *IndexService) CheckLiveAnchorParams(ctx context.Context, req *v2pb.GetAllListReq) (uid int64, relationPage int64, err error) {
	if req == nil {
		err = ecode.LiveAnchorReqV2ParamsNil
		return
	}
	uid = relationT.GetUIDFromHeader(ctx)
	if uid == 0 {
		err = errors.WithMessage(ecode.NeedLogIn, "GET SEA PATROL FAIL")
		return
	}
	if req.RelationPage <= 0 {
		log.Error("CallRelationLiveAnchorV2ParamsCheckError|relationPage:%d", req.RelationPage)
		err = errors.WithMessage(ecode.LiveAnchorReqV2ParamsError, "GET SEA PATROL FAIL")
		return
	}
	return
}

// GetLastLiveAnchorInfo ... implementation
// [app端关注首页]获取最新一次直播信息
func GetLastLiveAnchorInfo(ctx context.Context, roomResult map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo,
	RoomIDs []int64, RoomID2UID map[int64]int64, ExtraInfo *v2pb.MyIdolExtra) (uid int64, relationPage int64, err error) {
	if len(roomResult) == 0 || len(RoomIDs) == 0 || len(RoomID2UID) == 0 {
		return
	}
	lastLiveTime, _ := relationT.GetLastLiveTime(ctx, RoomIDs)
	_, _, sorted := relationV1.GetLastAnchorLiveTime(lastLiveTime)
	var firstRoom int64
	var firstValue int64
	if sorted.Len() > 0 {
		for _, v := range sorted {
			firstRoom = int64(v.Key)
			firstValue = int64(v.Value)
			break
		}
		firstUID := int64(RoomID2UID[firstRoom])
		tempTime := make(map[int64]int64)
		if firstValue > 0 {
			tempTime[firstUID] = firstValue
			if roomItem, exist := roomResult[firstUID]; exist {
				ExtraInfo.UnameDesc = roomItem.Uname
				ExtraInfo.TagsDesc = roomItem.AreaV2Name
				liveDesc, _ := relationV1.TimeLineRule(tempTime, nil)
				ExtraInfo.TimeDesc = liveDesc[firstUID]
			}
		}
	}
	return
}

// GetUserInfoData ...
// 调用account grpc接口cards获取用户信息
func (s *RelationService) GetUserInfoData(ctx context.Context, UIDs []int64) (userResult map[int64]*accountM.Card, err error) {

	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(relationT.AccountGRPC)
	params := relationT.ChunkCallInfo{ParamsName: "ufos", URLName: relationT.AccountGRPC, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
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
			ret, err := s.accountRPC.Cards3(ctx, &actmdl.ArgMids{Mids: chunkUfosIds})
			if err != nil {
				err = errors.WithMessage(ecode.AccountGRPCError, "GET SEA PATROL FAIL")
				log.Error("Call main.Account.Cards Error.Infos(%+v) error(%+v)", chunkUfosIds, err)
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
	// 整理数据
	for _, chunkItemList := range chunkResult {
		for _, item := range chunkItemList {
			if item != nil {
				userResult[item.Mid] = item
			}
		}
	}
	return
}

// App536ABTest ... hard code配置
// ABTest
func (s *IndexService) App536ABTest(mid int64, platform string) (grayType int64) {
	// 因为ipad屏幕尺寸较大,展示2卡会影响体验,故如果是ipad客户端则使用4卡样式,产品：tianyumo
	if platform == "ipad" {
		return 0
	}
	mUID := mid % 100
	if mUID >= 0 && mUID <= 89 {
		// 4卡
		return 0
	}
	return 1
}

// 后台配置
// // App536ABTest ...
// // ABTest
// func (s *IndexService) App536ABTest(content string, mid int64) (grayType int64) {
// 	if len(content) == 0 {
// 		grayType = 0
// 		return
// 	}
// 	resultMap := make(map[string]int64)
// 	resultMap["app536_4card_type"] = 0
// 	resultMap["app536_2card_type"] = 1
// 	typeMap := make([]string, 0)
// 	mr := &[]GrayRule{}
// 	if err := json.Unmarshal([]byte(content), mr); err != nil {
// 		grayType = 0
// 		return
// 	}
// 	ruleArr := *mr
// 	scoreMap := make(map[string]int)
//
// 	for _, v := range ruleArr {
// 		scoreMap[v.Mark] = int(RParseInt(v.Value, 100))
// 	}
// 	sortedScore := SortMapByValue(scoreMap)
// 	scoreEnd := make([]int, 0)
// 	for _, v := range sortedScore {
// 		scoreEnd = append(scoreEnd, v.Value)
// 		typeMap = append(typeMap, v.Key)
// 	}
// 	score1 := scoreEnd[0]
// 	score2 := 100
// 	section1 := make(map[int]bool)
// 	section2 := make(map[int]bool)
// 	for section1Loop := 0; section1Loop < score1; section1Loop++ {
// 		section1[section1Loop] = true
// 	}
// 	for sectionLoop2 := score1; sectionLoop2 < score2; sectionLoop2++ {
// 		section2[sectionLoop2] = true
// 	}
// 	result := int(mid % 100)
// 	if scoreEnd[0] != 0 {
// 		if _, exist := section1[result]; exist {
// 			grayType = resultMap[typeMap[0]]
// 			return
// 		}
// 	}
// 	if scoreEnd[1] != 0 {
// 		if _, exist := section2[result]; exist {
// 			grayType = resultMap[typeMap[1]]
// 			return
// 		}
// 	}
// 	grayType = 0
// 	return
// }
