package http

// RespIncomeAssetOverview .
type RespIncomeAssetOverview struct {
	Total         int64 `json:"total"`
	TotalBuyTimes int64 `json:"total_buy_times"`
	MonthNew      int64 `json:"month_new"`
	DayNew        int64 `json:"day_new"`
}

// ArgIncomeAssetList .
type ArgIncomeAssetList struct {
	Ver string `form:"ver"`
	PS  int64  `form:"ps"`
	PN  int64  `form:"pn"`
}

// RespIncomeAssetList .
type RespIncomeAssetList struct {
	List []*RespIncomeAsset `json:"list"`
	Page RespPage           `json:"page"`
}

// RespIncomeAsset .
type RespIncomeAsset struct {
	OID           int64  `json:"oid"`
	OType         string `json:"otype"`
	Title         string `json:"title"`
	Currency      string `json:"currency"`
	Price         int64  `json:"price"`
	TotalBuyTimes int64  `json:"total_buy_times"`
	NewBuyTimes   int64  `json:"new_buy_times"`
	TotalErrTimes int64  `json:"total_err_times"`
	NewErrTimes   int64  `json:"new_err_times"`
}

// RespPage .
type RespPage struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}
