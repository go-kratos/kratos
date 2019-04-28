package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/admin/openplatform/sug/model"
	"go-common/library/log"

	"gopkg.in/olivere/elastic.v5"
	"io/ioutil"
)

const (
	_seasonIndex = "%s_sug_job_season"
	_seasonType  = "sug_job_season"
)

// GetSeason get season from es.
func (d *Dao) GetSeason(ctx context.Context, seasonID int64) (season model.Season, err error) {
	seasonTermQuery := elastic.NewTermQuery("id", seasonID)
	searchResult, err := d.es.Search().Index(fmt.Sprintf(_seasonIndex, d.c.Env)).Type(_seasonType).Query(seasonTermQuery).From(0).Size(1).Timeout("1s").Do(ctx)
	if err != nil {
		log.Error("es search error(%v)", err)
		return
	}
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err = json.Unmarshal(*hit.Source, &season)
			if err != nil {
				log.Error("json.Unmarshal err(%v)", err)
				return season, err
			}
		}
	}
	return
}

// SeasonList search season list.
func (d *Dao) SeasonList(ctx context.Context, params *model.SourceSearch) (list []model.Season, err error) {
	query := elastic.NewBoolQuery()
	matchQuery := elastic.NewMatchQuery("title", params.Keyword).Fuzziness("40")
	sid, _ := strconv.Atoi(params.Keyword)
	termQuery := elastic.NewMatchQuery("id", sid).Boost(100)
	query.Should(matchQuery)
	query.Should(termQuery)
	searchResult, err := d.es.Search().Index(fmt.Sprintf(_seasonIndex, d.c.Env)).Type(_seasonType).Query(query).From(0).Size(10).Timeout("1s").Do(ctx)
	if err != nil {
		return
	}
	list = []model.Season{}
	if searchResult.Hits.TotalHits > 0 {
		var season model.Season
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &season)
			if err != nil {
				log.Error("json.Unmarshal error(%v)", err)
				continue
			}
			list = append(list, season)
		}
	}
	return
}

// ItemList mall items list from http.
func (d *Dao) ItemList(ctx context.Context, params *model.SourceSearch) (itemsList []model.Items, err error) {
	query := make(map[string]interface{})
	query["pageNum"] = params.PageNum
	query["pageSize"] = params.PageSize
	query["shopId"] = 0
	query["keyword"] = params.Keyword
	jsonQuery, _ := json.Marshal(query)
	resp, err := d.client.Post(d.c.URL.ItemSearch, "application/json", bytes.NewReader(jsonQuery))
	if err != nil {
		log.Error("Request error(%v)", err)
		return
	}
	HTTPResponse := model.HTTPResponse{}
	bodyJSON, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(bodyJSON, &HTTPResponse); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
	}
	if HTTPResponse.Code != 0 {
		log.Error("Request (%s) search error(%v)", d.c.URL.ItemSearch, err)
		return
	}
	itemsList = HTTPResponse.Data.List
	return
}
