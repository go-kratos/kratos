package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_search = "/cate/search"
)

// Dao is search dao.
type Dao struct {
	client    *httpx.Client
	searchURL string
	// search duration
	searchTick string
}

//New recommend dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:     httpx.NewClient(conf.Conf.HTTPClient),
		searchURL:  c.Host.Search + _search,
		searchTick: c.Duration.Search,
	}
	return
}

// SearchList
func (d *Dao) SearchList(c context.Context, rid, build, pn, ps int, mid int64, ts time.Time, ip, order, tagName, platform, mobiApp, device string) (arcids []int64, err error) {
	advance, _ := time.ParseDuration(d.searchTick) // three weeks -504h
	starttime := ts.Add(advance)
	params := url.Values{}
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("order", order)
	params.Set("page", strconv.Itoa(pn))
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("time_from", starttime.Format("20060102"))
	params.Set("time_to", ts.Format("20060102"))
	params.Set("build", strconv.Itoa(build))
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("search_type", "video")
	params.Set("view_type", "hot_rank")
	params.Set("clientip", ip)
	if tagName != "" {
		params.Set("keyword", tagName)
	}
	params.Set("cate_id", strconv.Itoa(rid))
	var res struct {
		Code int `json:"code"`
		Data []struct {
			Aid interface{} `json:"id"`
		} `json:"result"`
	}
	if err = d.client.Get(c, d.searchURL, "", params, &res); err != nil {
		log.Error("search news url(%s) error(%v)", d.searchURL+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("search list url(%v) response(%s)", d.searchURL+"?"+params.Encode(), b)
	if res.Code != 0 && res.Code != -1 {
		log.Error("search region news url(%s) error(%v)", d.searchURL+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("search region news api response code(%v)", res)
		return
	}
	for _, arcs := range res.Data {
		var aidInt int64
		switch aid := arcs.Aid.(type) {
		case string:
			aidInt = aidToInt(aid)
		case float64:
			aidInt = int64(aid)
		}
		arcids = append(arcids, aidInt)
	}
	return
}

func aidToInt(aidstr string) (aid int64) {
	aid, _ = strconv.ParseInt(aidstr, 10, 64)
	return
}
