package model

// ArgBind bind args.
type ArgBind struct {
	OpenID    string
	OutOpenID string
	AppID     int64
}

// ArgBindInfo bind info args.
type ArgBindInfo struct {
	Mid   int64
	AppID int64
}

// ArgThirdPrizeGrant prize grant args.
type ArgThirdPrizeGrant struct {
	Mid       int64  `form:"mid" validate:"required"`
	PrizeKey  int64  `form:"prize_key"`
	UniqueNo  string `form:"unique_no" validate:"required"`
	PrizeType int8   `form:"prize_type" validate:"required"`
	Appkey    string `form:"appkey" validate:"required"`
	Remark    string `form:"remark" validate:"required"`
	AppID     int64
}

// ArgBilibiliPrizeGrant args.
type ArgBilibiliPrizeGrant struct {
	PrizeKey string
	UniqueNo string
	OpenID   string
	AppID    int64
}

// BilibiliPrizeGrantResp resp.
type BilibiliPrizeGrantResp struct {
	Amount      float64
	FullAmount  float64
	Description string
}
