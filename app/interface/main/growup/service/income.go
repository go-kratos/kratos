package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"go-common/app/interface/main/growup/model"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_video = iota
	_audio
	_column
	_bgm
	_up
)

var (
	_layout = "2006-01-02"
)

// GetUpCharge get up daily sum(charge) last 30 days
func (s *Service) GetUpCharge(c context.Context, mid int64, t time.Time) (total int, err error) {
	begin := t.AddDate(0, 0, -30).Format(_layout)
	incs, err := s.dao.GetUpDailyCharge(c, mid, begin)
	if err != nil {
		return
	}
	for _, inc := range incs {
		total += inc
	}
	total = int(float64(total) * float64(0.6))
	return
}

// ArchiveIncome get archive income by mid
func (s *Service) ArchiveIncome(c context.Context, mid int64, typ, page, size, all int) (data map[string]interface{}, err error) {
	redisKey := fmt.Sprintf("growup-archive-income:%d+%d+%d+%d+%d", typ, mid, all, page, size)
	data, err = s.dao.GetIncomeCache(c, redisKey)
	if err != nil {
		log.Error("s.dao.GetIncomeCache error(%v)", err)
		return
	}
	if data != nil {
		return
	}

	data, err = s.archiveIncome(c, mid, typ, page, size, all)
	if err != nil {
		log.Error("s.archiveIncome error(%v)", err)
		return
	}

	err = s.dao.SetIncomeCache(c, redisKey, data)
	if err != nil {
		log.Error("s.dao.SetIncomeCache error(%v)", err)
	}
	return
}

// archiveIncome get archive income by mid
func (s *Service) archiveIncome(c context.Context, mid int64, typ, page, size, all int) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	if page == 0 {
		page = 1
	}
	start, end := (page-1)*size, page*size

	date := time.Now().AddDate(0, 0, -2)
	startMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local).Format(_layout)
	if all == 1 {
		startMonth = "2017-12-01"
	}
	var (
		archives []*model.ArchiveIncome
		total    int
	)
	switch typ {
	case _video:
		archives, total, err = s.listAvIncome(c, mid, startMonth, date.Format(_layout), start, end)
	case _audio:
	case _column:
		archives, total, err = s.listColumnIncome(c, mid, startMonth, date.Format(_layout), start, end)
	case _bgm:
		archives, total, err = s.listBgmIncome(c, mid, startMonth, date.Format(_layout), start, end)
	}
	if err != nil {
		log.Error("s.archiveIncome error(%v)", err)
		return nil, err
	}

	data["data"] = archives
	data["total_count"] = total
	data["page"] = page
	return
}

func (s *Service) listAvIncome(c context.Context, mid int64, startTime, endTime string, start, end int) (archiveList []*model.ArchiveIncome, total int, err error) {
	avs, err := s.dao.ListAvIncome(c, mid, startTime, endTime)
	if err != nil {
		log.Error("s.dao.ListAvIncome error(%v)", err)
		return
	}
	return s.handleArchiveIncome(c, avs, start, end, _video)
}

func (s *Service) listColumnIncome(c context.Context, mid int64, startTime, endTime string, start, end int) (archiveList []*model.ArchiveIncome, total int, err error) {
	columns, err := s.dao.ListColumnIncome(c, mid, startTime, endTime)
	if err != nil {
		log.Error("s.dao.ListColumnIncome error(%v)", err)
		return
	}
	return s.handleArchiveIncome(c, columns, start, end, _column)
}

func (s *Service) listBgmIncome(c context.Context, mid int64, startTime, endTime string, start, end int) (archiveList []*model.ArchiveIncome, total int, err error) {
	bgms, err := s.dao.ListBgmIncome(c, mid, startTime, endTime)
	if err != nil {
		log.Error("s.dao.ListBgmIncome error(%v)", err)
		return
	}
	return s.handleArchiveIncome(c, bgms, start, end, _bgm)
}

func calArchiveIncome(archives []*model.ArchiveIncome, avBMap map[int64]struct{}) (archiveMap map[int64]*model.ArchiveIncome) {
	endTime := time.Now().AddDate(0, 0, -2)
	archiveMap = make(map[int64]*model.ArchiveIncome)
	archMonthIncome := make(map[int64]int64)
	for _, arch := range archives {
		if _, ok := avBMap[arch.ArchiveID]; ok {
			continue
		}
		archiveDate := arch.Date.Time()
		// cal month income
		if archiveDate.Month() == endTime.Month() {
			archMonthIncome[arch.ArchiveID] += arch.Income
		}
		if archiveDate.Format(_layout) != endTime.Format(_layout) {
			arch.Income = 0
		}
		if old, ok := archiveMap[arch.ArchiveID]; !ok {
			archiveMap[arch.ArchiveID] = arch
		} else {
			if old.Date < arch.Date {
				archiveMap[arch.ArchiveID] = arch
			}
		}
		archiveMap[arch.ArchiveID].MonthIncome = archMonthIncome[arch.ArchiveID]
	}
	return
}

func (s *Service) handleArchiveIncome(c context.Context, archives []*model.ArchiveIncome, start, end, typ int) (archiveList []*model.ArchiveIncome, total int, err error) {
	archiveList = make([]*model.ArchiveIncome, 0)
	if len(archives) == 0 {
		return
	}
	aIDMap := make(map[int64]struct{})
	aIDList := []int64{}
	for _, arch := range archives {
		if _, ok := aIDMap[arch.ArchiveID]; !ok {
			aIDMap[arch.ArchiveID] = struct{}{}
			aIDList = append(aIDList, arch.ArchiveID)
		}
	}
	avBMap, err := s.dao.ListAvBlackList(c, aIDList, typ)
	if err != nil {
		log.Error("s.dao.ListAvBlackList error(%v)", err)
		return
	}
	avsMap := calArchiveIncome(archives, avBMap)
	for _, av := range avsMap {
		archiveList = append(archiveList, av)
	}

	sort.Slice(archiveList, func(i, j int) bool {
		if archiveList[i].Income == archiveList[j].Income {
			if archiveList[i].MonthIncome == archiveList[j].MonthIncome {
				return archiveList[i].TotalIncome > archiveList[j].TotalIncome
			}
			return archiveList[i].MonthIncome > archiveList[j].MonthIncome
		}
		return archiveList[i].Income > archiveList[j].Income
	})

	total = len(archiveList)
	if end > total {
		end = total
	}
	if start >= total || start > end {
		return
	}
	archiveList = archiveList[start:end]
	avIDs := []int64{}
	for _, av := range archiveList {
		avIDs = append(avIDs, av.ArchiveID)
	}
	titles := make(map[int64]string)
	switch typ {
	case _video:
		titles, err = s.getAvTitle(c, avIDs)
		if err != nil {
			log.Error("s.getAvTitle error(%v)", err)
			return
		}
	case _column:
		titles, err = s.getColumnTitle(c, avIDs)
		if err != nil {
			log.Error("s.getColumnTitle error(%v)", err)
			return
		}
	case _bgm:
		titles, err = s.getBgmTitle(c, avIDs)
		if err != nil {
			log.Error("s.getBgmTitle error(%v)", err)
			return
		}
	}

	icons, err := s.getAvIcon(c, avIDs)
	if err != nil {
		log.Error("s.getAvIcon error(%v)", err)
		return
	}
	breachs, err := s.getAvBreach(c, avIDs, typ)
	if err != nil {
		log.Error("s.getAvBreach error(%v)", err)
		return
	}
	for _, av := range archiveList {
		av.Title = titles[av.ArchiveID]
		av.Icon = icons[av.ArchiveID]
		av.Breach = breachs[av.ArchiveID]
	}
	return
}

func (s *Service) getColumnTitle(c context.Context, avs []int64) (titles map[int64]string, err error) {
	return s.dao.GetColumnTitle(c, avs)
}

func (s *Service) getBgmTitle(c context.Context, avs []int64) (titles map[int64]string, err error) {
	return s.dao.GetBgmTitle(c, avs)
}

func (s *Service) getAvTitle(c context.Context, avs []int64) (titles map[int64]string, err error) {
	req, err := http.NewRequest("GET", s.conf.Host.ArchiveURI, nil)
	if err != nil {
		log.Error("http.NewRequest error(%v)", err)
		return
	}
	q := req.URL.Query()
	q.Add("aids", xstr.JoinInts(avs))
	q.Add("appkey", s.conf.AppConf.Key)
	now := time.Now().Unix()
	q.Add("ts", strconv.FormatInt(now, 10))

	sign := q.Encode() + s.conf.AppConf.Secret
	q.Add("sign", fmt.Sprintf("%x", md5.Sum([]byte(sign))))

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("http.DefaultClient.Do error(%v)", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("ioutil.ReadAll error(%v)", err)
		return
	}

	res := model.ArchiveRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Error("json.Unmarshal body %s error(%v)", string(body), err)
		return
	}

	titles = make(map[int64]string)
	for _, archive := range res.Data {
		titles[archive.AID] = archive.Title
	}
	return
}

func (s *Service) getAvIcon(c context.Context, avs []int64) (acM map[int64]string, err error) {
	if len(avs) == 0 {
		return
	}
	activeM, err := s.dao.ListActiveInfo(c, avs)
	if err != nil {
		log.Error("s.dao.ListActiveInfo error(%v)", err)
		return
	}

	tagIDM := make(map[int64]struct{})
	for _, tagID := range activeM {
		tagIDM[tagID] = struct{}{}
	}
	tagIDList := make([]int64, 0)
	for tagID := range tagIDM {
		tagIDList = append(tagIDList, tagID)
	}

	if len(tagIDList) == 0 {
		return
	}
	tagIconM, err := s.dao.ListTagInfo(c, tagIDList)
	if err != nil {
		log.Error("s.dao.ListTagInfo error(%v)", err)
		return
	}

	acM = make(map[int64]string)
	for avID, tagID := range activeM {
		if _, ok := tagIconM[tagID]; ok {
			acM[avID] = tagIconM[tagID].Icon
		}
	}
	return
}

func (s *Service) getAvBreach(c context.Context, avs []int64, typ int) (breachs map[int64]*model.AvBreach, err error) {
	if len(avs) == 0 {
		return
	}
	return s.dao.GetAvBreachs(c, avs, typ)
}

// UpSummary summary up income
func (s *Service) UpSummary(c context.Context, mid int64) (data interface{}, err error) {
	redisKey := fmt.Sprintf("growup-up-summary:%d", mid)
	res, err := s.dao.GetIncomeCache(c, redisKey)
	if err != nil {
		log.Error("s.dao.GetIncomeCache error(%v)", err)
		return
	}
	if res != nil {
		data = res["data"]
		return
	}

	data, err = s.upSummary(c, mid)
	if err != nil {
		log.Error("s.upSummary error(%v)", err)
		return
	}

	err = s.dao.SetIncomeCache(c, redisKey, map[string]interface{}{"data": data})
	if err != nil {
		log.Error("s.dao.SetIncomeCache error(%v)", err)
	}
	return
}

func (s *Service) upSummary(c context.Context, mid int64) (data interface{}, err error) {
	summary := new(struct {
		DayIncome      string `json:"day_income"`
		Date           string `json:"date"`
		Income         string `json:"income"`
		TotalIncome    string `json:"total_income"`
		WaitWithdraw   string `json:"wait_withdraw"`
		BreachMoney    string `json:"breach_money"`
		UnwithdrawDate string `json:"unwithdraw_date"`
	})
	summary.DayIncome, summary.Income, summary.TotalIncome, summary.WaitWithdraw, summary.BreachMoney = "0", "0", "0", "0", "0"

	now := time.Now().AddDate(0, 0, -2)
	nowMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	count, err := s.dao.GetUpIncomeCount(c, now.Format(_layout))
	if err != nil {
		log.Error("s.dao.GetUpIncomeCount error(%v)", err)
		return
	}
	if count <= 0 {
		summary.DayIncome, summary.Income, summary.TotalIncome, summary.WaitWithdraw, summary.BreachMoney = "-1", "-1", "-1", "-1", "-1"
		data = summary
		return
	}

	upIncomes, err := s.dao.ListUpIncome(c, mid, "up_income", nowMonth.Format(_layout), now.AddDate(0, 0, 1).Format(_layout))
	if err != nil {
		log.Error("s.dao.ListUpIncome error(%v)", err)
		return
	}
	var monthIncome, lastDayIncome, dayIncome int64
	for _, up := range upIncomes {
		if up.Date.Time().Format(_layout) == now.Format(_layout) {
			dayIncome = up.Income
		}
		if up.Date.Time().Format(_layout) == now.AddDate(0, 0, 1).Format(_layout) {
			lastDayIncome = up.Income
		}
		monthIncome += up.Income
	}

	upAccount, err := s.dao.ListUpAccount(c, mid)
	if err != nil {
		log.Error("s.dao.ListUpAccount error(%v)", err)
		return
	}
	if upAccount == nil {
		data = summary
		return
	}

	breachs, err := s.dao.ListAvBreach(c, mid, nowMonth.Format(_layout), time.Now().Format(_layout))
	if err != nil {
		log.Error("s.dao.ListAvBreach error(%v)", err)
		return
	}
	var breachMoney int64
	for _, b := range breachs {
		breachMoney += b.Money
	}

	summary.DayIncome = fmt.Sprintf("%.2f", fromYuanToFen(dayIncome))
	summary.BreachMoney = fmt.Sprintf("%.2f", fromYuanToFen(breachMoney))
	summary.TotalIncome = fmt.Sprintf("%.2f", fromYuanToFen(upAccount.TotalIncome-lastDayIncome))
	summary.Date = now.Format(_layout)
	wdv, err := time.Parse("2006-01", upAccount.WithdrawDateVersion)
	if err != nil {
		log.Error("time.Parse error(%v)", err)
		return
	}
	summary.UnwithdrawDate = time.Date(wdv.Year(), wdv.Month()+1, 1, 0, 0, 0, 0, time.Local).Format("2006-01")

	// 如果T-1的month不等于T-2的month 当月新增不需要减去那一天的收入:lastDayIncome = 0
	if now.AddDate(0, 0, 1).Month() == now.Month() {
		monthIncome -= lastDayIncome
	}
	summary.Income = fmt.Sprintf("%.2f", fromYuanToFen(monthIncome))

	// 当月未提现，待结算不能减去昨天收入，当月已提现，需要减去昨日收入
	preMonth := time.Date(nowMonth.Year(), nowMonth.Month()-1, 1, 0, 0, 0, 0, time.Local).Format("2006-01")
	if preMonth != upAccount.WithdrawDateVersion || now.AddDate(0, 0, 1).Month() != now.Month() {
		lastDayIncome = 0
	}
	summary.WaitWithdraw = fmt.Sprintf("%.2f", fromYuanToFen(upAccount.TotalUnwithdrawIncome-lastDayIncome))
	data = summary
	return
}

func fromYuanToFen(income int64) float64 {
	return float64(income) / float64(100)
}

// ArchiveSummary get archive summary
func (s *Service) ArchiveSummary(c context.Context, typ int, mid int64) (data interface{}, err error) {
	redisKey := fmt.Sprintf("growup-archive-summary:%d+%d", typ, mid)
	res, err := s.dao.GetIncomeCache(c, redisKey)
	if err != nil {
		log.Error("s.dao.GetIncomeCache error(%v)", err)
		return
	}
	if res != nil {
		data = res["data"]
		return
	}

	data, err = s.archiveSummary(c, typ, mid)
	if err != nil {
		log.Error("s.archiveSummary error(%v)", err)
		return
	}

	err = s.dao.SetIncomeCache(c, redisKey, map[string]interface{}{"data": data})
	if err != nil {
		log.Error("s.dao.SetIncomeCache error(%v)", err)
	}
	return
}

func (s *Service) archiveSummary(c context.Context, typ int, mid int64) (data interface{}, err error) {
	date := time.Now().AddDate(0, 0, -2)
	startMonth := xtime.Time(time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local).Unix())
	upIncomes, err := s.dao.ListUpIncome(c, mid, "up_income", "2017-12-01", date.Format(_layout))
	if err != nil {
		log.Error("s.dao.ListUpIncome error(%v)", err)
		return
	}
	summary := new(struct {
		DayIncome   int64  `json:"day_income"`
		Date        string `json:"date"`
		MonthIncome int64  `json:"income"`
		TotalIncome int64  `json:"total_income"`
		Breach      int64  `json:"breach_money"`
	})
	if len(upIncomes) == 0 {
		data = summary
		return
	}
	sort.Slice(upIncomes, func(i, j int) bool {
		return upIncomes[i].Date > upIncomes[j].Date
	})
	summary.Date = date.Format(_layout)
	if upIncomes[0].Date.Time().Format(_layout) == summary.Date {
		switch typ {
		case _video:
			summary.DayIncome = upIncomes[0].AvIncome
		case _column:
			summary.DayIncome = upIncomes[0].ColumnIncome
		case _bgm:
			summary.DayIncome = upIncomes[0].BgmIncome
		}
	}

	var breachType []int64
	for _, up := range upIncomes {
		if up.Date >= startMonth {
			switch typ {
			case _video:
				summary.MonthIncome += up.AvIncome
			case _column:
				summary.MonthIncome += up.ColumnIncome
			case _bgm:
				summary.MonthIncome += up.BgmIncome
			}
		}
		switch typ {
		case _video:
			breachType = []int64{0}
			summary.TotalIncome += up.AvIncome
		case _column:
			breachType = []int64{2}
			summary.TotalIncome += up.ColumnIncome
		case _bgm:
			breachType = []int64{3}
			summary.TotalIncome += up.BgmIncome
		}
	}

	breach, err := s.dao.GetAvBreachByType(c, mid, "2017-12-01", time.Now().Format(_layout), breachType)
	if err != nil {
		log.Error("s.dao.GetAvBreachByType error(%v)", err)
		return
	}
	for d, money := range breach {
		if d >= startMonth {
			summary.Breach += money
		}
		summary.TotalIncome -= money
	}
	if summary.TotalIncome < 0 {
		summary.TotalIncome = 0
	}
	data = summary
	return
}

// ArchiveDetail cal archive detail
func (s *Service) ArchiveDetail(c context.Context, typ int, archiveID int64) (data interface{}, err error) {
	redisKey := fmt.Sprintf("growup-archive-detail:%d+%d", typ, archiveID)
	res, err := s.dao.GetIncomeCache(c, redisKey)
	if err != nil {
		log.Error("s.dao.GetIncomeCache error(%v)", err)
		return
	}
	if res != nil {
		data = res["data"]
		return
	}

	data, err = s.archiveDetail(c, typ, archiveID)
	if err != nil {
		log.Error("s.archiveDetail error(%v)", err)
		return
	}

	err = s.dao.SetIncomeCache(c, redisKey, map[string]interface{}{"data": data})
	if err != nil {
		log.Error("s.dao.SetIncomeCache error(%v)", err)
	}
	return
}

func (s *Service) archiveDetail(c context.Context, typ int, archiveID int64) (archives []*model.ArchiveIncome, err error) {
	archives = make([]*model.ArchiveIncome, 0)
	endTime := time.Now().AddDate(0, 0, -2).Format(_layout)
	switch typ {
	case _video:
		archives, err = s.dao.ListAvIncomeByID(c, archiveID, endTime)
		if err != nil {
			log.Error("s.dao.ListAvIncomeByID error(%v)", err)
			return
		}
	case _audio:
	case _column:
		archives, err = s.dao.ListColumnIncomeByID(c, archiveID, endTime)
		if err != nil {
			log.Error("s.dao.ListColumnIncomeByID error(%v)", err)
			return
		}
	case _bgm:
		archives, err = s.listBgmIncomeByID(c, archiveID, endTime)
		if err != nil {
			return
		}
	}

	sort.Slice(archives, func(i, j int) bool {
		return archives[i].Date > archives[j].Date
	})
	return
}

func (s *Service) listBgmIncomeByID(c context.Context, archiveID int64, endTime string) (archives []*model.ArchiveIncome, err error) {
	as, err := s.dao.ListBgmIncomeByID(c, archiveID, endTime)
	if err != nil {
		log.Error("s.dao.ListBgmIncomeByID error(%v)", err)
		return
	}
	am := make(map[xtime.Time][]*model.ArchiveIncome)
	for _, a := range as {
		if _, ok := am[a.Date]; ok {
			am[a.Date] = append(am[a.Date], a)
		} else {
			am[a.Date] = []*model.ArchiveIncome{a}
		}
	}
	archives = make([]*model.ArchiveIncome, 0)
	for date, ars := range am {
		a := &model.ArchiveIncome{}
		a.Date = date
		for _, ar := range ars {
			a.Income += ar.Income
			a.Avs = append(a.Avs, ar.ArchiveID)
		}
		archives = append(archives, a)
	}
	return
}

// ArchiveBreach get av_breach_record
func (s *Service) ArchiveBreach(c context.Context, mid int64, typ, page, size, all int) (data interface{}, err error) {
	if page == 0 {
		page = 1
	}
	start, end := (page-1)*size, page*size
	date := time.Now().AddDate(0, 0, -2)
	startMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local).Format(_layout)
	if all == 1 {
		startMonth = "2017-12-01"
	}
	archives, err := s.dao.ListAvBreach(c, mid, startMonth, date.Format(_layout))
	if err != nil {
		log.Error("s.dao.ListAvBreach error(%v)", err)
		return
	}
	typBreach := make([]*model.AvBreach, 0)
	for _, arch := range archives {
		if arch.CType == typ {
			typBreach = append(typBreach, arch)
		}
	}
	if len(typBreach) == 0 {
		return
	}

	breachs, err := s.breachInBlacklist(c, typBreach, typ)
	if err != nil {
		log.Error("s.breachInBlacklist error(%v)", err)
		return
	}
	sort.Slice(breachs, func(i, j int) bool {
		return breachs[i].CDate > breachs[j].CDate
	})
	if end > len(breachs) {
		end = len(breachs)
	}

	list := breachs[start:end]
	aIDs := make([]int64, 0)
	for _, b := range list {
		aIDs = append(aIDs, b.AvID)
	}
	var titles map[int64]string
	switch typ {
	case _video:
		titles, err = s.getAvTitle(c, aIDs)
		if err != nil {
			log.Error("s.getAvTitle error(%v)", err)
			return
		}
	case _column:
		titles, err = s.getColumnTitle(c, aIDs)
		if err != nil {
			log.Error("s.getColumnTitle error(%v)", err)
			return
		}
	case _bgm:
		titles, err = s.getBgmTitle(c, aIDs)
		if err != nil {
			log.Error("s.getBgmTitle error(%v)", err)
			return
		}
	}

	for _, b := range list {
		b.Title = titles[b.AvID]
	}
	data = map[string]interface{}{
		"data":        list,
		"total_count": len(breachs),
		"page":        page,
	}
	return
}

func (s *Service) breachInBlacklist(c context.Context, avs []*model.AvBreach, typ int) (breachs []*model.AvBreach, err error) {
	aIDList := make([]int64, 0)
	for _, a := range avs {
		aIDList = append(aIDList, a.AvID)
	}
	avBMap, err := s.dao.ListAvBlackList(c, aIDList, typ)
	if err != nil {
		return
	}
	breachs = make([]*model.AvBreach, 0)
	for _, av := range avs {
		if _, ok := avBMap[av.AvID]; ok {
			breachs = append(breachs, av)
		}
	}
	return
}

// UpIncomeStat get up income stat by month
func (s *Service) UpIncomeStat(c context.Context, typ int, mid int64, date time.Time) (data interface{}, err error) {
	redisKey := fmt.Sprintf("growup-income-stat:%d+%d", typ, mid)
	res, err := s.dao.GetIncomeCache(c, redisKey)
	if err != nil {
		log.Error("s.dao.GetIncomeCache error(%v)", err)
		return
	}
	if res != nil {
		data = res["data"]
		return
	}

	data, err = s.upIncomeStat(c, typ, mid, date)
	if err != nil {
		return
	}

	err = s.dao.SetIncomeCache(c, redisKey, map[string]interface{}{"data": data})
	if err != nil {
		log.Error("s.dao.SetIncomeCache error(%v)", err)
	}
	return
}

func (s *Service) upIncomeStat(c context.Context, typ int, mid int64, date time.Time) (stats []*model.UpIncomeStat, err error) {
	et := date.AddDate(0, 0, -2)
	end := et.Format(_layout)
	// last 30 days
	begin := et.AddDate(0, 0, -30).Format(_layout)

	upIncomes, err := s.dao.ListUpIncome(c, mid, "up_income", begin, end)
	if err != nil {
		log.Error("s.dao.ListUpIncome error(%v)", err)
		return
	}
	stats = make([]*model.UpIncomeStat, 0)
	if len(upIncomes) == 0 {
		return
	}
	var breachType []int64
	for _, up := range upIncomes {
		var income, baseIncome int64
		switch typ {
		case _video:
			income, baseIncome = up.AvIncome, up.AvBaseIncome
			breachType = []int64{0}
		case _column:
			income, baseIncome = up.ColumnIncome, up.ColumnBaseIncome
			breachType = []int64{2}
		case _bgm:
			income, baseIncome = up.BgmIncome, up.BgmBaseIncome
			breachType = []int64{3}
		case _up:
			income, baseIncome = up.Income, up.BaseIncome
			breachType = []int64{0, 1, 2, 3}
		}
		extra := income - baseIncome
		if extra < 0 {
			extra = 0
		}
		stats = append(stats, &model.UpIncomeStat{
			MID:         up.MID,
			Income:      income,
			BaseIncome:  baseIncome,
			ExtraIncome: extra,
			Date:        up.Date,
		})
	}

	rs, err := s.dao.GetAvBreachByType(c, mid, begin, end, breachType)
	if err != nil {
		log.Error("s.dao.GetAvBreachByType error(%v)", err)
		return
	}
	for _, stat := range stats {
		if _, ok := rs[stat.Date]; ok {
			stat.Breach = rs[stat.Date]
			delete(rs, stat.Date)
		}
	}
	for date, breach := range rs {
		stats = append(stats, &model.UpIncomeStat{
			MID:    mid,
			Date:   date,
			Breach: breach,
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Date < stats[j].Date
	})
	return
}
