package model

import (
	"encoding/json"

	"go-common/app/common/openplatform/encoding"
	"go-common/app/service/openplatform/ticket-sales/api/grpc/type"
	"go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/library/time"

	"github.com/gogo/protobuf/types"
)

//OrderMain 订单主表结构
type OrderMain struct {
	OrderID      int64                `json:"order_id"`
	UID          string               `json:"uid"`
	OrderType    int16                `json:"order_type"`
	ItemID       int64                `json:"item_id"`
	ItemInfo     *_type.OrderItemInfo `json:"item_info"`
	Count        int64                `json:"count"`
	TotalMoney   int64                `json:"total_money"`
	PayMoney     int64                `json:"pay_money"`
	ExpressFee   int64                `json:"express_fee"`
	PayChannel   int16                `json:"pay_channel"`
	PayTime      int64                `json:"pay_time"`
	Source       string               `json:"source"`
	Status       int16                `json:"status"`
	SubStatus    int16                `json:"sub_status"`
	RefundStatus int16                `json:"refund_status"`
	IsDeleted    int16                `json:"is_deleted"`
	CTime        time.Time            `json:"ctime"`
	MTime        time.Time            `json:"mtime"`
}

//OrderMainQuerier 订单表查询参数
type OrderMainQuerier v1.ListOrdersRequest

// OrderDetail 订单详情表信息
type OrderDetail struct {
	OrderID       int64               `json:"order_id"`
	Buyer         string              `json:"buyer"`
	Tel           string              `json:"tel"`
	PersonalID    string              `json:"personal_id"`
	ExpressCO     string              `json:"express_co"`
	ExpressNO     string              `json:"express_no"`
	ExpressType   int16               `json:"express_type"`
	Remark        string              `json:"remark"`
	DeviceType    int16               `json:"device_type"`
	IP            []byte              `json:"ip"`
	Coupon        *_type.OrderCoupon  `json:"coupon"`
	DeliverDetail *_type.OrderDeliver `json:"deliver_detail"`
	Detail        *_type.OrderExtra   `json:"detail"`
	MSource       string              `json:"msource"`
	CTime         time.Time           `json:"-"`
	MTime         time.Time           `json:"-"`
}

//OrderSKU order_sku表结构
type OrderSKU _type.OrderSKU

//OrderPayCharge 订单支付表结构
type OrderPayCharge _type.OrderPayCharge

//GetFields 获取order_main表所有字段,参数是需要排除的字段
func (o *OrderMain) GetFields(except *types.FieldMask) []string {
	fields := []string{
		"order_id", "uid", "order_type", "item_id", "item_info",
		"count", "total_money", "express_fee", "pay_money", "pay_channel",
		"pay_time", "source", "status", "sub_status", "refund_status",
		"is_deleted", "ctime", "mtime",
	}
	if except != nil {
		lp := len(except.Paths)
		mExcept := make(map[string]bool, lp)
		for _, v := range except.Paths {
			mExcept[v] = true
		}
		res := make([]string, len(fields)-lp)
		i := 0
		for _, v := range fields {
			if ok := mExcept[v]; !ok {
				res[i] = v
				i++
			}
		}
		return res
	}
	return fields
}

//GetFields 获取order_detail字段名称
func (o *OrderDetail) GetFields(except *types.FieldMask) []string {
	fields := []string{
		"order_id", "buyer", "tel", "personal_id", "express_co",
		"express_no", "express_type", "remark", "device_type", "ip",
		"coupon", "deliver_detail", "detail", "msource", "ctime",
		"mtime",
	}
	if except != nil {
		lp := len(except.Paths)
		mExcept := make(map[string]bool, lp)
		for _, v := range except.Paths {
			mExcept[v] = true
		}
		res := make([]string, len(fields)-lp)
		i := 0
		for _, v := range fields {
			if ok := mExcept[v]; !ok {
				res[i] = v
				i++
			}
		}
		return res
	}
	return fields
}

//GetFields 获取order_sku对象的字段名
func (o *OrderSKU) GetFields(except *types.FieldMask) []string {
	fields := []string{
		"order_id", "sku_id", "count", "origin_price", "price",
		"seat_ids", "ticket_type", "discounts", "ctime", "mtime",
	}
	if except != nil {
		lp := len(except.Paths)
		mExcept := make(map[string]bool, lp)
		for _, v := range except.Paths {
			mExcept[v] = true
		}
		res := make([]string, len(fields)-lp)
		i := 0
		for _, v := range fields {
			if ok := mExcept[v]; !ok {
				res[i] = v
				i++
			}
		}
		return res
	}
	return fields
}

//GetFields 获取order_pay_charge对象的字段名
func (o *OrderPayCharge) GetFields(except *types.FieldMask) []string {
	fields := []string{
		"order_id", "charge_id", "channel", "paid", "refunded",
		"ctime", "mtime",
	}
	if except != nil {
		lp := len(except.Paths)
		mExcept := make(map[string]bool, lp)
		for _, v := range except.Paths {
			mExcept[v] = true
		}
		res := make([]string, len(fields)-lp)
		i := 0
		for _, v := range fields {
			if ok := mExcept[v]; !ok {
				res[i] = v
				i++
			}
		}
		return res
	}
	return fields
}

//GetPtrs 获取order_main对象指针
// 如果设置vptr参数,会把struct指针替换成string指针,并在vptr保存原struct指针(as value)和它在返回数组中的下标(as key)
func (o *OrderMain) GetPtrs(fields *types.FieldMask, vptr map[int]interface{}) []interface{} {
	ptrs := map[string]interface{}{
		"order_id":      &o.OrderID,
		"uid":           &o.UID,
		"order_type":    &o.OrderType,
		"item_id":       &o.ItemID,
		"item_info":     &o.ItemInfo,
		"count":         &o.Count,
		"total_money":   &o.TotalMoney,
		"express_fee":   &o.ExpressFee,
		"pay_money":     &o.PayMoney,
		"pay_channel":   &o.PayChannel,
		"pay_time":      &o.PayTime,
		"source":        &o.Source,
		"status":        &o.Status,
		"sub_status":    &o.SubStatus,
		"refund_status": &o.RefundStatus,
		"is_deleted":    &o.IsDeleted,
		"ctime":         &o.CTime,
		"mtime":         &o.MTime,
	}
	if fields == nil {
		fields = &types.FieldMask{Paths: o.GetFields(nil)}
	}
	ret := make([]interface{}, len(fields.Paths))
	i := 0
	for _, f := range fields.Paths {
		if vptr != nil && f == "item_info" {
			var s string
			ret[i] = &s
			vptr[i] = ptrs[f]
		} else {
			ret[i] = ptrs[f]
		}
		i++
	}
	return ret
}

//GetPtrs 获取order_detail对象指针
func (o *OrderDetail) GetPtrs(fields *types.FieldMask, vptr map[int]interface{}) []interface{} {
	ptrs := map[string]interface{}{
		"order_id":       &o.OrderID,
		"buyer":          &o.Buyer,
		"tel":            &o.Tel,
		"personal_id":    &o.PersonalID,
		"express_co":     &o.ExpressCO,
		"express_no":     &o.ExpressNO,
		"express_type":   &o.ExpressType,
		"remark":         &o.Remark,
		"device_type":    &o.DeviceType,
		"ip":             &o.IP,
		"coupon":         &o.Coupon,
		"deliver_detail": &o.DeliverDetail,
		"detail":         &o.Detail,
		"msource":        &o.MSource,
		"ctime":          &o.CTime,
		"mtime":          &o.MTime,
	}
	if fields == nil {
		fields = &types.FieldMask{Paths: o.GetFields(nil)}
	}
	ret := make([]interface{}, len(fields.Paths))
	i := 0
	for _, f := range fields.Paths {
		if vptr != nil && (f == "coupon" || f == "deliver_detail" || f == "detail") {
			var s string
			ret[i] = &s
			vptr[i] = ptrs[f]
		} else {
			ret[i] = ptrs[f]
		}
		i++
	}
	return ret
}

//GetPtrs 获取order_sku对象的字段指针
func (o *OrderSKU) GetPtrs(fields *types.FieldMask, vptr map[int]interface{}) []interface{} {
	ptrs := map[string]interface{}{
		"order_id":     &o.OrderID,
		"sku_id":       &o.SKUID,
		"count":        &o.Count,
		"origin_price": &o.OriginPrice,
		"price":        &o.Price,
		"seat_ids":     &o.SeatIDs,
		"ticket_type":  &o.TicketType,
		"discounts":    &o.Discounts,
		"ctime":        &o.CTime,
		"mtime":        &o.MTime,
	}
	if fields == nil {
		fields = &types.FieldMask{Paths: o.GetFields(nil)}
	}
	ret := make([]interface{}, len(fields.Paths))
	i := 0
	for _, f := range fields.Paths {
		if vptr != nil && (f == "discounts" || f == "seat_ids") {
			var s string
			ret[i] = &s
			vptr[i] = ptrs[f]
		} else {
			ret[i] = ptrs[f]
		}
		i++
	}
	return ret
}

//GetPtrs 获取order_pay_charge对象的字段指针
func (o *OrderPayCharge) GetPtrs(fields *types.FieldMask, vptr map[int]interface{}) []interface{} {
	ptrs := map[string]interface{}{
		"order_id":  &o.OrderID,
		"charge_id": &o.ChargeID,
		"channel":   &o.Channel,
		"paid":      &o.Paid,
		"refunded":  &o.Refunded,
		"ctime":     &o.CTime,
		"mtime":     &o.MTime,
	}
	if fields == nil {
		fields = &types.FieldMask{Paths: o.GetFields(nil)}
	}
	ret := make([]interface{}, len(fields.Paths))
	i := 0
	for _, f := range fields.Paths {
		if vptr != nil && f == "discounts" {
			var s string
			ret[i] = &s
			vptr[i] = ptrs[f]
		} else {
			ret[i] = ptrs[f]
		}
		i++
	}
	return ret
}

//GetVals 获取order_main对象里的值
func (o *OrderMain) GetVals(fields *types.FieldMask, asString bool) []interface{} {
	vals := map[string]interface{}{
		"order_id":      o.OrderID,
		"uid":           o.UID,
		"order_type":    o.OrderType,
		"item_id":       o.ItemID,
		"item_info":     o.ItemInfo,
		"count":         o.Count,
		"total_money":   o.TotalMoney,
		"express_fee":   o.ExpressFee,
		"pay_money":     o.PayMoney,
		"pay_channel":   o.PayChannel,
		"pay_time":      o.PayTime,
		"source":        o.Source,
		"status":        o.Status,
		"sub_status":    o.SubStatus,
		"refund_status": o.RefundStatus,
		"is_deleted":    o.IsDeleted,
		"ctime":         o.CTime,
		"mtime":         o.MTime,
	}
	if fields == nil {
		fields = &types.FieldMask{Paths: o.GetFields(nil)}
	}
	ret := make([]interface{}, len(fields.Paths))
	i := 0
	for _, f := range fields.Paths {
		if asString && f == "item_info" {
			if b, err := json.Marshal(vals[f]); err == nil {
				ret[i] = string(b)
			}
		} else {
			ret[i] = vals[f]
		}
		i++
	}
	return ret
}

//GetVals 获取order_detail对象字段的值
func (o *OrderDetail) GetVals(fields *types.FieldMask, asString bool) []interface{} {
	vals := map[string]interface{}{
		"order_id":       o.OrderID,
		"buyer":          o.Buyer,
		"tel":            o.Tel,
		"personal_id":    o.PersonalID,
		"express_co":     o.ExpressCO,
		"express_no":     o.ExpressNO,
		"express_type":   o.ExpressType,
		"remark":         o.Remark,
		"device_type":    o.DeviceType,
		"ip":             o.IP,
		"coupon":         o.Coupon,
		"deliver_detail": o.DeliverDetail,
		"detail":         o.Detail,
		"msource":        o.MSource,
		"ctime":          o.CTime,
		"mtime":          o.MTime,
	}
	if fields == nil {
		fields = &types.FieldMask{Paths: o.GetFields(nil)}
	}
	ret := make([]interface{}, len(fields.Paths))
	i := 0
	for _, f := range fields.Paths {
		if asString && (f == "coupon" || f == "deliver_detail" || f == "detail") {
			if b, err := json.Marshal(vals[f]); err == nil {
				ret[i] = string(b)
			}
		} else {
			ret[i] = vals[f]
		}
		i++
	}
	return ret
}

//GetVals 获取order_sku字段的值
func (o *OrderSKU) GetVals(fields *types.FieldMask, asString bool) []interface{} {
	vals := map[string]interface{}{
		"order_id":     o.OrderID,
		"sku_id":       o.SKUID,
		"count":        o.Count,
		"origin_price": o.OriginPrice,
		"price":        o.Price,
		"seat_ids":     o.SeatIDs,
		"ticket_type":  o.TicketType,
		"discounts":    o.Discounts,
		"ctime":        o.CTime,
		"mtime":        o.MTime,
	}
	if fields == nil {
		fields = &types.FieldMask{Paths: o.GetFields(nil)}
	}
	ret := make([]interface{}, len(fields.Paths))
	i := 0
	for _, f := range fields.Paths {
		if asString && (f == "discounts" || f == "seat_ids") {
			if b, err := json.Marshal(vals[f]); err == nil {
				ret[i] = string(b)
			}
		} else {
			ret[i] = vals[f]
		}
		i++
	}
	return ret
}

func (o *OrderDetail) getEncryptPtrs() []*string {
	res := make([]*string, 3)
	res[0] = &o.Tel
	res[1] = &o.PersonalID
	if o.DeliverDetail != nil {
		res[2] = &o.DeliverDetail.Tel
	}
	return res
}

//Encrypt 加密order_detail的字段
func (o *OrderDetail) Encrypt(c *encoding.EncryptConfig) {
	for _, p := range o.getEncryptPtrs() {
		if p != nil {
			s, _ := encoding.Encrypt(*p, c)
			*p = s
		}
	}
}

//Decrypt 解密order_detail字段
func (o *OrderDetail) Decrypt(c *encoding.EncryptConfig) {
	for _, p := range o.getEncryptPtrs() {
		if p != nil {
			s, _ := encoding.Decrypt(*p, c)
			*p = s
		}
	}
}

//GetSettleOrdersRequest 获取结算订单请求
type GetSettleOrdersRequest struct {
	Date      string `form:"date" validate:"required"`
	Ref       byte   `form:"ref"`
	ExtParams string `form:"extParams" validate:"omitempty,numeric"`
	PageSize  int    `form:"pagesize"`
}

//SettleOrder 获取结算订单返回
type SettleOrder struct {
	ID              int64     `json:"-"`
	OrderID         int64     `json:"order_id"`
	RefID           int64     `json:"ref_id"`
	RefundApplyTime time.Time `json:"-"`
}

//SettleOrders 获取结算订单返回
type SettleOrders struct {
	Data      []*SettleOrder `json:"data"`
	ExtParams string         `json:"extParams"`
}
