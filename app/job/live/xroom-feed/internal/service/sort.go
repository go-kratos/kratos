package service

import (
	"go-common/library/log"
	"sort"
)

type roomInfo struct {
	roomId int64
	value  int64
}

type attrSortStruct struct {
	roomIdInfos []roomInfo
	value       int64
}

type RoomInfoSlice []roomInfo

func (a RoomInfoSlice) Len() int {
	return len(a)
}

func (a RoomInfoSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a RoomInfoSlice) Less(i, j int) bool {
	return a[j].value < a[i].value
}

func (s *Service) sort(attrsFilter map[int64]map[int64]*attrFilter, filterRoomData []*filterItem, Cond string, recId int) (roomIds []int64) {
	roomIds = make([]int64, 0)
	noSort := make([]int64, 0)
	attrSort := make(map[int64]*attrSortStruct)
	attrTag := make([]int64, 0)
	if len(attrsFilter) <= 0 {
		for _, attrResp := range filterRoomData {
			if attrResp == nil {
				continue
			}
			noSort = append(noSort, attrResp.item.RoomId)
		}
	} else {
		//需要排序的
		for _, attrResp := range filterRoomData {
			if attrResp == nil {
				continue
			}
			for _, attr := range attrResp.item.AttrList {
				if _, ok := attrsFilter[attr.AttrId][attr.AttrSubId]; ok {
					if _, ok := attrSort[attr.AttrId]; !ok {
						attrSort[attr.AttrId] = &attrSortStruct{}
					}
					// 如果有top

					// 小时榜特殊处理
					if attr.AttrId == attrType[_hourRankType] {
						attr.AttrValue = 10 - attr.AttrValue
					}
					attrSort[attr.AttrId].value = attrsFilter[attr.AttrId][attr.AttrSubId].top
					attrSort[attr.AttrId].roomIdInfos = append(attrSort[attr.AttrId].roomIdInfos, roomInfo{
						roomId: attrResp.item.RoomId,
						value:  attr.AttrValue,
					})
				}
			}
			//or 逻辑需要非attr类型合并
			if attrResp.isTagHit {
				attrTag = append(attrTag, attrResp.item.RoomId)
			}
		}
	}

	log.Info("[sort]recId:%d, attrSort:%+v, Cond:%s, attrsFilter:%+v, filterRoomData:%+v", recId, attrSort, Cond, attrsFilter, filterRoomData)
	log.Info("[sort]recId:%d, attrTag:%+v, Cond:%s, attrsFilter:%+v, filterRoomData:%+v", recId, attrTag, Cond, attrsFilter, filterRoomData)
	log.Info("[sort]recId:%d, noSort:%+v, Cond:%s, attrsFilter:%+v, filterRoomData:%+v", recId, noSort, Cond, attrsFilter, filterRoomData)

	sortedList := make([][]int64, 0)
	//有top 需要排序的
	for _, obj := range attrSort {
		sortedList = append(sortedList, s.sortRoomList(obj.roomIdInfos, obj.value))
	}

	if len(sortedList) <= 0 {
		sortedList = append(sortedList, noSort)
	}

	// and 求交集
	if Cond == _condAnd {
		// 如果不填top
		roomIds = s.sliceIntersect(sortedList)
	}

	if Cond == _condOr {
		// or 条件时有些roomData可能没有attr，需要合并没有attr的
		sortedList = append(sortedList, attrTag)
		roomIds = s.sliceMerge(sortedList)
	}

	log.Info("[sort]recId:%d, roomIds:%+v, Cond:%s, attrsFilter:%+v, filterRoomData:%+v", recId, roomIds, Cond, attrsFilter, filterRoomData)

	return
}

func (s *Service) sortRoomList(roomIdInfos []roomInfo, value int64) (roomList []int64) {
	//value <= 0 不排序,直接返回
	if value <= 0 {
		for _, info := range roomIdInfos {
			roomList = append(roomList, info.roomId)
		}
		return
	}
	sort.Sort(RoomInfoSlice(roomIdInfos))
	// top
	if len(roomIdInfos) >= int(value) {
		roomIdInfos = roomIdInfos[:value]
	}

	for _, info := range roomIdInfos {
		roomList = append(roomList, info.roomId)
	}

	return
}

// 带去重 slice合并
func (s *Service) sliceMerge(roomList [][]int64) (mergedSlice []int64) {
	duplicateMap := make(map[int64]bool)
	mergedSlice = make([]int64, 0)
	for _, list := range roomList {
		for _, roomId := range list {
			if _, ok := duplicateMap[roomId]; !ok {
				mergedSlice = append(mergedSlice, roomId)
				duplicateMap[roomId] = true
			}
		}
	}

	return
}

func (s *Service) sliceIntersect(roomList [][]int64) (mergedSlice []int64) {
	mergedSlice = make([]int64, 0)
	if len(roomList) == 0 {
		return
	}

	if len(roomList) == 1 {
		mergedSlice = roomList[0]
		return
	}

	if len(roomList) >= 2 {
		i := 3
		mergedSlice = s.intersect(roomList[0], roomList[1])
		for {
			if len(roomList) < i {
				break
			}
			mergedSlice = s.intersect(mergedSlice, roomList[int64(i-1)])
			i++
		}
	}
	return
}

// 带去重 slice交集 slice已排序
func (s *Service) intersect(nums1 []int64, nums2 []int64) (IntersectSlice []int64) {
	IntersectSlice = make([]int64, 0)
	x := 0
	y := 0
	for {
		if x < len(nums1) && y < len(nums2) {
			if nums1[x] == nums2[y] {
				IntersectSlice = append(IntersectSlice, nums1[x])
				x++
				y++
			} else if nums1[x] > nums2[y] {
				y++
			} else {
				x++
			}
		} else {
			break
		}

	}
	return
}
