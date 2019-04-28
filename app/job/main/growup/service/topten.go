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
	_topTen = "%d年%d月%d日UP主,稿件收入top 10 统计"
)

func (s *Service) execSendTopTen(c context.Context, date time.Time) (body string, err error) {
	var ups []*model.MIDInfo
	ups, err = s.getTopTenUps(c, date)
	if err != nil {
		log.Error("s.execSendTopTen s.getTopTenUps  error(%v)", err)
		return
	}
	avs, err := s.getTopTenAVs(c, date)
	if err != nil {
		log.Error("s.execSendTopTen s.getTopTenAVs  error(%v)", err)
		return
	}

	data := [][]string{
		{fmt.Sprintf(_topTen, date.Year(), date.Month(), date.Day())},
		{"日期", "名次", "UID", "昵称", "新增收入(元)", "累计收入(元)"},
	}
	for i, up := range ups {
		a := []string{
			fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day()), strconv.Itoa(i + 1), strconv.FormatInt(up.MID, 10), up.NickName, strconv.FormatFloat(float64(up.Income)/100, 'f', 2, 32), strconv.FormatFloat(float64(up.TotalIncome)/100, 'f', 2, 32),
		}
		data = append(data, a)
	}
	t := []string{"日期", "名次", "avid", "UID", "昵称", "新增收入(元)", "累计收入(元)"}
	data = append(data, t)
	for i, av := range avs {
		a := []string{
			fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day()), strconv.Itoa(i + 1), strconv.FormatInt(av.AVID, 10), strconv.FormatInt(av.MID, 10), av.NickName, strconv.FormatFloat(float64(av.Income)/100, 'f', 2, 32), strconv.FormatFloat(float64(av.TotalIncome)/100, 'f', 2, 32),
		}
		data = append(data, a)
	}

	f, err := os.Create("top10.csv")
	if err != nil {
		log.Error("s.execSendTopTen create top10.csv error(%v)", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(data)
	w.Flush()

	rfile, err := os.Open("top10.csv")
	if err != nil {
		log.Error("s.execSendTopTen open top10.csv error(%v)", err)
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
	err = os.Remove("top10.csv")
	if err != nil {
		log.Error("s.execSendTopTen remove top10.csv error(%v)", err)
	}
	return
}

func (s *Service) getTopTenUps(c context.Context, date time.Time) (infos []*model.MIDInfo, err error) {
	ups, err := s.getUpIncome(c, date)
	if err != nil {
		log.Error("s.getTopTenUps s.getUpIncome date(%v) error(%v)", date, err)
		return
	}
	infos = make([]*model.MIDInfo, 0)
	sort.SliceStable(ups, func(i, j int) bool {
		return ups[i].Income > ups[j].Income
	})
	cnt := 10
	if len(ups) < 10 {
		cnt = len(ups)
	}
	for i := 0; i < cnt; i++ {
		a := &model.MIDInfo{MID: ups[i].MID, Income: ups[i].Income, TotalIncome: ups[i].TotalIncome}
		a.NickName, err = s.dao.GetNickname(c, a.MID)
		if err != nil {
			log.Error("s.getTopTenUps dao.GetNickname mid(%v) error(%v)", a.MID, err)
			return
		}
		infos = append(infos, a)
	}
	return
}

func (s *Service) getTopTenAVs(c context.Context, date time.Time) (infos []*model.AVIDInfo, err error) {
	avs, err := s.getAvIncome(c, date)
	if err != nil {
		log.Error("s.getTopTenAVs s.getAvIncome error(%v)", err)
		return
	}
	sort.SliceStable(avs, func(i, j int) bool {
		return avs[i].Income > avs[j].Income
	})

	infos = make([]*model.AVIDInfo, 0)

	cnt := 10
	if len(avs) < 10 {
		cnt = len(avs)
	}
	for i := 0; i < cnt; i++ {
		a := &model.AVIDInfo{AVID: avs[i].AVID, MID: avs[i].MID, Income: avs[i].Income, TotalIncome: avs[i].TotalIncome}
		a.NickName, err = s.dao.GetNickname(c, a.MID)
		if err != nil {
			log.Error("s.getTopTenAVs dao.GetNickname mid(%v) error(%v)", a.MID, err)
			return
		}
		infos = append(infos, a)
	}
	return
}
