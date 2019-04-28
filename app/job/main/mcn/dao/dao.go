package dao

import (
	"context"
	"crypto/tls"

	"go-common/app/job/main/mcn/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"

	gomail "gopkg.in/gomail.v2"
)

// Dao dao
type Dao struct {
	c     *conf.Config
	mc    *memcache.Pool
	db    *xsql.DB
	email *gomail.Dialer
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:  c,
		mc: memcache.NewPool(c.Memcache),
		db: xsql.NewMySQL(c.MySQL),
		// mail
		email: gomail.NewDialer(c.MailConf.Host, c.MailConf.Port, c.MailConf.Username, c.MailConf.Password),
	}
	dao.email.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}
