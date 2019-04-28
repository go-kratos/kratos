package model

import (
	"fmt"
	"strings"

	"go-common/library/time"
)

//Order...主订单表状态
const (
	OrderPaid         = 2
	OrderRefunded     = 3
	OrderRefundPartly = 2
	OrderRefundedAll  = 4
)

//DistOrder...分销订单表状态
const (
	DistOrderNormal         = 1
	DistOrderRefunded       = 2
	DistOrderPartlyRefunded = 3
)

// SpecsSeparator sku 规格分隔符
const SpecsSeparator = "_"

//OrderInfo 订单同步字段
type OrderInfo struct {
	Oid       uint64    `json:"oid"`
	CmAmount  uint64    `json:"cm_amount"`
	CmMethod  int64     `json:"cm_method"`
	CmPrice   uint64    `json:"cm_price"`
	Duid      uint64    `json:"duid"`
	Stat      int64     `json:"status"`
	Pid       uint64    `json:"pid"`
	Count     uint64    `json:"count"`
	Sid       uint64    `json:"sid"`
	Type      int64     `json:"type"`
	PayAmount uint64    `json:"pay_amount"`
	Serial    string    `json:"serial_num"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

//OrderStockCnt 订单库存
type OrderStockCnt struct {
	OrderID int64
	Count   int64
}

//SkuCnt sku库存
type SkuCnt struct {
	SkuID int64
	Count int64
}

//Batch 批量库存操作
type Batch struct {
	ScreenID      int64 // 场次 ID
	TicketPriceID int64 // 票价 ID 就是 sku_stock SKUID
	SkAlert       int64 // 库存预警数
	TotalStock    int64 // 总库存数
}

//Specs 生成库存规格
func (b *Batch) Specs() (s string) {
	s = fmt.Sprintf("%d%s%d", b.ScreenID, SpecsSeparator, b.TicketPriceID)
	return
}

//InsPlHlds insert语句占位符
func InsPlHlds(colCnt int, rowCnt int) string {
	hLen := colCnt*2 + 2
	h := make([]byte, hLen)
	h[0] = '('
	h[hLen-2] = ')'
	h[hLen-1] = ','
	copy(h[1:], []byte(strings.Repeat(",?", colCnt)[1:]))
	bPlHlds := make([]byte, hLen*rowCnt-1)
	j := 0
	for i := 0; i < rowCnt; i++ {
		copy(bPlHlds[j:], h)
		j += hLen
	}
	return string(bPlHlds)
}
