package dao

import (
	"context"

	"go-common/app/service/main/passport/conf"
	"go-common/library/database/hbase.v2"
)

// Dao struct answer history of Dao
type Dao struct {
	c             *conf.Config
	hbase         *hbase.Client
	loginLogHBase *hbase.Client
	pwdLogHBase   *hbase.Client
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	var loginLogHBase *hbase.Client
	if c.Switch.LoginLogHBase {
		loginLogHBase = hbase.NewClient(c.HBase.LoginLog.Config)
	}

	d = &Dao{
		c:             c,
		hbase:         hbase.NewClient(c.HBase.FaceApply.Config),
		loginLogHBase: loginLogHBase,
		pwdLogHBase:   hbase.NewClient(c.HBase.PwdLog.Config),
	}
	return
}

// Close close connections.
func (d *Dao) Close() {
	if d.hbase != nil {
		d.hbase.Close()
	}
	if d.loginLogHBase != nil {
		d.loginLogHBase.Close()
	}
	if d.pwdLogHBase != nil {
		d.pwdLogHBase.Close()
	}
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
