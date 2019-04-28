package data

import (
	"context"
	"errors"

	"go-common/app/job/main/videoup-report/conf"
	"go-common/app/job/main/videoup-report/model/data"
	"go-common/library/ecode"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
)

const (
	_hotArc        = "/data/rank/reco-app-remen-pre.json"
	_monitorNotify = "/va/monitor/notify"
	_replyChange   = "/x/internal/v2/reply/subject/state"
	_replyInfo     = "/x/internal/v2/reply/subject"
	_profitUpState = "/allowance/api/x/admin/growup/up/account/state"
)

//Dao dao
type Dao struct {
	c                                                      *conf.Config
	moniNotifyURL                                          string
	hotArcURL, replyInfoURL, replyChangeURL, upProfitState string
	client, clientWriter                                   *xhttp.Client
}

//New new
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:              c,
		hotArcURL:      c.Host.Data + _hotArc,
		moniNotifyURL:  c.Host.Archive + _monitorNotify,
		replyInfoURL:   c.Host.API + _replyInfo,
		replyChangeURL: c.Host.API + _replyChange,
		upProfitState:  c.Host.Profit + _profitUpState,
		client:         xhttp.NewClient(c.HTTPClient.Read),
		clientWriter:   xhttp.NewClient(c.HTTPClient.Write),
	}
	return
}

// HotArchive get hot archives which need rechecking
func (d *Dao) HotArchive(c context.Context) (aids []int64, err error) {
	res := &data.HotArchiveRes{}
	if err = d.client.Get(c, d.hotArcURL, "", nil, &res); err != nil {
		log.Error("d.HotArchive() error(%v)", err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.New("data api bad response")
		log.Error("d.HotArchive() bad code(%d)", res.Code)
		return
	}
	for _, item := range res.List {
		aids = append(aids, item.Aid)
	}
	return
}
