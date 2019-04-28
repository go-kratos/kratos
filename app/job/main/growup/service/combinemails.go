package service

import (
	"context"
	"fmt"
	"time"

	"go-common/library/log"
)

const _combine = "%d年%d月%d日统计信息邮件"

// CombineMailsByHTTP send income mail by http.
func (s *Service) CombineMailsByHTTP(c context.Context, year int, month int, day int) (err error) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	err = s.execCombineMails(c, t)
	if err != nil {
		log.Error("s.CombineMailsByHTTP error(%v)", err)
	}
	return
}

// CombineMails combine mails.
func (s *Service) CombineMails() (err error) {
	c := context.TODO()
	date := time.Now().Add(-24 * time.Hour)
	t := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	err = s.execCombineMails(c, t)
	if err != nil {
		log.Error("s.CombineMails execCombineMails error(%v)", err)
	}
	return
}

func (s *Service) execCombineMails(c context.Context, date time.Time) (err error) {
	var body string
	singedUpBody, err := s.execSignedUps(c, date)
	if err != nil {
		log.Error("s.CombineMails s.execSignedUps error(%v)", err)
		return
	}
	body += singedUpBody
	incomeBody, err := s.execIncome(c, date.Add(-24*time.Hour))
	if err != nil {
		log.Error("s.CombineMails s.execIncome error(%v)", err)
		return
	}
	body += incomeBody

	uploadBody, err := s.execSendUpload(c, date.Add(-24*time.Hour))
	if err != nil {
		log.Error("s.CombineMails s.execSendUpload error(%v)", err)
		return
	}
	body += uploadBody

	topTenBody, err := s.execSendTopTen(c, date.Add(-24*time.Hour))
	if err != nil {
		log.Error("s.CombineMails s.execSendTopTen error(%v)", err)
		return
	}
	body += topTenBody

	var send []string
	for _, v := range s.conf.Mail.Send {
		if v.Type == 2 {
			send = v.Addr
		}
	}
	err = s.email.SendMail(date, fmt.Sprintf("<table border='1'>%s</table>", body), _combine, send...)
	if err != nil {
		log.Error("s.execSendUpload send upload.csv error(%v)", err)
		return
	}
	return
}
