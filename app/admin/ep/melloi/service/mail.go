package service

import (
	"gopkg.in/gomail.v2"
)

// SendMail send email
func (s *Service) SendMail(receiver string, subject string, content string) error {
	var (
		m = gomail.NewMessage()
	)

	m.SetHeader("To", receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	return s.dao.SendMail(m)
}
