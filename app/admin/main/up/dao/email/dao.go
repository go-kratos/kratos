package email

import (
	"crypto/tls"

	"go-common/app/admin/main/up/conf"

	gomail "gopkg.in/gomail.v2"
)

// Dao is redis dao.
type Dao struct {
	email *gomail.Dialer
}

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// mail
		email: gomail.NewDialer(c.MailConf.Host, c.MailConf.Port, c.MailConf.Username, c.MailConf.Password),
	}
	d.email.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return
}
