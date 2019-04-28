package service

import (
	"context"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"
)

// VenueSearch 场馆搜索
func (s *ItemService) VenueSearch(c context.Context, in *model.VenueSearchParam) (venues *model.VenueSearchList, err error) {
	venues, err = s.dao.VenueSearch(c, in)
	return
}

// VenueInfo 添加/修改场馆信息
func (s *ItemService) VenueInfo(c context.Context, info *item.VenueInfoRequest) (res *item.VenueInfoReply, err error) {
	var oriVenue = &model.Venue{
		ID:            info.ID,
		Name:          info.Name,
		City:          info.City,
		Province:      info.Province,
		District:      info.District,
		AddressDetail: info.AddressDetail,
		Status:        info.Status,
		Traffic:       info.Traffic,
		PlaceNum:      0,
	}

	if err = v.Struct(info); err != nil {
		err = ecode.RequestErr
		return
	}
	if info.ID == 0 {
		err = s.dao.AddVenue(c, oriVenue)
	} else {
		err = s.dao.UpdateVenue(c, oriVenue)
	}
	res = &item.VenueInfoReply{
		Success: true,
		ID:      oriVenue.ID,
	}
	return
}
