package feed

import (
	"context"
	"encoding/json"
	"net/url"

	"go-common/app/job/main/app/model/feed"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

const (
	_hot    = "/data/rank/reco-tmzb.json"
	_rcmdUp = "/x/feed/rcmd/up"
)

func (d *Dao) Hots(c context.Context) (aids []int64, err error) {
	var res struct {
		Code int `json:"code"`
		List []struct {
			Aid int64 `json:"aid"`
		} `json:"list"`
	}
	if err = d.clientAsyn.Get(c, d.hot, "", nil, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.hot)
		return
	}
	if len(res.List) == 0 {
		return
	}
	aids = make([]int64, 0, len(res.List))
	for _, list := range res.List {
		if list.Aid != 0 {
			aids = append(aids, list.Aid)
		}
	}
	return
}

func (d *Dao) UpRcmdCache(c context.Context, is []*feed.RcmdItem) (err error) {
	params := url.Values{}
	var b []byte
	if b, err = json.Marshal(is); err != nil {
		err = errors.Wrapf(err, "%v", is)
		return
	}
	params.Set("item", string(b))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.rcmdUp, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.rcmdUp+"?"+params.Encode())
	}
	return
}
