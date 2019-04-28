package dao

import (
	"context"
	"go-common/app/service/bbq/search/api/grpc/v1"
	"go-common/app/service/bbq/search/conf"
	"go-common/library/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/json-iterator/go"

	"gopkg.in/olivere/elastic.v5"
)

const (
	_bbqEsName    = "bbq"
	_videoIndex   = "video"
	_videoType    = "video_info"
	_videoMapping = `
{
	"settings":{
		"number_of_shards":5,
		"number_of_replicas":2,
		"index":{
			"analysis.analyzer.default.type":"ik_smart"
		}
	}
}
`
)

//SaveVideo 保存视频信息
func (d *Dao) SaveVideo(c context.Context, videos *v1.SaveVideoRequest) (err error) {
	d.createESIndex(_bbqEsName, _videoIndex, _videoMapping)
	bulkRequest := d.esPool[_bbqEsName].Bulk()
	for _, v := range videos.List {
		request := elastic.NewBulkUpdateRequest().Index(_videoIndex).Type(_videoType).Id(strconv.Itoa(int(v.SVID))).Doc(v).DocAsUpsert(true)
		bulkRequest.Add(request)
	}
	if _, err = bulkRequest.Do(c); err != nil {
		log.Error("save es [%d] err(%v)", 1, err)
	}
	return
}

//RecVideoDataElastic 获取视频信息
func (d *Dao) RecVideoDataElastic(c context.Context, query elastic.Query, script *elastic.ScriptSort, from, size int) (total int64, list []*v1.RecVideoInfo, err error) {
	search := d.esPool[_bbqEsName].Search().Index(_videoIndex).Type(_videoType)
	if query != nil {
		search.Query(query)
	}
	if script != nil {
		search.SortBy(script)
	}
	log.Error("start time(%d)", time.Now().UnixNano())
	res, err := search.From(from).Size(size).Timeout(d.c.Es[_bbqEsName].Timeout).Do(c)
	if err != nil {
		log.Error("video search es (%s) err(%v)", _bbqEsName, err)
		return
	}
	log.Error("do time(%d)", time.Now().UnixNano())
	total = res.TotalHits()
	list = []*v1.RecVideoInfo{}
	for _, value := range res.Hits.Hits {
		tmp := new(v1.RecVideoInfo)
		byte, _ := jsoniter.Marshal(value.Source)
		jsoniter.Unmarshal(byte, tmp)
		if value.Score != nil {
			tmp.ESScore = float64(*value.Score)
		}
		if value.Sort != nil {
			for _, v := range value.Sort {
				tmp.CustomScore = append(tmp.CustomScore, v.(float64))
			}
		}
		list = append(list, tmp)
	}
	return
}

//VideoData 获取视频信息
func (d *Dao) VideoData(c context.Context, query elastic.Query, from, size int) (total int64, list []*v1.VideoESInfo, err error) {
	res, err := d.esPool[_bbqEsName].Search().Index(_videoIndex).Type(_videoType).Query(query).From(from).Size(size).Timeout(d.c.Es[_bbqEsName].Timeout).Do(c)
	if err != nil {
		log.Error("video search es (%s) err(%v)", _bbqEsName, err)
		return
	}
	total = res.TotalHits()
	list = []*v1.VideoESInfo{}
	for _, value := range res.Hits.Hits {
		tmp := new(v1.VideoESInfo)
		byte, _ := jsoniter.Marshal(value.Source)
		jsoniter.Unmarshal(byte, tmp)
		list = append(list, tmp)
	}
	return
}

//ESVideoData 获取视频信息
func (d *Dao) ESVideoData(c context.Context, query string) (total int64, list []*v1.RecVideoInfo, err error) {
	i := rand.Intn(len(conf.Conf.Es["bbq"].Addr))
	req, err := http.NewRequest("POST", conf.Conf.Es["bbq"].Addr[i]+"/video/_search", strings.NewReader(query))
	if err != nil {
		log.Error("conn es err(%v)", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	j := rand.Intn(len(d.httpClient))
	res, err := d.httpClient[j].Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Error("conn es http err(%v)", err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("es read body err(%v)", err)
		return
	}
	log.Infov(c, log.KV("query", query), log.KV("response", string(body)))

	videos := new(elastic.SearchResult)
	jsoniter.Unmarshal(body, &videos)
	if videos == nil {
		return
	}

	list = make([]*v1.RecVideoInfo, 0)
	total = videos.TotalHits()
	for _, value := range videos.Hits.Hits {
		tmp := new(v1.RecVideoInfo)
		byte, _ := jsoniter.Marshal(value.Source)
		jsoniter.Unmarshal(byte, tmp)
		list = append(list, tmp)
	}
	return
}

// DelVideoDataBySVID 根据svid删除视频
func (d *Dao) DelVideoDataBySVID(c context.Context, svid int64) (err error) {
	i := rand.Intn(len(conf.Conf.Es["bbq"].Addr))
	url := conf.Conf.Es["bbq"].Addr[i] + "/video/video_info/" + strconv.Itoa(int(svid))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Error("conn es err(%v)", err)
		return
	}
	j := rand.Intn(len(d.httpClient))
	res, err := d.httpClient[j].Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Error("conn read body err(%v)", err)
	}
	return
}
