package dao

import (
	"context"
	"crypto/tls"

	"go-common/app/interface/main/laser/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	"gopkg.in/gomail.v2"
)

type Dao struct {
	c        *conf.Config
	db       *xsql.DB
	mc       *memcache.Pool
	mcExpire int32
	email    *gomail.Dialer
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:        c,
		db:       xsql.NewMySQL(c.Mysql),
		mc:       memcache.NewPool(c.Memcache.Laser.Config),
		mcExpire: 3600 * 6,
		email:    gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.Username, c.Mail.Password),
	}
	d.email.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return
}

func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// BeginTran BeginTran.
func (d *Dao) BeginTran(c context.Context) (*xsql.Tx, error) {
	return d.db.Begin(c)
}
