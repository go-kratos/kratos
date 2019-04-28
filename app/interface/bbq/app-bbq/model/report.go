package model

const (
	//TypeVideo video type
	TypeVideo = int16(10)
	//TypeComment comment type
	TypeComment = int16(20)
	//TypeDanmu danmu type
	TypeDanmu = int16(30)
	//TypeUser user type
	TypeUser = int16(40)
)

// ReportConfig .
type ReportConfig struct {
	Type    int16   `json:"type"`
	Reasons []int16 `json:"reasons,omitempty"`
}

// ReasonConfig .
type ReasonConfig struct {
	ReasonType int16  `json:"reason_type"`
	Name       string `json:"name"`
}

//Reports .
var Reports = []*ReportConfig{
	{
		Type:    10,
		Reasons: []int16{1, 2, 7, 3, 4, 100},
	},
	{
		Type:    20,
		Reasons: []int16{1, 2, 7, 3, 100},
	},
	{
		Type:    30,
		Reasons: []int16{1, 2, 7, 3, 100},
	},
	{
		Type:    40,
		Reasons: []int16{5, 6},
	},
}

//Reasons .
var Reasons = []*ReasonConfig{
	{
		ReasonType: 1,
		Name:       "违法违禁",
	},
	{
		ReasonType: 2,
		Name:       "色情",
	},
	{
		ReasonType: 3,
		Name:       "赌博诈骗",
	},
	{
		ReasonType: 4,
		Name:       "血腥暴力",
	},
	{
		ReasonType: 5,
		Name:       "昵称违规",
	},
	{
		ReasonType: 6,
		Name:       "头像违规",
	},
	{
		ReasonType: 7,
		Name:       "低俗",
	},
	{
		ReasonType: 100,
		Name:       "其他",
	},
}

// MapReasons map reasons
var MapReasons = map[int16]string{
	1:   "违法违禁",
	2:   "色情",
	3:   "赌博诈骗",
	4:   "血腥暴力",
	5:   "昵称违规",
	6:   "头像违规",
	7:   "低俗",
	100: "其他",
}

// BiliReasonsMap 主站评论举报类型映射, key bbq value bilibili
var BiliReasonsMap = map[int16]int16{
	1:   9,
	2:   2,
	3:   12,
	4:   0,
	5:   0,
	6:   0,
	7:   10,
	100: 0,
}
