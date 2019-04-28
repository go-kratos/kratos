package dao

import (
	"go-common/app/service/openplatform/ticket-item/model"
	"time"
)

const (
	// StatusWaitShelf 待上架
	StatusWaitShelf = 1
	// StatusUpShelf 已上架
	StatusUpShelf = 2
)

// HasPromotion 检查场次id或者票价id下是否有未开售,售卖中的待上架和已上架拼团 checkType 1-场次id 2-票价id
func (d *Dao) HasPromotion(ids []int64, checkType int32) bool {
	if ids == nil {
		return false
	}

	status := []int{StatusWaitShelf, StatusUpShelf}
	whereStr := "(begin_time > ? OR (begin_time <= ? AND end_time > ?)) AND status IN (?) AND "

	if checkType == 1 {
		// screen
		whereStr += "extra IN (?)"
	} else {
		// sku
		whereStr += "sku_id IN (?)"
	}

	var count int64
	currTime := time.Now().Unix()
	d.db.Model(&model.Promotion{}).Where(whereStr, currTime, currTime, currTime, status, ids).Count(&count)

	return count > 0
}
