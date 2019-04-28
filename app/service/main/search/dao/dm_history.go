package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/search/model"

	elastic "gopkg.in/olivere/elastic.v5"
)

func (d *Dao) DmHistory(c context.Context, p *model.DmHistoryParams) (res *model.SearchResult, err error) {
	var (
		query     = elastic.NewBoolQuery()
		indexName = fmt.Sprintf("dm_search_%03d", p.Oid%1000)
	)
	if p.Bsp.KW != "" {
		query = query.Must(elastic.NewMultiMatchQuery(p.Bsp.KW, p.Bsp.KwFields...).Type("best_fields").TieBreaker(0.6))
	}
	if p.Oid != -1 {
		query = query.Filter(elastic.NewTermQuery("oidstr", p.Oid))
	}
	if len(p.States) > 0 {
		interfaceSlice := make([]interface{}, len(p.States))
		for k, m := range p.States {
			interfaceSlice[k] = m
		}
		query = query.Filter(elastic.NewTermsQuery("state", interfaceSlice...))
	}
	if p.CtimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Gte(p.CtimeFrom))
	}
	if p.CtimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Lte(p.CtimeTo))
	}

	fmt.Println(indexName)
	if res, err = d.searchResult(c, "dmExternal", indexName, query, p.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
	}
	return
}
