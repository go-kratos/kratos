package audio

import (
	"context"
	"net/url"

	"go-common/app/interface/main/app-card/model/card/audio"
	"go-common/app/interface/main/app-feed/conf"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_audios = "/x/internal/v1/audio/menus/batch"
)

type Dao struct {
	client    *httpx.Client
	getAudios string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:    httpx.NewClient(c.HTTPClient),
		getAudios: c.Host.APICo + _audios,
	}
	return
}

func (d *Dao) Audios(c context.Context, ids []int64) (aum map[int64]*audio.Audio, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("ids", xstr.JoinInts(ids))
	var res struct {
		Code int                    `json:"code"`
		Data map[int64]*audio.Audio `json:"data"`
	}
	if err = d.client.Get(c, d.getAudios, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.getAudios+"?"+params.Encode())
		return
	}
	aum = res.Data
	return
}
