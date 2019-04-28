package show

import (
	"context"
	"net/url"

	"go-common/app/interface/main/app-card/model/card/show"
	"go-common/app/interface/main/app-feed/conf"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_getCard = "/api/ticket/project/getcard"
)

// Dao is show dao.
type Dao struct {
	// http client
	client *httpx.Client
	// live
	getCard string
}

// New new a bangumi dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client:  httpx.NewClient(c.HTTPShow),
		getCard: c.Host.Show + _getCard,
	}
	return d
}

func (d *Dao) Card(c context.Context, ids []int64) (rs map[int64]*show.Shopping, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("id", xstr.JoinInts(ids))
	params.Set("for", "1")
	params.Set("price", "1")
	var res struct {
		Code int              `json:"errno"`
		Data []*show.Shopping `json:"data"`
	}
	if err = d.client.Get(c, d.getCard, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.getCard+"?"+params.Encode())
		return
	}
	rs = make(map[int64]*show.Shopping, len(res.Data))
	for _, r := range res.Data {
		rs[r.ID] = r
	}
	return
}
