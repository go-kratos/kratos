package model

import "go-common/library/time"

// VipMonth vip month info.
type VipMonth struct {
	ID        int64     `json:"id"`
	Month     int8      `json:"month"`
	MonthType int8      `json:"month_type"`
	Operator  string    `json:"operator"`
	Status    int8      `json:"status"`
	Mtime     time.Time `json:"mtime"`
}

// VipMonthPrice month price info.
type VipMonthPrice struct {
	ID                 int64     `json:"id" form:"id"`
	MonthID            int8      `json:"month_id" form:"month_id"`
	MonthType          int8      `json:"month_type" form:"month_type"`
	Money              float64   `json:"money" form:"money"`
	FirstDiscountMoney float64   `json:"first_discount_money" form:"first_discount_money"`
	DiscountMoney      float64   `json:"discount_money" form:"discount_money"`
	Selected           int8      `json:"selected" form:"selected"`
	StartTime          time.Time `json:"start_time" form:"start_time"`
	EndTime            time.Time `json:"end_time" form:"end_time"`
	Remark             string    `json:"remark" form:"remark"`
	Operator           string    `json:"operator"`
	Status             int8      `json:"status"`
	Mtime              time.Time `json:"mtime"`
}
