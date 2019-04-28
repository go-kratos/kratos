package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/player/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const _viewPointsURI = "/x/internal/creative/video/viewpoints"

// ViewPoints get view points data from creative.
func (d *Dao) ViewPoints(c context.Context, aid, cid int64) (points []*model.Points, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("cid", strconv.FormatInt(cid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Points []*model.Points `json:"points"`
		} `json:"data"`
	}
	ip := metadata.String(c, metadata.RemoteIP)
	if err = d.client.Get(c, d.viewPointsURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.viewPointsURL+"?"+params.Encode())
		return
	}
	points = res.Data.Points
	return
}
