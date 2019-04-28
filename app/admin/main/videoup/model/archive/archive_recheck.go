package archive

import "time"

const (
	//TypeChannelRecheck 频道回查
	TypeChannelRecheck = 0 //频道回查
	// TypeHotRecheck 热门回查
	TypeHotRecheck = 1
	// TypeInspireRecheck 激励回查
	TypeInspireRecheck = 2
	//RecheckStateWait 频道回查待回查
	RecheckStateWait = -1 //待回查
	//RecheckStateDone 频道回查已回查
	RecheckStateDone = 0 //已回查
)

// Recheck archive recheck
type Recheck struct {
	ID     int64     `json:"id"`
	Type   int       `json:"type"`
	AID    int64     `json:"aid"`
	UID    int64     `json:"uid"`
	State  int8      `json:"state"`
	Remark string    `json:"remark"`
	CTime  time.Time `json:"ctime"`
	MTime  time.Time `json:"mtime"`
}
