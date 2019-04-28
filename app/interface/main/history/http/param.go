package http

// Histroy Histroy.
type Histroy struct {
	Pn int  `form:"pn"`
	Ps int  `form:"ps"`
	TP int8 `form:"type"`
}

// AddHistory AddHistory.
type AddHistory struct {
	Aid      int64  `form:"aid" validate:"required,gt=0"`
	Cid      int64  `form:"cid"`
	Epid     int64  `form:"epid"`
	TP       int8   `form:"type"`
	SubTP    int8   `form:"sub_type"`
	DT       int8   `form:"dt"`
	Sid      int64  `form:"sid"`
	Platform string `form:"platform"`
	Device   string `form:"device"`
}

// HistoryReport HistoryReport.
type HistoryReport struct {
	Mid       int64  `form:"mid"`
	Aid       int64  `form:"aid"`
	Type      int8   `form:"type"`
	Cid       int64  `form:"cid"`
	Epid      int64  `form:"epid"`
	Sid       int64  `form:"sid"`
	SubTP     int8   `form:"subtype"`
	SubType   int8   `form:"sub_type"`
	DT        int8   `form:"dt"`
	Realtime  int64  `form:"realtime"`
	Source    int64  `form:"source"`
	Progress  int64  `form:"progress"`
	Platform  string `form:"platform"`
	Device    string `form:"device"`
	PlayTime  int64  `form:"play_time"`
	MobileApp string `form:"mobi_app"`
}

// Page Page.
type Page struct {
	Pn int `form:"pn"`
	Ps int `form:"ps"`
}
