package relation

import (
	"context"
	"encoding/json"
	"go-common/library/net/metadata"
	"sort"
	"strconv"
)

// Pair ...
// 自定义map排序结构
type Pair struct {
	Key   int64
	Value int64
}

// Gray ...
// 自定义灰度策略结构
type Gray struct {
	Key   string
	Value int
}

// ErrLogStrut ...
// 自定义ErrLog结构
type ErrLogStrut struct {
	Code       int64
	Msg        string
	ErrDesc    string
	ErrType    string
	URLName    string
	RPCTimeout int64
	ChunkSize  int64
	ChunkNum   int64
	ErrorPtr   *error
}

// GrayRule ...
// 自定义灰度策略
type GrayRule struct {
	Name  string `json:"name"`
	Mark  string `json:"mark"`
	Value string `json:"value"`
}

// PairList ...
// 自定义灰度策略
type PairList []Pair

// GrayList ...
// 自定义灰度策略
type GrayList []Gray

// Swap
// 自定义排序
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Len
// 自定义排序
func (p PairList) Len() int { return len(p) }

// Less
// 自定义排序
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

// Swap
// 自定义排序
func (p GrayList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Len
// 自定义排序
func (p GrayList) Len() int { return len(p) }

// Less
// 自定义排序
func (p GrayList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// SortMap ...
// 自定义排序
func SortMap(input map[int64]int64) (sorted PairList) {

	p := make(PairList, len(input))
	i := 0
	for k, v := range input {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	sorted = p
	return
}

const (
	_androidBugBuildLeft  = 5332000
	_androidBugBuildRight = 5341000
)

// SortMapByValue ...
// 自定义排序
func SortMapByValue(m map[string]int) GrayList {
	p := make(GrayList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Gray{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

// RParseInt ...
// 转int
func RParseInt(inputStr string, defaultValue int64) (output int64) {
	if mid, err := strconv.ParseInt(inputStr, 10, 64); err == nil {
		output = mid
	} else {
		output = defaultValue
	}
	return
}

func RoleMap(role int8) (changeType int64) {
	switch role {
	case 0:
		{
			changeType = -1
		}
	case 1, 2:
		{
			changeType = 0
		}
	case 3, 4, 5, 6:
		{
			changeType = 1
		}
	default:
		{
			changeType = -1
		}
	}
	return
}

// CheckReturn ...
// 检查返回
func CheckReturn(err error, code int64, msg string, urlName string,
	rpcTimeout int64, chunkSize int64, chunkNum int64) (errLog *ErrLogStrut, success bool) {
	errInfo := ErrLogStrut{}
	errInfo.URLName = urlName
	errInfo.RPCTimeout = rpcTimeout
	errInfo.ChunkSize = chunkSize
	errInfo.ChunkNum = chunkNum
	success = true
	if err != nil {
		errInfo.Code = 1003000
		errInfo.Msg = ""
		errInfo.ErrDesc = "liveRpc调用失败"
		errInfo.ErrType = "LiveRpcFrameWorkCallError"
		errInfo.ErrorPtr = &err
		success = false
	} else if code != 0 {
		errInfo.Code = code
		errInfo.Msg = msg
		errInfo.ErrDesc = "调用直播服务" + urlName + "出错"
		errInfo.ErrType = "CallLiveRpcCodeError"
		success = false
	}
	errLog = &errInfo
	return
}

// App531ABTest ...
// ABTest
func App531ABTest(ctx context.Context, content string, build string, platform string) (grayType int64) {
	buildIntValue := RParseInt(build, 534000)
	if platform == "android" && buildIntValue > _androidBugBuildLeft && buildIntValue <= _androidBugBuildRight {
		grayType = 0
		return
	}
	if len(content) == 0 {
		grayType = 0
		return
	}
	resultMap := make(map[string]int64)
	resultMap["double_small_card"] = 0
	resultMap["card_not_auto_play"] = 1
	resultMap["card_auto_play"] = 2
	typeMap := make([]string, 0)
	mr := &[]GrayRule{}
	if err := json.Unmarshal([]byte(content), mr); err != nil {
		grayType = 0
		return
	}
	ruleArr := *mr
	scoreMap := make(map[string]int)

	for _, v := range ruleArr {
		scoreMap[v.Mark] = int(RParseInt(v.Value, 100))
	}
	sortedScore := SortMapByValue(scoreMap)
	scoreEnd := make([]int, 0)
	for _, v := range sortedScore {
		scoreEnd = append(scoreEnd, v.Value)
		typeMap = append(typeMap, v.Key)
	}
	score1 := scoreEnd[0]
	score2 := scoreEnd[0] + scoreEnd[1]
	score3 := 100
	section1 := make(map[int]bool)
	section2 := make(map[int]bool)
	section3 := make(map[int]bool)
	for section1Loop := 0; section1Loop < score1; section1Loop++ {
		section1[section1Loop] = true
	}
	for sectionLoop2 := score1; sectionLoop2 < score2; sectionLoop2++ {
		section2[sectionLoop2] = true
	}
	for sectionLoop3 := score2; sectionLoop3 < score3; sectionLoop3++ {
		section3[sectionLoop3] = true
	}
	mid := GetUIDFromHeader(ctx)
	result := int(mid % 100)
	if scoreEnd[0] != 0 {
		if _, exist := section1[result]; exist {
			grayType = resultMap[typeMap[0]]
			return
		}
	}
	if scoreEnd[1] != 0 {
		if _, exist := section2[result]; exist {
			grayType = resultMap[typeMap[1]]
			return
		}
	}
	if scoreEnd[2] != 0 {
		if _, exist := section3[result]; exist {
			grayType = resultMap[typeMap[2]]
			return
		}
	}
	grayType = 0
	return
}

// GetUIDFromHeader ...
// 获取uid
func GetUIDFromHeader(ctx context.Context) (uid int64) {
	midInterface, isUIDSet := metadata.Value(ctx, metadata.Mid).(int64) // 大多使用header里的mid解析, 框架已封装请求的header
	mid := int64(0)
	if isUIDSet {
		mid = midInterface
	}
	uid = mid
	return
}
