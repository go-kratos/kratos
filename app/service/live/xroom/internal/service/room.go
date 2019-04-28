package service

import (
	"context"

	dav1 "go-common/app/service/live/dao-anchor/api/grpc/v1"
	pb "go-common/app/service/live/xroom/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

type fieldInfo struct {
	fields []string
	attrs  map[string]bool
}

// Model维度Model定义
var (
	attrsMap = map[string][]string{
		"status": {"live_status", "live_screen_type", "live_mark", "lock_time", "lock_status", "hidden_status", "hidden_time", "live_type"},
		"show":   {"short_id", "title", "cover", "tags", "background", "description", "keyframe", "popularity_count", "tag_list", "live_start_time"},
		"area":   {"area_id", "area_name", "parent_area_id", "parent_area_name"},
		"anchor": {"anchor_profile_type", "anchor_data"},
	}
)

// GetMultiple implementation
// 批量根据room_ids获取房间信息
func (s *Service) GetMultiple(ctx context.Context, req *pb.RoomIDsReq) (resp *pb.RoomIDsInfosResp, err error) {
	resp = &pb.RoomIDsInfosResp{
		List: map[int64]*pb.Infos{},
	}

	var fi *fieldInfo
	if fi, err = attrsProcess(req.Attrs); err != nil {
		return
	}

	freq := &dav1.RoomByIDsReq{
		RoomIds:       req.RoomIds,
		Fields:        fi.fields,
		DefaultFields: 0,
	}

	fresp, ferr := s.daClient.FetchRoomByIDs(ctx, freq)
	if ferr != nil {
		log.Error("[GetMultiple]FetchRoomByIDs from dao-anchor err:%+v", ferr.Error())
		return nil, ferr
	}

	for k, v := range fresp.RoomDataSet {
		resp.List[k] = infoSet(v, len(req.RoomIds), fi.attrs)
	}

	return resp, nil
}

// GetMultipleByUids implementation
// 批量根据uids获取房间信息
func (s *Service) GetMultipleByUids(ctx context.Context, req *pb.UIDsReq) (resp *pb.UIDsInfosResp, err error) {
	resp = &pb.UIDsInfosResp{
		List: map[int64]*pb.Infos{},
	}

	var fi *fieldInfo
	if fi, err = attrsProcess(req.Attrs); err != nil {
		return
	}

	freq := &dav1.RoomByIDsReq{
		Uids:          req.Uids,
		Fields:        fi.fields,
		DefaultFields: 0,
	}

	fresp, ferr := s.daClient.FetchRoomByIDs(ctx, freq)
	if ferr != nil {
		log.Error("[GetMultipleByUids]FetchRoomByIDs from dao-anchor err:%+v", ferr.Error())
		return resp, ferr
	}

	for _, v := range fresp.RoomDataSet {
		resp.List[v.Uid] = infoSet(v, len(req.Uids), fi.attrs)
	}
	return
}

// IsAnchor implementation
// 批量根据uids判断是否是主播，如果是返回主播的room_id，否则返回0
func (s *Service) IsAnchor(ctx context.Context, req *pb.IsAnchorUIDsReq) (resp *pb.IsAnchorUIDsResp, err error) {
	resp = &pb.IsAnchorUIDsResp{
		List: map[int64]int64{},
	}

	freq := &dav1.RoomByIDsReq{
		Uids:          req.Uids,
		Fields:        []string{"room_id", "uid"},
		DefaultFields: 1,
	}

	fresp, ferr := s.daClient.FetchRoomByIDs(ctx, freq)
	if ferr != nil {
		log.Error("[IsAnchor] FetchRoomByIDs from dao-anchor err:%+v]", ferr.Error())
		return resp, ferr
	}

	uidMap := make(map[int64]int64)
	for _, v := range fresp.RoomDataSet {
		uidMap[v.Uid] = v.RoomId
	}

	for _, v := range req.Uids {
		if roomID, ok := uidMap[v]; ok {
			resp.List[v] = roomID
			continue
		}
		resp.List[v] = 0
	}
	return
}

func attrsProcess(attrs []string) (fi *fieldInfo, err error) {
	fi = &fieldInfo{
		fields: []string{},
		attrs:  map[string]bool{},
	}

	tidyFields := make([]string, 0)
	fieldsExistMap := make(map[string]bool)

	for _, dimension := range attrs {
		if fields, ok := attrsMap[dimension]; ok {
			for _, field := range fields {
				if _, exist := fieldsExistMap[field]; !exist {
					tidyFields = append(tidyFields, field)
					fieldsExistMap[field] = true
				}
			}
			fi.attrs[dimension] = true
		}
	}

	if len(tidyFields) == 0 {
		return nil, ecode.Error(-400, "attrs error value")
	}
	fi.fields = tidyFields

	return fi, nil
}

func tagList(tag []*dav1.TagData, len int) (pbTags []*pb.TagData) {
	pbTags = make([]*pb.TagData, 0, len)
	for k, v := range tag {
		pbTags[k] = &pb.TagData{
			TagId:    v.GetTagId(),
			TagSubId: v.GetTagSubId(),
			TagValue: v.GetTagValue(),
			TagExt:   v.GetTagExt(),
		}
	}
	return nil
}

func infoSet(rs *dav1.RoomData, len int, attrs map[string]bool) (info *pb.Infos) {
	infos := &pb.Infos{
		RoomId: rs.GetRoomId(),
		Uid:    rs.GetUid(),
	}

	if attrs["status"] {
		infos.Status = &pb.RoomStatusInfo{
			LiveStatus:     rs.GetLiveStatus(),
			LiveScreenType: rs.GetLiveScreenType(),
			LiveMark:       rs.GetLiveMark(),
			LockTime:       rs.GetLockTime(),
			LockStatus:     rs.GetLockStatus(),
			HiddenStatus:   rs.GetHiddenStatus(),
			HiddenTime:     rs.GetHiddenTime(),
			LiveType:       rs.GetLiveType(),
		}
	}

	if attrs["show"] {
		infos.Show = &pb.RoomShowInfo{
			ShortId:         rs.GetShortId(),
			Title:           rs.GetTitle(),
			Cover:           rs.GetCover(),
			Tags:            rs.GetTags(),
			Background:      rs.GetBackground(),
			Description:     rs.GetDescription(),
			Keyframe:        rs.GetKeyframe(),
			PopularityCount: rs.GetPopularityCount(),
			TagList:         tagList(rs.GetTagList(), len),
			LiveStartTime:   rs.GetLiveStartTime(),
		}
	}

	if attrs["area"] {
		infos.Area = &pb.RoomAreaInfo{
			AreaId:         rs.GetAreaId(),
			AreaName:       rs.GetAreaName(),
			ParentAreaId:   rs.GetParentAreaId(),
			ParentAreaName: rs.GetParentAreaName(),
		}
	}

	if attrs["anchor"] {
		infos.Anchor = &pb.RoomAnchorInfo{
			AnchorProfileType: rs.GetAnchorProfileType(),
			AnchorLevel: &pb.AnchorLevel{
				Level:    rs.GetAnchorLevel().GetLevel(),
				Color:    rs.GetAnchorLevel().GetColor(),
				Score:    rs.GetAnchorLevel().GetScore(),
				Left:     rs.GetAnchorLevel().GetLeft(),
				Right:    rs.GetAnchorLevel().GetRight(),
				MaxLevel: rs.GetAnchorLevel().GetMaxLevel(),
			},
		}
	}

	return infos
}
