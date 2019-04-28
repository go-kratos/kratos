package relation

import (
	"context"
	v1pb "go-common/app/interface/live/app-interface/api/http/v1"
	"go-common/library/log"
	"sort"
)

const (
	// AppSortDefaultT ...
	// 默认排序
	AppSortDefaultT = 0
	// AppSortRuleLiveTimeT ...
	// 开播时间倒序
	AppSortRuleLiveTimeT = 1
	// AppSortRuleOnlineT ...
	// 人气值倒序
	AppSortRuleOnlineT = 2
	// AppSortRuleGoldT ...
	// 金瓜子倒序
	AppSortRuleGoldT = 3
)

// SendGift ...
// [app端关注二级页]按照金瓜子排序结构
type SendGift struct {
	Mid  int64
	gold int64
}

// SortLiveTime ... implementation
// [app端关注二级页]按照开播时间排序
type SortLiveTime []*v1pb.LiveAnchorResp_Rooms

// SortOnlineTime ... implementation
// [app端关注二级页]按照开播时间排序
type SortOnlineTime []*v1pb.LiveAnchorResp_Rooms

// // AddGoldRoomInfo ... implementation
// // [app端关注二级页]按照开播时间排序
// type AddGoldRoomInfo []*LiveAnchorRespRoomsAddGold

// SortUIDGift ... implementation
// [app端关注二级页]按照开播时间排序
type SortUIDGift []SendGift

// AppSortRuleLiveTime implementation
// [app端关注二级页]按照开播时间排序
func (p SortLiveTime) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// AppSortRuleLiveTime implementation
// [app端关注二级页]按照开播时间排序
func (p SortLiveTime) Len() int { return len(p) }

// AppSortRuleLiveTime implementation
// [app端关注二级页]按照开播时间排序
func (p SortLiveTime) Less(i, j int) bool { return p[i].LiveTime > p[j].LiveTime }

// AppSortRuleLiveTime implementation
// [app端关注二级页]按照开播时间排序
func (p SortOnlineTime) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// AppSortRuleLiveTime implementation
// [app端关注二级页]按照开播时间排序
func (p SortOnlineTime) Len() int { return len(p) }

// AppSortRuleLiveTime implementation
// [app端关注二级页]按照开播时间排序
func (p SortOnlineTime) Less(i, j int) bool { return p[i].Online > p[j].Online }

// Swap
// [app端关注二级页]自定义排序结构
func (p SortUIDGift) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Len
// [app端关注二级页]自定义排序结构
func (p SortUIDGift) Len() int { return len(p) }

// Less
// [app端关注二级页]自定义排序结构
func (p SortUIDGift) Less(i, j int) bool { return p[i].gold > p[j].gold }

// AppSortRuleLiveTime implementation
// [app端关注二级页]按照开播时间排序
func AppSortRuleLiveTime(originResult *v1pb.LiveAnchorResp) (resp []*v1pb.LiveAnchorResp_Rooms) {
	resp = make([]*v1pb.LiveAnchorResp_Rooms, 0)
	if originResult == nil {
		return
	}
	p := make(SortLiveTime, len(originResult.Rooms))
	i := 0
	for _, v := range originResult.Rooms {
		p[i] = &v1pb.LiveAnchorResp_Rooms{
			Roomid:           v.Roomid,
			Uid:              v.Uid,
			Uname:            v.Uname,
			Face:             v.Face,
			Title:            v.Title,
			LiveTagName:      v.LiveTagName,
			LiveTime:         v.LiveTime,
			Online:           v.Online,
			Playurl:          v.Playurl,
			AcceptQuality:    v.AcceptQuality,
			CurrentQuality:   v.CurrentQuality,
			PkId:             v.PkId,
			SpecialAttention: v.SpecialAttention,
			Area:             v.Area,
			AreaName:         v.AreaName,
			AreaV2Id:         v.AreaV2Id,
			AreaV2Name:       v.AreaV2Name,
			AreaV2ParentName: v.AreaV2ParentName,
			AreaV2ParentId:   v.AreaV2ParentId,
			BroadcastType:    v.BroadcastType,
			OfficialVerify:   v.OfficialVerify,
			Link:             v.Link,
			Cover:            v.Cover,
			PendentRu:        v.PendentRu,
			PendentRuColor:   v.PendentRuColor,
			PendentRuPic:     v.PendentRuPic}
		i++
	}
	sort.Sort(p)
	resp = p
	return
}

// AppSortRuleOnline implementation
// [app端关注二级页]按照人气值排序
func AppSortRuleOnline(originResult *v1pb.LiveAnchorResp) (resp []*v1pb.LiveAnchorResp_Rooms) {
	resp = make([]*v1pb.LiveAnchorResp_Rooms, 0)
	if originResult == nil {
		return
	}
	p := make(SortOnlineTime, len(originResult.Rooms))
	i := 0
	for _, v := range originResult.Rooms {
		p[i] = &v1pb.LiveAnchorResp_Rooms{
			Roomid:           v.Roomid,
			Uid:              v.Uid,
			Uname:            v.Uname,
			Face:             v.Face,
			Title:            v.Title,
			LiveTagName:      v.LiveTagName,
			LiveTime:         v.LiveTime,
			Online:           v.Online,
			Playurl:          v.Playurl,
			AcceptQuality:    v.AcceptQuality,
			CurrentQuality:   v.CurrentQuality,
			PkId:             v.PkId,
			SpecialAttention: v.SpecialAttention,
			Area:             v.Area,
			AreaName:         v.AreaName,
			AreaV2Id:         v.AreaV2Id,
			AreaV2Name:       v.AreaV2Name,
			AreaV2ParentName: v.AreaV2ParentName,
			AreaV2ParentId:   v.AreaV2ParentId,
			BroadcastType:    v.BroadcastType,
			OfficialVerify:   v.OfficialVerify,
			Link:             v.Link,
			Cover:            v.Cover,
			PendentRu:        v.PendentRu,
			PendentRuColor:   v.PendentRuColor,
			PlayUrlH265:      v.PlayUrlH265,
			PendentRuPic:     v.PendentRuPic}
		i++
	}
	sort.Sort(p)
	resp = p
	return
}

// AppSortRuleGold implementation
// [app端关注二级页]按照送礼排序
func AppSortRuleGold(ctx context.Context, originResult *v1pb.LiveAnchorResp) (resp []*v1pb.LiveAnchorResp_Rooms) {
	resp = make([]*v1pb.LiveAnchorResp_Rooms, 0)
	if originResult == nil {
		return
	}

	giftInfo, err := GetGiftInfo(ctx)
	if err != nil {
		log.Error("[LiveAnchor][FilterType][AppSortRuleGold]get_RelationGift_rpc_error")
		resp = originResult.Rooms
		return
	}
	if len(giftInfo) == 0 {
		resp = AppSortRuleOnline(originResult)
		return
	}

	respHasGold := make([]*v1pb.LiveAnchorResp_Rooms, 0)
	respNoGold := make([]*v1pb.LiveAnchorResp_Rooms, 0)
	GiftRank := make(map[int64]int64)
	GiftNoGold := make([]int64, 0)
	// 计算金瓜子排行榜,uid分key
	for _, v := range originResult.Rooms {
		roomUID := v.Uid
		if _, exist := giftInfo[roomUID]; exist {
			GiftRank[roomUID] += giftInfo[roomUID]
		}
	}
	sorted := SortMap(GiftRank)

	// 没有送礼的用户
	for _, v := range originResult.Rooms {
		if _, exist := GiftRank[v.Uid]; !exist {
			GiftNoGold = append(GiftNoGold, v.Uid)
		}
	}

	for _, vv := range sorted {
		for _, v := range originResult.Rooms {
			if v.Uid == vv.Key {
				respHasGold = append(respHasGold, v)
			}
		}
	}

	for _, v := range originResult.Rooms {
		for _, vv := range GiftNoGold {
			if v.Uid == vv {
				respNoGold = append(respNoGold, v)
			}
		}
	}
	tempLiveAnchor := &v1pb.LiveAnchorResp{}
	tempLiveAnchor.Rooms = respNoGold
	respNoGoldSorted := AppSortRuleOnline(tempLiveAnchor)
	resp = append(resp, respHasGold...)
	resp = append(resp, respNoGoldSorted...)
	return
}
