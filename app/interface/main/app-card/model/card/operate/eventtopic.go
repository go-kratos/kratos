package operate

type EventTopic struct {
	ID      int64  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Desc    string `json:"desc,omitempty"`
	Cover   string `json:"cover,omitempty"`
	ReType  int8   `json:"re_type,omitempty"`
	ReValue string `json:"re_value,omitempty"`
	Corner  string `json:"corner,omitempty"`
}
