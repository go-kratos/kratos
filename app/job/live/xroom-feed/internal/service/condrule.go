package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/live/xroom-feed/internal/model"
	daoAnchorV1 "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/library/log"
)

const (
	_isAll           = -1
	_roomStatusTagId = 3
	_onlineCurrent   = 3
)

type attrFilter struct {
	attrId    int64
	attrSubId int64
	max       int64
	min       int64
	top       int64
}

var attrType = map[string]int64{
	_onlineType:   1,
	_incomeType:   2,
	_dmsType:      3,
	_hourRankType: 4,
	_liveDaysType: 5,
}

var tagType = map[string]bool{
	_areaType:       true,
	_roomStatusType: true,
	_anchorCateType: true,
}

type filterItem struct {
	isTagHit bool
	item     *daoAnchorV1.AttrResp
}

// ParentAreaIds ...
var ParentAreaIds []int64

// AreaIds ...
var AreaIds [][]int64

// getConditionTypeIndex ...
func (s *Service) getConditionTypeIndex(conds []*model.RuleProtocol, sType string) (i int64) {
	i = -1
	for index, cond := range conds {
		if cond == nil {
			continue
		}
		if cond.ConfType == sType {
			i = int64(index)
			break
		}
	}
	return
}

// genCondConfRoomList ... recId for log
func (s *Service) genCondConfRoomList(ctx context.Context, Condition []*model.RuleProtocol, Cond string, recId int) (roomIds []int64) {
	ParentAreaIds := make([]int64, 0)
	AreaIds := make([][]int64, 0)
	areaFilter := make(map[int64]map[int64]bool)
	roomStatus := make(map[int64]bool)
	anchorCate := make(map[int64]bool)
	attrs := make([]*daoAnchorV1.AttrReq, 0)
	attrsFilter := make(map[int64]map[int64]*attrFilter)
	for _, cond := range Condition {
		// 统计型标签
		if attrTypeV, ok := attrType[cond.Key]; ok {
			isAppend := true
			param := &daoAnchorV1.AttrReq{
				AttrId: attrTypeV,
			}
			if len(cond.Condition) <= 0 || cond.Condition[0] == nil {
				continue
			}

			if cond.Key == _onlineType {
				switch cond.Condition[0].StringV {
				case "current":
					param.AttrSubId = 1
				case "last7day":
					param.AttrSubId = 2
				case "last30day":
					param.AttrSubId = 3
				}
			}

			if cond.Key == _incomeType || cond.Key == _dmsType {
				switch cond.Condition[0].StringV {
				case "range15min":
					param.AttrSubId = 1
				case "range30min":
					param.AttrSubId = 2
				case "range45min":
					param.AttrSubId = 3
				case "range60min":
					param.AttrSubId = 4
				}
			}

			if cond.Key == _liveDaysType {
				switch cond.Condition[0].StringV {
				case "last7day":
					param.AttrSubId = 3
				case "last30day":
					param.AttrSubId = 4
				}
			}

			if cond.Key == _hourRankType {
				if cond.Condition[0].TopV <= 0 {
					isAppend = false
				}
				param.AttrSubId = 1 // 目前小时榜固定是1 @orca
			}

			// 排除后台选择了小时榜，但是top写了0，导致dao anchor不返回这个attr，导致下面and or 过滤出错
			if isAppend {
				attrsFilter[attrTypeV] = make(map[int64]*attrFilter)
				attrsFilter[attrTypeV][param.AttrSubId] = &attrFilter{}
				if i := s.getConditionTypeIndex(cond.Condition, _confTypeRange); i >= 0 && cond.Condition[i] != nil {
					attrsFilter[attrTypeV][param.AttrSubId].max = cond.Condition[i].Max
					attrsFilter[attrTypeV][param.AttrSubId].min = cond.Condition[i].Min
				}
				if i := s.getConditionTypeIndex(cond.Condition, _confTypeTop); i >= 0 && cond.Condition[i] != nil {
					attrsFilter[attrTypeV][param.AttrSubId].top = cond.Condition[i].TopV
				}
				attrs = append(attrs, param)
			}
		}
		// 展示型标签
		if _, ok := tagType[cond.Key]; ok {
			if cond.Key == _areaType {
				for index, areaCond := range cond.Condition {
					// parent area id dimension
					if index == 0 {
						err := json.Unmarshal([]byte(areaCond.StringV), &ParentAreaIds)
						if err != nil {
							log.Error("[genCondConfRoomList]recId:%d, unmarshalParentAreaIdErr:%+v", recId, err)
							continue
						}
						if len(ParentAreaIds) == 1 && ParentAreaIds[0] == _isAll {
							areaFilter[_isAll] = make(map[int64]bool)
							break
						}

						for _, pId := range ParentAreaIds {
							if _, ok := areaFilter[pId]; !ok {
								areaFilter[pId] = make(map[int64]bool)
							}
						}
					}
					// area id dimension
					if index == 1 {
						err := json.Unmarshal([]byte(areaCond.StringV), &AreaIds)
						if err != nil {
							log.Error("[genCondConfRoomList]recId:%d, unmarshalAreaIdErr:%+v", recId, err)
							continue
						}
						for pIdIndex, ids := range AreaIds {
							// 后台保证长度对应 ParentAreaId[pIdIndex]
							if len(ParentAreaIds) <= pIdIndex {
								continue
							}
							if _, ok := areaFilter[ParentAreaIds[pIdIndex]]; !ok {
								continue
							}
							for _, id := range ids {
								if _, ok := areaFilter[ParentAreaIds[pIdIndex]][id]; !ok {
									areaFilter[ParentAreaIds[pIdIndex]][id] = true
								}
							}
						}
					}
				}
			}

			if cond.Key == _roomStatusType && len(cond.Condition) >= 0 && cond.Condition[0] != nil {
				switch cond.Condition[0].StringV {
				case "lottery":
					roomStatus[2] = true
				case "pk":
					roomStatus[1] = true
				}
			}

			if cond.Key == _anchorCateType && len(cond.Condition) >= 0 && cond.Condition[0] != nil {
				switch cond.Condition[0].StringV {
				case "normal":
					anchorCate[0] = true
				case "sign":
					anchorCate[1] = true
				case "union":
					anchorCate[2] = true
				}
			}
		}
	}

	log.Info("[genCondConfRoomList]recId:%d, attrs: %+v", recId, attrs)
	roomList, err := s.getOnlineListByAttrs(ctx, attrs)
	if err != nil {
		log.Error("[getOnlineListByAttrs]recId:%d, getOnlineListByAttrsErr:%+v, resp:%+v", recId, err, roomList)
		return
	}

	filterRoomData := make([]*filterItem, 0)
	for _, roomData := range roomList {
		isHit := false
		isTagHit := false
		if Cond == _condAnd {
			isHit = s.andFilter(attrsFilter, areaFilter, roomStatus, anchorCate, roomData, recId)
		}

		if Cond == _condOr {
			isHit, isTagHit = s.orFilter(attrsFilter, areaFilter, roomStatus, anchorCate, roomData, recId)
		}

		if !isHit {
			continue
		}

		// 命中逻辑
		filterRoomData = append(filterRoomData, &filterItem{
			isTagHit: isTagHit,
			item:     roomData,
		})
	}

	return s.sort(attrsFilter, filterRoomData, Cond, recId)
}

func (s *Service) andFilter(attrsFilters map[int64]map[int64]*attrFilter, areaFilter map[int64]map[int64]bool, roomStatus, anchorCate map[int64]bool, roomData *daoAnchorV1.AttrResp, recId int) (isHit bool) {
	// area filter
	if len(areaFilter) > 0 {
		_, isParentAll := areaFilter[_isAll]
		// 不是全选继续判断
		if !isParentAll {
			if _, ok := areaFilter[roomData.ParentAreaId]; !ok {
				log.Info("[andFilter]recId:%d, ParentAreaIdNotMatch:areaFilter:%+v, parentId:%d, roomId:%d", recId, areaFilter, roomData.ParentAreaId, roomData.RoomId)
				return false
			}

			_, isAreaAll := areaFilter[roomData.ParentAreaId][_isAll]
			_, isAreaIdExist := areaFilter[roomData.ParentAreaId][roomData.AreaId]
			if !isAreaAll && !isAreaIdExist {
				log.Info("[andFilter]recId:%d, AreaIdNotMatch:areaFilter:%+v, Id:%d, roomId:%d", recId, areaFilter, roomData.AreaId, roomData.RoomId)
				return false
			}
		}
	}

	//如果配置了房间状态筛选
	if len(roomStatus) > 0 {
		roomStatusTag := &daoAnchorV1.TagData{}
		for _, tag := range roomData.TagList {
			if tag.TagId == _roomStatusTagId {
				roomStatusTag = tag
				break
			}
		}
		if _, ok := roomStatus[roomStatusTag.TagSubId]; !ok {
			log.Info("[andFilter]recId:%d, roomStatusNotMatch:roomStatus:%+v, tagSubId:%d, roomId:%d", recId, roomStatus, roomStatusTag.TagSubId, roomData.RoomId)
			return false
		}
	}

	//如果配置了主播类型筛选
	if len(anchorCate) > 0 {
		if _, ok := anchorCate[roomData.AnchorProfileType]; !ok {
			log.Info("[andFilter]recId:%d, anchorCateNotMatch:anchorCate:%+v, anchorProfileType:%d, roomId:%d", recId, anchorCate, roomData.AnchorProfileType, roomData.RoomId)
			return false
		}
	}

	// attr(统计类) filter
	attrsMap := make(map[int64]*daoAnchorV1.AttrData)
	for _, attr := range roomData.AttrList {
		attrsMap[attr.AttrId] = attr
	}
	for attrId, attrFilter := range attrsFilters {
		if _, ok := attrsMap[attrId]; !ok {
			log.Info("[andFilter]recId:%d, attrNotExist:attrId:%d, attrFilter:%d, roomID:%d", recId, attrId, attrsMap, roomData.RoomId)
			return false
		}
		for attrSubId, attrInfo := range attrFilter {
			if attrInfo == nil {
				continue
			}
			if attrsMap[attrId].AttrSubId == attrSubId {
				if attrInfo.min > 0 && attrsMap[attrId].AttrValue < attrInfo.min {
					log.Info("[andFilter]recId:%d, attrMinNotMatch:attrId:%d, value:%d, max:%d, min:%d, roomID:%d", recId, attrId, attrsMap[attrId].AttrValue, attrInfo.max, attrInfo.min, roomData.RoomId)
					return false
				}

				if attrInfo.max > 0 && attrsMap[attrId].AttrValue > attrInfo.max {
					log.Info("[andFilter]recId:%d, attrMaxNotMatch:attrId:%d, value:%d, max:%d, min:%d, roomID:%d", recId, attrId, attrsMap[attrId].AttrValue, attrInfo.max, attrInfo.min, roomData.RoomId)
					return false
				}
			}
		}
	}
	log.Info("[andFilter]recId:%d, successMatch:roomData:%+v,attrs:%+v, areaFilter:%+v, roomStatus:%+v, anchorState:%+v, roomID:%d", recId, roomData, attrsFilters, areaFilter, roomStatus, anchorCate, roomData.RoomId)

	return true
}

func (s *Service) orFilter(attrsFilters map[int64]map[int64]*attrFilter, areaFilter map[int64]map[int64]bool, roomStatus, anchorCate map[int64]bool, roomData *daoAnchorV1.AttrResp, recId int) (isHit bool, isTagHit bool) {
	// area filter
	isTagHit = false
	if len(areaFilter) > 0 {
		if _, ok := areaFilter[_isAll]; ok {
			log.Info("[orFilter]recId:%d, ParentAreaIdAllMatch:areaFilter:%+v, parentId:%d, roomID:%d", recId, areaFilter, roomData.ParentAreaId, roomData.RoomId)
			isTagHit = true
			return true, isTagHit
		}

		if _, ok := areaFilter[roomData.ParentAreaId]; ok {
			log.Info("[orFilter]recId:%d, ParentAreaIdMatch:areaFilter:%+v, parentId:%d, roomID:%d", recId, areaFilter, roomData.ParentAreaId, roomData.RoomId)
			isTagHit = true
			return true, isTagHit
		}

		_, isAll := areaFilter[roomData.ParentAreaId][_isAll]
		_, isAreaIdExist := areaFilter[roomData.ParentAreaId][roomData.AreaId]
		if isAll || isAreaIdExist {
			log.Info("[orFilter]recId:%d, AreaIdMatch:areaFilter:%+v, Id:%d, roomID:%d", recId, areaFilter, roomData.AreaId, roomData.RoomId)
			isTagHit = true
			return true, isTagHit
		}
	}

	//如果配置了房间状态筛选
	if len(roomStatus) > 0 {
		roomStatusTag := &daoAnchorV1.TagData{}
		for _, tag := range roomData.TagList {
			if tag.TagId == _roomStatusTagId {
				roomStatusTag = tag
				break
			}
		}
		if _, ok := roomStatus[roomStatusTag.TagSubId]; ok {
			log.Info("[orFilter]recId:%d, roomStatusMatch:roomStatus:%+v, tagSubId:%d, roomID:%d", recId, roomStatus, roomStatusTag.TagSubId, roomData.RoomId)
			isTagHit = true
			return true, isTagHit
		}
	}

	//如果配置了主播类型筛选
	if len(anchorCate) > 0 {
		if _, ok := anchorCate[roomData.AnchorProfileType]; ok {
			log.Info("[orFilter]recId:%d, anchorCateMatch:anchorCate:%+v, anchorProfileType:%d, roomID:%d", recId, anchorCate, roomData.AnchorProfileType, roomData.RoomId)
			isTagHit = true
			return true, isTagHit
		}
	}

	// attr filter
	attrDataMap := make(map[int64]*daoAnchorV1.AttrData)
	for _, attr := range roomData.AttrList {
		attrDataMap[attr.AttrId] = attr
	}
	for attrId, attrFilter := range attrsFilters {
		if _, ok := attrDataMap[attrId]; !ok {
			continue
		}

		//小时榜不需要判断max min
		if attrId == attrType[_hourRankType] {
			return true, isTagHit
		}
		for attrSubId, attrFilterInfo := range attrFilter {
			attrData, ok := attrDataMap[attrId]
			if !ok {
				continue
			}
			if attrData.AttrSubId == attrSubId {
				if attrFilterInfo.min > 0 && attrData.AttrValue < attrFilterInfo.min {
					continue
				}
				if attrFilterInfo.max > 0 && attrData.AttrValue > attrFilterInfo.max {
					continue
				}
				log.Info("[orFilter]recId:%d, attrMaxMinMatch:attrId:%d, value:%d, max:%d, min:%d, roomID:%d", recId, attrId, attrDataMap[attrId].AttrValue, attrFilterInfo.max, attrFilterInfo.min, roomData.RoomId)

				return true, isTagHit
			}
		}
	}

	log.Info("[orFilter]recId:%d, failMatch:roomData:%+v,attrs:%+v, areaFilter:%+v, roomStatus:%+v, anchorState:%+v, roomId:%d", recId, roomData, attrsFilters, areaFilter, roomStatus, anchorCate, roomData.RoomId)
	return false, isTagHit
}
