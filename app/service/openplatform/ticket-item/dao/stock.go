package dao

import (
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/log"
)

// StockChanged 检查票价下库存是否有变动
func (d *Dao) StockChanged(ids []int64) bool {
	if ids == nil {
		return false
	}

	var stocks []model.Stock
	if err := d.db.Select("total_stock, stock").Where("sku_id IN (?)", ids).Find(&stocks).Error; err != nil {
		log.Error("获取票价库存信息失败:%s", err)
		return true
	}

	for _, v := range stocks {
		if (v.TotalStock - v.Stock) != 0 {
			log.Error("票价存在库存有变动")
			return true
		}
	}
	return false
}
