package service

import "gopkg.in/gomail.v2"

func (s *Service) sendMail(receiver []string, header, body string) {
	m := gomail.NewMessage()
	m.SetHeader("To", receiver...)
	m.SetHeader("Subject", header)
	m.SetBody("text/html", body)
	s.dao.SendMail(m)
}

// TapdMailNotice Tapd Mail Notice.
func (s *Service) TapdMailNotice(header, body string) {
	s.sendMail(s.c.Mail.NoticeOwner, header, body)
}
