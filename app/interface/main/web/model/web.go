package model

import (
	bcmdl "go-common/app/service/main/broadcast/api/grpc/v1"
)

// Coin add type.
const (
	CoinAddArcType  = 1
	CoinAddArtType  = 2
	CoinArcBusiness = "archive"
	CoinArtBusiness = "article"
)

var (
	// RankType rank type params
	RankType = map[int]string{
		1: "all",
		2: "origin",
		3: "rookie",
	}
	// DayType day params
	DayType = map[int]int{
		1:  1,
		3:  3,
		7:  7,
		30: 30,
	}
	// ArcType arc params type all:0 and recent:1
	ArcType = map[int]int{
		0: 0,
		1: 1,
	}
	// IndexDayType rank index day type
	IndexDayType = []int{
		1,
		3,
		7,
	}
	// OriType original or not
	OriType = []string{
		0: "",
		1: "_origin",
	}
	// AllType all or origin type
	AllType = []string{
		0: "all",
		1: "origin",
	}
	// TagIDs feedback tag ids
	TagIDs = []int64{
		300, //播放卡顿
		301, //进度条君无法调戏
		302, //画音不同步
		303, //弹幕无法加载/弹幕延迟
		304, //出现浮窗广告
		305, //无限小电视
		306, //黑屏
		307, //其他
		354, //校园网无法访问
	}
	// LimitTypeIDs view limit type id
	LimitTypeIDs = []int16{13, 32, 33, 94, 120}
	// RecSpecTypeName recommend data special type name
	RecSpecTypeName = map[int32]string{
		28: "原创",
		30: "V家",
		31: "翻唱",
		59: "演奏",
	}
	// LikeType thumbup like type
	LikeType = map[int8]string{
		1: "like",
		2: "like_cancel",
		3: "dislike",
		4: "dislike_cancel",
	}
	// NewListRid new list need more rids
	NewListRid = map[int32]int32{
		177: 37,
		23:  147,
		11:  185,
	}
	// DefaultServer  broadcst servers default value.
	DefaultServer = &bcmdl.ServerListReply{
		Domain:    "broadcast.chat.bilibili.com",
		TcpPort:   7821,
		WsPort:    7822,
		WssPort:   7823,
		Heartbeat: 30,
		Nodes:     []string{"broadcast.chat.bilibili.com"},
		Backoff: &bcmdl.Backoff{
			MaxDelay:  300,
			BaseDelay: 3,
			Factor:    1.8,
			Jitter:    0.3,
		},
		HeartbeatMax: 3,
	}
)

// CheckFeedTag check if tagID in TagIDs
func CheckFeedTag(tagID int64) bool {
	check := false
	for _, id := range TagIDs {
		if tagID == id {
			check = true
			break
		}
	}
	return check
}
