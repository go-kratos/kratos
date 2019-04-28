package dao

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/esports/model"
	arcMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_esports        = "esports"
	_contest        = "esports_contests"
	_videomap       = "esports_map"
	_matchmap       = "esports_contests_map"
	_calendar       = "esports_contests_date"
	_contestFav     = "esports_fav"
	_searchPlatform = "web"
	_fromSource     = "esports_search"
	_searchType     = "vesports"
	_searchVer      = "v3"
	_orderRank      = "totalrank"
	_orderHot       = "hot"
	_orderPub       = "pubdate"
	_active         = 1
	_pageNum        = 1
	_pageSize       = 1000
)

// Search search api.
func (d *Dao) Search(c context.Context, mid int64, p *model.ParamSearch, buvid string) (rs *model.SearchEsp, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("keyword", p.Keyword)
	params.Set("platform", _searchPlatform)
	params.Set("from_source", _fromSource)
	params.Set("search_type", _searchType)
	params.Set("main_ver", _searchVer)
	params.Set("clientip", ip)
	params.Set("userid", strconv.FormatInt(mid, 10))
	params.Set("buvid", buvid)
	params.Set("page", strconv.Itoa(p.Pn))
	params.Set("pagesize", strconv.Itoa(p.Ps))
	if p.Sort == 0 {
		params.Set("order", _orderRank)
	} else if p.Sort == 1 {
		params.Set("order", _orderPub)
	} else if p.Sort == 2 {
		params.Set("order", _orderHot)
	}

	if err = d.http.Get(c, d.searchURL, ip, params, &rs); err != nil {
		log.Error("Search接口错误 Search d.http.Get(%s) error(%v)", d.searchURL+"?"+params.Encode(), err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("Search接口code错误 Search d.http.Do(%s) code error(%d)", d.searchURL, rs.Code)
		err = ecode.Int(rs.Code)
	}
	return
}

// SearchVideo search video.
func (d *Dao) SearchVideo(c context.Context, p *model.ParamVideo) (rs []*model.SearchVideo, total int, err error) {
	var res struct {
		Page   model.Page           `json:"page"`
		Result []*model.SearchVideo `json:"result"`
	}
	states := []int64{arcMdl.StateForbidFixed, arcMdl.StateOpen, arcMdl.StateOrange}
	r := d.ela.NewRequest(_esports).WhereIn("state", states).WhereEq("is_deleted", 0).Fields("aid").Index(_esports).Pn(p.Pn).Ps(p.Ps)
	if p.Gid > 0 {
		r.WhereEq("gid", p.Gid)
	}
	if p.Mid > 0 {
		r.WhereEq("matchs", p.Mid)
	}
	if p.Year > 0 {
		r.WhereEq("year", p.Year)
	}
	if p.Tid > 0 {
		r.WhereEq("teams", p.Tid)
	}
	if p.Tag > 0 {
		r.WhereEq("tags", p.Tag)
	}
	if p.Sort == 0 {
		r.Order("score", elastic.OrderDesc)
	} else if p.Sort == 1 {
		r.Order("pubtime", elastic.OrderDesc)
	} else if p.Sort == 2 {
		r.Order("click", elastic.OrderDesc)
	}
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	total = res.Page.Total
	rs = res.Result
	return
}

// SearchContest search contest.
func (d *Dao) SearchContest(c context.Context, p *model.ParamContest) (rs []*model.Contest, total int, err error) {
	var res struct {
		Page   model.Page       `json:"page"`
		Result []*model.Contest `json:"result"`
	}
	r := d.ela.NewRequest(_contest).Index(_contest).WhereEq("status", 0).Pn(p.Pn).Ps(p.Ps)
	r.Fields("id", "game_stage", "stime", "etime", "home_id", "away_id", "home_score", "away_score", "live_room", "aid", "collection", "game_state", "dic", "ctime", "mtime", "status", "sid", "mid")
	if p.Sort == 1 {
		r.Order("stime", elastic.OrderDesc)
	} else {
		r.Order("stime", elastic.OrderAsc)
	}
	if p.Gid > 0 {
		r.WhereEq("gid", p.Gid)
	}
	if p.Mid > 0 {
		r.WhereEq("mid", p.Mid)
	}
	if p.Tid > 0 {
		r.WhereOr("home_id", p.Tid).WhereOr("away_id", p.Tid)
	}
	if p.Stime != "" && p.Etime != "" {
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Stime+" 00:00:00", time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Etime+" 23:59:59", time.Local)
		r.WhereRange("stime", start.Unix(), end.Unix(), elastic.RangeScopeLcRc)
	} else if p.Stime != "" && p.Etime == "" {
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Stime+" 00:00:00", time.Local)
		r.WhereRange("stime", start.Unix(), "", elastic.RangeScopeLcRo)
	} else if p.Stime == "" && p.Etime != "" {
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Etime+" 23:59:59", time.Local)
		r.WhereRange("stime", "", end.Unix(), elastic.RangeScopeLoRc)
	}
	if p.GState != "" {
		r.WhereIn("game_state", strings.Split(p.GState, ","))
	}
	if len(p.Sids) > 0 {
		r.WhereIn("sid", p.Sids)
	}
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	total = res.Page.Total
	rs = res.Result
	return
}

// SearchContestQuery search contest.
func (d *Dao) SearchContestQuery(c context.Context, p *model.ParamContest) (rs []*model.Contest, total int, err error) {
	var res struct {
		Page   model.Page       `json:"page"`
		Result []*model.Contest `json:"result"`
	}
	r := d.ela.NewRequest(_contest).Index(_contest).WhereEq("status", 0).Pn(p.Pn).Ps(p.Ps)
	r.Fields("id", "game_stage", "stime", "etime", "home_id", "away_id", "home_score", "away_score", "live_room", "aid", "collection", "game_state", "dic", "ctime", "mtime", "status", "sid", "mid")
	if p.Sort == 1 {
		r.Order("stime", elastic.OrderDesc)
	} else {
		r.Order("stime", elastic.OrderAsc)
	}
	if p.Gid > 0 {
		r.WhereEq("gid", p.Gid)
	}
	if p.Mid > 0 {
		r.WhereEq("mid", p.Mid)
	}
	if p.Tid > 0 {
		r.WhereOr("home_id", p.Tid).WhereOr("away_id", p.Tid)
	}
	if p.Stime != "" && p.Etime != "" {
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Stime, time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Etime, time.Local)
		r.WhereRange("stime", start.Unix(), end.Unix(), elastic.RangeScopeLcRc)
	} else if p.Stime != "" && p.Etime == "" {
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Stime, time.Local)
		r.WhereRange("stime", start.Unix(), "", elastic.RangeScopeLcRo)
	} else if p.Stime == "" && p.Etime != "" {
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Etime, time.Local)
		r.WhereRange("stime", "", end.Unix(), elastic.RangeScopeLoRc)
	}
	if p.GState != "" {
		r.WhereIn("game_state", strings.Split(p.GState, ","))
	}
	if len(p.Sids) > 0 {
		r.WhereIn("sid", p.Sids)
	}
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	total = res.Page.Total
	rs = res.Result
	return
}

// FilterVideo video filter.
func (d *Dao) FilterVideo(c context.Context, p *model.ParamFilter) (rs *model.FilterES, err error) {
	var res struct {
		Page   model.Page      `json:"page"`
		Result *model.FilterES `json:"result"`
	}
	r := d.ela.NewRequest(_videomap).Index(_videomap).Pn(_pageNum).Ps(_pageSize)
	r.WhereEq("active", _active).GroupBy("group_by", "gid", nil).GroupBy("group_by", "match", nil).GroupBy("group_by", "tag", nil).GroupBy("group_by", "team", nil).GroupBy("group_by", "year", nil)
	if p.Gid > 0 {
		r.WhereEq("gid", p.Gid).WhereEq("active", _active)
	}
	if p.Mid > 0 {
		r.WhereEq("match", p.Mid).WhereEq("active", _active)
	}
	if p.Tid > 0 {
		r.WhereOr("team", p.Tid).WhereEq("active", _active)
	}
	if p.Tag > 0 {
		r.WhereEq("tag", p.Tag).WhereEq("active", _active)
	}
	if p.Year > 0 {
		r.WhereEq("year", p.Year).WhereEq("active", _active)
	}
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	rs = res.Result
	return
}

// FilterMatch match filter.
func (d *Dao) FilterMatch(c context.Context, p *model.ParamFilter) (rs *model.FilterES, err error) {
	var res struct {
		Page   model.Page      `json:"page"`
		Result *model.FilterES `json:"result"`
	}
	r := d.ela.NewRequest(_matchmap).Index(_matchmap).Pn(_pageNum).Ps(_pageSize)
	r.WhereEq("active", _active).GroupBy("group_by", "match", nil).GroupBy("group_by", "gid", nil).GroupBy("group_by", "team", nil)
	if p.Mid > 0 {
		r.WhereEq("match", p.Mid).WhereEq("active", _active)
	}
	if p.Gid > 0 {
		r.WhereEq("gid", p.Gid).WhereEq("active", _active)
	}
	if p.Tid > 0 {
		r.WhereOr("team", p.Tid).WhereEq("active", _active)
	}
	if p.Stime != "" {
		r.WhereOr("stime", p.Stime).WhereEq("active", _active)
	}
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	rs = res.Result
	return
}

// FilterCale Calendar filter.
func (d *Dao) FilterCale(c context.Context, p *model.ParamFilter) (rs map[string]int64, err error) {
	var res struct {
		Page   model.Page       `json:"page"`
		Result map[string]int64 `json:"result"`
	}
	r := d.ela.NewRequest(_calendar).Index(_matchmap).Pn(_pageNum).Ps(_pageSize).WhereEq("active", _active).WhereRange("stime", p.Stime, p.Etime, elastic.RangeScopeLcRc)
	r.GroupBy("group_by", "stime", nil)
	if p.Mid > 0 {
		r.WhereEq("match", p.Mid).WhereEq("active", _active)
	}
	if p.Gid > 0 {
		r.WhereEq("gid", p.Gid).WhereEq("active", _active)
	}
	if p.Tid > 0 {
		r.WhereOr("team", p.Tid).WhereEq("active", _active)
	}
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	rs = res.Result
	return
}

// SearchFav search app fav contest.
func (d *Dao) SearchFav(c context.Context, mid int64, p *model.ParamFav) (rs []int64, total int, err error) {
	var res struct {
		Page   model.Page      `json:"page"`
		Result []*model.ElaSub `json:"result"`
	}
	r := d.ela.NewRequest(_contestFav).Index(_contestFav).WhereEq("mid", mid).WhereEq("state", 0).Pn(p.Pn).Ps(p.Ps).Fields("oid")
	if p.Sort == 1 {
		r.Order("stime", elastic.OrderDesc)
	} else {
		r.Order("stime", elastic.OrderAsc)
	}
	if p.Stime != "" && p.Etime != "" {
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Stime+" 00:00:00", time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Etime+" 23:59:59", time.Local)
		r.WhereRange("stime", start.Unix(), end.Unix(), elastic.RangeScopeLcRc)
	} else if p.Stime != "" && p.Etime == "" {
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Stime+" 00:00:00", time.Local)
		r.WhereRange("stime", start.Unix(), "", elastic.RangeScopeLcRo)
	} else if p.Stime == "" && p.Etime != "" {
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", p.Etime+" 23:59:59", time.Local)
		r.WhereRange("stime", "", end.Unix(), elastic.RangeScopeLoRc)
	}
	if len(p.Sids) > 0 {
		r.WhereIn("sid", p.Sids)
	}
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	total = res.Page.Total
	if total == 0 {
		return
	}
	for _, contest := range res.Result {
		rs = append(rs, contest.Oid)
	}
	return
}

// SeasonFav  fav season list.
func (d *Dao) SeasonFav(c context.Context, mid int64, p *model.ParamSeason) (rs []*model.ElaSub, total int, err error) {
	var res struct {
		Page   model.Page      `json:"page"`
		Result []*model.ElaSub `json:"result"`
	}
	r := d.ela.NewRequest(_contestFav).Index(_contestFav).WhereEq("mid", mid).WhereEq("state", 0).Pn(p.Pn).Ps(p.Ps).Fields("sid", "oid")
	if p.Sort == 1 {
		r.Order("season_stime", elastic.OrderDesc)
	} else {
		r.Order("season_stime", elastic.OrderAsc)
	}
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	total = res.Page.Total
	if total == 0 {
		return
	}
	rs = res.Result
	return
}

// StimeFav  fav contest stime list.
func (d *Dao) StimeFav(c context.Context, mid int64, p *model.ParamSeason) (rs []*model.ElaSub, total int, err error) {
	var res struct {
		Page   model.Page      `json:"page"`
		Result []*model.ElaSub `json:"result"`
	}
	r := d.ela.NewRequest(_contestFav).Index(_contestFav).WhereEq("mid", mid).WhereEq("state", 0).Pn(p.Pn).Ps(p.Ps).Fields("stime", "oid")
	if p.Sort == 1 {
		r.Order("stime", elastic.OrderDesc)
	} else {
		r.Order("stime", elastic.OrderAsc)
	}
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	total = res.Page.Total
	if total == 0 {
		return
	}
	rs = res.Result
	return
}
