package special

type Card struct {
	ID      int64  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Desc    string `json:"desc,omitempty"`
	Cover   string `json:"cover,omitempty"`
	ReType  int    `json:"re_type,omitempty"`
	ReValue string `json:"re_value,omitempty"`
	Badge   string `json:"badge,omitempty"`
}
