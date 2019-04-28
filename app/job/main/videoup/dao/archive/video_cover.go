package archive

import (
	"context"

	"go-common/library/ecode"
	"go-common/library/log"
	"net/url"
)

// AICover get covers from ai
func (d *Dao) AICover(c context.Context, filename string) (covers []string, err error) {
	params := url.Values{}
	params.Set("filename", filename)
	params.Set("from", "videoup-job")
	var res struct {
		Code int      `json:"code"`
		Data []string `json:"data"`
	}

	if err = d.client.Get(c, d.recommendURI, "", params, &res); err != nil {
		log.Error("AICover error(%v), url(%s)", err, d.recommendURI+"?"+params.Encode())
		return
	}
	if res.Code != 0 {
		err = ecode.CreativeDataErr
		log.Error("AICover code not 0, url(%s) res(%v)", d.recommendURI+"?"+params.Encode(), err)
		return
	}

	covers = res.Data
	return
}
