package v2

import (
	"context"
	"fmt"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
	"go-common/app/interface/live/app-interface/dao"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	"go-common/library/log"
	"go-common/library/xstr"
)

type multiModuleInfo struct {
	areaName     string
	parentAreaId int64
	areaId       int64
}

const _areaModuleLink = "https://live.bilibili.com/app/area?parent_area_id=%d&parent_area_name=%s&area_id=%d&area_name=%s"

func (s *IndexService) getCommonRoomListForIndex(ctx context.Context, build int64, platform string, quality int64) (respCommonRoomList []*v2pb.MRoomBlock) {
	respCommonRoomList = make([]*v2pb.MRoomBlock, 0)
	moduleInfoMap := s.GetAllModuleInfoMapFromCache(ctx)
	roomModuleIDs := make([]int64, 0)
	for _, moduleType := range s.commonType {
		if _, ok := moduleInfoMap[moduleType]; ok {
			for _, moduleInfo := range moduleInfoMap[moduleType] {
				respCommonRoomList = append(respCommonRoomList, &v2pb.MRoomBlock{
					ModuleInfo: moduleInfo,
				})
				roomModuleIDs = append(roomModuleIDs, moduleInfo.Id)
			}
		}
	}
	s.getCommonRoomList(ctx, roomModuleIDs, respCommonRoomList)

	return
}

func (s *IndexService) getCommonRoomList(ctx context.Context, moduleIds []int64, respCommonRoomList []*v2pb.MRoomBlock) (err error) {
	if len(moduleIds) <= 0 || len(respCommonRoomList) <= 0 {
		return
	}

	roomListMap := make(map[int64]*roomV2.AppIndexGetRoomListByIdsResp_RoomList)
	roomListMap = s.getCommonListFromCache(moduleIds)

	if roomListMap == nil {
		log.Error("[commonList]GetListByIds error, res:%+v", roomListMap)
		return
	}

	for _, moduleInfo := range respCommonRoomList {
		moduleInfo.List = make([]*v2pb.CommonRoomItem, 0)
		if _, ok := roomListMap[moduleInfo.ModuleInfo.Id]; ok {
			for _, roomList := range roomListMap[moduleInfo.ModuleInfo.Id].List {
				if moduleInfo.ModuleInfo.Type == _parentAreaFormType ||
					moduleInfo.ModuleInfo.Type == _parentAreaSquareType ||
					moduleInfo.ModuleInfo.Type == _yunyingRecSquareType ||
					moduleInfo.ModuleInfo.Type == _yunyingRecFormType {

					if len(moduleInfo.List) >= 4 {
						break
					}
				}
				moduleInfo.List = append(moduleInfo.List, &v2pb.CommonRoomItem{
					Roomid:           roomList.Roomid,
					Title:            roomList.Title,
					Uname:            roomList.Uname,
					AreaV2Id:         roomList.AreaV2Id,
					AreaV2Name:       roomList.AreaV2Name,
					AreaV2ParentId:   roomList.AreaV2ParentId,
					AreaV2ParentName: roomList.AreaV2ParentName,
					Online:           roomList.Online,
					Face:             roomList.Face,
					Cover:            roomList.Cover,
					BroadcastType:    roomList.BroadcastType,
					CurrentQuality:   roomList.CurrentQuality,
					AcceptQuality:    roomList.AcceptQuality,
					RecType:          roomList.RecType,
					PendentRu:        roomList.PendentRu,
					PendentRuColor:   roomList.PendentRuColor,
					PendentRuPic:     roomList.PendentRuPic,
					PlayUrl:          roomList.PlayUrl,
					PlayUrlH265:      roomList.PlayUrlH265,
					PkId:             roomList.PkId,
				})
			}
		}
	}

	return
}

func (s *IndexService) getCommonRoomListByID(ctx context.Context, moduleID int64, build int64, platform string, quality int64, device string, duplicates []int64) (respCommonRoomList []*v2pb.CommonRoomItem, err error) {
	respCommonRoomList = make([]*v2pb.CommonRoomItem, 0)
	duplicatesMap := make(map[int64]bool)

	for _, roomID := range duplicates {
		duplicatesMap[roomID] = true
	}
	moduleExistIds := make([]int64, 0)
	if moduleID == 0 {
		return
	}

	roomListMap, err := s.roomDao.GetListByIds(ctx, []int64{moduleID})
	if err != nil {
		log.Error("[commonList]GetAllModuleListByIds error, error:%+v", err)
		return
	}

	// 24个对关注去重的item
	if _, ok := roomListMap[moduleID]; ok {
		list := roomListMap[moduleID].List
		for _, info := range list {
			if len(respCommonRoomList) >= 24 {
				break
			}
			if _, ok := duplicatesMap[info.Roomid]; !ok {
				respCommonRoomList = append(respCommonRoomList, &v2pb.CommonRoomItem{
					Roomid:           info.Roomid,
					Title:            info.Title,
					Uname:            info.Uname,
					AreaV2Id:         info.AreaV2Id,
					AreaV2Name:       info.AreaV2Name,
					AreaV2ParentId:   info.AreaV2ParentId,
					AreaV2ParentName: info.AreaV2ParentName,
					Online:           info.Online,
					Face:             info.Face,
					Cover:            info.Cover,
					BroadcastType:    info.BroadcastType,
					CurrentQuality:   info.CurrentQuality,
					AcceptQuality:    info.AcceptQuality,
					RecType:          info.RecType,
					PendentRu:        info.PendentRu,
					PendentRuColor:   info.PendentRuColor,
					PendentRuPic:     info.PendentRuPic,
					PlayUrl:          info.PlayUrl,
					PlayUrlH265:      info.PlayUrlH265,
					PkId:             info.PkId,
				})
				moduleExistIds = append(moduleExistIds, info.Roomid)
			}
		}

		filterList := make([]*v2pb.CommonRoomItem, 0)
		if device != "pad" {
			//投放位覆盖
			recPool, recPoolRooms := s.fourTimeRecPoolForYumo(ctx, moduleExistIds, duplicates)
			log.Info("[getCommonRoomListByID]投放位 recPool: %+v, moduleExistIds: %+v, otherExistIds:%+v", recPoolRooms, xstr.JoinInts(moduleExistIds), xstr.JoinInts(duplicates))

			duplicateMap := make(map[int64]bool)
			for _, id := range recPoolRooms {
				duplicateMap[id] = true
			}
			newRecFilterCommonList := make([]*v2pb.CommonRoomItem, 0)
			for _, room := range respCommonRoomList {
				if _, ok := duplicateMap[room.Roomid]; !ok {
					newRecFilterCommonList = append(newRecFilterCommonList, room)
				}
			}

			for i := 0; i < 24; i++ {
				position := int64(i) + 1
				if item, ok := recPool[position]; ok {
					filterList = append(filterList, item)
				} else {
					if len(newRecFilterCommonList) <= 0 {
						continue
					}
					filterList = append(filterList, newRecFilterCommonList[0:1][0])
					newRecFilterCommonList = newRecFilterCommonList[1:]
				}
			}
			respCommonRoomList = filterList
		}

		//获取playurl
		s.getPlayUrl(ctx, respCommonRoomList, quality, build, platform)
	}
	return
}

func (s *IndexService) getMultiRoomList(ctx context.Context, myTag []*v2pb.MMyTag, platform string, build int64, quality int64) (respMyTagRoomList []*v2pb.MRoomBlock, existAreaMap map[int64]bool) {
	// 未登陆
	respMyTagRoomList = make([]*v2pb.MRoomBlock, 0)
	existAreaMap = make(map[int64]bool) // for duplicate yunying rec
	if len(myTag) <= 0 {
		log.Warn("[getMultiRoomList]my tag empty!")
		return
	}

	parentName := map[int64]string{
		1: "娱乐",
		2: "游戏",
		3: "手游",
		4: "绘画",
		5: "电台",
	}
	areaIds := make([]int64, 0)
	multiModuleInfos := make([]*multiModuleInfo, 0)
	for _, tag := range myTag {
		// 默认tag没有list 跳过
		if tag.ExtraInfo.IsGray == 0 {
			continue
		}

		for _, item := range tag.List {
			if item.AreaV2Id == 0 {
				// 过滤异常或全部标签
				continue
			}
			areaIds = append(areaIds, item.AreaV2Id)
			multiModuleInfos = append(multiModuleInfos, &multiModuleInfo{
				areaId:       item.AreaV2Id,
				areaName:     item.AreaV2Name,
				parentAreaId: item.AreaV2ParentId,
			})
		}
	}

	if len(areaIds) <= 0 {
		log.Info("[getMultiRoomList]no gray tag, so no room list")
		return
	}
	multiRoomListMap, err := s.roomDao.GetMultiRoomList(ctx, xstr.JoinInts(areaIds), platform)
	if err != nil {
		log.Error("[getMultiRoomList]roomDao.GetMultiRoomList get error:%+v", err)
		return
	}

	moduleInfoMap := s.GetAllModuleInfoMapFromCache(ctx)
	myTagRoomListTypeMap, exist := moduleInfoMap[_myAreaTagListType]
	if !exist || myTagRoomListTypeMap == nil || len(myTagRoomListTypeMap) <= 0 {
		log.Info("[getMultiRoomList]my tag room list module not exist, all: %+v", moduleInfoMap)
		return
	}
	roomIds := make([]int64, 0)
	for index, moduleInfo := range myTagRoomListTypeMap {
		if len(multiModuleInfos) <= index {
			// 后台多配的case，防止溢出
			continue
		}
		mInfo := multiModuleInfos[index]
		if mInfo == nil {
			continue
		}
		moduleInfo.Title = mInfo.areaName
		moduleInfo.Link = fmt.Sprintf(_areaModuleLink, mInfo.parentAreaId, parentName[mInfo.parentAreaId], mInfo.areaId, mInfo.areaName)
		item := &v2pb.MRoomBlock{
			ModuleInfo: moduleInfo,
		}

		l, ok := multiRoomListMap[mInfo.areaId]
		if !ok || l == nil {
			continue
		}
		existAreaMap[mInfo.areaId] = true

		innerList := make([]*v2pb.CommonRoomItem, 0)
		for _, v := range l {
			roomIds = append(roomIds, v.Roomid)
			innerList = append(innerList, &v2pb.CommonRoomItem{
				Roomid:           v.Roomid,
				Title:            v.Title,
				Uname:            v.Uname,
				AreaV2Id:         v.AreaV2Id,
				AreaV2Name:       v.AreaV2Name,
				AreaV2ParentId:   v.AreaV2ParentId,
				AreaV2ParentName: v.AreaV2ParentName,
				Online:           v.Online,
				Face:             v.Face,
				Cover:            v.Cover,
				BroadcastType:    v.BroadcastType,
				CurrentQuality:   v.CurrentQuality,
				AcceptQuality:    v.AcceptQuality,
				RecType:          v.RecType,
				PendentRu:        v.PendentRu,
				PendentRuColor:   v.PendentRuColor,
				PendentRuPic:     v.PendentRuPic,
				PlayUrl:          v.PlayUrl,
				PlayUrlH265:      v.PlayUrlH265,
				PkId:             v.PkId,
			})
		}
		item.List = innerList
		respMyTagRoomList = append(respMyTagRoomList, item)
	}

	// 拼playurl
	changeRoomListPlayURLMap := dao.BvcApi.GetPlayUrlMulti(ctx, roomIds, 0, quality, build, platform)

	for _, v := range respMyTagRoomList {
		for _, vv := range v.List {
			if changeRoomListPlayURLMap[vv.Roomid] != nil {
				vv.AcceptQuality = changeRoomListPlayURLMap[vv.Roomid].AcceptQuality
				vv.CurrentQuality = changeRoomListPlayURLMap[vv.Roomid].CurrentQuality
				vv.PlayUrl = changeRoomListPlayURLMap[vv.Roomid].Url["h264"]
				vv.PlayUrlH265 = changeRoomListPlayURLMap[vv.Roomid].Url["h265"]
			}
		}
	}

	return
}
