package service

import (
	"context"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"

	"github.com/jinzhu/gorm"
)

// PlaceInfo 添加/修改场地信息
func (s *ItemService) PlaceInfo(c context.Context, info *item.PlaceInfoRequest) (res *item.PlaceInfoReply, err error) {
	var (
		oriPlace *model.Place
		tx       *gorm.DB
	)
	oriPlace = &model.Place{
		ID:      info.ID,
		Name:    info.Name,
		BasePic: info.BasePic,
		Status:  info.Status,
		Venue:   info.Venue,
		DWidth:  info.DWidth,
		DHeight: info.DHeight,
	}
	if err = v.Struct(info); err != nil {
		err = ecode.RequestErr
		return
	}
	// 开启事务
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = ecode.NotModified
		return
	}
	if info.ID == 0 {
		if err = s.dao.TxAddPlace(c, tx, oriPlace); err != nil {
			tx.Rollback()
			return
		}
	} else {
		var place *model.Place
		if place, err = s.dao.TxRawPlace(c, tx, info.ID); err != nil {
			tx.Rollback()
			return
		}
		if err = s.dao.TxUpdatePlace(c, tx, oriPlace); err != nil {
			tx.Rollback()
			return
		}
		if err = s.dao.TxDecPlaceNum(c, tx, place.Venue); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = s.dao.TxIncPlaceNum(c, tx, info.Venue); err != nil {
		tx.Rollback()
		return
	}
	// 提交事务
	if err = s.dao.CommitTran(c, tx); err != nil {
		err = ecode.NotModified
		return
	}
	res = &item.PlaceInfoReply{
		Success: true,
		ID:      oriPlace.ID,
	}
	return
}
