package service

import (
	"context"
	"go-common/library/ecode"
	"go-common/library/xstr"
	"net/url"
	"strconv"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

const (
	_searchBusinessQAVideo       = "task_qa"
	_searchBusinessQAVideoRandom = "task_qa_random"
	_searchIndexQAVideo          = "task_qa"
	_searchLogURL                = "/x/admin/search/log"
)

func (s *Service) searchQAVideo(c context.Context, pm *model.ListParams) (list *model.QAVideoList, err error) {
	needRandom := pm.Limit > 0 && pm.Seed != ""
	business := _searchBusinessQAVideo
	if needRandom {
		business = _searchBusinessQAVideoRandom
	}

	req := s.es.NewRequest(business).Index(_searchIndexQAVideo).Ps(pm.Ps).Pn(pm.Pn)
	if pm.CTimeFrom != "" || pm.CTimeTo != "" {
		req.WhereRange("ctime", pm.CTimeFrom, pm.CTimeTo, elastic.RangeScopeLcRc)
	}
	if pm.FTimeFrom != "" || pm.FTimeTo != "" {
		req.WhereRange("ftime", pm.FTimeFrom, pm.FTimeTo, elastic.RangeScopeLcRc)
	}
	if pm.FansFrom > 0 || pm.FansTo > 0 {
		req.WhereRange("fans", pm.FansFrom, pm.FansTo, elastic.RangeScopeLcRc)
	}
	if len(pm.UID) > 0 {
		req.WhereIn("uid", pm.UID)
	}
	if len(pm.TaskID) > 0 {
		req.WhereIn("task_id", pm.TaskID)
	}
	if len(pm.TagID) > 0 {
		req.WhereIn("audit_tagid", pm.TagID)
	}
	if len(pm.UPGroup) > 0 {
		req.WhereIn("up_groups", pm.UPGroup)
	}
	if len(pm.ArcTypeID) > 0 {
		req.WhereIn("arc_typeid", pm.ArcTypeID)
	}
	if len(pm.AuditStatus) > 0 {
		req.WhereIn("audit_status", pm.AuditStatus)
	}
	if len(pm.Keyword) > 0 {
		req.WhereLike([]string{"arc_title"}, pm.Keyword, true, elastic.LikeLevelLow)
	}
	if needRandom {
		req.WhereEq("seed", pm.Seed)
	} else {
		req.Order(pm.Order, pm.Sort)
	}
	if pm.State == model.QAStateWait || pm.State == model.QAStateFinish {
		req.WhereEq("state", pm.State)
	}

	if err = req.Scan(c, &list); err != nil {
		log.Error("searchQAVideo elastic scan error(%v) params(%+v)", err, pm)
		return
	}
	if needRandom && list != nil && list.Page.Total > pm.Limit {
		list.Page.Total = pm.Limit
		//移除多余部分
		addition := list.Page.Num*list.Page.Size - pm.Limit
		if addition > 0 {
			list.Result = list.Result[:(list.Page.Size - addition)]
		}
	}
	return
}

func (s *Service) lastInTime(c context.Context, ids []int64) (mcases map[int64][]interface{}, err error) {
	return s.lastTime(c, model.ActionHandsUP, ids)
}

func (s *Service) lastOutTime(c context.Context, ids []int64) (mcases map[int64][]interface{}, err error) {
	return s.lastTime(c, model.ActionHandsOFF, ids)
}

// lastInOutTime
func (s *Service) lastTime(c context.Context, action int8, ids []int64) (mcases map[int64][]interface{}, err error) {
	mcases = make(map[int64][]interface{})
	params := url.Values{}
	uri := s.c.Host.Search + _searchLogURL
	params.Set("appid", "log_audit_group")
	params.Set("group", "uid")
	params.Set("uid", xstr.JoinInts(ids))
	params.Set("business", strconv.Itoa(model.LogClientConsumer))
	params.Set("action", strconv.Itoa(int(action)))
	params.Set("ps", strconv.Itoa(len(ids)))
	res := &model.SearchLogResult{}
	if err = s.httpClient.Get(c, uri, "", params, &res); err != nil {
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
	log.Info("log_audit_group get: %s params:%s ret:%v", uri, params.Encode(), res)
	return
}
