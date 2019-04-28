package kfc

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_kfcDelURI = "/x/internal/activity/kfc/deliver"
)

// KfcDelver .
func (d *Dao) KfcDelver(c context.Context, id, mid int64) (err error) {
	var httpRes struct {
		Code    int              `json:"code"`
		Data    map[string]int64 `json:"data"`
		Message string           `json:"message"`
	}
	params := url.Values{}
	params.Set("id", strconv.FormatInt(id, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	if err = d.httpClient.Post(c, d.kfcDelURL, metadata.String(c, metadata.RemoteIP), params, &httpRes); err != nil {
		err = errors.Wrap(err, "KfcDelver http")
		return
	}
	if httpRes.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(httpRes.Code), "KfcDelver msg(%s)", httpRes.Message)
	}
	return
}
