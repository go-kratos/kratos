package live

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_liveURL   = "/gift/simGift"
	_addVipURL = "/user/v0/Vip/addVip"
)

// Dao is live dao
type Dao struct {
	client    *httpx.Client
	liveURL   string
	addVipURL string
}

// New live dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:    httpx.NewClient(c.HTTPClient),
		liveURL:   c.Host.Live + _liveURL,
		addVipURL: c.Host.APILive + _addVipURL,
	}
	return
}

// Pack
func (d *Dao) Pack(c context.Context, mid int64, cardType int) (err error) {
	var res struct {
		Code int `json:"code"`
	}
	upack := map[string]map[string]interface{}{
		"header": map[string]interface{}{},
		"body": map[string]interface{}{
			"uid":       mid,
			"card_type": cardType,
		},
	}
	bytesData, err := json.Marshal(upack)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return
	}
	req, err := http.NewRequest("POST", d.liveURL, bytes.NewReader(bytesData))
	if err != nil {
		log.Error("http.NewRequest error(%v)", err)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("X-BACKEND-BILI-REAL-IP", "")
	log.Info("unicom pack mid(%d) card_type(%d)", mid, cardType)
	if err = d.client.Do(c, req, &res); err != nil || res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("unicom pack d.client.Do(%s) mid(%d) card_type(%d) error(%v)", d.liveURL, mid, cardType, err)
		return
	}
	return
}

// AddVip add live vip
func (d *Dao) AddVip(c context.Context, mid int64, day int) (msg string, err error) {
	params := url.Values{}
	params.Set("vip_type", "1")
	params.Set("day", strconv.Itoa(day))
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("platform", "main")
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}
	if err = d.client.Post(c, d.addVipURL, "", params, &res); err != nil {
		log.Error("live add vip url(%v) error(%v)", d.addVipURL+"?"+params.Encode(), err)
		return
	}
	msg = res.Msg
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("live add vip url(%v) res code(%d)", d.addVipURL+"?"+params.Encode(), res.Code)
		return
	}
	return
}
