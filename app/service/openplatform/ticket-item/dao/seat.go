package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// 管理area_seats、area_seatmap、seat_set、seat_order座位相关表
// area_seats为区域座位表，每个座位一行
// area_seatmap为区域座位图表，每个区域一行，与area_seats对应
// seat_order为场次的座位订单表，每个座位一行，创建场次座位图时生成
// seat_set为场次的座位设置图表，每个区域一行，基于area_seatmap生成，针对场次含有不同的票价标记

const (
	// StatusCansale 可售
	StatusCansale = 0
	// StatusIssue 已出票
	StatusIssue = 2
	// StatusLocked 已锁定
	StatusLocked = 3
	// StatusBooked 已预订
	StatusBooked = 4
)

// TxUpdateSeat 编辑区域座位信息（事务）
func (d *Dao) TxUpdateSeat(c context.Context, tx *gorm.DB, area *model.Area) (err error) {
	if err = tx.Table("area").Where("id = ?", area.ID).Updates(
		map[string]interface{}{
			"seats_num":      area.SeatsNum,
			"width":          area.Width,
			"height":         area.Height,
			"deleted_status": area.DeletedStatus,
			"col_start":      area.ColStart,
			"col_type":       area.ColType,
			"col_direction":  area.ColDirection,
			"row_list":       area.RowList,
			"seat_start":     area.SeatStart,
		}).Error; err != nil {
		log.Error("更新区域座位信息（ID:%d）失败:%s", area.ID, err)
		err = ecode.NotModified
		return
	}
	return
}

// TxGetAreaSeats 批量获取区域对应的区域座位信息（事务）
func (d *Dao) TxGetAreaSeats(c context.Context, tx *gorm.DB, area int64) (areaSeats []*model.AreaSeats, err error) {
	if err = tx.Find(&areaSeats, "area = ?", area).Error; err != nil {
		log.Error("批量获取区域座位信息（area:%d）失败:%s", area, err)
		return
	}
	return
}

// TxBatchAddAreaSeats 批量添加区域座位信息（事务）
func (d *Dao) TxBatchAddAreaSeats(c context.Context, tx *gorm.DB, areaSeats []*model.AreaSeats) (err error) {
	if len(areaSeats) == 0 {
		return
	}
	var values = make([]string, len(areaSeats))
	for i, areaSeat := range areaSeats {
		values[i] = fmt.Sprintf("(%d,%d,'%s','%s',%d,%d)", areaSeat.X, areaSeat.Y, areaSeat.Label, areaSeat.Bgcolor, areaSeat.Area, 0)
	}
	var sql = fmt.Sprintf("INSERT INTO `area_seats` (`x`, `y`, `label`, `bgcolor`, `area`, `dstatus`) VALUES %s;", strings.Join(values, ","))
	if err = tx.Exec(sql).Error; err != nil {
		log.Error("批量添加区域座位信息（%s）失败:%s", sql, err)
		err = ecode.NotModified
		return
	}
	return
}

// TxBatchDeleteAreaSeats 软删除区域对应的区域座位表信息
func (d *Dao) TxBatchDeleteAreaSeats(c context.Context, tx *gorm.DB, area int64) (err error) {
	if err = tx.Table("area_seats").Where("area = ?", area).Updates(
		map[string]interface{}{
			"dstatus": 1,
		}).Error; err != nil {
		log.Error("删除区域座位信息（area:%d）失败:%s", area, err)
		err = ecode.NotModified
		return
	}
	return
}

// TxBatchRecoverAreaSeats 恢复软删除的区域座位表信息
func (d *Dao) TxBatchRecoverAreaSeats(c context.Context, tx *gorm.DB, ids []int64) (err error) {
	if err = tx.Table("area_seats").Where("id in (?)", ids).Updates(
		map[string]interface{}{
			"dstatus": 0,
		}).Error; err != nil {
		log.Error("批量恢复区域座位信息失败:%s", err)
		err = ecode.NotModified
		return
	}
	return
}

// TxRawAreaSeatmap 获取区域座位图信息（事务）
func (d *Dao) TxRawAreaSeatmap(c context.Context, tx *gorm.DB, id int64) (areaSeatmap *model.AreaSeatmap, err error) {
	areaSeatmap = new(model.AreaSeatmap)
	if err = tx.First(&areaSeatmap, id).Error; err != nil {
		log.Error("获取区域座位信息（ID:%d）失败:%s", id, err)
		return
	}
	return
}

// TxSaveAreaSeatmap 添加/修改区域座位图信息（事务）
func (d *Dao) TxSaveAreaSeatmap(c context.Context, tx *gorm.DB, areaSeatmap *model.AreaSeatmap) (err error) {
	if res := tx.Save(areaSeatmap); res.Error != nil {
		log.Error("添加区域座位信息失败:%s", res.Error)
		err = ecode.NotModified
		return
	}
	return
}

// TxGetSeatChart 根据场次ID和区域ID查询ID和票价设置图（事务）
func (d *Dao) TxGetSeatChart(c context.Context, tx *gorm.DB, screen int64, area int64) (seatSet *model.SeatSet, err error) {
	seatSet = new(model.SeatSet)
	if res := tx.Select("id, seat_chart").Where("screen_id = ? AND area_id = ? AND deleted_at = 0", screen, area).First(seatSet); res.Error != nil {
		if res.RecordNotFound() {
			return
		}
		err = res.Error
		log.Error("TxGetSeatChart error(%v)", err)
	}
	return
}

// TxGetSeatCharts 根据场次ID和多个区域ID批量查询多个票价设置ID和票价设置图（事务）
func (d *Dao) TxGetSeatCharts(c context.Context, tx *gorm.DB, screen int64, areas []int64) (seatSets []*model.SeatSet, err error) {
	if err = tx.Select("id, seat_chart").Where("screen_id = ? AND area_id in (?) AND deleted_at = 0", screen, areas).Find(&seatSets).Error; err != nil {
		log.Error("TxGetSeatCharts error(%v)", err)
	}
	return
}

// TxGetSeatSets 根据区域ID批量查询多个票价设置ID和场次ID（事务）
func (d *Dao) TxGetSeatSets(c context.Context, tx *gorm.DB, area int64) (seatSets []*model.SeatSet, err error) {
	if err = tx.Select("id, screen_id").Where("area_id = ? AND deleted_at = 0", area).Find(&seatSets).Error; err != nil {
		log.Error("TxGetSeatSets error(%v)", err)
	}
	return
}

// TxAddSeatChart 添加票价设置图（事务）
func (d *Dao) TxAddSeatChart(c context.Context, tx *gorm.DB, seatSet *model.SeatSet) (err error) {
	if res := tx.Create(seatSet); res.Error != nil {
		log.Error("添加票价设置图失败:%s", res.Error)
		err = ecode.NotModified
		return
	}
	return
}

// TxUpdateSeatChart 更新票价设置图（事务）
func (d *Dao) TxUpdateSeatChart(c context.Context, tx *gorm.DB, id int64, seatChart string) (err error) {
	if err = tx.Table("seat_set").Where("id = ? AND deleted_at = 0", id).Updates(
		map[string]interface{}{
			"seat_chart": seatChart,
		}).Error; err != nil {
		log.Error("更新票价设置图（ID:%d）失败:%s", id, err)
	}
	return
}

// TxClearSeatCharts 清空票价设置图（事务）
func (d *Dao) TxClearSeatCharts(c context.Context, tx *gorm.DB, ids []int64) (err error) {
	if err = tx.Table("seat_set").Where("id IN (?) AND deleted_at = 0", ids).Updates(
		map[string]interface{}{
			"seat_chart": "",
		}).Error; err != nil {
		log.Error("清空票价设置图失败:%s", err)
	}
	return
}

// TxGetUnsaleableSeatOrders 根据场次和区域ID查询不可售座位订单信息（事务）
func (d *Dao) TxGetUnsaleableSeatOrders(c context.Context, tx *gorm.DB, screen int64, area int64) (seatOrders []*model.SeatOrder, err error) {
	if err = tx.Unscoped().Where("screen_id = ? AND area_id = ? AND status in (?) AND deleted_at = 0", screen, area, []int32{StatusIssue, StatusLocked, StatusBooked}).Find(&seatOrders).Error; err != nil {
		log.Error("TxGetUnsaleableSeatOrders error(%v)", err)
	}
	return
}

// TxGetSaleableSeatOrders 根据场次ID和票价ID查询可售座位订单ID和区域ID信息（事务）
func (d *Dao) TxGetSaleableSeatOrders(c context.Context, tx *gorm.DB, screen int64, price int64) (seatOrders []*model.SeatOrder, err error) {
	if err = tx.Select("id, area_id").Unscoped().Where("screen_id = ? AND price_id = ? AND status = ? AND deleted_at = 0", screen, price, StatusCansale).Find(&seatOrders).Error; err != nil {
		log.Error("TxGetSaleableSeatOrders error(%v)", err)
	}
	return
}

// TxBatchDelUnsoldSeatOrders 根据区域ID清空未售出的座位订单（事务）
func (d *Dao) TxBatchDelUnsoldSeatOrders(c context.Context, tx *gorm.DB, area int64) (err error) {
	if err = tx.Table("seat_order").Where("area_id = ? AND status IN (?) AND deleted_at = 0", area, []int32{StatusCansale, StatusLocked}).Updates(
		map[string]interface{}{
			"deleted_at": time.Now(),
		}).Error; err != nil {
		log.Error("批量删除座位订单信息失败:%s", err)
		err = ecode.NotModified
		return
	}
	return
}

// TxAddSeatOrder 添加座位订单信息（事务，暂未使用）
func (d *Dao) TxAddSeatOrder(c context.Context, tx *gorm.DB, seatOrder *model.SeatOrder) (err error) {
	if res := tx.Create(seatOrder); res.Error != nil {
		log.Error("添加座位订单信息失败:%s", res.Error)
		err = ecode.NotModified
		return
	}
	return
}

// TxUpdateSeatOrder 编辑座位订单信息（事务，暂未使用）
// TODO: 具体字段未指定
func (d *Dao) TxUpdateSeatOrder(c context.Context, tx *gorm.DB, seatOrder *model.SeatOrder) (err error) {
	if err = tx.Table("seat_order").Where("id = ? AND deleted_at = 0", seatOrder.ID).Updates(
		map[string]interface{}{}).Error; err != nil {
		log.Error("更新座位订单信息（ID:%d）失败:%s", seatOrder.ID, err)
		err = ecode.NotModified
		return
	}
	return
}

// TxBatchDeleteSeatOrder 批量软删除座位订单信息（事务）
func (d *Dao) TxBatchDeleteSeatOrder(c context.Context, tx *gorm.DB, ids []int64) (err error) {
	if err = tx.Table("seat_order").Where("id in (?)", ids).Updates(
		map[string]interface{}{
			"deleted_at": time.Now(),
		}).Error; err != nil {
		log.Error("批量删除座位订单信息失败:%s", err)
		err = ecode.NotModified
		return
	}
	return
}

// TxBatchAddSeatOrder 批量添加座位订单（事务）
func (d *Dao) TxBatchAddSeatOrder(c context.Context, tx *gorm.DB, seatOrders []*model.SeatOrder) (err error) {
	if len(seatOrders) == 0 {
		return
	}
	var values = make([]string, len(seatOrders))
	for i, so := range seatOrders {
		values[i] = fmt.Sprintf("(%d,%d,%d,%d,%d,%d)", so.AreaID, so.ScreenID, so.Row, so.Col, so.PriceID, so.Price)
	}
	var sql = fmt.Sprintf("INSERT INTO `seat_order` (`area_id`, `screen_id`, `row`, `col`, `price_id`, `price`) VALUES %s;", strings.Join(values, ","))
	if err = tx.Exec(sql).Error; err != nil {
		log.Error("批量添加区域座位信息（%s）失败:%s", sql, err)
		err = ecode.NotModified
		return
	}
	return
}
