package order

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/order"
	"go-common/library/ecode"
	"go-common/library/log"
)

// UpValidate fn
func (d *Dao) UpValidate(c context.Context, mid int64, ip string) (uv *order.UpValidate, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("appkey", conf.Conf.HTTPClient.UpMng.Key)
	params.Set("appsecret", conf.Conf.HTTPClient.UpMng.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	var (
		sign, _ = tool.Sign(params)
		res     struct {
			Status      string            `json:"status"`
			CurrentTime int64             `json:"current_time"`
			Result      *order.UpValidate `json:"result"`
		}
		_upValidateURL = d.upValidateURI + "?" + sign
	)
	log.Info("upValidate url(%s)", _upValidateURL)
	req, err := http.NewRequest("GET", _upValidateURL, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v); mid(%d), ip(%s)", _upValidateURL, err, mid, ip)
		err = ecode.CreativeOrderAPIErr
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do upValidate url(%s)|mid(%d)|ip(%s)|error(%v)", _upValidateURL, mid, ip, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	if res.Status != "success" {
		log.Error("upValidate url(%s)|mid(%d) res(%v)", _upValidateURL, mid, res)
		return
	}
	uv = res.Result
	return
}

// GrowAccountState 获取up主状态 type 类型 0 视频 2 专栏 3 素材.
func (d *Dao) GrowAccountState(c context.Context, mid int64, ty int) (result *order.UpValidate, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("type", strconv.Itoa(ty))
	var res struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    *order.UpValidate `json:"data"`
	}
	if err = d.client.Get(c, d.accountStateURI, "", params, &res); err != nil {
		log.Error("GrowAccountState url(%s) response(%v) error(%v)", d.accountStateURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	log.Info("GrowAccountState url(%s)", d.accountStateURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("GrowAccountState url(%s),code(%d) msg(%s) res(%v)", d.accountStateURI, res.Code, res.Message, res)
		err = ecode.Int(res.Code)
		return
	}
	result = res.Data
	return
}
