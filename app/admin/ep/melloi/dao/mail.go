package dao

import (
	"go-common/library/log"

	"gopkg.in/gomail.v2"
)

// SendMail asynchronous send mail.
func (d *Dao) SendMail(message *gomail.Message) (err error) {
	message.SetAddressHeader("From", d.email.Username, "melloi")
	err = d.email.DialAndSend(message)
	if err != nil {
		log.Error("send email error :(%v)", err)
	}
	return
}
