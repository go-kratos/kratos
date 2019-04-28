package archive

import xtime "go-common/library/time"

// Pconfig str
type Pconfig struct {
	ID   int64  `json:"id"`
	Tp   int8   `json:"type"`
	Name string `json:"name"`
}

// PorderArc str
type PorderArc struct {
	AID        int64      `json:"aid"`
	IndustryID int64      `json:"industry_id"`
	BrandID    int64      `json:"brand_id"`
	BrandName  string     `json:"brand_name"`
	Official   int8       `json:"official"`
	ShowType   string     `json:"show_type"`
	Advertiser string     `json:"advertiser"`
	Agent      string     `json:"agent"`
	State      int8       `json:"state"`
	ShowFront  int8       `json:"show_front"`
	Ctime      xtime.Time `json:"ctime"`
	Mtime      xtime.Time `json:"mtime"`
}
