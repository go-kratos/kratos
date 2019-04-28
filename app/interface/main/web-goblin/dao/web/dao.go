package web

import (
	"context"

	"go-common/app/interface/main/web-goblin/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

const (
	_pgcFullURL  = "/ext/internal/archive/channel/content"
	_pgcIncreURL = "/ext/internal/archive/channel/content/change"
)

// Dao dao .
type Dao struct {
	c                       *conf.Config
	db                      *sql.DB
	showDB                  *sql.DB
	httpR                   *bm.Client
	pgcFullURL, pgcIncreURL string
	ela                     *elastic.Elastic
}

// New init mysql db .
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:           c,
		db:          sql.NewMySQL(c.DB.Goblin),
		showDB:      sql.NewMySQL(c.DB.Show),
		httpR:       bm.NewClient(c.SearchClient),
		pgcFullURL:  c.Host.PgcURI + _pgcFullURL,
		pgcIncreURL: c.Host.PgcURI + _pgcIncreURL,
		ela:         elastic.NewElastic(c.Es),
	}
	return
}

// Close close the resource .
func (d *Dao) Close() {
}

// Ping dao ping .
func (d *Dao) Ping(c context.Context) error {
	return nil
}

// PromError stat and log .
func PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}
