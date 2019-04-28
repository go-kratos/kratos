package model

// ReplyHistory reply history
type ReplyHistory struct {
	Page struct {
		Num   int `json:"num"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"page"`
	Records []*Record `json:"records"`
}

// Record record
type Record struct {
	ID      int       `json:"id"`
	Oid     int64     `json:"oid"`
	OidStr  string    `json:"oid_str"` // oid 前端会溢出改用 string
	Type    int64     `json:"type"`
	Floor   int       `json:"floor"`
	Like    int       `json:"like"`
	Rcount  int       `json:"rcount"`
	Mid     int64     `json:"mid"`
	State   int       `json:"state"`
	Message string    `json:"message"`
	Ctime   string    `json:"ctime"`
	Members []*Member `json:"members"`
	RecordAppend
}

// RecordAppend record append
type RecordAppend struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Member member
type Member struct {
	Mid   int64  `json:"mid"`
	Uname string `json:"uname"`
}
