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

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	_upload = "%d年%d月%d日投稿转换统计"
)

func (s *Service) execSendUpload(c context.Context, date time.Time) (body string, err error) {
	// total signed up count
	ts, err := s.dao.GetDateSignedUps(c, time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), date.Add(24*time.Hour))
	if err != nil {
		log.Error("s.execSendUpload GetDateSignedUps error(%v)", date, err)
		return
	}
	// total_income > 0 up count
	iu, err := s.dao.GetUpTotalIncomeCnt(c)
	if err != nil {
		log.Error("s.execSendUpload dao.GetUpTotalIncomeCnt date(%v) error(%v)", date, err)
		return
	}
	// total_income > 0 archive count
	ia, err := s.avSatisCount(c)
	if err != nil {
		log.Error("s.execSendUpload s.avSatisCount date(%v) error(%v)", date, err)
		return
	}
	// upcnt, avcnt
	upCnt, avCnt, err := s.recvSignedUpInfo(c, date)
	if err != nil {
		log.Error("s.execSendUpload s.recvSignedUpInfo  date(%v) error(%v)", date, err)
		return
	}
	data := [][]string{
		{fmt.Sprintf(_upload, date.Year(), date.Month(), date.Day())},
		{"日期", "累计签约人数", "投稿UP数", "获得收入UP数", "UP收入率", "投稿总数", "获得收入稿件数", "稿件收入率"},
		{fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day()), strconv.Itoa(ts), strconv.Itoa(upCnt), strconv.Itoa(iu), strconv.FormatFloat(float64(iu)/float64(upCnt)*100, 'f', 2, 32), strconv.Itoa(avCnt), strconv.Itoa(ia), strconv.FormatFloat(float64(ia)/float64(avCnt)*100, 'f', 2, 32)},
	}

	f, err := os.Create("upload.csv")
	if err != nil {
		log.Error("s.execSendUpload create upload.csv error(%v)", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(data)
	w.Flush()

	rfile, err := os.Open("upload.csv")
	if err != nil {
		log.Error("s.execSendUpload open upload.csv  error(%v)", err)
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
	err = os.Remove("upload.csv")
	if err != nil {
		log.Error("s.execSendUpload remove upload.csv error(%v)", err)
	}
	return
}

func (s *Service) avSatisCount(c context.Context) (avCnt int, err error) {
	return s.dao.GetAvStatisCount(c)
}

// recvSignedUpInfo recv signed up upload archive cnt; archive cnt > 0, up cnt.
func (s *Service) recvSignedUpInfo(c context.Context, date time.Time) (upCnt, avCnt int, err error) {
	query := "{\"select\": [{\"name\": \"up_cnt\",\"as\": \"up_cnt\"},{\"name\": \"archive_cnt\",\"as\": \"archive_cnt\"},{\"name\": \"log_date\",\"as\": \"log_date\"}],\"where\": {\"log_date\": {\"in\": [\"%s\"]}}}"
	t := strconv.Itoa(date.Year())
	if int(date.Month()) < 10 {
		t += "0"
	}
	t += strconv.Itoa(int(date.Month()))
	if date.Day() < 10 {
		t += "0"
	}
	t += strconv.Itoa(date.Day())

	var res []*model.ArchiveInfo
	res, err = s.dp.Send(c, fmt.Sprintf(query, t))
	if err != nil {
		log.Error("s.recvSignedUpInfo error(%v)", err)
		return
	}
	for _, i := range res {
		if i.UploadDate == t {
			upCnt = i.UpCnt
			avCnt = i.ArchiveCnt
		}
	}
	return
}
