package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"go-common/app/admin/main/search/model"
	"go-common/library/log"

	"gopkg.in/olivere/elastic.v5"
)

const (
	_queryConfSQL = `select appid,es_name,index_prefix,index_type,index_id,index_mapping,query_max_indexes from digger_app`
)

// QueryConf query conf
func (d *Dao) QueryConf(ctx context.Context) (res map[string]*model.QueryConfDetail, err error) {
	rows, err := d.queryConfStmt.Query(ctx)
	if err != nil {
		log.Error("d.queryConfStmt.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[string]*model.QueryConfDetail)
	for rows.Next() {
		var (
			appid string
			qcd   = new(model.QueryConfDetail)
		)
		if err = rows.Scan(&appid, &qcd.ESCluster, &qcd.IndexPrefix, &qcd.IndexType, &qcd.IndexID, &qcd.IndexMapping, &qcd.MaxIndicesNum); err != nil {
			log.Error("d.QueryConf() rows.Scan() error(%v)", err)
			return
		}
		res[appid] = qcd
	}
	err = rows.Err()
	return
}

type querysModel struct {
	field     string
	whereKind string
	esQuery   elastic.Query
}

// QueryBasic 其中boolQuery方便定制化业务传参过来.
func (d *Dao) QueryBasic(c context.Context, sp *model.QueryParams) (mixedQuery *elastic.BoolQuery, qbDebug *model.QueryDebugResult) {
	mixedQuery = elastic.NewBoolQuery()
	qbDebug = &model.QueryDebugResult{}
	querys := []*querysModel{}
	netstedQuerys := map[string]*elastic.BoolQuery{} // key: path  value: boolQuery
	//fields
	if len(sp.QueryBody.Fields) == 0 {
		sp.QueryBody.Fields = []string{}
	}
	//from done
	//where
	if sp.QueryBody.Where == nil {
		sp.QueryBody.Where = &model.QueryBodyWhere{} //要给个默认值
	}
	//where - eq
	for k, v := range sp.QueryBody.Where.EQ {
		querys = append(querys, &querysModel{
			field:     k,
			whereKind: "eq",
			esQuery:   elastic.NewTermQuery(k, v),
		})
	}
	//where - or
	for k, v := range sp.QueryBody.Where.Or {
		querys = append(querys, &querysModel{
			field:     k,
			whereKind: "or",
			esQuery:   elastic.NewTermQuery(k, v),
		})
	}
	//where - in
	for k, v := range sp.QueryBody.Where.In {
		if len(v) > 1024 {
			e := fmt.Sprintf("where in 超过1024 business(%s) error(%v)", sp.Business, v)
			log.Error(e)
			qbDebug.AddErrMsg(e)
			continue
		}
		querys = append(querys, &querysModel{
			field:     k,
			whereKind: "in",
			esQuery:   elastic.NewTermsQuery(k, v...),
		})
	}
	//where - range
	ranges, err := d.queryBasicRange(sp.QueryBody.Where.Range)
	if err != nil {
		qbDebug.AddErrMsg(err.Error())
	}
	for k, v := range ranges {
		querys = append(querys, &querysModel{
			field:     k,
			whereKind: "range",
			esQuery:   v,
		})
	}
	//where - combo
	for _, v := range sp.QueryBody.Where.Combo {
		//外面用bool+should+minimum包裹
		combo := elastic.NewBoolQuery()
		//里面每个子项也是bool+should+minimum
		cmbEQ := elastic.NewBoolQuery()
		cmbIn := elastic.NewBoolQuery()
		cmbRange := elastic.NewBoolQuery()
		cmbNotEQ := elastic.NewBoolQuery()
		cmbNotIn := elastic.NewBoolQuery()
		cmbNotRange := elastic.NewBoolQuery()
		//所有的minumum
		if v.Min.Min == 0 {
			v.Min.Min = 1
		}
		if v.Min.EQ == 0 {
			v.Min.EQ = 1
		}
		if v.Min.In == 0 {
			v.Min.In = 1
		}
		if v.Min.Range == 0 {
			v.Min.Range = 1
		}
		if v.Min.NotEQ == 0 {
			v.Min.NotEQ = 1
		}
		if v.Min.NotIn == 0 {
			v.Min.NotIn = 1
		}
		if v.Min.NotRange == 0 {
			v.Min.NotRange = 1
		}
		//子项should
		for _, vEQ := range v.EQ {
			for eqK, eqV := range vEQ {
				cmbEQ.Should(elastic.NewTermQuery(eqK, eqV))
			}
		}
		for _, vIn := range v.In {
			for inK, inV := range vIn {
				cmbIn.Should(elastic.NewTermsQuery(inK, inV...))
			}
		}
		for _, vRange := range v.Range {
			ranges, _ := d.queryBasicRange(vRange)
			for _, rangeV := range ranges {
				cmbRange.Should(rangeV)
			}
		}
		for _, notEQ := range v.NotEQ {
			for k, v := range notEQ {
				cmbNotEQ.Should(elastic.NewTermQuery(k, v))
			}
		}
		for _, notIn := range v.NotIn {
			for k, v := range notIn {
				cmbNotIn.Should(elastic.NewTermsQuery(k, v...))
			}
		}
		for _, notRange := range v.NotRange {
			ranges, _ := d.queryBasicRange(notRange)
			for _, v := range ranges {
				cmbNotRange.Should(v)
			}
		}
		//子项minimum
		if len(v.EQ) > 0 {
			combo.Should(cmbEQ.MinimumNumberShouldMatch(v.Min.EQ))
		}
		if len(v.In) > 0 {
			combo.Should(cmbIn.MinimumNumberShouldMatch(v.Min.In))
		}
		if len(v.Range) > 0 {
			combo.Should(cmbRange.MinimumNumberShouldMatch(v.Min.Range))
		}
		if len(v.NotEQ) > 0 {
			combo.MustNot(elastic.NewBoolQuery().Should(cmbNotEQ.MinimumNumberShouldMatch(v.Min.NotEQ)))
		}
		if len(v.NotIn) > 0 {
			combo.MustNot(elastic.NewBoolQuery().Should(cmbNotIn.MinimumNumberShouldMatch(v.Min.NotIn)))
		}
		if len(v.NotRange) > 0 {
			combo.MustNot(elastic.NewBoolQuery().Should(cmbNotRange.MinimumNumberShouldMatch(v.Min.NotRange)))
		}
		//合并子项
		mixedQuery.Filter(combo.MinimumNumberShouldMatch(v.Min.Min))
	}
	//where - like
	like, err := d.queryBasicLike(sp.QueryBody.Where.Like, sp.Business)
	if err != nil {
		qbDebug.AddErrMsg(err.Error())
	}
	for _, v := range like {
		querys = append(querys, &querysModel{
			whereKind: "like",
			esQuery:   v,
		})
	}
	//mixedQuery
	for _, q := range querys {
		// like  TODO like的map型字段也要支持must not和 nested
		if q.field == "" && q.whereKind == "like" {
			mixedQuery.Must(q.esQuery)
			continue
		}
		if q.field == "" {
			continue
		}
		// prepare nested 一个DSL只能出现一个nested，不然会有问题
		if mapField := strings.Split(q.field, "."); len(mapField) > 1 && mapField[0] != "" {
			if _, ok := netstedQuerys[mapField[0]]; !ok {
				netstedQuerys[mapField[0]] = elastic.NewBoolQuery()
			}
			if bl, ok := sp.QueryBody.Where.Not[q.whereKind][q.field]; ok && bl {
				// mixedQuery.Must(elastic.NewNestedQuery(mapField[0], elastic.NewBoolQuery().MustNot(q.esQuery)))
				netstedQuerys[mapField[0]].MustNot(q.esQuery)
				continue
			}
			// mixedQuery.Must(elastic.NewNestedQuery(mapField[0], elastic.NewBoolQuery().Must(q.esQuery)))
			netstedQuerys[mapField[0]].Must(q.esQuery)
			continue
		}
		// must not
		if bl, ok := sp.QueryBody.Where.Not[q.whereKind][q.field]; ok && bl {
			mixedQuery.MustNot(q.esQuery)
			continue
		}
		// should
		if q.whereKind == "or" {
			mixedQuery.Should(q.esQuery)
			mixedQuery.MinimumShouldMatch("1") // 暂时为1
			continue
		}
		// default
		mixedQuery.Filter(q.esQuery)
		// random order with seed
		if sp.QueryBody.OrderRandomSeed != "" {
			random := elastic.NewRandomFunction().Seed(sp.QueryBody.OrderRandomSeed)
			score := elastic.NewFunctionScoreQuery().Add(elastic.NewBoolQuery(), random)
			mixedQuery.Must(score)
		}
	}
	// insert nested
	for k, n := range netstedQuerys {
		mixedQuery.Must(elastic.NewNestedQuery(k, n))
	}
	// DSL
	if sp.DebugLevel != 0 {
		if src, e := mixedQuery.Source(); e == nil {
			if data, er := json.Marshal(src); er == nil {
				qbDebug.DSL = string(data)
			}
		}
	}

	return
}

// queryBasicRange .
func (d *Dao) queryBasicRange(rangeMap map[string]string) (rangeQuery map[string]*elastic.RangeQuery, err error) {
	rangeQuery = make(map[string]*elastic.RangeQuery)
	for k, v := range rangeMap {
		if r := strings.Trim(v, " "); r != "" {
			if rs := []rune(r); len(rs) > 3 {
				firstStr := string(rs[0:1])
				endStr := string(rs[len(rs)-1:])
				rangeStr := strings.Trim(v, "[]() ")
				FromTo := strings.Split(rangeStr, ",")
				if len(FromTo) != 2 {
					err = fmt.Errorf("sp.QueryBody.Where.Range Fromto err")
					continue
				}
				rQuery := elastic.NewRangeQuery(k)
				rc := 0
				if firstStr == "(" && strings.Trim(FromTo[0], " ") != "" {
					rQuery.Gt(strings.Trim(FromTo[0], " "))
					rc++
				}
				if firstStr == "[" && strings.Trim(FromTo[0], " ") != "" {
					rQuery.Gte(strings.Trim(FromTo[0], " "))
					rc++
				}
				if endStr == ")" && strings.Trim(FromTo[1], " ") != "" {
					rQuery.Lt(strings.Trim(FromTo[1], " "))
					rc++
				}
				if endStr == "]" && strings.Trim(FromTo[1], " ") != "" {
					rQuery.Lte(strings.Trim(FromTo[1], " "))
					rc++
				}
				if rc == 0 {
					continue
				}
				rangeQuery[k] = rQuery
			} else {
				// 范围格式有问题
				err = fmt.Errorf("sp.QueryBody.Where.Range range format err. error(%v)", v)
				continue
			}
		}
	}
	return
}

func (d *Dao) queryBasicLike(likeMap []model.QueryBodyWhereLike, business string) (likeQuery []elastic.Query, err error) {
	for _, v := range likeMap {
		if len(v.KW) == 0 {
			continue
		}
		switch v.Level {
		case model.LikeLevelHigh:
			kw := []string{}
			r := []rune(v.KW[0])
			for i := 0; i < len(r); i++ {
				if k := string(r[i : i+1]); !strings.ContainsAny(k, "~[](){}^?:\"\\/!+-=&* ") { //去掉特殊符号
					kw = append(kw, k)
				} else if len(kw) > 1 && kw[len(kw)-1:][0] != "*" {
					kw = append(kw, "*", " ", "*")
				}
			}
			if len(kw) == 0 || strings.Join(kw, "") == "* *" {
				continue
			}
			qs := elastic.NewQueryStringQuery("*" + strings.Trim(strings.Join(kw, ""), "* ") + "*").AllowLeadingWildcard(true) //默认是or
			if !v.Or {
				qs.DefaultOperator("AND")
			}
			for _, v := range v.KWFields {
				qs.Field(v)
			}
			likeQuery = append(likeQuery, qs)
		case model.LikeLevelMiddel:
			// 单个字要特殊处理
			if r := []rune(v.KW[0]); len(r) == 1 && len(v.KW) == 1 {
				qs := elastic.NewQueryStringQuery("*" + string(r[:]) + "*").AllowLeadingWildcard(true) //默认是or
				if !v.Or {
					qs.DefaultOperator("AND")
				}
				for _, v := range v.KWFields {
					qs.Field(v)
				}
				likeQuery = append(likeQuery, qs)
				continue
			}
			// 自定义analyzer时，multi_match无法使用minimum_should_match，默认为至少一个满足，导致结果集还是很大
			// ngram(2,2)
			for _, kw := range v.KW {
				rn := []rune(kw)
				for i := 0; i+1 < len(rn); i++ {
					kwStr := string(rn[i : i+2])
					for _, kwField := range v.KWFields {
						likeQuery = append(likeQuery, elastic.NewTermQuery(kwField, kwStr))
					}
				}
			}
		case "", model.LikeLevelLow:
			qs := elastic.NewMultiMatchQuery(strings.Join(v.KW, " "), v.KWFields...).Type("best_fields").TieBreaker(0.6).MinimumShouldMatch("90%") //默认是and
			// TODO 业务自定义match
			if business == "copyright" {
				qs.MinimumShouldMatch("10%")
			}
			if business == "academy_archive" {
				qs.MinimumShouldMatch("50%")
			}
			if v.Or {
				qs.Operator("OR")
			}
			likeQuery = append(likeQuery, qs)
		}
	}
	return
}

func (d *Dao) Scroll(c context.Context, sp *model.QueryParams) (res *model.QueryResult, debug *model.QueryDebugResult, err error) {
	var (
		tList    []json.RawMessage
		tLen     int
		ScrollID = ""
	)
	res = &model.QueryResult{}
	esCluster := sp.AppIDConf.ESCluster
	query, _ := d.QueryBasic(c, sp)
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
		searchResult, err := eSearch.Scroll().Index(sp.QueryBody.From).Type("base").
			Query(query).FetchSourceContext(fsc).Size(5000).Scroll("1m").ScrollId(ScrollID).SortBy(sorterSlice...).Do(c)
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
