package assist

import "time"

type ArgAssists struct {
	Mid    int64
	RealIP string
}

type ArgAssist struct {
	Mid       int64
	AssistMid int64
	Type      int64
	RealIP    string
}

// ArgAssistLogAdd add log
type ArgAssistLogAdd struct {
	Mid       int64
	AssistMid int64
	Type      int64
	Action    int64
	SubjectID int64
	ObjectID  string
	Detail    string
	RealIP    string
}

type ArgAssistLog struct {
	Mid       int64
	AssistMid int64
	LogID     int64
	RealIP    string
}

type ArgAssistLogs struct {
	Mid       int64
	AssistMid int64
	Stime     time.Time
	Etime     time.Time
	Pn        int
	Ps        int
	RealIP    string
}

type ArgAssistUps struct {
	AssistMid int64
	Pn        int64
	Ps        int64
	RealIP    string
}
