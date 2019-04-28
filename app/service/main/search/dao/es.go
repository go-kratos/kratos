package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-common/app/service/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"

	elastic "gopkg.in/olivere/elastic.v5"
)

// UpdateMapBulk .
func (d *Dao) UpdateMapBulk(c context.Context, esName string, bulkData []BulkMapItem) (err error) {
	if _, ok := d.esPool[esName]; !ok {
		PromError(fmt.Sprintf("es:集群不存在%s", esName), "s.dao.searchResult indexName:%s", esName)
		err = ecode.SearchUpdateIndexFailed
		return
	}
	bulkRequest := d.esPool[esName].Bulk()
	for _, b := range bulkData {
		request := elastic.NewBulkUpdateRequest().Index(b.IndexName()).Type(b.IndexType()).Id(b.IndexID()).Doc(b.PField()).DocAsUpsert(true)
		bulkRequest.Add(request)
	}
	if _, err = bulkRequest.Do(context.TODO()); err != nil {
		log.Error("esName(%s) bulk error(%v)", esName, err)
	}
	return
}

func (d *Dao) UpdateBulk(c context.Context, esName string, bulkData []BulkItem) (err error) {
	if _, ok := d.esPool[esName]; !ok {
		PromError(fmt.Sprintf("es:集群不存在%s", esName), "s.dao.searchResult indexName:%s", esName)
		err = ecode.SearchUpdateIndexFailed
		return
	}
	bulkRequest := d.esPool[esName].Bulk()
	for _, b := range bulkData {
		request := elastic.NewBulkUpdateRequest().Index(b.IndexName()).Type(b.IndexType()).Id(b.IndexID()).Doc(b).DocAsUpsert(true)
		//fmt.Println(request)
		bulkRequest.Add(request)
	}
	if bulkRequest.NumberOfActions() == 0 {
		return
	}
	if _, err = bulkRequest.Do(context.TODO()); err != nil {
		log.Error("esName(%s) bulk error(%v)", esName, err)
	}
	return
}

// searchResult get result from ES.
func (d *Dao) searchResult(c context.Context, esClusterName, indexName string, query elastic.Query, bsp *model.BasicSearchParams) (res *model.SearchResult, err error) {
	res = &model.SearchResult{Debug: ""}
	if bsp.Debug {
		var src interface{}
		if src, err = query.Source(); err == nil {
			var data []byte
			if data, err = json.Marshal(src); err == nil {
				res = &model.SearchResult{Debug: string(data)}
			}
		}
	}
	if _, ok := d.esPool[esClusterName]; !ok {
		PromError(fmt.Sprintf("es:集群不存在%s", esClusterName), "s.dao.searchResult indexName:%s", indexName)
		res = &model.SearchResult{Debug: fmt.Sprintf("es:集群不存在%s, %s", esClusterName, res.Debug)}
		return
	}
	// multi sort
	sorterSlice := []elastic.Sorter{}
	if bsp.KW != "" {
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
	fsc := elastic.NewFetchSourceContext(true).Include(bsp.Source...)
	searchResult, err := d.esPool[esClusterName].
		Search().Index(indexName).
		Query(query).
		SortBy(sorterSlice...).
		From((bsp.Pn - 1) * bsp.Ps).
		Size(bsp.Ps).
		Pretty(true).
		FetchSourceContext(fsc).
		Do(c)
	if err != nil {
		PromError(fmt.Sprintf("es:执行查询失败%s ", esClusterName), "%v", err)
		res = &model.SearchResult{Debug: res.Debug + "es:执行查询失败"}
		return
	}
	var data []json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		var t json.RawMessage
		e := json.Unmarshal(*hit.Source, &t)
		if e != nil {
			PromError(fmt.Sprintf("es:%s 索引有脏数据", esClusterName), "s.dao.SearchArchiveCheck(%d,%d) error(%v) ", bsp.Pn*bsp.Ps, bsp.Ps, e)
			continue
		}
		data = append(data, t)
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
