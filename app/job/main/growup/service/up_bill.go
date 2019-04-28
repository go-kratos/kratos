package service

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/app/job/main/growup/model/income"

	"go-common/library/log"
)

var (
	_dbLimit     = 2000
	_dbBatchSize = 2000
)

// CreativeUpBill creative up bill
func (s *Service) CreativeUpBill(c context.Context, startDate, endDate time.Time) (err error) {
	// up_info_video
	ups, err := s.signed(c, int64(_dbLimit))
	if err != nil {
		log.Error("s.signed error(%v)", err)
		return
	}
	upBills := handleUps(ups, endDate)

	// up_income
	upIncome, err := s.getUpIncomeByDate(c, "up_income", startDate, endDate)
	if err != nil {
		log.Error("s.getUpIncomeByDate error(%v)", err)
		return
	}
	handleUpIncome(upBills, upIncome)

	// up_signed_avs
	upAvs, err := s.upSignedAvs(c, _dbLimit)
	if err != nil {
		log.Error("s.upSignedAvs error(%v)", err)
		return
	}
	for _, up := range upBills {
		up.AvCount = upAvs[up.MID]
	}

	// av_income_statis
	avs, err := s.getAvIncomeStatis(c, int64(_dbLimit), endDate)
	if err != nil {
		log.Error("s.getAvIncomeStatis error(%v)", err)
		return
	}
	handleAvIncomeStatis(upBills, avs)

	upQuality, err := s.getUpQuality(c, int(endDate.Day()), _dbLimit)
	if err != nil {
		log.Error("s.getUpQualities error(%v)", err)
		return
	}
	if len(upQuality) == 0 {
		err = fmt.Errorf("Error: get 0 ups from up_quality_info_%d", endDate.Day())
		return
	}
	handleUpQuality(upBills, upQuality)
	handleUpBills(upBills)

	// insert
	err = s.upBillDBStore(c, upBills)
	if err != nil {
		log.Error("s.upBillDBStore error(%v)", err)
	}
	return
}

func randStr(strs []string) string {
	return strs[rand.Intn(len(strs))%len(strs)]
}

func handleUpBills(upBills map[int64]*model.UpBill) {
	for _, up := range upBills {
		titles := []string{}
		shareItem := ""
		switch {
		case up.TotalIncome >= 500000:
			titles = append(titles, "掘金小能手")
			shareItem = randStr([]string{"98亿手办", "圣地巡礼机票"})
		case up.TotalIncome >= 100000:
			shareItem = randStr([]string{"BML现场门票", "老婆的演唱会门票", "购物车里的“老婆”"})
		case up.TotalIncome >= 50000:
			shareItem = randStr([]string{"超大堆小电视抱枕", "肥宅快乐桶吃到吐", "老婆的应援周边"})
		case up.TotalIncome >= 10000:
			shareItem = randStr([]string{"一堆“2233”挂件", "N个月大会员", "一暑假肥宅快乐水"})
		case up.TotalIncome < 10000:
			shareItem = randStr([]string{"创作补给餐", "自我打call棒", "承包几部番剧", "老婆的海报"})
		}
		if up.AvCount >= 30 {
			titles = append(titles, "B站劳模")
		}
		if up.Fans >= 10000 {
			titles = append(titles, "万人迷")
		}
		if up.TotalPlayCount >= 500000 {
			titles = append(titles, "流量王")
		}
		if len(titles) == 0 {
			titles = []string{"社会人", "快乐肥宅", "9percent"}
		}
		up.Title = randStr(titles)
		up.ShareItems = shareItem
	}
}

func handleUps(ups map[int64]*model.UpInfoVideo, end time.Time) (upBills map[int64]*model.UpBill) {
	upBills = make(map[int64]*model.UpBill)
	for _, up := range ups {
		upBills[up.MID] = &model.UpBill{
			MID:            up.MID,
			SignedAt:       up.SignedAt.Time().Format(_layout),
			Fans:           up.Fans,
			TotalPlayCount: up.TotalPlayCount,
			EndAt:          end.Format(_layout),
		}
	}
	return
}

func handleUpIncome(upBills map[int64]*model.UpBill, upIncome []*model.UpIncome) {
	for _, up := range upIncome {
		upB, ok := upBills[up.MID]
		if !ok {
			continue
		}
		if up.Date.Time().Format(_layout) >= upB.SignedAt {
			upB.TotalIncome += up.Income
			if upB.FirstTime == "" {
				upB.FirstIncome = up.Income
				upB.FirstTime = up.Date.Time().Format(_layout)
			}
			if upB.MaxIncome < up.Income {
				upB.MaxIncome = up.Income
				upB.MaxTime = up.Date.Time().Format(_layout)
			}
		}
	}
}

func handleAvIncomeStatis(upBills map[int64]*model.UpBill, avs map[int64]*income.AvIncomeStat) {
	for _, av := range avs {
		upB, ok := upBills[av.MID]
		if !ok {
			continue
		}
		if av.CTime.Time().Format(_layout) < upB.SignedAt {
			continue
		}
		income := av.TotalIncome
		if income > upB.AvMaxIncome {
			upB.AvMaxIncome = income
			upB.AvID = av.AvID
		}
	}
}

func handleUpQuality(upBills map[int64]*model.UpBill, upQualities []*model.UpQuality) {
	total := len(upQualities)
	sort.Slice(upQualities, func(i, j int) bool {
		return upQualities[i].Quality > upQualities[j].Quality
	})
	for i := 0; i < len(upQualities); i++ {
		mid := upQualities[i].MID
		if _, ok := upBills[mid]; ok {
			upBills[mid].QualityValue = upQualities[i].Quality
			rank := i
			for rank > 0 && upQualities[rank].Quality == upQualities[rank-1].Quality {
				rank--
			}
			upBills[mid].DefeatNum = (10000 * (total - rank)) / total
		}
	}
}

func (s *Service) signed(c context.Context, limit int64) (m map[int64]*model.UpInfoVideo, err error) {
	var id int64
	m = make(map[int64]*model.UpInfoVideo)
	for {
		var us map[int64]*model.UpInfoVideo
		id, us, err = s.dao.UpInfoVideo(c, id, limit)
		if err != nil {
			return
		}
		for k, v := range us {
			if v.AccountState == 3 && v.IsDeleted == 0 {
				m[k] = v
			}
		}
		if len(us) < _dbLimit {
			break
		}
	}
	return
}

func (s *Service) getUpIncomeByDate(c context.Context, table string, start, end time.Time) (ups []*model.UpIncome, err error) {
	ups = make([]*model.UpIncome, 0)
	end = end.AddDate(0, 0, 1)
	for start.Before(end) {
		var up []*model.UpIncome
		up, err = s.GetUpIncome(c, table, start.Format("2006-01-02"))
		if err != nil {
			return
		}
		ups = append(ups, up...)
		start = start.AddDate(0, 0, 1)
	}
	return
}

func (s *Service) getAvIncomeStatis(c context.Context, limit int64, endDate time.Time) (m map[int64]*income.AvIncomeStat, err error) {
	m = make(map[int64]*income.AvIncomeStat)
	var id int64
	for {
		var am map[int64]*income.AvIncomeStat
		am, id, err = s.income.AvIncomeStat(c, id, limit)
		if err != nil {
			return
		}
		for avID, stat := range am {
			if stat.CTime.Time().Before(endDate.AddDate(0, 0, 1)) {
				m[avID] = stat
			}
		}
		if len(am) < int(limit) {
			break
		}
	}
	return
}

func (s *Service) upSignedAvs(c context.Context, limit int) (upAvs map[int64]int64, err error) {
	upAvs = make(map[int64]int64)
	var id int64
	for {
		var ups map[int64]int64
		ups, id, err = s.dao.ListUpSignedAvs(c, id, limit)
		if err != nil {
			return
		}
		for mid, avCount := range ups {
			upAvs[mid] = avCount
		}
		if len(ups) < limit {
			break
		}
	}
	return
}

func (s *Service) getUpQuality(c context.Context, day, limit int) (ups []*model.UpQuality, err error) {
	ups = make([]*model.UpQuality, 0)
	var (
		id    int64
		table string
		up    []*model.UpQuality
	)
	if day < 10 {
		table = fmt.Sprintf("up_quality_info_0%d", day)
	} else {
		table = fmt.Sprintf("up_quality_info_%d", day)
	}
	for {
		up, id, err = s.dao.GetUpQuality(c, table, id, limit)
		if err != nil {
			return
		}
		ups = append(ups, up...)
		if len(up) < limit {
			break
		}
	}
	return
}

func (s *Service) upBillDBStore(c context.Context, upBill map[int64]*model.UpBill) (err error) {
	var (
		buff    = make([]*model.UpBill, _dbBatchSize)
		buffEnd = 0
	)

	for _, u := range upBill {
		buff[buffEnd] = u
		buffEnd++

		if buffEnd >= _dbBatchSize {
			_, err = s.upBillBatchInsert(c, buff[:buffEnd])
			if err != nil {
				return
			}
			buffEnd = 0
		}
	}
	if buffEnd > 0 {
		_, err = s.upBillBatchInsert(c, buff[:buffEnd])
		if err != nil {
			return
		}
		buffEnd = 0
	}
	return
}

func assembleUpBill(upBill []*model.UpBill) (vals string) {
	var buf bytes.Buffer
	for _, row := range upBill {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.FirstIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MaxIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.AvCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.AvMaxIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.AvID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.QualityValue, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(row.DefeatNum))
		buf.WriteByte(',')
		buf.WriteString("\"" + strings.Replace(row.Title, "\"", "\\\"", -1) + "\"")
		buf.WriteByte(',')
		buf.WriteString("'" + row.ShareItems + "'")
		buf.WriteByte(',')
		buf.WriteString("'" + row.FirstTime + "'")
		buf.WriteByte(',')
		buf.WriteString("'" + row.MaxTime + "'")
		buf.WriteByte(',')
		buf.WriteString("'" + row.SignedAt + "'")
		buf.WriteByte(',')
		buf.WriteString("'" + row.EndAt + "'")
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals = buf.String()
	buf.Reset()
	return
}

func (s *Service) upBillBatchInsert(c context.Context, upBill []*model.UpBill) (rows int64, err error) {
	vals := assembleUpBill(upBill)
	rows, err = s.dao.InsertUpBillBatch(c, vals)
	return
}
