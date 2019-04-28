package topic

type Topic struct {
	Page     int     `json:"page"`
	PageSize int     `json:"pagesize"`
	Total    int     `json:"total"`
	Lists    []*List `json:"list"`
}

type List struct {
	ID       int64  `json:"id"`
	TpID     int64  `json:"tp_id"`
	MID      int64  `json:"mid"`
	FavAt    int64  `json:"fav_at"`
	State    int64  `json:"state"`
	Stime    string `json:"stime"`
	Etime    string `json:"etime"`
	Ctime    string `json:"ctime"`
	Mtime    string `json:"mtime"`
	Name     string `json:"name"`
	Author   string `json:"author"`
	PCUrl    string `json:"pc_url"`
	H5Url    string `json:"h5_url"`
	PCCover  string `json:"pc_cover"`
	H5Cover  string `json:"h5_cover"`
	Rank     int64  `json:"rank"`
	PageName string `json:"page_name"`
	Plat     int64  `json:"plat"`
	Desc     string `json:"desc"`
	Click    int64  `json:"click"`
	TPType   int64  `json:"type"`
	Mold     int64  `json:"mold"`
	Series   int64  `json:"series"`
	Dept     int64  `json:"dept"`
	ReplyID  int64  `json:"reply_id"`
}
