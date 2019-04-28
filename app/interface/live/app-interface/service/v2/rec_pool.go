package v2

import (
	"math"
	"strconv"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
	"go-common/app/interface/live/app-interface/dao"
	xrf "go-common/app/service/live/xroom-feed/api"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"

	"context"

	"go-common/library/log"
)

const _appModuleType = 1

// getRecPoolList
func (s *IndexService) getRecPoolList(ctx context.Context, module int64, num int64, moduleExists []int64, otherExists []int64) (list map[int64]*xrf.RoomItem) {
	list = make(map[int64]*xrf.RoomItem)
	req := &xrf.RecPoolReq{
		ModuleType:       module,
		PositionNum:      num,
		ModuleExistRooms: xstr.JoinInts(moduleExists),
		OtherExistRooms:  xstr.JoinInts(otherExists),
		From:             "app-interface",
	}
	resp, ferr := s.xrfClient.RecPoolClient.GetList(ctx, req)
	if ferr != nil {
		log.Error("[rec_pool]GetList from xroom-feed err:%+v", ferr.Error())
		return
	}

	list = resp.List
	return
}

func (s *IndexService) getRecPoolAllPosition(ctx context.Context, moduleRoomIDs, otherRoomIDs []int64) (recPool map[int64]*v2pb.CommonRoomItem) {
	recPool = make(map[int64]*v2pb.CommonRoomItem)

	uids := make([]int64, 0)
	roomIds := make([]int64, 0)
	mapRoomIdToUid := make(map[int64]int64)

	roomItem := s.getRecPoolList(ctx, _appModuleType, 6, moduleRoomIDs, otherRoomIDs)
	if len(roomItem) <= 0 {
		log.Info("[getRecPoolAllPosition]getRecPoolList roomItem empty")
		return
	}

	for i, item := range roomItem {
		uids = append(uids, item.Uid)
		roomIds = append(roomIds, item.RoomId)
		mapRoomIdToUid[item.RoomId] = item.Uid
		recPool[i] = &v2pb.CommonRoomItem{
			Roomid:           item.RoomId,
			Title:            item.Title,
			Link:             "/" + strconv.Itoa(int(item.RoomId)),
			AreaV2Id:         item.AreaId,
			AreaV2Name:       item.AreaName,
			AreaV2ParentId:   item.ParentAreaId,
			AreaV2ParentName: item.ParentAreaName,
			Online:           item.PopularityCount,
			Cover:            s.getCover(item.Cover, item.Keyframe),
			RecType:          item.RecType,
			AcceptQuality:    []int64{2, 4},
		}
	}

	if len(recPool) <= 0 {
		return
	}
	extraInfo := s.getExtraDataForRoom(ctx, roomIds, uids, mapRoomIdToUid)

	for _, pool := range recPool {
		if _, ok := extraInfo[pool.Roomid]; ok {
			pool.Uname = extraInfo[pool.Roomid].UName
			pool.Face = extraInfo[pool.Roomid].Face
			pool.PendentRu = extraInfo[pool.Roomid].PendentRu
			pool.PendentRuPic = extraInfo[pool.Roomid].PendentRuPic
			pool.PendentRuColor = extraInfo[pool.Roomid].PendentRuColor
		}
	}

	return
}

func (s *IndexService) handleCommonRoomList(ctx context.Context, respMyIdol *v2pb.MMyIdol, respCommonRoomList []*v2pb.MRoomBlock, quality, build int64, platform string, device string) []*v2pb.MRoomBlock {
	moduleExistIds := make([]int64, 0)
	otherExistIds := make([]int64, 0)
	duplicateMap := make(map[int64]bool)
	for _, idol := range respMyIdol.List {
		otherExistIds = append(otherExistIds, idol.Roomid)
		duplicateMap[idol.Roomid] = true
	}
	for _, roomBlock := range respCommonRoomList {
		if roomBlock.ModuleInfo.Type == _recFormType || roomBlock.ModuleInfo.Type == _recSquareType {
			newRecRoomList := make([]*v2pb.CommonRoomItem, 0)
			for _, item := range roomBlock.List {
				if len(newRecRoomList) >= 24 {
					break
				}
				if _, ok := duplicateMap[item.Roomid]; !ok {
					newRecRoomList = append(newRecRoomList, item)
					moduleExistIds = append(moduleExistIds, item.Roomid)
				}
			}

			// 投放位
			if device != "pad" {
				recPool, recPoolRooms := s.fourTimeRecPoolForYumo(ctx, moduleExistIds, otherExistIds)

				log.Info("[handleCommonRoomList]投放位 recPool: %+v, moduleExistIds: %+v, otherExistIds:%+v", recPoolRooms, xstr.JoinInts(moduleExistIds), xstr.JoinInts(otherExistIds))

				duplicateRecMap := make(map[int64]bool)
				for _, id := range recPoolRooms {
					duplicateRecMap[id] = true
				}
				newRecFilterCommonList := make([]*v2pb.CommonRoomItem, 0)
				for _, room := range newRecRoomList {
					if room == nil {
						continue
					}
					if _, ok := duplicateRecMap[room.Roomid]; !ok {
						newRecFilterCommonList = append(newRecFilterCommonList, room)
					}
				}
				filterList := make([]*v2pb.CommonRoomItem, 0)

				for i := 0; i < 24; i++ {
					position := int64(i) + 1
					if item, ok := recPool[position]; ok {
						filterList = append(filterList, item)
						continue
					}
					if len(newRecFilterCommonList) <= 0 {
						continue
					}
					filterList = append(filterList, newRecFilterCommonList[0:1][0])
					newRecFilterCommonList = newRecFilterCommonList[1:]
				}
				roomBlock.List = filterList
			} else {
				roomBlock.List = newRecRoomList
			}
		}
	}

	// 拼playurl
	roomIds := make([]int64, 0)
	for _, commRoomBlock := range respCommonRoomList {
		for _, roomList := range commRoomBlock.List {
			roomIds = append(roomIds, roomList.Roomid)
		}
	}
	changeRoomListPlayURLMap := dao.BvcApi.GetPlayUrlMulti(ctx, roomIds, 0, quality, build, platform)

	for _, v := range respCommonRoomList {
		for _, vv := range v.List {
			if changeRoomListPlayURLMap[vv.Roomid] != nil {
				vv.AcceptQuality = changeRoomListPlayURLMap[vv.Roomid].AcceptQuality
				vv.CurrentQuality = changeRoomListPlayURLMap[vv.Roomid].CurrentQuality
				vv.PlayUrl = changeRoomListPlayURLMap[vv.Roomid].Url["h264"]
				vv.PlayUrlH265 = changeRoomListPlayURLMap[vv.Roomid].Url["h265"]
			}
		}
	}

	return respCommonRoomList
}

// 一次请求 变为四次请求xroom feed
func (s *IndexService) fourTimeRecPoolForYumo(ctx context.Context, moduleExistIds []int64, otherExistIds []int64) (recPool map[int64]*v2pb.CommonRoomItem, recPoolRooms map[int64]int64) {
	recPool = make(map[int64]*v2pb.CommonRoomItem)
	recPoolRooms = make(map[int64]int64)
	pageSize := 6
	if len(moduleExistIds) < pageSize {
		return
	}
	page := int64(math.Floor(float64(len(moduleExistIds)) / float64(pageSize)))

	result := make([]map[int64]*v2pb.CommonRoomItem, page)
	wg := errgroup.Group{}
	var i int64
	for i = 1; i <= page; i++ {
		p := i - 1
		start := p * int64(pageSize)
		end := start + int64(pageSize)
		duplicateIds := moduleExistIds[start:end]

		wg.Go(func() (err error) {
			result[p] = s.getRecPoolAllPosition(ctx, duplicateIds, otherExistIds)
			return nil
		})
	}
	err := wg.Wait()
	if err != nil {
		log.Error("[fourTimeRecPoolForYumo]moduleExistIds: %+v, otherExistIds:%+v", xstr.JoinInts(moduleExistIds), xstr.JoinInts(otherExistIds))
		return
	}
	for page, recPoolMap := range result {
		for position, item := range recPoolMap {
			if item == nil {
				continue
			}
			p := position + int64(page*pageSize)
			recPool[p] = item
			recPoolRooms[p] = item.Roomid
		}
	}

	return
}
