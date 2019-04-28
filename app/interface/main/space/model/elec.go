package model

import "encoding/json"

// ElecInfo elec info.
type ElecInfo struct {
	Start   int64           `json:"start"`
	Show    bool            `json:"show"`
	Total   int             `json:"total"`
	Count   int             `json:"count"`
	State   int             `json:"state"`
	List    json.RawMessage `json:"list,omitempty"`
	User    json.RawMessage `json:"user,omitempty"`
	ElecSet json.RawMessage `json:"elec_set,omitempty"`
}
