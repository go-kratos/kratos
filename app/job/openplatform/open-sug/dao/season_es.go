package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/job/openplatform/open-sug/model"
	"go-common/library/log"

	"gopkg.in/olivere/elastic.v5"
)

// SeasonData
func (d *Dao) SeasonData(c context.Context, item *model.Item) (scoreSlice []model.Score, err error) {
	var searchResult *elastic.SearchResult
	query := elastic.NewMultiMatchQuery(item.Keywords, "title", "alias", "alias_search", "actors^1.25")
	searchResult, err = d.es.Search().
		Index(fmt.Sprintf("%s_%s", d.c.Env, d.c.ElasticSearch.Season.Index)).
		Type(d.c.ElasticSearch.Season.Type).
		Query(query).
		Timeout(d.c.ElasticSearch.Timeout).
		Do(c)
	if err != nil {
		return
	}
	if searchResult.TotalHits() > 0 {
		wishCount, _ := d.WishCount(c, item)
		commentCount, _ := d.CommentCount(c, item)
		salesCount, _ := d.SalesCount(c, item)
		for _, s := range searchResult.Hits.Hits {
			seasonJson, _ := s.Source.MarshalJSON()
			season := model.EsSeason{}
			json.Unmarshal(seasonJson, &season)
			scoreSlice = append(scoreSlice, model.Score{SeasonID: s.Id, Score: *s.Score, SeasonName: season.Title})
			switch {
			case wishCount > d.ItemWishMax[s.Id]:
				d.ItemWishMax[s.Id] = wishCount
			case wishCount != 0 && wishCount < d.ItemWishMin[s.Id]:
				d.ItemWishMin[s.Id] = wishCount
			case d.ItemWishMin[s.Id] == 0:
				d.ItemWishMin[s.Id] = wishCount
			}
			switch {
			case commentCount > d.ItemCommentMax[s.Id]:
				d.ItemCommentMax[s.Id] = commentCount
			case commentCount != 0 && commentCount < d.ItemCommentMin[s.Id]:
				d.ItemCommentMin[s.Id] = commentCount
			case d.ItemCommentMin[s.Id] == 0:
				d.ItemCommentMin[s.Id] = commentCount
			}
			switch {
			case salesCount > d.ItemSalesMax[s.Id]:
				d.ItemSalesMax[s.Id] = salesCount
			case salesCount != 0 && salesCount < d.ItemSalesMin[s.Id]:
				d.ItemSalesMin[s.Id] = salesCount
			case d.ItemSalesMin[s.Id] == 0:
				d.ItemSalesMin[s.Id] = salesCount
			}
		}
	}
	return
}

// Index ...
func (d *Dao) Index(ctx context.Context, index, typ string, id string, data interface{}) {
	resp, err := d.es.Index().Index(index).Type(typ).BodyJson(data).Id(id).Do(ctx)
	if err != nil {
		log.Error("索引写入失败 index(%s) type(%s) error(%v)", index, typ, err)
		return
	}
	if resp.Result != "" {
		log.Info("index(%s) type(%s) 创建成功", resp.Index, resp.Type)
	} else {
		log.Info("index(%s) type(%s) 更新成功", resp.Index, resp.Type)
	}
}

// IndexExists ...
func (d *Dao) IndexExists(ctx context.Context, index string) bool {
	e, err := d.es.IndexExists(index).Do(ctx)
	if err != nil {
		log.Error("检查索引是否存在出错 IndexExists(%s) error(%v)", index, err)
	}

	return e
}

// CreateIndex ...
func (d *Dao) CreateIndex(ctx context.Context, name string, mapping string) bool {
	resp, err := d.es.CreateIndex(name).BodyString(mapping).Do(ctx)
	if err != nil {
		log.Error("创建索引出错 CreateIndex(%s) error(%v)", name, err)
		return false
	}
	if !resp.Acknowledged {
		log.Error("创建索引失败 index(%s)", name)
	}
	return resp.Acknowledged
}
