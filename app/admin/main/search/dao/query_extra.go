package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"gopkg.in/olivere/elastic.v5"
)

// ArchiveVideoScore 稿件一审打分排序.
func (d *Dao) ArchiveVideoScore(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	query, qbDebug := d.QueryBasic(c, sp)
	// query append
	diffs := time.Now().Unix() - 1420041600
	days := fmt.Sprintf("%dd", diffs/(3600*24))
	score := elastic.NewFunctionScoreQuery().Add(elastic.NewTermQuery("user_type", 1), elastic.NewExponentialDecayFunction().FieldName("arc_senddate").Origin("2015-01-01 00:00:00").Scale(days).Offset("1d").Decay(0.8).Weight(float64(10000))).Add(nil, elastic.NewExponentialDecayFunction().FieldName("arc_senddate").Origin("2015-01-01 00:00:00").Scale(days).Offset("1d").Decay(0.8).Weight(float64(1)))
	query = query.Must(score)
	sp.QueryBody.Order = []map[string]string{}
	// do
	if res, debug, err = d.QueryResult(c, query, sp, qbDebug); err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Business), "%v", err)
	}
	return
}

// ArchiveScore 稿件二审打分排序.
func (d *Dao) ArchiveScore(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	query, qbDebug := d.QueryBasic(c, sp)
	// query append
	diffs := time.Now().Unix() - 1420041600
	days := fmt.Sprintf("%dd", diffs/(3600*24))
	score := elastic.NewFunctionScoreQuery().Add(elastic.NewTermQuery("user_type", 1), elastic.NewExponentialDecayFunction().FieldName("ctime").Origin("2015-01-01 00:00:00").Scale(days).Offset("1d").Decay(0.8).Weight(float64(10000))).Add(nil, elastic.NewExponentialDecayFunction().FieldName("ctime").Origin("2015-01-01 00:00:00").Scale(days).Offset("1d").Decay(0.8).Weight(float64(1)))
	query = query.Must(score)
	sp.QueryBody.Order = []map[string]string{}
	// do
	if res, debug, err = d.QueryResult(c, query, sp, qbDebug); err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Business), "%v", err)
	}
	return
}

// TaskQaRandom .
func (d *Dao) TaskQaRandom(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	random := elastic.NewRandomFunction()
	if sp != nil && sp.QueryBody != nil && sp.QueryBody.Where != nil && sp.QueryBody.Where.EQ != nil {
		if seed, ok := sp.QueryBody.Where.EQ["seed"]; ok {
			random = elastic.NewRandomFunction().Seed(seed)
			delete(sp.QueryBody.Where.EQ, "seed")
		}
	}
	query, qbDebug := d.QueryBasic(c, sp)
	if err != nil {
		PromError(fmt.Sprintf("es basic:%s ", sp.Business), "%v", err)
	}
	score := elastic.NewFunctionScoreQuery().Add(elastic.NewBoolQuery(), random)
	qy := elastic.NewBoolQuery().Must(query, score)
	if res, debug, err = d.QueryResult(c, qy, sp, qbDebug); err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Business), "%v", err)
	}
	return
}

// EsportsContestsDate 电竞右侧日历联动.
func (d *Dao) EsportsContestsDate(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	res = &model.QueryResult{}
	// query basic
	query, qbDebug := d.QueryBasic(c, sp)
	debug = qbDebug
	esCluster := sp.AppIDConf.ESCluster
	if _, ok := d.esPool[esCluster]; !ok {
		debug.AddErrMsg("es:集群不存在" + esCluster)
		return
	}
	aggs := elastic.NewTermsAggregation()
	fsc := elastic.NewFetchSourceContext(true).Include("ids")
	aggs = aggs.Field("stime").Size(1000).SubAggregation("top_ids_hits", elastic.NewTopHitsAggregation().FetchSourceContext(fsc).Size(1000))
	searchPrepare := d.esPool[esCluster].Search().Index(sp.QueryBody.From).Query(query).Aggregation("group_by_stime", aggs).Size(0)
	if sp.DebugLevel == 2 {
		searchPrepare.Profile(true)
	}
	searchResult, err := searchPrepare.Do(context.Background())
	if err != nil {
		debug.AddErrMsg(fmt.Sprintf("es:执行查询失败%s. %v", esCluster, err))
		PromError(fmt.Sprintf("es:执行查询失败%s ", esCluster), "%v", err)
		return
	}
	result, ok := searchResult.Aggregations.Terms("group_by_stime")
	if !ok {
		return
	}
	type hitDoc struct {
		Hits []struct {
			Source struct {
				IDs []string `json:"ids"`
			} `json:"_source"`
		} `json:"hits"`
	}
	type idsRes struct {
		Date string
		IDs  []string
	}
	ids := []idsRes{}
	for _, b := range result.Buckets {
		var hit hitDoc
		//b.KeyAsString
		if list, ok := b.Terms("top_ids_hits"); ok {
			a, _ := list.Aggregations["hits"].MarshalJSON()
			if err = json.Unmarshal(a, &hit); err != nil {
				return
			}
			for _, h := range hit.Hits {
				ids = append(ids, idsRes{
					Date: *b.KeyAsString,
					IDs:  h.Source.IDs,
				})
			}
		}
	}
	resDoc := map[string]int{}
	resDocTmp := map[string]map[string]bool{}
	for _, v := range ids {
		if _, ok := resDocTmp[v.Date]; !ok {
			resDocTmp[v.Date] = map[string]bool{}
		}
		for _, id := range v.IDs {
			resDocTmp[v.Date][id] = true
		}
	}
	for date, idList := range resDocTmp {
		resDoc[date] = len(idList)
	}
	if doc, er := json.Marshal(resDoc); er != nil {
		debug.AddErrMsg(fmt.Sprintf("es:Unmarshal docBuckets es:Unmarshal%v", er))
	} else {
		res.Result = doc
	}
	return
}

var (
	_pubed    = []interface{}{-40, 0, 10000, 1, 1001, 15000, 20000, 30000}
	_notpubed = []interface{}{-2, -4, -5, -11, -12, -16}
	_ispubing = []interface{}{-1, -6, -7, -8, -9, -10, -13, -15, -30}
	_all      = append(append(_pubed, _notpubed...), _ispubing...)
)

// CreativeArchiveSearch 创作中心
func (d *Dao) CreativeArchiveSearch(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	var (
		mid interface{}
		ok  bool
	)
	docBuckets := map[string]interface{}{}
	if sp == nil && sp.QueryBody == nil && sp.QueryBody.Where == nil && sp.QueryBody.Where.EQ == nil {
		return res, debug, ecode.RequestErr
	}
	if mid, ok = sp.QueryBody.Where.EQ["mid"]; !ok {
		return res, debug, ecode.RequestErr
	}
	// 列表
	if state, ok := sp.QueryBody.Where.EQ["state"]; ok {
		if sp.QueryBody.Where.In == nil {
			sp.QueryBody.Where.In = map[string][]interface{}{}
		}
		switch state {
		case "pubed":
			sp.QueryBody.Where.In["state"] = _pubed
		case "not_pubed":
			sp.QueryBody.Where.In["state"] = _notpubed
		case "is_pubing":
			sp.QueryBody.Where.In["state"] = _ispubing
		default:
			sp.QueryBody.Where.In["state"] = _all
		}
		delete(sp.QueryBody.Where.EQ, "state")
	} else {
		if sp.QueryBody.Where.In == nil {
			sp.QueryBody.Where.In = map[string][]interface{}{}
		}
		sp.QueryBody.Where.In["state"] = _all
	}
	query, qbDebug := d.QueryBasic(c, sp)
	if res, debug, err = d.QueryResult(c, query, sp, qbDebug); err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Business), "%v", err)
		return
	}
	docBuckets["vlist"] = res.Result
	// 类型统计
	typeFilter := elastic.NewBoolQuery().Must(elastic.NewTermsQuery("mid", mid))
	typeFilter = typeFilter.Filter(elastic.NewTermsQuery("state", _all...))
	for _, v := range sp.QueryBody.Where.Like {
		typeFilter = typeFilter.Filter(elastic.NewMultiMatchQuery(strings.Join(v.KW, " "), v.KWFields...).Type("best_fields").TieBreaker(0.6).MinimumShouldMatch("100%"))
	}
	typeAgg := elastic.NewTermsAggregation().Field("pid")
	request1 := elastic.NewSearchRequest().Index(sp.QueryBody.From).Type("base").Source(elastic.NewSearchSource().Query(typeFilter).Aggregation("pid", typeAgg))
	// 状态统计
	stateFilter := elastic.NewBoolQuery().Filter(elastic.NewTermsQuery("mid", mid))
	if pid, ok := sp.QueryBody.Where.EQ["pid"]; ok {
		stateFilter = stateFilter.Filter(elastic.NewTermsQuery("pid", pid))
	}
	for _, v := range sp.QueryBody.Where.Like {
		stateFilter = typeFilter.Filter(elastic.NewMultiMatchQuery(strings.Join(v.KW, " "), v.KWFields...).Type("best_fields").TieBreaker(0.6).MinimumShouldMatch("100%"))
	}
	stateAgg := elastic.NewFiltersAggregation().
		FilterWithName("pubed", elastic.NewTermsQuery("state", _pubed...)).
		FilterWithName("not_pubed", elastic.NewTermsQuery("state", _notpubed...)).
		FilterWithName("is_pubing", elastic.NewTermsQuery("state", _ispubing...))
	request2 := elastic.NewSearchRequest().Index(sp.QueryBody.From).Type("base").Source(elastic.NewSearchSource().Query(stateFilter).Aggregation("state", stateAgg))
	MultiRes, err := d.esPool[sp.AppIDConf.ESCluster].MultiSearch().Add(request1, request2).Do(c)
	if err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Business), "%v", err)
		return
	}
	// 取得数据
	tmp := map[string]interface{}{}
	json.Unmarshal(*MultiRes.Responses[0].Aggregations["pid"], &tmp)
	docBuckets["tlist"] = tmp["buckets"]
	tmp = map[string]interface{}{}
	json.Unmarshal(*MultiRes.Responses[1].Aggregations["state"], &tmp)
	docBuckets["plist"] = tmp["buckets"]
	if resResult, e := json.Marshal(docBuckets); e != nil {
		log.Error("CreativeArchiveSearch.json.error(%v)", e)
	} else {
		res.Result = resResult
	}
	return
}

// CreativeArchiveStaff 创作中心
func (d *Dao) CreativeArchiveStaff(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	docBuckets := map[string]interface{}{}
	if sp == nil || sp.QueryBody == nil || sp.QueryBody.Where == nil || sp.QueryBody.Where.Combo == nil || len(sp.QueryBody.Where.Combo) != 1 {
		return res, debug, ecode.RequestErr
	}
	combo := sp.QueryBody.Where.Combo[0]
	if len(combo.EQ) == 0 {
		return res, debug, ecode.RequestErr
	}
	queryListParams := &model.QueryParams{
		QueryBody: &model.QueryBody{
			Where: &model.QueryBodyWhere{
				Combo: sp.QueryBody.Where.Combo,
			},
		},
	}
	queryList, _ := d.QueryBasic(c, queryListParams)
	// 列表
	if state, ok := sp.QueryBody.Where.EQ["state"]; ok {
		if sp.QueryBody.Where.In == nil {
			sp.QueryBody.Where.In = map[string][]interface{}{}
		}
		switch state {
		case "pubed":
			sp.QueryBody.Where.In["state"] = _pubed
		case "not_pubed":
			sp.QueryBody.Where.In["state"] = _notpubed
		case "is_pubing":
			sp.QueryBody.Where.In["state"] = _ispubing
		default:
			sp.QueryBody.Where.In["state"] = _all
		}
		delete(sp.QueryBody.Where.EQ, "state")
	} else {
		if sp.QueryBody.Where.In == nil {
			sp.QueryBody.Where.In = map[string][]interface{}{}
		}
		sp.QueryBody.Where.In["state"] = _all
	}
	query, qbDebug := d.QueryBasic(c, sp)
	if res, debug, err = d.QueryResult(c, query, sp, qbDebug); err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Business), "%v", err)
		return
	}
	docBuckets["vlist"] = res.Result
	// 类型统计
	typeFilter := elastic.NewBoolQuery().Filter(queryList)
	typeFilter = typeFilter.Filter(elastic.NewTermsQuery("state", _all...))
	for _, v := range sp.QueryBody.Where.Like {
		typeFilter = typeFilter.Filter(elastic.NewMultiMatchQuery(strings.Join(v.KW, " "), v.KWFields...).Type("best_fields").TieBreaker(0.6).MinimumShouldMatch("90%"))
	}
	typeAgg := elastic.NewTermsAggregation().Field("pid")
	request1 := elastic.NewSearchRequest().Index(sp.QueryBody.From).Type("base").Source(elastic.NewSearchSource().Query(typeFilter).Aggregation("pid", typeAgg).Size(0))
	// 状态统计
	stateFilter := elastic.NewBoolQuery().Filter(queryList)
	if pid, ok := sp.QueryBody.Where.EQ["pid"]; ok {
		stateFilter = stateFilter.Filter(elastic.NewTermsQuery("pid", pid))
	}
	for _, v := range sp.QueryBody.Where.Like {
		stateFilter = typeFilter.Filter(elastic.NewMultiMatchQuery(strings.Join(v.KW, " "), v.KWFields...).Type("best_fields").TieBreaker(0.6).MinimumShouldMatch("90%"))
	}
	stateAgg := elastic.NewFiltersAggregation().
		// 稿件状态
		FilterWithName("pubed", elastic.NewTermsQuery("state", _pubed...)).
		FilterWithName("not_pubed", elastic.NewTermsQuery("state", _notpubed...)).
		FilterWithName("is_pubing", elastic.NewTermsQuery("state", _ispubing...))
	request2 := elastic.NewSearchRequest().Index(sp.QueryBody.From).Type("base").Source(elastic.NewSearchSource().Query(stateFilter).Aggregation("state", stateAgg).Size(0))
	MultiRes, err := d.esPool[sp.AppIDConf.ESCluster].MultiSearch().Add(request1, request2).Do(c)
	if err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Business), "%v", err)
		return
	}
	// 取得数据
	tmp := map[string]interface{}{}
	json.Unmarshal(*MultiRes.Responses[0].Aggregations["pid"], &tmp)
	docBuckets["tlist"] = tmp["buckets"]
	tmp = map[string]interface{}{}
	json.Unmarshal(*MultiRes.Responses[1].Aggregations["state"], &tmp)
	docBuckets["plist"] = tmp["buckets"]
	if resResult, e := json.Marshal(docBuckets); e != nil {
		log.Error("CreativeArchiveSearch.json.error(%v)", e)
	} else {
		res.Result = resResult
	}
	return
}

// CreativeArchiveStaff 创作中心
func (d *Dao) CreativeArchiveApply(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	var (
		applyStaffMid interface{}
		ok            bool
	)
	docBuckets := map[string]interface{}{}
	if sp == nil || sp.QueryBody == nil || sp.QueryBody.Where == nil || sp.QueryBody.Where.EQ == nil {
		return res, debug, ecode.RequestErr
	}
	if applyStaffMid, ok = sp.QueryBody.Where.EQ["apply_staff.apply_staff_mid"]; !ok {
		return res, debug, ecode.RequestErr
	}
	// 列表
	if state, ok := sp.QueryBody.Where.EQ["apply_staff.deal_state"]; ok {
		if sp.QueryBody.Where.In == nil {
			sp.QueryBody.Where.In = map[string][]interface{}{}
		}
		switch state {
		case "pending": //待处理
			sp.QueryBody.Where.In["apply_staff.deal_state"] = []interface{}{1}
		case "processed": //已处理
			sp.QueryBody.Where.In["apply_staff.deal_state"] = []interface{}{2}
		case "neglected": //已忽略
			sp.QueryBody.Where.In["apply_staff.deal_state"] = []interface{}{3}
		default:
			sp.QueryBody.Where.In["apply_staff.deal_state"] = []interface{}{1, 2, 3}
		}
		delete(sp.QueryBody.Where.EQ, "apply_staff.deal_state")
	} else {
		if sp.QueryBody.Where.In == nil {
			sp.QueryBody.Where.In = map[string][]interface{}{}
		}
		sp.QueryBody.Where.In["apply_staff.deal_state"] = []interface{}{1, 2, 3}
	}
	sp.QueryBody.Where.In["state"] = _all
	query, qbDebug := d.QueryBasic(c, sp)
	if res, debug, err = d.QueryResult(c, query, sp, qbDebug); err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Business), "%v", err)
		return
	}
	docBuckets["vlist"] = res.Result
	// 类型统计
	typeFilter := elastic.NewBoolQuery().Filter(
		elastic.NewTermsQuery("state", _all...),
		elastic.NewNestedQuery("apply_staff", elastic.NewBoolQuery().Must(
			elastic.NewTermQuery("apply_staff.apply_staff_mid", applyStaffMid),
			elastic.NewTermsQuery("apply_staff.deal_state", []interface{}{1, 2, 3}...),
		)),
	)
	for _, v := range sp.QueryBody.Where.Like {
		typeFilter = typeFilter.Filter(elastic.NewMultiMatchQuery(strings.Join(v.KW, " "), v.KWFields...).Type("best_fields").TieBreaker(0.6).MinimumShouldMatch("90%"))
	}
	typeAgg := elastic.NewTermsAggregation().Field("pid")
	request1 := elastic.NewSearchRequest().Index(sp.QueryBody.From).Type("base").Source(elastic.NewSearchSource().Query(typeFilter).Aggregation("pid", typeAgg).Size(0))
	// 状态统计
	stateFilter := elastic.NewBoolQuery().Filter(
		elastic.NewTermsQuery("state", _all...),
		elastic.NewNestedQuery("apply_staff", elastic.NewBoolQuery().Must(elastic.NewTermQuery("apply_staff.apply_staff_mid", applyStaffMid))),
	)
	if pid, ok := sp.QueryBody.Where.EQ["pid"]; ok {
		stateFilter = stateFilter.Filter(elastic.NewTermsQuery("pid", pid))
	}
	for _, v := range sp.QueryBody.Where.Like {
		stateFilter = typeFilter.Filter(elastic.NewMultiMatchQuery(strings.Join(v.KW, " "), v.KWFields...).Type("best_fields").TieBreaker(0.6).MinimumShouldMatch("90%"))
	}
	stateAgg := elastic.NewFiltersAggregation().
		FilterWithName("pending", elastic.NewNestedQuery("apply_staff", elastic.NewBoolQuery().Must(elastic.NewTermQuery("apply_staff.apply_staff_mid", applyStaffMid), elastic.NewTermQuery("apply_staff.deal_state", 1)))).
		FilterWithName("processed", elastic.NewNestedQuery("apply_staff", elastic.NewBoolQuery().Must(elastic.NewTermQuery("apply_staff.apply_staff_mid", applyStaffMid), elastic.NewTermQuery("apply_staff.deal_state", 2)))).
		FilterWithName("neglected", elastic.NewNestedQuery("apply_staff", elastic.NewBoolQuery().Must(elastic.NewTermQuery("apply_staff.apply_staff_mid", applyStaffMid), elastic.NewTermQuery("apply_staff.deal_state", 3))))
	request2 := elastic.NewSearchRequest().Index(sp.QueryBody.From).Type("base").Source(elastic.NewSearchSource().Query(stateFilter).Aggregation("state", stateAgg).Size(0))
	MultiRes, err := d.esPool[sp.AppIDConf.ESCluster].MultiSearch().Add(request1, request2).Do(c)
	if err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Business), "%v", err)
		return
	}
	// 取得数据
	tmp := map[string]interface{}{}
	json.Unmarshal(*MultiRes.Responses[0].Aggregations["pid"], &tmp)
	docBuckets["tlist"] = tmp["buckets"]
	tmp = map[string]interface{}{}
	json.Unmarshal(*MultiRes.Responses[1].Aggregations["state"], &tmp)
	docBuckets["plist"] = tmp["buckets"]
	if resResult, e := json.Marshal(docBuckets); e != nil {
		log.Error("CreativeArchiveSearch.json.error(%v)", e)
	} else {
		res.Result = resResult
	}
	return
}
