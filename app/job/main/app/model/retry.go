package model

import (
	"go-common/app/job/main/app/model/space"
	xtime "go-common/library/time"
)

const (
	ActionUpView           = "upView"
	ActionUpStat           = "upStat"
	ActionUpContribute     = "upContribute"
	ActionUpContributeAid  = "upContributeAid"
	ActionUpViewContribute = "upViewContribute"
	ActionUpAccount        = "upAccount"
)

type Retry struct {
	Action string `json:"action,omitempty"`
	Data   struct {
		Mid    int64         `json:"mid,omitempty"`
		Aid    int64         `json:"aid,omitempty"`
		Attrs  *space.Attrs  `json:"attrs,omitempty"`
		Items  []*space.Item `json:"item,omitempty"`
		Time   xtime.Time    `json:"time,omitempty"`
		IP     string        `json:"ip,omitempty"`
		Action string        `json:"action,omitempty"`
	} `json:"data,omitempty"`
}
