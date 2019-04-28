package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// PassportDetail get detail.
func (d *Dao) PassportDetail(c context.Context, mid int64) (res *model.PassportDetail, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	// params.Set("access_key", accessKey)
	params.Set("mid", strconv.FormatInt(mid, 10))
	var resp struct {
		Code int                   `json:"code"`
		Info *model.PassportDetail `json:"data"`
	}
	req, err := d.httpR.NewRequest("GET", d.detailURI, ip, params)
	if err != nil {
		err = errors.Wrap(err, "dao passport detail")
		return
	}
	// req.Header.Set("Cookie", cookie)
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	if err = d.httpR.Do(c, req, &resp); err != nil {
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

// PassportProfile is.
func (d *Dao) PassportProfile(c context.Context, mid int64, ip string) (res *model.PassportProfile, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var resp struct {
		Code int                    `json:"code"`
		Data *model.PassportProfile `json:"data"`
	}
	if err = d.httpP.Get(c, d.profileURI, ip, params, &resp); err != nil {
		err = errors.Wrap(err, "dao passport profile")
		return nil, err
	}
	if resp.Code != 0 {
		err = ecode.Int(resp.Code)
		err = errors.WithStack(err)
		return
	}
	res = resp.Data
	return
}
