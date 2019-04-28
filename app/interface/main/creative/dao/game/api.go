package game

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/game"
	"go-common/library/ecode"

	"go-common/library/log"
)

const (
	_gameListURI = "/game/list"
	_gameInfoURI = "/game/info"
)

// List fn
func (d *Dao) List(c context.Context, keywordStr, ip string) (gameList []*game.ListItem, err error) {
	params := url.Values{}
	params.Set("appkey", conf.Conf.Game.App.Key)
	params.Set("appsecret", conf.Conf.Game.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	params.Set("keyword", keywordStr)
	var res struct {
		Code int              `json:"code"`
		Data []*game.ListItem `json:"data"`
	}
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.gameListURL + "?" + query
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v), ip(%s)", url, err, ip)
		err = ecode.CreativeGameOpenAPIErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); ip(%s)", url, err, ip)
		err = ecode.CreativeGameOpenAPIErr
		return
	}
	log.Info("game list url(%+v)", url)
	if res.Code != 0 {
		log.Error("game open api url(%s) res(%v); ip(%s), code(%d)", url, res, ip, res.Code)
		err = ecode.CreativeGameOpenAPIErr
	}
	gameList = res.Data
	return
}

// Info fn
func (d *Dao) Info(c context.Context, gameBaseID int64, platForType int, ip string) (info *game.Info, err error) {
	params := url.Values{}
	params.Set("game_base_id", strconv.FormatInt(gameBaseID, 10))
	params.Set("platform_type", strconv.Itoa(platForType))
	params.Set("appkey", conf.Conf.Game.App.Key)
	params.Set("appsecret", conf.Conf.Game.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	var res struct {
		Code int        `json:"code"`
		Msg  string     `json:"message"`
		Data *game.Info `json:"data"`
	}
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.gameInfoURL + "?" + query
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v), ip(%s)", url, err, ip)
		err = ecode.CreativeGameOpenAPIErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); ip(%s)", url, err, ip)
		err = ecode.CreativeGameOpenAPIErr
		return
	}
	log.Info("game info url(%+v)", url)
	if res.Code != 0 {
		log.Error("game open api url(%s) res(%v); ip(%s), code(%d)", url, res, ip, res.Code)
		err = ecode.CreativeGameOpenAPIErr
	}
	info = res.Data
	return
}
