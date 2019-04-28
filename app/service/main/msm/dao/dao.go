package dao

import (
	"go-common/app/service/main/msm/conf"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao.
type Dao struct {
	client     *bm.Client
	db         *sql.DB
	treeHost   string
	platformID string
}

// New new dao.
func New(c *conf.Config) *Dao {
	d := &Dao{
		db:         sql.NewMySQL(c.Mysql),
		client:     bm.NewClient(c.HTTPClient),
		treeHost:   c.Tree.Host,
		platformID: c.Tree.PlatformID,
	}
	return d
}

// Close close mysql resource.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
