package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/passport-sns/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_qqAuthorizeUrl   = "https://graph.qq.com/oauth2.0/authorize"
	_qqAccessTokenUrl = "https://graph.qq.com/oauth2.0/token"
	_qqOpenIDUrl      = "https://graph.qq.com/oauth2.0/me"

	_respCodeSuccess = 0
)

// QQAuthorize .
func (d *Dao) QQAuthorize(c context.Context, appID, redirectURL, display string) (url string) {
	scope := "do_like,get_user_info,get_simple_userinfo,get_vip_info,get_vip_rich_info,add_one_blog,list_album,upload_pic,add_album,list_photo,get_info,add_t,del_t,add_pic_t,get_repost_list,get_other_info,get_fanslist,get_idollist,add_idol,del_idol,get_tenpay_addr"
	displayParam := ""
	if display != "" {
		displayParam = "&display=" + display
	}
	return fmt.Sprintf(_qqAuthorizeUrl+"?response_type=code&state=authorize%s&client_id=%s&redirect_uri=%s&scope=%s", displayParam, appID, redirectURL, scope)
}

// QQOauth2Info .
func (d *Dao) QQOauth2Info(c context.Context, code, redirectUrl string, app *model.SnsApps) (res *model.Oauth2Info, err error) {
	accessResp, err := d.qqAccessToken(c, code, app.AppID, app.AppSecret, redirectUrl)
	if err != nil {
		return nil, err
	}
	openIdResp, err := d.qqOpenID(c, accessResp.Token, app.Business)
	if err != nil {
		return nil, err
	}
	// TODO 保证能获取到UnionID的情况可以考虑去掉
	if openIdResp.UnionID == "" {
		openIdResp.UnionID = openIdResp.OpenID
	}
	res = &model.Oauth2Info{
		UnionID: openIdResp.UnionID,
		OpenID:  openIdResp.OpenID,
		Token:   accessResp.Token,
		Refresh: accessResp.Refresh,
		Expires: accessResp.Expires,
	}
	return
}

// qqAccessToken .
func (d *Dao) qqAccessToken(c context.Context, code, appID, appSecret, redirectUrl string) (resp *model.QQAccessResp, err error) {
	var (
		res     *http.Response
		bs      []byte
		params  = url.Values{}
		value   = url.Values{}
		expires int64
	)
	params.Set("client_id", appID)
	params.Set("client_secret", appSecret)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("redirect_uri", redirectUrl)

	if res, err = d.client.Get(_qqAccessTokenUrl + "?" + params.Encode()); err != nil {
		log.Error("d.qqAccessToken error(%+v) code(%s) appID(%s)", err, code, appID)
		return nil, err
	}
	defer res.Body.Close()

	if bs, err = ioutil.ReadAll(res.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%+v) code(%s) appID(%s)", err, code, appID)
		return nil, err
	}
	respStr := string(bs)
	if strings.HasPrefix(respStr, "callback") {
		resp = new(model.QQAccessResp)
		start := strings.Index(respStr, "{")
		end := strings.Index(respStr, "}")
		respStr = respStr[start : end+1]
		if err = json.Unmarshal([]byte(respStr), resp); err != nil {
			return nil, err
		}
		log.Error("request qq token failed with code(%d) desc(%s)", resp.Code, resp.Description)
		return nil, ecode.PassportSnsRequestErr
	}
	value, err = url.ParseQuery(respStr)
	expires, err = strconv.ParseInt(value.Get("expires_in"), 10, 64)

	resp = &model.QQAccessResp{
		Token:   value.Get("access_token"),
		Refresh: value.Get("refresh_token"),
		Expires: time.Now().Unix() + expires,
	}
	return resp, nil
}

// qqOpenID .
func (d *Dao) qqOpenID(c context.Context, token string, business int) (resp *model.QQOpenIDResp, err error) {
	var (
		res    *http.Response
		bs     []byte
		params = url.Values{}
	)
	params.Set("access_token", token)
	params.Set("unionid", "1")
	// TODO 如果后续要支持没有unionid权限的appid，可以考虑在sns_apps表增加unionid权限标识的字段
	//if business == model.BusinessMall {
	//	params.Set("unionid", "1")
	//}

	if res, err = d.client.Get(_qqOpenIDUrl + "?" + params.Encode()); err != nil {
		log.Error("d.qqOpenID error(%+v) token(%d) business(%d)", err, token, business)
		return nil, err
	}
	defer res.Body.Close()

	if bs, err = ioutil.ReadAll(res.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%+v) token(%d) business(%d)", err, token, business)
		return nil, err
	}
	respStr := string(bs)
	if strings.HasPrefix(respStr, "callback") {
		start := strings.Index(respStr, "{")
		end := strings.Index(respStr, "}")
		respStr = respStr[start : end+1]
	}
	resp = new(model.QQOpenIDResp)
	if err = json.Unmarshal([]byte(respStr), resp); err != nil {
		return nil, err
	}
	if resp.Code == _respCodeSuccess {
		return resp, nil
	}
	log.Error("request qq openid failed with code(%d) desc(%s)", resp.Code, resp.Description)
	return nil, ecode.PassportSnsRequestErr
}
