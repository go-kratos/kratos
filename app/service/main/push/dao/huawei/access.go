package huawei

// http://developer.huawei.com/consumer/cn/service/hms/catalog/huaweipush.html?page=hmssdk_huaweipush_api_reference_s1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	_accessTokenURL  = "https://login.vmall.com/oauth2/token"
	_grantType       = "client_credentials"
	_respCodeSuccess = 0
)

// Access huawei access token.
type Access struct {
	AppID  string
	Token  string
	Expire int64
}

type accessResponse struct {
	Token  string `json:"access_token"`
	Expire int64  `json:"expires_in"` // Access Token的有效期，以秒为单位
	Scope  string `json:"scope"`      // Access Token的访问范围，即用户实际授予的权限列表
	Code   int    `json:"error"`
	Desc   string `json:"error_description"`
}

// NewAccess get token.
func NewAccess(clientID, clientSecret string) (a *Access, err error) {
	params := url.Values{}
	params.Add("grant_type", _grantType)
	params.Add("client_id", clientID)
	params.Add("client_secret", clientSecret)
	res, err := http.PostForm(_accessTokenURL, params)
	if err != nil {
		return
	}
	defer res.Body.Close()
	dc := json.NewDecoder(res.Body)
	resp := new(accessResponse)
	if err = dc.Decode(resp); err != nil {
		return
	}
	if resp.Code == _respCodeSuccess {
		a = &Access{
			AppID:  clientID,
			Token:  resp.Token,
			Expire: time.Now().Unix() + resp.Expire,
		}
		return
	}
	err = fmt.Errorf("new access error, code(%d) description(%s)", resp.Code, resp.Desc)
	return
}

// IsExpired judge that whether privilige expired.
func (a *Access) IsExpired() bool {
	return a.Expire <= time.Now().Add(8*time.Hour).Unix() // 提前8小时过期，for renew auth
}
