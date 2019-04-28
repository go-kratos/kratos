package model

// Season .
type Season struct {
	ID        int64  `json:"id" form:"id"`
	Mid       int64  `json:"mid" form:"mid" validate:"required"`
	Title     string `json:"title" form:"title" validate:"required"`
	SubTitle  string `json:"sub_title" form:"sub_title"`
	Stime     int64  `json:"stime" form:"stime"`
	Etime     int64  `json:"etime" form:"etime"`
	Sponsor   string `json:"sponsor" form:"sponsor"`
	Logo      string `json:"logo" form:"logo" validate:"required"`
	Dic       string `json:"dic" form:"dic"`
	Status    int    `json:"status"  form:"is_deleted"`
	IsApp     int    `json:"is_app" form:"is_app"`
	Rank      int    `json:"rank" form:"rank" validate:"min=0,max=10"`
	URL       string `json:"url" form:"url"`
	DataFocus string `json:"data_focus" form:"data_focus"`
	FocusURL  string `json:"focus_url" form:"focus_url"`
}

// SeasonInfo .
type SeasonInfo struct {
	*Season
	Games []*Game `json:"games"`
}

// TableName .
func (s Season) TableName() string {
	return "es_seasons"
}
