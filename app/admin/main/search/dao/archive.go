package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/admin/main/search/model"

	"gopkg.in/olivere/elastic.v5"
)

// ArchiveCheck search archive check from ES.
func (d *Dao) ArchiveCheck(c context.Context, p *model.ArchiveCheckParams) (res *model.SearchResult, err error) {
	query := elastic.NewBoolQuery()
	if len(p.Bsp.KWs) > 0 {
		for _, v := range p.Bsp.KWs {
			if p.Bsp.Pattern == "equal" {
				query = query.Must(elastic.NewMultiMatchQuery(v, p.Bsp.KwFields...).Type("best_fields").TieBreaker(0.3).MinimumShouldMatch("100%"))
			} else {
				query = query.Should(elastic.NewMultiMatchQuery(v, p.Bsp.KwFields...).Type("best_fields").TieBreaker(0.3).MinimumShouldMatch("80%")).MinimumNumberShouldMatch(1)
			}
		}
	} else if p.Bsp.KW != "" { //高级搜索比下面的高
		query = query.Must(elastic.NewMultiMatchQuery(p.Bsp.KW, p.Bsp.KwFields...).Type("best_fields").TieBreaker(0.3).MinimumShouldMatch("100%"))
	}
	if p.FromIP != "" {
		query = query.Must(elastic.NewQueryStringQuery("*" + p.FromIP + "*").AllowLeadingWildcard(true).Field("from_ip"))
	}
	if len(p.Aids) > 0 {
		interfaceSlice := make([]interface{}, len(p.Aids))
		for i, d := range p.Aids {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("aid", interfaceSlice...))
	}
	if len(p.TypeIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.TypeIds))
		for i, d := range p.TypeIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("typeid", interfaceSlice...))
	}
	if len(p.Attrs) > 0 {
		interfaceSlice := make([]interface{}, len(p.Attrs))
		for i, d := range p.Attrs {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("attribute", interfaceSlice...))
	}
	if len(p.States) > 0 {
		interfaceSlice := make([]interface{}, len(p.States))
		for i, d := range p.States {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("state", interfaceSlice...))
	}
	if len(p.Mids) > 0 {
		interfaceSlice := make([]interface{}, len(p.Mids))
		for i, d := range p.Mids {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("mid", interfaceSlice...))
	}
	if p.MidFrom > 0 {
		query = query.Filter(elastic.NewRangeQuery("mid").Gte(p.MidFrom))
	}
	if p.MidTo > 0 {
		query = query.Filter(elastic.NewRangeQuery("mid").Lte(p.MidTo))
	}
	if p.DurationFrom > 0 {
		query = query.Filter(elastic.NewRangeQuery("duration").Gte(p.DurationFrom))
	}
	if p.DurationTo > 0 {
		query = query.Filter(elastic.NewRangeQuery("duration").Lte(p.DurationTo))
	}
	if p.TimeFrom != "" && (p.Time == "ctime" || p.Time == "mtime" || p.Time == "pubtime") {
		query = query.Filter(elastic.NewRangeQuery(p.Time).Gte(p.TimeFrom))
	}
	if p.TimeTo != "" && (p.Time == "ctime" || p.Time == "mtime" || p.Time == "pubtime") {
		query = query.Filter(elastic.NewRangeQuery(p.Time).Lte(p.TimeTo))
	}
	if res, err = d.searchResult(c, "ssd_archive", "archivecheck", query, p.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
	}
	return
}

// Video search video from ES (deprecated).
func (d *Dao) Video(c context.Context, p *model.VideoParams) (res *model.SearchResult, err error) {
	query := elastic.NewBoolQuery()
	if p.Bsp.KW != "" {
		query = query.Must(elastic.NewMultiMatchQuery(p.Bsp.KW, p.Bsp.KwFields...).Type("best_fields").TieBreaker(0.3))
	}
	if len(p.VIDs) > 0 {
		interfaceSlice := make([]interface{}, len(p.VIDs))
		for i, d := range p.VIDs {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("vid", interfaceSlice...))
	}
	if len(p.AIDs) > 0 {
		interfaceSlice := make([]interface{}, len(p.AIDs))
		for i, d := range p.AIDs {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("aid", interfaceSlice...))
	}
	if len(p.CIDs) > 0 {
		interfaceSlice := make([]interface{}, len(p.CIDs))
		for i, d := range p.CIDs {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("cid", interfaceSlice...))
	}
	if len(p.TIDs) > 0 {
		interfaceSlice := make([]interface{}, len(p.TIDs))
		for i, d := range p.TIDs {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("arc_typeid", interfaceSlice...))
	}
	if len(p.FileNames) > 0 {
		interfaceSlice := make([]interface{}, len(p.FileNames))
		for i, d := range p.FileNames {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("filename", interfaceSlice...))
	}
	if len(p.RelationStates) > 0 {
		interfaceSlice := make([]interface{}, len(p.RelationStates))
		for i, d := range p.RelationStates {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("relation_state", interfaceSlice...))
	}
	if len(p.ArcMids) > 0 {
		interfaceSlice := make([]interface{}, len(p.ArcMids))
		for i, d := range p.ArcMids {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("arc_mid", interfaceSlice...))
	}
	if len(p.ArcMids) > 0 {
		interfaceSlice := make([]interface{}, len(p.ArcMids))
		for i, d := range p.ArcMids {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("arc_mid", interfaceSlice...))
	}
	if p.TagID > 0 {
		query = query.Filter(elastic.NewTermQuery("tag_id", p.TagID))
	}
	if len(p.Status) > 0 {
		interfaceSlice := make([]interface{}, len(p.Status))
		for i, d := range p.Status {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("status", interfaceSlice...))
	}
	if len(p.XCodeState) > 0 {
		interfaceSlice := make([]interface{}, len(p.XCodeState))
		for i, d := range p.XCodeState {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("xcode_state", interfaceSlice...))
	}
	// 不再查库过滤arc_mid
	if p.UserType > 0 {
		query = query.Filter(elastic.NewTermQuery("user_type", p.UserType))
	}
	if p.DurationFrom > 0 {
		query = query.Filter(elastic.NewRangeQuery("duration").Gte(p.DurationFrom))
	}
	if p.DurationTo > 0 {
		query = query.Filter(elastic.NewRangeQuery("duration").Lte(p.DurationTo))
	}
	if p.OrderType == 1 {
		diffs := time.Now().Unix() - 1420041600
		days := fmt.Sprintf("%dd", diffs/(3600*24))
		score := elastic.NewFunctionScoreQuery().Add(elastic.NewTermQuery("user_type", 1), elastic.NewExponentialDecayFunction().FieldName("arc_senddate").Origin("2015-01-01 00:00:00").Scale(days).Offset("1d").Decay(0.8).Weight(float64(10000))).Add(nil, elastic.NewExponentialDecayFunction().FieldName("arc_senddate").Origin("2015-01-01 00:00:00").Scale(days).Offset("1d").Decay(0.8).Weight(float64(1)))
		query = query.Must(score)
		p.Bsp.Order = []string{}
	}
	if res, err = d.searchResult(c, "ssd_archive", "archive_video", query, p.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
	}
	return
}

// TaskQa .
func (d *Dao) TaskQa(c context.Context, p *model.TaskQa) (res *model.SearchResult, err error) {
	query := elastic.NewBoolQuery()
	if p.Bsp.KW != "" {
		query = query.Must(elastic.NewMultiMatchQuery(p.Bsp.KW, p.Bsp.KwFields...).Type("best_fields").TieBreaker(0.3))
	}
	if len(p.Ids) > 0 {
		interfaceSlice := make([]interface{}, len(p.Ids))
		for i, d := range p.Ids {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("id", interfaceSlice...))
	}
	if len(p.TaskIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.TaskIds))
		for i, d := range p.TaskIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("task_id", interfaceSlice...))
	}
	if len(p.Uids) > 0 {
		interfaceSlice := make([]interface{}, len(p.Uids))
		for i, d := range p.Uids {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("uid", interfaceSlice...))
	}
	if len(p.ArcTagIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.ArcTagIds))
		for i, d := range p.ArcTagIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("arc_tagid", interfaceSlice...))
	}
	if len(p.AuditTagIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.AuditTagIds))
		for i, d := range p.AuditTagIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("audit_tagid", interfaceSlice...))
	}
	if len(p.UpGroups) > 0 {
		interfaceSlice := make([]interface{}, len(p.UpGroups))
		for i, d := range p.UpGroups {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("up_groups", interfaceSlice...))
	}
	if len(p.ArcTitles) > 0 {
		interfaceSlice := make([]interface{}, len(p.ArcTitles))
		for i, d := range p.ArcTitles {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("arc_title", interfaceSlice...))
	}
	if len(p.ArcTypeIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.ArcTypeIds))
		for i, d := range p.ArcTypeIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("arc_typeid", interfaceSlice...))
	}
	if len(p.States) > 0 {
		interfaceSlice := make([]interface{}, len(p.States))
		for i, d := range p.States {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("state", interfaceSlice...))
	}
	if len(p.AuditStatuses) > 0 {
		interfaceSlice := make([]interface{}, len(p.AuditStatuses))
		for i, d := range p.AuditStatuses {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("audit_status", interfaceSlice...))
	}
	if p.FansFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("fans").Gte(p.FansFrom))
	}
	if p.FansTo != "" {
		query = query.Filter(elastic.NewRangeQuery("fans").Lte(p.FansTo))
	}
	if p.CtimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Gte(p.CtimeFrom))
	}
	if p.CtimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Lte(p.CtimeTo))
	}
	if p.FtimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("ftime").Gte(p.FtimeFrom))
	}
	if p.FtimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("ftime").Lte(p.FtimeTo))
	}
	if res, err = d.searchResult(c, "ssd_archive", p.Bsp.AppID, query, p.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
	}
	return
}

// ArchiveCommerce .
func (d *Dao) ArchiveCommerce(c context.Context, p *model.ArchiveCommerce) (res *model.SearchResult, err error) {
	query := elastic.NewBoolQuery()
	if p.Bsp.KW != "" {
		query = query.Must(elastic.NewMultiMatchQuery(p.Bsp.KW, p.Bsp.KwFields...).Type("best_fields").TieBreaker(0.3))
	}
	if len(p.Ids) > 0 {
		interfaceSlice := make([]interface{}, len(p.Ids))
		for i, d := range p.Ids {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("id", interfaceSlice...))
	}
	if len(p.Mids) > 0 {
		interfaceSlice := make([]interface{}, len(p.Mids))
		for i, d := range p.Mids {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("mid", interfaceSlice...))
	}
	if len(p.PTypeIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.PTypeIds))
		for i, d := range p.PTypeIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("ptypeid", interfaceSlice...))
	}
	if len(p.TypeIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.TypeIds))
		for i, d := range p.TypeIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("typeid", interfaceSlice...))
	}
	if len(p.States) > 0 {
		interfaceSlice := make([]interface{}, len(p.States))
		for i, d := range p.States {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("state", interfaceSlice...))
	}
	if len(p.Copyrights) > 0 {
		interfaceSlice := make([]interface{}, len(p.Copyrights))
		for i, d := range p.Copyrights {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("copyright", interfaceSlice...))
	}
	if len(p.OrderIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.OrderIds))
		for i, d := range p.OrderIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("order_id", interfaceSlice...))
	}
	if p.IsOrder == 1 {
		query = query.Filter(elastic.NewRangeQuery("order_id").Gt(0))
	}
	if p.IsOrder == 0 {
		query = query.MustNot(elastic.NewRangeQuery("order_id").Gt(0))
	}
	if p.IsOriginal == 1 {
		query = query.Filter(elastic.NewTermsQuery("copyright", 1))
	}
	if p.IsOriginal == 0 {
		query = query.MustNot(elastic.NewTermsQuery("copyright", 1))
	}
	if p.Action == "get_ptypeids" {
		if res, err = d.ArchiveCommercePTypeIds(c, query); err != nil {
			PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
		}
		return
	}
	if res, err = d.searchResult(c, "ssd_archive", "archive_commerce_v", query, p.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
	}
	return
}

// ArchiveCommercePTypeIds .
func (d *Dao) ArchiveCommercePTypeIds(c context.Context, query *elastic.BoolQuery) (res *model.SearchResult, err error) {
	res = &model.SearchResult{
		Result: []json.RawMessage{},
		Page:   &model.Page{},
	}
	aggs := elastic.NewTermsAggregation()
	aggs = aggs.Field("ptypeid").Size(1000)
	if _, ok := d.esPool["ssd_archive"]; !ok {
		PromError(fmt.Sprintf("es:集群不存在%s", "ssd_archive"), "s.dao.searchResult indexName:%s", "ssd_archive")
		res = &model.SearchResult{Debug: fmt.Sprintf("es:集群不存在%s, %s", "ssd_archive", res.Debug)}
		return
	}
	searchResult, err := d.esPool["ssd_archive"].Search().Index("archive_commerce_v").Query(query).Aggregation("group_by_ptypeid", aggs).Size(0).Do(context.Background())
	if err != nil {
		PromError(fmt.Sprintf("es:执行查询失败%s ", "ArchiveCommercePTypeIds"), "dao.log.ArchiveCommercePTypeIds(%v)", err)
		return
	}
	result, ok := searchResult.Aggregations.Terms("group_by_ptypeid")
	if !ok {
		PromError(fmt.Sprintf("es:Unmarshal%s ", "log"), "es:Unmarshal%v", err)
		return
	}
	for _, v := range result.Buckets {
		res.Result = append(res.Result, []byte(v.Key.(string)))
	}
	res.Page.Pn = 1
	res.Page.Ps = 1000
	res.Page.Total = int64(len(res.Result))
	return
}
