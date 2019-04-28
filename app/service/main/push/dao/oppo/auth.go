package oppo

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Auth oppo auth token.
type Auth struct {
	Token  string
	Expire int64
}

type authResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token      string `json:"auth_token"`
		CreateTime int64  `json:"create_time"` // auth token的授权时间，单位为毫秒
	} `json:"data"`
}

// NewAuth get auth token.
func NewAuth(key, secret string) (a *Auth, err error) {
	tm := strconv.FormatInt(time.Now().UnixNano()/1000000, 10) // 用毫秒
	params := url.Values{}
	params.Add("app_key", key)
	params.Add("timestamp", tm)
	params.Add("sign", sign(key, secret, tm))
	res, err := http.PostForm(_apiAuth, params)
	if err != nil {
		return
	}
	defer res.Body.Close()
	dc := json.NewDecoder(res.Body)
	resp := new(authResponse)
	if err = dc.Decode(resp); err != nil {
		return
	}
	if resp.Code == ResponseCodeSuccess {
		a = &Auth{
			Token:  resp.Data.Token,
			Expire: resp.Data.CreateTime/1000 + _authExpire,
		}
		return
	}
	err = fmt.Errorf("new access error, code(%d) description(%s)", resp.Code, resp.Message)
	return
}

func sign(key, secret, timestamp string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(key+timestamp+secret)))
}

// IsExpired judge that whether privilige expired.
func (a *Auth) IsExpired() bool {
	return a.Expire <= time.Now().Add(4*time.Hour).Unix() // 提前4小时过期，for renew auth
}
