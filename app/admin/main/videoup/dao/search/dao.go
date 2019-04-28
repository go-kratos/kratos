package search

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/admin/main/videoup/conf"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/search"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const (
	_searchURL = "/x/admin/search/log"
)

// Dao is search dao
type Dao struct {
	c          *bm.ClientConfig
	httpClient *bm.Client
	URI        string
	es         *elastic.Config
}

var (
	d *Dao
)

// New new search dao
func New(c *conf.Config) *Dao {
	return &Dao{
		c:          c.HTTPClient.Read,
		httpClient: bm.NewClient(c.HTTPClient.Read),
		URI:        c.Host.MngSearch + _searchURL,
		es: &elastic.Config{
			Host:       c.Host.Manager,
			HTTPClient: c.HTTPClient.Search,
		},
	}
}

// OutTime 退出时间,es的group by查询,最大1000条
func (d *Dao) OutTime(c context.Context, ids []int64) (mcases map[int64][]interface{}, err error) {
	mcases = make(map[int64][]interface{})
	params := url.Values{}
	params.Set("appid", "log_audit_group")
	params.Set("group", "uid")
	params.Set("uid", xstr.JoinInts(ids))
	params.Set("business", strconv.Itoa(archive.LogClientConsumer))
	params.Set("action", strconv.Itoa(int(archive.ActionHandsOFF)))
	params.Set("ps", strconv.Itoa(len(ids)))
	res := &archive.SearchLogResult{}
	if err = d.httpClient.Get(c, d.URI, "", params, &res); err != nil {
		log.Error("log_audit_group d.httpClient.Get error(%v)", err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("log_audit_group ecode:%v", res.Code)
		return
	}
	for _, item := range res.Data.Result {
		mcases[item.UID] = []interface{}{item.Ctime}
	}
	log.Info("log_audit_group get: %s params:%s ret:%v", d.URI, params.Encode(), res)
	return
}

// InQuitList 登入登出日志
func (d *Dao) InQuitList(c context.Context, uids []int64, bt, et string) (l []*archive.InQuit, err error) {
	params := url.Values{}
	params.Set("appid", "log_audit")
	params.Set("business", strconv.Itoa(archive.LogClientConsumer))
	if len(uids) > 0 {
		params.Set("uid", xstr.JoinInts(uids))
	}
	if len(bt) > 0 && len(et) > 0 {
		params.Set("ctime_from", bt)
		params.Set("ctime_to", et)
	}
	params.Set("order", "ctime")
	params.Set("sort", "desc")
	params.Set("ps", "10000")

	res := &archive.SearchLogResult{}
	if err = d.httpClient.Get(c, d.URI, "", params, res); err != nil {
		log.Error("InQuitList d.httpClient.Get error(%v)", err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("InQuitList ecode:%v", res.Code)
		return
	}

	mapHelp := make(map[int64]*archive.InQuit)
	for i := len(res.Data.Result) - 1; i >= 0; i-- {
		item := res.Data.Result[i]
		if item.Action == "0" {
			ctime, _ := time.Parse(archive.TimeFormatSec, item.Ctime)
			iqlog := &archive.InQuit{
				Date:   ctime.Format("2006-01-02"),
				UID:    item.UID,
				Uname:  item.Uname,
				InTime: ctime.Format("15:04:05"),
			}
			mapHelp[item.UID] = iqlog
			l = append([]*archive.InQuit{iqlog}, l[:]...)
		}
		if item.Action == "1" {
			if iqlog, ok := mapHelp[item.UID]; ok {
				ctime, _ := time.Parse(archive.TimeFormatSec, item.Ctime)
				if date := ctime.Format("2006-01-02"); date == iqlog.Date {
					iqlog.OutTime = ctime.Format("15:04:05")
				} else {
					iqlog.OutTime = ctime.Format(archive.TimeFormatSec)
				}
			}
		}
	}

	return
}

// SearchCopyright search video copyright
func (d *Dao) SearchCopyright(c context.Context, kw string) (result *search.CopyrightResultData, err error) {
	var (
		ps = 30 //copyright不需要翻页，产品（计晓峰）说返回30条数据就可以
	)
	if kw == "" {
		return
	}
	es := elastic.NewElastic(d.es)
	eReq := es.NewRequest("copyright")
	eReq.Ps(ps)
	eReq.Index("copyright")
	eReq.WhereLike([]string{"name", "oname", "aka_names"}, []string{kw}, true, elastic.LikeLevelLow)
	log.Info("SearchCopyright(%s)", eReq.Params())
	if err = eReq.Scan(c, &result); err != nil {
		log.Error("s.SearchCopyright(%s) error(%v)", kw, err)
		return
	}
	if result == nil {
		result = &search.CopyrightResultData{}
	}
	if result.Result == nil {
		result.Result = []*search.Copyright{}
	}
	return
}
