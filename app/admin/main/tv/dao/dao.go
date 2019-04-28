package dao

import (
	"net/http"
	"time"

	"go-common/app/admin/main/tv/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/elastic"
	"go-common/library/database/orm"
	httpx "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao struct user of Dao.
type Dao struct {
	c *conf.Config
	// db
	DB *gorm.DB
	// dbshow
	DBShow     *gorm.DB
	fullURL    string
	httpSearch *httpx.Client
	client     *httpx.Client
	bfsClient  *http.Client
	esClient   *elastic.Elastic
	// memcache
	mc        *memcache.Pool
	cmsExpire int32
}

// New create a instance of Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// db
		DB: orm.NewMySQL(c.ORM),
		// dbshow
		DBShow: orm.NewMySQL(c.ORMShow),
		// http client
		fullURL:    c.HTTPSearch.FullURL,
		httpSearch: httpx.NewClient(c.HTTPSearch.ClientConfig),
		client:     httpx.NewClient(c.HTTPClient),
		bfsClient:  &http.Client{Timeout: time.Duration(c.Bfs.Timeout) * time.Millisecond},
		esClient: elastic.NewElastic(&elastic.Config{
			Host:       c.Cfg.Hosts.Manager,
			HTTPClient: c.HTTPClient,
		}),
		mc:        memcache.NewPool(c.Memcache.Config),
		cmsExpire: int32(time.Duration(c.Memcache.CmsExpire) / time.Second),
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	d.DB.LogMode(true)
	d.DBShow.LogMode(true)
}
