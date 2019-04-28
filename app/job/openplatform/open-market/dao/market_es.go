package dao

import (
	"context"
	"gopkg.in/olivere/elastic.v5"
	"strconv"
	"time"

	"go-common/app/job/openplatform/open-market/model"
	"go-common/library/log"
)

var (
	orderIndex   = "product_order"
	orderType    = "product_order_info"
	commentIndex = "open_mall"
	commentType  = "ugc"
	marketIndex  = "product_market"
	marketType   = "product_market_info"
)

const (
	_orderStatusComplete = 2
	_dateFormat          = "2006-01-02"
	_timeFormat          = "2006-01-02 15:04:05"
)

//OrderData fetch ordercount by project and days
func (d *Dao) OrderData(c context.Context, projectID int32, startTimeUnix int64) (orderData map[int32]int64, err error) {
	var (
		searchResult *elastic.SearchResult
		startTime    string
		firstTime    string
		firstDay     time.Time
		startDay     time.Time
		daysBefore   = -1
	)
	orderData = make(map[int32]int64)
	startTime = time.Unix(startTimeUnix, 0).Add(time.Hour * 24).Format(_dateFormat)
	firstTime = time.Unix(startTimeUnix, 0).Add(time.Hour * 24).Add(time.Hour * 24 * -30).Format(_dateFormat)
	startDay, _ = time.Parse(_dateFormat, startTime)
	firstDay, _ = time.Parse(_dateFormat, firstTime)
	for {
		daysBefore++
		startDay = startDay.Add(time.Hour * -24)
		if !(startDay.Before(firstDay)) {
			rangeQuery := elastic.NewRangeQuery("pay_time")
			rangeQuery.Gte(startDay.Format(_timeFormat))
			rangeQuery.Lte(startDay.Add(time.Hour * 24).Format(_timeFormat))
			query := elastic.NewBoolQuery()
			query.Must(elastic.NewMatchQuery("project_id", projectID))
			query.Must(elastic.NewMatchQuery("status", _orderStatusComplete))
			query.Must(elastic.NewExistsQuery("pay_time"))
			query.Must(rangeQuery)
			searchResult, err = d.es.Search().
				Index(orderIndex).Type(orderType).
				Query(query).
				Size(0).
				Timeout(d.c.ElasticSearch.Timeout).
				Do(c)
			orderData[int32(daysBefore)] = searchResult.TotalHits()
			continue
		}
		break
	}
	return
}

// CommentData get comment info from ugc es
func (d *Dao) CommentData(c context.Context, projectID int32, startTimeUnix int64) (commentData map[int32]int64, err error) {
	var (
		searchResult *elastic.SearchResult
		startTime    string
		firstTime    string
		firstDay     time.Time
		startDay     time.Time
		daysBefore   = -1
	)
	commentData = make(map[int32]int64)
	startTime = time.Unix(startTimeUnix, 0).Add(time.Hour * 24).Format(_dateFormat)
	firstTime = time.Unix(startTimeUnix, 0).Add(time.Hour * 24).Add(time.Hour * 24 * -30).Format(_dateFormat)
	startDay, _ = time.Parse(_dateFormat, startTime)
	firstDay, _ = time.Parse(_dateFormat, firstTime)
	for {
		daysBefore++
		startDay = startDay.Add(time.Hour * -24)
		if !(startDay.Before(firstDay)) {
			rangeQuery := elastic.NewRangeQuery("ctime")
			rangeQuery.Gte(startDay.Unix() * 1000)
			rangeQuery.Lte(startDay.Add(time.Hour*24).Unix() * 1000)
			query := elastic.NewBoolQuery()
			query.Must(elastic.NewMatchQuery("subjectId", projectID))
			query.Must(elastic.NewMatchQuery("subjectType", 2))
			query.Must(rangeQuery)
			searchResult, err = d.esUgc.Search().
				Index(commentIndex).Type(commentType).
				Query(query).
				Size(0).
				Timeout(d.c.ElasticSearch.Timeout).
				Do(c)
			commentData[int32(daysBefore)] = searchResult.TotalHits()
			continue
		}
		break
	}
	return
}

//SaveData save result to es
func (d *Dao) SaveData(c context.Context, project *model.Project) (err error) {
	exists, err := d.es.IndexExists(marketIndex).Do(context.Background())
	if err != nil {
		log.Error("check if index exists error (%v)", err)
		return
	}
	if !exists {
		if _, err = d.es.CreateIndex(marketIndex).Do(c); err != nil {
			log.Error("index name(%s) create err(%v)", marketIndex, err)
			return
		}
	}
	_, err = d.es.Index().
		Index(marketIndex).
		Type(marketType).
		Id(strconv.Itoa(int(project.ID))).
		BodyJson(project).
		Refresh("true").
		Do(c)
	if err != nil {
		log.Error("put es [%d] err(%v)", project.ID, err)
	}
	return
}
