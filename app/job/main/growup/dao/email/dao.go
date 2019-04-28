package email

import (
	"crypto/tls"

	"go-common/app/job/main/growup/conf"

	"gopkg.in/gomail.v2"
)

// Dao is redis dao.
type Dao struct {
	c     *conf.Config
	email *gomail.Dialer
}

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// mail
		email: gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.Username, c.Mail.Password),
	}
	d.email.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return
}
