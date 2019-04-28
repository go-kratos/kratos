package model

import (
	"fmt"

	"go-common/library/time"
)

const (
	// ViewCountType view count type.
	ViewCountType = "view"
	// FavCountType fav count type.
	FavCountType = "favorite"
	// ReplyCountType reply count type.
	ReplyCountType = "reply"
	// ShareCountType share count type.
	ShareCountType = "share"
)

// StatM  playlist's topic stat message in databus.
type StatM struct {
	Type      string    `json:"type"`
	ID        int64     `json:"id"`
	Aid       int64     `json:"aid"`
	Count     *int64    `json:"count"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
}

// StatMsg means playlist's stat message in databus.
type StatMsg struct {
	Type     string    `json:"type"`
	Pid      int64     `json:"pid"`
	Mid      int64     `json:"mid"`
	Fid      int64     `json:"fid"`
	Aid      int64     `json:"aid"`
	View     *int64    `json:"view"`
	Favorite *int64    `json:"fav"`
	Reply    *int64    `json:"reply"`
	Share    *int64    `json:"share"`
	MTime    time.Time `json:"mtime"`
	IP       string    `json:"ip"`
}

// String format sm
func (sm *StatM) String(tp string) (res string) {
	if sm == nil {
		res = "<nil>"
		return
	}
	res = fmt.Sprintf("pid: %v, aid: %v, ip: %v, "+tp+"(%s) count(%d)", sm.ID, sm.Aid, sm.IP, formatPInt(sm.Count))
	return
}

func formatPInt(s *int64) (res string) {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%d", *s)
}

// Merge merges stat.
func Merge(last, m *StatMsg) {
	if m.View != nil && *m.View >= 0 {
		*last.View = *m.View
	}
	if m.Share != nil && *m.Share >= 0 {
		*last.Share = *m.Share
	}
	if m.Favorite != nil && *m.Favorite >= 0 {
		*last.Favorite = *m.Favorite
	}
	if m.Reply != nil && *m.Reply >= 0 {
		*last.Reply = *m.Reply
	}
}
