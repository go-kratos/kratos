package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/log"
)

// TagIds get tag ids from tag
func (dao *Dao) TagIds(c context.Context, aid int64) (ids []int64, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int
		Data []*struct {
			ID int64 `json:"tag_id"`
		}
	}
	if err = dao.httpClient.Get(c, dao.tagURI, "", params, &res); err != nil {
		log.Error("dao.Ding api(%s) fail,err(%v)", dao.tagURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("dao.Ding api(%s) fail Code(%v)", dao.tagURI+"?"+params.Encode(), res.Code)
		return
	}
	ids = make([]int64, 0, len(res.Data))
	for _, d := range res.Data {
		ids = append(ids, d.ID)
	}
	return
}
