package http

// ArgTradeCallback .
type ArgTradeCallback struct {
	MSGID      int64  `form:"msgId" validate:"required"`
	MSGContent string `form:"msgContent" validate:"required"`
}

// ArgTradeRefund .
type ArgTradeRefund struct {
	OrderID string `form:"order_id" validate:"required"`
}

// ArgTradeRefunds .
type ArgTradeRefunds struct {
	OrderIDs []string `form:"order_ids,split" validate:"required"`
}
