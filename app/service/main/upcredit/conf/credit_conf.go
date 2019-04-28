package conf

import (
	"go-common/library/log"
	"strconv"
)

//CreditConf 信用分数配置
type CreditConf struct {
	CalculateConf *CreditScoreCalculateConfig
	ArticleRule   *ArticleRuleConf
}

//AfterLoad load之后进行一些计算和整理
func (c *CreditConf) AfterLoad() {
	c.CalculateConf.AfterLoad()
	c.ArticleRule.AfterLoad()
}

//CreditScoreCalculateConfig 分数计算配置
type CreditScoreCalculateConfig struct {
	// 时间衰减因子
	// [离今年的差值] = 权重值
	TimeWeight  map[string]int
	TimeWeight2 map[int]int
}

//AfterLoad after load
func (c *CreditScoreCalculateConfig) AfterLoad() {
	c.TimeWeight2 = make(map[int]int)
	for k, v := range c.TimeWeight {
		key, _ := strconv.Atoi(k)
		c.TimeWeight2[key] = v
	}
}

//StateMachineState 状态机数据
type StateMachineState struct {
	State  int
	Round  int
	Reason int
}

//ArticleRuleConf 稿件记分规则
type ArticleRuleConf struct {
	AcceptOptypeData []int
	RejectOpTypeData []int

	// [score] ->[ optype list ]
	OptypeScoreData map[string][]int

	AcceptOptypeMap map[int]struct{}
	RejectOptypeMap map[int]struct{}

	// [optype] -> score
	OptypeScoreMap      map[int]int
	InitState           StateMachineState
	ArticleMaxOpenCount int
}

//AfterLoad after load
func (a *ArticleRuleConf) AfterLoad() {
	a.AcceptOptypeMap = make(map[int]struct{})
	a.RejectOptypeMap = make(map[int]struct{})
	a.OptypeScoreMap = make(map[int]int)
	for _, v := range a.AcceptOptypeData {
		a.AcceptOptypeMap[v] = struct{}{}
	}

	for _, v := range a.RejectOpTypeData {
		a.RejectOptypeMap[v] = struct{}{}
	}

	for k, varr := range a.OptypeScoreData {
		var score, err = strconv.ParseInt(k, 10, 64)
		if err != nil {
			log.Error("score config wrong, k=%s", k)
			return
		}
		for _, v := range varr {
			a.OptypeScoreMap[v] = int(score)
		}
	}
}

//GetScore get score by type, opType, reason
func (a *ArticleRuleConf) GetScore(typ int, opType int, reason int) (score int) {
	var s, ok = a.OptypeScoreMap[opType]
	if !ok {
		s = 0
	}
	score = s
	return
}

//IsRejected is article reject
func (a *ArticleRuleConf) IsRejected(typ int, opType int, reason int) (res bool) {
	_, res = a.RejectOptypeMap[opType]
	return
}

//IsAccepted is article accepted
func (a *ArticleRuleConf) IsAccepted(typ int, opType int, reason int) (res bool) {
	_, res = a.AcceptOptypeMap[opType]
	return
}
