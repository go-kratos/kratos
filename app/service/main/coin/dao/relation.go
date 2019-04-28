package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const _passportURL = "http://passport.bilibili.co/intranet/acc/bindDetail"

// PassportDetail .
type PassportDetail struct {
	BindEmail bool  `json:"bind_email"`
	BindTel   bool  `json:"bind_tel"`
	Mid       int64 `json:"mid"`
}

// PassportDetail get detail.
func (d *Dao) PassportDetail(c context.Context, mid int64) (res *PassportDetail, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var resp struct {
		Code int             `json:"code"`
		Info *PassportDetail `json:"data"`
	}
	req, err := d.httpClient.NewRequest("GET", _passportURL, ip, params)
	if err != nil {
		err = errors.Wrap(err, "dao passport detail")
		return
	}
	// req.Header.Set("Cookie", cookie)
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	if err = d.httpClient.Do(c, req, &resp); err != nil {
		err = errors.Wrap(err, "dao passport detail")
		return
	}
	if resp.Code != 0 {
		err = ecode.Int(resp.Code)
		err = errors.Wrap(err, "dao passport detail")
		return
	}
	res = resp.Info
	return
}
