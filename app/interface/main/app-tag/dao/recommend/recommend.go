package recommend

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-tag/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

const (
	_recommand = "/recommand"
	_top       = "/feed/tag/top"
)

// Dao is show dao.
type Dao struct {
	// bigdata
	top       string
	recommand string
	client    *httpx.Client
}

// New new a show dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:    httpx.NewClient(c.HTTPClient),
		top:       c.Host.Data + _top,
		recommand: c.Host.Data + _recommand,
	}
	return
}

// TagRecommand feed tag
func (d *Dao) TagRecommand(c context.Context, plat int8, rid, build int, tid int64, rn int, mid int64, buvid string) (dataAids []int64,
	ctop, cbottom xtime.Time, err error) {
	// param
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("from", "54")
	params.Set("cmd", "video")
	params.Set("timeout", "100")
	params.Set("plat", strconv.Itoa(int(plat)))
	params.Set("build", strconv.Itoa(build))
	params.Set("buvid", buvid)
	if rid != 0 {
		params.Set("region", strconv.Itoa(rid))
	}
	if tid != 0 {
		params.Set("tag", strconv.FormatInt(tid, 10))
	}
	params.Set("request_cnt", strconv.Itoa(rn))
	var res struct {
		Code int `json:"code"`
		Data []struct {
			Aid  int64  `json:"id"`
			Goto string `json:"goto"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.recommand, "", params, &res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.recommand+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("recommand tag url(%s) response(%s)", d.recommand+"?"+params.Encode(), b)
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("url(%s) res code(%d)", d.recommand+"?"+params.Encode(), res.Code)
		return
	}
	for _, data := range res.Data {
		if data.Goto == "av" {
			dataAids = append(dataAids, data.Aid)
		}
	}
	ctop = 0
	cbottom = 0
	return
}

// FeedDynamic feed tag
func (d *Dao) FeedDynamic(c context.Context, pull bool, rid int, tid int64, rn int, ctime, mid int64, now time.Time) (hotAids, newAids []int64, ctop, cbottom xtime.Time, err error) {
	var pn string
	if pull {
		pn = "1"
	} else {
		pn = "2"
	}
	// param
	params := url.Values{}
	params.Set("src", "2")
	params.Set("pn", pn)
	params.Set("mid", strconv.FormatInt(mid, 10))
	if ctime != 0 {
		params.Set("ctime", strconv.FormatInt(ctime, 10))
	}
	if rid != 0 {
		params.Set("rid", strconv.Itoa(rid))
	}
	if tid != 0 {
		params.Set("tag", strconv.FormatInt(tid, 10))
	}
	if rn != 0 {
		params.Set("rn", strconv.Itoa(rn))
	}
	var res struct {
		Code    int        `json:"code"`
		Hot     []int64    `json:"hot"`
		Data    []int64    `json:"data"`
		CTop    xtime.Time `json:"ctop"`
		CBottom xtime.Time `json:"cbottom"`
	}
	if err = d.client.Get(c, d.top, "", params, &res); err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.top+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("feed tag url(%s) response(%s)", d.top+"?"+params.Encode(), b)
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("url(%s) res code(%d)", d.top+"?"+params.Encode(), res.Code)
		return
	}
	hotAids = res.Hot
	newAids = res.Data
	ctop = res.CTop
	cbottom = res.CBottom
	return
}
