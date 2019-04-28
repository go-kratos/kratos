package audio

type Audio struct {
	MenuID      int64  `json:"menu_id,omitempty"`
	Title       string `json:"title,omitempty"`
	CoverURL    string `json:"cover_url,omitempty"`
	RecordNum   int    `json:"record_num,omitempty"`
	PlayNum     int32  `json:"play_num,omitempty"`
	FavoriteNum int32  `json:"favorite_num,omitempty"`
	Face        string `json:"face,omitempty"`
	Songs       []*struct {
		Title string `json:"title,omitempty"`
	} `json:"songs,omitempty"`
	PaTime int64  `json:"pa_time,omitempty"`
	Type   int8   `json:"type,omitempty"`
	Ctgs   []*Ctg `json:"ctgs,omitempty"`
}

type Ctg struct {
	ItemID  int64  `json:"item_id,omitempty"`
	ItemVal string `json:"item_val,omitempty"`
	Schema  string `json:"schema,omitempty"`
}
