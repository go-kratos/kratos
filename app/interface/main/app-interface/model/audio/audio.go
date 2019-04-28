package audio

import xtime "go-common/library/time"

type Audio struct {
	ID       int64      `json:"id"`
	Aid      int64      `json:"aid"`
	UID      int64      `json:"uid"`
	Title    string     `json:"title"`
	Cover    string     `json:"cover"`
	Author   string     `json:"author"`
	Schema   string     `json:"schema"`
	Duration int64      `json:"duration"`
	Play     int        `json:"play"`
	Reply    int        `json:"reply"`
	IsOff    int        `json:"isOff"`
	AuthType int        `json:"authType"`
	CTime    xtime.Time `json:"ctime"`
}

type FavAudio struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	ImgURL     string `json:"img_url"`
	RecordsNum int    `json:"records_num"`
	IsOpen     int    `json:"is_open"`
}

type UpperCert struct {
	Cert *struct {
		Type int    `json:"type,omitempty"`
		Desc string `json:"desc,omitempty"`
	} `json:"cert,omitempty"`
}

type Card struct {
	Type   int `json:"type,omitempty"`
	Status int `json:"status,omitempty"`
}

type Fav struct {
	Song    bool `json:"song,omitempty"`
	Menu    bool `json:"menu,omitempty"`
	PGCMenu bool `json:"pgc_menu,omitempty"`
}
