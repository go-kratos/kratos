package geetest

import (
	"context"
	"crypto/tls"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/geetest"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	_register = "/register.php"
	_validate = "/validate.php"
)

// Dao is account dao.
type Dao struct {
	c *conf.Config
	// url
	registerURI string
	validateURI string
	// http client
	client  *http.Client
	clientx *httpx.Client
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		registerURI: c.Host.Geetest + _register,
		validateURI: c.Host.Geetest + _validate,
		// http client
		client:  NewClient(c.HTTPClient),
		clientx: httpx.NewClient(c.HTTPClient.Slow),
	}
	return
}

// PreProcess preprocessing the geetest and get to challenge
func (d *Dao) PreProcess(c context.Context, mid int64, ip, clientType string, newCaptcha int) (challenge string, err error) {
	var (
		req    *http.Request
		res    *http.Response
		bs     []byte
		params url.Values
	)
	params = url.Values{}
	params.Set("user_id", strconv.FormatInt(mid, 10))
	params.Set("new_captcha", strconv.Itoa(newCaptcha))
	params.Set("client_type", clientType)
	params.Set("ip_address", ip)
	params.Set("gt", d.c.Geetest.CaptchaID)
	if req, err = http.NewRequest("GET", d.registerURI+"?"+params.Encode(), nil); err != nil {
		log.Error("d.preprocess uri(%s) params(%s) error(%v)", d.registerURI, params.Encode(), err)
		err = ecode.CreativeGeetestAPIErr
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if res, err = d.client.Do(req); err != nil {
		log.Error("client.Do(%s) error(%v)", d.registerURI+"?"+params.Encode(), err)
		err = ecode.CreativeGeetestAPIErr
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusInternalServerError {
		log.Error("gtServerErr uri(%s) error(%v)", d.registerURI+"?"+params.Encode(), err)
		err = ecode.CreativeGeetestAPIErr
		return
	}
	if bs, err = ioutil.ReadAll(res.Body); err != nil {
		log.Error("ioutil.ReadAll(%s) uri(%s) error(%v)", bs, d.registerURI+"?"+params.Encode(), err)
		return
	}
	if len(bs) != 32 {
		log.Error("d.preprocess len(%s) the length not equate 32byte", string(bs))
		return
	}
	challenge = string(bs)
	return
}

// Validate recheck the challenge code and get to seccode
func (d *Dao) Validate(c context.Context, challenge, seccode, clientType, ip, captchaID string, mid int64) (res *geetest.ValidateRes, err error) {
	params := url.Values{}
	params.Set("seccode", seccode)
	params.Set("challenge", challenge)
	params.Set("captchaid", captchaID)
	params.Set("client_type", clientType)
	params.Set("ip_address", ip)
	params.Set("json_format", "1")
	params.Set("sdk", "golang_3.0.0")
	params.Set("user_id", strconv.FormatInt(mid, 10))
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	req, err := http.NewRequest("POST", d.validateURI, strings.NewReader(params.Encode()))
	if err != nil {
		log.Error("http.NewRequest error(%v) | uri(%s) params(%s)", err, d.validateURI, params.Encode())
		err = ecode.CreativeGeetestAPIErr
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err = d.clientx.Do(c, req, &res); err != nil {
		log.Error("d.client.Do error(%v)", err)
		err = ecode.CreativeGeetestAPIErr
		return
	}
	return
}

// NewClient new a http client.
func NewClient(c *conf.HTTPClient) (client *http.Client) {
	var (
		transport *http.Transport
		dialer    *net.Dialer
	)
	dialer = &net.Dialer{
		Timeout:   time.Duration(c.Slow.Timeout),
		KeepAlive: time.Duration(c.Slow.KeepAlive),
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
