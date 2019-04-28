package model

// PayNotifyContent def.
type PayNotifyContent struct {
	TxID         int64  `json:"txId"`
	OrderID      string `json:"orderId"`
	PayAmount    int64  `json:"payAmount"`
	PayChannel   string `json:"payChannel"`
	PayChannelID int32  `json:"payChannelId"`
	PayStatus    string `json:"payStatus"`
	CustomerID   int64  `json:"customerId"`
	ExpiredTime  int64  `json:"expiredTime"`
}

// PayNotifyContentOld .
type PayNotifyContentOld struct {
	TradeNO      string `json:"txId"`
	OrderID      string `json:"orderId"`
	PayAmount    int64  `json:"payAmount"`
	PayChannel   string `json:"payChannel"`
	PayChannelID int32  `json:"payChannelId"`
	PayStatus    string `json:"payStatus"`
	CustomerID   int64  `json:"customerId"`
}

// PayCallBackResult def.
type PayCallBackResult struct {
	TradeNO     string  `json:"trade_no" form:"trade_no"`
	OutTradeNO  string  `json:"out_trade_no" form:"out_trade_no"`
	TradeStatus int8    `json:"trade_status" form:"trade_status"`
	Bp          float64 `json:"bp" form:"bp"`
}

//PaySignNotify .
type PaySignNotify struct {
	ChangeType string `json:"changeType"`
	PayChannel string `json:"payChannel"`
	UID        int64  `json:"uid"`
	CustomerID int64  `json:"customerId"`
}

//PayRefundNotify pay refund notify.
type PayRefundNotify struct {
	CustomerID      int64            `json:"customerId"`
	OrderID         string           `json:"orderId"`
	TxID            int64            `json:"txId"`
	RefundCount     int64            `json:"refundCount"`
	PayChannel      int64            `json:"payChannel"`
	PayChannelID    int32            `json:"payChannelId"`
	BatchRefundList []*PayRefundList `json:"batchRefundList"`
}

//PayRefundList pay refund list.
type PayRefundList struct {
	CustomerRefundID string `json:"customerRefundId"`
	RefundStatus     string `json:"refundStatus"`
	RefundStatusDesc string `json:"refundStatusDesc"`
	RefundAmount     int64  `json:"refundAmount"`
}

// trade status.
const (
	TradeSuccess int8 = iota
	TradeFailed
)
