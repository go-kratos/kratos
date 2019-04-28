package model

import (
	"encoding/json"
	"math"
)

// 各种状态枚举
const (
	OrderStatePaid            = "paid"
	OrderStateSettled         = "settled"
	OrderStateBadDebt         = "bad_debt"
	OrderStateRefundFinished  = "ref_finished"
	OrderStateSettledRefunded = "st_refunded"

	BizAsset = "asset"

	CurrencyBP = "bp"

	StateRunning = "running"
	StateValid   = "valid"

	AccountStateIncome   = "income"
	AccountStateWithdraw = "withdraw"
	AccountStateProfit   = "profit"
	AccountStateLoss     = "fill_loss"

	PayCheckOrderStateING     = "WAIT_RECONCILIATION"
	PayCheckOrderStateSuccess = "RECONCILIATION_SUCCESS"
	PayCheckOrderStateFail    = "RECONCILIATION_FAIL"

	DefaultUserSetting = math.MaxInt32
)

// Message binlog databus msg.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// BinlogOrderUser .
type BinlogOrderUser struct {
	OrderID string `json:"order_id"`
}

// BinlogAsset .
type BinlogAsset struct {
	OID      int64  `json:"oid"`
	OType    string `json:"otype"`
	Currency string `json:"currency"`
}

// BinlogAssetRelation .
type BinlogAssetRelation struct {
	OID   int64  `json:"oid"`
	OType string `json:"otype"`
	MID   int64  `json:"mid"`
}

// PayCheckRefundOrder .
type PayCheckRefundOrder struct {
	Elements []*PayCheckRefundOrderEle `json:"batchRefundBillVOS"`
	TXID     string                    `json:"txId"`
}

// PayCheckRefundOrderEle .
type PayCheckRefundOrderEle struct {
	RefundNO         string `json:"refundNo"`
	RefundAmount     int64  `json:"refundAmount"`
	CustomerRefundID string `json:"customerRefundId"`
	RecoStatusDesc   string `json:"recoStatusDesc"`
	TXID             string `json:"txId"`
}

// PayCheckOrder .
type PayCheckOrder struct {
	PayChannelOrderNo string `json:"payChannelOrderNo"` //第三方支付渠道支付流水号
	TxID              string `json:"txId"`
	BankAmount        int64  `json:"bankAmount"`     // 订单支付金额
	PayTime           int64  `json:"payTime"`        // 订单支付时间，毫秒值
	RecoStatusDesc    string `json:"recoStatusDesc"` // 对账状态 WAIT_RECONCILIATION（对账中）,RECONCILIATION_SUCCESS（对账成功）,RECONCILIATION_FAIL（对账失败）
}

// PayQuery .
type PayQuery struct {
	Orders []*PayOrder `json:"orders"`
}

// PayOrder .
type PayOrder struct {
	TXID          int64  `json:"txId"`
	OrderID       string `json:"orderId"`
	PayStatus     string `json:"payStatus"`
	PayStatusDesc string `json:"payStatusDesc"`
	FailReason    string `json:"failReason"`
}
