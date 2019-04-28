package dao

import (
	"go-common/library/log"
	"gopkg.in/gomail.v2"
)

func (d *Dao) SendEmail(subject string, to string, body string) (err error) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", d.c.Mail.Username)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)
	if err = d.email.DialAndSend(msg); err != nil {
		log.Error("s.email.DialAndSend error(%v)", err)
		return
	}
	return
}
