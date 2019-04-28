package model

import (
	"go-common/app/service/main/archive/model/archive"
)

// ArcToView toview video.
type ArcToView struct {
	*archive.Archive3
	Page     *archive.Page3 `json:"page,omitempty"`
	Count    int            `json:"count"`
	Cid      int64          `json:"cid"`
	Progress int64          `json:"progress"`
	AddTime  int64          `json:"add_at"`
}

// WebArcToView toview video.
type WebArcToView struct {
	*archive.View3
	BangumiInfo *Bangumi `json:"bangumi,omitempty"`
	Cid         int64    `json:"cid"`
	Progress    int64    `json:"progress"`
	AddTime     int64    `json:"add_at"`
}

// ToView toview.
type ToView struct {
	Aid  int64 `json:"aid,omitempty"`
	Unix int64 `json:"now,omitempty"`
}

// ToViews toview sorted.
type ToViews []*ToView

func (h ToViews) Len() int           { return len(h) }
func (h ToViews) Less(i, j int) bool { return h[i].Unix > h[j].Unix }
func (h ToViews) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
