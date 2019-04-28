package audio

import (
	"context"
	"go-common/app/interface/main/app-view/conf"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"net/url"

	"go-common/app/interface/main/app-view/model/view"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	// _audioByCids = "/x/internal/v1/audio/songs/by_cids"
	_audioByCids = "/audio/music-service-c/internal/songs-by-cids"
)

// Dao is archive dao.
type Dao struct {
	// http client
	client         *bm.Client
	audioByCidsURL string
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client:         bm.NewClient(c.HTTPAudio),
		audioByCidsURL: c.Host.APICo + _audioByCids,
	}
	return
}

func (d *Dao) AudioByCids(c context.Context, cids []int64) (vam map[int64]*view.Audio, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("cids", xstr.JoinInts(cids))
	var res struct {
		Code int                   `json:"code"`
		Data map[int64]*view.Audio `json:"data"`
	}
	if err = d.client.Get(c, d.audioByCidsURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.audioByCidsURL+"?"+params.Encode())
		return
	}
	vam = res.Data
	return
}
