package dao

import (
	"context"
	"strings"

	"gopkg.in/gomail.v2"
)

const (
	_MailBoxNotFound = "Mailbox not found"
)

// SendMail asynchronous send mail.
func (d *Dao) SendMail(message *gomail.Message) {
	message.SetAddressHeader("From", d.email.Username, "merlin")
	d.cache.Do(context.TODO(), func(ctx context.Context) {
		d.SendMailIfFailed(message)
	})
}

// SendMailIfFailed Send Mail If Failed
func (d *Dao) SendMailIfFailed(message *gomail.Message) {
	if err := d.email.DialAndSend(message); err != nil {
		if strings.Contains(err.Error(), _MailBoxNotFound) {
			headerMsg := message.GetHeader("Subject")
			headerMsg = append(headerMsg, "Mail Send Error:"+err.Error()+",Receiver:")
			headerMsg = append(headerMsg, message.GetHeader("To")...)

			message.SetHeader("To", d.c.Mail.NoticeOwner...)
			message.SetHeader("Subject", headerMsg...)
			d.email.DialAndSend(message)
		}
	}
}
