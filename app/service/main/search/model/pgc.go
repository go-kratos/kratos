package model

import "fmt"

// PgcMediaParams .
type PgcMediaParams struct {
	Bsp             *BasicSearchParams
	MediaIds        []int64  `form:"media_ids,split" params:"media_ids"`
	SeasonIds       []int64  `form:"season_ids,split" params:"season_ids"`
	SeasonTypes     []int64  `form:"season_types,split" params:"season_types"`
	StyleIds        []int64  `form:"style_ids,split" params:"style_ids"`
	Status          int      `form:"status" params:"status" default:"-1000"`
	ReleaseDateFrom string   `form:"release_date_from" params:"release_date_from"`
	ReleaseDateTo   string   `form:"release_date_to" params:"release_date_to"`
	SeasonIDFrom    int      `form:"season_id_from" params:"season_id_from"`
	SeasonIDTo      int      `form:"season_id_to" params:"season_id_to"`
	ProducerIds     []int64  `form:"producer_ids,split" params:"producer_ids"`
	IsDeleted       int      `form:"is_deleted" params:"is_deleted" default:"0"`
	AreaIds         []string `form:"area_ids,split" params:"area_ids"`
	ScoreFrom       int      `form:"score_from" params:"score_from"`
	ScoreTo         int      `form:"score_to" params:"score_to"`
	IsFinish        string   `form:"is_finish" params:"is_finish"`
	SeasonVersions  []int64  `form:"season_versions,split" params:"season_versions"`
	SeasonStatuses  []int64  `form:"season_statuses,split" params:"season_statuses"`
	PubTimeFrom     string   `form:"pub_time_from" params:"pub_time_from"`
	PubTimeTo       string   `form:"pub_time_to" params:"pub_time_to"`
	SeasonMonths    []int64  `form:"season_months,split" params:"season_months"`
	LatestTimeFrom  string   `form:"latest_time_from" params:"latest_time_from"`
	LatestTimeTo    string   `form:"latest_time_to" params:"latest_time_to"`
	CopyrightInfos  []string `form:"copyright_infos,split" params:"copyright_infos"`
	CTimeFrom       string   `form:"ctime_from" params:"ctime_from"`
	CTimeTo         string   `form:"ctime_to" params:"ctime_to"`
	MTimeFrom       string   `form:"mtime_from" params:"mtime_from"`
	MTimeTo         string   `form:"mtime_to" params:"mtime_to"`
}

// PgcMediaUptParams .
type PgcMediaUptParams struct {
	MediaID int64 `json:"media_id"`
	Field   map[string]interface{}
}

// IndexName .
func (m *PgcMediaUptParams) IndexName() string {
	return "pgc_media"
}

// IndexType .
func (m *PgcMediaUptParams) IndexType() string {
	return "base"
}

// IndexID .
func (m *PgcMediaUptParams) IndexID() string {
	return fmt.Sprintf("%d", m.MediaID)
}

// PField .
func (m *PgcMediaUptParams) PField() map[string]interface{} {
	return m.Field
}
