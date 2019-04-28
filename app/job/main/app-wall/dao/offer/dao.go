package offer

import (
	"context"
	"net/url"

	"go-common/app/job/main/app-wall/conf"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_active = "/x/wall/offer/active2"
)

// Dao dao
type Dao struct {
	c      *conf.Config
	redis  *redis.Pool
	client *bm.Client
	active string
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:      c,
		redis:  redis.NewPool(c.Redis.Feed.Config),
		client: bm.NewClient(c.HTTPClient),
		active: c.Host.APP + _active,
	}
	return
}

// Close close the resource.
func (dao *Dao) Close() (err error) {
	return dao.redis.Close()
}

// Ping dao ping
func (dao *Dao) Ping(c context.Context) (err error) {
	return dao.PingRedis(c)
}

func (d *Dao) Active(c context.Context, os, imei, androidid, mac, ip string) (err error) {
	params := url.Values{}
	params.Set("os", os)
	params.Set("imei", imei)
	params.Set("androidid", androidid)
	params.Set("mac", mac)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.active, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.active+"?"+params.Encode())
		if res.Code == ecode.RequestErr.Code() {
			log.Error("%+v", err)
			err = nil
			return
		}
	}
	return
}
