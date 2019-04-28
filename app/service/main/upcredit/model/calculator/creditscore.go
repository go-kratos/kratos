package calculator

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-common/app/service/main/upcredit/common/fsm"
	"go-common/app/service/main/upcredit/conf"
	"go-common/app/service/main/upcredit/mathutil"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/library/log"
	xtime "go-common/library/time"
	"math"
	"sort"
	"time"
)

/*
	计算过程
	1.读取up的信用日志记录
	2.按照时间维度，分为1年内、2年内...的记录，并计算每年的分数（并不是自然年，是以当前时间为起点的365天内算1年）,calcUpCreditScoreByYear
	3.总分为每年分数的加权平均（没有记录的分数为0的年份不记入加权计算）,calcWeightedCreditScore, -> 标准化 calcNormalizedCreditScore
	4.写入对应的分数db
	5.全部完成后写入task info db
	其他：
	1.每个分数计算分为绝对分与相对分的加权平均
*/

const (
	//ScoreRange score range, max - min
	ScoreRange = float32(2000)
	//MaxScore max score
	MaxScore = float32(1000)
	//WeightRelative relative weight
	WeightRelative = float32(0.5)
)

//CreditScoreCalculator score calculator
type CreditScoreCalculator struct {
	CreditScoreOutputChan chan<- *upcrmmodel.UpScoreHistory
}

//New create
func New(channel chan<- *upcrmmodel.UpScoreHistory) *CreditScoreCalculator {
	var c = &CreditScoreCalculator{
		CreditScoreOutputChan: channel,
	}
	return c
}

type upLog struct {
	Mid        int64
	CreditLogs []upcrmmodel.SimpleCreditLog
}

type creditStat struct {
	RejectArcNum  int64
	AcceptArcNum  int64
	AbsoluteScore float32
	RelativeScore float32
	TotalScore    float32
}

func (c *creditStat) onLogResult(e *fsm.Event, l *upcrmmodel.SimpleCreditLog, a *ArticleStateMachine) {

	var score = 0
	var ruleConfig = conf.CreditConfig.ArticleRule
	switch e.Dst {
	case StateClose:
		if ruleConfig.IsRejected(l.Type, l.OpType, l.Reason) {
			c.RejectArcNum++
			fmt.Printf("article rejected, state=%d, round=%d\n", a.State, a.Round)
		} else {
			fmt.Printf("article close, but not rejected, state=%d, round=%d\n", a.State, a.Round)
		}

		score = ruleConfig.GetScore(l.Type, l.OpType, l.Reason)
		c.AbsoluteScore += float32(score)
		fmt.Printf("reject, score=%d\n", score)
	case StateOpen:
		c.AcceptArcNum++
		score = conf.CreditConfig.ArticleRule.GetScore(l.Type, l.OpType, l.Reason)
		c.AbsoluteScore += float32(score)
		fmt.Printf("accept, score=%d\n", score)
	}

}

//CalcRelativeScore relative score
func (c *creditStat) CalcRelativeScore() {
	c.RelativeScore = relativeScore(c.RejectArcNum, c.AcceptArcNum)
}

//CalcTotalScore total score
func (c *creditStat) CalcTotalScore() {
	// 绝对分数限制在 (-SCORE_RANGE/2, SCORE_RANGE/2)
	c.AbsoluteScore = ScoreRange * float32(1/(1+math.Exp(-float64(c.AbsoluteScore/300)))-0.5)
	c.TotalScore = WeightRelative*c.RelativeScore + (1-WeightRelative)*c.AbsoluteScore
}

//AppendLog append log
func (u *upLog) AppendLog(log upcrmmodel.SimpleCreditLog) {
	u.CreditLogs = append(u.CreditLogs, log)
}

//SortLog sort log by ctime asc
func (u *upLog) SortLog() {
	sort.SliceStable(u.CreditLogs, func(i, j int) bool {
		return u.CreditLogs[i].CTime < u.CreditLogs[j].CTime
	})
}

func relativeScore(rejectArcNum int64, acceptedArcnum int64) float32 {
	var total = rejectArcNum + acceptedArcnum
	if total == 0 {
		return 0
	}

	var factor = float32(1.0)
	if total <= 5 {
		factor = 0.2
	} else if total <= 10 {
		factor = 0.5
	}
	const middle = 0.6
	var accRatio = float32(acceptedArcnum)/float32(total) - middle
	if accRatio < -0.5 {
		accRatio = -0.5
	}
	return factor * accRatio * ScoreRange
}

//

const (
	//Day day
	Day = time.Hour * 24 / time.Second
	//Month month
	Month = 30 * Day
	//Year year
	Year = 365 * Day
)

func calcUpCreditScoreByYear(uplog *upLog, currentTime time.Time) (yearStat map[int]*creditStat) {
	// 以时间为key，creditStat为value的分类
	yearStat = make(map[int]*creditStat)
	var articleMachine = make(map[int]*ArticleStateMachine)
	var now = currentTime.Unix()
	for _, l := range uplog.CreditLogs {
		var difftime = now - int64(l.CTime)
		// 不计算来自未来的数据
		if difftime < 0 {
			continue
		}
		var index = int(difftime / int64(Year))
		v := getOrCreateCreditScore(index, yearStat)
		for _, r := range RuleList {
			r(l, v, articleMachine)
		}
	}

	// 计算每年的分数
	for k, v := range yearStat {
		v.CalcRelativeScore()
		v.CalcTotalScore()
		log.Info("score for mid:%d, [%d]s=%+v", uplog.Mid, k, v)
	}

	return
}

func calcWeightedCreditScore(yearStat map[int]*creditStat) (score float32) {
	// 每年的分数加权平均，如果当年分数为0，那不进行加权平均
	var totalWeight = float32(0)
	var totalScore = float32(0)
	for diff, weight := range conf.CreditConfig.CalculateConf.TimeWeight2 {
		s, o := yearStat[diff]
		if !o {
			continue
		}
		totalWeight += float32(weight)
		totalScore += float32(s.TotalScore) * float32(weight)
	}
	log.Info("total score=%f, total weight=%f", totalScore, totalWeight)
	if !mathutil.FloatEquals(totalWeight, 0) {
		score = totalScore / totalWeight
	}
	return
}

func calcNormalizedCreditScore(weightedScore float32) float32 {
	return (weightedScore/ScoreRange + 0.5) * MaxScore
}

func getOrCreateArticleFSM(aid int, m map[int]*ArticleStateMachine) (artFsm *ArticleStateMachine) {
	artFsm, ok := m[aid]
	if !ok {
		artFsm = CreateArticleStateMachineWithInitState()
		m[aid] = artFsm
	}
	return
}
func getOrCreateCreditScore(index int, m map[int]*creditStat) (c *creditStat) {
	c, ok := m[index]
	if !ok {
		c = new(creditStat)
		m[index] = c
	}
	return
}

//OverAllStatistic score statistics
type OverAllStatistic struct {
	MaxScore     int
	SectionCount int
	// 分数段人数统计， map[分数段:0~9]int人数
	ScoreSection map[int]int
	sectionScore int
}

//NewOverAllStatistic create empty
func NewOverAllStatistic(maxScore int, sectionCount int) *OverAllStatistic {
	if maxScore <= 0 || sectionCount <= 0 {
		panic(errors.New("max score or sction count must > 0"))
	}
	return &OverAllStatistic{
		MaxScore:     maxScore,
		SectionCount: sectionCount,
		ScoreSection: map[int]int{},
		sectionScore: maxScore / sectionCount,
	}
}

//AddScore add score
/*params:
score, 分数
exceptScore，不统计的分数，主要是默认分数，不进行统计
*/
func (s *OverAllStatistic) AddScore(score int, exceptScore int) {
	if score == exceptScore {
		return
	}
	var section = score / s.sectionScore
	if section > s.SectionCount-1 {
		section = s.SectionCount - 1
	}
	s.ScoreSection[section]++
}

// GetScore no record will return default int(0)
func (s *OverAllStatistic) GetScore(section int) int {
	return s.ScoreSection[section]
}

//CalcLogTable calculate all table
func (c *CreditScoreCalculator) CalcLogTable(tableNum int, currentDate time.Time, overall *OverAllStatistic) (err error) {
	var crmdb, e = gorm.Open("mysql", conf.Conf.DB.UpcrmReader.DSN)
	err = e
	if e != nil {
		log.Error("fail to open crm db, for table=%d", tableNum)
		return
	}
	crmdb.LogMode(true)
	defer crmdb.Close()
	var startTime = time.Now()
	var upLogMap = make(map[int64]*upLog)
	var lastID uint
	var limit = 1000
	var tableName = fmt.Sprintf("credit_log_%02d", tableNum)
	log.Info("table[%s] start load player's data", tableName)
	var total = 0
	for {
		var users []upcrmmodel.SimpleCreditLog
		e = crmdb.Table(tableName).Where("id > ?", lastID).Limit(limit).Find(&users).Error
		if e != nil {
			err = e
			log.Error("fail to get users from db, err=%v", e)
			break
		}
		// 加入到列表中
		for _, l := range users {
			up, ok := upLogMap[l.Mid]
			if !ok {
				up = &upLog{}
				up.Mid = l.Mid
				upLogMap[l.Mid] = up
			}
			lastID = l.ID
			up.AppendLog(l)
		}

		var thisCount = len(users)
		total += thisCount
		if thisCount < limit {
			log.Info("table[%s] total read record, num=%d", tableName, thisCount)
			break
		}
	}
	//crmdb.Close()
	if err != nil {
		log.Error("table[%s] error happen, exit calc, err=%v", tableName, err)
		return
	}
	log.Info("table[%s] start calculate player's data, total mid=%d", tableName, len(upLogMap))

	var date = currentDate
	for _, v := range upLogMap {
		v.SortLog()
		var yearStat = calcUpCreditScoreByYear(v, date)
		var weightedScore = calcWeightedCreditScore(yearStat)
		var finalScore = calcNormalizedCreditScore(weightedScore)
		log.Info("mid=%d, weightscore=%f, finalscore=%f", v.Mid, weightedScore, finalScore)
		c.writeUpCreditScore(v.Mid, finalScore, date, crmdb, overall)
	}
	var elapsed = time.Since(startTime)

	log.Info("table[%s] finish calculate player's data, total mid=%d, total logs=%d, duration=%s, avg=%0.2f/s", tableName, len(upLogMap), total, elapsed, float64(total)/elapsed.Seconds())

	return
}

func (c *CreditScoreCalculator) writeUpCreditScore(mid int64, score float32, date time.Time, crmdb *gorm.DB, overall *OverAllStatistic) {
	if c.CreditScoreOutputChan == nil {
		log.Error("output chan is nil, fail to output credit score")
		return
	}
	if overall != nil {
		overall.AddScore(int(score), 500)
	}
	var creditScore = &upcrmmodel.UpScoreHistory{
		Mid:          mid,
		Score:        int(score),
		ScoreType:    upcrmmodel.ScoreTypeCredit,
		GenerateDate: xtime.Time(date.Unix()),
		CTime:        xtime.Time(time.Now().Unix()),
	}
	c.CreditScoreOutputChan <- creditScore
	log.Info("output credit score, mid=%d", mid)
}
