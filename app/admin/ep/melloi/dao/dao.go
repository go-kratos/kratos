package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/orm"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"

	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
)

// Prometheus
var (
	errorsCount = prom.BusinessErrCount
)

// Dao dao
type Dao struct {
	c          *conf.Config
	DB         *gorm.DB
	httpClient *xhttp.Client
	client     *http.Client
	MC         *memcache.Pool
	email      *gomail.Dialer
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:          c,
		httpClient: xhttp.NewClient(c.HTTPClient),
		email:      gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.Username, c.Mail.Password),
		client:     new(http.Client),
		DB:         orm.NewMySQL(c.ORM),
		MC:         memcache.NewPool(c.Memcache),
	}
	return
}

// Ping verify server is ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.DB.DB().Ping(); err != nil {
		log.Error("dao.cloudDB.Ping() error(%v)", err)
		return
	}
	return
}

func (d *Dao) newRequest(method, url string, v interface{}) (req *http.Request, err error) {
	body := &bytes.Buffer{}
	if method != http.MethodGet {
		if err = json.NewEncoder(body).Encode(v); err != nil {
			log.Error("json encode value(%s), error(%v) ", v, err)
			return
		}
	}
	if req, err = http.NewRequest(method, url, body); err != nil {
		log.Error("http new request url(%s), error(%v)", url, err)
	}
	return
}

// PromError prom error
func PromError(name string) {
	errorsCount.Incr(name)
}

// Close close the resource.
func (d *Dao) Close() {
	d.DB.Close()
}
