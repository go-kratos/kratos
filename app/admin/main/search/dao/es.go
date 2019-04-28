package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"go-common/app/admin/main/search/model"
	"go-common/library/log"

	"gopkg.in/olivere/elastic.v5"
)

// UpdateMapBulk (Deprecated).
func (d *Dao) UpdateMapBulk(c context.Context, esName string, bulkData []BulkMapItem) (err error) {
	bulkRequest := d.esPool[esName].Bulk()
	for _, b := range bulkData {
		request := elastic.NewBulkUpdateRequest().Index(b.IndexName()).Type(b.IndexType()).Id(b.IndexID()).Doc(b.PField()).DocAsUpsert(true)
		bulkRequest.Add(request)
	}
	if _, err = bulkRequest.Do(c); err != nil {
		log.Error("esName(%s) bulk error(%v)", esName, err)
	}
	return
}

// UpdateBulk (Deprecated).
func (d *Dao) UpdateBulk(c context.Context, esName string, bulkData []BulkItem) (err error) {
	bulkRequest := d.esPool[esName].Bulk()
	for _, b := range bulkData {
		request := elastic.NewBulkUpdateRequest().Index(b.IndexName()).Type(b.IndexType()).Id(b.IndexID()).Doc(b).DocAsUpsert(true)
		bulkRequest.Add(request)
	}
	if _, err = bulkRequest.Do(c); err != nil {
		log.Error("esName(%s) bulk error(%v)", esName, err)
	}

	return
}

// UpsertBulk 为了替换UpdateMapBulk和UpdateBulk .
func (d *Dao) UpsertBulk(c context.Context, esCluster string, up *model.UpsertParams) (err error) {
	es, ok := d.esPool[esCluster]
	if !ok {
		log.Error("esCluster(%s) not exists", esCluster)
		return
	}
	bulkRequest := es.Bulk()
	for _, b := range up.UpsertBody {
		request := elastic.NewBulkUpdateRequest().Index(b.IndexName).Type(b.IndexType).Id(b.IndexID).Doc(b.Doc)
		if up.Insert {
			request.DocAsUpsert(true)
		}
		//fmt.Println(request)
		bulkRequest.Add(request)
	}
	if _, err = bulkRequest.Do(c); err != nil {
		log.Error("esCluster(%s) bulk error(%v)", esCluster, err)
	}
	return
}

// searchResult get result from ES. (Deprecated) v3迁移完要删掉.
func (d *Dao) searchResult(c context.Context, esClusterName, indexName string, query elastic.Query, bsp *model.BasicSearchParams) (res *model.SearchResult, err error) {
	res = &model.SearchResult{Debug: ""}
	if bsp.Debug {
		if src, e := query.Source(); e == nil {
			if data, er := json.Marshal(src); er == nil {
				res = &model.SearchResult{Debug: string(data)}
			} else {
				err = er
				log.Error("searchResult query.Source.json.Marshal error(%v)", err)
				return
			}
		} else {
			err = e
			log.Error("searchResult query.Source error(%v)", err)
			return
		}
	}
	if _, ok := d.esPool[esClusterName]; !ok {
		PromError(fmt.Sprintf("es:集群不存在%s", esClusterName), "s.dao.searchResult indexName:%s", indexName)
		res = &model.SearchResult{Debug: fmt.Sprintf("es:集群不存在%s, %s", esClusterName, res.Debug)}
		return
	}
	// multi sort
	sorterSlice := []elastic.Sorter{}
	if bsp.KW != "" && bsp.ScoreFirst {
		sorterSlice = append(sorterSlice, elastic.NewScoreSort().Desc())
	}
	for i, d := range bsp.Order {
		if len(bsp.Sort) < i+1 {
			if bsp.Sort[0] == "desc" {
				sorterSlice = append(sorterSlice, elastic.NewFieldSort(d).Desc())
			} else {
				sorterSlice = append(sorterSlice, elastic.NewFieldSort(d).Asc())
			}
		} else {
			if bsp.Sort[i] == "desc" {
				sorterSlice = append(sorterSlice, elastic.NewFieldSort(d).Desc())
			} else {
				sorterSlice = append(sorterSlice, elastic.NewFieldSort(d).Asc())
			}
		}
	}
	if bsp.KW != "" && !bsp.ScoreFirst {
		sorterSlice = append(sorterSlice, elastic.NewScoreSort().Desc())
	}
	// source
	fsc := elastic.NewFetchSourceContext(true).Include(bsp.Source...)
	// highlight
	hl := elastic.NewHighlight()
	if bsp.Highlight && len(bsp.KwFields) > 0 {
		for _, v := range bsp.KwFields {
			hl = hl.Fields(elastic.NewHighlighterField(v))
		}
		hl = hl.PreTags("<em class=\"keyword\">").PostTags("</em>")
	}
	// from + size = 10,000
	from := (bsp.Pn - 1) * bsp.Ps
	size := bsp.Ps
	if (from + size) > 10000 {
		from = 10000 - size
	}
	// do
	searchResult, err := d.esPool[esClusterName].
		Search().Index(indexName).
		Highlight(hl).
		Query(query).
		SortBy(sorterSlice...).
		From(from).
		Size(size).
		Pretty(true).
		FetchSourceContext(fsc).
		Do(context.Background())
	if err != nil {
		PromError(fmt.Sprintf("es:执行查询失败%s ", esClusterName), "%v", err)
		res = &model.SearchResult{Debug: res.Debug + "es:执行查询失败"}
		return
	}
	var data []json.RawMessage
	b := bytes.Buffer{}
	b.WriteString("{")
	b.WriteString("}")
	for _, hit := range searchResult.Hits.Hits {
		var t json.RawMessage
		e := json.Unmarshal(*hit.Source, &t)
		if e != nil {
			PromError(fmt.Sprintf("es:%s 索引有脏数据", esClusterName), "s.dao.SearchArchiveCheck(%d,%d) error(%v) ", bsp.Pn*bsp.Ps, bsp.Ps, e)
			continue
		}
		data = append(data, t)
		// highlight
		if len(hit.Highlight) > 0 {
			b, _ := json.Marshal(hit.Highlight)
			h := []byte(string(b))
			data = append(data, h)
		} else if bsp.Highlight {
			data = append(data, b.Bytes()) //保证在高亮情况下，肯定有一对数据
		}
	}
	if len(data) == 0 {
		data = []json.RawMessage{}
	}
	res = &model.SearchResult{
		Order:  strings.Join(bsp.Order, ","),
		Sort:   strings.Join(bsp.Sort, ","),
		Result: data,
		Debug:  res.Debug,
		Page: &model.Page{
			Pn:    bsp.Pn,
			Ps:    bsp.Ps,
			Total: searchResult.Hits.TotalHits,
		},
	}
	return
}

// QueryResult query result from ES.
func (d *Dao) QueryResult(c context.Context, query elastic.Query, sp *model.QueryParams, qbDebug *model.QueryDebugResult) (res *model.QueryResult, qrDebug *model.QueryDebugResult, err error) {
	qrDebug = &model.QueryDebugResult{}
	if qbDebug != nil {
		qrDebug = qbDebug
	}
	esCluster := sp.AppIDConf.ESCluster
	if _, ok := d.esPool[esCluster]; !ok {
		qrDebug.AddErrMsg("es:集群不存在" + esCluster)
		return
	}
	if sp.DebugLevel != 0 {
		qrDebug.Mapping, err = d.esPool[esCluster].GetMapping().Index(sp.QueryBody.From).Do(context.Background())
	}
	// 低级别debug，在dsl执行前退出
	if sp.DebugLevel == 1 {
		return
	}
	// multi sort
	sorterSlice := []elastic.Sorter{}
	if len(sp.QueryBody.Where.Like) > 0 && sp.QueryBody.OrderScoreFirst { // like 长度 > 0，但里面是空的也是个问题
		sorterSlice = append(sorterSlice, elastic.NewScoreSort().Desc())
	}
	for _, i := range sp.QueryBody.Order {
		for k, v := range i {
			if v == "asc" {
				sorterSlice = append(sorterSlice, elastic.NewFieldSort(k).Asc())
			} else {
				sorterSlice = append(sorterSlice, elastic.NewFieldSort(k).Desc())
			}
		}
	}
	if len(sp.QueryBody.Where.Like) > 0 && sp.QueryBody.OrderScoreFirst {
		sorterSlice = append(sorterSlice, elastic.NewScoreSort().Desc())
	}
	// source
	fsc := elastic.NewFetchSourceContext(true).Include(sp.QueryBody.Fields...)
	// highlight
	hl := elastic.NewHighlight()
	if sp.QueryBody.Highlight && len(sp.QueryBody.Where.Like) > 0 {
		for _, v := range sp.QueryBody.Where.Like {
			for _, field := range v.KWFields {
				hl = hl.Fields(elastic.NewHighlighterField(field))
			}
		}
		hl = hl.PreTags("<em class=\"keyword\">").PostTags("</em>")
	}
	// from + size = 10,000
	maxRows := 10000
	if b, ok := model.PermConf["oht"][sp.Business]; ok && b == "true" {
		maxRows = 100000
	}
	from := (sp.QueryBody.Pn - 1) * sp.QueryBody.Ps
	size := sp.QueryBody.Ps
	if (from + size) > maxRows {
		from = maxRows - size
	}
	// Scroll
	if sp.QueryBody.Scroll == true {
		var (
			tList    []json.RawMessage
			tLen     int
			ScrollID = ""
		)
		res = &model.QueryResult{}
		esCluster := sp.AppIDConf.ESCluster
		eSearch, ok := d.esPool[esCluster]
		if !ok {
			PromError(fmt.Sprintf("es:集群不存在%s", esCluster), "s.dao.searchResult indexName:%s", esCluster)
			return
		}
		fsc := elastic.NewFetchSourceContext(true).Include(sp.QueryBody.Fields...)
		// multi sort
		sorterSlice := []elastic.Sorter{}
		if len(sp.QueryBody.Where.Like) > 0 && sp.QueryBody.OrderScoreFirst { // like 长度 > 0，但里面是空的也是个问题
			sorterSlice = append(sorterSlice, elastic.NewScoreSort().Desc())
		}
		for _, i := range sp.QueryBody.Order {
			for k, v := range i {
				if v == "asc" {
					sorterSlice = append(sorterSlice, elastic.NewFieldSort(k).Asc())
				} else {
					sorterSlice = append(sorterSlice, elastic.NewFieldSort(k).Desc())
				}
			}
		}
		if len(sp.QueryBody.Where.Like) > 0 && !sp.QueryBody.OrderScoreFirst {
			sorterSlice = append(sorterSlice, elastic.NewScoreSort().Desc())
		}
		for {
			searchResult, err := eSearch.Scroll().Index(sp.QueryBody.From).
				Query(query).FetchSourceContext(fsc).Size(sp.QueryBody.Ps).Scroll("1m").ScrollId(ScrollID).SortBy(sorterSlice...).Do(c)
			if err == io.EOF {
				break
			} else if err != nil {
				PromError(fmt.Sprintf("es:执行查询失败%s ", "Scroll"), "es:执行查询失败%v", err)
				break
			}
			ScrollID = searchResult.ScrollId
			for _, hit := range searchResult.Hits.Hits {
				var t json.RawMessage
				if err = json.Unmarshal(*hit.Source, &t); err != nil {
					PromError(fmt.Sprintf("es:Unmarshal%s ", "Scroll"), "es:Unmarshal%v", err)
					break
				}
				tList = append(tList, t)
				tLen++
				if tLen >= sp.QueryBody.Pn*sp.QueryBody.Ps {
					goto ClearScroll
				}
			}
		}
	ClearScroll:
		go eSearch.ClearScroll().ScrollId(ScrollID).Do(context.Background())
		if res.Result, err = json.Marshal(tList); err != nil {
			PromError(fmt.Sprintf("es:Unmarshal%s ", "Scroll"), "es:Unmarshal%v", err)
			return
		}
		return
	}
	// do
	searchPrepare := d.esPool[esCluster].
		Search().Index(sp.QueryBody.From).
		Highlight(hl).
		Query(query).
		SortBy(sorterSlice...).
		From(from).
		Size(size).
		FetchSourceContext(fsc).IgnoreUnavailable(true).AllowNoIndices(true)
	if ec, ok := model.PermConf["es_cache"][sp.Business]; ok && ec == "true" {
		searchPrepare.RequestCache(true)
	}
	if rt, ok := model.PermConf["routing"][sp.Business]; ok {
		routing := make([]string, 0, 1)
		if sp.QueryBody.Where.EQ != nil {
			if eq, ok := sp.QueryBody.Where.EQ[rt]; ok {
				routing = append(routing, fmt.Sprintf("%v", eq))
			}
		}
		if sp.QueryBody.Where.In != nil {
			if in, ok := sp.QueryBody.Where.In[rt]; ok {
				for _, v := range in {
					routing = append(routing, fmt.Sprintf("%v", v))
				}
			}
		}
		if len(routing) == 0 {
			qrDebug.AddErrMsg("es:路由不存在" + rt)
			return
		}
		searchPrepare.Routing(routing...)
	}
	if sp.DebugLevel == 2 {
		searchPrepare.Profile(true)
	}
	// Enhanced
	for _, v := range sp.QueryBody.Where.Enhanced {
		aggKey := v.Mode + "_" + v.Field
		switch v.Mode {
		case model.EnhancedModeGroupBy:
			aggs := elastic.NewTermsAggregation()
			aggs = aggs.Field(v.Field).Size(1000) //要和业务方确定具体值
			searchPrepare.Aggregation(aggKey, aggs)

		case model.EnhancedModeCollapse, model.EnhancedModeDistinct:
			collapse := elastic.NewCollapseBuilder(v.Field).MaxConcurrentGroupRequests(1)
			innerHit := elastic.NewInnerHit().Name("last_one").Size(1)
			for _, v := range v.Order {
				for field, sort := range v {
					if sort == "desc" {
						innerHit.Sort(field, false)
					} else {
						innerHit.Sort(field, true)
					}
				}
			}
			if len(v.Order) > 0 {
				collapse.InnerHit(innerHit)
			}
			searchPrepare.Collapse(collapse)
		case model.EnhancedModeSum:
			aggs := elastic.NewSumAggregation()
			aggs = aggs.Field(v.Field)
			searchPrepare.Aggregation(aggKey, aggs)
		case model.EnhancedModeDistinctCount:
			aggs := elastic.NewCardinalityAggregation()
			aggs = aggs.Field(v.Field)
			searchPrepare.Aggregation(aggKey, aggs)
		}
	}

	searchResult, err := searchPrepare.Do(context.Background())
	if err != nil {
		qrDebug.AddErrMsg(fmt.Sprintf("es:执行查询失败%s. %v", esCluster, err))
		PromError(fmt.Sprintf("es:执行查询失败%s ", esCluster), "%v", err)
		return
	}
	// data
	data := json.RawMessage{}
	docHits := []json.RawMessage{}
	docBuckets := map[string][]map[string]*json.RawMessage{}
	b := bytes.Buffer{}
	b.WriteString("{")
	b.WriteString("}")
	for _, hit := range searchResult.Hits.Hits {
		var t json.RawMessage
		e := json.Unmarshal(*hit.Source, &t)
		if e != nil {
			PromError(fmt.Sprintf("es:%s 索引有脏数据", esCluster), "s.dao.SearchArchiveCheck(%d,%d) error(%v) ", sp.QueryBody.Pn*sp.QueryBody.Ps, sp.QueryBody.Ps, e)
			continue
		}
		docHits = append(docHits, t)
		// highlight
		if len(hit.Highlight) > 0 {
			b, _ := json.Marshal(hit.Highlight)
			docHits = append(docHits, b)
		} else if sp.QueryBody.Highlight {
			docHits = append(docHits, b.Bytes()) //保证在高亮情况下，肯定有一对数据
		}
	}
	if len(docHits) > 0 {
		if doc, er := json.Marshal(docHits); er != nil {
			qrDebug.AddErrMsg(fmt.Sprintf("es:Unmarshal docHits es:Unmarshal%v ", er))
			PromError(fmt.Sprintf("es:Unmarshal%s ", "docHits"), "es:Unmarshal%v", er)
		} else {
			data = doc
		}
	} else {
		h := bytes.Buffer{}
		h.WriteString("[")
		h.WriteString("]")
		data = h.Bytes()
	}
	// data overwrite
	for _, v := range sp.QueryBody.Where.Enhanced {
		key := v.Mode + "_" + v.Field
		switch v.Mode {
		case model.EnhancedModeGroupBy:
			result, ok := searchResult.Aggregations.Terms(key)
			if !ok {
				PromError(fmt.Sprintf("es:Unmarshal%s ", key), "es:Unmarshal%v", err)
				continue
			}
			for _, b := range result.Buckets {
				docBuckets[key] = append(docBuckets[key], b.Aggregations)
			}
			data = b.Bytes() //保证无数据情况下，有正常返回
		case model.EnhancedModeSum:
			result, ok := searchResult.Aggregations.Sum(key)
			if !ok {
				PromError(fmt.Sprintf("es:Unmarshal%s ", key), "es:Unmarshal%v", err)
				continue
			}
			docBuckets[key] = append(docBuckets[key], result.Aggregations)
			data = b.Bytes() //保证无数据情况下，有正常返回
		case model.EnhancedModeDistinctCount:
			result, ok := searchResult.Aggregations.Cardinality(key)
			if !ok {
				PromError(fmt.Sprintf("es:Unmarshal%s ", key), "es:Unmarshal%v", err)
				continue
			}
			docBuckets[key] = append(docBuckets[key], result.Aggregations)
			data = b.Bytes() //保证无数据情况下，有正常返回
		default:
			// other modes...
		}
	}
	if len(docBuckets) > 0 {
		if doc, er := json.Marshal(docBuckets); er != nil {
			qrDebug.AddErrMsg(fmt.Sprintf("es:Unmarshal docBuckets es:Unmarshal%v", er))
			PromError(fmt.Sprintf("es:Unmarshal%s ", "docBuckets"), "es:Unmarshal%v", er)
		} else {
			data = doc
		}
	}
	order := []string{}
	sort := []string{}
	for _, i := range sp.QueryBody.Order {
		for k, v := range i {
			order = append(order, k)
			sort = append(sort, v)
		}
	}
	res = &model.QueryResult{
		Order:  strings.Join(order, ","),
		Sort:   strings.Join(sort, ","),
		Result: data,
		Page: &model.Page{
			Pn:    sp.QueryBody.Pn,
			Ps:    sp.QueryBody.Ps,
			Total: searchResult.Hits.TotalHits,
		},
	}
	//（默认的debug）高级别debug，在dsl执行后退出
	if sp.DebugLevel == 2 {
		qrDebug.Profile = searchResult.Profile
		return
	}
	return
}

// BulkIndex .
func (d *Dao) BulkIndex(c context.Context, esName string, bulkData []BulkItem) (err error) {
	bulkRequest := d.esPool[esName].Bulk()
	for _, b := range bulkData {
		request := elastic.NewBulkIndexRequest().Index(b.IndexName()).Type(b.IndexType()).Id(b.IndexID()).Doc(b)
		bulkRequest.Add(request)
	}
	if _, err = bulkRequest.Do(c); err != nil {
		log.Error("esName(%s) bulk error(%v)", esName, err)
	}
	return
}

// ExistIndex .
func (d *Dao) ExistIndex(c context.Context, esClusterName, indexName string) (exist bool, err error) {
	if _, ok := d.esPool[esClusterName]; !ok {
		PromError(fmt.Sprintf("es:集群不存在%s", esClusterName), "s.dao.searchResult indexName:%s", indexName)
		err = fmt.Errorf("集群不存在")
		return
	}
	exist, err = d.esPool[esClusterName].IndexExists(indexName).Do(c)
	return
}
