package model

import "time"

// TagUpInfo used when get avid_tagca
type TagUpInfo struct {
	TagID      int64 `json:"id"`
	MID        int64 `json:"mid"`
	Category   int   `json:"category"`
	Business   int   `json:"business"`
	Ratio      int   `json:"ratio"`
	ActivityID int64 `json:"activity_id"`
	IsCommon   int64 `json:"is_common"`
}

// AVInfo used when calculate income
type AVInfo struct {
	MID         int64 `json:"mid"`
	Category    int   `json:"category"`
	Income      int   `json:"income"`
	TotalIncome int   `json:"totalIncome"`
	TaxMoney    int   `json:"tax_money"`
}

// AID used for aid query
type AID struct {
	ID        int64
	AvID      int64
	IncCharge int64
	IsDeleted int
}

// ActivityAVInfo active_id -> avid
type ActivityAVInfo struct {
	ActivityID int64 `json:"mission_id"`
	AVID       int64 `json:"id"`
	MID        int64 `json:"mid"`
	Category   int   `json:"typeid"`
	TagID      int64 `json:"-"`
	Ratio      int   `json:"-"`
}

// TypesInfo category info.
type TypesInfo struct {
	PID int16 `json:"pid"`
	ID  int16 `json:"id"`
}

// TagAvIncome tag av info.
type TagAvIncome struct {
	TagID       int64
	MID         int64
	AVID        int64
	Income      int
	TotalIncome int
	Date        time.Time
}
