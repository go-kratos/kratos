package archive

import (
	"go-common/library/time"
)

// Archive is archive model.
type Archive struct {
	Aid          int64     `json:"aid"`
	Mid          int64     `json:"mid"`
	TypeID       int16     `json:"tid"`
	HumanRank    int       `json:"-"`
	Title        string    `json:"title"`
	Author       string    `json:"-"`
	Cover        string    `json:"cover"`
	RejectReason string    `json:"reject_reason"`
	Tag          string    `json:"tag"`
	Duration     int64     `json:"duration"`
	Copyright    int8      `json:"copyright"`
	Desc         string    `json:"desc"`
	MissionID    int64     `json:"mission_id"`
	Round        int8      `json:"-"`
	Forward      int64     `json:"-"`
	Attribute    int32     `json:"attribute"`
	Access       int16     `json:"-"`
	State        int8      `json:"state"`
	Source       string    `json:"source"`
	NoReprint    int32     `json:"no_reprint"`
	OrderID      int64     `json:"order_id"`
	Dynamic      string    `json:"dynamic"`
	DTime        time.Time `json:"dtime"`
	PTime        time.Time `json:"ptime"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"-"`
	Tnames       []string  `json:"tid_names"`
}

// Addit is archive addit info
type Addit struct {
	Aid           int64  `json:"aid"`
	MissionID     int64  `json:"mission_id"`
	UpFrom        int8   `json:"up_from"`
	FromIP        int64  `json:"from_ip"`
	Source        string `json:"source"`
	OrderID       int64  `json:"order_id"`
	RecheckReason string `json:"recheck_reason"`
	RedirectURL   string `json:"redirect_url"`
	FlowID        int64  `json:"flow_id"`
	Advertiser    string `json:"advertiser"`
	DescFormatID  int64  `json:"desc_format_id"`
	Dynamic       string `json:"dynamic"`
	InnerAttr     int64  `json:"inner_attr"`
}

// Delay is archive delay info
type Delay struct {
	Aid   int64
	Mid   int64
	State int16
	DTime time.Time
}

// Type is archive type info
type Type struct {
	ID   int16  `json:"id"`
	PID  int16  `json:"pid"`
	Name string `json:"name"`
	Desc string `json:"description"`
}

//ChannelInfo channel info
type ChannelInfo struct {
	CheckBack int32      `json:"check_back"`
	Channels  []*Channel `json:"channels"`
}

//Channel channe & tag hit rule
type Channel struct {
	TID         int64    `json:"tid"`        //频道id
	Tname       string   `json:"tname"`      //频道名称
	HitRules    []string `json:"hit_rules"`  //命中的频道规则
	HitTagNames []string `json:"hit_tnames"` //命中频道的所有tag名称
}
