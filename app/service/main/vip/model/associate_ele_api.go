package model

// ArgEleAccessToken ele access token args.
type ArgEleAccessToken struct {
	AuthCode string `json:"auth_code"`
}

// EleAccessTokenResp ele access token resp.
type EleAccessTokenResp struct {
	OpenID string `json:"open_id"`
}

// ArgEleReceivePrizes receive prizes args.
type ArgEleReceivePrizes struct {
	ElemeOpenID string `json:"eleme_open_id"`
	BliOpenID   string `json:"bli_open_id"`
	SourceID    string `json:"source_id"`
}

// EleReceivePrizesResp  receive prizes resp.
type EleReceivePrizesResp struct {
	Amount       float64 `json:"amount"`
	SumCondition float64 `json:"sum_condition"`
	Description  string  `json:"description"`
}

// ArgEleUnionUpdateOpenID union update open id args.
type ArgEleUnionUpdateOpenID struct {
	ElemeOpenID string `json:"eleme_open_id"`
	BliOpenID   string `json:"bli_open_id"`
}

// EleUnionUpdateOpenIDResp union update resp.
type EleUnionUpdateOpenIDResp struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}

// ArgEleBindUnion ele salary bind vip args.
type ArgEleBindUnion struct {
	ElemeOpenID string `json:"eleme_open_id"`
	BliOpenID   string `json:"bli_open_id"`
	VipType     int32  `json:"vip_type"`
	SourceID    string `json:"source_id"`
	UserIP      string `json:"user_ip"`
}

// EleBindUnionResp ele bind union resp.
type EleBindUnionResp struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}

// ArgEleCanPurchase ele can purchase args.
type ArgEleCanPurchase struct {
	ElemeOpenID string `json:"eleme_open_id"`
	BliOpenID   string `json:"bli_open_id"`
	UserIP      string `json:"user_ip"`
	VipType     int32  `json:"vip_type"`
}

// EleCanPurchaseResp ele can purchase resp.
type EleCanPurchaseResp struct {
	CanPurchase bool   `json:"can_purchase"`
	Status      int32  `json:"status"`
	Message     string `json:"message"`
}

// ArgEleUnionMobile ele union mobile.
type ArgEleUnionMobile struct {
	ElemeOpenID string `json:"eleme_open_id"`
	BliOpenID   string `json:"bli_open_id"`
}

// EleUnionMobileResp ele get union mobile resp.
type EleUnionMobileResp struct {
	Status     int32  `json:"status"`
	Message    string `json:"message"`
	BlurMobile string `json:"blur_mobile"`
}

// EleRedPackagesResp ele red packages.
type EleRedPackagesResp struct {
	Name         string  `json:"name"`
	Amount       float64 `json:"amount"`
	SumCondition float64 `json:"sum_condition"`
}

// EleSpecailFoodsResp ele specail foods resp.
type EleSpecailFoodsResp struct {
	RestaurantName string  `json:"restaurant_name"`
	FoodName       string  `json:"food_name"`
	FoodURL        string  `json:"food_url"`
	Discount       float64 `json:"discount"`
	Amount         float64 `json:"amount"`
	OriginalAmount float64 `json:"original_amount"`
	RatingPoint    float64 `json:"rating_point"`
}
