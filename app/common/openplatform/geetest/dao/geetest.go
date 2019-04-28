package dao

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"go-common/app/common/openplatform/geetest/model"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	_register = "http://api.geetest.com/register.php"
	_validate = "http://api.geetest.com/validate.php"
)

// PreProcess preprocessing the geetest and get to challenge
func PreProcess(c context.Context, mid int64, ip, clientType string, newCaptcha int, captchaID string) (challenge string, err error) {
	var (
		bs     []byte
		params url.Values
	)
	params = url.Values{}
	params.Set("user_id", strconv.FormatInt(mid, 10))
	params.Set("new_captcha", strconv.Itoa(newCaptcha))
	params.Set("client_type", clientType)
	params.Set("ip_address", ip)
	params.Set("gt", captchaID)

	if bs, err = Get(c, _register, params); err != nil {
		return
	}

	if len(bs) != 32 {
		return
	}
	challenge = string(bs)
	return
}

// Validate recheck the challenge code and get to seccode
func Validate(c context.Context, challenge, seccode, clientType, ip, captchaID string, mid int64) (res *model.ValidateRes, err error) {
	var (
		bs     []byte
		params url.Values
	)
	params = url.Values{}
	params.Set("seccode", seccode)
	params.Set("challenge", challenge)
	params.Set("captchaid", captchaID)
	params.Set("client_type", clientType)
	params.Set("ip_address", ip)
	params.Set("json_format", "1")
	params.Set("sdk", "golang_3.0.0")
	params.Set("user_id", strconv.FormatInt(mid, 10))
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	if bs, err = Post(c, _validate, params); err != nil {
		return
	}
	if err = json.Unmarshal(bs, &res); err != nil {
		return
	}
	return
}

// NewRequest new a http request.
func NewRequest(method, uri string, params url.Values) (req *http.Request, err error) {
	if method == "GET" {
		req, err = http.NewRequest(method, uri+"?"+params.Encode(), nil)
	} else {
		req, err = http.NewRequest(method, uri, strings.NewReader(params.Encode()))
	}
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return
}

// Do handler request
func Do(c context.Context, req *http.Request) (body []byte, err error) {
	var res *http.Response
	dialer := &net.Dialer{
		Timeout:   time.Duration(1 * int64(time.Second)),
		KeepAlive: time.Duration(1 * int64(time.Second)),
	}
	transport := &http.Transport{
		DialContext:     dialer.DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
	}
	req = req.WithContext(c)
	if res, err = client.Do(req); err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusInternalServerError {
		err = errors.New("http status code 5xx")
		return
	}
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}
	return
}

// Get client.Get send GET request
func Get(c context.Context, uri string, params url.Values) (body []byte, err error) {
	req, err := NewRequest("GET", uri, params)
	if err != nil {
		return
	}
	body, err = Do(c, req)
	return
}

// Post client.Get send POST request
func Post(c context.Context, uri string, params url.Values) (body []byte, err error) {
	req, err := NewRequest("POST", uri, params)
	if err != nil {
		return
	}
	body, err = Do(c, req)
	return
}
