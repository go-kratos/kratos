package view

import (
	"encoding/json"
	"strconv"

	arcwar "go-common/app/service/main/archive/api"
)

// View view+relates
type View struct {
	*Static           // relate data
	ReqUser *ReqUser  `json:"req_user,omitempty"` // now user
	History *History  `json:"history,omitempty"`
	Relates []*Relate `json:"relates,omitempty"`
	PID     int32     `json:"category"` // father level partition ID
}

// ReqUser struct
type ReqUser struct {
	Attention int  `json:"attention"`
	Favorite  int8 `json:"favorite"`
	Like      int8 `json:"like"`
	Dislike   int8 `json:"dislike"`
	Coin      int8 `json:"coin"`
}

// Static .
type Static struct {
	*arcwar.Arc
	Pages []*Page `json:"pages,omitempty"`
}

// Page .
type Page struct {
	*arcwar.Page
	Metas []*Meta `json:"metas"`
}

// Meta .
type Meta struct {
	Quality int    `json:"quality"`
	Format  string `json:"format"`
	Size    int64  `json:"size"`
}

// Relate .
type Relate struct {
	Aid        int64         `json:"aid,omitempty"`
	Pic        string        `json:"pic,omitempty"`
	Title      string        `json:"title,omitempty"`
	Author     arcwar.Author `json:"owner,omitempty"`
	Stat       arcwar.Stat   `json:"stat,omitempty"`
	Duration   int64         `json:"duration,omitempty"`
	Goto       string        `json:"goto,omitempty"`
	Param      string        `json:"param,omitempty"`
	URI        string        `json:"uri,omitempty"`
	Rating     float64       `json:"rating,omitempty"`
	Reserve    string        `json:"reserve,omitempty"`
	From       string        `json:"from,omitempty"`
	Desc       string        `json:"desc,omitempty"`
	RcmdReason string        `json:"rcmd_reason,omitempty"`
	Badge      string        `json:"badge,omitempty"`
	Cid        int64         `json:"cid,omitempty"`
	// cm ad
	AdIndex      int             `json:"ad_index,omitempty"`
	CmMark       int             `json:"cm_mark,omitempty"`
	SrcID        int64           `json:"src_id,omitempty"`
	RequestID    string          `json:"request_id,omitempty"`
	CreativeID   int64           `json:"creative_id,omitempty"`
	CreativeType int64           `json:"creative_type,omitempty"`
	Type         int             `json:"type,omitempty"`
	Cover        string          `json:"cover,omitempty"`
	ButtonTitle  string          `json:"button_title,omitempty"`
	View         int             `json:"view,omitempty"`
	Danmaku      int             `json:"danmaku,omitempty"`
	IsAd         bool            `json:"is_ad,omitempty"`
	IsAdLoc      bool            `json:"is_ad_loc,omitempty"`
	AdCb         string          `json:"ad_cb,omitempty"`
	ShowURL      string          `json:"show_url,omitempty"`
	ClickURL     string          `json:"click_url,omitempty"`
	ClientIP     string          `json:"client_ip,omitempty"`
	Extra        json.RawMessage `json:"extra,omitempty"`
	Button       *Button         `json:"button,omitempty"`
	CardIndex    int             `json:"card_index,omitempty"`
}

// Button .
type Button struct {
	Title string `json:"title,omitempty"`
	URI   string `json:"uri,omitempty"`
}

// FromAv treatment
func (r *Relate) FromAv(a *arcwar.Arc, from string) {
	r.Aid = a.Aid
	r.Title = a.Title
	r.Pic = a.Pic
	r.Author = a.Author
	r.Stat = a.Stat
	r.Duration = a.Duration
	r.Cid = a.FirstCid
	r.Goto = GotoAv
	r.Param = strconv.FormatInt(a.Aid, 10)
	r.URI = FillURI(r.Goto, r.Param, AvHandler(a))
	r.From = from
}

// History struct
type History struct {
	Cid      int64 `json:"cid"`
	Progress int64 `json:"progress"`
}
