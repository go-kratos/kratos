package v2

import (
	"context"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
	recommendV1 "go-common/app/service/live/recommend/api/grpc/v1"
	"go-common/library/log"
)

const (
	_recTypeRecommend = 5
	_recNum           = 6 //获取推荐的数量
)

func (s *IndexService) getLiveRecRoomList(ctx context.Context, respMyIdol *v2pb.MMyIdol, mid int64, build int64, platform string, recPage int64, quality int64) (respLiuGangRecRoomList []*v2pb.CommonRoomItem, err error) {
	duplicate := make([]int64, 0)
	// ctx可以换带cancel或timeout的
	for _, idol := range respMyIdol.List {
		duplicate = append(duplicate, idol.Roomid)
	}

	return s.getRecInfo(ctx, mid, duplicate, build, platform, recPage, quality)
}

func (s *IndexService) getLiveRecRoomListForChange(ctx context.Context, mid int64, build int64, platform string, duplicate []int64, recPage int64, quality int64) (respLiuGangRecRoomList []*v2pb.CommonRoomItem, err error) {
	return s.getRecInfo(ctx, mid, duplicate, build, platform, recPage, quality)
}

func (s *IndexService) getRecInfo(ctx context.Context, mid int64, duplicate []int64, build int64, platform string, recPage int64, quality int64) (respLiveRecRoomList []*v2pb.CommonRoomItem, err error) {
	// 天马对关注去重
	duplicates := duplicate

	idolDuplicateMap := make(map[int64]bool)

	for _, id := range duplicates {
		if _, ok := idolDuplicateMap[id]; !ok {
			idolDuplicateMap[id] = true
		}
	}
	// 获取强推
	strongRecLen := 0
	//不考虑位置好的
	recPool := s.getRecPoolAllPosition(ctx, nil, duplicates)
	// 获取强推
	if len(recPool) > 0 {
		for _, strongInfo := range recPool {
			if strongInfo.Roomid == 0 {
				continue
			}
			if _, ok := idolDuplicateMap[strongInfo.Roomid]; !ok {
				duplicates = append(duplicates, strongInfo.Roomid)
				strongRecLen++
			}
		}
	}

	respLiveRecRoomList = make([]*v2pb.CommonRoomItem, 0)
	count := _recNum - strongRecLen
	if count <= 0 {
		count = _recNum
	}
	GetRandomRecResp, err := s.recommendConn.RandomRecsByUser(ctx, &recommendV1.GetRandomRecReq{
		Uid:      mid,
		Count:    uint32(count), // 首页6个推荐
		ExistIds: duplicates,
	})

	if err != nil {
		log.Error("[GetLiveRoomList]GetLiveRecResp err, err:%+v", err)
		return
	}

	if GetRandomRecResp == nil {
		log.Error("[GetLiveRoomList]GetLiveRecResp empty err")
		return
	}

	if len(GetRandomRecResp.RoomIds) < count {
		log.Info("[GetLiveRoomList]GetLiveRecResp not enough num:%d,mid:%d", len(GetRandomRecResp.RoomIds), mid)
		return
	}

	respLiveRecRoomList, err = s.getRecRoomList(ctx, GetRandomRecResp.RoomIds, recPool, build, platform, idolDuplicateMap, _recTypeRecommend, quality)
	if err != nil {
		log.Error("[GetLiveRoomList]FillLiveRecRoomList err:%+v", err)
	}
	return
}
