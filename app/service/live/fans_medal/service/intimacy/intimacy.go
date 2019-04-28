package intimacy

// MaxLevel 大航海最高等级20级
var MaxLevel = 20

// levelConf 等级权益
type levelConf struct {
	levelUpExp int // 升级所需
	levelLimit int // 等级上限
	sumScore   int // 亲密度总和
}

// LevelConf 粉丝勋章等级配置
var LevelConf = []levelConf{
	{0, 0, -1},
	{201, 500, 0}, // 1
	{300, 500, 201},
	{500, 500, 501},
	{700, 500, 1001},
	{1000, 500, 1701}, //5
	{1500, 1000, 2701},
	{1600, 1000, 4201},
	{1700, 1000, 5801},
	{1900, 1000, 7501},
	{5500, 1000, 9401}, //10
	{10000, 1500, 14901},
	{10000, 1500, 24901},
	{10000, 1500, 34901},
	{15000, 1500, 44901},
	{40000, 1500, 59901}, //15
	{50000, 2000, 99901},
	{100000, 2000, 149901},
	{250000, 2000, 249901},
	{500000, 2000, 499901},
	{1000000, 2000, 999901}, //20
}

// LimitAddition 大航海不同等级时的加成
// 0 没有大航海，不加成
// 1 总督 3倍加成
// 2 提督 2倍加成
// 3 舰长 1.5倍加成
var LimitAddition = [4]float32{1.0, 3.0, 2.0, 1.5}

// Increaser 亲密度增加前的状态
type Increaser struct {
	currScore int  // 当前等级
	currLimit int  // 当前消耗的上限
	incrScore int  // 本次增加的亲密度
	skipLimit bool // 是否跳过每日上限的限制
}

// Result 亲密度增加后的结果
type Result struct {
	incr  int // 实际增加的亲密度
	score int // 增加亲密度之后
	limit int // 当前上限
	level int // 亲密度增加后的等级
}

// IncreaseResult 亲密度增加结果，包含是否超过最低上线的标志
type IncreaseResult struct {
	// 是否超过最低的上限,加上这个字段是为了提升性能
	// 如果某一次增加亲密度没有超过最低的上限,就不需要请求该用户的大航海信息
	// 减少一次网络请求
	passLimit bool
	retList   [4]Result
}

// Increase 根据当前 score 计算增加亲密度之后的结果
// 粉丝勋章服务中最为核心的逻辑
func Increase(medal *Increaser) (ret IncreaseResult) {
	// 用户在5上限是最大上限500，在6级是最大上限是1000
	// 如果他在当前上限为499的情况下增加两点亲密度
	// 按增加之前的等级计算上限，会浪费掉一点亲密度
	// 所以按增加之后的等级计算上限
	sum := medal.currScore + medal.incrScore
	newLevel, _ := GetLevel(sum)

	// 把几种大航海等级对应的上限加成结果都计算出来
	for i, addition := range LimitAddition {
		incr := medal.incrScore

		// 不消耗今日上限，直接按原数值加上
		if medal.skipLimit {
			ptr := &ret.retList[i]
			ptr.incr = incr
			ptr.limit = medal.currLimit
			ptr.score = medal.currScore + incr
			ptr.level, _ = GetLevel(ptr.score)

			continue
		}

		// 新等级对应的上限，考虑大航海加成情况
		maxLimit := addition * float32(LevelConf[newLevel].levelLimit)

		// 今日还能加多少亲密度
		leftLimit := int(maxLimit) - medal.currLimit

		// 避免极端情况下该值为负
		// 比如产品把大航海加成调低，某些土豪用户已有上限又较高
		// 该值就会为负
		if leftLimit < 0 {
			leftLimit = 0
		}

		// 累加不能超过上限
		if incr > leftLimit {
			ret.passLimit = true
			incr = leftLimit
		}

		ptr := &ret.retList[i]
		ptr.incr = incr
		ptr.limit = medal.currLimit + incr
		ptr.score = medal.currScore + incr
		ptr.level, _ = GetLevel(ptr.score)
	}

	return
}

// GetLevel 根据 score 计算等级和当前等级亲密度
func GetLevel(score int) (currLevel int, currIntimacy int) {
	currLevel = MaxLevel
	for currLevel > 0 {
		conf := LevelConf[currLevel]
		if score >= conf.sumScore {
			currIntimacy = score - conf.sumScore
			break
		}

		currLevel -= 1
	}
	return
}
