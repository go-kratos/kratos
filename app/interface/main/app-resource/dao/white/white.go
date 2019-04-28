package white

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// Dao bplus
type Dao struct {
	client *httpx.Client
}

// New bplus
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPClient),
	}
	return
}

// WhiteVerify white verify
func (d *Dao) WhiteVerify(c context.Context, mid int64, urlStr string) (ok bool, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Status int `json:"status"`
		}
	}
	if err = d.client.Get(c, urlStr, "", params, &res); err != nil {
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("WhiteVerify url(%s) response(%s)", urlStr+"?"+params.Encode(), b)
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), urlStr+"?"+params.Encode())
		return
	}
	if res.Data.Status == 1 {
		ok = true
	}
	return
}
