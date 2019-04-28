package email

import (
	"fmt"
	"os"
	"time"

	"go-common/library/log"

	gomail "gopkg.in/gomail.v2"
)

//SendMail send the email.
func (d *Dao) SendMail(date time.Time, body string, subject string, send ...string) (err error) {
	log.Info("send mail send:%v", send)
	msg := gomail.NewMessage()
	msg.SetHeader("From", d.c.Mail.Username)
	msg.SetHeader("To", send...)
	msg.SetHeader("Subject", fmt.Sprintf(subject, date.Year(), date.Month(), date.Day()))
	msg.SetBody("text/html", body)
	if err = d.email.DialAndSend(msg); err != nil {
		log.Error("s.email.DialAndSend error(%v)", err)
		return
	}
	return
}

//SendMailAttach send the email.
func (d *Dao) SendMailAttach(filename string, subject string, send []string) (err error) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", d.c.Mail.Username)
	msg.SetHeader("To", send...)
	msg.SetHeader("Subject", subject)
	msg.Attach(filename)
	if err = d.email.DialAndSend(msg); err != nil {
		log.Error("s.email.DialAndSend error(%v)", err)
		return
	}
	err = os.Remove(filename)
	return
}
