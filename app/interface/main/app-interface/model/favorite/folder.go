package favorite

import "go-common/app/service/main/archive/api"

type Folder struct {
	MediaID    int64   `json:"media_id"`
	Fid        int     `json:"fid"`
	Mid        int     `json:"mid"`
	Name       string  `json:"name"`
	MaxCount   int     `json:"max_count"`
	CurCount   int     `json:"cur_count"`
	AttenCount int     `json:"atten_count"`
	State      int     `json:"state"`
	CTime      int     `json:"ctime"`
	MTime      int     `json:"mtime"`
	Cover      []Cover `json:"cover,omitempty"`
	Videos     []Cover `json:"videos,omitempty"` // NOTE: old favourite
}

type Cover struct {
	Aid  int    `json:"aid"`
	Pic  string `json:"pic"`
	Type int32  `json:"type"`
}

type Video struct {
	Seid           string `json:"seid"`
	Page           int    `json:"page"`
	Pagesize       int    `json:"pagesize"`
	PageCount      int    `json:"pagecount"`
	Total          int    `json:"total"`
	SuggestKeyword string `json:"suggest_keyword"`
	Mid            int64  `json:"mid"`
	Fid            int64  `json:"fid"`
	Tid            int    `json:"tid"`
	Order          string `json:"order"`
	Keyword        string `json:"keyword"`
	Tlist          []struct {
		Tid   int16  `json:"tid"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	} `json:"tlist,omitempty"`
	Archives []*Archive `json:"archives"`
}

type Archive struct {
	*api.Arc
	FavAt          int64  `json:"fav_at"`
	PlayNum        string `json:"play_num"`
	HighlightTitle string `json:"highlight_title"`
}
