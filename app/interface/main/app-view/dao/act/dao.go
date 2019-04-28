package act

import (
	"context"
	"net/url"
	"strconv"

	actmdl "go-common/app/interface/main/activity/model/like"
	actrpc "go-common/app/interface/main/activity/rpc/client"
	"go-common/app/interface/main/app-view/conf"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_actInfo      = "/matsuri/api/get/videoviewinfo"
	_lotteryTimes = "/matsuri/api/get/act/mylotterytimes"
)

// Dao is elec dao.
type Dao struct {
	client       *httpx.Client
	actInfo      string
	lotteryTimes string
	actRPC       *actrpc.Service
}

// New elec dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:       httpx.NewClient(c.HTTPClient),
		actInfo:      c.Host.Activity + _actInfo,
		lotteryTimes: c.Host.Activity + _lotteryTimes,
		actRPC:       actrpc.New(c.ActivityRPC),
	}
	return
}

var _emptyList = []string{}

// Info mid+aid total elec info
func (d *Dao) Info(c context.Context, actID int64, randomCnt int64) (gifts, winners []string, err error) {
	var res struct {
		Code int `json:"code"`
		Data struct {
			Gifts []struct {
				Img string `json:"sponsors_logo"`
			} `json:"gifts"`
			Winner []struct {
				Gift  string `json:"gift"`
				UName string `json:"uname"`
			} `json:"winner"`
		} `json:"data"`
	}
	params := url.Values{}
	params.Set("act_id", strconv.FormatInt(actID, 10))
	params.Set("random_count", strconv.FormatInt(randomCnt, 10))
	if err = d.client.Get(c, d.actInfo, "", params, &res); err != nil {
		err = errors.Wrapf(err, "d.client.Get(%s)", d.actInfo+"?"+params.Encode())
		return _emptyList, _emptyList, err
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "d.client.Get(%s)", d.actInfo+"?"+params.Encode())
		return _emptyList, _emptyList, err
	}
	for _, v := range res.Data.Gifts {
		gifts = append(gifts, v.Img)
	}
	for _, v := range res.Data.Winner {
		winners = append(winners, v.UName+" 抽到了 "+v.Gift)
	}
	if gifts == nil {
		gifts = _emptyList
	}
	if winners == nil {
		winners = _emptyList
	}
	return
}

// LeftLotteryTimes 剩余抽奖次数
func (d *Dao) LeftLotteryTimes(c context.Context, actID, mid int64) (times int64, err error) {
	var res struct {
		Code int `json:"code"`
		Data struct {
			Times int64 `json:"times"`
		} `json:"data"`
	}
	params := url.Values{}
	params.Set("act_id", strconv.FormatInt(actID, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	if err = d.client.Get(c, d.lotteryTimes, "", params, &res); err != nil {
		err = errors.Wrapf(err, "d.client.Get(%s)", d.lotteryTimes+"?"+params.Encode())
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(res.Code), "d.client.Get(%s)", d.lotteryTimes+"?"+params.Encode())
		return
	}
	times = res.Data.Times
	return
}

// ActProtocol get act subject & protocol
func (d *Dao) ActProtocol(c context.Context, messionID int64) (protocol *actmdl.SubProtocol, err error) {
	arg := &actmdl.ArgActProtocol{Sid: messionID}
	protocol, err = d.actRPC.ActProtocol(c, arg)
	return
}
