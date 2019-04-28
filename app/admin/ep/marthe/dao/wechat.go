package dao

import (
	"context"
	"net/url"

	"go-common/app/admin/ep/marthe/model"
)

const (
	_qyWechatURL  = "https://qyapi.weixin.qq.com"
	_corpID       = "wx0833ac9926284fa5" // 企业微信：Bilibili的企业ID
	_departmentID = "12"                 // 公司统一用部门ID
	_corpsecret   = "WveODxk3xpT9box48wcxkmArx3mu6d4vJHdJkNy_iTk"

	_getToken = "/cgi-bin/gettoken"
	_userList = "/cgi-bin/user/list"
)

// WechatAccessToken query access token with the specified secret 企业微信api获取公司token
func (d *Dao) WechatAccessToken(c context.Context) (token string, err error) {
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
	u = _qyWechatURL + _getToken
	params.Set("corpid", _corpID)
	params.Set("corpsecret", _corpsecret)

	if err = d.httpClient.Get(c, u, "", params, &res); err != nil {
		return
	}

	if res.ErrCode != 0 {
		return
	}

	token = res.AccessToken
	return
}

// WechatContacts Wechat Contacts 获取用户信息列表
func (d *Dao) WechatContacts(c context.Context) (contacts []*model.WechatContact, err error) {
	var (
		token  string
		u      string
		params = url.Values{}
		res    struct {
			ErrCode  int                    `json:"errcode"`
			ErrMsg   string                 `json:"errmsg"`
			UserList []*model.WechatContact `json:"userlist"`
		}
	)
	//get token
	if token, err = d.WechatAccessToken(c); err != nil {
		return
	}

	u = _qyWechatURL + _userList
	params.Set("access_token", token)
	params.Set("department_id", _departmentID)
	params.Set("fetch_child", "1")

	if err = d.httpClient.Get(c, u, "", params, &res); err != nil {
		return
	}

	if res.ErrCode != 0 {
		return
	}

	contacts = res.UserList
	return
}
