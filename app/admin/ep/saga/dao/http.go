package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	xhttp "net/http"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/admin/ep/saga/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	qyWechatURL  = "https://qyapi.weixin.qq.com"
	corpID       = "wx0833ac9926284fa5" // 企业微信：Bilibili的企业ID
	departmentID = "12"                 // 公司统一用部门ID
	//corpID       = "wwa24b497e5efdd78c" // 香满庭
	//departmentID = "2"                  // 香满庭
	_ajSessionID = "_AJSESSIONID"
)

// WechatPushMsg wechat push text message to specified user 发送企业微信信息
func (d *Dao) WechatPushMsg(c context.Context, token string, txtMsg *model.TxtNotification) (invalidUser string, err error) {
	var (
		u      string
		params = url.Values{}
		res    struct {
			ErrCode      int    `json:"errcode"`
			ErrMsg       string `json:"errmsg"`
			InvalidUser  string `json:"invaliduser"`
			InvalidParty string `json:"invalidparty"`
			InvalidTag   string `json:"invalidtag"`
		}
	)
	u = qyWechatURL + "/cgi-bin/message/send"
	params.Set("access_token", token)

	if err = d.PostJSON(c, u, "", params, &res, txtMsg); err != nil {
		log.Info("WechatPushMsg PostJSON err (%+v)", err)
		return
	}

	if res.ErrCode != 0 || res.InvalidUser != "" || res.InvalidParty != "" || res.InvalidTag != "" {
		invalidUser = res.InvalidUser
		err = errors.Errorf("WechatPushMsg: errcode: %d, errmsg: %s, invalidUser: %s, invalidParty: %s, invalidTag: %s",
			res.ErrCode, res.ErrMsg, res.InvalidUser, res.InvalidParty, res.InvalidTag)
		log.Info("WechatPushMsg err (%+v)", err)
		return
	}

	return
}

// PostJSON post http request with json params.
func (d *Dao) PostJSON(c context.Context, uri, ip string, params url.Values, res interface{}, v interface{}) (err error) {
	var (
		body = &bytes.Buffer{}
		req  *xhttp.Request
		url  string
		en   string
	)

	if err = json.NewEncoder(body).Encode(v); err != nil {
		return
	}

	url = uri
	if en = params.Encode(); en != "" {
		url += "?" + en
	}

	if req, err = xhttp.NewRequest(xhttp.MethodPost, url, body); err != nil {
		return
	}

	if err = d.httpClient.Do(c, req, &res); err != nil {
		return
	}
	return
}

// WechatAccessToken query access token with the specified secret 企业微信api获取公司token
func (d *Dao) WechatAccessToken(c context.Context, secret string) (token string, expire int32, err error) {
	var (
		u      string
		params = url.Values{}
		res    struct {
			ErrCode     int    `json:"errcode"`
			ErrMsg      string `json:"errmsg"`
			AccessToken string `json:"access_token"`
			ExpiresIn   int32  `json:"expires_in"`
		}
	)
	u = qyWechatURL + "/cgi-bin/gettoken"
	params.Set("corpid", corpID)
	params.Set("corpsecret", secret)

	if err = d.httpClient.Get(c, u, "", params, &res); err != nil {
		return
	}

	if res.ErrCode != 0 {
		err = errors.Errorf("WechatAccessToken: errcode: %d, errmsg: %s", res.ErrCode, res.ErrMsg)
		return
	}

	token = res.AccessToken
	expire = res.ExpiresIn

	return
}

// WechatContacts query all the contacts
func (d *Dao) WechatContacts(c context.Context, accessToken string) (contacts []*model.ContactInfo, err error) {
	var (
		u      string
		params = url.Values{}
		res    struct {
			ErrCode  int                  `json:"errcode"`
			ErrMsg   string               `json:"errmsg"`
			UserList []*model.ContactInfo `json:"userlist"`
		}
	)
	u = qyWechatURL + "/cgi-bin/user/list"
	params.Set("access_token", accessToken)
	params.Set("department_id", departmentID)
	params.Set("fetch_child", "1")

	if err = d.httpClient.Get(c, u, "", params, &res); err != nil {
		return
	}

	if res.ErrCode != 0 {
		err = errors.Errorf("WechatContacts: errcode: %d, errmsg: %s", res.ErrCode, res.ErrMsg)
		return
	}

	contacts = res.UserList

	return
}

// WechatSagaVisible get all the user ids who can visiable saga 获取用户ID列表
func (d *Dao) WechatSagaVisible(c context.Context, accessToken string, agentID int) (users []*model.UserInfo, err error) {
	var (
		u      string
		params = url.Values{}
		res    struct {
			ErrCode      int                 `json:"errcode"`
			ErrMsg       string              `json:"errmsg"`
			VisibleUsers model.AllowUserInfo `json:"allow_userinfos"`
		}
	)
	u = qyWechatURL + "/cgi-bin/agent/get"
	params.Set("access_token", accessToken)
	params.Set("agentid", strconv.Itoa(agentID))

	if err = d.httpClient.Get(c, u, "", params, &res); err != nil {
		return
	}

	if res.ErrCode != 0 {
		err = errors.Errorf("WechatSagaVisible: errcode: %d, errmsg: %s", res.ErrCode, res.ErrMsg)
		return
	}

	users = res.VisibleUsers.Users
	return
}

// WechatParams ...
func (d *Dao) WechatParams(c context.Context, u string, params url.Values, resp interface{}) (err error) {
	var (
		req *xhttp.Request
		en  string
	)
	if en = params.Encode(); en != "" {
		u += "?" + en
	}
	if req, err = xhttp.NewRequest(xhttp.MethodGet, u, nil); err != nil {
		return
	}

	return d.httpClient.Do(c, req, &resp)
}

// NewRequest ...
func (d *Dao) NewRequest(method, url string, v interface{}) (req *http.Request, err error) {
	body := &bytes.Buffer{}
	if method != http.MethodGet {
		if err = json.NewEncoder(body).Encode(v); err != nil {
			log.Error("json encode value(%s), error(%v) ", v, err)
			return
		}
	}
	if req, err = http.NewRequest(method, url, body); err != nil {
		log.Error("http new request url(%s), error(%v)", url, err)
	}
	return
}

// QueryAllConfigFile ...
func (d *Dao) QueryAllConfigFile(c context.Context, sessionID, url string) (resp *model.ConfigData, err error) {
	var (
		req       *http.Request
		respValue = &model.SvenResp{}
	)

	log.Info("QueryAllConfigFile: sessionID: %s, url: %s", sessionID, url)
	if req, err = d.NewRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	req.Header.Set("Cookie", _ajSessionID+"="+sessionID)

	if err = d.httpClient.Do(c, req, &respValue); err != nil {
		return
	}
	if respValue.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(respValue.Code), "QueryAllConfigFile failed, sessionID(%s), url(%s)", sessionID, url)
		return
	}
	resp = respValue.Data
	return
}

// QueryConfigFileContent ...
func (d *Dao) QueryConfigFileContent(c context.Context, sessionID, url string) (content string, err error) {
	var (
		req       *http.Request
		respValue = &model.ConfigValueResp{}
	)

	if req, err = d.NewRequest(http.MethodGet, url, nil); err != nil {
		return
	}

	req.Header.Set("Cookie", _ajSessionID+"="+sessionID)
	if err = d.httpClient.Do(c, req, respValue); err != nil {
		return
	}
	if respValue.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(respValue.Code), "QueryConfigFileContent failed, sessionID(%s), url(%s)", sessionID, url)
		return
	}
	if respValue.Data != nil {
		content = respValue.Data.Comment
	}

	return
}

// RequestConfig ...
func (d *Dao) RequestConfig(c context.Context, sessionID, reqUrl string, params url.Values) (resp *model.CommonResp, err error) {
	var req *http.Request

	if req, err = http.NewRequest("POST", reqUrl, strings.NewReader(params.Encode())); err != nil {
		log.Error("http.NewRequest error(%v) | uri(%s) params(%s)", err, reqUrl, params.Encode())
		return
	}
	log.Info("RequestConfig url: %v", req.URL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", _ajSessionID+"="+sessionID)
	if err = d.httpClient.Do(c, req, &resp); err != nil {
		log.Error("RequestConfig err：%+v", err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(resp.Code), "RequestConfig failed, sessionID(%s), url(%s)", sessionID, reqUrl)
		return
	}
	return
}
