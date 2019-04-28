package archive

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_realteURL = "/recsys/related"
)

// RelateAids get relate by aid
func (d *Dao) RelateAids(c context.Context, aid int64, ip string) (aids []int64, err error) {
	params := url.Values{}
	params.Set("key", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
		Data []*struct {
			Value string `json:"value"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.relateURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.relateURL+"?"+params.Encode())
		return
	}
	if len(res.Data) != 0 {
		if aids, err = xstr.SplitInts(res.Data[0].Value); err != nil {
			err = errors.Wrap(err, res.Data[0].Value)
		}
	}
	return
}
