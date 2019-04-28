package live

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/live"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_list        = "/room/v1/RoomMng/allLivingRoomInfo"
	_bnj2019Conf = "/activity/v0/bainian/config"
)

// Dao is space dao
type Dao struct {
	client  *httpx.Client
	list    string
	bnj2019 string
}

// New initial space dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:  httpx.NewClient(c.HTTPAsync),
		list:    c.Host.APILiveCo + _list,
		bnj2019: c.Host.APILiveCo + _bnj2019Conf,
	}
	return
}

// Living is get living rooms from api
func (d *Dao) Living(c context.Context) (ll []*live.Live, err error) {
	params := url.Values{}
	params.Set("filter_user_cover", "0")
	params.Set("need_broadcast_type", "1")
	params.Set("extra_fields[]", "title")
	var res struct {
		Code int              `json:"code"`
		Data []*live.RoomInfo `json:"data"`
	}
	if err = d.client.Get(c, d.list, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.list)
		return
	}
	for _, info := range res.Data {
		l := &live.Live{}
		l.LiveChange(info)
		ll = append(ll, l)
	}
	return
}

// Bnj2019Conf 直播控制白名单
func (d *Dao) Bnj2019Conf(c context.Context) (greyStatus int, mids []int64, err error) {
	var res struct {
		Code int `json:"code"`
		Data struct {
			GreyStatus int    `json:"grey_status"`
			GreyUids   string `json:"grey_uids"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.bnj2019, "", nil, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.bnj2019)
		return
	}
	greyStatus = res.Data.GreyStatus
	if greyStatus == 1 {
		midsStr := strings.Split(res.Data.GreyUids, ",")
		for _, midStr := range midsStr {
			var mid int64
			if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
				err = errors.New(fmt.Sprintf("live grey_uids(%s)", res.Data.GreyUids))
				return
			}
			mids = append(mids, mid)
		}
	}
	return
}
