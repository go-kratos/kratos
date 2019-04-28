package model

// TradeOrder .
type TradeOrder struct {
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

// IncomeAssetOverview .
type IncomeAssetOverview struct {
	Total         int64
	TotalBuyTimes int64
	MonthNew      int64
	DayNew        int64
}

// IncomeAssetMonthly .
type IncomeAssetMonthly struct {
	List []*IncomeAssetMonthlyByContent
	Page *Page
}

// IncomeAssetMonthlyByContent .
type IncomeAssetMonthlyByContent struct {
	OID           int64
	OType         string
	Currency      string
	Title         string
	Price         int64
	TotalBuyTimes int64
	NewBuyTimes   int64
	TotalErrTimes int64
	NewErrTimes   int64
}

// Page .
type Page struct {
	Num   int64
	Size  int64
	Total int64
}
