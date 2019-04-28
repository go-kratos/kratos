package model

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"sort"
	"strconv"

	"go-common/library/log"
	xtime "go-common/library/time"
)

// MsgContent 支付回调 msgContent 字段的结构
// 由于 MsgContent 可能增/减字段故使用 map 兼容
type MsgContent map[string]string

// Charge 支付回调结构体
type Charge struct {
	ID         string
	Paid       bool
	Refunded   bool
	OrderID    int64
	DeviceType int64
	Channel    int64
	Amount     int64
	TimePaid   xtime.Time
}

// ValidSign 签名校验
func (m MsgContent) ValidSign() (ok bool) {
	sign := m["sign"]
	delete(m, "sign")
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	buf := bytes.Buffer{}
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(m[k])
	}

	fmt.Println(buf.String())
	h := md5.New()
	ok = string(h.Sum(buf.Bytes())) == sign
	return
}

// ToCharge 生成 charge 结构
func (m MsgContent) ToCharge() (charge *Charge, err error) {
	charge = &Charge{}
	charge.ID = m["txId"]
	charge.Paid = true
	charge.Refunded = false

	orderID, err := strconv.ParseInt(m["orderId"], 10, 64)
	if err != nil {
		log.Error("MsgContent.ToCharge() strconv.ParseInt(orderId: %s) error(%v)", m["orderId"], err)
		return
	}
	charge.OrderID = orderID

	deviceType, err := strconv.ParseInt(m["deviceType"], 10, 64)
	if err != nil {
		log.Error("MsgContent.ToCharge() strconv.ParseInt(deviceType: %s) error(%v)", m["deviceType"], err)
		return
	}
	charge.DeviceType = deviceType

	payChannel, err := strconv.ParseInt(m["payChannel"], 10, 64)
	if err != nil {
		log.Error("MsgContent.ToCharge() strconv.ParseInt(payChannel: %s) error(%v)", m["payChannel"], err)
		return
	}
	charge.Channel = payChannel

	amount, err := strconv.ParseInt(m["amount"], 10, 64)
	if err != nil {
		log.Error("MsgContent.ToCharge() strconv.ParseInt(amount: %s) error(%v)", m["amount"], err)
		return
	}
	charge.Amount = amount

	if err = charge.TimePaid.Scan(m["orderPayTime"]); err != nil {
		log.Error("MsgContent.ToCharge() TimePaid.Scan(%v) error(%v)", m["orderPayTime"], err)
		return
	}
	return
}
