package goblin

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// VerUpdate gets upgrade info.
func (d *Dao) VerUpdate(c context.Context, ver *model.VerUpdate) (result *model.HTTPData, errCode ecode.Codes, err error) {
	var (
		appURL = d.conf.Host.ReqURL
		res    struct {
			Code    int             `json:"code"`
			Data    *model.HTTPData `json:"data"`
			Message string          `json:"message"`
		}
	)
	params := url.Values{}
	params.Set("mobi_app", ver.MobiApp)
	params.Set("build", fmt.Sprintf("%d", ver.Build))
	params.Set("channel", ver.Channel)
	params.Set("seed", fmt.Sprintf("%d", ver.Seed))
	params.Set("sdkint", fmt.Sprintf("%d", ver.Sdkint))
	params.Set("model", ver.Model)
	params.Set("old_id", ver.OldID)
	if err = d.client.Get(c, appURL, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), appURL+"?"+params.Encode())
		errCode = ecode.Int(res.Code)
		return
	}
	result = res.Data
	return
}
