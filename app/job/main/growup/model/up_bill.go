package model

// UpBill up_bill
type UpBill struct {
	MID            int64
	FirstIncome    int64
	MaxIncome      int64
	TotalIncome    int64
	AvCount        int64
	AvMaxIncome    int64
	AvID           int64
	QualityValue   int64
	DefeatNum      int
	Fans           int64
	TotalPlayCount int64
	Title          string
	ShareItems     string
	FirstTime      string
	MaxTime        string
	SignedAt       string
	EndAt          string
}

// UpQuality up quality
type UpQuality struct {
	MID     int64
	Quality int64
}
