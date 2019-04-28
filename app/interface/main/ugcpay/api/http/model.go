package http

// ArgTradeOrder .
type ArgTradeOrder struct {
	OrderID string `form:"order_id" validate:"required"`
}

// RespTradeOrder .
type RespTradeOrder struct {
	OrderID  string `json:"order_id"`
	MID      int64  `json:"mid"`
	Biz      string `json:"biz"`
	Platform string `json:"platform"`
	OID      int64  `json:"oid"`
	OType    string `json:"otype"`
	Fee      int64  `json:"fee"`
	Currency string `json:"currency"`
	PayID    string `json:"pay_id"`
	State    string `json:"state"`
	Reason   string `json:"reason"`
}

// ArgTradeCreate .
type ArgTradeCreate struct {
	OID      int64  `form:"oid" validate:"required"`
	OType    string `form:"otype" validate:"required"`
	Currency string `form:"currency" validate:"required"`
}

// RespTradeCreate .
type RespTradeCreate struct {
	OrderID string `json:"order_id"`
	PayData string `json:"pay_data"`
}
