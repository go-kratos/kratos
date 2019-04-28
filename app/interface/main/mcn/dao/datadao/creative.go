package datadao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// HTTPDataHandle .
func (d *Dao) HTTPDataHandle(c context.Context, params url.Values, key string) (data interface{}, err error) {
	var (
		uri string
		res struct {
			Code    int             `json:"code"`
			Data    json.RawMessage `json:"data"`
			Message string          `json:"message"`
		}
	)
	if uri, err = d.getURI(key); err != nil {
		return
	}
	if err = d.bmClient.Get(c, uri, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		log.Error("d.bmClient.Get(%s,%d)", uri+"?"+params.Encode(), res.Code)
		err = ecode.Error(ecode.Int(res.Code), res.Message)
		return
	}
	data = res.Data
	return
}

// getURI .
func (d *Dao) getURI(key string) (uri string, err error) {
	var (
		ok  bool
		url struct {
			host string
			uri  string
		}
		_url = map[string]struct {
			host string
			uri  string
		}{
			"archives": {
				uri:  "/x/internal/creative/archives",
				host: d.Conf.Host.API,
			},
			"archiveHistoryList": {
				uri:  "/x/internal/creative/archive/history/list",
				host: d.Conf.Host.API,
			},
			"archiveVideos": {
				uri:  "/x/internal/creative/archive/videos",
				host: d.Conf.Host.API,
			},
			"dataArchive": {
				uri:  "/x/internal/creative/data/archive",
				host: d.Conf.Host.API,
			},
			"dataVideoQuit": {
				uri:  "/x/internal/creative/data/videoquit",
				host: d.Conf.Host.API,
			},
			"danmuDistri": {
				uri:  "/x/internal/creative/danmu/distri",
				host: d.Conf.Host.API,
			},
			"dataBase": {
				uri:  "/x/internal/creative/data/base",
				host: d.Conf.Host.API,
			},
			"dataTrend": {
				uri:  "/x/internal/creative/data/trend",
				host: d.Conf.Host.API,
			},
			"dataAction": {
				uri:  "/x/internal/creative/data/action",
				host: d.Conf.Host.API,
			},
			"dataFan": {
				uri:  "/x/internal/creative/data/fan",
				host: d.Conf.Host.API,
			},
			"dataPandect": {
				uri:  "/x/internal/creative/data/pandect",
				host: d.Conf.Host.API,
			},
			"dataSurvey": {
				uri:  "/x/internal/creative/data/survey",
				host: d.Conf.Host.API,
			},
			"dataPlaySource": {
				uri:  "/x/internal/creative/data/playsource",
				host: d.Conf.Host.API,
			},
			"dataPlayAnalysis": {
				uri:  "/x/internal/creative/data/playanalysis",
				host: d.Conf.Host.API,
			},
			"dataArticleRank": {
				uri:  "/x/internal/creative/data/article/rank",
				host: d.Conf.Host.API,
			},
		}
	)
	if url, ok = _url[key]; !ok {
		return uri, fmt.Errorf("url(%s) not exist", key)
	}
	uri = url.host + url.uri
	return uri, err
}
