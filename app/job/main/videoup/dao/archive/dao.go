package archive

import (
	"context"
	"go-common/app/job/main/videoup/conf"
	xredis "go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	xhttp "go-common/library/net/http/blademaster"
)

// Dao is redis dao.
type Dao struct {
	c            *conf.Config
	db           *xsql.DB
	rdb          *xsql.DB
	coverRds     *xredis.Pool
	coverExpire  int32
	client       *xhttp.Client
	statURI      string
	recommendURI string
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		db:           xsql.NewMySQL(c.DB.Archive),
		rdb:          xsql.NewMySQL(c.DB.ArchiveRead),
		coverRds:     xredis.NewPool(c.Redis),
		coverExpire:  86400 * 15,
		client:       xhttp.NewClient(c.HTTPClient),
		statURI:      c.Host.API + "/x/internal/v2/archive/stat",
		recommendURI: c.Host.RecCover + "/cover_recomm",
	}
	return d
}

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(c)
}

// Close fn
func (d *Dao) Close() {
	d.coverRds.Close()
}

// Ping hbase
func (d *Dao) Ping(c context.Context) (err error) {
	return nil
}
