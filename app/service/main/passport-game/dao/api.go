package dao

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/service/main/passport-game/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// MyInfo get info from api.
func (d *Dao) MyInfo(c context.Context, accessKey string) (accountInfo *model.Info, err error) {
	params := url.Values{}
	params.Set("access_key", accessKey)
	params.Set("type", "json")
	var res struct {
		Code int `json:"code"`
		model.Info
	}
	if err = d.client.Get(c, d.myInfoURI, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("myInfo url(%s) error(%v)", d.myInfoURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("myInfo url(%s) error(%v)", d.myInfoURI+"?"+params.Encode(), err)
		return
	}
	accountInfo = &res.Info
	return
}

// Oauth oauth via passport api.
func (d *Dao) Oauth(c context.Context, uri, accessKey, from string) (token *model.Token, err error) {
	params := url.Values{}
	params.Set("access_key", accessKey)
	params.Set("from", from)
	var res struct {
		Code  int `json:"code"`
		Token *struct {
			Mid         string `json:"mid"`
			AppID       int32  `json:"appid"`
			AccessToken string `json:"access_key"`
			CreateAt    int64  `json:"create_at"`
			UserID      string `json:"userid"`
			Uname       string `json:"uname"`
			Expires     string `json:"expires"`
			Permission  string `json:"permission"`
		} `json:"access_info,omitempty"`
		Data *model.Token `json:"data,omitempty"`
	}
	if err = d.client.Get(c, uri, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("model oauth url(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("model oauth url(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Token != nil {
		t := res.Token
		var mid int64
		if mid, err = strconv.ParseInt(t.Mid, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s, 10, 64) error(%v)", t.Mid, err)
			return
		}
		var expires int64
		if expires, err = strconv.ParseInt(t.Expires, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s, 10, 64) error(%v)", t.Expires, err)
			return
		}
		token = &model.Token{
			Mid:         mid,
			AppID:       t.AppID,
			AccessToken: t.AccessToken,
			CreateAt:    t.CreateAt,
			UserID:      t.UserID,
			Uname:       t.Uname,
			Expires:     expires,
			Permission:  t.Permission,
		}
	} else {
		token = res.Data
	}
	return
}

// Login login via model api.
func (d *Dao) Login(c context.Context, query, cookie string) (loginToken *model.LoginToken, err error) {
	req, err := http.NewRequest("GET", d.loginURI+"?"+query, nil)
	if err != nil {
		log.Error("http.NewRequest(GET, %s) error(%v)", d.loginURI+"?"+query, err)
		return
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("X-BACKEND-BILI-REAL-IP", metadata.String(c, metadata.RemoteIP))
	var res struct {
		Code      int    `json:"code"`
		Mid       int64  `json:"mid"`
		AccessKey string `json:"access_key"`
		Expires   int64  `json:"expires"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("model login url(%s) error(%v)", d.loginURI+"?"+query, err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("model login url(%s) error(%v)", d.loginURI+"?"+query, err)
		return
	}
	loginToken = &model.LoginToken{
		Mid:       res.Mid,
		AccessKey: res.AccessKey,
		Expires:   res.Expires,
	}
	return
}

// LoginOrigin login via passport api.
func (d *Dao) LoginOrigin(c context.Context, userid, rsaPwd string) (loginToken *model.LoginToken, err error) {
	params := url.Values{}
	params.Set("userid", userid)
	params.Set("pwd", rsaPwd)
	var res struct {
		Code      int    `json:"code"`
		Mid       int64  `json:"mid"`
		AccessKey string `json:"access_key"`
		Expires   int64  `json:"expires"`
	}
	if err = d.client.Get(c, d.loginURI, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("login url(%s) error(%v)", d.loginURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("login url(%s) error(%v)", d.loginURI+"?"+params.Encode(), err)
		return
	}

	str, _ := json.Marshal(res)
	log.Info("login url(%s) res(%s)", d.loginURI+"?"+params.Encode(), str)

	loginToken = &model.LoginToken{
		Mid:       res.Mid,
		AccessKey: res.AccessKey,
		Expires:   res.Expires,
	}
	return
}

// RSAKeyOrigin get rsa pub key and ts hash via passport api.
func (d *Dao) RSAKeyOrigin(c context.Context) (key *model.RSAKey, err error) {
	var res struct {
		*model.RSAKey
		Code int `json:"code"`
	}
	params := url.Values{}
	if err = d.client.Get(c, d.getKeyURI, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("key url(%s) error(%v)", d.getKeyURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("key url(%s) error(%v)", d.getKeyURI+"?"+params.Encode(), err)
		return
	}
	key = res.RSAKey
	return
}

// RenewToken renew token via passport api.
func (d *Dao) RenewToken(c context.Context, uri, ak, from string) (renewToken *model.RenewToken, err error) {
	params := url.Values{}
	params.Set("access_key", ak)
	params.Set("from", from)
	var res struct {
		Code    int   `json:"code"`
		Expires int64 `json:"expires"`
		Data    struct {
			Expires int64 `json:"expires"`
		}
	}
	if err = d.client.Get(c, uri, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("renewtoken url(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("renewtoken url(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	expires := res.Expires
	if expires == 0 {
		expires = res.Data.Expires
	}
	renewToken = &model.RenewToken{
		Expires: expires,
	}
	return
}

// RegV3 RegV3
func (d *Dao) RegV3(c context.Context, tdoRegV3 model.TdoRegV3) (regV3 *model.ResRegV3, err error) {
	p := url.Values{}
	p.Add("userpwd", tdoRegV3.Arg.Pwd)
	p.Add("userid", tdoRegV3.Arg.User)
	p.Add("captcha", tdoRegV3.Arg.Captcha)
	p.Add("ctoken", tdoRegV3.Arg.Ctoken)
	req, err := d.client.NewRequest(http.MethodPost, d.regV3URI, tdoRegV3.IP, p)
	if err != nil {
		log.Error("client.NewRequest(GET, %s) error(%v)", req.URL.String(), err)
		return
	}

	req.Header.Set("X-Forwarded-For", tdoRegV3.IP)
	req.Header.Set("Cookie", tdoRegV3.Cookie)
	var response struct {
		Code int `json:"code"`
		Mid  int `json:"mid"`
	}
	if err = d.client.Do(c, req, &response); err != nil {
		log.Error("client.Do(%s) error(%v)", req.URL.String(), err)
		return
	}
	if response.Code != ecode.OK.Code() {
		log.Warn("regv3 url(%s) code(%d)", req.URL.String(), response.Code)
		err = ecode.Int(response.Code)
		return
	}
	regV3 = new(model.ResRegV3)
	regV3.Mid = response.Mid
	return
}

// RegV2 RegV2
func (d *Dao) RegV2(c context.Context, tdoRegV2 model.TdoRegV2) (regV2 *model.ResRegV2, err error) {
	p := url.Values{}
	p.Add("captcha", tdoRegV2.Arg.Captcha)
	p.Add("ctoken", tdoRegV2.Arg.Ctoken)
	req, err := d.client.NewRequest(http.MethodPost, d.regV2URI, tdoRegV2.IP, p)
	if err != nil {
		log.Error("client.NewRequest(GET, %s) error(%v)", req.URL.String(), err)
		return
	}

	req.Header.Set("X-Forwarded-For", tdoRegV2.IP)
	req.Header.Set("Cookie", tdoRegV2.Cookie)
	var response struct {
		Code      int    `json:"code"`
		Mid       int    `json:"mid"`
		AccessKey string `json:"access_key"`
	}
	if err = d.client.Do(c, req, &response); err != nil {
		log.Error("client.Do(%s) error(%v)", req.URL.String(), err)
		return
	}
	if response.Code != ecode.OK.Code() {
		log.Warn("reg v2 url(%s) code(%d)", req.URL.String(), response.Code)
		err = ecode.Int(response.Code)
		return
	}
	regV2 = new(model.ResRegV2)
	regV2.Mid = response.Mid
	regV2.AccessKey = response.AccessKey
	return
}

// Reg Reg
func (d *Dao) Reg(c context.Context, tdoReg model.TdoReg) (reg *model.ResReg, err error) {
	p := url.Values{}
	p.Add("userpwd", tdoReg.Arg.Userpwd)
	p.Add("user", tdoReg.Arg.User)
	p.Add("email", tdoReg.Arg.Email)

	// new request
	req, err := d.client.NewRequest(http.MethodPost, d.regURI, tdoReg.IP, p)
	if err != nil {
		log.Error("client.NewRequest(GET, %s) error(%v)", req.URL.String(), err)
		return
	}

	req.Header.Set("X-Forwarded-For", tdoReg.IP)
	req.Header.Set("Cookie", tdoReg.Cookie)
	var response struct {
		Code int `json:"code"`
		Mid  int `json:"mid"`
	}
	if err = d.client.Do(c, req, &response); err != nil {
		log.Error("client.Do(%s) error(%v)", req.URL.String(), err)
		return
	}
	if response.Code != ecode.OK.Code() {
		log.Warn("reg url(%s) code(%d)", req.URL.String(), response.Code)
		err = ecode.Int(response.Code)
		return
	}
	reg = new(model.ResReg)
	reg.Mid = response.Mid
	return
}

// ByTel ByTel
func (d *Dao) ByTel(c context.Context, tdoByTel model.TdoByTel) (byTel *model.ResByTel, err error) {

	p := url.Values{}
	p.Add("userpwd", tdoByTel.Arg.Userpwd)
	p.Add("tel", tdoByTel.Arg.Tel)
	p.Add("captcha", tdoByTel.Arg.Captcha)
	p.Add("country_id", tdoByTel.Arg.CountryID)
	p.Add("uname", tdoByTel.Arg.Uname)

	// new request
	req, err := d.client.NewRequest(http.MethodPost, d.byTelURI, tdoByTel.IP, p)
	if err != nil {
		log.Error("client.NewRequest(GET, %s) error(%v)", req.URL.String(), err)
		return
	}

	req.Header.Set("X-Forwarded-For", tdoByTel.IP)
	req.Header.Set("Cookie", tdoByTel.Cookie)
	var response struct {
		Code int `json:"code"`
		Data struct {
			Mid       int    `json:"mid"`
			AccessKey string `json:"access_key"`
		}
	}

	if err = d.client.Do(c, req, &response); err != nil {
		log.Error("client.Do(%s) error(%v)", req.URL.String(), err)
		return
	}
	if response.Code != ecode.OK.Code() {
		log.Warn("byTel url(%s) code(%d)", req.URL.String(), response.Code)
		err = ecode.Int(response.Code)
		return
	}
	byTel = new(model.ResByTel)
	byTel.Mid = response.Data.Mid
	byTel.AccessKey = response.Data.AccessKey
	return
}

// Captcha Captcha
func (d *Dao) Captcha(c context.Context, ip string) (captchaData *model.CaptchaData, err error) {
	p := url.Values{}
	req, err := d.client.NewRequest(http.MethodGet, d.captchaURI, ip, p)
	if err != nil {
		log.Error("client.NewRequest(GET, %s) error(%v)", req.URL.String(), err)
		return
	}
	resCaptcha := new(model.ResCaptcha)

	req.Header.Set("X-Forwarded-For", ip)
	if err = d.client.Do(c, req, &resCaptcha); err != nil {
		log.Error("client.Do(%s) error(%v)", req.URL.String(), err)
		return
	}

	err = ecode.Int(resCaptcha.Code)
	captchaData = &resCaptcha.Data
	return
}

// SendSms SendSms
func (d *Dao) SendSms(c context.Context, tdoSendSms model.TdoSendSms) (err error) {
	p := url.Values{}
	p.Add("captcha", tdoSendSms.Arg.Captcha)
	p.Add("tel", tdoSendSms.Arg.Tel)
	p.Add("country_id", tdoSendSms.Arg.CountryID)
	p.Add("ctoken", tdoSendSms.Arg.Ctoken)
	p.Add("reset_pwd", strconv.FormatBool(tdoSendSms.Arg.ResetPwd))
	req, err := d.client.NewRequest(http.MethodGet, d.sendSmsURI, tdoSendSms.IP, p)
	if err != nil {
		log.Error("client.NewRequest(GET, %s) error(%v)", req.URL.String(), err)
		return
	}

	req.Header.Set("X-Forwarded-For", tdoSendSms.IP)
	req.Header.Set("Cookie", tdoSendSms.Cookie)
	resCaptcha := new(model.ResCaptcha)
	if err = d.client.Do(c, req, &resCaptcha); err != nil {
		log.Error("client.Do(%s) error(%v)", req.URL.String(), err)
		return
	}

	err = ecode.Int(resCaptcha.Code)
	return
}
