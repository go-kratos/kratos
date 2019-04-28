package geetest

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_register = "http://api.geetest.com/register.php"
	_validate = "http://api.geetest.com/validate.php"
)

// Dao is account dao.
type Dao struct {
	c *conf.Config
	// url
	registerURI string
	validateURI string
	// http client
	client  *http.Client
	clientx *bm.Client
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		registerURI: _register,
		validateURI: _validate,
		// http client
		client:  NewClient(c.HTTPClient),
		clientx: bm.NewClient(c.HTTPClient.Slow),
	}
	return
}

// PreProcess preprocessing the geetest and get to challenge
func (d *Dao) PreProcess(c context.Context, mid int64, geeType string, gc conf.GeetestConfig, newCaptcha int) (challenge string, err error) {
	var (
		req    *http.Request
		res    *http.Response
		bs     []byte
		params url.Values
	)
	params = url.Values{}
	params.Set("user_id", strconv.FormatInt(mid, 10))
	params.Set("new_captcha", strconv.Itoa(newCaptcha))
	params.Set("ip_address", metadata.String(c, metadata.RemoteIP))
	params.Set("client_type", geeType)
	params.Set("gt", gc.CaptchaID)

	if req, err = http.NewRequest("GET", d.registerURI+"?"+params.Encode(), nil); err != nil {
		log.Error("d.preprocess uri(%s) params(%s) error(%v)", d.registerURI, params.Encode(), err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if res, err = d.client.Do(req); err != nil {
		log.Error("client.Do(%s) error(%v)", d.registerURI+"?"+params.Encode(), err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusInternalServerError {
		err = errors.New("http status code 5xx")
		log.Error("gtServerErr uri(%s) error(%v)", d.validateURI+"?"+params.Encode(), err)
		return
	}
	if bs, err = ioutil.ReadAll(res.Body); err != nil {
		log.Error("ioutil.ReadAll(%s) uri(%s) error(%v)", bs, d.validateURI+"?"+params.Encode(), err)
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
func (d *Dao) Validate(c context.Context, challenge, seccode, clientType, captchaID string, mid int64) (res *model.ValidateRes, err error) {
	params := url.Values{}
	params.Set("seccode", seccode)
	params.Set("challenge", challenge)
	params.Set("captchaid", captchaID)
	params.Set("client_type", clientType)
	params.Set("ip_address", metadata.String(c, metadata.RemoteIP))
	params.Set("json_format", "1")
	params.Set("sdk", "golang_3.0.0")
	params.Set("user_id", strconv.FormatInt(mid, 10))
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	req, err := http.NewRequest("POST", d.validateURI, strings.NewReader(params.Encode()))
	if err != nil {
		log.Error("http.NewRequest error(%v) | uri(%s) params(%s)", err, d.validateURI, params.Encode())
		return
	}
	log.Info("Validate(%v) start", params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res = new(model.ValidateRes)
	if err = d.clientx.Do(c, req, &res); err != nil {
		log.Error("d.client.Do error(%v) | uri(%s) params(%s) res(%v)", err, d.validateURI, params.Encode(), res)
		return
	}
	log.Info("Validate(%v) suc", res)
	return
}

// NewClient new a http client.
func NewClient(c *conf.HTTPClient) (client *http.Client) {
	var (
		transport *http.Transport
		dialer    *net.Dialer
	)
	dialer = &net.Dialer{
		Timeout:   time.Duration(c.Slow.Dial),
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

// GeeConfig get geetest config.
func (d *Dao) GeeConfig(ct string, c *conf.Geetest) (gc conf.GeetestConfig, clientType string) {
	if ct == "" {
		return c.PC, "web"
	}
	switch ct {
	case "web":
		gc = c.PC
		clientType = "web"
	case "h5":
		gc = c.H5
		clientType = "h5"
	default:
		gc = c.PC
		clientType = "web"
	}
	return
}
