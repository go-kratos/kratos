package feed

import (
	"context"
	"net/url"
	"time"

	"go-common/app/job/main/app/model/feed"
	"go-common/library/ecode"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_tags = "/x/internal/tag/archive/multi/tags"
)

func (d *Dao) Tags(c context.Context, aids []int64, now time.Time) (tm map[string][]*feed.Tag, err error) {
	params := url.Values{}
	params.Set("aids", xstr.JoinInts(aids))
	var res struct {
		Code int                    `json:"code"`
		Data map[string][]*feed.Tag `json:"data"`
	}
	if err = d.clientAsyn.Get(c, d.tags, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.tags+"?"+params.Encode())
		return
	}
	tm = res.Data
	return
}
