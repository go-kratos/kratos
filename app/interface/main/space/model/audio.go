package model

// AudioCard .
type AudioCard struct {
	Type   int `json:"type"`
	Status int `json:"status"`
}

// AudioUpperCert .
type AudioUpperCert struct {
	Cert *struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"cert"`
}

// Audio .
type Audio struct {
	MenuID    int64  `json:"menu_id"`
	Title     string `json:"title"`
	CoverURL  string `json:"cover_url"`
	PlayNum   int    `json:"play_num"`
	RecordNum int    `json:"record_num"`
	Ctgs      []*struct {
		ItemID  int64  `json:"item_id"`
		ItemVal string `json:"item_val"`
	} `json:"ctgs"`
	Type int `json:"type"`
}
