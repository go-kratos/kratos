package model

import (
	"time"
)

// AggrIncomeUser .
type AggrIncomeUser struct {
	ID         int64
	MID        int64
	Currency   string
	PaySuccess int64
	PayError   int64
	TotalIn    int64
	TotalOut   int64
	CTime      time.Time
	MTime      time.Time
}

// AggrIncomeUserAsset .
type AggrIncomeUserAsset struct {
	ID         int64
	MID        int64
	Currency   string
	Ver        int64
	OID        int64
	OType      string
	PaySuccess int64
	PayError   int64
	TotalIn    int64
	TotalOut   int64
	CTime      time.Time
	MTime      time.Time
}

// AggrIncomeUserAssetList .
type AggrIncomeUserAssetList struct {
	MID    int64
	Ver    int64
	Assets []*AggrIncomeUserAsset
	Page   *Page
}

// Page .
type Page struct {
	Num   int64
	Size  int64
	Total int64
}
