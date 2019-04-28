package dao

import (
	"context"
	"net/http"

	"go-common/app/interface/main/player/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Dao dao.
type Dao struct {
	// config
	c *conf.Config
	// mysql
	showDB *sql.DB
	// stmt
	paramStmt *sql.Stmt
	// client
	client   *xhttp.Client
	vsClient *http.Client
	// API URL
	blockTimeURL   string
	onlineCountURL string
	viewPointsURL  string
}

// New return new dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:        c,
		showDB:   sql.NewMySQL(c.MySQL.Show),
		client:   xhttp.NewClient(c.HTTPClient),
		vsClient: http.DefaultClient,
	}
	d.paramStmt = d.showDB.Prepared(_param)
	d.blockTimeURL = c.Host.AccCo + _blockTimeURI
	d.onlineCountURL = c.Host.APICo + _onlineCountURI
	d.viewPointsURL = c.Host.APICo + _viewPointsURI
	return
}

// Ping check service health
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.showDB.Ping(c); err != nil {
		log.Error("d.show.Ping() err(%v)", err)
	}
	return
}

// PromError stat and log.
func PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}
