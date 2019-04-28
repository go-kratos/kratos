package model

import (
	"encoding/json"
)

// ViewDm .
type ViewDm struct {
	Closed     bool            `json:"closed"`
	ViewDmSeg  *ViewDmSeg      `json:"dm_seg"`                // 分段弹幕规则
	Flag       json.RawMessage `json:"flag"`                  // flag
	Subtitle   *ViewSubtitle   `json:"subtitle,omitempty"`    // 字幕
	ViewDmMask *Mask           `json:"mask,omitempty"`        // 蒙版
	SpecialDms []string        `json:"special_dms,omitempty"` // 高级弹幕
}

// ViewDmSeg .
type ViewDmSeg struct {
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
}

// ViewSubtitle .
type ViewSubtitle struct {
	Lan       string               `json:"lan"`
	LanDoc    string               `json:"lan_doc"`
	Subtitles []*ViewVideoSubtitle `json:"subtitles"`
}
