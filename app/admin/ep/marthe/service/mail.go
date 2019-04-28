package service

import "gopkg.in/gomail.v2"

// SendMail Send Mail.
func (s *Service) SendMail(receiver []string, header, body string) {
	m := gomail.NewMessage()
	m.SetHeader("To", receiver...)
	m.SetHeader("Subject", header)
	m.SetBody("text/html", body)
	s.dao.SendMail(m)
}
