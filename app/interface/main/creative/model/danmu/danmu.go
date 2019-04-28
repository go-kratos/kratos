package danmu

import (
	"go-common/library/time"
)

// AdvanceDanmu str
type AdvanceDanmu struct {
	ID        int64  `json:"id"`
	Cid       int64  `json:"cid"`
	Mid       int64  `json:"mid"`
	Aid       int64  `json:"aid"`
	Type      string `json:"type"`
	Mode      string `json:"mode"`
	UName     string `json:"uname"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Timestamp int64  `json:"timestamp"`
}

// DmList str
type DmList struct {
	List        []*MemberDM `json:"list"`
	Page        int64       `json:"page"`
	Size        int64       `json:"page_size"`
	TotalItems  int64       `json:"total_items"`
	TotalPages  int         `json:"TotalPages"`
	NormalCount int         `json:"normal_count"`
	SubCount    int         `json:"sub_count"`
	SpecCount   int         `json:"spec_count"`
}

// MemberDM str
type MemberDM struct {
	ID       int64     `json:"id"`
	FontSize int32     `json:"fontsize"`
	Color    string    `json:"color"`
	Mode     int32     `json:"mode"`
	Msg      string    `json:"msg"`
	VTitle   string    `json:"vtitle"`
	Oid      int64     `json:"oid"`
	Aid      int64     `json:"aid"`
	ArcTitle string    `json:"atitle"`
	Cover    string    `json:"cover"`
	Attrs    string    `json:"attrs"`
	Mid      int64     `json:"mid"`
	Playtime float64   `json:"playtime"`
	Pool     int32     `json:"pool"`
	State    int32     `json:"state"`
	Ctime    time.Time `json:"ctime"`
	Uname    string    `json:"uname"`
	Uface    string    `json:"uface"`
	Relation int       `json:"relation"`
	IsElec   int       `json:"is_elec"`
}

// Recent str
type Recent struct {
	ID       int64     `json:"id"`
	Aid      int64     `json:"aid"`
	Type     int32     `json:"type"`
	Oid      int64     `json:"oid"`
	Mid      int64     `json:"mid"`
	Msg      string    `json:"msg"`
	Cover    string    `json:"cover"`
	FontSize int32     `json:"font_size"`
	Color    string    `json:"color"`
	Attrs    string    `json:"attrs"`
	Mode     int32     `json:"mode"`
	Playtime float64   `json:"playtime"`
	Pool     int32     `json:"pool"`
	State    int32     `json:"state"`
	Title    string    `json:"title"` // oid所对应的稿件的标题
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
	Uname    string    `json:"uname"`
	Uface    string    `json:"uface"`
	Relation int       `json:"relation"`
	IsElec   int       `json:"is_elec"`
}

// DmRecent str
type DmRecent struct {
	List        []*Recent `json:"list"`
	Page        int64     `json:"page"`
	Size        int64     `json:"page_size"`
	TotalItems  int64     `json:"total_items"`
	TotalPages  int       `json:"TotalPages"`
	NormalCount int       `json:"normal_count"`
	SubCount    int       `json:"sub_count"`
	SpecCount   int       `json:"spec_count"`
}

// DmReport str
type DmReport struct {
	RpID       int64  `json:"rp_id"`
	DmInID     int64  `json:"dm_inid"`
	AID        int64  `json:"aid"`
	Pic        string `json:"pic"`
	ReportTime int64  `json:"reporttime"`
	Title      string `json:"title"`
	Reason     string `json:"reason"`
	DmID       int64  `json:"dmid"`
	DmIDStr    string `json:"dmid_str"`
	UpUID      int64  `json:"up_uid"`
	Content    string `json:"content"`
	UID        int64  `json:"uid"`
	UserName   string `json:"username"`
}

// DmArc str
type DmArc struct {
	Aid   int64  `json:"aid"`
	Title string `json:"title"`
}

// Pager str
type Pager struct {
	Total      int `json:"total"`
	Current    int `json:"current"`
	Size       int `json:"size"`
	TotalCount int `json:"total_count"`
}

// Apply str
type Apply struct {
	ID       int64   `json:"id"`
	IDStr    string  `json:"id_str"`
	AID      int64   `json:"aid"`
	CID      int64   `json:"cid"`
	Title    string  `json:"title"`
	ApplyUID int64   `json:"-"`
	Pic      string  `json:"pic"`
	Uname    string  `json:"uname"`
	Msg      string  `json:"msg"`
	Playtime float32 `json:"playtime"`
	Ctime    string  `json:"ctime"`
}

// ApplyListFromDM str
type ApplyListFromDM struct {
	Pager *Pager
	List  []*Apply
}

// ApplyList str
type ApplyList struct {
	Pager *Pager   `json:"pager"`
	List  []*Apply `json:"list"`
}

// ------------------- danmu2 upgrade -------------------//

// DMMember str
type DMMember struct {
	ID       int64     `json:"id"`
	Type     int32     `json:"type"`
	Aid      int64     `json:"aid"`
	Oid      int64     `json:"oid"`
	Mid      int64     `json:"mid"`
	MidHash  string    `json:"mid_hash"`
	Pool     int32     `json:"pool"`
	Attrs    string    `json:"attrs"`
	Progress int32     `json:"progress"`
	Mode     int32     `json:"mode"`
	Msg      string    `json:"msg"`
	State    int32     `json:"state"`
	FontSize int32     `json:"fontsize"`
	Color    string    `json:"color"`
	Ctime    time.Time `json:"ctime"`
	Uname    string    `json:"uname"`
	Title    string    `json:"title"`
}

// RecentPage str
type RecentPage struct {
	Pn    int64 `json:"num"`
	Ps    int64 `json:"size"`
	Total int64 `json:"total"`
}

// ResNewRecent str
type ResNewRecent struct {
	Result []*DMMember `json:"result"`
	Page   *RecentPage `json:"page"`
}

//SearchDMResult dm list
type SearchDMResult struct {
	Page struct {
		Num   int64 `json:"num"`
		Size  int64 `json:"size"`
		Total int64 `json:"total"`
	} `json:"page"`
	Result []*DMMember `json:"result"`
}

// SubtitleSubjectReply str
type SubtitleSubjectReply struct {
	AllowSubmit bool   `json:"allow"`
	Lan         string `json:"lan"`
	LanDoc      string `json:"lan_doc"`
}
