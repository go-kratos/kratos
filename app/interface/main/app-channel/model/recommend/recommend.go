package recommend

import (
	"encoding/json"

	tag "go-common/app/interface/main/tag/model"
	"go-common/app/service/main/archive/api"
)

type Item struct {
	ID         int64           `json:"id,omitempty"`
	Name       string          `json:"name,omitempty"`
	Goto       string          `json:"goto,omitempty"` // goto: av, live, bangumi, topic, activity, ad
	TagID      int64           `json:"tid,omitempty"`
	From       int8            `json:"from,omitempty"`
	Source     string          `json:"source,omitempty"`
	AvFeature  json.RawMessage `json:"av_feature,omitempty"`
	Config     *Config         `json:"config,omitempty"`
	RcmdReason *RcmdReason     `json:"rcmd_reason,omitempty"`
	StatType   int8            `json:"stat_type,omitempty"`
	Items      []*Item         `json:"-"`
	Archive    *api.Arc        `json:"archive,omitempty"`
	Tag        *tag.Tag        `json:"tag,omitempty"`
	Limit      int             `json:"-"`
}

type Config struct {
	URI      string `json:"uri,omitempty"`
	Title    string `json:"title,omitempty"`
	Cover    string `json:"cover,omitempty"`
	Content  string `json:"content,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}

type RcmdReason struct {
	ID           int    `json:"id"`
	Content      string `json:"content"`
	BgColor      string `json:"bg_color"`
	IconLocation string `json:"icon_location"`
	Style        int    `json:"style"`
	Font         int    `jsn:"font"`
	Position     string `jsn:"position"`
	Grounding    string `jsn:"grounding"`
}
