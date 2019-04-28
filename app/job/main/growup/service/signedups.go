package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
)

const (
	_SignedUps = "%d年%d月%d日加入人数统计"
)

func (s *Service) execSignedUps(c context.Context, date time.Time) (body string, err error) {
	tm1, err := s.dao.GetDateSignedUps(c, date, date.Add(24*time.Hour))
	if err != nil {
		log.Error("s.execSignedUps dao.GetDateSignedUps T-1(%v) signed ups error(%v)", date, err)
		return
	}
	tm2, err := s.dao.GetDateSignedUps(c, date.Add(-24*time.Hour), date)
	if err != nil {
		log.Error("s.execSignedUps dao.GetDateSignedUps T-2(%v) signed ups error(%v)", date.Add(-24*time.Hour), err)
		return
	}
	tm8, err := s.dao.GetDateSignedUps(c, date.Add(-24*7*time.Hour), date.Add(-24*6*time.Hour))
	if err != nil {
		log.Error("s.execSignedUps dao.GetDateSignedUps T-8(%v) signed ups error(%v)", date.Add(-24*6*time.Hour), err)
		return
	}
	log.Info("t-8:%d, %d-%d-%d ", tm8, date.Add(-24*7*time.Hour).Year(), date.Add(-24*7*time.Hour).Month(), date.Add(-24*7*time.Hour).Day())

	sa, err := s.dao.GetAllSignedUps(c, date)
	if err != nil {
		log.Error("s.execSignedUps dao.GetAllSignedUps date(%v) error(%v)", date, err)
		return
	}

	applyDm1, err := s.dao.GetVideoApplyUpCount(c, date, date.Add(24*time.Hour))
	if err != nil {
		log.Error("s.execSignedUps dao.GetVideoUpCount applyDm1 error(%v)", err)
		return
	}
	applyDm2, err := s.dao.GetVideoApplyUpCount(c, date.Add(-24*time.Hour), date)
	if err != nil {
		log.Error("s.execSignedUps dao.GetVideoUpCount applyDm2 error(%v)", err)
		return
	}
	var du, wu float64
	if tm2 > 0 {
		du = float64(tm1-tm2) / float64(tm2) * 100
	}
	if tm8 > 0 {
		wu = float64(tm1-tm8) / float64(tm8) * 100
	}
	var applyDu float64
	if applyDm2 > 0 {
		applyDu = float64(applyDm1-applyDm2) / float64(applyDm2) * 100
	}
	data := [][]string{
		{fmt.Sprintf(_SignedUps, date.Year(), date.Month(), date.Day())},
		{"日期", "已签约新增加人数", "已签约人数日增长率", "已签约人数同比上周增长率", "累计签约人数", "已申请新增加人数", "已申请人数日增长率"},
		{fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day()), strconv.Itoa(tm1), strconv.FormatFloat(du, 'f', 2, 32), strconv.FormatFloat(wu, 'f', 2, 32), strconv.Itoa(sa), strconv.Itoa(applyDm1), strconv.FormatFloat(applyDu, 'f', 2, 32)},
	}

	f, err := os.Create("signedups.csv")
	if err != nil {
		log.Error("s.execSignedUps create signedups.csv error(%v)", err)
		return
	}
	log.Info("tm1:%d, tm2:%d, tm8:%d, sa:%d", tm1, tm2, tm8, sa)
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(data)
	w.Flush()

	rfile, err := os.Open("signedups.csv")
	if err != nil {
		log.Error("s.execSignedUps open signedups.csv error(%v)", err)
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
	err = os.Remove("signedups.csv")
	if err != nil {
		log.Error("growup-job s.execSignedUps remove signedups.csv error(%v)", err)
	}
	return
}
