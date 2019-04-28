package dao

import (
	"context"

	"go-common/app/job/main/passport/conf"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Dao struct info of Dao.
type Dao struct {
	c             *conf.Config
	logDB         *sql.DB
	asoDB         *sql.DB
	client        *bm.Client
	gameClient    *bm.Client
	loginLogHBase *hbase.Client
	pwdLogHBase   *hbase.Client

	setTokenURI     string
	delCacheURI     string
	delGameCacheURI string
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:             c,
		logDB:         sql.NewMySQL(c.DB.Log),
		asoDB:         sql.NewMySQL(c.DB.ASO),
		client:        bm.NewClient(c.HTTPClient),
		gameClient:    bm.NewClient(c.Game.Client),
		loginLogHBase: hbase.NewClient(c.HBase.LoginLog.Config),
		pwdLogHBase:   hbase.NewClient(c.HBase.PwdLog.Config),

		setTokenURI:     c.URI.SetToken,
		delCacheURI:     c.URI.DelCache,
		delGameCacheURI: c.Game.DelCacheURI,
	}
	return d
}

// Ping ping check dao health.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.logDB.Ping(c); err != nil {
		log.Info("dao.logDB.Ping() error(%v)", err)
	}
	return
}

// Close close dao.
func (d *Dao) Close() (err error) {
	if d.logDB != nil {
		d.logDB.Close()
	}
	if d.loginLogHBase != nil {
		d.loginLogHBase.Close()
	}
	if d.pwdLogHBase != nil {
		d.pwdLogHBase.Close()
	}
	return
}
