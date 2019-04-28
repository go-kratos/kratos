package archive

import (
	"time"
)

const (
	//TypeHotRecheck 热门回查
	TypeHotRecheck = 1
	//TypeChannelRecheck 频道回查
	TypeChannelRecheck = 0
	//TypeExcitationRecheck 激励回查
	TypeExcitationRecheck = 2
	//RecheckStateWait  待回查
	RecheckStateWait = int8(-1)
	//RecheckStateNoForbid 已回查，且没有禁止(热门)  已回查(频道)
	RecheckStateNoForbid = int8(0)
	//RecheckStateForbid 已回查，且禁止(热门)
	RecheckStateForbid = int8(1)

	//RecheckStateIgnore 被忽略不需要回查的状态
	RecheckStateIgnore = int8(-2)

	// FromListChannelReview 从频道回查列表提交的数据
	FromListChannelReview = "channel_review"
	// FromListHotReview 从热门回查列表提交的数据
	FromListHotReview = "hot_review"

	// FromListExcitation 从激励回查列表提交的数据
	FromListExcitation = "excitation_list"
)

var (
	_recheckTypes = map[int]string{
		//TypeChannelRecheck: "频道回查",
		TypeHotRecheck:        "热门回查",
		TypeExcitationRecheck: "激励回查",
	}
)

// Recheck archive recheck
type Recheck struct {
	ID     int64     `json:"id"`
	Type   int       `json:"type"`
	Aid    int64     `json:"aid"`
	UID    int64     `json:"uid"`
	State  int8      `json:"state"`
	Remark string    `json:"remark"`
	CTime  time.Time `json:"ctime"`
	MTime  time.Time `json:"mtime"`
}

//RecheckType get recheck type name
func RecheckType(tp int) (str string) {
	return _recheckTypes[tp]
}
