package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-common/app/service/ep/footman/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/orm"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"

	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
)

// Dao dao
type Dao struct {
	c          *conf.Config
	httpClient *xhttp.Client
	email      *gomail.Dialer
	db         *gorm.DB
	cache      *fanout.Fanout
	mc         *memcache.Pool
	expire     int32
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c: c,
		//email:      gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.Username, c.Mail.Password),
		httpClient: xhttp.NewClient(c.HTTPClient),
		cache:      fanout.New("mcCache", fanout.Worker(1), fanout.Buffer(1024)),
		mc:         memcache.NewPool(c.Memcache.Config),
		expire:     int32(time.Duration(c.Memcache.Expire) / time.Second),
	}
	if c.ORM != nil {
		dao.db = orm.NewMySQL(c.ORM)
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}

	if d.mc != nil {
		d.mc.Close()
	}
}

// Ping verify server is ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if d.db != nil {
		if err = d.db.DB().Ping(); err != nil {
			log.Info("dao.cloudDB.Ping() error(%v)", err)
		}
	}
	return
}

func (d *Dao) newRequest(method, url string, v interface{}) (req *http.Request, err error) {
	body := &bytes.Buffer{}
	if method != http.MethodGet {
		if err = json.NewEncoder(body).Encode(v); err != nil {
			log.Error("json encode value(%s) err(%v) ", v, err)
			return
		}
	}
	if req, err = http.NewRequest(method, url, body); err != nil {
		log.Error("http new request url(%s) err(%v)", url, err)
	}
	return
}

// cacheSave cache Save.
func (d *Dao) cacheSave(c context.Context, cacheItem *memcache.Item) {
	var f = func(ctx context.Context) {
		var (
			conn = d.mc.Get(c)
			err  error
		)
		defer conn.Close()
		if err = conn.Set(cacheItem); err != nil {
			log.Error("Add Cache conn.Set(%s) error(%v)", cacheItem.Key, err)
		}
	}
	if err := d.cache.Do(c, f); err != nil {
		log.Error("ReleaseName cache save err(%v)", err)
	}
}
