package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/service/main/search/model"

	elastic "gopkg.in/olivere/elastic.v5"
)

func (d *Dao) DmDateSearch(c context.Context, p *model.DmDateParams) (res *model.SearchResult, err error) {
	query := elastic.NewBoolQuery()
	indexName := "dm_date_" + strings.Replace(p.Month, "-", "_", -1)
	if p.Bsp.KW != "" {
		query = query.Must(elastic.NewRegexpQuery(p.Bsp.KwFields[0], ".*"+p.Bsp.KW+".*"))
	}
	if p.Oid != -1 {
		query = query.Filter(elastic.NewTermQuery("oid", p.Oid))
	}
	if p.Month != "" {
		query = query.Filter(elastic.NewTermQuery("month", p.Month))
	}
	if p.MonthFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("month").Gte(p.MonthFrom))
	}
	if p.MonthTo != "" {
		query = query.Filter(elastic.NewRangeQuery("month").Lte(p.MonthTo))
	}
	if res, err = d.searchResult(c, "dmExternal", indexName, query, p.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
	}
	return
}
