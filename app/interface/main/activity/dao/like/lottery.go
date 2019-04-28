package like

import (
	"context"
	"net/url"
	"strconv"

	l "go-common/app/interface/main/activity/model/like"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// AddLotteryTimes .
func (d *Dao) AddLotteryTimes(c context.Context, sid, mid int64) (err error) {
	params := url.Values{}
	params.Set("act_id", strconv.FormatInt(sid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Get(c, d.addLotteryTimesURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		err = errors.Wrapf(err, "d.client.Get(%s)", d.addLotteryTimesURL+"?"+params.Encode())
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
	}
	return
}

// LotteryIndex .
func (d *Dao) LotteryIndex(c context.Context, actID, platform, source, mid int64) (res *l.Lottery, err error) {
	params := url.Values{}
	params.Set("act_id", strconv.FormatInt(actID, 10))
	params.Set("platform", strconv.FormatInt(platform, 10))
	params.Set("source", strconv.FormatInt(source, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	res = new(l.Lottery)
	if err = d.client.Get(c, d.lotteryIndexURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		err = errors.Wrapf(err, "d.client.NewRequest(%s)", d.lotteryIndexURL)
	}
	return
}
