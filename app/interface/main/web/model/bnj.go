package model

import (
	arcmdl "go-common/app/service/main/archive/api"
)

// Bnj2019 .
type Bnj2019 struct {
	*Bnj2019View
	Elec    *ElecShow         `json:"elec"`
	Related []*Bnj2019Related `json:"related"`
	ReqUser *ReqUser          `json:"req_user"`
}

// Bnj2019View .
type Bnj2019View struct {
	*arcmdl.Arc
	Pages []*arcmdl.Page `json:"pages"`
}

// Bnj2019Related .
type Bnj2019Related struct {
	*arcmdl.Arc
	Pages []*arcmdl.Page `json:"pages"`
}

// ReqUser req user.
type ReqUser struct {
	Attention bool  `json:"attention"`
	Favorite  bool  `json:"favorite"`
	Like      bool  `json:"like"`
	Dislike   bool  `json:"dislike"`
	Coin      int64 `json:"coin"`
}

// Timeline bnj timeline.
type Timeline struct {
	Name    string `json:"name"`
	Start   int64  `json:"start"`
	End     int64  `json:"end"`
	Cover   string `json:"cover"`
	H5Cover string `json:"h5_cover"`
}
