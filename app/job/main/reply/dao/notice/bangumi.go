package notice

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

type bangumi struct {
	EpisodeID  int64  `json:"episode_id,string"`
	SeasonID   int64  `json:"season_id"`
	Title      string `json:"title"`
	IndexTitle string `json:"index_title"`
}

// Bangumi return link.
func (d *Dao) Bangumi(c context.Context, oid int64) (title, link string, epid int64, err error) {
	params := url.Values{}
	params.Set("aids", strconv.FormatInt(oid, 10))
	params.Set("platform", "reply")
	params.Set("build", "0")
	var res struct {
		Code   int                `json:"code"`
		Result map[int64]*bangumi `json:"result"`
	}
	if err = d.httpClient.Get(c, d.urlBangumi, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s?%s) error(%v)", d.urlBangumi, params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Result == nil {
		err = fmt.Errorf("url:%s?%s code:%d", d.urlBangumi, params.Encode(), res.Code)
		return
	}
	if r := res.Result[oid]; r != nil {
		epid = r.EpisodeID
	}
	return
}
