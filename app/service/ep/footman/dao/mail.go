package dao

import (
	"context"

	"gopkg.in/gomail.v2"
)

// SendMail asynchronous send mail.
func (d *Dao) SendMail(message *gomail.Message) {
	message.SetAddressHeader("From", d.email.Username, "footman")
	d.cache.Do(context.TODO(), func(ctx context.Context) {
		d.email.DialAndSend(message)
	})
}
