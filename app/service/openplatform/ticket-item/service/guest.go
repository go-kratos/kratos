package service

import (
	"context"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"
)

// GuestInfo Add or Update Guest
func (s *ItemService) GuestInfo(c context.Context, info *item.GuestInfoRequest) (ret *item.GuestInfoReply, err error) {
	if paramErr := v.Struct(info); paramErr != nil {
		return ret, ecode.RequestErr
	}

	var result bool
	var daoErr error
	if info.ID == 0 {
		result, daoErr = s.dao.AddGuest(c, info)
	} else {
		result, daoErr = s.dao.UpdateGuest(c, info)
	}
	return &item.GuestInfoReply{Success: result}, daoErr
}

// GuestStatus Change Guest Status
func (s *ItemService) GuestStatus(c context.Context, info *item.GuestStatusRequest) (ret *item.GuestInfoReply, err error) {
	if paramErr := v.Struct(info); paramErr != nil {
		return ret, ecode.RequestErr
	}
	result, daoError := s.dao.GuestStatus(c, info.ID, int8(info.Status))

	return &item.GuestInfoReply{Success: result}, daoError

}

// GetGuests 获取单个宾客信息 测试用
/**func (s *ItemService) GetGuests(c context.Context, id *int64) (res []*model.Guest, err error) {
	res, err = s.dao.GetGuests(c, *id)
	return
}**/

// GuestSearch 场馆搜索
func (s *ItemService) GuestSearch(c context.Context, arg *model.GuestSearchParam) (res *model.GuestSearchList, err error) {
	res, err = s.dao.GuestSearch(c, arg)
	return
}
