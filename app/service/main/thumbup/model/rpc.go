package model

// ArgLike .
type ArgLike struct {
	Mid       int64
	UpMid     int64
	Business  string
	OriginID  int64
	MessageID int64
	Type      int8
	RealIP    string
}

// ArgHasLike .
type ArgHasLike struct {
	Business   string
	MessageIDs []int64
	Mid        int64
	RealIP     string
}

// ArgStats .
type ArgStats struct {
	Business   string
	OriginID   int64
	MessageIDs []int64
	RealIP     string
}

// ArgStatsWithLike .
type ArgStatsWithLike struct {
	Mid        int64
	Business   string
	OriginID   int64
	MessageIDs []int64
	RealIP     string
}

// ArgUserLikes .
type ArgUserLikes struct {
	Business string
	Mid      int64
	Pn, Ps   int
	RealIP   string
}

// ArgItemLikes .
type ArgItemLikes struct {
	Business  string
	OriginID  int64
	MessageID int64
	Mid       int64
	Pn, Ps    int
	RealIP    string
}

// ArgUpdateCount .
type ArgUpdateCount struct {
	Business      string
	OriginID      int64
	MessageID     int64
	LikeChange    int64
	DislikeChange int64
	Operator      string
	RealIP        string
}

// ArgRawStats .
type ArgRawStats struct {
	Business  string
	OriginID  int64
	MessageID int64
	RealIP    string
}
