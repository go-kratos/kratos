package model

import (
	"time"
)

// Asset .
type Asset struct {
	ID       int64
	MID      int64
	OID      int64
	OType    string
	Currency string
	Price    int64
	State    string
	CTime    time.Time
	MTime    time.Time
}

// PickPrice 获得平台相应的价格
func (a Asset) PickPrice(platform string, pp map[string]int64) (price int64) {
	var ok bool
	if price, ok = pp[platform]; ok {
		return
	}
	return a.Price
}

// AssetRelation .
type AssetRelation struct {
	ID    int64
	OID   int64
	OType string
	MID   int64
	State string
	CTime time.Time
	MTime time.Time
}
