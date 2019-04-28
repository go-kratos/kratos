package sp

type Specil struct {
	Results int              `json:"results"`
	Pages   int              `json:"pages"`
	Items   map[string]*Item `json:"list"`
}

type Item struct {
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	MCover    string `json:"m_cover"`
	SCover    string `json:"s_cover"`
	CTime     string `json:"create_at"`
	SpID      int64  `json:"spid"`
	Attention int    `json:"attention"`
}
