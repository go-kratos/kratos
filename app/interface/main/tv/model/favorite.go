package model

import (
	arcwar "go-common/app/service/main/archive/api"
)

// FormFav is the form validation for favorites display
type FormFav struct {
	AccessKey string `form:"access_key" validate:"required"`
	Pn        int    `form:"pn" default:"1"`
}

// ReqFav is request for favorites function
type ReqFav struct {
	MID int64
	Pn  int
}

// ToReq def.
func (f *FormFav) ToReq(mid int64) *ReqFav {
	return &ReqFav{
		MID: mid,
		Pn:  f.Pn,
	}
}

// FormFavAct is the form validation for favorite action
type FormFavAct struct {
	AccessKey string `form:"access_key" validate:"required"`
	AID       int64  `form:"aid" validate:"required"`
	Action    int    `form:"action" validate:"min=1,max=2"`
}

// ReqFavAct is request for favorites action ( add/del ) function
type ReqFavAct struct {
	MID    int64
	AID    int64 // resource id ( ugc avid )
	Action int   // 1=add,2=delete
}

// ToReq def.
func (f *FormFavAct) ToReq(mid int64) *ReqFavAct {
	return &ReqFavAct{
		MID:    mid,
		AID:    f.AID,
		Action: f.Action,
	}
}

// FavMList def.
type FavMList struct {
	Page struct {
		Num   int `json:"num"`
		Size  int `json:"size"`
		Count int `json:"count"`
	} `json:"page"`
	List []*arcwar.Arc `json:"list"`
}

// RespFavAct is response strure for favorite actions
type RespFavAct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
