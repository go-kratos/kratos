package model

import (
	xtime "go-common/library/time"
)

// PayOrder represents pay order.
type PayOrder struct {
	ID           int32      `json:"id"`             // 订单表自增ID
	OrderNo      string     `json:"order_no"`       // 订单号
	Platform     int8       `json:"platform"`       // 设备平台,1:tv安卓 2:公众号
	OrderType    int8       `json:"order_type"`     // 订单类型0-普通订单 1-自动续费订单
	Mid          int64      `json:"mid"`            // 下单支付的用户mid
	BuyMonths    int8       `json:"buy_months"`     // 购买vip时长
	ProductId    string     `json:"product_id"`     // 产品id
	Money        int32      `json:"money"`          // vip单价，单位分
	Quantity     int32      `json:"quantity"`       // 购买数量
	RefundAmount int32      `json:"refund_amount"`  // 退款金额，单位分
	Status       int8       `json:"status"`         // 订单状态，1.消费中 2.消费成功 3.消费失败
	ThirdTradeNo string     `json:"third_trade_no"` // 第三方订单号（yst订单号）
	PaymentMoney int32      `json:"payment_money"`  // 真正支付金额，单位分
	PaymentType  string     `json:"payment_type"`   // 支付方式：alipay,wechat
	PaymentTime  xtime.Time `json:"payment_time"`   // 支付时间
	Ver          int32      `json:"ver"`            // 版本号，用于乐观锁
	AppChannel   string     `json:"app_channel"`    // 应用渠道
	Token        string
	Ctime        xtime.Time `json:"ctime"` // 创建时间
	Mtime        xtime.Time `json:"mtime"` // 修改时间
}

// CopyFromPayParam copies fiels from pay param.
func (p *PayOrder) CopyFromPayParam(pp *PayParam) {
	p.OrderNo = pp.OrderNo
	p.Quantity = pp.BuyNum
	p.AppChannel = pp.AppChannel
}

// CopyFromPanel copies field from panel.
func (p *PayOrder) CopyFromPanel(panel *PanelPriceConfig) {
	if panel.SubType == 0 {
		p.OrderType = 0
	}
	if panel.SubType == 1 {
		p.OrderType = 1
	}
	p.ProductId = panel.ProductId
	p.Money = panel.Price
	p.BuyMonths = int8(panel.Month * p.Quantity)
	p.PaymentMoney = panel.Price * p.Quantity
}
