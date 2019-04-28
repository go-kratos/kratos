package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/search/model"

	elastic "gopkg.in/olivere/elastic.v5"
)

// DmSearch .
func (d *Dao) DmSearch(c context.Context, p *model.DmSearchParams) (res *model.SearchResult, err error) {
	var (
		query     = elastic.NewBoolQuery()
		indexName = fmt.Sprintf("dm_search_%03d", p.Oid%1000)
	)
	if p.Bsp.KW != "" {
		query = query.Must(elastic.NewRegexpQuery(p.Bsp.KwFields[0], ".*"+p.Bsp.KW+".*"))
	}
	if p.Oid != -1 {
		query = query.Filter(elastic.NewTermQuery("oid", p.Oid))
	}
	if p.Mid != -1 {
		query = query.Filter(elastic.NewTermQuery("mid", p.Mid))
	}
	if p.Mode != -1 {
		query = query.Filter(elastic.NewTermQuery("mode", p.Mode))
	}
	if p.Pool != -1 {
		query = query.Filter(elastic.NewTermQuery("pool", p.Pool))
	}
	if p.Progress != -1 {
		query = query.Filter(elastic.NewTermQuery("progress", p.Progress))
	}
	if len(p.States) > 0 {
		interfaceSlice := make([]interface{}, len(p.States))
		for k, m := range p.States {
			interfaceSlice[k] = m
		}
		query = query.Filter(elastic.NewTermsQuery("state", interfaceSlice...))
	}
	if p.Type != -1 {
		query = query.Filter(elastic.NewTermQuery("type", p.Type))
	}
	if len(p.AttrFormat) > 0 {
		interfaceSlice := make([]interface{}, len(p.AttrFormat))
		for k, m := range p.AttrFormat {
			interfaceSlice[k] = m
		}
		query = query.Filter(elastic.NewTermsQuery("attr_format", interfaceSlice...))
	}
	if p.CtimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Gte(p.CtimeFrom))
	}
	if p.CtimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Lte(p.CtimeTo))
	}
	if res, err = d.searchResult(c, "dmExternal", indexName, query, p.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
	}
	return
}
