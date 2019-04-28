package ai

import (
	"encoding/json"

	"go-common/app/interface/main/app-card/model/card/banner"
	"go-common/app/interface/main/app-card/model/card/cm"
	tag "go-common/app/interface/main/tag/model"
	"go-common/app/service/main/archive/model/archive"
)

type Item struct {
	ID         int64           `json:"id,omitempty"`
	TrackID    string          `json:"trackid,omitempty"`
	Name       string          `json:"name,omitempty"`
	Goto       string          `json:"goto,omitempty"`
	Tid        int64           `json:"tid,omitempty"`
	From       int8            `json:"from,omitempty"`
	Source     string          `json:"source,omitempty"`
	AvFeature  json.RawMessage `json:"av_feature,omitempty"`
	Config     *Config         `json:"config,omitempty"`
	RcmdReason *RcmdReason     `json:"rcmd_reason,omitempty"`
	StatType   int8            `json:"stat_type,omitempty"`
	Style      int8            `json:"style,omitempty"`
	// extra
	Archive    *archive.Archive3 `json:"archive,omitempty"`
	Tag        *tag.Tag          `json:"tag,omitempty"`
	Ad         *cm.AdInfo        `json:"-"`
	Banners    []*banner.Banner  `json:"-"`
	Version    string            `json:"-"`
	HideButton bool              `json:"-"`
	CornerMark int8              `json:"corner_mark,omitempty"`
}

type Config struct {
	URI   string `json:"uri,omitempty"`
	Title string `json:"title,omitempty"`
	Cover string `json:"cover,omitempty"`
}

type RcmdReason struct {
	ID           int    `json:"id,omitempty"`
	Content      string `json:"content,omitempty"`
	BgColor      string `json:"bg_color,omitempty"`
	IconLocation string `json:"icon_location,omitempty"`
	Style        int    `json:"style,omitempty"`
	Font         int    `json:"font,omitempty"`
	Position     string `json:"position,omitempty"`
	Grounding    string `json:"grounding,omitempty"`
	CornerMark   int8   `json:"corner_mark,omitempty"`
	FollowedMid  int64  `json:"followed_mid,omitempty"`
}
