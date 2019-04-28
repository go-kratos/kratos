package model

import (
	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"strconv"
)

// TicketInfo 票价综合字段
type TicketInfo struct {
	TicketPrice
	BuyNumLimit map[string]*TicketPriceExtra
}

// FormatTicketBuyLimit 格式化票价购票限制
func (t *TicketInfo) FormatTicketBuyLimit(limit *item.TicketBuyNumLimit) {
	limit.Normal = make(map[int32]int64)
	limit.Vip = make(map[int32]int64)
	limit.AnnualVip = make(map[int32]int64)
	if ext, ok := t.BuyNumLimit[TkBuyNumLimitNormal]; ok {
		limit.Normal = ext.ParseBuyLimit()
	}
	if ext, ok := t.BuyNumLimit[TkBuyNumLimitVip]; ok {
		limit.Vip = ext.ParseBuyLimit()
	}

	if ext, ok := t.BuyNumLimit[TkBuyNumLimitAnnualVip]; ok {
		limit.AnnualVip = ext.ParseBuyLimit()
	}
}

// ParseBuyLimit parse 购票限制成map
func (ext *TicketPriceExtra) ParseBuyLimit() (m map[int32]int64) {
	var (
		i   int32
		max int32
		l   int32
	)
	l = 2
	max = 6
	r := []rune(ext.Value)
	m = make(map[int32]int64)
	for i = 0; i < max+1; i++ {
		m[i] = ext.SliceBuyLimit(r, i*l, (i+1)*l)
	}
	return m
}

// SliceBuyLimit 分割等级购票限制
func (ext *TicketPriceExtra) SliceBuyLimit(r []rune, start int32, end int32) int64 {
	slice := string(r[start:end])
	if i, err := strconv.ParseInt(slice, 10, 64); err == nil {
		return i
	} else if slice == "**" {
		return -1
	} else {
		return 0
	}
}
