package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_isHotURL = "http://api.bilibili.co/x/internal/v2/reply/ishot"
)

// IsOriginHot return is origin hot reply.
func (d *Dao) IsOriginHot(ctx context.Context, oid, rpID int64, tp int) (isHot bool, err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("rpid", strconv.FormatInt(rpID, 10))
	params.Set("type", strconv.Itoa(tp))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			IsHot bool `json:"isHot"`
		} `json:"data"`
	}
	if err = d.httpCli.Get(ctx, _isHotURL, "", params, &res); err != nil {
		log.Error("d.httpCli.Get(%s, %s) error(%v)", _isHotURL, params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		return
	}
	if res.Data != nil {
		isHot = res.Data.IsHot
	}
	return
}
