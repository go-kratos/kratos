package model

import (
	"encoding/json"
	"fmt"

	"go-common/app/admin/main/member/model/block"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// review state const.
const (
	ReviewStateWait = iota
	ReviewStatePass
	ReviewStateNoPass
	ReviewStateArchived
	ReviewStateQueuing = 10
)

// review property const.
const (
	ReviewProperty = iota
	ReviewPropertyFace
	ReviewPropertySign
	ReviewPropertyName
)

// all
var (
	AllReviewStates = []int8{
		ReviewStateWait,
		ReviewStatePass,
		ReviewStateNoPass,
		ReviewStateQueuing,
	}
)

// UserPropertyReview is
type UserPropertyReview struct {
	ID        int64      `json:"id" gorm:"column:id"`
	Mid       int64      `json:"mid" gorm:"column:mid"`
	Old       string     `json:"old" gorm:"column:old"`
	New       string     `json:"new" gorm:"column:new"`
	State     int8       `json:"state" gorm:"column:state"`
	Property  int8       `json:"property" gorm:"column:property"`
	Remark    string     `json:"remark" gorm:"column:remark"`
	Operator  string     `json:"operator" gorm:"column:operator"`
	IsMonitor bool       `json:"is_monitor" gorm:"column:is_monitor"`
	Extra     string     `json:"extra" gorm:"column:extra"`
	CTime     xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime     xtime.Time `json:"mtime" gorm:"column:mtime"`

	// 昵称，展示用
	Name       string             `json:"name" gorm:"-"`
	FaceReject int64              `json:"face_reject" gorm:"-"`
	Block      *block.BlockDetail `json:"block" gorm:"-"`
	Follower   int64              `json:"follower" gorm:"-"`
}

// Extra is.
type Extra struct {
	NickFree bool `json:"nick_free"`
}

// NickFree nick free.
func (r *UserPropertyReview) NickFree() bool {
	if len(r.Extra) == 0 {
		return false
	}
	ext := Extra{}
	if err := json.Unmarshal([]byte(r.Extra), &ext); err != nil {
		log.Error("Failed to unmarshal extra, userPropertyReview: %+v error: %v", r, err)
		return false
	}
	return ext.NickFree
}

// FaceCheckRes is.
type FaceCheckRes struct {
	Blood    float64 `json:"blood,omitempty"`
	Violent  float64 `json:"violent,omitempty"`
	Sex      float64 `json:"sex,omitempty"`
	Politics float64 `json:"politics,omitempty"`
}

// Valid is.
func (fcr *FaceCheckRes) Valid() bool {
	return fcr.Sex < 0.19 && fcr.Politics < 0.5 && fcr.Blood < 0.5 && fcr.Violent < 0.5
}

// String is.
func (fcr *FaceCheckRes) String() string {
	return fmt.Sprintf("Sex: %.4f, Politics: %.4f", fcr.Sex, fcr.Politics)
}

//BuildFaceURL buildFaceUrl.
func (r *UserPropertyReview) BuildFaceURL() {
	r.Old = BuildFaceURL(r.Old)
	r.New = BuildFaceURL(r.New)
}
