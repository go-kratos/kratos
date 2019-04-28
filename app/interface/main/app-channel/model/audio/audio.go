package audio

type Audio struct {
	MenuID      int64  `json:"menu_id,omitempty"`
	Title       string `json:"title,omitempty"`
	CoverURL    string `json:"cover_url,omitempty"`
	PlayNum     int    `json:"play_num,omitempty"`
	RecordNum   int    `json:"record_num,omitempty"`
	FavoriteNum int    `json:"favorite_num,omitempty"`
	Author      string `json:"author,omitempty"`
	Face        string `json:"face,omitempty"`
	Songs       []*struct {
		Title string `json:"title,omitempty"`
	} `json:"songs,omitempty"`
	PaTime int64 `json:"pa_time,omitempty"`
	Ctgs   []*struct {
		ItemID  int64  `json:"item_id,omitempty"`
		ItemVal string `json:"item_val,omitempty"`
	} `json:"ctgs,omitempty"`
	Type int `json:"type,omitempty"`
}

type Song struct {
	SongID   int64  `json:"song_id,omitempty"`
	Title    string `json:"title,omitempty"`
	CoverURL string `json:"cover_url,omitempty"`
	PlayNum  int    `json:"play_num,omitempty"`
	ReplyNum int    `json:"reply_num,omitempty"`
	Ctgs     []*struct {
		ItemID  int64  `json:"item_id,omitempty"`
		ItemVal string `json:"item_val,omitempty"`
	} `json:"ctgs,omitempty"`
}
