package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	xhttp "net/http"
	"net/url"
	"strconv"

	"go-common/app/tool/saga/model"

	"github.com/pkg/errors"
)

const (
	qyWechatURL  = "https://qyapi.weixin.qq.com"
	corpID       = "wx0833ac9926284fa5" // 企业微信：Bilibili的企业ID
	departmentID = "12"                 // 公司统一用部门ID
)

// WechatPushMsg wechat push text message to specified user
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
		return
	}

	if res.ErrCode != 0 || res.InvalidUser != "" || res.InvalidParty != "" || res.InvalidTag != "" {
		invalidUser = res.InvalidUser
		err = errors.Errorf("WechatPushMsg: errcode: %d, errmsg: %s, invalidUser: %s, invalidParty: %s, invalidTag: %s",
			res.ErrCode, res.ErrMsg, res.InvalidUser, res.InvalidParty, res.InvalidTag)
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

// WechatAccessToken query access token with the specified secret
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

// WechatSagaVisible get all the user ids who can visiable saga
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

// RepoFiles ...TODO 该方法放在gitlab里比较好
func (d *Dao) RepoFiles(c context.Context, Host string, token string, repo *model.RepoInfo) (files []string, err error) {
	var (
		u   string
		req *xhttp.Request
	)
	u = fmt.Sprintf("http://%s/%s/%s/files/%s?format=json", Host, repo.Group, repo.Name, repo.Branch)
	if req, err = d.newRequest("GET", u, nil); err != nil {
		return
	}
	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if err = d.httpClient.Do(c, req, &files); err != nil {
		return
	}
	return
}
