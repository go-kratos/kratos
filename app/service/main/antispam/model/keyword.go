package model

import (
	"fmt"

	"go-common/app/service/main/antispam/util"
)

const (
	// ParamKeywordHitCounts .
	ParamKeywordHitCounts = "show_up_counts"

	// KeywordTag .
	KeywordTag = "tag"
	// KeywordTagBlack .
	KeywordTagBlack = "black"
	// KeywordTagWhite .
	KeywordTagWhite = "white"
	// KeywordTagDefaultLimit .
	KeywordTagDefaultLimit = "limit"
	// KeywordTagRestrictLimit .
	KeywordTagRestrictLimit = "restrict"
	// KeywordContent .
	KeywordContent = "content"
	// KeywordHitCounts .
	KeywordHitCounts = "hit_counts"
)

// SenderList .
type SenderList struct {
	SenderIDs []int64 `json:"sender_ids"`
	Counts    int     `json:"counts"`
}

// Keyword .
type Keyword struct {
	ID            int64         `json:"id"`
	Area          string        `json:"-"`
	Content       string        `json:"content"`
	SenderID      int64         `json:"-"`
	OriginContent string        `json:"origin_content"`
	SenderCounts  int64         `json:"sender_counts"`
	RegexpName    string        `json:"reg_name"`
	Tag           string        `json:"tag"`
	State         string        `json:"state"`
	HitCounts     int64         `json:"show_up_counts"`
	CTime         util.JSONTime `json:"ctime"`
	MTime         util.JSONTime `json:"mtime"`
}

func (k *Keyword) String() string {
	return fmt.Sprintf("id: %d, area: %s, content: %s, tag: %s, state: %s, hitCounts %d\n",
		k.ID, k.Area, k.Content, k.Tag, k.State, k.HitCounts)
}
