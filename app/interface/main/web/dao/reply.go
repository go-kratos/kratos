package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_hotURI      = "/x/internal/v2/reply/hot"
	_archiveType = 1
)

// Hot reply hot.
func (d *Dao) Hot(c context.Context, aid int64) (rs *model.ReplyHot, err error) {
	var (
		params   = url.Values{}
		remoteIP = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("type", strconv.FormatInt(_archiveType, 10))
	var res struct {
		Code int             `json:"code"`
		Data *model.ReplyHot `json:"data"`
	}
	if err = d.httpR.Get(c, d.replyHotURL, remoteIP, params, &res); err != nil {
		err = errors.Wrapf(err, "replyHot url(%s)", d.replyHotURL+"?"+params.Encode())
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		return
	}
	rs = &model.ReplyHot{Replies: res.Data.Replies, Page: res.Data.Page}
	return
}
