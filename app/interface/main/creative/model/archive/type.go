package archive

import (
	"sort"
)

// RulesForApp  var
var RulesForApp = map[int16]int16{
	160: 1,  // 生活
	4:   2,  // 游戏
	129: 3,  // 舞蹈
	3:   4,  // 音乐
	155: 5,  // 时尚
	5:   6,  // 娱乐
	36:  7,  // 科技
	188: 8,  // 数码
	1:   9,  // 动画
	181: 10, // 影视
	119: 11, // 鬼畜
	167: 12, // 国创
	165: 13, // 广告
	177: 14, // 纪录片
	11:  15, // 电视剧
	13:  16, // 番剧
	23:  17, // 电影
}

// RulesForWeb var
var RulesForWeb = map[int16]int16{
	4:   1,  // 游戏
	160: 2,  // 生活
	5:   3,  // 娱乐
	181: 4,  // 影视
	3:   5,  // 音乐
	36:  6,  // 科技
	188: 7,  // 数码
	1:   8,  // 动画
	155: 9,  // 时尚
	129: 10, // 舞蹈
	13:  11, // 番剧
	177: 12, // 纪录片
	119: 13, // 鬼畜
	165: 14, // 广告
	167: 15, // 国创
	11:  16, // 电视剧
	23:  17, // 电影
}

// const client type values
const (
	WebType = int8(1)
	AppType = int8(2)
)

// Type is archive type.
type Type struct {
	ID        int16   `json:"id"`
	Lang      string  `json:"-"`
	Parent    int16   `json:"parent"`
	Name      string  `json:"name"`
	Desc      string  `json:"description"`
	Descapp   string  `json:"desc"`
	Count     int64   `json:"-"`              // top type count
	Original  string  `json:"intro_original"` // second type original
	IntroCopy string  `json:"intro_copy"`     // second type copy
	Notice    string  `json:"notice"`         // second type notice
	AppNotice string  `json:"-"`              // app notice
	CopyRight int8    `json:"copy_right"`
	Show      bool    `json:"show"`
	Rank      int16   `json:"-"`
	Children  []*Type `json:"children,omitempty"`
}

// ForbidSubTypesForAppAdd fn
// 不允许投稿的分区: 连载剧集->15,完结剧集->34,电视剧相关->128,电影相关->82
// 原不能投多P的分区，创作姬统一不允许投稿: 欧美电影->145,日本电影->146,国产电影->147,其他国家->83
func ForbidSubTypesForAppAdd(tid int16) bool {
	return tid == 15 || tid == 34 || tid == 128 || tid == 82 || tid == 145 || tid == 146 || tid == 147 || tid == 83
}

// ForbidTopTypesForAll fn 175=>ASMR
func ForbidTopTypesForAll(tid int16) bool {
	return tid == 175
}

// ForbidTopTypesForAppAdd fn
// 不允许APP投稿顶级分区: 纪录片->177,电视剧->11,番剧->13
func ForbidTopTypesForAppAdd(tid int16) bool {
	return tid == 177 || tid == 11 || tid == 13
}

// CopyrightForCreatorAdd fn 0,不限制；1,只允许自制；2,只允许转载
// 音乐-OP/ED/OST 54 番剧-完结动画 32 番剧-连载动画 33 番剧-资讯 51 电影-其他国家 83 电影-欧美电影 145 电影-日本电影 146 电影-国产电影 147  纪录片-人文历史 37
func CopyrightForCreatorAdd(tid int16) bool {
	return tid == 54 || tid == 32 || tid == 33 || tid == 51 || tid == 83 || tid == 145 || tid == 146 || tid == 147 || tid == 37
}

// SortRulesForTopTypes fn
func SortRulesForTopTypes(types []*Type, clientType int8) {
	var rules map[int16]int16
	if clientType == WebType {
		rules = RulesForWeb
	} else if clientType == AppType {
		rules = RulesForApp
	}
	for _, t := range types {
		if rank, ok := rules[t.ID]; ok {
			t.Rank = rank
		} else {
			t.Rank = 32767
		}
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i].Rank < types[j].Rank
	})
}
