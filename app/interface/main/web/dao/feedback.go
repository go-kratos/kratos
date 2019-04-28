package dao

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const _platform = "web"

// Feedback feedback http request.
func (d *Dao) Feedback(c context.Context, feedParams *model.Feedback) (err error) {
	var (
		content []byte
		params  = url.Values{}
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	if feedParams.Aid > 0 {
		params.Set("aid", strconv.FormatInt(feedParams.Aid, 10))
	}
	params.Set("mid", strconv.FormatInt(feedParams.Mid, 10))
	params.Set("tag_id", strconv.FormatInt(feedParams.TagID, 10))
	if content, err = json.Marshal(feedParams.Content); err != nil {
		log.Error("content json.Marshal(%+v) error(%v)", feedParams.Content, err)
	} else {
		params.Set("content", string(content))
	}
	params.Set("browser", feedParams.Browser)
	params.Set("version", feedParams.Version)
	params.Set("buvid", feedParams.Buvid)
	params.Set("platform", _platform)
	if feedParams.Email != "" {
		params.Set("email", feedParams.Email)
	}
	if feedParams.QQ != "" {
		params.Set("qq", feedParams.QQ)
	}
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.feedbackURL, ip, params, &res); err != nil {
		log.Error("d.client.Post(%s) error(%v)", d.feedbackURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.client.Post(%s) code(%d)", d.feedbackURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	return
}
