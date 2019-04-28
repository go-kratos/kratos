package model

// TokenInfo represents pay token.
type TokenInfo struct {
	Token   string
	OrderNo string
	Status  int8
	Mid     int64
}

// CopyFromPayParam copies fields from pay param.
func (t *TokenInfo) CopyFromPayParam(pp *PayParam) {
	t.Mid = pp.Mid
	t.Status = pp.Status
	t.OrderNo = pp.OrderNo
}
