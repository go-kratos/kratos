package model

import "go-common/library/time"

// SKUStockLog sku 库存操作日志
type SKUStockLog struct {
	ID         int64     `db:"id"`
	SKUID      int64     `db:"sku_id"`
	OpType     int16     `db:"op_type"`
	SrcID      int64     `db:"src_id"`
	Stock      int64     `db:"stock"`
	CanceledAt time.Time `db:"canceled_at"`
	Ctime      time.Time `db:"ctime,autogen"`
	Mtime      time.Time `db:"mtime,autogen"`
}

// SKUStock sku 库存
type SKUStock struct {
	SKUID       int64     `db:"sku_id"`
	ParentSKUID int64     `db:"parent_sku_id"`
	ItemID      int64     `db:"item_id"`
	Specs       string    `db:"specs"`
	TotalStock  int64     `db:"total_stock"`
	Stock       int64     `db:"stock"`
	LockedStock int64     `db:"locked_stock"`
	SkAlert     int64     `db:"sk_alert"`
	Ctime       time.Time `db:"ctime,autogen"`
	Mtime       time.Time `db:"mtime,autogen"`
}
