package service

import (
	"context"
	"fmt"
	"go-common/app/service/bbq/search/api/grpc/v1"
	"go-common/app/service/bbq/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"time"

	"github.com/json-iterator/go"
	"gopkg.in/olivere/elastic.v5"
)

//SaveVideo 保存视频信息
func (s *Service) SaveVideo(c context.Context, req *v1.SaveVideoRequest) (res *v1.SaveVideoResponse, err error) {
	if err = s.dao.SaveVideo(c, req); err != nil {
		err = ecode.SearchCreateIndexErr
	}
	res = &v1.SaveVideoResponse{}
	return
}

//RecVideoData 获取视频信息,给推荐使用
func (s *Service) RecVideoData(c context.Context, req *v1.RecVideoDataRequest) (res *v1.RecVideoDataResponse, err error) {
	time1 := time.Now().UnixNano()
	esParams := new(model.EsParam)
	esParams.Query = make(map[string]map[string]interface{})

	params := new(model.Query)
	jsoniter.Unmarshal([]byte(req.Query), params)

	esParams.From = params.From
	esParams.Size = 5000
	if params.Size > 0 {
		esParams.Size = params.Size
	}

	b := make(map[string]interface{})
	if params.Where != nil {
		must := make([]interface{}, 0)
		if params.Where.In != nil {
			for k, v := range params.Where.In {
				ms := make(map[string]interface{})
				temp := make(map[string]interface{})
				switch v[0].(type) {
				case string:
					k += ".keyword"
				}
				temp[k] = v
				ms["terms"] = temp
				must = append(must, ms)
			}
		}
		if params.Where.Gte != nil {
			for k, v := range params.Where.Gte {
				rg := make(map[string]interface{})
				it := make(map[string]interface{})
				kw := make(map[string]interface{})
				kw["gte"] = v
				it[k] = kw
				rg["range"] = it
				must = append(must, rg)
			}
		}
		if params.Where.Lte != nil {
			for k, v := range params.Where.Lte {
				rg := make(map[string]interface{})
				it := make(map[string]interface{})
				kw := make(map[string]interface{})
				kw["lte"] = v
				it[k] = kw
				rg["range"] = it
				must = append(must, rg)
			}
		}
		b["must"] = must
		if params.Where.NotIn != nil {
			mustNot := make([]interface{}, 0)
			for k, v := range params.Where.NotIn {
				ms := make(map[string]interface{})
				temp := make(map[string]interface{})
				switch v[0].(type) {
				case string:
					k += ".keyword"
				}
				temp[k] = v
				ms["terms"] = temp
				mustNot = append(mustNot, ms)
			}
			b["must_not"] = mustNot
		}
	}

	if len(params.Filter) != 0 {
		ft := make(map[string]interface{})
		sc := make(map[string]interface{})
		sc1 := make(map[string]interface{})
		sc1["source"] = "bloom_filter"
		sc1["lang"] = "native"
		sc1["params"] = params.Filter
		sc["script"] = sc1
		ft["script"] = sc
		b["filter"] = ft
	}

	esParams.Query["bool"] = b

	esParams.Sort = make([]map[string]*model.Script, 0)
	if params.Calc != nil && params.Calc.Open == 1 {
		tmp1 := make(map[string]float64)
		tmp1["play_ratio"] = params.Calc.PlayRatio
		tmp1["fav_ratio"] = params.Calc.FavRatio
		tmp1["like_ratio"] = params.Calc.LikeRatio
		tmp1["share_ratio"] = params.Calc.ShareRatio
		tmp1["coin_ratio"] = params.Calc.CoinRatio
		tmp1["reply_ratio"] = params.Calc.ReplyRatio

		sc := new(model.Script)
		sc.Order = "desc"
		sc.Type = "number"
		sc.Script = make(map[string]interface{})
		sc.Script["inline"] = "params.play_ratio * Math.log1p(Math.min(10000, doc['play_hive'].value))/4.0 + 1.0 * params.fav_ratio * (1+doc['fav_hive'].value)/(100+doc['play_hive'].value) + 1.0 * params.like_ratio * (1+doc['likes_hive'].value)/(100+doc['play_hive'].value) + 1.0 * params.share_ratio * (1+doc['share_hive'].value)/(100+doc['play_hive'].value) + 1.0 * params.coin_ratio * (1+doc['coin_hive'].value)/(100+doc['play_hive'].value) + 1.0 * params.reply_ratio * (1+doc['reply_hive'].value)/(100+doc['play_hive'].value)"
		sc.Script["params"] = tmp1
		sc2 := make(map[string]*model.Script)
		sc2["_script"] = sc
		esParams.Sort = append(esParams.Sort, sc2)
	}

	b1, _ := jsoniter.Marshal(esParams)

	res = new(v1.RecVideoDataResponse)
	//res.Total, res.List, err = s.dao.RecVideoData2(c, string(b1))

	if res.Total, res.List, err = s.dao.ESVideoData(c, string(b1)); err != nil {
		err = ecode.SearchVideoDataErr
	}
	log.Infov(c, log.KV("method", "RecVideoData"), log.KV("time", time.Now().UnixNano()-time1), log.KV("total", res.Total), log.KV("list", len(res.List)))
	return
}

//VideoData 获取视频信息,给mis端使用
func (s *Service) VideoData(c context.Context, req *v1.VideoDataRequest) (res *v1.VideoDataResponse, err error) {
	params := new(model.Query)
	jsoniter.Unmarshal([]byte(req.Query), params)

	boolQuery := elastic.NewBoolQuery()
	if params.Where != nil {
		if params.Where.In != nil {
			for k, v := range params.Where.In {

				boolQuery.Must(elastic.NewTermsQuery(k, v...))
			}
		}
		if params.Where.NotIn != nil {
			for k, v := range params.Where.NotIn {
				boolQuery.MustNot(elastic.NewTermsQuery(k, v...))
			}
		}
		if params.Where.Lte != nil {
			for k, v := range params.Where.Lte {
				r := elastic.NewRangeQuery(k)
				r.Lte(v)
				boolQuery.Must(r)
			}
		}
		if params.Where.Gte != nil {
			for k, v := range params.Where.Gte {
				fmt.Println(k, v)
				r := elastic.NewRangeQuery(k)
				r.Gte(v)
				boolQuery.Must(r)
			}
		}
	}

	size := 10
	if params.Size > 0 {
		size = params.Size
	}

	res = new(v1.VideoDataResponse)
	if res.Total, res.List, err = s.dao.VideoData(c, boolQuery, params.From, size); err != nil {
		err = ecode.SearchVideoDataErr
	}
	return
}

// ESVideoData es原生查找
func (s *Service) ESVideoData(c context.Context, req *v1.ESVideoDataRequest) (res *v1.ESVideoDataResponse, err error) {
	res = new(v1.ESVideoDataResponse)
	if res.Total, res.List, err = s.dao.ESVideoData(c, req.Query); err != nil {
		err = ecode.SearchVideoDataErr
	}
	return
}

// DelVideoBySVID 批量删除视频
func (s *Service) DelVideoBySVID(c context.Context, req *v1.DelVideoBySVIDRequest) (res *v1.DelVideoBySVIDResponse, err error) {
	res = new(v1.DelVideoBySVIDResponse)
	err = nil
	if len(req.SVIDs) == 0 {
		return
	}
	for _, svid := range req.SVIDs {
		s.dao.DelVideoDataBySVID(c, svid)
	}
	return
}
