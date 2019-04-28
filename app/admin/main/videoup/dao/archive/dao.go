package archive

import (
	"context"

	"go-common/app/admin/main/videoup/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao is redis dao.
type Dao struct {
	c *conf.Config
	// db
	db   *sql.DB
	rddb *sql.DB
	// redis
	redis *redis.Pool
	// hbase
	hbase            *hbase.Client
	userCardURL      string
	addQAVideoURL    string
	clientW, clientR *bm.Client
	creativeDB       *sql.DB
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:             c,
		db:            sql.NewMySQL(c.DB.Archive),
		rddb:          sql.NewMySQL(c.DB.ArchiveRead),
		redis:         redis.NewPool(c.Redis.Track.Config),
		hbase:         hbase.NewClient(&c.HBase.Config),
		userCardURL:   c.Host.Account + "/api/member/getCardByMid",
		addQAVideoURL: c.Host.Task + "/vt/video/add",
		clientW:       bm.NewClient(c.HTTPClient.Write),
		clientR:       bm.NewClient(c.HTTPClient.Read),
		creativeDB:    sql.NewMySQL(c.DB.Creative),
	}
	return d
}

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}

// Close close dao.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.creativeDB != nil {
		d.creativeDB.Close()
	}
	d.redis.Close()
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}
