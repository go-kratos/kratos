package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	_income = "%d年%d月%d日收入统计"
)

func (s *Service) execIncome(c context.Context, date time.Time) (body string, err error) {
	upInfo, err := s.getUpIncome(c, date)
	if err != nil {
		log.Error("s.execIncome s.getUpIncome error(%v)", err)
		return
	}
	avInfo, err := s.getAvIncome(c, date)
	if err != nil {
		log.Error("s.execIncome s.getAvIncome error(%v)", err)
		return
	}
	totalIncome, err := s.getTotalIncome(c)
	if err != nil {
		log.Error("s.execIncome s.getTotalIncome error(%v)", err)
		return
	}

	var income, avgUpIncome, avgAVIncome, upCnt int64
	for _, up := range upInfo {
		income += up.Income
		upCnt++
	}
	if upCnt > 0 {
		avgUpIncome = income / upCnt // avg up income
	}
	if len(avInfo) > 0 {
		avgAVIncome = income / int64(len(avInfo)) // avg av income
	}

	midAVIncome := getMIDAVIncome(avInfo)
	midUpIncome := getMIDUpIncome(upInfo)
	data := [][]string{
		{fmt.Sprintf(_income, date.Year(), date.Month(), date.Day())},
		{"日期", "新增收入(元)", "累计收入(元)", "获得收入UP数", "UP平均收入(元)", "UP收入中位数(元)", "获得收入稿件数", "稿件平均收入(元)", "稿件收入中位数(元)"},
		{fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day()), strconv.FormatFloat(float64(income)/100, 'f', 2, 32), strconv.FormatFloat(float64(totalIncome)/100, 'f', 2, 32), strconv.FormatInt(upCnt, 10), strconv.FormatFloat(float64(avgUpIncome)/100, 'f', 2, 32), strconv.FormatFloat(float64(midUpIncome)/100, 'f', 2, 32), strconv.Itoa(len(avInfo)), strconv.FormatFloat(float64(avgAVIncome)/100, 'f', 2, 32), strconv.FormatFloat(float64(midAVIncome)/100, 'f', 2, 32)},
	}
	f, err := os.Create("upincome.csv")
	if err != nil {
		log.Error("growup-job create upincome.csv error(%v)", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(data)
	w.Flush()

	rfile, err := os.Open("upincome.csv")
	if err != nil {
		log.Error("growup-job open upincome.csv error(%v)", err)
		return
	}
	defer rfile.Close()
	r := csv.NewReader(rfile)
	for {
		var strs []string
		strs, err = r.Read()
		if err == io.EOF {
			break
		}
		body += fmt.Sprintf("<tr><td>%s</td></tr>", strings.Join(strs, "</td><td>"))
	}
	err = os.Remove("upincome.csv")
	if err != nil {
		log.Error("growup-job s.execIncome remove upincome.csv error(%v)", err)
	}
	return
}

func (s *Service) getTotalIncome(c context.Context) (totalIncome int64, err error) {
	var (
		totalInfos, infos []*model.MIDInfo
		from, limit       int64
	)
	limit = 2000

	for {
		infos, err = s.dao.GetUpTotalIncome(c, from, limit)
		if err != nil {
			log.Error("s.getTotalIncome email.GetUpTotalIncome error(%v)", err)
			return
		}
		totalInfos = append(totalInfos, infos...)
		if int64(len(infos)) < limit {
			break
		}
		from = infos[len(infos)-1].ID
	}

	for _, v := range totalInfos {
		if v.IsDeleted == 0 {
			totalIncome += v.TotalIncome
		}
	}
	return
}

func (s *Service) getUpIncome(c context.Context, date time.Time) (totalInfos []*model.MIDInfo, err error) {
	var (
		infos       []*model.MIDInfo
		from, limit int64
	)
	limit = 2000
	for {
		infos, err = s.dao.GetUpIncome(c, date, from, limit)
		if err != nil {
			log.Error("s.getUpIncome dao.GetUpIncome error(%v)", err)
			return
		}
		totalInfos = append(totalInfos, infos...)
		if int64(len(infos)) < limit {
			break
		}
		from = infos[len(infos)-1].ID
	}
	return
}

func (s *Service) getAvIncome(c context.Context, date time.Time) (totalInfos []*model.IncomeInfo, err error) {
	var (
		infos       []*model.IncomeInfo
		from, limit int64
	)
	limit = 3000
	for {
		infos, err = s.dao.GetAvIncome(c, date, from, limit)
		if err != nil {
			log.Error("s.getAvIncome dao.GetAvIncome error(%v)", err)
			return
		}
		for _, i := range infos {
			if i.Income > 0 {
				totalInfos = append(totalInfos, i)
			}
		}
		if int64(len(infos)) < limit {
			break
		}
		from = infos[len(infos)-1].ID
	}
	return
}

func getMIDAVIncome(avs []*model.IncomeInfo) (income int64) {
	var cnt = len(avs)
	if cnt <= 0 {
		return
	}
	if cnt == 1 {
		income = avs[0].Income
		return
	}

	sort.SliceStable(avs, func(i, j int) bool {
		return avs[i].Income < avs[j].Income
	})
	if cnt%2 == 0 {
		income = (avs[cnt/2].Income + avs[cnt/2-1].Income) / 2
	} else {
		income = avs[cnt/2].Income
	}
	return
}

func getMIDUpIncome(ups []*model.MIDInfo) (income int64) {
	var cnt = len(ups)
	if cnt == 0 {
		return
	}
	if cnt == 1 {
		income = ups[0].Income
		return
	}
	sort.SliceStable(ups, func(i, j int) bool {
		return ups[i].Income < ups[j].Income
	})
	if cnt%2 == 0 {
		income = (ups[cnt/2].Income + ups[cnt/2-1].Income) / 2
	} else {
		income = ups[cnt/2].Income
	}
	return
}
