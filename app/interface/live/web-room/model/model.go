package model

//History 存档信息
type History struct {
	Text       string        `json:"text"`
	NickName   string        `json:"nickname"`
	UnameColor string        `json:"uname_color"`
	UID        int64         `json:"uid"`
	TimeLine   string        `json:"timeline"`
	Isadmin    int32         `json:"isadmin"`
	Vip        int           `json:"vip"`
	SVip       int           `json:"svip"`
	Medal      []interface{} `json:"medal"`
	Title      []interface{} `json:"title"`
	UserLevel  []interface{} `json:"user_level"`
	Rank       int32         `json:"rank"`
	Teamid     int64         `json:"teamid"`
	RND        string        `json:"rnd"`
	UserTitle  string        `json:"user_title"`
	GuardLevel int           `json:"guard_level"`
	Bubble     int64         `json:"bubble"`
}
