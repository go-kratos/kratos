package model

import (
	"go-common/library/time"
)

//CondIntNil cond int nil
const CondIntNil = -10516

// DmRecentResponse .
type DmRecentResponse struct {
	Page *Page       `json:"page"`
	Data []*DMMember `json:"result"`
}

// Recent recent dm
type Recent struct {
	ID       int64  `json:"id"`
	Type     int32  `json:"type"`
	Aid      int64  `json:"pid"`
	Oid      int64  `json:"oid"`
	Mid      int64  `json:"mid"`
	Pool     int32  `json:"pool"`
	Attr     int32  `json:"attr"`
	Progress int32  `json:"progress"`
	Mode     int32  `json:"mode"`
	Msg      string `json:"msg"`
	State    int32  `json:"state"`
	FontSize int32  `json:"fontsize"`
	Color    int32  `json:"color"`
	Ctime    string `json:"ctime"`
}

// DMMember dm struct used in member
type DMMember struct {
	ID       int64     `json:"id"`
	IDStr    string    `json:"id_str"`
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

// SearchDMParams dm search params
type SearchDMParams struct {
	Type         int32
	Oid          int64
	Keyword      string
	Mids         string
	Mode         string
	Pool         string
	Attrs        string
	ProgressFrom int64
	ProgressTo   int64
	CtimeFrom    string
	CtimeTo      string
	Pn           int64
	Ps           int64
	Sort         string
	Order        string
	State        string
}

// SearchRecentDMParam .
type SearchRecentDMParam struct {
	Type   int32
	UpMid  int64
	States []int32
	Ps     int
	Pn     int
	Sort   string
	Field  string
}

// SearchRecentDMResult .
type SearchRecentDMResult struct {
	Page   *Page     `json:"page"`
	Result []*Recent `json:"result"`
}

// SearchDMData dm meta data from search
type SearchDMData struct {
	Result []*struct {
		ID int64 `json:"id"`
	} `json:"result"`
	Page *SearchPage
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

// UptSearchDMState update search dm state
type UptSearchDMState struct {
	ID    int64  `json:"id"`
	Oid   int64  `json:"oid"`
	Type  int32  `json:"type"`
	State int32  `json:"state"`
	Mtime string `json:"mtime"`
}

// UptSearchDMPool update search dm pool
type UptSearchDMPool struct {
	ID    int64  `json:"id"`
	Oid   int64  `json:"oid"`
	Type  int32  `json:"type"`
	Pool  int32  `json:"pool"`
	Mtime string `json:"mtime"`
}

// UptSearchDMAttr update search dm attr
type UptSearchDMAttr struct {
	ID         int64   `json:"id"`
	Oid        int64   `json:"oid"`
	Type       int32   `json:"type"`
	Attr       int32   `json:"attr"`
	Mtime      string  `json:"mtime"`
	AttrFormat []int64 `json:"attr_format"`
}
