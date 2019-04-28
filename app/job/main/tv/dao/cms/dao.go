package cms

import (
	"go-common/app/job/main/tv/conf"
	"go-common/library/database/sql"
	httpx "go-common/library/net/http/blademaster"
)

// Dao dao.
type Dao struct {
	conf   *conf.Config
	DB     *sql.DB
	client *httpx.Client
}

// New create a instance of Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:   c,
		DB:     sql.NewMySQL(c.Mysql),
		client: httpx.NewClient(conf.Conf.HTTPClient),
	}
	return
}
