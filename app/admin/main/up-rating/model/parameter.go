package model

// RatingParameter rating weight args
type RatingParameter struct {
	WDP      int64 // dp weight
	WDC      int64 // dc weight
	WDV      int64 // dv weight
	WMDV     int64 // mdv weight
	WCS      int64
	WCSR     int64
	WMAAFans int64
	WMAHFans int64
	WIS      int64
	WISR     int64
	// 信用分
	HBASE int64
	HR    int64
	HV    int64
	HVM   int64
	HL    int64
	HLM   int64
}
