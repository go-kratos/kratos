package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/search/model"

	elastic "gopkg.in/olivere/elastic.v5"
)

// ReplyRecord search reply record from ES.
func (d *Dao) ReplyRecord(c context.Context, p *model.ReplyRecordParams) (res *model.SearchResult, err error) {
	query := elastic.NewBoolQuery()
	if p.Mid > 0 {
		query = query.Must(elastic.NewTermQuery("mid", p.Mid))
	} else {
		return
	}
	if len(p.Types) > 0 {
		interfaceSlice := make([]interface{}, len(p.Types))
		for i, d := range p.Types {
			interfaceSlice[i] = d
		}
		query = query.Must(elastic.NewTermsQuery("type", interfaceSlice...))
	}
	if len(p.States) > 0 {
		interfaceSlice := make([]interface{}, len(p.States))
		for i, d := range p.States {
			interfaceSlice[i] = d
		}
		query = query.Must(elastic.NewTermsQuery("state", interfaceSlice...))
	}
	if p.CTimeFrom != "" {
		query = query.Must(elastic.NewRangeQuery("ctime").Gte(p.CTimeFrom))
	}
	if p.CTimeTo != "" {
		query = query.Must(elastic.NewRangeQuery("ctime").Lte(p.CTimeTo))
	}
	indexName := fmt.Sprintf("replyrecord_%d", p.Mid%100)
	if res, err = d.searchResult(c, "replyExternal", indexName, query, p.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
		return
	}
	return
}
