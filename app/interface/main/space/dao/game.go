package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_lastPlayGameURI  = "/user/games.mid"
	_appPlayedGameURI = "/game/recent/play"
	_platformAndroid  = "android"
	_platformIOS      = "ios"
	_platTypeAndroid  = 1
	_platTypeIOS      = 2
)

// LastPlayGame get last play game.
func (d *Dao) LastPlayGame(c context.Context, mid int64) (data []*model.Game, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int           `json:"code"`
		Data []*model.Game `json:"games"`
	}
	if err = d.httpR.Get(c, d.lastPlayGameURL, ip, params, &res); err != nil {
		log.Error("d.httpR.Get(%s,%d) error(%v)", d.lastPlayGameURL, mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s,%d) code error(%d)", d.lastPlayGameURL, mid, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Data
	return
}

// AppPlayedGame get app player games.
func (d *Dao) AppPlayedGame(c context.Context, mid int64, platform string, pn, ps int) (data []*model.AppGame, count int, err error) {
	var platformType int
	switch platform {
	case _platformAndroid:
		platformType = _platTypeAndroid
	case _platformIOS:
		platformType = _platTypeIOS
	}
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("platform_type", strconv.Itoa(platformType))
	params.Set("page_num", strconv.Itoa(pn))
	params.Set("page_size", strconv.Itoa(ps))
	params.Set("ts", strconv.FormatInt(time.Now().Unix()*1000, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			List       []*model.AppGame `json:"list"`
			TotalCount int              `json:"total_count"`
		}
	}
	if err = d.httpGame.Get(c, d.appPlayedGameURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("AppPlayedGame d.httpR.Get(%s,%d) error(%v)", d.appPlayedGameURL+params.Encode(), mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("AppPlayedGame d.httpR.Get(%s,%d) code error(%d)", d.appPlayedGameURL+params.Encode(), mid, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Data.List
	count = res.Data.TotalCount
	return
}
