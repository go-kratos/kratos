package dao

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_mediaURL    = "http://api.bilibili.co/pgc/internal/review/score/card"
	_setScoreURL = "http://api.bilibili.co/pgc/internal/review/score/rating"
	_delScoreURL = "http://api.bilibili.co/pgc/internal/review/score/delete"
)

// Media get media score.
func (d *Dao) Media(c context.Context, mediaID, mid int64) (res *model.MediaResult, err error) {
	if mediaID == 0 {
		return
	}
	params := url.Values{}
	params.Set("media_id", strconv.FormatInt(mediaID, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	resp := &model.MediaResp{}
	err = d.httpClient.Get(c, _mediaURL, "", params, &resp)
	if err != nil {
		PromError("media:番剧信息接口")
		log.Error("media: d.client.Get(%s) error(%+v)", _mediaURL+"?"+params.Encode(), err)
		return
	}
	if resp.Code != 0 {
		PromError("media:番剧信息接口")
		log.Error("media: url(%s) res: %+v", _mediaURL+"?"+params.Encode(), resp)
		err = ecode.Int(resp.Code)
		return
	}
	res = resp.Result
	return
}

// SetScore set media score.
func (d *Dao) SetScore(c context.Context, score, aid, mediaID, mid int64) (err error) {
	if score < 1 || score > 10 || score%2 != 0 {
		err = errors.New("评分分值错误")
		return
	}
	params := url.Values{}
	params.Set("media_id", strconv.FormatInt(mediaID, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("score", strconv.FormatInt(score, 10))
	params.Set("from", "1")
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	err = d.httpClient.Post(c, _setScoreURL, "", params, &resp)
	if err != nil {
		PromError("media:番剧评分接口")
		log.Error("media: d.client.Post(%s, data:%s) error(%+v)", _setScoreURL, params.Encode(), err)
		return
	}
	if resp.Code != 0 {
		PromError("media:番剧评分接口")
		log.Error("media: d.client.Post(%s, data:%s) res: %+v", _setScoreURL, params.Encode(), resp)
		err = ecode.Int(resp.Code)
		return
	}
	log.Info("media: set score success(media_id: %d, mid: %d, oid: %d, score: %d)", mediaID, mid, aid, score)
	return
}

// DelScore get media score.
func (d *Dao) DelScore(c context.Context, aid, mediaID, mid int64) (err error) {
	if mediaID == 0 || mid == 0 || aid == 0 {
		return
	}
	params := url.Values{}
	params.Set("media_id", strconv.FormatInt(mediaID, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("from", "1")
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	err = d.httpClient.Post(c, _delScoreURL, "", params, &resp)
	if err != nil {
		PromError("media:番剧删除评分接口")
		log.Error("media: d.client.Post(%s, data:%s) error(%+v)", _delScoreURL, params.Encode(), err)
		return
	}
	if resp.Code != 0 {
		PromError("media:番剧删除评分接口")
		log.Error("media: d.client.Post(%s, data:%s) res: %+v", _delScoreURL, params.Encode(), resp)
		err = ecode.Int(resp.Code)
	}
	log.Info("media: del score success(media_id: %d, mid: %d, oid: %d)", mediaID, mid, aid)
	return
}
