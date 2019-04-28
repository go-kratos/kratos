package tag

type Tag struct {
	Tid     int64  `json:"tag_id"`
	Name    string `json:"tag_name"`
	IsAtten int8   `json:"is_atten"`
	Count   struct {
		Atten int `json:"atten,omitempty"`
	} `json:"count,omitempty"`
}
