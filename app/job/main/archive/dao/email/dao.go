package email

import (
	gomail "gopkg.in/gomail.v2"

	"go-common/app/job/main/archive/conf"
)

// Dao is redis dao.
type Dao struct {
	c     *conf.Config
	email *gomail.Dialer
}

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:     c,
		email: gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.Username, c.Mail.Password),
	}
	return d
}
