package model

import (
	"time"

	"go-common/library/log"
)

// LogOrder .
type LogOrder struct {
	ID        int64
	OrderID   string
	FromState string
	ToState   string
	Desc      string
	CTime     time.Time
	MTime     time.Time
}

// Order .
type Order struct {
	ID         int64
	OrderID    string
	MID        int64
	Biz        string
	Platform   string
	OID        int64
	OType      string
	Fee        int64
	RealFee    int64
	Currency   string
	PayID      string
	PayReason  string
	PayTime    time.Time
	RefundTime time.Time
	Version    int64
	State      string // created 已创建, paying 支付中, paid 已支付, failed 支付失败, closed 已关闭, expired 已超时, finished 已完成(支付成功且对账成功)
	CTime      time.Time
	MTime      time.Time
}

// IsPay 是否已支付
func (o Order) IsPay() bool {
	return o.State == OrderStatePaid || o.State == OrderStateSettled
}

// IsRefunding 是否为退款中
func (o Order) IsRefunding() bool {
	return o.State == OrderStateRefunding
}

// ParsePaidTime 解析支付成功时间
func (o *Order) ParsePaidTime(t string) {
	o.PayTime = o.parsePayTime(t)
}

// ParseRefundedTime 解析退款成功时间
func (o *Order) ParseRefundedTime(t string) {
	o.RefundTime = o.parsePayTime(t)
}

func (o Order) parsePayTime(tstr string) (t time.Time) {
	var err error
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", tstr, time.Local); err != nil {
		log.Error("Order parse pay time from timestr: %s, err: %+v", t, err)
		t = time.Now() // 兜底
	}
	return
}

// ReturnState 返回前端显示state，重复下单的退费状态，被认为是已付费
func (o Order) ReturnState() string {
	if o.State == OrderStateDupRefunded {
		return OrderStatePaid
	}
	return o.State
}

// UpdateState 更新状态，return true 更新成功，false 更新失败（无需更新）
func (o *Order) UpdateState(state string) (changed bool) {
	switch o.State { // 当前state
	case OrderStateCreated, OrderStatePaying:
		switch state {
		case OrderStatePaying, OrderStatePaid, OrderStateFailed, OrderStateClosed, OrderStateExpired, OrderStateDupRefunded:
			changed = true
		}
	case OrderStatePaid:
		if state == OrderStateSettled || state == OrderStateRefunding || state == OrderStateDupRefunded || state == OrderStateRefunded || state == OrderStateBadDebt {
			changed = true
		}
	case OrderStateFailed, OrderStateClosed, OrderStateExpired:
		if state == OrderStatePaid {
			changed = true
		}
	case OrderStateRefunding:
		if state == OrderStateRefunded {
			changed = true
		}
	case OrderStateSettled:
		if state == OrderStateSettledRefunding || state == OrderStateSettledRefunded {
			changed = true
		}
	case OrderStateSettledRefunding:
		if state == OrderStateSettledRefunded {
			changed = true
		}
	case OrderStateSettledRefunded:
		if state == OrderStateRefundFinished {
			changed = true
		}
	}

	if changed {
		o.State = state
	} else {
		log.Error("order update state from : %s , to : %s", o.State, state)
	}
	return
}

// OrderRefund .
type OrderRefund struct {
	ID      int64
	OrderID int64
	State   string
	CTime   time.Time
	MTime   time.Time
}

// OrderBadDebt .
type OrderBadDebt struct {
	ID      int64
	OrderID int64
	Type    string
	State   string
	CTime   time.Time
	MTime   time.Time
}
