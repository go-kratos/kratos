package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_accErrWasFormal = -659
)

// BeFormal be formal member.
func (d *Dao) BeFormal(c context.Context, mid int64, cookie, ip string) (err error) {
	params := url.Values{}
	params.Set("Cookie", cookie)
	params.Set("mid", strconv.FormatInt(mid, 10))
	// request
	req, err := d.httpClient.NewRequest("POST", d.beFormalURI, ip, params)
	if err != nil {
		log.Error("account beformal uri(%s) error(%v)", d.beFormalURI+params.Encode(), err)
		return
	}
	req.Header.Set("Cookie", cookie)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != 0 {
		if res.Code == _accErrWasFormal {
			return
		}
		err = ecode.Int(res.Code)
		log.Error("account beformal uri(%s) error(%v)", d.beFormalURI+params.Encode(), err)
	}
	return
}
