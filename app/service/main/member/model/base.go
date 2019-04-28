package model

import (
	"fmt"
	"math/rand"
)

// consts
const (
	URLNoFace           = "http://static.hdslb.com/images/member/noface.gif"
	ActUpdateByAdmin    = "updateByAdmin"
	ActUpdatePersonInfo = "updatePersonInfo"
	ActUpdateFace       = "updateFace"
	ActUpdateUname      = "updateUname"
	ActBlockUser        = "blockUser"
	CertNO              = -1      // 未认证
	DefaultRank         = 5000    // default rank
	DefaultTime         = -28800  // default time
	DefaultMoral        = 7000    // default moral
	MaxMoral            = 10000   // max moral
	CacheKeyBase        = "bs_%d" // key of baseInfo
)

// RandFaceURL get face URL
func (b *BaseInfo) RandFaceURL() {
	if b.Face == "" {
		b.Face = URLNoFace
		return
	}
	b.Face = fmt.Sprintf("http://i%d.hdslb.com%s", rand.Int63n(3), b.Face)
}

// SexStr get sex str
func (b *BaseInfo) SexStr() string {
	switch b.Sex {
	case 0:
		return "保密"
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return "保密"
	}
}

// NotifyInfo notify info.
type NotifyInfo struct {
	Uname   string `json:"uname"`
	Mid     int64  `json:"mid"`
	Type    string `json:"type"`
	NewName string `json:"newName"`
	Action  string `json:"action"`
}

// Equal is.
func (of *OfficialInfo) Equal(cof *OfficialInfo) bool {
	return of.Role == cof.Role && of.Title == cof.Title && of.Desc == cof.Desc
}

// BaseExp exp and base info.
type BaseExp struct {
	*BaseInfo
	*LevelInfo
}

// Member is the full information within member-service.
type Member struct {
	*BaseInfo
	*LevelInfo
	*OfficialInfo
}
