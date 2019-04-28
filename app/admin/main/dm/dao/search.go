package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/dm/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_dmMointorFields = []string{"id", "type", "pid", "oid", "state", "attr", "mcount", "ctime", "mtime", "mid", "title", "author"}
	_dmReportFields  = []string{"id", "dmid", "cid", "arc_aid", "arc_typeid", "dm_owner_uid", "dm_msg", "count", "content", "up_op",
		"state", "uid", "rp_time", "reason", "arc_title", "dm_deleted", "arc_mid", "pool_id", "model", "score", "dm_ctime", "ctime", "mtime"}
)

// return recent two years report search index.
func (d *Dao) rptSearchIndex() string {
	year := time.Now().Year()
	return fmt.Sprintf("dmreport_%d,dmreport_%d", year-1, year)
}

// SearchMonitor get monitor list from search
func (d *Dao) SearchMonitor(c context.Context, tp int32, pid, oid, mid int64, attr int32, kw, sort, order string, page, size int64) (data *model.SearchMonitorResult, err error) {
	req := d.esCli.NewRequest("dm_monitor_list").Index("dm_monitoring").Fields(_dmMointorFields...).Pn(int(page)).Ps(int(size))
	if tp > 0 {
		req.WhereEq("type", tp)
	}
	if pid > 0 {
		req.WhereEq("pid", pid)
	}
	if oid > 0 {
		req.WhereEq("oid", oid)
	}
	if mid > 0 {
		req.WhereEq("mid", mid)
	}
	if len(kw) > 0 {
		req.WhereLike([]string{"title"}, []string{kw}, false, elastic.LikeLevelLow)
	}
	if attr != 0 {
		req.WhereEq("attr_format", attr)
	} else {
		req.WhereIn("attr_format", []int64{int64(model.AttrSubMonitorBefore + 1), int64(model.AttrSubMonitorAfter + 1)})
	}
	if len(sort) > 0 && len(order) > 0 {
		req.Order(order, sort)
	}
	if err = req.Scan(c, &data); err != nil {
		log.Error("SearchMonitor:Scan params(%s) error(%v)", req.Params(), err)
		return
	}
	if data == nil || data.Page == nil {
		err = fmt.Errorf("data or data.Page nil")
		log.Error("SearchMonitor params(%s) error(%v)", req.Params(), err)
		return
	}
	return
}

// SearchReport2 .
func (d *Dao) SearchReport2(c context.Context, params *model.ReportListParams) (data *model.SearchReportResult, err error) {
	req := d.esCli.NewRequest("dmreport").Index(d.rptSearchIndex()).Fields(_dmReportFields...).Ps(int(params.PageSize)).Pn(int(params.Page))
	if len(params.Tids) > 0 {
		req.WhereIn("arc_typeid", params.Tids)
	}
	if len(params.RpTypes) > 0 {
		req.WhereIn("reason", params.RpTypes)
	}
	if params.Aid > 0 {
		req.WhereEq("arc_aid", params.Aid)
	}
	if params.Cid > 0 {
		req.WhereEq("cid", params.Cid)
	}
	if params.UID > 0 {
		req.WhereEq("dm_owner_uid", params.UID)
	} else {
		req.WhereNot(elastic.NotTypeEq, "dm_owner_uid").WhereEq("dm_owner_uid", 0)
	}
	if params.RpUID > 0 {
		req.WhereEq("uid", params.RpUID)
	}
	if len(params.States) > 0 {
		req.WhereIn("state", params.States)
	}
	if len(params.UpOps) > 0 {
		req.WhereIn("up_op", params.UpOps)
	}
	if params.Start != "" || params.End != "" {
		req.WhereRange("rp_time", params.Start, params.End, elastic.RangeScopeLcRc)
	}
	if params.Keyword != "" {
		req.WhereLike([]string{"dm_msg"}, []string{params.Keyword}, false, elastic.LikeLevelLow)
	}
	if len(params.Sort) > 0 && len(params.Order) > 0 {
		req.Order(params.Order, params.Sort)
	}
	if err = req.Scan(c, &data); err != nil {
		log.Error("SearchReport:Scan params(%s) error(%v)", req.Params(), err)
		return
	}
	if data == nil || data.Page == nil {
		err = fmt.Errorf("data or data.Page nil")
		log.Error("SearchReport params(%s) error(%v)", req.Params(), err)
		return
	}
	return
}

// SearchReport get report list from search
func (d *Dao) SearchReport(c context.Context, page, size int64, start, end, order, sort, keyword string, tid, rpID, state, upOp []int64, rt *model.Report) (data *model.SearchReportResult, err error) {
	req := d.esCli.NewRequest("dmreport").Index(d.rptSearchIndex()).Fields(_dmReportFields...).Ps(int(size)).Pn(int(page))
	if len(tid) > 0 {
		req.WhereIn("arc_typeid", tid)
	}
	if len(rpID) > 0 {
		req.WhereIn("reason", rpID)
	}
	if rt.Aid != -1 {
		req.WhereEq("arc_aid", rt.Aid)
	}
	if rt.Cid != -1 {
		req.WhereEq("cid", rt.Cid)
	}
	if rt.UID != -1 {
		req.WhereEq("dm_owner_uid", rt.UID)
	} else {
		req.WhereNot(elastic.NotTypeEq, "dm_owner_uid").WhereEq("dm_owner_uid", 0)
	}
	if rt.RpUID != -1 {
		req.WhereEq("uid", rt.RpUID)
	}
	if len(state) > 0 {
		req.WhereIn("state", state)
	}
	if len(upOp) > 0 {
		req.WhereIn("up_op", (upOp))
	}
	if start != "" || end != "" {
		req.WhereRange("rp_time", start, end, elastic.RangeScopeLcRc)
	}
	if keyword != "" {
		req.WhereLike([]string{"dm_msg"}, []string{keyword}, false, elastic.LikeLevelLow)
	}
	if len(sort) > 0 && len(order) > 0 {
		req.Order(order, sort)
	}
	if err = req.Scan(c, &data); err != nil {
		log.Error("SearchReport:Scan params(%s) error(%v)", req.Params(), err)
		return
	}
	if data == nil || data.Page == nil {
		err = fmt.Errorf("data or data.Page nil")
		log.Error("SearchReport params(%s) error(%v)", req.Params(), err)
		return
	}
	return
}

// SearchReportByID search report by cid and dmids.
func (d *Dao) SearchReportByID(c context.Context, dmids []int64) (data *model.SearchReportResult, err error) {
	req := d.esCli.NewRequest("dmreport").Index(d.rptSearchIndex()).Fields(_dmReportFields...).Ps(100)
	req.WhereIn("dmid", dmids)
	if err = req.Scan(c, &data); err != nil {
		log.Error("SearchReportByID:Scan params(%s) error(%v)", req.Params(), err)
		return
	}
	if data == nil || data.Page == nil {
		err = fmt.Errorf("data or data.Page nil")
		log.Error("SearchReportByID params(%s) error(%v)", req.Params(), err)
		return
	}
	return
}

// UptSearchReport 强制更新举报搜索索引 使用 v3
func (d *Dao) UptSearchReport(c context.Context, uptRpts []*model.UptSearchReport) (err error) {
	upt := d.esCli.NewUpdate("dmreport")
	var t time.Time
	for _, rpt := range uptRpts {
		t, err = time.ParseInLocation("2006-01-02 15:04:05", rpt.Ctime, time.Local)
		if err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", rpt.Ctime, err)
			return
		}
		upt.AddData(fmt.Sprintf("dmreport_%d", t.Year()), rpt)
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() error(%v)", err)
	}
	return
}

// SearchDM 搜索弹幕
func (d *Dao) SearchDM(c context.Context, p *model.SearchDMParams) (data *model.SearchDMData, err error) {
	var (
		order = "ctime"
		sort  = "desc"
		req   *elastic.Request
	)
	req = d.esCli.NewRequest("dm").Index(fmt.Sprintf("dm_%03d", p.Oid%_indexSharding)).Fields("id").Ps(int(p.Size)).Pn(int(p.Page))
	req.WhereEq("oid", p.Oid)
	if p.Mid != model.CondIntNil {
		req.WhereEq("mid", p.Mid)
	}
	if p.State != "" {
		if states, err1 := xstr.SplitInts(p.State); err1 == nil {
			req.WhereIn("state", states)
		}
	}
	if p.Pool != "" {
		if pools, err1 := xstr.SplitInts(p.Pool); err1 == nil {
			req.WhereIn("pool", pools)
		}
	}
	if p.Attrs != "" {
		if attrs, err1 := xstr.SplitInts(p.Attrs); err1 == nil {
			req.WhereIn("attr_format", attrs)
		}
	}
	if p.IP != "" {
		req.WhereEq("ip_format", p.IP)
	}

	switch {
	case p.ProgressFrom != model.CondIntNil && p.ProgressTo != model.CondIntNil:
		req.WhereRange("progress", p.ProgressFrom, p.ProgressTo, elastic.RangeScopeLcRc)
	case p.ProgressFrom != model.CondIntNil:
		req.WhereRange("progress", p.ProgressFrom, nil, elastic.RangeScopeLcRc)
	case p.ProgressTo != model.CondIntNil:
		req.WhereRange("progress", nil, p.ProgressTo, elastic.RangeScopeLcRc)
	}

	switch {
	case p.CtimeFrom != model.CondIntNil && p.CtimeTo != model.CondIntNil:
		req.WhereRange("ctime", time.Unix(p.CtimeFrom, 0).Format("2006-01-02 15:04:05"), time.Unix(p.CtimeTo, 0).Format("2006-01-02 15:04:05"), elastic.RangeScopeLcRc)
	case p.CtimeFrom != model.CondIntNil:
		req.WhereRange("ctime", time.Unix(p.CtimeFrom, 0).Format("2006-01-02 15:04:05"), nil, elastic.RangeScopeLcRc)
	case p.CtimeTo != model.CondIntNil:
		req.WhereRange("ctime", nil, time.Unix(p.CtimeTo, 0).Format("2006-01-02 15:04:05"), elastic.RangeScopeLcRc)
	}

	if p.Keyword != "" {
		req.WhereLike([]string{"kwmsg"}, []string{p.Keyword}, false, elastic.LikeLevelHigh)
	}
	if p.Order != "" {
		order = p.Order
	}
	if p.Sort == "asc" {
		sort = p.Sort
	}
	req.Order(order, sort)
	if err = req.Scan(c, &data); err != nil {
		log.Error("SearchDM:Scan params(%s) error(%v)", req.Params(), err)
		return
	}
	if data == nil || data.Page == nil {
		err = fmt.Errorf("data or data.Page nil")
		log.Error("SearchDM params(%s) error(%v)", req.Params(), err)
		return
	}
	return
}

// SearchProtectCount get protected dm count.
func (d *Dao) SearchProtectCount(c context.Context, tp int32, oid int64) (count int64, err error) {
	var res struct {
		Result map[string][]struct {
			Key   string `json:"key"`
			Count int64  `json:"doc_count"`
		} `json:"result"`
	}
	order := []map[string]string{{"attr": "desc"}}
	req := d.esCli.NewRequest("dm").Fields("attr").Index(fmt.Sprintf("dm_%03d", d.hitIndex(oid)))
	req.WhereEq("oid", oid).WhereIn("attr_format", "1").GroupBy(elastic.EnhancedModeGroupBy, "attr", order)
	req.Pn(1).Ps(10)
	if err = req.Scan(c, &res); err != nil {
		log.Error("req.Scan() error(%v)", err)
		return
	}
	if values, ok := res.Result["group_by_attr"]; ok {
		for _, v := range values {
			count = count + v.Count
		}
	}
	return
}

// UpSearchDMState 通知搜索服务更新弹幕状态
func (d *Dao) UpSearchDMState(c context.Context, tp int32, state int32, dmidM map[int64][]int64) (err error) {
	upt := d.esCli.NewUpdate("dm")
	for oid, dmids := range dmidM {
		for _, dmid := range dmids {
			data := &model.UptSearchDMState{
				ID:    dmid,
				Oid:   oid,
				State: state,
				Type:  tp,
				Mtime: time.Now().Format("2006-01-02 15:04:05"),
			}
			upt.AddData(fmt.Sprintf("dm_%03d", oid%_indexSharding), data)
		}
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// UpSearchDMPool  通知搜索服务更新弹幕池
func (d *Dao) UpSearchDMPool(c context.Context, tp int32, oid int64, pool int32, dmids []int64) (err error) {
	upt := d.esCli.NewUpdate("dm")
	for _, dmid := range dmids {
		data := &model.UptSearchDMPool{
			ID:    dmid,
			Oid:   oid,
			Pool:  pool,
			Type:  tp,
			Mtime: time.Now().Format("2006-01-02 15:04:05"),
		}
		upt.AddData(fmt.Sprintf("dm_%03d", oid%_indexSharding), data)
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// UpSearchDMAttr  通知搜索服务更新弹幕属性
func (d *Dao) UpSearchDMAttr(c context.Context, tp int32, oid int64, attr int32, dmids []int64) (err error) {
	var bits []int64
	for k, v := range strconv.FormatInt(int64(attr), 2) {
		if v == 49 {
			bits = append(bits, int64(k+1))
		}
	}
	upt := d.esCli.NewUpdate("dm")
	for _, dmid := range dmids {
		data := &model.UptSearchDMAttr{
			ID:         dmid,
			Oid:        oid,
			Attr:       attr,
			AttrFormat: bits,
			Type:       tp,
			Mtime:      time.Now().Format("2006-01-02 15:04:05"),
		}
		upt.AddData(fmt.Sprintf("dm_%03d", oid%_indexSharding), data)
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// SearchSubjectLog get subject log
func (d *Dao) SearchSubjectLog(c context.Context, tp int32, oid int64) (data []*model.SubjectLog, err error) {
	req := d.esCli.NewRequest("log_audit").Index("log_audit_31_all").Fields("uid", "uname", "oid", "ctime", "action", "extra_data").WhereEq("oid", oid).WhereEq("type", tp)
	req.Ps(20).Order("ctime", "desc")
	res := &model.SearchSubjectLog{}
	if err = req.Scan(c, &res); err != nil || res == nil {
		log.Error("SearchSubcetLog:Scan params(%s) error(%v)", req.Params(), err)
		return
	}
	data = make([]*model.SubjectLog, 0)
	s := new(struct {
		Comment string `json:"comment"`
	})
	for _, v := range res.Result {
		if err = json.Unmarshal([]byte(v.ExtraData), s); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", v.ExtraData, err)
			return
		}
		log := &model.SubjectLog{
			UID:     v.UID,
			Uname:   v.Uname,
			Oid:     v.Oid,
			Action:  v.Action,
			Comment: s.Comment,
			Ctime:   v.Ctime,
		}
		data = append(data, log)
	}
	return
}

// UpSearchRecentDMState .
func (d *Dao) UpSearchRecentDMState(c context.Context, tp int32, state int32, dmidM map[int64][]int64) (err error) {
	upt := d.esCli.NewUpdate("dm_home")
	year, month, _ := time.Now().Date()
	yearPre, monthPre, _ := time.Now().AddDate(0, -1, 0).Date()
	for oid, dmids := range dmidM {
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
	}
	if err = upt.Do(c); err != nil {
		log.Error("update.Do() params(%s) error(%v)", upt.Params(), err)
	}
	return
}

// UpSearchRecentDMPool .
func (d *Dao) UpSearchRecentDMPool(c context.Context, tp int32, oid int64, pool int32, dmids []int64) (err error) {
	upt := d.esCli.NewUpdate("dm_home")
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

// UpSearchRecentDMAttr .
func (d *Dao) UpSearchRecentDMAttr(c context.Context, tp int32, oid int64, attr int32, dmids []int64) (err error) {
	var bits []int64
	for k, v := range strconv.FormatInt(int64(attr), 2) {
		if v == 49 {
			bits = append(bits, int64(k+1))
		}
	}
	upt := d.esCli.NewUpdate("dm_home")
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

// SearchSubject get subject log list from search
func (d *Dao) SearchSubject(c context.Context, req *model.SearchSubjectReq) (data []int64, page *model.Page, err error) {
	r := d.esCli.NewRequest("dm_monitor_list").Index("dm_monitoring").Fields("oid")
	if len(req.Oids) > 0 {
		r.WhereIn("oid", req.Oids)
	}
	if len(req.Mids) > 0 {
		r.WhereIn("mid", req.Mids)
	}
	if len(req.Aids) > 0 {
		r.WhereIn("pid", req.Aids)
	}
	if len(req.Attrs) > 0 {
		r.WhereIn("attr_format", req.Attrs)
	}
	if req.State != model.CondIntNil {
		r.WhereEq("state", req.State)
	}
	r.Ps(int(req.Ps)).Pn(int(req.Pn)).Order(req.Order, req.Sort)
	res := &model.SearchSubjectResult{}
	if err = r.Scan(c, &res); err != nil {
		log.Error("SearchSubject:Scan params(%s) error(%v)", r.Params(), err)
		return
	}
	if res == nil || res.Page == nil {
		err = fmt.Errorf("data or data.Page nil")
		log.Error("SearchSubject params(%s) error(%v)", r.Params(), err)
		return
	}
	for _, v := range res.Result {
		data = append(data, v.Oid)
	}
	page = res.Page
	return
}
