package v1

import (
	"context"
	"math"
	"strconv"
	"time"

	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"

	v1pb "go-common/app/interface/live/app-interface/api/http/v1"
	"go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	relationT "go-common/app/interface/live/app-interface/service/v1/relation"
	avV1 "go-common/app/service/live/av/api/liverpc/v1"
	relationV1 "go-common/app/service/live/relation/api/liverpc/v1"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	roomExV1 "go-common/app/service/live/room_ex/api/liverpc/v1"
	playurlbvc "go-common/app/service/live/third_api/bvc"
	userExV1 "go-common/app/service/live/userext/api/liverpc/v1"
	accountM "go-common/app/service/main/account/model"
	actmdl "go-common/app/service/main/account/model"
	account "go-common/app/service/main/account/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/rpc/liverpc"
	liveConText "go-common/library/net/rpc/liverpc/context"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
)

// RelationService struct
type RelationService struct {
	conf       *conf.Config
	accountRPC *account.Service3
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

// NewRelationService init
func NewRelationService(c *conf.Config) (s *RelationService) {
	s = &RelationService{
		conf:       c,
		accountRPC: account.New3(nil),
	}
	return s
}

const (
	// RoomStatusLive ...
	RoomStatusLive = 1
	// MobileIndexBadgeColorDefault ...
	MobileIndexBadgeColorDefault = "#FB9E60"
)

// UnliveAnchor ... implementation
// 直播二级页暂未开播接口
func (s *RelationService) UnliveAnchor(ctx context.Context, req *v1pb.UnLiveAnchorReq) (resp *v1pb.UnLiveAnchorResp, err error) {
	resp = &v1pb.UnLiveAnchorResp{}
	config := conf.GetDummyUidConf()
	if config == relationT.DummyUIDEnable {
		dummyHeader := &liverpc.Header{Uid: relationT.RParseInt(req.Buyaofangqizhiliao, relationT.SelfUID)}
		ctx = liveConText.WithHeader(ctx, dummyHeader)
	}
	MakeUnLiveDefaultResult(resp)
	uid := relationT.GetUIDFromHeader(ctx)
	if uid <= 0 && config == 0 {
		return
	}
	wg, _ := errgroup.WithContext(ctx)
	pass, page, pageSize, uid, err := CheckUnLiveAnchorParams(ctx, req)
	if !pass {
		log.Error("[UnLiveAnchor]CheckParamsError,page:%d,pageSize:%d,uid:%d", page, pageSize, uid)
		return
	}
	relationInfo, groupList, mapUfos2Rolaids, mapRolaids2Ufos, setRolaids, err := GetAttentionListAndGroup(ctx)
	if err != nil {
		log.Error("[LiveAnchor]get_attentionList_rpc_error")
		return
	}

	// 获取有效(曾经直播过)主播,剪枝roomIDs
	lastLiveTime, _ := relationT.GetLastLiveTime(ctx, setRolaids)
	alienableRolaids, liberateInfo, sorted := FilterEverLived(lastLiveTime)
	alienableUfos := GetUID(mapRolaids2Ufos, alienableRolaids)

	roomExReq := &roomExV1.RoomNewsMultiGetReq{RoomIds: alienableRolaids, IsDecoded: 1}
	roomParams := &roomV1.RoomGetStatusInfoByUidsReq{Uids: groupList["all"], FilterOffline: 0}

	userInfo := make(map[int64]*accountM.Card)
	roomResp := make(map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo)
	roomExResp := make(map[int64]*roomExV1.RoomNewsMultiGetResp_Data)
	userfcResp := make(map[int64]*relationV1.FeedGetUserFcBatchResp_RelationList)

	// room
	wg.Go(func() error {
		roomResp, err = relationT.GetRoomInfo(ctx, roomParams)
		return err
	})

	// user信息
	wg.Go(func() error {
		userInfo, err = s.GetUserInfoData(ctx, alienableUfos)
		return err
	})

	// roomEx
	wg.Go(func() error {
		roomExResp, err = relationT.GetRoomNewsInfo(ctx, roomExReq)
		return err
	})

	// fansNum
	wg.Go(func() error {
		userfcResp, err = GetUserFc(ctx, alienableUfos)
		return err
	})

	waitErr := wg.Wait()
	if waitErr != nil {
		log.Error("[UnLiveAnchor][step2] rpc error: %s", waitErr)
		return
	}

	mapSp := make(map[int64]bool)
	normalSp := make(map[int64]bool)
	for _, v := range groupList["special"] {
		mapSp[v] = true
	}
	for _, v := range groupList["normal"] {
		normalSp[v] = true
	}
	specialRoomed, normalRoomed := s.GroupByRule(ctx, mapSp, normalSp, liberateInfo, sorted, mapRolaids2Ufos)
	specialUID := GetUID(mapRolaids2Ufos, specialRoomed)
	normalUID := GetUID(mapRolaids2Ufos, normalRoomed)
	liveDesc, newsDesc := CalcTimeLine(liberateInfo, roomExResp)
	LiveCount := CountLiveRooms(roomResp)
	resp.Rooms = AdaptField(roomResp, userfcResp, userInfo, roomExResp, relationInfo, specialUID, normalUID, mapUfos2Rolaids, liveDesc, newsDesc)
	resp.TotalCount = int64(len(resp.Rooms))
	resp.NoRoomCount = int64(len(setRolaids) - int(resp.TotalCount) - int(LiveCount))
	resp.Rooms = UnLiveAnchorSlice(resp.Rooms, page, pageSize)
	if (page * pageSize) >= resp.TotalCount {
		resp.HasMore = 0
	} else {
		resp.HasMore = 1
	}
	return
}

// CountLiveRooms 计算正在直播数目
func CountLiveRooms(input map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo) (count int) {
	if len(input) <= 0 {
		count = 0
		return
	}
	for _, v := range input {
		if v.LiveStatus == RoomStatusLive {
			count++
		}
	}
	return
}

// CheckUnLiveAnchorParams implementation
// 入参校验
func CheckUnLiveAnchorParams(ctx context.Context, req *v1pb.UnLiveAnchorReq) (pass bool, page int64, pageSize int64, uid int64, err error) {
	if req == nil {
		pass = false
		return
	}
	config := conf.GetDummyUidConf()
	uid = relationT.GetUIDFromHeader(ctx)
	if uid == 0 && config == 0 {
		err = errors.WithMessage(ecode.NeedLogIn, "GET SEA PATROL FAIL")
		pass = false
		return
	}
	page = req.Page
	pageSize = req.Pagesize
	if page <= 0 || pageSize <= 0 {
		pass = false
		log.Error("CallRelationUnLiveAnchorParamsCheckError|page:%d,pageSize:%d", page, pageSize)
		err = errors.WithMessage(ecode.UnliveAnchorReqParamsError, "GET SEA PATROL FAIL")
		return
	}
	pass = true
	return
}

// CheckLiveAnchorParams implementation
// 入参校验
func CheckLiveAnchorParams(ctx context.Context, req *v1pb.LiveAnchorReq) (sortRule int64, filterRule int64, uid int64, err error) {
	if req == nil {
		err = ecode.LiveAnchorReqParamsNil
		return
	}
	sortRule = req.SortRule
	filterRule = req.FilterRule
	uid = relationT.GetUIDFromHeader(ctx)
	config := conf.GetDummyUidConf()
	if uid == 0 && config == 0 {
		err = errors.WithMessage(ecode.NeedLogIn, "GET SEA PATROL FAIL")
		return
	}
	if sortRule < 0 || filterRule < 0 {
		log.Error("CallRelationLiveAnchorParamsCheckError|page:%d,pageSize:%d", sortRule, filterRule)
		err = errors.WithMessage(ecode.LiveAnchorReqParamsError, "GET SEA PATROL FAIL")
		return
	}
	return
}

// UnLiveAnchorSlice implementation
// 分页逻辑
func UnLiveAnchorSlice(req []*v1pb.UnLiveAnchorResp_Rooms, page int64, pageSize int64) (resp []*v1pb.UnLiveAnchorResp_Rooms) {
	resp = make([]*v1pb.UnLiveAnchorResp_Rooms, 0)
	start := (page - 1) * pageSize
	end := start + pageSize
	length := int64(len(req))
	if start >= length {
		return
	}
	if end >= length {
		resp = req[start:]
	} else {
		resp = req[start:end]
	}
	return
}

// MakeUnLiveDefaultResult implementation
// 缺省返回
func MakeUnLiveDefaultResult(resp *v1pb.UnLiveAnchorResp) {
	if resp != nil {
		resp.HasMore = 0
		resp.NoRoomCount = 0
		resp.TotalCount = 0
		resp.Rooms = make([]*v1pb.UnLiveAnchorResp_Rooms, 0)
	}
}

// GroupByRule implementation
// 按照规则排序,组间按照特别关注优先,组内按照上次关播时间倒序
func (s *RelationService) GroupByRule(ctx context.Context, special map[int64]bool, normal map[int64]bool,
	liberate map[int64]int64, sorted relationT.PairList, mapRolaids2Ufos map[int64]int64) (specialRoomed []int64, normalRoomed []int64) {
	specialRoomed = make([]int64, 0)
	normalRoomed = make([]int64, 0)
	if len(liberate) == 0 || liberate == nil {
		return
	}
	for _, v := range sorted {
		if _, ok := special[mapRolaids2Ufos[v.Key]]; ok {
			specialRoomed = append(specialRoomed, v.Key)
		} else {
			normalRoomed = append(normalRoomed, v.Key)
		}
	}
	return
}

// AdaptField implementation
// 填充逻辑
func AdaptField(roomInfo map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo,
	fansInfo map[int64]*relationV1.FeedGetUserFcBatchResp_RelationList,
	userResult map[int64]*accountM.Card,
	roomedInfo map[int64]*roomExV1.RoomNewsMultiGetResp_Data,
	relationInfo map[int64]*relationV1.BaseInfoGetFollowTypeResp_UidInfo,
	specialUID []int64, normalUID []int64,
	mapUfos2Rolaids map[int64]int64, liveDesc map[int64]string, newsDesc map[int64]string) (resp []*v1pb.UnLiveAnchorResp_Rooms) {
	var item []*v1pb.UnLiveAnchorResp_Rooms
	resp = make([]*v1pb.UnLiveAnchorResp_Rooms, 0)
	if len(specialUID) > 0 {
		item = FireField(roomInfo, fansInfo, userResult, roomedInfo, relationInfo, specialUID, mapUfos2Rolaids, liveDesc, newsDesc)
		resp = append(resp, item...)
	}
	if len(normalUID) > 0 {
		item = FireField(roomInfo, fansInfo, userResult, roomedInfo, relationInfo, normalUID, mapUfos2Rolaids, liveDesc, newsDesc)
		resp = append(resp, item...)
	}
	return
}

// CalcTimeLine ...
// 计算时间规则
func CalcTimeLine(liberateInfo map[int64]int64,
	roomNewsInfo map[int64]*roomExV1.RoomNewsMultiGetResp_Data) (liveDesc map[int64]string, newsDesc map[int64]string) {
	liveDesc, newsDesc = TimeLineRule(liberateInfo, roomNewsInfo)
	return
}

// TimeLineRule ...
// 计算时间规则
func TimeLineRule(liberateInfo map[int64]int64, roomNewsInfo map[int64]*roomExV1.RoomNewsMultiGetResp_Data) (liveDesc map[int64]string, newsDesc map[int64]string) {
	liveDesc = make(map[int64]string)
	newsDesc = make(map[int64]string)
	if len(liberateInfo) <= 0 {
		return
	}
	for livedRoomed, lastLiveTime := range liberateInfo {
		now := time.Now()
		currentYear, currentMonth, currentDay := now.Date()
		currentLocation := now.Location()
		firstOfMonth := time.Date(currentYear, 1, 1, 0, 0, 0, 0, currentLocation)
		thisYearUnixTimeStamp := firstOfMonth.Unix()
		todayUnixTimeStamp := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation).Unix()
		today24 := math.Abs(float64(todayUnixTimeStamp - lastLiveTime))
		liveTime := math.Abs(float64(now.Unix() - lastLiveTime))
		if lastLiveTime == 0 {
			liveDesc[livedRoomed] = "上次"
		}
		if liveTime < 60 {
			liveDesc[livedRoomed] = "刚刚"
		} else if liveTime >= 60 && liveTime < 3600 {
			text := int(math.Floor(liveTime / 60))
			liveDesc[livedRoomed] = strconv.Itoa(text) + "分钟前"
		} else if liveTime >= 3600 && liveTime < 86400 {
			text := int(math.Floor(liveTime / 3600))
			liveDesc[livedRoomed] = strconv.Itoa(text) + "小时前"
		} else if liveTime >= 86400 && today24 <= 86400 {
			liveDesc[livedRoomed] = "昨天"
		} else if liveTime >= 86400 && lastLiveTime >= thisYearUnixTimeStamp {
			tm := time.Unix(lastLiveTime, 0)
			text := tm.Format("1-2")
			liveDesc[livedRoomed] = text
		} else {
			if lastLiveTime < thisYearUnixTimeStamp && liveTime >= 86400 {
				tm := time.Unix(lastLiveTime, 0)
				text := tm.Format("2006-1-2")
				liveDesc[livedRoomed] = text
			} else {
				tm := time.Unix(lastLiveTime, 0)
				text := tm.Format("2006-1-2")
				liveDesc[livedRoomed] = text
			}
		}
	}

	if len(roomNewsInfo) <= 0 {
		return
	}
	for livedRoomed, lastNewsTime := range roomNewsInfo {
		lastLiveTimeStr := lastNewsTime.Ctime
		now := time.Now()
		timeFmt, _ := time.ParseInLocation("2006-01-02 15:04:05", lastLiveTimeStr, time.Local)
		lastLiveTime := timeFmt.Unix()
		currentYear, currentMonth, currentDay := now.Date()
		currentLocation := now.Location()
		todayUnixTimeStamp := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation).Unix()
		today24 := math.Abs(float64(todayUnixTimeStamp - lastLiveTime))
		liveTime := math.Abs(float64(now.Unix() - lastLiveTime))
		if lastLiveTime == 0 {
			newsDesc[livedRoomed] = ""
		}
		if liveTime < 60 {
			newsDesc[livedRoomed] = "刚刚"
		} else if liveTime >= 60 && liveTime < 3600 {
			text := int(math.Floor(liveTime / 60))
			newsDesc[livedRoomed] = strconv.Itoa(text) + "分钟前"
		} else if liveTime >= 3600 && liveTime < 86400 {
			text := int(math.Floor(liveTime / 3600))
			newsDesc[livedRoomed] = strconv.Itoa(text) + "小时前"
		} else if liveTime >= 86400 && today24 <= 86400 {
			newsDesc[livedRoomed] = "昨天"
		}
	}
	return
}

// FireField ...
// 适配返回值
func FireField(roomInfo map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo,
	fansInfo map[int64]*relationV1.FeedGetUserFcBatchResp_RelationList,
	userResult map[int64]*accountM.Card,
	roomedInfo map[int64]*roomExV1.RoomNewsMultiGetResp_Data,
	relationInfo map[int64]*relationV1.BaseInfoGetFollowTypeResp_UidInfo,
	ufos []int64,
	mapUfos2Rolaids map[int64]int64, liveDesc map[int64]string, newsDesc map[int64]string) (resp []*v1pb.UnLiveAnchorResp_Rooms) {
	for _, v := range ufos {
		item := v1pb.UnLiveAnchorResp_Rooms{}
		roomID, roomIDExist := mapUfos2Rolaids[v]
		if !roomIDExist {
			continue
		}
		roomItem := roomInfo[v]
		userItem := userResult[v]
		fansItem := fansInfo[v]
		relationItem := relationInfo[v]
		roomedItem := roomedInfo[roomID]
		roomNewsDesc := newsDesc[roomID]
		liveDescItem := liveDesc[roomID]
		roomNewsContent := ""
		roomNewsDescText := ""
		if roomItem == nil || userItem == nil || relationItem == nil {
			continue
		}
		if roomItem.LiveStatus == RoomStatusLive {
			continue
		}
		if roomedItem != nil {
			roomNewsContent = roomedItem.NewsContent
			roomNewsDescText = roomNewsDesc
		}
		item.Roomid = roomItem.RoomId
		item.Uid = roomItem.Uid
		item.Uname = userItem.Name
		item.Face = userItem.Face
		item.LiveStatus = roomItem.LiveStatus
		item.Area = roomItem.Area
		item.AreaName = roomItem.AreaName
		item.AreaV2Id = roomItem.AreaV2Id
		item.AreaV2Name = roomItem.AreaV2Name
		item.AreaV2ParentId = roomItem.AreaV2ParentId
		item.AreaV2ParentName = roomItem.AreaV2ParentName
		item.BroadcastType = roomItem.BroadcastType
		item.Link = relationT.LiveDomain + strconv.Itoa(int(roomID)) + relationT.BoastURL + strconv.Itoa(int(item.BroadcastType))
		item.OfficialVerify = int64(relationT.RoleMap(userItem.Official.Role))
		item.Attentions = fansItem.Fc
		item.SpecialAttention = relationItem.Special
		item.AnnouncementContent = roomNewsContent
		item.AnnouncementTime = roomNewsDescText
		item.LiveDesc = liveDescItem

		resp = append(resp, &item)
	}
	return
}

// LiveFireField ...
// 适配返回值
func LiveFireField(roomInfo map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo,
	roomPendentInfo map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result,
	userResult map[int64]*accountM.Card,
	pkIDInfo map[string]int64, playURLInfo map[int64]*playurlbvc.PlayUrlItem,
	relationInfo map[int64]*relationV1.BaseInfoGetFollowTypeResp_UidInfo,
	ufos []int64, mapUfos2Rolaids map[int64]int64) (resp []*v1pb.LiveAnchorResp_Rooms) {

	for _, v := range ufos {
		item := v1pb.LiveAnchorResp_Rooms{}
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
		if PendentRuColor == "" {
			PendentRuColor = MobileIndexBadgeColorDefault

		}
		item.Roomid = roomItem.RoomId
		item.Uid = roomItem.Uid
		item.Uname = userItem.Name
		item.Face = userItem.Face
		item.Title = roomItem.Title
		item.LiveTagName = roomItem.AreaV2Name
		item.LiveTime = roomItem.LiveTime
		item.Online = roomItem.Online
		item.Playurl = PlayURL
		item.AcceptQuality = PlayURLAcc
		item.CurrentQuality = int64(PlayURLCur)
		item.PkId = pkItem
		item.Area = roomItem.Area
		item.AreaName = roomItem.AreaName
		item.AreaV2Id = roomItem.AreaV2Id
		item.PlayUrlH265 = PlayURL265
		item.AreaV2Name = roomItem.AreaV2Name
		item.AreaV2ParentId = roomItem.AreaV2ParentId
		item.AreaV2ParentName = roomItem.AreaV2ParentName
		item.BroadcastType = roomItem.BroadcastType
		item.Link = relationT.LiveDomain + strconv.Itoa(int(roomID)) + relationT.BoastURL + strconv.Itoa(int(item.BroadcastType))
		item.OfficialVerify = int64(relationT.RoleMap(userItem.Official.Role))
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

// GroupUfos ...
// 按照关注类型分组
func GroupUfos(input map[int64]*relationV1.BaseInfoGetFollowTypeResp_UidInfo) (resp map[string][]int64, err error) {
	if input == nil {
		return nil, nil
	}
	resp = make(map[string][]int64)
	for k, v := range input {
		if v.Special == 0 {
			resp["normal"] = append(resp["normal"], k)
		} else {
			resp["special"] = append(resp["special"], k)
		}
		resp["all"] = append(resp["all"], k)
	}
	return resp, nil
}

// GetUID ...
// 获取uid
func GetUID(idsMap map[int64]int64, input []int64) (resp []int64) {
	if idsMap == nil || input == nil {
		return nil
	}
	for _, v := range input {
		resp = append(resp, idsMap[int64(v)])
	}
	return resp
}

// FilterEverLived ...
// 过滤未开播
func FilterEverLived(lastLiveTime map[string]string) (rolaids []int64, lifetime map[int64]int64, sorted relationT.PairList) {
	rolaids = make([]int64, 0)
	lifetime = make(map[int64]int64)
	for roomed, v := range lastLiveTime {
		timeFmt, _ := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
		if !timeFmt.IsZero() {
			if mid, err := strconv.ParseInt(roomed, 10, 64); err == nil {
				lifetime[mid] = timeFmt.Unix()
				rolaids = append(rolaids, mid)
			}
		}
	}
	sorted = make([]relationT.Pair, 0)
	sorted = relationT.SortMap(lifetime)
	return rolaids, lifetime, sorted
}

// GetLastAnchorLiveTime ...
// 获取上一个主播信息
func GetLastAnchorLiveTime(lastLiveTime map[string]string) (rolaids []int64, lifetime map[int64]int64, sorted relationT.PairList) {
	rolaids = make([]int64, 0)
	lifetime = make(map[int64]int64)
	for roomed, v := range lastLiveTime {
		timeFmt, _ := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
		if mid, err := strconv.ParseInt(roomed, 10, 64); err == nil {
			lifetime[mid] = timeFmt.Unix()
			rolaids = append(rolaids, mid)
		}
	}
	sorted = make([]relationT.Pair, 0)
	sorted = relationT.SortMap(lifetime)
	return rolaids, lifetime, sorted
}

// MakeLiveAnchorDefaultResult ...
// 正在直播默认返回
func MakeLiveAnchorDefaultResult(resp *v1pb.LiveAnchorResp) {
	if resp != nil {
		resp.TotalCount = 0
		// [历史原因]cardType只能为1,否则客户端报错,见 https://www.tapd.cn/20082211/prong/stories/view/1120082211001086997
		resp.CardType = relationT.App533CardType
		resp.BigCardType = 0
		resp.Rooms = make([]*v1pb.LiveAnchorResp_Rooms, 0)
	}
}

// GetAttentionListAndGroup ...
// 关注分组
func GetAttentionListAndGroup(ctx context.Context) (relationInfo map[int64]*relationV1.BaseInfoGetFollowTypeResp_UidInfo, groupList map[string][]int64,
	mapUfos2Rolaids map[int64]int64, mapRolaids2Ufos map[int64]int64, setRolaids []int64, attentionErr error) {
	relationTimeout := conf.GetTimeout("relation", 200)
	attentionErr = nil
	attentionData, attentionErr := dao.RelationApi.V1BaseInfo.GetFollowType(
		rpcCtx.WithTimeout(ctx, time.Duration(relationTimeout)*time.Millisecond),
		&relationV1.BaseInfoGetFollowTypeReq{})
	if attentionErr != nil || attentionData == nil {
		attentionErr = ecode.AttentionListRPCError
		return
	}
	relationInfo = attentionData.Data
	groupList, _ = GroupUfos(attentionData.Data)
	// 转换ids
	mapUfos2Rolaids, err := relationT.UIDs2roomIDs(ctx, groupList["all"])
	if err != nil {
		attentionErr = ecode.RoomGetRoomIDCodeRPCError
		return
	}
	mapRolaids2Ufos, setRolaids = TransRoomedUUID(mapUfos2Rolaids)
	return
}

// TransRoomedUUID ...
// 转换ids
func TransRoomedUUID(mapUfos2Rolaids map[int64]int64) (mapRolaids2Ufos map[int64]int64, setRolaids []int64) {
	mapRolaids2Ufos = make(map[int64]int64)
	for k, v := range mapUfos2Rolaids {
		mapRolaids2Ufos[v] = k
	}
	setRolaids = make([]int64, 0)
	for _, v := range mapUfos2Rolaids {
		setRolaids = append(setRolaids, v)
	}
	return
}

// AdaptLivingField ...
// 填充逻辑
func AdaptLivingField(roomInfo map[int64]*roomV1.RoomGetStatusInfoByUidsResp_RoomInfo,
	roomPendentInfo map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result,
	userResult map[int64]*accountM.Card,
	relationInfo map[int64]*relationV1.BaseInfoGetFollowTypeResp_UidInfo,
	pkIDInfo map[string]int64, playURLInfo map[int64]*playurlbvc.PlayUrlItem, specialUID []int64, normalUID []int64,
	mapUfos2Rolaids map[int64]int64) (resp []*v1pb.LiveAnchorResp_Rooms) {

	var item []*v1pb.LiveAnchorResp_Rooms
	resp = make([]*v1pb.LiveAnchorResp_Rooms, 0)
	normalResp := make([]*v1pb.LiveAnchorResp_Rooms, 0)
	if len(specialUID) > 0 {
		item = LiveFireField(roomInfo, roomPendentInfo, userResult, pkIDInfo, playURLInfo, relationInfo, specialUID, mapUfos2Rolaids)
		tempResp := &v1pb.LiveAnchorResp{}
		tempResp.Rooms = item
		resp = relationT.AppSortRuleOnline(tempResp)
	}
	if len(normalUID) > 0 {
		item = LiveFireField(roomInfo, roomPendentInfo, userResult, pkIDInfo, playURLInfo, relationInfo, normalUID, mapUfos2Rolaids)
		tempResp := &v1pb.LiveAnchorResp{}
		tempResp.Rooms = item
		normalResp = relationT.AppSortRuleOnline(tempResp)
	}
	if len(normalResp) > 0 {
		resp = append(resp, normalResp...)
	}
	return
}

// LiveAnchor implementation
// [app端关注二级页][全量]正在直播接口
func (s *RelationService) LiveAnchor(ctx context.Context, req *v1pb.LiveAnchorReq) (resp *v1pb.LiveAnchorResp, err error) {
	resp = &v1pb.LiveAnchorResp{}
	MakeLiveAnchorDefaultResult(resp)
	sortRule, filterRule, uid, err := CheckLiveAnchorParams(ctx, req)
	wg, _ := errgroup.WithContext(ctx)
	if err != nil {
		log.Error("[LiveAnchor]CheckParamsError,page:%d,pageSize:%d,uid:%d", sortRule, filterRule, uid)
		return
	}
	relationInfo, groupList, mapUfos2Rolaids, _, _, err := GetAttentionListAndGroup(ctx)
	if err != nil {
		log.Error("[LiveAnchor]get_attentionList_rpc_error")
		return
	}
	// 获取有效(正在直播中)主播,剪枝roomIDs
	roomParams := &roomV1.RoomGetStatusInfoByUidsReq{Uids: groupList["all"], FilterOffline: 1, NeedBroadcastType: 1}
	// room
	roomResp, err := relationT.GetRoomInfo(ctx, roomParams)
	if err != nil {
		log.Error("[LiveAnchor]get_room_rpc_error")
		return
	}
	livingUfos := make([]int64, 0)
	livingRolaids := make([]int64, 0)
	// 没有人直播
	if len(roomResp) == 0 {
		return
	}
	for k, v := range roomResp {
		livingUfos = append(livingUfos, k)
		livingRolaids = append(livingRolaids, v.RoomId)
	}

	userResp := make(map[int64]*accountM.Card)
	roomCornerResp := make(map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result)
	pkResp := make(map[string]int64)
	attentionRoomListPlayURLMap := make(map[int64]*playurlbvc.PlayUrlItem)
	build, _ := strconv.ParseInt(req.Build, 10, 64)
	roomPendentParams := &roomV1.RoomPendantGetPendantByIdsReq{Ids: livingRolaids, Type: relationT.PendentMobileBadge, Position: relationT.PendentPosition}
	pkParams := &avV1.PkGetPkIdsByRoomIdsReq{RoomIds: livingRolaids, Platform: req.Platform}
	if err != nil {
		log.Error("[LiveAnchor]get_roomPendant_rpc_error")
		return
	}
	// user信息
	wg.Go(func() error {
		userResp, err = s.GetUserInfoData(ctx, livingUfos)
		return err
	})
	// room
	wg.Go(func() error {
		roomCornerResp, err = relationT.GetRoomPendantInfo(ctx, roomPendentParams)
		return err
	})
	// pk_id
	wg.Go(func() error {
		pkResp, err = relationT.GetPkID(ctx, pkParams)
		return err
	})

	quality := req.Quality
	if quality <= 0 {
		quality = 4
	}
	// playurl
	wg.Go(func() error {
		attentionRoomListPlayURLMap = dao.BvcApi.GetPlayUrlMulti(ctx, livingRolaids, 0, quality, build, req.Platform)
		return err
	})

	waitErr := wg.Wait()
	if waitErr != nil {
		log.Error("[LiveAnchor][step2] rpc error: %s", waitErr)
		return
	}

	// 下游数据收集完成
	mapSp := make([]int64, 0)
	normalSp := make([]int64, 0)
	mapSp = append(mapSp, groupList["special"]...)
	normalSp = append(normalSp, groupList["normal"]...)
	resp.Rooms = AdaptLivingField(roomResp, roomCornerResp, userResp, relationInfo, pkResp, attentionRoomListPlayURLMap, mapSp, normalSp, mapUfos2Rolaids)
	resp.TotalCount = int64(len(resp.Rooms))
	userExtParams := &userExV1.GrayRuleGetByMarkReq{Mark: relationT.App531GrayRule}
	grayRule, err := relationT.GetGrayRule(ctx, userExtParams)
	if err != nil {
		log.Error("[LiveAnchor]get_GrayRule_rpc_error")
		resp.BigCardType = 0
	} else if grayRule != nil {
		resp.BigCardType = relationT.App531ABTest(ctx, grayRule.Content, req.Build, req.Platform)
	}
	FilterType(ctx, livingUfos, resp, filterRule)
	SortType(ctx, resp, sortRule)
	return
}

// FilterType implementation
// [app端关注二级页]按照规则过滤结果集
func FilterType(ctx context.Context, targetUIDs []int64, originResult *v1pb.LiveAnchorResp, filterType int64) {
	if originResult == nil || len(originResult.Rooms) == 0 {
		return
	}
	switch filterType {
	case relationT.AppFilterDefault:
		{

		}
	case relationT.AppFilterFansMedal:
		{
			filteredRooms, _ := relationT.AppFilterRuleFansMedal(ctx, originResult, targetUIDs)
			originResult.Rooms = filteredRooms
			originResult.TotalCount = int64(len(originResult.Rooms))
		}
	case relationT.AppFilterGoldType:
		{
			filteredRooms, _ := relationT.AppFilterGold(ctx, originResult)
			originResult.Rooms = filteredRooms
			originResult.TotalCount = int64(len(originResult.Rooms))
		}
	}
}

// SortType implementation
// [app端关注二级页]按照规则排序结果集
// 规则见https://www.tapd.cn/20082211/prong/stories/view/1120082211001067961
func SortType(ctx context.Context, originResult *v1pb.LiveAnchorResp, sortType int64) (resp *v1pb.LiveAnchorResp, err error) {
	if originResult == nil || len(originResult.Rooms) == 0 {
		return
	}
	switch sortType {
	// 组间特别关注、组内人气值
	case relationT.AppSortDefaultT:
		{

		}
	case relationT.AppSortRuleLiveTimeT:
		{
			originResult.Rooms = relationT.AppSortRuleLiveTime(originResult)
		}
	case relationT.AppSortRuleOnlineT:
		{
			originResult.Rooms = relationT.AppSortRuleOnline(originResult)
		}
	case relationT.AppSortRuleGoldT:
		{
			originResult.Rooms = relationT.AppSortRuleGold(ctx, originResult)
		}
	default:
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

// GetUserFc ...
// 获取用户粉丝
func GetUserFc(ctx context.Context, UIDs []int64) (userResult map[int64]*relationV1.FeedGetUserFcBatchResp_RelationList, err error) {

	rpcChunkSize, RPCTimeout, err := relationT.GetChunkInfo(relationT.FansNum)
	params := relationT.ChunkCallInfo{ParamsName: "ufos", URLName: relationT.FansNum, ChunkSize: rpcChunkSize, RPCTimeout: RPCTimeout}
	userResult = make(map[int64]*relationV1.FeedGetUserFcBatchResp_RelationList)
	lens := len(UIDs)
	if lens <= 0 {
		return
	}
	// 批次
	params.ChunkNum = int64(math.Ceil(float64(lens) / float64(params.ChunkSize)))
	chunkResult := make([]map[int64]*relationV1.FeedGetUserFcBatchResp_RelationList, params.ChunkNum)
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
			ret, err := dao.RelationApi.V1Feed.GetUserFcBatch(ctx, &relationV1.FeedGetUserFcBatchReq{Uids: chunkUfosIds})
			if err != nil {
				err = errors.WithMessage(ecode.AccountGRPCError, "GET SEA PATROL FAIL")
				log.Error("Call main.Account.Cards Error.Infos(%+v) error(%+v)", chunkUfosIds, err)
			}
			chunkResult[x-1] = ret.Data
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		erelongInfo := relationT.ErrLogStrut{}
		erelongInfo.ErrType = "GoRoutingWaitError"
		erelongInfo.URLName = relationT.FansNum
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
				userResult[item.Uid] = item
			}
		}
	}
	return
}
