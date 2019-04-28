package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
)

// RecommendAuthors .
func (d *Dao) RecommendAuthors(c context.Context, platform string, mobiApp string, device string, build int, clientIP string, userID int64, buvid string, recType string, serviceArea string, _rapagesizen int, mid int64) (res []*model.RecommendAuthor, err error) {
	params := url.Values{}
	params.Set("platform", platform)
	params.Set("mobi_app", mobiApp)
	params.Set("device", device)
	params.Set("clientip", clientIP)
	params.Set("buvid", buvid)
	params.Set("rec_type", recType)
	params.Set("service_area", serviceArea)
	params.Set("userid", strconv.FormatInt(userID, 10))
	params.Set("build", strconv.Itoa(build))
	params.Set("context_id", strconv.FormatInt(mid, 10))
	var r struct {
		Code int `json:"code"`
		Data []*model.RecommendAuthor
	}
	if err = d.httpClient.Get(c, d.c.Article.RecommendAuthorsURL, "", params, &r); err != nil {
		log.Error("activity: RecommendAuthors url(%s) error(%+v)", d.c.Article.RecommendAuthorsURL+"?"+params.Encode(), err)
		return
	}
	if r.Code != 0 {
		log.Error("activity: RecommendAuthors url(%s) res(%d) error(%+v)", d.c.Article.RecommendAuthorsURL+"?"+params.Encode(), r.Code, err)
		return
	}
	res = r.Data
	return
}
