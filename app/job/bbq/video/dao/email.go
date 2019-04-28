package dao

import (
	gomail "gopkg.in/gomail.v2"
)

// SendMail asynchronous send mail.
func (d *Dao) SendMail(message *gomail.Message) (err error) {
	message.SetAddressHeader("From", d.c.Mail.From, "bbq")
	err = d.email.DialAndSend(message)
	return
}
