package model

import "go-common/app/service/main/archive/api"

// is
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
	Click   int32  `json:"click"`
	DM      int32  `json:"dm"`
	Reply   int32  `json:"reply"`
	Fav     int32  `json:"fav"`
	Coin    int32  `json:"coin"`
	Share   int32  `json:"share"`
	NowRank int32  `json:"now_rank"`
	HisRank int32  `json:"his_rank"`
	Like    int32  `json:"like"`
	Type    string `json:"-"`
	Ts      int64  `json:"-"`
}

// StatCount is
type StatCount struct {
	Type      string `json:"type"`
	Aid       int64  `json:"id"`
	Count     int32  `json:"count"`
	TimeStamp int64  `json:"timestamp"`
}

// Merge merge message and stat from db.
func Merge(m *StatMsg, s *api.Stat) {
	if m.Click >= 0 && m.Type == TypeForView {
		s.View = m.Click
	}
	if m.Coin >= 0 && m.Type == TypeForCoin {
		s.Coin = m.Coin
	}
	if m.DM >= 0 && m.Type == TypeForDm {
		s.Danmaku = m.DM
	}
	if m.Fav >= 0 && m.Type == TypeForFav {
		s.Fav = m.Fav
	}
	if m.Reply >= 0 && m.Type == TypeForReply {
		s.Reply = m.Reply
	}
	if m.Share >= 0 && m.Type == TypeForShare && m.Share > s.Share {
		s.Share = m.Share
	}
	if m.NowRank >= 0 && m.Type == TypeForRank {
		s.NowRank = m.NowRank
	}
	if m.HisRank >= 0 && m.Type == TypeForRank {
		s.HisRank = m.HisRank
	}
	if m.Like >= 0 && m.Type == TypeForLike {
		s.Like = m.Like
	}
}
