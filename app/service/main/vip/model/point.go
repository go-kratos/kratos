package model

//PointExchangePrice .
type PointExchangePrice struct {
	ID             int64  `json:"id"`
	OriginPoint    int32  `json:"originPoint"`
	CurrentPoint   int32  `json:"currentPoint"`
	Month          int16  `json:"month"`
	PromotionTip   string `json:"promotionTip"`
	PromotionColor string `json:"promotionColor"`
	OperatorID     string `json:"operatorId"`
}

// point consume status
const (
	PointConsumeSuc   = 1
	PointConsumeFaild = 1
)
