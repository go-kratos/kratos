package model

import "go-common/library/time"

// BlockedUserCard usr blocked info.
type BlockedUserCard struct {
	UID             int64  `json:"uid"`
	Uname           string `json:"uname"`
	Face            string `json:"face"`
	BlockedSum      int    `json:"blockedSum"`
	MoralBlockedSum int    `json:"moralBlockedSum"`
	MoralNum        int    `json:"moralNum"`
	BlockedStatus   int    `json:"blockedStatus"`
	BlockedForever  bool   `json:"blockedForever"`
	BlockedRestDay  int64  `json:"blockedRestDays"`
	AnsWerStatus    bool   `json:"answerStatus"`
	BlockedEndTime  int64  `json:"blockedEndTime"`
}

// BlockedAnnouncement blocked publish info.
type BlockedAnnouncement struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	SubTitle      string    `json:"subTitle"`
	PublishStatus uint8     `json:"-"`
	StickStatus   uint8     `json:"stickStatus"`
	Content       string    `json:"content"`
	URL           string    `json:"url"`
	Ptype         int8      `json:"ptype"`
	CTime         time.Time `json:"ctime"`
	MTime         time.Time `json:"mtime"`
}

// AnnounceList announce list.
type AnnounceList struct {
	List  []*BlockedAnnouncement `json:"list"`
	Count int64                  `json:"count"`
}
