package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_dataTimeFormat = "2006-01-02 15:03:04"

	_subtitleFields = []string{"oid", "id"}

	_dmRecentFields = []string{"attr", "color", "ctime", "fontsize", "id", "mid", "mode", "msg", "mtime", "oid", "pool", "progress", "state", "type", "pid"}
)

func hisDateIndex(month string) string {
	return "dm_date_" + strings.Replace(month, "-", "_", -1)
}

// SearchDMHisIndex get dm date index by oid from search.
func (d *Dao) SearchDMHisIndex(c context.Context, tp int32, oid int64, month string) (dates []string, err error) {
	var (
		pn, ps = 1, 31
		res    model.SearchHistoryIdxResult
	)
	req := d.elastic.NewRequest("dm_date")
	req.Fields("date").WhereEq("oid", oid)
	req.Index(hisDateIndex(month)).Pn(pn).Ps(ps).Order("date", "asc")
	if err = req.Scan(c, &res); err != nil {
		log.Error("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	for _, v := range res.Result {
		dates = append(dates, v.Date)
	}
	return
}

// SearchDMHistory get history dmid from search.
// 搜索定制api，改动需要沟通
func (d *Dao) SearchDMHistory(c context.Context, tp int32, oid, ctimeTo int64, pn, ps int) (dmids []int64, err error) {
	var (
		res model.SearchHistoryResult
		end = time.Unix(ctimeTo, 0).Format("2006-01-02 15:04:05")
	)
	req := d.elastic.NewRequest("dm_history")
	req.Index(fmt.Sprintf("dm_search_%03d", oid%_indexSharding))
	req.Fields("id").WhereEq("oid", oid).WhereIn("state", []int64{0, 2, 6})
	req.WhereRange("ctime", nil, end, elastic.RangeScopeLcRc)
	req.Pn(pn).Ps(ps).Order("ctime", "desc")
	if err = req.Scan(c, &res); err != nil {
		log.Error("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	for _, v := range res.Result {
		dmids = append(dmids, v.ID)
	}
	return
}

// SearchDM 搜索弹幕
func (d *Dao) SearchDM(c context.Context, p *model.SearchDMParams) (res *model.SearchDMData, err error) {
	req := d.elastic.NewRequest("dm_search")
	req.Fields("id").Index(fmt.Sprintf("dm_search_%03d", p.Oid%_indexSharding)).WhereEq("oidstr", p.Oid)
	if p.Mids != "" {
		mids, _ := xstr.SplitInts(p.Mids)
		req.WhereIn("mid", mids)
	}
	if p.State != "" {
		states, _ := xstr.SplitInts(p.State)
		req.WhereIn("state", states)
	}
	if p.Mode != "" {
		modes, _ := xstr.SplitInts(p.Mode)
		req.WhereIn("mode", modes)
	}
	if p.Pool != "" {
		pools, _ := xstr.SplitInts(p.Pool)
		req.WhereIn("pool", pools)
	}
	if p.Attrs != "" {
		attrs, _ := xstr.SplitInts(p.Attrs)
		req.WhereIn("attr_format", attrs)
	}
	req.WhereEq("type", p.Type)
	switch {
	case p.ProgressFrom != model.CondIntNil && p.ProgressTo != model.CondIntNil:
		req.WhereRange("progress_long", p.ProgressFrom, p.ProgressTo, elastic.RangeScopeLcRc)
	case p.ProgressFrom != model.CondIntNil:
		req.WhereRange("progress_long", p.ProgressFrom, nil, elastic.RangeScopeLcRc)
	case p.ProgressTo != model.CondIntNil:
		req.WhereRange("progress_long", nil, p.ProgressTo, elastic.RangeScopeLcRc)
	}
	req.WhereRange("ctime", p.CtimeFrom, p.CtimeTo, elastic.RangeScopeLcRc)
	if p.Keyword != "" {
		req.WhereLike([]string{"msg"}, []string{p.Keyword}, true, elastic.LikeLevelHigh)
		req.OrderScoreFirst(true)
	}
	if p.Order == "progress" {
		p.Order = "progress_long"
	}
	req.Order(p.Order, p.Sort)
	req.Pn(int(p.Pn)).Ps(int(p.Ps))
	res = &model.SearchDMData{}
	if err = req.Scan(c, &res); err != nil {
		log.Error("search params(%s), error(%v)", req.Params(), err)
	}
	return
}

// UptSearchDMState update dm search state
func (d *Dao) UptSearchDMState(c context.Context, dmids []int64, oid int64, state, tp int32) (err error) {
	upt := d.elastic.NewUpdate("dm_search")
	for _, dmid := range dmids {
		data := &model.UptSearchDMState{
			ID:    dmid,
			Oid:   oid,
			State: state,
			Type:  tp,
			Mtime: time.Now().Format("2006-01-02 15:04:05"),
		}
		upt.AddData(fmt.Sprintf("dm_search_%03d", oid%_indexSharding), data)
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// UptSearchDMPool update dm search pool
func (d *Dao) UptSearchDMPool(c context.Context, dmids []int64, oid int64, pool, tp int32) (err error) {
	upt := d.elastic.NewUpdate("dm_search")
	for _, dmid := range dmids {
		data := &model.UptSearchDMPool{
			ID:    dmid,
			Oid:   oid,
			Pool:  pool,
			Type:  tp,
			Mtime: time.Now().Format("2006-01-02 15:04:05"),
		}
		upt.AddData(fmt.Sprintf("dm_search_%03d", oid%_indexSharding), data)
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// UptSearchDMAttr update dm search attr
func (d *Dao) UptSearchDMAttr(c context.Context, dmids []int64, oid int64, attr, tp int32) (err error) {
	var bits []int64
	for k, v := range strconv.FormatInt(int64(attr), 2) {
		if v == 49 {
			bits = append(bits, int64(k+1))
		}
	}
	upt := d.elastic.NewUpdate("dm_search")
	for _, dmid := range dmids {
		data := &model.UptSearchDMAttr{
			ID:         dmid,
			Oid:        oid,
			Attr:       attr,
			AttrFormat: bits,
			Type:       tp,
			Mtime:      time.Now().Format("2006-01-02 15:04:05"),
		}
		upt.AddData(fmt.Sprintf("dm_search_%03d", oid%_indexSharding), data)
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// SearchSubtitles .
func (d *Dao) SearchSubtitles(c context.Context, page, size int32, mid int64, upMids []int64, aid, oid int64, tp int32, status []int64) (res *model.SearchSubtitleResult, err error) {
	var (
		req    *elastic.Request
		fields []string
	)
	fields = _subtitleFields
	req = d.elastic.NewRequest("dm_subtitle").Index("subtitle").Fields(fields...).Pn(int(page)).Ps(int(size))
	if mid > 0 {
		req.WhereEq("mid", mid)
	}
	if aid > 0 {
		req.WhereEq("aid", aid)
	}
	if oid > 0 {
		req.WhereEq("oid", oid)
		req.WhereEq("type", tp)
	}

	switch {
	case len(upMids) > 0 && len(status) > 0:
		cmbs := &elastic.Combo{}
		cmbu := &elastic.Combo{}
		var (
			statusInf []interface{}
			upMidsInf []interface{}
		)
		for _, s := range status {
			statusInf = append(statusInf, s)
		}
		for _, s := range upMids {
			upMidsInf = append(upMidsInf, s)
		}
		cmbs.ComboIn([]map[string][]interface{}{
			{"status": statusInf},
		})
		cmbu.ComboIn([]map[string][]interface{}{
			{"up_mid": upMidsInf},
		})
		req = req.WhereCombo(cmbs.MinIn(1).MinAll(1), cmbu.MinIn(1).MinAll(1))
	case len(upMids) > 0:
		req.WhereIn("up_mid", upMids)
	case len(status) > 0:
		req.WhereIn("status", status)
	}
	req.Order("mtime", "desc")
	if err = req.Scan(c, &res); err != nil {
		log.Error("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	return
}

// CountSubtitles .
func (d *Dao) CountSubtitles(c context.Context, mid int64, upMids []int64, aid, oid int64, tp int32) (countSubtitle *model.CountSubtitleResult, err error) {
	var (
		req           *elastic.Request
		fields        []string
		res           map[string]interface{}
		_searchStatus = []string{
			fmt.Sprint(model.SubtitleStatusDraft),
			fmt.Sprint(model.SubtitleStatusToAudit),
			fmt.Sprint(model.SubtitleStatusAuditBack),
			fmt.Sprint(model.SubtitleStatusPublish),
			fmt.Sprint(model.SubtitleStatusCheckToAudit),
			fmt.Sprint(model.SubtitleStatusCheckPublish),
			fmt.Sprint(model.SubtitleStatusManagerBack),
		}
		result      map[string]interface{}
		groupStatus []interface{}
		itemStatus  map[string]interface{}
		ok          bool
	)
	req = d.elastic.NewRequest("dm_subtitle").Index("subtitle").Fields(fields...).Pn(0).Ps(0)
	if mid > 0 {
		req.WhereEq("mid", mid)
	}
	if aid > 0 {
		req.WhereEq("aid", aid)
	}
	if oid > 0 {
		req.WhereEq("oid", oid)
		req.WhereEq("type", tp)
	}
	if len(upMids) > 0 {
		req.WhereIn("up_mid", upMids)
	}
	req = req.GroupBy(elastic.EnhancedModeGroupBy, "status", nil).WhereIn("status", _searchStatus)
	if err = req.Scan(c, &res); err != nil {
		log.Error("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	countSubtitle = &model.CountSubtitleResult{}
	if result, ok = res["result"].(map[string]interface{}); !ok {
		return
	}
	if groupStatus, ok = result["group_by_status"].([]interface{}); !ok {
		return
	}
	for _, item := range groupStatus {
		if itemStatus, ok = item.(map[string]interface{}); ok {
			docCount, _ := itemStatus["doc_count"].(float64)
			switch itemStatus["key"] {
			case fmt.Sprint(model.SubtitleStatusDraft):
				countSubtitle.Draft += int64(docCount)
			case fmt.Sprint(model.SubtitleStatusToAudit):
				countSubtitle.ToAudit += int64(docCount)
			case fmt.Sprint(model.SubtitleStatusAuditBack):
				countSubtitle.AuditBack += int64(docCount)
			case fmt.Sprint(model.SubtitleStatusPublish):
				countSubtitle.Publish += int64(docCount)
			case fmt.Sprint(model.SubtitleStatusCheckToAudit):
				countSubtitle.ToAudit += int64(docCount)
			case fmt.Sprint(model.SubtitleStatusCheckPublish):
				countSubtitle.Publish += int64(docCount)
			case fmt.Sprint(model.SubtitleStatusManagerBack):
				countSubtitle.AuditBack += int64(docCount)
			}
		}
	}
	return
}

// UptSearchRecentState .
func (d *Dao) UptSearchRecentState(c context.Context, dmids []int64, oid int64, state, tp int32) (err error) {
	upt := d.elastic.NewUpdate("dm_home")
	year, month, _ := time.Now().Date()
	yearPre, monthPre, _ := time.Now().AddDate(0, -1, 0).Date()
	for _, dmid := range dmids {
		data := &model.UptSearchDMState{
			ID:    dmid,
			Oid:   oid,
			State: state,
			Type:  tp,
			Mtime: time.Now().Format("2006-01-02 15:04:05"),
		}
		upt.AddData(fmt.Sprintf("dm_home_%v",
			fmt.Sprintf("%d_%02d", yearPre, int(monthPre)),
		), data)
		upt.AddData(fmt.Sprintf("dm_home_%v",
			fmt.Sprintf("%d_%02d", year, int(month)),
		), data)
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// UptSearchRecentAttr .
func (d *Dao) UptSearchRecentAttr(c context.Context, dmids []int64, oid int64, attr, tp int32) (err error) {
	var bits []int64
	for k, v := range strconv.FormatInt(int64(attr), 2) {
		if v == 49 {
			bits = append(bits, int64(k+1))
		}
	}
	upt := d.elastic.NewUpdate("dm_home")
	year, month, _ := time.Now().Date()
	yearPre, monthPre, _ := time.Now().AddDate(0, -1, 0).Date()
	for _, dmid := range dmids {
		data := &model.UptSearchDMAttr{
			ID:         dmid,
			Oid:        oid,
			Attr:       attr,
			AttrFormat: bits,
			Type:       tp,
			Mtime:      time.Now().Format("2006-01-02 15:04:05"),
		}
		upt.AddData(fmt.Sprintf("dm_home_%v",
			fmt.Sprintf("%d_%02d", yearPre, int(monthPre)),
		), data)
		upt.AddData(fmt.Sprintf("dm_home_%v",
			fmt.Sprintf("%d_%02d", year, int(month)),
		), data)
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// UptSearchRecentPool .
func (d *Dao) UptSearchRecentPool(c context.Context, dmids []int64, oid int64, pool, tp int32) (err error) {
	upt := d.elastic.NewUpdate("dm_home")
	year, month, _ := time.Now().Date()
	yearPre, monthPre, _ := time.Now().AddDate(0, -1, 0).Date()
	for _, dmid := range dmids {
		data := &model.UptSearchDMPool{
			ID:    dmid,
			Oid:   oid,
			Pool:  pool,
			Type:  tp,
			Mtime: time.Now().Format("2006-01-02 15:04:05"),
		}
		upt.AddData(fmt.Sprintf("dm_home_%v",
			fmt.Sprintf("%d_%02d", yearPre, int(monthPre)),
		), data)
		upt.AddData(fmt.Sprintf("dm_home_%v",
			fmt.Sprintf("%d_%02d", year, int(month)),
		), data)
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// SearhcDmRecent .
func (d *Dao) SearhcDmRecent(c context.Context, param *model.SearchRecentDMParam) (res *model.SearchRecentDMResult, err error) {
	var (
		req *elastic.Request
	)
	year, month, _ := time.Now().Date()
	yearPre, monthPre, _ := time.Now().AddDate(0, -1, 0).Date()
	req = d.elastic.NewRequest("dm_home").Index(fmt.Sprintf("dm_home_%v,dm_home_%v",
		fmt.Sprintf("%d_%02d", year, int(month)),
		fmt.Sprintf("%d_%02d", yearPre, int(monthPre)),
	)).Fields(_dmRecentFields...).Pn(param.Pn).Ps(param.Ps)
	if param.Type > 0 {
		req.WhereEq("type", param.Type)
	}
	if param.UpMid > 0 {
		req.WhereEq("o_mid", param.UpMid)
	}
	if len(param.States) > 0 {
		req.WhereIn("state", param.States)
	}
	req.WhereRange("ctime", time.Now().Local().AddDate(0, 0, -30).Format(_dataTimeFormat), time.Now().Local().Format(_dataTimeFormat), elastic.RangeScopeLcRo)
	req.Order(param.Field, param.Sort)
	res = &model.SearchRecentDMResult{}
	if err = req.Scan(c, &res); err != nil {
		log.Error("search params(%s), error(%v)", req.Params(), err)
	}
	return
}
