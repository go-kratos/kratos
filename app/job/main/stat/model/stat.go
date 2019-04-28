package model

import "go-common/app/service/main/archive/api"

const (
	TypeForView  = "view"
	TypeForDm    = "dm"
	TypeForReply = "reply"
	TypeForFav   = "fav"
	TypeForCoin  = "coin"
	TypeForShare = "share"
	TypeForRank  = "rank"
	TypeForLike  = "like"
)

// StatMsg stat info.
type StatMsg struct {
	Aid     int64  `json:"aid"`
	Click   int    `json:"click"`
	DM      int    `json:"dm"`
	Reply   int    `json:"reply"`
	Fav     int    `json:"fav"`
	Coin    int    `json:"coin"`
	Share   int    `json:"share"`
	NowRank int    `json:"now_rank"`
	HisRank int    `json:"his_rank"`
	Like    int    `json:"like"`
	DisLike int    `json:"dislike_count"`
	Type    string `json:"-"`
	Ts      int64  `json:"-"`
}

type StatCount struct {
	Type      string `json:"type"`
	Aid       int64  `json:"id"`
	Count     int    `json:"count"`
	DisLike   int    `json:"dislike_count"`
	TimeStamp int64  `json:"timestamp"`
}

// Merge merge message and stat from db.
func Merge(m *StatMsg, s *api.Stat) {
	if m.Click >= 0 && m.Type == TypeForView {
		s.View = int32(m.Click)
	}
	if m.Coin >= 0 && m.Type == TypeForCoin {
		s.Coin = int32(m.Coin)
	}
	if m.DM >= 0 && m.Type == TypeForDm {
		s.Danmaku = int32(m.DM)
	}
	if m.Fav >= 0 && m.Type == TypeForFav {
		s.Fav = int32(m.Fav)
	}
	if m.Reply >= 0 && m.Type == TypeForReply {
		s.Reply = int32(m.Reply)
	}
	if m.Share >= 0 && m.Type == TypeForShare && int32(m.Share) > s.Share {
		s.Share = int32(m.Share)
	}
	if m.NowRank >= 0 && m.Type == TypeForRank {
		s.NowRank = int32(m.NowRank)
	}
	if m.HisRank >= 0 && m.Type == TypeForRank {
		s.HisRank = int32(m.HisRank)
	}
	if m.Like >= 0 && m.Type == TypeForLike {
		s.Like = int32(m.Like)
	}
	if m.DisLike >= 0 && m.Type == TypeForLike {
		s.DisLike = int32(m.DisLike)
	}
}
