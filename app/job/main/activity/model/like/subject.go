package like

import xtime "go-common/library/time"

// Subject subject
type Subject struct {
	ID       int64      `json:"id"`
	Name     string     `json:"name"`
	Dic      string     `json:"dic"`
	Cover    string     `json:"cover"`
	Stime    xtime.Time `json:"stime"`
	Etime    xtime.Time `json:"etime"`
	Interval int64      `json:"interval"`
	Tlimit   int64      `json:"tlimit"`
	Ltime    int64      `json:"ltime"`
	List     []*Like    `json:"list"`
}

// ActSubject .
type ActSubject struct {
	ID        int64     `json:"id"`
	Oid       int64     `json:"oid"`
	Type      int       `json:"type"`
	State     int       `json:"state"`
	Stime     wocaoTime `json:"stime"`
	Etime     wocaoTime `json:"etime"`
	Ctime     wocaoTime `json:"ctime"`
	Mtime     wocaoTime `json:"mtime"`
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	ActURL    string    `json:"act_url"`
	Lstime    wocaoTime `json:"lstime"`
	Letime    wocaoTime `json:"letime"`
	Cover     string    `json:"cover" `
	Dic       string    `json:"dic"`
	Flag      int64     `json:"flag"`
	Uetime    wocaoTime `json:"uetime"`
	Ustime    wocaoTime `json:"ustime"`
	Level     int       `json:"level"`
	H5Cover   string    `json:"h5_cover"`
	Rank      int64     `json:"rank"`
	LikeLimit int       `json:"like_limit"`
}

// SubjectTotalStat .
type SubjectTotalStat struct {
	SumCoin int64 `json:"sum_coin"`
	SumFav  int64 `json:"sum_fav"`
	SumLike int64 `json:"sum_like"`
	SumView int64 `json:"sum_view"`
	Count   int   `json:"count"`
}

// VipActOrder .
type VipActOrder struct {
	ID             int64     `json:"id"`
	Mid            int64     `json:"mid"`
	OrderNo        string    `json:"order_no"`
	ProductID      string    `json:"product_id"`
	Ctime          wocaoTime `json:"ctime"`
	Mtime          wocaoTime `json:"mtime"`
	PanelType      string    `json:"panel_type"`
	Months         int       `json:"months"`
	AssociateState int       `json:"associate_state"`
}
