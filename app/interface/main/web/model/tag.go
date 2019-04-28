package model

import (
	"go-common/app/interface/main/tag/model"
	v1 "go-common/app/service/main/archive/api"
)

const (
	// TagStateOK means normal state
	TagStateOK = 0
	// TagStateDeleted means tag was deleted
	TagStateDeleted = 1
	// TagStateBlocked means tag was blocked
	TagStateBlocked = 2
)

// TagAids .
type TagAids struct {
	Code  int     `json:"code"`
	Total int64   `json:"total"`
	Data  []int64 `json:"data"`
}

// TagDetail .
type TagDetail struct {
	Total int       `json:"total"`
	List  []*v1.Arc `json:"list"`
	*model.TagTop
}
