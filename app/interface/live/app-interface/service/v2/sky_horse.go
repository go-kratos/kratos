package v2

import (
	"context"

	"github.com/pkg/errors"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
	"go-common/app/interface/live/app-interface/dao"
	"go-common/library/ecode"
)

const (
	_skyHorseRecTimeOut           = 100
	_recTypeForce                 = 3
	_mobileIndexBadgeColorDefault = "#FB9E60"
	_recTypeSkyHorse              = 4
)

func (s *IndexService) getSkyHorseRoomListForIndex(ctx context.Context, respMyIdol *v2pb.MMyIdol, mid int64, buvid string, build int64, platform string, recPage int64, quality int64) (respSkyHorseRoomList []*v2pb.CommonRoomItem, err error) {
	respSkyHorseRoomList = make([]*v2pb.CommonRoomItem, 0)
	duplicate := make([]int64, 0)
	// ctx可以换带cancel或timeout的
	for _, idol := range respMyIdol.List {
		duplicate = append(duplicate, idol.Roomid)
	}
	respSkyHorseRoomList, err = s.getSkyHorseRoomList(ctx, mid, buvid, build, platform, duplicate, recPage, quality)
	if err != nil {
		return
	}
	return
}

func (s *IndexService) getSkyHorseRoomList(ctx context.Context, uid int64, buvid string, build int64, platform string, idolIds []int64, recPage int64, quality int64) (respSkyHorseRoomList []*v2pb.CommonRoomItem, err error) {
	// 天马对关注去重
	duplicates := idolIds

	idolDuplicateMap := make(map[int64]bool)

	for _, id := range duplicates {
		if _, ok := idolDuplicateMap[id]; !ok {
			idolDuplicateMap[id] = true
		}
	}

	strongRecLen := 0
	//天马不考虑位置好的
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

	skyHorseRec, skyHorseErr := dao.SkyHorseApi.GetSkyHorseRec(ctx, uid, buvid, build, platform, duplicates, strongRecLen, _skyHorseRecTimeOut)
	if skyHorseErr != nil {
		err = errors.WithMessage(ecode.SkyHorseError, "")
		return
	}

	roomIds := make([]int64, 0)
	for _, skyHorseInfo := range skyHorseRec.Data {
		roomIds = append(roomIds, int64(skyHorseInfo.Id))
	}

	return s.getRecRoomList(ctx, roomIds, recPool, build, platform, idolDuplicateMap, _recTypeSkyHorse, quality)
}
