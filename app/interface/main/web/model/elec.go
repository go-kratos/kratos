package model

import "encoding/json"

// ElecShow elec show
type ElecShow struct {
	ShowInfo   *ShowInfo       `json:"show_info"`
	AvCount    int             `json:"av_count"`
	Count      int             `json:"count"`
	TotalCount int64           `json:"total_count"`
	SpecialDay int             `json:"special_day"`
	DisplayNum int             `json:"display_num"`
	AvList     json.RawMessage `json:"av_list,omitempty"`
	AvUser     json.RawMessage `json:"av_user,omitempty"`
	List       json.RawMessage `json:"list,omitempty"`
	User       json.RawMessage `json:"user,omitempty"`
}

// ShowInfo show info
type ShowInfo struct {
	Show   bool   `json:"show"`
	State  int8   `json:"state"`
	Reason string `json:"reason,omitempty"`
}
