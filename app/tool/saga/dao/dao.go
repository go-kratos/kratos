package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-common/app/tool/saga/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/orm"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Dao def
type Dao struct {
	// cache
	httpClient       *bm.Client
	mysql            *gorm.DB
	mcMR             *memcache.Pool
	redis            *redis.Pool
	hbase            *hbase.Client
	mcMRRecordExpire int32
}

// New create instance of Dao
func New() (d *Dao) {
	d = &Dao{
		httpClient:       bm.NewClient(conf.Conf.HTTPClient),
		mysql:            orm.NewMySQL(conf.Conf.ORM),
		mcMR:             memcache.NewPool(conf.Conf.Memcache.MR),
		redis:            redis.NewPool(conf.Conf.Redis),
		hbase:            hbase.NewClient(conf.Conf.HBase.Config),
		mcMRRecordExpire: int32(time.Duration(conf.Conf.Memcache.MRRecordExpire) / time.Second),
	}
	return
}

// Ping dao.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		return
	}
	if err = d.pingMC(c); err != nil {
		return
	}
	return d.mysql.DB().Ping()
}

// Close dao.
func (d *Dao) Close() {
	if d.mcMR != nil {
		d.mcMR.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	if d.mysql != nil {
		d.mysql.Close()
	}
	if d.hbase != nil {
		d.hbase.Close()
	}
}

func (d *Dao) newRequest(method, url string, v interface{}) (req *http.Request, err error) {
	body := &bytes.Buffer{}
	if method != http.MethodGet {
		if err = json.NewEncoder(body).Encode(v); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if req, err = http.NewRequest(method, url, body); err != nil {
		err = errors.WithStack(err)
	}
	return
}
