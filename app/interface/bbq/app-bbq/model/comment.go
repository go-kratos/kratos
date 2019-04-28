package model

import (
	user "go-common/app/service/bbq/user/api"
	xtime "go-common/library/time"
)

const (
	//DefaultCmType bbq评论区
	DefaultCmType = 23
)

// Int64Bytes implements the Scanner interface.
type Int64Bytes []int64

// Mstr Mstr
type Mstr []string

// ReplyCount /mcounts 返回字段
type ReplyCount struct {
	SubState int64 `json:"sub_state"`
	Count    int64 `json:"count"`
}

// ReplyList /reply 返回字段
type ReplyList struct {
	Page struct {
		Num    int `json:"num"`
		Size   int `json:"size"`
		Count  int `json:"count"`
		Acount int `json:"acount"`
	}
	Config struct {
		ShowEntry int `json:"showentry"`
		ShowAdmin int `json:"showadmin"`
	}
	Replies []*Reply `json:"replies,omitempty"`
	Hots    []*Reply `json:"hots,omitempty"`
	Upper   struct {
		MID int64  `json:"mid"`
		Top *Reply `json:"top,omitempty"`
	}
	Top       *Reply  `json:"top,omitempty"`
	Notice    *Notice `json:"notice,omitempty"`
	BlackList int16   `json:"blacklist"`
	Assist    int16   `json:"assist"`
}

// Reply .
type Reply struct {
	RpID      int64      `json:"rpid"`
	Oid       uint64     `json:"oid"`
	Type      int8       `json:"type"`
	Mid       int64      `json:"mid"`
	Root      int64      `json:"root"`
	Parent    int64      `json:"parent"`
	Dialog    int64      `json:"dialog"`
	Count     int        `json:"count"`
	RCount    int        `json:"rcount"`
	Floor     int        `json:"floor"`
	State     int8       `json:"state"`
	FansGrade int8       `json:"fansgrade"`
	Attr      int8       `json:"attr"`
	CTime     xtime.Time `json:"ctime"`
	MTime     xtime.Time `json:"-"`
	// string
	RpIDStr   string `json:"rpid_str"`
	RootStr   string `json:"root_str"`
	ParentStr string `json:"parent_str"`
	DialogStr string `json:"dialog_str"`
	// action count, from ReplyAction count
	Like   int  `json:"like"`
	Hate   int  `json:"-"`
	Action int8 `json:"action"`
	// member info
	Member *Member `json:"member,omitempty"`
	// other
	Content *Content `json:"content,omitempty"`
	Replies []*Reply `json:"replies,omitempty"`
	Assist  int      `json:"assist"`
}

// Content .
type Content struct {
	RpID    int64      `json:"-"`
	Message string     `json:"message"`
	Ats     Int64Bytes `json:"ats"`
	Topics  Mstr       `json:"topics"`
	IP      uint32     `json:"ipi"`
	Plat    int8       `json:"plat"`
	Device  string     `json:"device"`
	Version string     `json:"version"`
	CTime   xtime.Time `json:"-"`
	MTime   xtime.Time `json:"-"`
	// ats member info
	Members []*Info `json:"members,omitempty"`
}

// Info Info
type Info struct {
	Mid         string `json:"mid"`
	Name        string `json:"uname"`
	Sex         string `json:"sex"`
	Sign        string `json:"sign"`
	Avatar      string `json:"avatar"`
	Rank        string `json:"rank"`
	DisplayRank string `json:"DisplayRank"`
	LevelInfo   struct {
		Cur     int `json:"current_level"`
		Min     int `json:"current_min"`
		NowExp  int `json:"current_exp"`
		NextExp int `json:"next_exp"`
	} `json:"level_info"`
	Pendant        PendantInfo   `json:"pendant"`
	Nameplate      NameplateInfo `json:"nameplate"`
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Vip struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
	BInfo *user.UserBase `json:"b_info,omitempty"`
}

// Member .
type Member struct {
	Mid         string `json:"mid"`
	Name        string `json:"uname"`
	Sex         string `json:"sex"`
	Sign        string `json:"sign"`
	Avatar      string `json:"avatar"`
	Rank        string `json:"rank"`
	DisplayRank string `json:"DisplayRank"`
	LevelInfo   struct {
		CurrentLevel int `json:"current_level"`
		Min          int `json:"current_min"`
		NowExp       int `json:"current_exp"`
		NextExp      int `json:"next_exp"`
	} `json:"level_info"`
	Pendant        PendantInfo   `json:"pendant"`
	Nameplate      NameplateInfo `json:"nameplate"`
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Vip struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
	FansDetail *FansDetail    `json:"fans_detail,omitempty"`
	Following  int16          `json:"following"` //是否关注
	BInfo      *user.UserBase `json:"b_info,omitempty"`
}

// FansDetail .
type FansDetail struct {
	UID       int64  `json:"uid"`
	MedalID   int32  `json:"medal_id"`      //勋章id
	MedalName string `json:"medal_name"`    //勋章名称
	Score     int32  `json:"score"`         //当前总经验值
	Level     int8   `json:"level"`         //level等级
	Intimacy  int32  `json:"intimacy"`      //当前亲密度
	Status    int8   `json:"master_status"` //佩戴状态1:佩戴中0:未佩戴
	Received  int8   `json:"is_receive"`    //是否领取0:未领取1:已领取
}

// NameplateInfo .
type NameplateInfo struct {
	Nid        int    `json:"nid"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ImageSmall string `json:"image_small"`
	Level      string `json:"level"`
	Condition  string `json:"condition"`
}

// PendantInfo .
type PendantInfo struct {
	PID    int    `json:"pid"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Expire int    `json:"expire"`
}

//AddRes 添加评论返回字段
type AddRes struct {
	RpID        int64  `json:"rpid"`
	NeedCaptcha bool   `json:"need_captcha"`
	URL         string `json:"url"`
}

//CursorRes /reply/cusor response
type CursorRes struct {
	Assist    int64    `json:"assist"`
	BlackList int64    `json:"blacklist"`
	Replies   []*Reply `json:"replies,omitempty"`
	Upper     *Upper   `json:"upper,omitempty"`
	Cursor    *Cursor  `json:"cursor,omitempty"`
	Config    *Config  `json:"config,omitempty"`
	Hots      []*Reply `json:"hots,omitempty"`
	Notice    *Notice  `json:"notice,omitempty"`
	Top       *Top     `json:"top,omitempty"`
}

//SubCursorRes /reply/reply/cusor response
type SubCursorRes struct {
	Assist    int64   `json:"assist"`
	BlackList int64   `json:"blacklist"`
	Root      *Reply  `json:"root,omitempty"`
	Upper     *Upper  `json:"upper,omitempty"`
	Cursor    *Cursor `json:"cursor,omitempty"`
}

// Notice Notice
type Notice struct {
	ID         int64  `json:"id"`
	Plat       int8   `json:"-"`
	Condition  int8   `json:"-"`
	Build      int64  `json:"-"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Link       string `json:"link"`
	ClientType string `json:"-"`
}

//Top .
type Top struct {
	Admin *Reply `json:"admin,omitempty"`
	Upper *Reply `json:"upper,omitempty"`
}

//Config .
type Config struct {
	ShowEntry int16 `json:"showentry"`
	ShowAdmin int16 `json:"showadmin"`
}

//Cursor .
type Cursor struct {
	ACount         int64 `json:"all_count"`
	CursorRangeMax int64 `json:"max_id"`
	CursorRangeMin int64 `json:"min_id"`
	Size           int64 `json:"size"`
}

//Upper .
type Upper struct {
	MID int64 `json:"mid"`
}
