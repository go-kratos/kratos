package relation

import (
	"context"
	v1pb "go-common/app/interface/live/app-interface/api/http/v1"
	fansMedalV1 "go-common/app/service/live/fans_medal/api/liverpc/v1"
	"go-common/library/log"
)

const (
	// AppFilterDefault implementation
	// 默认排序
	AppFilterDefault = 0
	// AppFilterFansMedal implementation
	// 只看我有粉丝勋章的
	AppFilterFansMedal = 1
	// AppFilterGoldType implementation
	// 按照金瓜子排序
	AppFilterGoldType = 2
)

// AppFilterRuleFansMedal implementation
// [app端关注二级页]过滤粉丝勋章
func AppFilterRuleFansMedal(ctx context.Context, originResult *v1pb.LiveAnchorResp, targetUIDs []int64) (resp []*v1pb.LiveAnchorResp_Rooms, err error) {
	uid := GetUIDFromHeader(ctx)
	resp = make([]*v1pb.LiveAnchorResp_Rooms, 0)
	if originResult == nil || len(originResult.Rooms) == 0 {
		return
	}
	fansParams := &fansMedalV1.FansMedalTargetsWithMedalReq{Uid: uid, TargetIds: targetUIDs}
	hasMedalUIDs, err := GetFansMedal(ctx, fansParams)
	if err != nil {
		log.Error("[LiveAnchor][FilterType]get_FansMedal_rpc_error")
		resp = originResult.Rooms
		return
	}
	for _, v := range originResult.Rooms {
		if _, exist := hasMedalUIDs[v.Uid]; exist {
			resp = append(resp, v)
		}
	}
	return
}

// AppFilterGold implementation
// [app端关注二级页]过滤送礼
func AppFilterGold(ctx context.Context, originResult *v1pb.LiveAnchorResp) (resp []*v1pb.LiveAnchorResp_Rooms, err error) {
	giftInfo, err := GetGiftInfo(ctx)
	resp = make([]*v1pb.LiveAnchorResp_Rooms, 0)
	if originResult == nil || len(originResult.Rooms) == 0 {
		return
	}
	if err != nil {
		log.Error("[LiveAnchor][FilterType]get_RelationGift_rpc_error")
		resp = originResult.Rooms
		return
	}
	for _, v := range originResult.Rooms {
		if _, exist := giftInfo[v.Uid]; exist {
			resp = append(resp, v)
		}
	}
	return
}
