package member

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

type Dao struct {
	c      *conf.Config
	client *httpx.Client
	member string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: httpx.NewClient(c.HTTPClient),
		member: c.Host.APICo + "/x/internal/creative/app/pre",
	}
	return
}

// Creative get user bcoin  doc:http://info.bilibili.co/display/coding/internal-creative#internal-creative-APP%E4%B8%AA%E4%BA%BA%E4%B8%AD%E5%BF%83
func (d *Dao) Creative(c context.Context, mid int64) (isUp, show int, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			IsUp int `json:"is_up"`
			Show int `json:"show"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.member, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.member+"?"+params.Encode())
		return
	}
	isUp = res.Data.IsUp
	show = res.Data.Show
	return
}
