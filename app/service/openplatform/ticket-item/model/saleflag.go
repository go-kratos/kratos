package model

import (
	"time"
)

// SaleFlag...售卖状态状态
const (
	SaleFlagNotBegin = 1 // 未开售
	SaleFlagBegin    = 2 // 预售中
	SaleFlagEnd      = 3 // 已停售
	SaleFlagNotSale  = 5 // 不可售
	SaleFlagOut      = 4 // 已售罄
	SaleFlagTight    = 6 // 库存紧张
)

// CalTkSaleFlag 计算SaleFlag
func (tk *TicketInfo) CalTkSaleFlag() (flag int32) {

	current := time.Now().Unix()
	if tk.IsSale == 0 {
		flag = SaleFlagNotSale
	} else if int64(tk.SaleStart) > current {
		flag = SaleFlagNotBegin
	} else if int64(tk.SaleEnd) < current {
		flag = SaleFlagEnd
	} else {
		flag = SaleFlagBegin
	}
	return
}
