package dao

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"go-common/app/service/live/xcaptcha/conf"
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

//ValidateRes 验证返回值
type ValidateRes struct {
	Seccode string `json:"seccode"`
}

// GeeClient ...
type GeeClient struct {
	client     *http.Client
	clientPost *http.Client
}

// NewGeeClient new httpClient
func NewGeeClient(c *conf.GeeTestConfig) (client *GeeClient) {
	client = &GeeClient{
		client:     NewHttpClient(c.Get.Timeout, c.Get.KeepAlive),
		clientPost: NewHttpClient(c.Post.Timeout, c.Post.KeepAlive),
	}
	return client
}

// Register preprocessing the geetest and get to challenge
func (d *Dao) Register(c context.Context, ip, clientType string, newCaptcha int, captchaID string) (challenge string, err error) {
	challenge = ""

	var (
		bs     []byte
		params url.Values
	)
	params = url.Values{}
	params.Set("new_captcha", strconv.Itoa(newCaptcha))
	params.Set("client_type", clientType)
	params.Set("ip_address", ip)
	params.Set("gt", captchaID)

	if bs, err = d.Get(c, _register, params); err != nil {
		return
	}

	if len(bs) != 32 {
		return
	}
	challenge = string(bs)
	return
}

// Validate recheck the challenge code and get to seccode
func (d *Dao) Validate(c context.Context, challenge string, seccode string, clientType string, ip string, mid int64, captchaID string) (res *ValidateRes, err error) {
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

	res = &ValidateRes{}
	if bs, err = d.Post(c, _validate, params); err != nil {
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
func (d *Dao) Do(c context.Context, req *http.Request, client *http.Client) (body []byte, err error) {
	var res *http.Response
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
func (d *Dao) Get(c context.Context, uri string, params url.Values) (body []byte, err error) {
	req, err := NewRequest("GET", uri, params)
	if err != nil {
		return
	}
	body, err = d.Do(c, req, d.geeClient.client)
	return
}

// Post client.Get send POST request
func (d *Dao) Post(c context.Context, uri string, params url.Values) (body []byte, err error) {
	req, err := NewRequest("POST", uri, params)
	if err != nil {
		return
	}
	body, err = d.Do(c, req, d.geeClient.clientPost)
	return
}

// NewHttpClient new a http client.
func NewHttpClient(timeout int64, keepAlive int64) (client *http.Client) {
	var (
		transport *http.Transport
		dialer    *net.Dialer
	)

	dialer = &net.Dialer{
		Timeout:   time.Duration(time.Duration(timeout) * time.Millisecond),
		KeepAlive: time.Duration(time.Duration(keepAlive) * time.Second),
	}
	transport = &http.Transport{
		DialContext:     dialer.DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: transport,
	}
	return
}
