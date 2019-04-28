package model

// ArgReport .
type ArgReport struct {
	ID           int64
	APPID        int64
	PlatformID   int
	Mid          int64
	Buvid        string
	DeviceToken  string
	Build        int
	TimeZone     int
	NotifySwitch int
	DeviceBrand  string
	DeviceModel  string
	OSVersion    string
	Extra        string
	RealIP       string
}

// ArgReports .
type ArgReports struct {
	Reports []*Report
}

// ArgUserReports .
type ArgUserReports struct {
	Mid     int64
	Reports []*Report
}

// ArgToken .
type ArgToken struct {
	Token  string
	RealIP string
}

// ArgMidToken .
type ArgMidToken struct {
	Mid    int64
	Token  string
	RealIP string
}

// ArgDelInvalidReport .
type ArgDelInvalidReport struct {
	Type   int
	RealIP string
}

// ArgMid .
type ArgMid struct {
	Mid    int64
	RealIP string
}

// ArgSetting .
type ArgSetting struct {
	Mid    int64
	Type   int
	Value  int
	RealIP string
}

// ArgCallback .
type ArgCallback struct {
	Task     string
	APP      int64
	Platform int
	Mid      int64
	Pid      int // mobi_app ID
	Token    string
	Buvid    string
	Click    uint8 // 是否被点击
	Extra    *CallbackExtra
}

// ArgMidProgress .
type ArgMidProgress struct {
	Task     string
	MidTotal int64
	MidValid int64
	RealIP   string
}
