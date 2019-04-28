package shop

type Info struct {
	Shop *struct {
		ID     int64  `json:"id,omitempty"`
		Name   string `json:"name,omitempty"`
		Status int    `json:"status,omitempty"`
	} `json:"shop,omitempty"`
}
