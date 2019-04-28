package tag

import (
	"context"

	"go-common/app/job/main/growup/conf"

	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao is redis dao.
type Dao struct {
	c            *conf.Config
	db           *sql.DB
	client       *bm.Client
	archiveURL   string
	typeURL      string
	columnURL    string
	columnActURL string
}

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		db:           sql.NewMySQL(c.Mysql.Growup),
		client:       bm.NewClient(c.HTTPClient),
		archiveURL:   c.Host.Archive + "/manager/search",
		typeURL:      c.Host.VideoType + "/videoup/types",
		columnURL:    c.Host.ColumnType,
		columnActURL: c.Host.ColumnAct,
	}
	return
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close close connections
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// BeginTran begin transcation
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}
