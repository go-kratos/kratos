package game

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model"
	"go-common/app/interface/main/app-view/model/game"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_infoURL = "/game/info"
)

type Dao struct {
	client  *httpx.Client
	infoURL string
	key     string
	secret  string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:  httpx.NewClient(c.HTTPGame),
		infoURL: c.Host.Game + _infoURL,
		key:     c.HTTPGame.Key,
		secret:  c.HTTPGame.Secret,
	}
	return
}

func (d *Dao) Info(c context.Context, gameID int64, plat int8) (info *game.Info, err error) {
	var platType int
	if model.IsAndroid(plat) {
		platType = 1
	} else if model.IsIOS(plat) {
		platType = 2
	}
	if platType == 0 {
		return
	}
	var req *http.Request
	params := url.Values{}
	params.Set("appkey", d.key)
	params.Set("game_base_id", strconv.FormatInt(gameID, 10))
	params.Set("platform_type", strconv.Itoa(platType))
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	mh := md5.Sum([]byte(params.Encode() + d.secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	if req, err = d.client.NewRequest("GET", d.infoURL, "", params); err != nil {
		return
	}
	var res struct {
		Code int        `json:"code"`
		Data *game.Info `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.infoURL+"?"+params.Encode())
		return
	}
	info = res.Data
	return
}
