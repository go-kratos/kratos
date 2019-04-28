package search

import (
	"go-common/app/admin/main/videoup/model/manager"
	account "go-common/app/service/main/account/api"
)

// VideoResultData search video return struct
type VideoResultData struct {
	Page struct {
		Num   int `json:"num"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"page"`
	Result []*Video `json:"result"`
}

type ArchiveResultData struct {
	Page struct {
		Num   int `json:"num"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"page"`
	Result   []*Archive    `json:"result"`
	Tips     string        `json:"_tips"`
	MoniAids map[int64]int `json:"moni_aids"`
}

// CopyrightResultData search copyright return struct
type CopyrightResultData struct {
	Page struct {
		Num   int `json:"num"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"page"`
	Result []*Copyright `json:"result"`
}

// Video search return video item struct
type Video struct {
	ID            int64              `json:"id"`
	Aid           int64              `json:"aid"`
	Cid           int64              `json:"cid"`
	Vid           int64              `json:"vid"`
	ArcTitle      string             `json:"arc_title"`
	ArcState      int                `json:"arc_state"` //稿件状态。
	RelationState int                `json:"relation_state"`
	State         int                `json:"state"`
	Status        int                `json:"status"` //视频状态。如果archive_video_relation的state被删除，则此Status为-100；否则此Status为video表的status
	ArcTypeID     int64              `json:"arc_typeid"`
	ArcMid        int64              `json:"arc_mid"`
	ArcAuthor     string             `json:"arc_author"`
	ArcSendDate   string             `json:"arc_senddate"`
	Duration      int64              `json:"duration"`
	Filename      string             `json:"filename"`
	MTime         string             `json:"mtime"`
	TagID         int64              `json:"tag_id"`
	TagName       string             `json:"tag_name"`
	UserType      []int64            `json:"user_type"`
	UserGroup     []*manager.UpGroup `json:"user_group"`
	CTime         string             `json:"ctime"`
	VCTime        string             `json:"v_ctime"`
	VMTime        string             `json:"v_mtime"`
	XcodeState    int8               `json:"xcode_state"`
}
type Archive struct {
	ID        int64                `json:"id"`
	Mid       int64                `json:"mid"`
	Official  account.OfficialInfo `json:"official_verify"`
	TagNames  []string             `json:"tid_names"`
	Access    int16                `json:"access"`
	Attribute []int                `json:"attribute"`
	Attrs     []int                `json:"attrs"`
	State     int8                 `json:"state"`
	Author    string               `json:"author"`
	Cover     string               `json:"cover"`
	CTime     string               `json:"ctime"`
	MTime     string               `json:"mtime"`
	PubDate   string               `json:"pubtime"`
	Copyright int8                 `json:"copyright"`
	FlowID    int64                `json:"flow_id"`
	MissionID int64                `json:"mission_id"`
	OrderID   int64                `json:"order_id"`
	Round     int                  `json:"round"`
	Title     string               `json:"title"`
	Content   string               `json:"content"`
	TypeID    int64                `json:"typeid"`
	UpFrom    int8                 `json:"up_from"`
	UserType  []int64              `json:"user_type"`
	UserGroup []*manager.UpGroup2  `json:"user_group"`
}

// Copyright search return copyright item struct
type Copyright struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	OName    string `json:"oname"`
	AkaNames string `json:"aka_names"`
	Level    string `json:"level"`
	Avoid    string `json:"avoid"`
	Plan     string `json:"plan"`
	Desc     string `json:"description"`
	URL      string `json:"url"`
}

// VideoParams search video params
type VideoParams struct {
	Action       string `form:"action"`
	Status       string `form:"status"`
	ArcTitle     string `form:"arc_title"`
	ArcMids      string `form:"arc_mids"`
	Order        string `form:"order"`
	Sort         int8   `form:"sort_order"`
	Keywords     string `form:"keywords"`
	Aids         string `form:"aids"`
	Cids         string `form:"cids"`
	Vids         string `form:"vids"`
	TypeID       string `form:"typeid"`
	Filename     string `form:"filename"`
	TagID        string `form:"tag_id"`
	Pn           int    `form:"pn"`
	Ps           int    `form:"ps"`
	Xcode        string `form:"xcode_state"`
	UserType     string `form:"user_type"`
	OrderType    string `form:"order_type"`
	DurationFrom string `form:"duration_from"`
	DurationTo   string `form:"duration_to"`
	MonitorList  string `form:"monitor_list"`
}

// ArchiveParams search archive params
type ArchiveParams struct {
	TypeID      string `form:"typeid"`
	SpecialType string `form:"special_arctype"`
	Round       string `form:"round"`
	Aids        string `form:"aids"`
	Mids        string `form:"mids"`
	Pn          int    `form:"page"`
	Ps          int    `form:"pagesize"`
	OrderType   string `form:"order_type"`
	Keywords    string `form:"keywords"`
	KwFields    string `form:"kw_fields"`
	IsFirst     string `form:"is_first"`
	IsOrder     int8   `form:"execute_order"`
	State       string `form:"state"`
	Access      string `form:"access"`
	UpFroms     string `form:"up_froms"`
	PGCList     string `form:"pgc_list"`
	OrderId     string `form:"order_id"`
	Attr        string `form:"attribute"`
	//ChannelReview string `form:"channel_review"`
	//HotReview     string `form:"hot_review"`
	Review      string `form:"review"`
	ReviewState string `form:"review_state"`
	MissionID   string `form:"mission_id"`
	NoMission   string `form:"no_mission"`
	UserType    string `form:"user_type"`
	Copyright   string `form:"copyright"`
	Order       string `form:"order"`
	ScoreFirst  string `form:"score_first"` //是否按关键字匹配优先
	Sort        string `form:"sort_order"`
	MonitorList string `form:"monitor_list"`
}

// ArcPGCConfig
type ArcPGCConfig struct {
	UPFrom  []int8
	Rounds  []int8
	States  []int8
	InState bool
	Auth    string
}
