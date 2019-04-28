package dao

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/reply/conf"
	"go-common/app/admin/main/reply/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

const (
	// api
	_apiSearch       = "/api/reply/internal/search"
	_apiSearchUpdate = "/api/reply/internal/update"
	// index
	_searchIdxReply      = "reply"
	_searchIdxReport     = "replyreport"
	_searchIdxMonitor    = "replymonitor"
	_searchIdxTimeFormat = "2006-01-02 15:03:04"
)

var zeroTime = time.Time{}

func (d *Dao) SearchReplyV3(c context.Context, sp *model.SearchParams, page, pageSize int64) (res *model.SearchResult, err error) {
	var (
		end      = sp.End
		begin    = sp.Begin
		business = "reply_list"
	)
	if end == zeroTime {
		end = time.Now()
	}
	if begin == zeroTime {
		begin = end.Add(-time.Hour * 24 * 30)
	}
	r := d.es.NewRequest(business).IndexByTime("reply_list", elastic.IndexTypeWeek, begin, end).
		WhereEq("type", fmt.Sprint(sp.Type)).Order(sp.Order, sp.Sort).Pn(int(page)).Ps(int(pageSize)).
		WhereRange("ctime", begin.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), elastic.RangeScopeLcRc)
	if sp.Oid != 0 {
		r = r.WhereEq("oid", strconv.FormatInt(sp.Oid, 10))
	}
	if sp.TypeIds != "" {
		r = r.WhereIn("typeid", strings.Split(sp.TypeIds, ","))
	}
	if sp.Keyword != "" {
		r = r.WhereLike([]string{"message"}, []string{sp.Keyword}, true, elastic.LikeLevelLow)
	}
	if sp.KeywordHigh != "" {
		r = r.WhereLike([]string{"message_middle"}, []string{sp.KeywordHigh}, true, elastic.LikeLevelMiddle)
		r = r.OrderScoreFirst(false)
	}
	if sp.UID != 0 {
		r = r.WhereEq("mid", sp.UID)
	}
	if sp.Uname != "" {
		r = r.WhereEq("replier", sp.Uname)
	}
	if sp.AdminID != 0 {
		r = r.WhereEq("adminid", sp.AdminID)
	}
	if sp.States != "" {
		r = r.WhereIn("state", strings.Split(sp.States, ","))
	}
	if sp.IP != 0 {
		var ip = make([]byte, 4)
		binary.BigEndian.PutUint32(ip, uint32(sp.IP))
		r = r.WhereEq("ip", net.IPv4(ip[0], ip[1], ip[2], ip[3]).String())
	}
	if sp.Attr != "" {
		r = r.WhereIn("attr", strings.Split(sp.Attr, ","))
	}
	if sp.AdminName != "" {
		r = r.WhereEq("admin_name", sp.AdminName)
	}
	result := new(struct {
		Code    int
		Page    *model.Page
		Order   string
		Sort    string
		Result  []*model.SearchReply
		Message string
	})
	log.Warn("search params: %s", r.Params())
	err = r.Scan(c, &result)
	if err != nil || result.Code != 0 {
		log.Error("SearchReplyV3 r.Scan(%v) error:(%v)", c, err)
		return
	}
	res = new(model.SearchResult)
	res.Result = result.Result
	res.Code = result.Code
	res.Page = result.Page.Num
	res.PageSize = result.Page.Size
	if res.PageSize > 0 {
		res.PageCount = result.Page.Total / result.Page.Size
	}
	res.Total = result.Page.Total
	res.Order = result.Order
	res.Message = result.Message
	return
}

// SearchAdminLog search adminlog
func (d *Dao) SearchAdminLog(c context.Context, rpids []int64) (res []*model.SearchAdminLog, err error) {
	if len(rpids) == 0 {
		return
	}
	r := d.es.NewRequest("reply_admin_log").Index("replyadminlog").Pn(int(1)).Ps(len(rpids)).WhereIn("rpid", rpids)
	result := new(struct {
		Code    int
		Page    *model.Page
		Order   string
		Sort    string
		Result  []*model.SearchAdminLog
		Message string
	})
	log.Warn("search params: %s", r.Params())
	err = r.Scan(c, &result)
	if err != nil || result.Code != 0 {
		log.Error("SearchAdminLog r.Scan(%v) error:(%v)", c, err)
		return
	}
	res = result.Result
	return
}

// SearchReply search reply from ES.
func (d *Dao) SearchReply(c context.Context, p *model.SearchParams, page, pageSize int64) (res *model.SearchResult, err error) {
	params := url.Values{}
	params.Set("appid", _searchIdxReply)
	params.Set("type", fmt.Sprint(p.Type))
	params.Set("sort", p.Sort)
	params.Set("order", p.Order)
	params.Set("page", fmt.Sprint(page))
	params.Set("pagesize", fmt.Sprint(pageSize))
	if p.Oid != 0 {
		params.Set("oid", strconv.FormatInt(p.Oid, 10))
	}
	if p.TypeIds != "" {
		params.Set("typeids", p.TypeIds)
	}
	if p.Keyword != "" {
		params.Set("keyword", p.Keyword)
	}
	if p.UID != 0 {
		params.Set("uid", strconv.FormatInt(p.UID, 10))
	}
	if p.Uname != "" {
		params.Set("nickname", p.Uname)
	}
	if p.AdminID != 0 {
		params.Set("adminid", strconv.FormatInt(p.AdminID, 10))
	}
	if p.Begin != zeroTime {
		params.Set("start_time", p.Begin.Format(model.DateFormat))
	}
	if p.End != zeroTime {
		params.Set("end_time", p.End.Format(model.DateFormat))
	}
	if p.States != "" {
		params.Set("states", p.States)
	}
	if p.IP != 0 {
		params.Set("ip", strconv.FormatInt(p.IP, 10))
	}
	if p.Attr != "" {
		params.Set("attr", p.Attr)
	}
	if p.AdminName != "" {
		params.Set("admin_name", p.AdminName)
	}
	res = &model.SearchResult{}
	uri := conf.Conf.Host.Search + _apiSearch
	if err = d.httpClient.Get(c, uri, "", params, res); err != nil {
		log.Error("searchReply error(%v)", err)
		return
	}
	if res.Code != 0 {
		err = model.ErrSearchReply
		log.Error("searchReply:%+v error(%v)", res, err)
	}
	return
}

// SearchMonitor return search monitor reply from ES.
func (d *Dao) SearchMonitor(c context.Context, sp *model.SearchMonitorParams, page, pageSize int64) (res *model.SearchMonitorResult, err error) {
	var (
		fields   []string
		keywords []string
		order    string
		sort     = "desc"
	)
	// NOTE:这里之前order 跟 sort 搞反了
	if sp.Sort != "" {
		order = sp.Sort
	}
	if sp.Order != "" {
		sort = sp.Order
	}
	r := d.es.NewRequest("reply_monitor").Index(_searchIdxMonitor).
		WhereEq("type", fmt.Sprint(sp.Type)).
		Order(order, sort).Pn(int(page)).Ps(int(pageSize))
	// mode=0 所有监控方式, mode=1 monitor, mode=2, 先审后发
	if sp.Mode == 0 {
		r = r.WhereOr("monitor", true).WhereOr("audit", true)
	} else if sp.Mode == 1 {
		r = r.WhereEq("monitor", true)
	} else if sp.Mode == 2 {
		r = r.WhereEq("audit", true)
	}
	if sp.Oid > 0 {
		r = r.WhereEq("oid", fmt.Sprint(sp.Oid))
	}
	if sp.UID > 0 {
		r = r.WhereEq("mid", fmt.Sprint(sp.UID))
	}
	if sp.NickName != "" {
		fields = append(fields, "uname")
		keywords = append(keywords, sp.NickName)
	}
	if sp.Keyword != "" {
		fields = append(fields, "title")
		keywords = append(keywords, sp.Keyword)
	}
	if fields != nil && keywords != nil {
		r = r.WhereLike(fields, keywords, true, elastic.LikeLevelLow)
	}
	result := new(struct {
		Code    int
		Page    *model.Page
		Order   string
		Sort    string
		Result  []*model.SearchMonitor
		Message string
	})
	res = &model.SearchMonitorResult{}
	log.Warn(r.Params())
	err = r.Scan(c, &result)
	if err != nil || result.Code != 0 {
		log.Error("r.Scan(%v) error:(%v)", c, err)
		return
	}
	res.Result = result.Result
	res.Code = result.Code
	res.Page = result.Page.Num
	res.PageSize = result.Page.Size
	if res.PageSize > 0 {
		res.PageCount = result.Page.Total / result.Page.Size
	}
	res.Total = result.Page.Total
	res.Order = result.Order
	res.Message = result.Message
	oids := make([]int64, len(res.Result))
	var tp int32
	for idx, r := range res.Result {
		oids[idx] = r.Oid
		tp = int32(r.Type)
	}
	results, err := d.SubMCount(c, oids, tp)
	if err != nil {
		log.Error("SubMCount(%v,%v) error", oids, tp)
		return
	}
	for i, reply := range res.Result {
		res.Result[i].MCount = results[reply.Oid]
		res.Result[i].OidStr = strconv.FormatInt(res.Result[i].Oid, 10)
	}
	return
}

// UpSearchMonitor update monitor to search data.
func (d *Dao) UpSearchMonitor(c context.Context, sub *model.Subject, remark string) (err error) {
	m := make(map[string]interface{})
	m["oid"] = sub.Oid
	m["type"] = sub.Type
	if sub.AttrVal(model.SubAttrMonitor) == model.AttrYes {
		m["monitor"] = true
	} else {
		m["monitor"] = false
	}
	if sub.AttrVal(model.SubAttrAudit) == model.AttrYes {
		m["audit"] = true
	} else {
		m["audit"] = false
	}
	m["remark"] = remark
	us := d.es.NewUpdate("reply_monitor").Insert()
	us.AddData("replymonitor", m)
	err = us.Do(c)
	if err != nil {
		err = model.ErrSearchReport
		log.Error("upSearchMonitor error(%v)", err)
		return
	}
	return
}

// SearchReport search reports from ES.
func (d *Dao) SearchReport(c context.Context, sp *model.SearchReportParams, page, pageSize int64) (res *model.SearchReportResult, err error) {
	params := url.Values{}
	params.Set("appid", "replyreport")
	params.Set("type", fmt.Sprint(sp.Type))
	params.Set("page", fmt.Sprint(page))
	params.Set("pagesize", fmt.Sprint(pageSize))
	if sp.Oid != 0 {
		params.Set("oid", fmt.Sprint(sp.Oid))
	}
	if sp.UID != 0 {
		params.Set("uid", fmt.Sprint(sp.UID))
	}
	if sp.Reason != "" {
		params.Set("reason", sp.Reason)
	}
	if sp.Typeids != "" {
		params.Set("typeids", sp.Typeids)
	}
	if sp.Keyword != "" {
		params.Set("keyword", sp.Keyword)
	}
	if sp.Nickname != "" {
		params.Set("nickname", sp.Nickname)
	}
	if sp.States != "" {
		params.Set("states", sp.States)
	}
	if sp.StartTime != "" {
		params.Set("start_time", sp.StartTime)
	}
	if sp.EndTime != "" {
		params.Set("end_time", sp.EndTime)
	}
	if sp.Order != "" {
		params.Set("order", sp.Order)
	}
	if sp.Sort != "" {
		params.Set("sort", sp.Sort)
	}
	res = &model.SearchReportResult{}
	uri := conf.Conf.Host.Search + _apiSearch
	if err = d.httpClient.Get(c, uri, "", params, res); err != nil {
		log.Error("searchReport error(%v)", err)
		return
	}
	if res.Code != 0 {
		err = model.ErrSearchReport
		log.Error("searchReport:%+v error(%v)", res, err)
	}
	return
}

// MonitorStats return search monitor stats from ES.
func (d *Dao) MonitorStats(c context.Context, mode, page, pageSize int64, adminIDs, sort, order, startTime, endTime string) (res *model.StatsMonitorResult, err error) {
	params := url.Values{}
	params.Set("appid", "replymonista")
	params.Set("mode", fmt.Sprint(mode))
	params.Set("page", fmt.Sprint(page))
	params.Set("pagesize", fmt.Sprint(pageSize))
	if adminIDs != "" {
		params.Set("adminids", adminIDs)
		params.Set("typeid", fmt.Sprint(model.MonitorStatsUser))
	} else {
		params.Set("typeid", fmt.Sprint(model.MonitorStatsAll))
	}
	if sort != "" {
		params.Set("sort", sort)
	}
	if order != "" {
		params.Set("order", order)
	}
	if startTime != "" {
		params.Set("start_time", startTime)
	}
	if endTime != "" {
		params.Set("end_time", endTime)
	}
	res = &model.StatsMonitorResult{}
	uri := conf.Conf.Host.Search + _apiSearch
	if err = d.httpClient.Get(c, uri, "", params, res); err != nil {
		log.Error("monitorStats error(%v)", err)
		return
	}
	if res.Code != 0 {
		err = model.ErrSearchMonitor
		log.Error("searchStats:%+v error(%v)", res, err)
	}
	return
}

// UpSearchReply update search reply index.
func (d *Dao) UpSearchReply(c context.Context, rps map[int64]*model.Reply, newState int32) (err error) {
	if len(rps) <= 0 {
		return
	}
	stales := d.es.NewUpdate("reply_list")
	for _, rp := range rps {
		m := make(map[string]interface{})
		m["id"] = rp.ID
		m["state"] = newState
		m["mtime"] = rp.MTime.Time().Format("2006-01-02 15:04:05")
		m["oid"] = rp.Oid
		m["type"] = rp.Type
		if rp.Content != nil {
			m["message"] = rp.Content.Message
		}
		stales = stales.AddData(d.es.NewUpdate("reply_list").IndexByTime("reply_list", elastic.IndexTypeWeek, rp.CTime.Time()), m)

	}
	err = stales.Do(c)
	if err != nil {
		log.Error("upSearchReply update stales(%s) failed!err:=%v", stales.Params(), err)
		return
	}
	log.Info("upSearchReply:stale:%s ret:%+v", stales.Params(), err)
	return
}

// UpSearchReport update search report index.
func (d *Dao) UpSearchReport(c context.Context, rpts map[int64]*model.Report, rpState *int32) (err error) {
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	params := url.Values{}
	params.Set("appid", _searchIdxReport)
	values := make([]map[string]interface{}, 0)
	rps := make(map[int64]*model.Reply)
	for _, rpt := range rpts {
		if int64(rpt.ReplyCtime) != 0 && rpState != nil {
			rps[rpt.RpID] = &model.Reply{
				ID:    rpt.RpID,
				Oid:   rpt.Oid,
				Type:  rpt.Type,
				MTime: rpt.MTime,
				CTime: rpt.ReplyCtime,
				State: *rpState,
			}
		}
		v := make(map[string]interface{})
		v["id"] = fmt.Sprintf("%d_%d_%d", rpt.RpID, rpt.Oid, rpt.Type)
		v["content"] = rpt.Content
		v["reason"] = rpt.Reason
		v["state"] = rpt.State
		v["mtime"] = rpt.MTime.Time().Format(_searchIdxTimeFormat)
		v["index_time"] = rpt.CTime.Time().Format(_searchIdxTimeFormat)
		if rpt.Attr == 1 {
			v["attr"] = []int{1}
		} else {
			v["attr"] = []int{}
		}
		if rpState != nil {
			v["reply_state"] = *rpState
		}
		values = append(values, v)
	}
	b, _ := json.Marshal(values)
	params.Set("val", string(b))
	// http post
	uri := conf.Conf.Host.Search + _apiSearchUpdate
	if err = d.httpClient.Post(c, uri, "", params, &res); err != nil {
		log.Error("upSearchReport error(%v)", err)
	}
	log.Info("upSearchReport:%s post:%s ret:%+v", uri, params.Encode(), res)
	if len(rps) != 0 && rpState != nil {
		err = d.UpSearchReply(c, rps, *rpState)
	}
	return
}
