package model

const (
	// BangumiTyp 番剧
	BangumiTyp = 1
	// ComicTyp 漫画
	ComicTyp = 2
	// ArchiveTyp 稿件
	ArchiveTyp = 3
	// PlaylistTyp 播单
	PlaylistTyp = 4
	// ArchiveMsgTyp .
	ArchiveMsgTyp = "archive"
	// BangumiMsgTyp .
	BangumiMsgTyp = "bangumi_share"
	// ComicMsgTyp .
	ComicMsgTyp = "comic_share"
	// PlaylistMsgTyp .
	PlaylistMsgTyp = "playlist"
)

// ShareParams .
type ShareParams struct {
	OID int64  `json:"oid" form:"oid" validate:"required,gt=0"`
	MID int64  `json:"mid" form:"mid" validate:"required"`
	TP  int    `json:"tp" form:"tp" validate:"required,gt=0"`
	IP  string `json:"ip"`
}

// Share share item
type Share struct {
	OID   int64 `json:"oid"`
	Tp    int   `json:"tp"`
	Count int64 `json:"count"`
}

// MIDShare .
type MIDShare struct {
	OID  int64 `json:"oid"`
	MID  int64 `json:"mid"`
	TP   int   `json:"tp"`
	Time int64 `json:"time"`
}

// ArchiveShare .
type ArchiveShare struct {
	Type  string `json:"type"`
	ID    int64  `json:"id"`
	Count int    `json:"count"`
	Ts    int64  `json:"timestamp"`
}
