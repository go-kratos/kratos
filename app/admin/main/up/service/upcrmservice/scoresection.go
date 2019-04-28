package upcrmservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-common/app/admin/main/up/dao/upcrm"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/jinzhu/gorm"
)

func convertValueList(score upcrmmodel.ScoreSectionHistory) []upcrmmodel.ScoreSection {
	var values []int
	values = append(values, score.Section0)
	values = append(values, score.Section1)
	values = append(values, score.Section2)
	values = append(values, score.Section3)
	values = append(values, score.Section4)
	values = append(values, score.Section5)
	values = append(values, score.Section6)
	values = append(values, score.Section7)
	values = append(values, score.Section8)
	values = append(values, score.Section9)
	var list []upcrmmodel.ScoreSection
	var total = 0
	for i, v := range values {
		var newsection = upcrmmodel.ScoreSection{
			Section: i,
			Value:   v,
		}
		list = append(list, newsection)
		total += v
	}

	if total > 0 {
		for i := range list {
			list[i].Percent = int(float32(list[i].Value) * 10000.0 / float32(total))
		}
	}

	return list
}

func generateScoreQueryXAxis(num int) (axis []string) {
	if num == 0 {
		return
	}
	axis = append(axis, fmt.Sprintf("%d-%d", 0, 100))
	for i := 1; i < num; i++ {
		axis = append(axis, fmt.Sprintf("%d-%d", i*100+1, (i+1)*100))
	}
	return
}

//ScoreQuery query score
func (s *Service) ScoreQuery(context context.Context, arg *upcrmmodel.ScoreQueryArgs) (result upcrmmodel.ScoreQueryResult, err error) {
	var today = time.Now()
	var yesterday, e = s.crmdb.GetLastHistory(arg.ScoreType)
	result = upcrmmodel.NewEmptyScoreQueryResult()
	err = e
	if err == ecode.NothingFound {
		log.Info("no record found in db, arg=%+v", arg)
		err = nil
		return
	}
	if err != nil {
		log.Error("fail get last history, arg=%+v, err=%+v", arg, err)
		return
	}
	yesterdayData, e := s.crmdb.ScoreQueryHistory(arg.ScoreType, yesterday)
	err = e
	// orm.go 里修改了gorm.ErrRecordNotFound！！
	if err != nil {
		log.Error("fail query score history, err=%+v", err)
		return
	}
	result.XAxis = generateScoreQueryXAxis(10)
	result.YAxis = convertValueList(yesterdayData)
	var comparedate time.Time
	// 获取对比数据
	switch arg.CompareType {
	case upcrmmodel.CompareTypeNothing:
		break
	case upcrmmodel.CompareType7day:
		comparedate = yesterday.AddDate(0, 0, -7)
	case upcrmmodel.CompareType30day:
		comparedate = yesterday.AddDate(0, 0, -30)
	case upcrmmodel.CompareTypeMonthFirstDay:
		y, m, _ := today.Date()
		comparedate = time.Date(y, m, 1, 0, 0, 0, 0, today.Location())
	}
	if !comparedate.IsZero() {
		var compareData, e = s.crmdb.ScoreQueryHistory(arg.ScoreType, comparedate)
		err = e
		// orm.go 里修改了gorm.ErrRecordNotFound！！
		if err != nil {
			log.Error("fail query score history, err=%+v", err)
			return
		}
		result.CompareAxis = convertValueList(compareData)
	}
	log.Info("result=%+v", result)
	return
}

func calcScoreInfo(datamap map[int8]map[time.Time]upcrmmodel.UpScoreHistory, stype int8, todate time.Time, fromdate time.Time) (info upcrmmodel.ScoreInfo) {
	var typeMap map[time.Time]upcrmmodel.UpScoreHistory
	var ok bool
	if typeMap, ok = datamap[stype]; !ok {
		log.Error("score type not found, type=%d", stype)
		return
	}
	var currentData upcrmmodel.UpScoreHistory
	if currentData, ok = typeMap[todate]; !ok {
		log.Error("score type for todate not exist, todate=%v", todate)
		return
	}
	info.Current = currentData.Score
	var lastScore = 0
	if lastData, ok := typeMap[fromdate]; ok {
		lastScore = lastData.Score
	}

	info.DiffLastDay = info.Current - lastScore
	return
}

func generateDataMap(scoreHistory []upcrmmodel.UpScoreHistory) map[int8]map[time.Time]upcrmmodel.UpScoreHistory {
	var dataMap = map[int8]map[time.Time]upcrmmodel.UpScoreHistory{}
	for _, v := range scoreHistory {
		var typeMap map[time.Time]upcrmmodel.UpScoreHistory
		var ok bool
		if typeMap, ok = dataMap[v.ScoreType]; !ok {
			typeMap = map[time.Time]upcrmmodel.UpScoreHistory{}
			dataMap[v.ScoreType] = typeMap
		}
		typeMap[GetDateStamp(v.GenerateDate.Time())] = v
	}
	return dataMap
}

func getDataFromMap(dataMap map[int8]map[time.Time]upcrmmodel.UpScoreHistory, scoreType int, date time.Time) (data upcrmmodel.UpScoreHistory, ok bool) {
	var mapDate, o = dataMap[int8(scoreType)]
	ok = o
	if !ok {
		return
	}
	data, ok = mapDate[date]
	return
}

//ScoreQueryUp query up score
func (s *Service) ScoreQueryUp(context context.Context, arg *upcrmmodel.ScoreQueryUpArgs) (result upcrmmodel.ScoreQueryUpResult, err error) {

	var now = time.Now()
	var todate = GetDateStamp(now)
	if arg.Date != "" {
		todate, err = time.ParseInLocation(upcrmmodel.TimeFmtDate, arg.Date, now.Location())
		if err != nil {
			log.Error("fail to parse date, date=%s", arg.Date)
			return
		}
	}
	var latestDate, e = s.crmdb.GetLatestUpScoreDate(arg.Mid, upcrm.ScoreTypePr, todate)
	if e != nil && e != gorm.ErrRecordNotFound {
		err = e
		return
	}
	todate = latestDate

	var fromdate = GetDateStamp(todate.AddDate(0, 0, -1))
	log.Info("query up's score all")
	var scoreHistory []upcrmmodel.UpScoreHistory
	scoreHistory, err = s.crmdb.QueryUpScoreHistory(arg.Mid, []int{upcrm.ScoreTypePr, upcrm.ScoreTypeQuality, upcrm.ScoreTypeCredit}, fromdate, todate)
	if err != nil {
		log.Error("query up score error, err=%+v", err)
		return
	}
	var dataMap = generateDataMap(scoreHistory)

	// 计算数据
	result.PrScore = calcScoreInfo(dataMap, upcrm.ScoreTypePr, todate, fromdate)
	result.QualityScore = calcScoreInfo(dataMap, upcrm.ScoreTypeQuality, todate, fromdate)
	result.CreditScore = calcScoreInfo(dataMap, upcrm.ScoreTypeCredit, todate, fromdate)
	result.Date = xtime.Time(todate.Unix())
	log.Info("score history ok req=%+v, result=%+v", arg, result)
	return
}

//ScoreQueryUpHistory query up history
func (s *Service) ScoreQueryUpHistory(context context.Context, arg *upcrmmodel.ScoreQueryUpHistoryArgs) (result upcrmmodel.ScoreQueryUpHistoryResult, err error) {
	var now = time.Now()
	var todate = now
	if arg.Date != "" {
		todate, err = time.ParseInLocation(upcrmmodel.TimeFmtDate, arg.Date, now.Location())
		if err != nil {
			log.Error("parse time fail, err=%+v", err)
			return
		}
	}
	switch arg.Day {
	case 7, 30, 90:
		break
	default:
		arg.Day = 7
	}
	latestDate, e := s.crmdb.GetLatestUpScoreDate(arg.Mid, upcrm.ScoreTypePr, todate)
	if e != nil && e != gorm.ErrRecordNotFound {
		err = e
		return
	}
	todate = latestDate

	var fromdate = todate.AddDate(0, 0, 1-arg.Day)
	var types []int
	switch arg.ScoreType {
	case 0:
		types = append(types, upcrm.ScoreTypePr, upcrm.ScoreTypeQuality, upcrm.ScoreTypeCredit)
	case upcrm.ScoreTypePr, upcrm.ScoreTypeQuality, upcrm.ScoreTypeCredit:
		types = append(types, arg.ScoreType)
	default:
		err = errors.New("score type not support")
		return
	}
	scoreHistory, e := s.crmdb.QueryUpScoreHistory(arg.Mid, types, fromdate, todate)
	err = e
	if err != nil {
		log.Error("query up score error, err=%+v", err)
		return
	}
	var dataMap = generateDataMap(scoreHistory)
	// 产生历史列表数据
	var dateSeries []xtime.Time
	var origDateSeries []time.Time
	for start := fromdate; !start.After(todate); start = start.AddDate(0, 0, 1) {
		var onlyDate = GetDateStamp(start)
		origDateSeries = append(origDateSeries, onlyDate)
		dateSeries = append(dateSeries, xtime.Time(onlyDate.Unix()))
	}
	//result.ScoreData = []
	for _, t := range types {
		var typehistory upcrmmodel.ScoreHistoryInfo
		typehistory.Type = t
		typehistory.Date = dateSeries
		for _, date := range origDateSeries {
			// 如果没有找到，就用默认的score = 0
			var score, _ = getDataFromMap(dataMap, t, date)
			typehistory.Score = append(typehistory.Score, score.Score)
		}
		result.ScoreData = append(result.ScoreData, typehistory)
	}
	log.Info("query up history sucessful, arg=%+v, result=%+v", arg, result)
	return
}
