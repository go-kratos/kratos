package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/admin/main/workflow/model/manager"
	"go-common/app/admin/main/workflow/model/search"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_unameURI       = "http://manager.bilibili.co/x/admin/manager/users/unames"
	_srhAuditLogURI = "http://bili-search.bilibili.co/x/admin/search/log"
)

// SearchGroup will search group by given conditions
func (d *Dao) SearchGroup(c context.Context, cond *search.GroupSearchCommonCond) (resp *search.GroupSearchCommonResp, err error) {
	start := time.Now()
	var r *elastic.Request
	defer func() {
		log.Info("SearchGroup params(%s) group search ts %s err_or(%v)", r.Params(), time.Since(start).String(), err)
	}()

	r = d.es.NewRequest(search.GroupSrhComID).Index(search.GroupSrhComID).Fields(cond.Fields...).
		WhereEq("business", cond.Business).WhereIn("round", cond.Rounds).WhereIn("tid", cond.Tids).
		WhereIn("state", cond.States).WhereIn("mid", cond.Mids).WhereIn("oid", cond.Oids).WhereIn("typeid", cond.TypeIDs).
		WhereIn("fid", cond.FID).WhereIn("rid", cond.RID).WhereIn("eid", cond.EID).
		WhereIn("report_mid", cond.ReportMID).WhereIn("first_user_tid", cond.FirstUserTid).
		Order(cond.Order, cond.Sort).
		Pn(int(cond.PN)).Ps(int(cond.PS))

	// 是否关键字匹配优先
	if cond.KWPriority == true {
		r.OrderScoreFirst(true)
	} else {
		r.OrderScoreFirst(false)
	}
	if len(cond.KWFields) > 0 && len(cond.KWFields) == len(cond.KW) {
		r.WhereLike(cond.KWFields, cond.KW, true, elastic.LikeLevelMiddle)
	}
	r.WhereRange("ctime", cond.CTimeFrom, cond.CTimeTo, elastic.RangeScopeLcRc)

	if err = r.Scan(c, &resp); err != nil {
		log.Error("r.Scan(%+v) error(%v) params(%s)", resp, err, r.Params())
	}
	return
}

// SearchGroupMultiPage .
func (d *Dao) SearchGroupMultiPage(c context.Context, cond *search.GroupSearchCommonCond) (result []*search.GroupSearchCommonData, err error) {
	var resp *search.GroupSearchCommonResp
	cond.PS = 1000
	cond.PN = 1
	result = make([]*search.GroupSearchCommonData, 0, len(cond.IDs))
	for {
		if resp, err = d.SearchGroup(c, cond); err != nil {
			return
		}
		result = append(result, resp.Result...)
		if len(resp.Result) < resp.Page.Size {
			break
		}
		cond.PN++
	}
	return
}

// SearchChallenge will search challenge by given conditions
func (d *Dao) SearchChallenge(c context.Context, cond *search.ChallSearchCommonCond) (resp *search.ChallSearchCommonResp, err error) {
	start := time.Now()
	var r *elastic.Request
	defer func() {
		log.Info("SearchChallenge params(%s) challenge search ts %s err_or(%v)", r.Params(), time.Since(start).String(), err)
	}()

	r = d.es.NewRequest(search.ChallSrhComID).Fields(cond.Fields...).
		WhereIn("id", cond.IDs).WhereIn("round", cond.Rounds).WhereIn("tid", cond.Tids).WhereIn("state", cond.States).
		WhereIn("business_state", cond.BusinessStates).WhereIn("mid", cond.Mids).WhereIn("oid", cond.Oids).WhereIn("typeid", cond.TypeIDs).
		WhereIn("gid", cond.Gids).WhereIn("assignee_adminid", cond.AssigneeAdminIDs).WhereIn("adminid", cond.AdminIDs)
	if cond.Business > 0 {
		r.WhereEq("business", cond.Business)
	}
	if len(cond.KWFields) > 0 && len(cond.KWFields) == len(cond.KW) {
		r.WhereLike(cond.KWFields, cond.KW, true, elastic.LikeLevelLow)
	}
	if cond.Order == "" {
		cond.Order = "id"
	}
	if cond.Sort == "" {
		cond.Sort = "desc"
	}
	r.Order(cond.Order, cond.Sort)
	r.WhereRange("ctime", cond.CTimeFrom, cond.CTimeTo, elastic.RangeScopeLcRc)
	if len(cond.Distinct) > 0 {
		for _, g := range cond.Distinct {
			r.GroupBy("distinct", g, nil)
		}
	}
	r.Index(search.ChallSrhComID)
	r.Pn(1)
	r.Ps(50)
	if cond.PN != 0 {
		r.Pn(int(cond.PN))
	}
	if cond.PS != 0 {
		r.Ps(int(cond.PS))
	}
	if err = r.Scan(c, &resp); err != nil {
		log.Error("r.Scan(%+v) error(%v) params(%s)", resp, err, r.Params())
	}
	return
}

// SearchChallengeMultiPage .
func (d *Dao) SearchChallengeMultiPage(c context.Context, cond *search.ChallSearchCommonCond) (result []*search.ChallSearchCommonData, err error) {
	var resp *search.ChallSearchCommonResp
	cond.PS = 1000
	cond.PN = 1
	result = make([]*search.ChallSearchCommonData, 0, len(cond.IDs))
	for {
		if resp, err = d.SearchChallenge(c, cond); err != nil {
			return
		}
		result = append(result, resp.Result...)
		if len(resp.Result) < resp.Page.Size {
			break
		}
		cond.PN++
		// return if result too long
		if cond.PN > 10 {
			log.Warn("cond(%+v) result is too long to degrade", cond)
			return
		}
	}
	return
}

// BatchUNameByUID will search unames by uids
func (d *Dao) BatchUNameByUID(c context.Context, uids []int64) (UNames map[int64]string, err error) {
	//todo: local cache uname
	uri := _unameURI
	uv := url.Values{}
	UNames = make(map[int64]string)
	if len(uids) == 0 {
		return
	}
	uv.Set("uids", xstr.JoinInts(uids))
	unameSchRes := new(manager.UNameSearchResult)
	if err = d.httpRead.Get(c, uri, "", uv, unameSchRes); err != nil {
		return
	}
	if unameSchRes.Code != ecode.OK.Code() {
		log.Error("search uname failed: %s?%s, error code(%d)", uri, uv.Get("uids"), unameSchRes.Code)
		err = ecode.Int(unameSchRes.Code)
		return
	}
	UNames = unameSchRes.Data
	return
}

// SearchAuditLogGroup search archive audit log from log platform
func (d *Dao) SearchAuditLogGroup(c context.Context, cond *search.AuditLogGroupSearchCond) (auditLogSchRes *search.AuditLogSearchResult, err error) {
	uri := _srhAuditLogURI
	uv := cond.Query()
	auditLogSchRes = new(search.AuditLogSearchResult)
	if err = d.httpRead.Get(c, uri, "", uv, auditLogSchRes); err != nil {
		log.Error("call search audit log %s error(%v)", uri, err)
		return
	}
	if auditLogSchRes.Code != ecode.OK.Code() {
		log.Error("call search audit log %s result error code(%d), message(%s)", uri, auditLogSchRes.Code, auditLogSchRes.Message)
		err = ecode.Int(auditLogSchRes.Code)
	}
	return
}

// SearchAuditReportLog .
func (d *Dao) SearchAuditReportLog(c context.Context, cond *search.AuditReportSearchCond) (resp *search.AuditLogSearchCommonResult, err error) {
	if len(cond.Fields) == 0 {
		return
	}
	r := d.es.NewRequest(search.LogAuditAction).Fields(cond.Fields...).WhereIn("uid", cond.UID).WhereIn("oid", cond.Oid).
		WhereEq("business", cond.Business).WhereIn("type", cond.Type).Order(cond.Order, cond.Sort).Pn(1).Ps(1000)

	indexPrefix := search.LogAuditAction + "_" + strconv.Itoa(cond.Business)
	if cond.IndexTimeType != "" {
		switch cond.IndexTimeType {
		case "year":
			r.IndexByTime(indexPrefix, elastic.IndexTypeYear, cond.IndexTimeFrom, cond.IndexTimeEnd)
		case "month":
			r.IndexByTime(indexPrefix, elastic.IndexTypeMonth, cond.IndexTimeFrom, cond.IndexTimeEnd)
		case "week":
			r.IndexByTime(indexPrefix, elastic.IndexTypeWeek, cond.IndexTimeFrom, cond.IndexTimeEnd)
		case "day":
			r.IndexByTime(indexPrefix, elastic.IndexTypeDay, cond.IndexTimeFrom, cond.IndexTimeEnd)
		default:
			r.Index(indexPrefix + "_all")
		}
	} else {
		r.Index(indexPrefix + "_all")
	}

	r.WhereIn("int_0", cond.Int0).WhereIn("int_1", cond.Int1).WhereIn("int_2", cond.Int2)
	if cond.Str0 != "" {
		r.WhereEq("str_0", cond.Str0)
	}
	if cond.Str1 != "" {
		r.WhereEq("str_1", cond.Str1)
	}
	if cond.Str2 != "" {
		r.WhereEq("str_2", cond.Str2)
	}

	if cond.Group != "" {
		r.GroupBy(elastic.EnhancedModeGroupBy, cond.Group, []map[string]string{{"ctime": "desc"}})
	}

	if cond.Distinct != "" {
		r.GroupBy(elastic.EnhancedModeDistinct, cond.Distinct, []map[string]string{{"ctime": "desc"}})
	}

	if err = r.Scan(c, &resp); err != nil {
		log.Error("r.Scan(%+v) error(%v) params(%s)", resp, err, r.Params())
	}
	log.Info("SearchAuditReportLog end param(%v) err(%v)", r.Params(), err)
	return
}
