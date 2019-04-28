package service

import (
	"context"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"

	"github.com/jinzhu/gorm"
)

// AreaInfo 添加/修改区域信息
func (s *ItemService) AreaInfo(c context.Context, info *item.AreaInfoRequest) (res *item.AreaInfoReply, err error) {
	var (
		oriArea *model.Area
		oldArea *model.Area
		tx      *gorm.DB
	)
	oriArea = &model.Area{
		ID:            info.ID,
		AID:           info.AID,
		Name:          info.Name,
		Place:         info.Place,
		DeletedStatus: 0,
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
		// 查询是否存在被删除的旧区域
		if oldArea, err = s.dao.TxRawDeletedAreaByAID(c, tx, oriArea.AID, oriArea.Place); err != nil {
			tx.Rollback()
			return
		}
		if oldArea == nil {
			// 添加主体信息
			if err = s.dao.TxAddArea(c, tx, oriArea); err != nil {
				tx.Rollback()
				return
			}
		} else {
			// 更新主体信息
			oriArea.ID = oldArea.ID
			if err = s.dao.TxUpdateArea(c, tx, oriArea); err != nil {
				tx.Rollback()
				return
			}
		}
	} else {
		// 查询并保存旧的区域信息
		if oldArea, err = s.dao.TxRawArea(c, tx, oriArea.ID); err != nil {
			tx.Rollback()
			return
		}
		// 更新主体信息
		if err = s.dao.TxUpdateArea(c, tx, oriArea); err != nil {
			tx.Rollback()
			return
		}
		// 删除旧场地的区域位置
		if oldArea.Place != info.Place {
			if err = s.dao.TxDelAreaPolygon(c, tx, oriArea.Place, info.ID); err != nil {
				tx.Rollback()
				return
			}

		}
	}
	// 更新新场地的区域位置
	if err = s.dao.TxAddOrUpdateAreaPolygon(c, tx, info.Place, oriArea.ID, &info.Coordinate); err != nil {
		tx.Rollback()
		return
	}
	// 提交事务
	if err = s.dao.CommitTran(c, tx); err != nil {
		err = ecode.NotModified
		return
	}
	res = &item.AreaInfoReply{
		Success:    true,
		ID:         oriArea.ID,
		Coordinate: info.Coordinate,
	}
	return
}

// DeleteArea 软删除区域信息
func (s *ItemService) DeleteArea(c context.Context, info *item.DeleteAreaRequest) (res *item.DeleteAreaReply, err error) {
	var (
		tx       *gorm.DB
		area     *model.Area
		seatSets []*model.SeatSet
	)
	if err = v.Struct(info); err != nil {
		err = ecode.RequestErr
		return
	}
	// 开启事务
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = ecode.NotModified
		return
	}
	if area, err = s.dao.TxRawArea(c, tx, info.ID); err != nil {
		tx.Rollback()
		return
	}
	// 清除区域座位图信息（area_seatmap表）
	if err = s.dao.TxSaveAreaSeatmap(c, tx, &model.AreaSeatmap{
		ID:      info.ID,
		SeatMap: "",
	}); err != nil {
		tx.Rollback()
		return
	}
	// 清除区域座位信息（area_seats表）
	if err = s.dao.TxBatchDeleteAreaSeats(c, tx, info.ID); err != nil {
		tx.Rollback()
		return
	}
	// 清空未售出的座位订单（seat_order表）
	if err = s.dao.TxBatchDelUnsoldSeatOrders(c, tx, info.ID); err != nil {
		tx.Rollback()
		return
	}

	if seatSets, err = s.dao.TxGetSeatSets(c, tx, info.ID); err != nil {
		tx.Rollback()
		return
	}
	delIDs := make([]int64, 0)
	for _, ss := range seatSets {
		delIDs = append(delIDs, ss.ID)
	}
	// 清空票价设置图（seat_set表）
	if err = s.dao.TxClearSeatCharts(c, tx, delIDs); err != nil {
		tx.Rollback()
		return
	}
	// 软删除area表
	if err = s.dao.TxDelArea(c, tx, info.ID); err != nil {
		tx.Rollback()
		return
	}
	// 删除区域的场地坐标信息
	if err = s.dao.TxDelAreaPolygon(c, tx, area.Place, info.ID); err != nil {
		tx.Rollback()
		return
	}

	// TODO: 删除所有相关缓存

	// 提交事务
	if err = s.dao.CommitTran(c, tx); err != nil {
		err = ecode.NotModified
		return
	}
	res = &item.DeleteAreaReply{Success: true}
	return
}
