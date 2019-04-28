package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/service/main/relation/model"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// PassportDetail get passport detail from passport service through http
func (d *Dao) PassportDetail(c context.Context, mid int64, ip string) (res *model.PassportDetail, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var resp struct {
		Code int                   `json:"code"`
		Data *model.PassportDetail `json:"data"`
	}
	req, err := d.client.NewRequest("GET", d.detailURI, ip, params)
	if err != nil {
		err = errors.Wrap(err, "dao passport detail")
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	if err = d.client.Do(c, req, &resp); err != nil {
		err = errors.Wrap(err, "dao passport detail")
		return
	}
	if resp.Code != 0 {
		err = ecode.Int(resp.Code)
		err = errors.Wrap(err, "dao passport detail")
		return
	}
	res = resp.Data
	return
}
