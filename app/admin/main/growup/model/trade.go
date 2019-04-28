package model

import (
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/growup/util"
)

// GoodsInfo .
type GoodsInfo struct {
	// internal
	ID            int64         `json:"id"`
	ProductID     string        `json:"product_id"`
	ResourceID    int64         `json:"-"`
	GoodsType     GoodsType     `json:"-"`
	Discount      int           `json:"discount"`
	IsDisplay     DisplayStatus `json:"is_display"`
	DisplayOnTime time.Time     `json:"-"`
	// derived
	GoodsTypeDesc string `json:"goods_type"` // 商品类型描述
	// external
	ProductName  string `json:"product_name"`  // 商品名称
	OriginPrice  int64  `json:"origin_price"`  // 实时成本, 单位分
	CurrentPrice int64  `json:"current_price"` // 实时售价, 单位分
	Month        int32  `json:"month"`         //有效期
}

// MergeExternal information from src to target
func MergeExternal(target *GoodsInfo, src *GoodsInfo) error {
	switch target.GoodsType {
	case GoodsVIP:
		target.OriginPrice = src.OriginPrice
		target.ProductName = src.ProductName
		target.CurrentPrice = int64(util.DivWithRound(float64(target.OriginPrice*int64(target.Discount)), 100, 0))
		target.Month = src.Month
		return nil
	default:
		return fmt.Errorf("illegal type of goods(%v)", target)
	}
}

// OrderInfo .
type OrderInfo struct {
	// internal
	ID         int64     `json:"-"`
	MID        int64     `json:"mid"`
	OrderNo    string    `json:"order_no"`
	OrderTime  time.Time `json:"-"`
	GoodsType  GoodsType `json:"-"`
	GoodsID    string    `json:"goods_id"`
	GoodsName  string    `json:"goods_name"`
	GoodsPrice int64     `json:"goods_price"`
	GoodsCost  int64     `json:"goods_cost"`
	// desc for front end
	GoodsTypeDesc string `json:"goods_type"` // 商品类型描述
	OrderTimeDesc string `json:"order_time"` // 订单时间
	// derived
	TotalPrice int64 `json:"total_price"`
	TotalCost  int64 `json:"total_cost"`
	GoodsNum   int64 `json:"goods_num"`
	// external
	Nickname string `json:"nickname"`
}

// OrderExportFields .
func OrderExportFields() []string {
	return []string{"订单ID", "时间", "商品ID", "商品名称", "售价", "成本", "数量", "总实收", "总成本", "UID", "昵称"}
}

// ExportStrings .
func (v *OrderInfo) ExportStrings() []string {
	return []string{
		v.OrderNo,
		v.OrderTimeDesc,
		v.GoodsID,
		v.GoodsName,
		strconv.FormatFloat(util.Div(float64(v.GoodsPrice), float64(100)), 'f', 2, 64),
		strconv.FormatFloat(util.Div(float64(v.GoodsCost), float64(100)), 'f', 2, 64),
		strconv.FormatInt(v.GoodsNum, 10),
		strconv.FormatFloat(util.Div(float64(v.TotalPrice), float64(100)), 'f', 2, 64),
		strconv.FormatFloat(util.Div(float64(v.TotalCost), float64(100)), 'f', 2, 64),
		strconv.FormatInt(v.MID, 10),
		v.Nickname,
	}
}

// GenDerived generates derived information
func (v *OrderInfo) GenDerived() *OrderInfo {
	v.GoodsNum = 1
	v.TotalPrice = v.GoodsPrice
	v.TotalCost = v.GoodsCost
	return v
}

// GenDesc generates descriptions
func (v *OrderInfo) GenDesc() *OrderInfo {
	v.GoodsTypeDesc = v.GoodsType.Desc()
	v.OrderTimeDesc = v.OrderTime.Format("2006-01-02 15:04:05")
	return v
}

// DisplayStatus .
type DisplayStatus int

// DisplayStatuses enum
const (
	DisplayOff DisplayStatus = 1
	DisplayOn  DisplayStatus = 2
)

// GoodsType .
type GoodsType int

// GoodsTypes enum
const (
	GoodsVIP GoodsType = 1
)

// Desc of GoodsType
func (t GoodsType) Desc() string {
	switch t {
	case GoodsVIP:
		return "大会员"
	default:
		return "未定义商品类型 " + string(t)
	}
}

// TimeType .
type TimeType int

// TimeTypes enum
const (
	Daily TimeType = 1 + iota
	Weekly
	Monthly
)

// RangeStart returns the included startTime
func (t TimeType) RangeStart(date time.Time) time.Time {
	if t == Weekly {
		n := int(date.Weekday() - time.Monday)
		if n < 0 {
			n += 7
		}
		return time.Date(date.Year(), date.Month(), date.Day()-n, 0, 0, 0, 0, time.Local)
	} else if t == Monthly {
		return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	}
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
}

// RangeEnd returns the excluded endTime
func (t TimeType) RangeEnd(date time.Time) time.Time {
	if t == Weekly {
		n := int(time.Monday - date.Weekday())
		if n <= 0 {
			n += 7
		}
		return time.Date(date.Year(), date.Month(), date.Day()+n, 0, 0, 0, 0, time.Local)
	} else if t == Monthly {
		return time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.Local)
	} else if t == Daily {
		return time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, time.Local)
	}
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
}

// RangeDesc .
func (t TimeType) RangeDesc(start time.Time, end time.Time) string {
	if t == Daily {
		return start.Format("2006-01-02")
	}
	return start.Format("2006-01-02") + "~" + end.AddDate(0, 0, -1).Format("2006-01-02")
}

// Next returns time on next range
func (t TimeType) Next() func(time.Time) time.Time {
	return func(start time.Time) time.Time {
		switch t {
		case Daily:
			return start.AddDate(0, 0, 1)
		case Weekly:
			return start.AddDate(0, 0, 7)
		case Monthly:
			return start.AddDate(0, 1, 0)
		default:
			return start.AddDate(0, 0, 1)
		}
	}
}

// OrderQueryArg .
type OrderQueryArg struct {
	TimeType  TimeType `form:"time_type" default:"1"`
	FromTime  int64    `form:"from_time" validate:"required,min=1"`
	ToTime    int64    `form:"to_time" validate:"required,min=1"`
	GoodsType int      `form:"goods_type"`
	GoodsID   string   `form:"goods_id"`
	GoodsName string   `form:"goods_name"`
	OrderNO   string   `form:"order_no"`
	MID       int64    `form:"mid"`
	Nickname  string   `form:"nickname"`
	From      int      `form:"from" validate:"min=0" default:"0"`
	Limit     int      `form:"limit" validate:"min=1" default:"20"`
	// fromTime + toTime + timeType => (included) startTime & (excluded) endTime
	StartTime time.Time `form:"-"`
	EndTime   time.Time `form:"-"`
}
