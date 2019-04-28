package dao

import (
	"context"

	"go-common/app/job/main/passport-user-compare/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	safeQuestionSegment = 100
)

// Dao dao
type Dao struct {
	c          *conf.Config
	httpClient *bm.Client
	originDB   *sql.DB
	userDB     *sql.DB
	secretDB   *sql.DB
}

// New new dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		originDB:   sql.NewMySQL(c.DB.Origin),
		userDB:     sql.NewMySQL(c.DB.User),
		secretDB:   sql.NewMySQL(c.DB.Secret),
		httpClient: bm.NewClient(c.HTTPClient),
	}
	return
}

// Ping check dao ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.originDB.Ping(c); err != nil {
		log.Info("dao.originDB.Ping() error(%v)", err)
	}
	if err = d.userDB.Ping(c); err != nil {
		log.Info("dao.newDB.Ping() error(%v)", err)
	}
	if err = d.secretDB.Ping(c); err != nil {
		log.Info("dao.secretDB.Ping() error(%v)", err)
	}
	return
}

// Close close connections.
func (d *Dao) Close() (err error) {
	if d.originDB != nil {
		d.originDB.Close()
	}
	if d.userDB != nil {
		d.userDB.Close()
	}
	if d.secretDB != nil {
		d.secretDB.Close()
	}
	return
}
