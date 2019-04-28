package dynamic

import (
	"context"
	"encoding/json"

	"go-common/app/interface/main/app-feed/conf"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

var (
	_dynamicGet     = "/dynamic_svr/v1/dynamic_svr/dynamic_new"
	_dynamicHistory = "/dynamic_svr/v1/dynamic_svr/dynamic_history"
)

// Dao .
type Dao struct {
	client *httpx.Client
	host   string
}

// New .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPDynamic),
		host:   c.Host.Dynamic,
	}
	return
}

func (d *Dao) dynamicSrv(c context.Context, reqPath string, params string) (res json.RawMessage, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	uri := d.host + reqPath + "?" + params
	if err = d.client.Get(c, uri, ip, nil, &res); err != nil {
		err = errors.Wrap(err, uri)
	}
	return
}

// DynamicHistory .
func (d *Dao) DynamicHistory(c context.Context, params string) (res json.RawMessage, err error) {
	res, err = d.dynamicSrv(c, _dynamicHistory, params)
	return
}

// DynamicCount .
func (d *Dao) DynamicCount(c context.Context, params string) (res json.RawMessage, err error) {
	res, err = d.dynamicSrv(c, _dynamicGet, params)
	return
}

// DynamicNew .
func (d *Dao) DynamicNew(c context.Context, params string) (res json.RawMessage, err error) {
	res, err = d.dynamicSrv(c, _dynamicGet, params)
	return
}
