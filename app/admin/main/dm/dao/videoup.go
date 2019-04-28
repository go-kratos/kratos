package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/dm/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_views  = "/videoup/views"
	_types  = "/videoup/types"
	_season = "/pgc/admin/season/dm/aids"
)

// TypeInfo TypeInfo
func (d *Dao) TypeInfo(c context.Context) (types map[int64]*model.ArchiveType, err error) {
	var (
		res struct {
			Code    int64                        `json:"code"`
			Data    map[int64]*model.ArchiveType `json:"data"`
			Message string                       `json:"message"`
		}
	)
	v := make(url.Values)
	if err = d.httpCli.Get(c, d.typesURI, "", v, &res); err != nil {
		log.Error("d.httpCli.Get(%s) error(%v)", d.typesURI, err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("%v", res)
		log.Error("d.httpClient.Get(%s) code(%d)", d.typesURI, res.Code)
	}
	types = res.Data
	return
}

// ArchiveVideos return archive and video info.
func (d *Dao) ArchiveVideos(c context.Context, aids []int64) (avm map[int64]*model.ArcVideo, err error) {
	var (
		res struct {
			Code    int64                     `json:"code"`
			Data    map[int64]*model.ArcVideo `json:"data"`
			Message string                    `json:"message"`
		}
		v = make(url.Values)
	)
	v.Set("aids", xstr.JoinInts(aids))
	if err = d.httpCli.Get(c, d.viewsURI, "", v, &res); err != nil {
		log.Error("d.httpClient.Get(%s) error(%v)", d.viewsURI, err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("%v", res)
		log.Error("d.httpClient.Get(%s) code(%d)", d.viewsURI, res.Code)
		return
	}
	avm = res.Data
	return
}

// SeasonInfos return season infos
func (d *Dao) SeasonInfos(c context.Context, IDType string, id int64) (aids, oids []int64, err error) {
	var (
		res struct {
			Code    int64               `json:"code"`
			Message string              `json:"message"`
			Data    []*model.SeasonInfo `json:"result"`
		}
		params = make(url.Values, 1)
	)
	switch IDType {
	case "ep":
		params.Set("epid", fmt.Sprint(id))
	case "ss":
		params.Set("season_id", fmt.Sprint(id))
	default:
		err = fmt.Errorf("season type(%s) error", IDType)
		log.Error("d.SeasonInfos error(%v)", err)
		return
	}
	if err = d.httpSearch.Get(c, d.seasonURI, "", params, &res); err != nil {
		log.Error("d.httpSearch.Get(uri:%s,params:%s) error(%v)", d.seasonURI, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("uri:%s,code:%d", d.seasonURI, res.Code)
		log.Error("d.SeasonInfos error(%v)", err)
		return
	}
	for _, v := range res.Data {
		aids = append(aids, v.Aid)
		oids = append(oids, v.Cid)
	}
	return
}
