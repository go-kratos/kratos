package up

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	upclient "go-common/app/service/main/up/api/gorpc"
	upapi "go-common/app/service/main/up/api/v1"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_upSpecialGroupURI = "/x/internal/uper/special/get_by_mid"
)

// Dao is search dao.
type Dao struct {
	c          *conf.Config
	up         *upclient.Service
	httpClient *httpx.Client
	UpClient   upapi.UpClient
}

// New new a search dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		up:         upclient.New(c.UPRPC),
		httpClient: httpx.NewClient(c.HTTPClient.Normal),
	}
	var err error
	if d.UpClient, err = upapi.NewClient(c.UpClient); err != nil {
		panic(err)
	}
	return d
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
