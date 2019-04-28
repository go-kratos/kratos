package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const _delScoreURL = "http://api.bilibili.co/pgc/internal/review/score/delete"

// DelScore .
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
