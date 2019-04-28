package model

import (
	"encoding/json"

	"go-common/app/job/main/app/model/space"
	xtime "go-common/library/time"
)

const (
	_gotoAv      = 0
	_gotoArticle = 1
	_gotoClip    = 2
	_gotoAlbum   = 3
	_gotoAudio   = 4

	TypeArchive    = "archive"
	TypeArchiveHis = "archive_his"

	TypeForView  = "view"
	TypeForDm    = "dm"
	TypeForReply = "reply"
	TypeForFav   = "fav"
	TypeForCoin  = "coin"
	TypeForShare = "share"
	TypeForLike  = "like"
	TypeForRank  = "rank"
)

type ArcMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
}

type AccMsg struct {
	Mid    int64  `json:"mid"`
	Action string `json:"action"`
}

type StatMsg struct {
	Type         string `json:"type,omitempty"`
	ID           int64  `json:"id,omitempty"`
	Count        int32  `json:"count,omitempty"`
	DislikeCount int32  `json:"dislike_count,omitempty"`
	Timestamp    int64  `json:"timestamp,omitempty"`
	BusType      string
}

type ContributeMsg struct {
	Vmid  int64        `json:"vmid"`
	CTime xtime.Time   `json:"ctime"`
	Attrs *space.Attrs `json:"attrs"`
	IP    string       `json:"ip"`
}

func FormatKey(id int64, gt string) int64 {
	switch gt {
	case GotoAv:
		return id<<6 | _gotoAv
	case GotoArticle:
		return id<<6 | _gotoArticle
	case GotoClip:
		return id<<6 | _gotoClip
	case GotoAlbum:
		return id<<6 | _gotoAlbum
	case GotoAudio:
		return id<<6 | _gotoAudio
	}
	return id
}
