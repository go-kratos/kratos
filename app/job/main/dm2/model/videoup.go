package model

const (
	// RouteSecondRound 稿件二审消息
	RouteSecondRound = "second_round"
	// RouteAutoOpen = 稿件自动开放浏览
	RouteAutoOpen = "auto_open"
	// RouteForceSync 稿件强制同步
	RouteForceSync = "force_sync"
	// RouteDelayOpen 稿件定时开放浏览
	RouteDelayOpen = "delay_open"
	// VideoStatusOpen 视频开放浏览
	VideoStatusOpen = int32(0)
	//VideoXcodeHDFinish  高清转码完成
	VideoXcodeHDFinish = int32(4)
	//VideoXcodeFinish   视频转码
	VideoXcodeFinish = int32(2)
)

// VideoupMsg second round msg from VideoupBvc.
type VideoupMsg struct {
	Route string `json:"route"`
	Aid   int64  `json:"aid"`
}

// Archive archive info.
type Archive struct {
	Aid int64 `json:"aid"`
	Mid int64 `json:"mid"`
}

// Video video info.
type Video struct {
	Aid        int64 `json:"aid"`
	Cid        int64 `json:"cid"`
	Mid        int64 `json:"mid"`
	Duration   int64 `json:"duration"`
	Status     int32 `json:"status"`
	XCodeState int32 `json:"xcode_state"`
}
