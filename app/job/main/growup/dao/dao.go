package dao

import (
	"context"

	"go-common/app/job/main/growup/conf"

	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
)

// Dao dao
type Dao struct {
	c            *conf.Config
	db           *sql.DB
	rddb         *sql.DB
	archiveURL   string
	typeURL      string
	breachURL    string
	forbidURL    string
	dismissURL   string
	passURL      string
	avLotteryURL string
	avVoteURL    string
	client       *xhttp.Client
	// hbase
	hbase *hbase.Client
}

// New fn
func New(c *conf.Config) (d *Dao) {
	log.Info("dao start")
	d = &Dao{
		c:            c,
		db:           sql.NewMySQL(c.Mysql.Growup),
		rddb:         sql.NewMySQL(c.Mysql.Allowance),
		client:       xhttp.NewClient(c.HTTPClient),
		archiveURL:   c.Host.Archive + "/manager/search",
		typeURL:      c.Host.VideoType + "/videoup/types",
		breachURL:    c.Host.Profit + "/allowance/api/x/admin/growup/auto/archive/breach",
		forbidURL:    c.Host.Profit + "/allowance/api/x/admin/growup/auto/up/forbid",
		dismissURL:   c.Host.Profit + "/allowance/api/x/admin/growup/auto/up/dismiss",
		passURL:      c.Host.Profit + "/allowance/api/x/admin/growup/up/pass",
		avLotteryURL: c.Host.VC + "/lottery_svr/v0/lottery_svr/export_rids",
		avVoteURL:    c.Host.API + "/x/internal/creative/archive/vote",
		// hbase
		hbase: hbase.NewClient(c.HBase.Config),
	}
	log.Info("data init end")
	//d.db.State = prom.LibClient
	return
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// BeginTran begin transcation
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}
