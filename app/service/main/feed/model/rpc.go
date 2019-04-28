package model

type ArgFeed struct {
	Mid    int64
	Pn     int
	Ps     int
	RealIP string
}

type ArgArc struct {
	Aid     int64
	Mid     int64
	PubDate int64
	RealIP  string
}

type ArgAidMid struct {
	Aid    int64
	Mid    int64
	RealIP string
}

type ArgMid struct {
	Mid    int64
	RealIP string
}

type ArgUnreadCount struct {
	Mid            int64
	WithoutBangumi bool
	RealIP         string
}

type ArgFold struct {
	Aid    int64
	Mid    int64
	RealIP string
}

type ArgChangeUpper struct {
	Aid    int64
	OldMid int64
	NewMid int64
	RealIP string
}
