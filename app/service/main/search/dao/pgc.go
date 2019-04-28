package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/search/model"

	elastic "gopkg.in/olivere/elastic.v5"
)

// PgcMedia .
func (d *Dao) PgcMedia(c context.Context, p *model.PgcMediaParams) (res *model.SearchResult, err error) {
	query := elastic.NewBoolQuery()
	if p.Bsp.KW != "" {
		query = query.Must(elastic.NewMultiMatchQuery(p.Bsp.KW, "title").Type("best_fields").TieBreaker(0.3))
	}
	if len(p.MediaIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.MediaIds))
		for i, d := range p.MediaIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("media_id", interfaceSlice...))
	}
	if len(p.SeasonIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.SeasonIds))
		for i, d := range p.SeasonIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("season_id", interfaceSlice...))
	}
	if len(p.SeasonTypes) > 0 {
		interfaceSlice := make([]interface{}, len(p.SeasonTypes))
		for i, d := range p.SeasonTypes {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("season_type", interfaceSlice...))
	}
	if len(p.StyleIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.StyleIds))
		for i, d := range p.StyleIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("style_id", interfaceSlice...))
	}
	if p.Status > -1000 {
		query = query.Filter(elastic.NewTermQuery("status", p.Status))
	}
	if p.ReleaseDateFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("release_date").Gte(p.ReleaseDateFrom))
	}
	if p.ReleaseDateTo != "" {
		query = query.Filter(elastic.NewRangeQuery("release_date").Lte(p.ReleaseDateTo))
	}
	if p.ReleaseDateFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("release_date").Gte(p.ReleaseDateFrom))
	}
	if p.ReleaseDateTo != "" {
		query = query.Filter(elastic.NewRangeQuery("release_date").Lte(p.ReleaseDateTo))
	}
	if p.SeasonIDFrom > 0 {
		query = query.Filter(elastic.NewRangeQuery("season_id").Gte(p.SeasonIDFrom))
	}
	if p.SeasonIDTo > 0 {
		query = query.Filter(elastic.NewRangeQuery("season_id").Lte(p.SeasonIDTo))
	}
	if len(p.ProducerIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.ProducerIds))
		for i, d := range p.ProducerIds {
			interfaceSlice[i] = d
		}
		query = query.Filter(elastic.NewTermsQuery("producer_id", interfaceSlice...))
	}
	if p.IsDeleted == 0 {
		query = query.MustNot(elastic.NewTermQuery("is_deleted", 1))
	}
	if len(p.AreaIds) > 0 {
		interfaceSlice := make([]interface{}, len(p.AreaIds))
		for i, o := range p.AreaIds {
			interfaceSlice[i] = o
		}
		query = query.Filter(elastic.NewTermsQuery("area_id", interfaceSlice...))
	}
	if p.ScoreFrom > 0 {
		query = query.Filter(elastic.NewRangeQuery("score_from").Gte(p.ScoreFrom))
	}
	if p.ScoreTo > 0 {
		query = query.Filter(elastic.NewRangeQuery("score_to").Lte(p.ScoreTo))
	}
	if p.IsFinish != "" {
		query = query.Filter(elastic.NewTermsQuery("is_finish", p.IsFinish))
	}
	if len(p.SeasonVersions) > 0 {
		interfaceSlice := make([]interface{}, len(p.SeasonVersions))
		for i, o := range p.SeasonVersions {
			interfaceSlice[i] = o
		}
		query = query.Filter(elastic.NewTermsQuery("season_version", interfaceSlice...))
	}
	if len(p.SeasonStatuses) > 0 {
		interfaceSlice := make([]interface{}, len(p.SeasonStatuses))
		for i, o := range p.SeasonStatuses {
			interfaceSlice[i] = o
		}
		query = query.Filter(elastic.NewTermsQuery("season_status", interfaceSlice...))
	}
	if p.PubTimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("pub_time").Gte(p.PubTimeFrom))
	}
	if p.PubTimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("pub_time").Lte(p.PubTimeTo))
	}
	if len(p.SeasonMonths) > 0 {
		interfaceSlice := make([]interface{}, len(p.SeasonMonths))
		for i, o := range p.SeasonMonths {
			interfaceSlice[i] = o
		}
		query = query.Filter(elastic.NewTermsQuery("season_month", interfaceSlice...))
	}
	if p.LatestTimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("latest_time").Gte(p.LatestTimeFrom))
	}
	if p.LatestTimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("latest_time").Lte(p.LatestTimeTo))
	}
	if len(p.CopyrightInfos) > 0 {
		interfaceSlice := make([]interface{}, len(p.CopyrightInfos))
		for i, o := range p.CopyrightInfos {
			interfaceSlice[i] = o
		}
		query = query.Filter(elastic.NewTermsQuery("copyright_info", interfaceSlice...))
	}
	p.Bsp.Source = []string{"media_id", "season_id", "season_type", "dm_count", "play_count", "fav_count", "score", "latest_time", "pub_time", "release_date"}
	if res, err = d.searchResult(c, "externalPublic", "pgc_media", query, p.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", p.Bsp.AppID), "%v", err)
	}
	return
}
