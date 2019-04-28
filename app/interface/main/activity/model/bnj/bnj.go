package bnj

import "go-common/app/service/main/archive/api"

// RewardTypes .
const (
	RewardTypePendant = "pendant"
	RewardTypeCoupon  = "coupon"
)

// Reward .
type Reward struct {
	Step       int
	Condition  int64
	RewardID   string
	RewardType string
	Expire     int64
}

// ResetMsg .
type ResetMsg struct {
	Mid int64 `json:"mid"`
	Ts  int64 `json:"ts"`
}

// PreviewInfo .
type PreviewInfo struct {
	ActID           int64   `json:"act_id"`
	SubID           int64   `json:"sub_id"`
	Info            []*Info `json:"info"`
	TimelinePic     string  `json:"timeline_pic"`
	H5TimelinePic   string  `json:"h5_timeline_pic"`
	GameCancel      int64   `json:"game_cancel"`
	RewardStep      []int64 `json:"reward_step"`
	HasRewardFirst  int     `json:"has_reward_first"`
	HasRewardSecond int     `json:"has_reward_second"`
}

// Timeline .
type Timeline struct {
	TimelinePic   string `json:"timeline_pic"`
	H5TimelinePic string `json:"h5_timeline_pic"`
	GameCancel    int64  `json:"game_cancel"`
	LikeCount     int64  `json:"like_count"`
}

// Info .
type Info struct {
	Nav      string   `json:"nav"`
	Pic      string   `json:"pic"`
	H5Pic    string   `json:"h5_pic"`
	Arc      *api.Arc `json:"arc"`
	Detail   string   `json:"detail"`
	H5Detail string   `json:"h5_detail"`
}
