package elec

import "encoding/json"

type Info struct {
	Start   int64           `json:"start,omitempty"`
	Show    bool            `json:"show"`
	Total   int             `json:"total,omitempty"`
	Count   int             `json:"count,omitempty"`
	State   int             `json:"state,omitempty"`
	List    json.RawMessage `json:"list,omitempty"`
	User    json.RawMessage `json:"user,omitempty"`
	ElecSet json.RawMessage `json:"elec_set,omitempty"`
}
