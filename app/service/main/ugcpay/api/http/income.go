package http

// ArgIncomeAssetOverview .
type ArgIncomeAssetOverview struct {
	MID int64 `form:"mid" validate:"required"`
}

// RespIncomeAssetOverview .
type RespIncomeAssetOverview struct {
	Total         int64 `json:"total"`
	TotalBuyTimes int64 `json:"total_buy_times"`
	MonthNew      int64 `json:"month_new"`
	DayNew        int64 `json:"day_new"`
}

// ArgIncomeAssetMonthly .
type ArgIncomeAssetMonthly struct {
	MID int64  `form:"mid" validate:"required"`
	Ver string `form:"ver"`
}

// RespIncomeAssetMonthly .
type RespIncomeAssetMonthly struct {
	List []*RespIncomeAssetMonthlyByContent `json:"list"`
}

// RespIncomeAssetMonthlyByContent .
type RespIncomeAssetMonthlyByContent struct {
	OID           int64  `json:"oid"`
	OType         string `json:"otype"`
	Currency      string `json:"currency"`
	Price         int64  `json:"price"`
	TotalBuyTimes int64  `json:"total_buy_times"`
	NewBuyTimes   int64  `json:"new_buy_times"`
	TotalErrTimes int64  `json:"total_err_times"`
	NewErrTimes   int64  `json:"new_err_times"`
}
