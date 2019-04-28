package search

// Info struct
type Info struct {
	ID      int64 `json:"id"`
	OPID    int64 `json:"oper_id"`
	Status  int8  `json:"status" default:"-1"`
	PStatus int8  `json:"publish_status" default:"-1"`
}

// Publish struct
type Publish struct {
	ID       int64  `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	SubTitle string `json:"sub_title,omitempty"`
	ShowTime string `json:"show_time,omitempty"`
	OPID     int64  `json:"oper_id,omitempty"`
	PType    int8   `json:"ptype,omitempty"`
	Status   int8   `json:"status"`
}

// Case struct
type Case struct {
	ID        int64  `json:"id,omitempty"`
	MID       int64  `json:"mid,omitempty"`
	OPID      int64  `json:"oper_id,omitempty"`
	OType     int8   `json:"origin_type,omitempty"`
	Status    int8   `json:"status,omitempty"`
	CaseType  int8   `json:"case_type"`
	StartTime string `json:"start_time,omitempty"`
}

// Jury struct
type Jury struct {
	ID      int64  `json:"id,omitempty"`
	OPID    int64  `json:"oper_id,omitempty"`
	Expired string `json:"expired,omitempty"`
	Status  int8   `json:"status,omitempty"`
	Black   int8   `json:"black"`
}

// Opinion struct
type Opinion struct {
	ID    int64 `json:"id,omitempty"`
	OPID  int64 `json:"oper_id,omitempty"`
	State int8  `json:"state,omitempty"`
}

// Update struct
type Update struct {
	AppID string
	IP    string
	Data  interface{}
}

// Page struct
type Page struct {
	PN    int `json:"num"`
	PS    int `json:"size"`
	Total int `json:"total"`
}

// ResSearch result
type ResSearch struct {
	ID int64 `json:"id"`
}

// ReSearchData search result detail.
type ReSearchData struct {
	Order  string       `json:"order"`
	Sort   string       `json:"sort"`
	Page   *Page        `json:"page"`
	Result []*ResSearch `json:"result"`
}
