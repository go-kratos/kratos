package blademaster

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	xhttp "net/http"

	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/net/netutil/breaker"
	"go-common/library/stat"
	xtime "go-common/library/time"

	"github.com/gogo/protobuf/proto"
	pkgerr "github.com/pkg/errors"
)

const (
	_minRead = 16 * 1024 // 16kb

	_appKey    = "appkey"
	_appSecret = "appsecret"
	_ts        = "ts"
)

var (
	_noKickUserAgent = "haoguanwei@bilibili.com "
	clientStats      = stat.HTTPClient
)

func init() {
	n, err := os.Hostname()
	if err == nil {
		_noKickUserAgent = _noKickUserAgent + runtime.Version() + " " + n
	}
}

// App bilibili intranet authorization.
type App struct {
	Key    string
	Secret string
}

// ClientConfig is http client conf.
type ClientConfig struct {
	*App
	Dial      xtime.Duration
	Timeout   xtime.Duration
	KeepAlive xtime.Duration
	Breaker   *breaker.Config
	URL       map[string]*ClientConfig
	Host      map[string]*ClientConfig
}

// Client is http client.
type Client struct {
	conf      *ClientConfig
	client    *xhttp.Client
	dialer    *net.Dialer
	transport xhttp.RoundTripper

	urlConf  map[string]*ClientConfig
	hostConf map[string]*ClientConfig
	mutex    sync.RWMutex
	breaker  *breaker.Group
}

// NewClient new a http client.
func NewClient(c *ClientConfig) *Client {
	client := new(Client)
	client.conf = c
	client.dialer = &net.Dialer{
		Timeout:   time.Duration(c.Dial),
		KeepAlive: time.Duration(c.KeepAlive),
	}

	originTransport := &xhttp.Transport{
		DialContext:     client.dialer.DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// wraps RoundTripper for tracer
	client.transport = &TraceTransport{RoundTripper: originTransport}
	client.client = &xhttp.Client{
		Transport: client.transport,
	}
	client.urlConf = make(map[string]*ClientConfig)
	client.hostConf = make(map[string]*ClientConfig)
	client.breaker = breaker.NewGroup(c.Breaker)
	// check appkey
	if c.Key == "" || c.Secret == "" {
		panic("http client must config appkey and appsecret")
	}
	if c.Timeout <= 0 {
		panic("must config http timeout!!!")
	}
	for uri, cfg := range c.URL {
		client.urlConf[uri] = cfg
	}
	for host, cfg := range c.Host {
		client.hostConf[host] = cfg
	}
	return client
}

// SetTransport set client transport
func (client *Client) SetTransport(t xhttp.RoundTripper) {
	client.transport = t
	client.client.Transport = t
}

// SetConfig set client config.
func (client *Client) SetConfig(c *ClientConfig) {
	client.mutex.Lock()
	if c.App != nil {
		client.conf.App.Key = c.App.Key
		client.conf.App.Secret = c.App.Secret
	}
	if c.Timeout > 0 {
		client.conf.Timeout = c.Timeout
	}
	if c.KeepAlive > 0 {
		client.dialer.KeepAlive = time.Duration(c.KeepAlive)
		client.conf.KeepAlive = c.KeepAlive
	}
	if c.Dial > 0 {
		client.dialer.Timeout = time.Duration(c.Dial)
		client.conf.Timeout = c.Dial
	}
	if c.Breaker != nil {
		client.conf.Breaker = c.Breaker
		client.breaker.Reload(c.Breaker)
	}
	for uri, cfg := range c.URL {
		client.urlConf[uri] = cfg
	}
	for host, cfg := range c.Host {
		client.hostConf[host] = cfg
	}
	client.mutex.Unlock()
}

// NewRequest new http request with method, uri, ip, values and headers.
// TODO(zhoujiahui): param realIP should be removed later.
func (client *Client) NewRequest(method, uri, realIP string, params url.Values) (req *xhttp.Request, err error) {
	enc, err := client.sign(params)
	if err != nil {
		err = pkgerr.Wrapf(err, "uri:%s,params:%v", uri, params)
		return
	}
	ru := uri
	if enc != "" {
		ru = uri + "?" + enc
	}
	if method == xhttp.MethodGet {
		req, err = xhttp.NewRequest(xhttp.MethodGet, ru, nil)
	} else {
		req, err = xhttp.NewRequest(xhttp.MethodPost, uri, strings.NewReader(enc))
	}
	if err != nil {
		err = pkgerr.Wrapf(err, "method:%s,uri:%s", method, ru)
		return
	}
	const (
		_contentType = "Content-Type"
		_urlencoded  = "application/x-www-form-urlencoded"
		_userAgent   = "User-Agent"
	)
	if method == xhttp.MethodPost {
		req.Header.Set(_contentType, _urlencoded)
	}
	if realIP != "" {
		req.Header.Set(_httpHeaderRemoteIP, realIP)
	}
	req.Header.Set(_userAgent, _noKickUserAgent+" "+env.AppID)
	return
}

// Get issues a GET to the specified URL.
func (client *Client) Get(c context.Context, uri, ip string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodGet, uri, ip, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// Post issues a Post to the specified URL.
func (client *Client) Post(c context.Context, uri, ip string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodPost, uri, ip, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// RESTfulGet issues a RESTful GET to the specified URL.
func (client *Client) RESTfulGet(c context.Context, uri, ip string, params url.Values, res interface{}, v ...interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodGet, fmt.Sprintf(uri, v...), ip, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res, uri)
}

// RESTfulPost issues a RESTful Post to the specified URL.
func (client *Client) RESTfulPost(c context.Context, uri, ip string, params url.Values, res interface{}, v ...interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodPost, fmt.Sprintf(uri, v...), ip, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res, uri)
}

// Raw sends an HTTP request and returns bytes response
func (client *Client) Raw(c context.Context, req *xhttp.Request, v ...string) (bs []byte, err error) {
	var (
		ok      bool
		code    string
		cancel  func()
		resp    *xhttp.Response
		config  *ClientConfig
		timeout time.Duration
		uri     = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.Path)
	)
	// NOTE fix prom & config uri key.
	if len(v) == 1 {
		uri = v[0]
	}
	// breaker
	brk := client.breaker.Get(uri)
	if err = brk.Allow(); err != nil {
		code = "breaker"
		clientStats.Incr(uri, code)
		return
	}
	defer client.onBreaker(brk, &err)
	// stat
	now := time.Now()
	defer func() {
		clientStats.Timing(uri, int64(time.Since(now)/time.Millisecond))
		if code != "" {
			clientStats.Incr(uri, code)
		}
	}()
	// get config
	// 1.url config 2.host config 3.default
	client.mutex.RLock()
	if config, ok = client.urlConf[uri]; !ok {
		if config, ok = client.hostConf[req.Host]; !ok {
			config = client.conf
		}
	}
	client.mutex.RUnlock()
	// timeout
	deliver := true
	timeout = time.Duration(config.Timeout)
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < timeout {
			// deliver small timeout
			timeout = ctimeout
			deliver = false
		}
	}
	if deliver {
		c, cancel = context.WithTimeout(c, timeout)
		defer cancel()
	}
	setTimeout(req, timeout)
	req = req.WithContext(c)
	setCaller(req)
	if color := metadata.String(c, metadata.Color); color != "" {
		setColor(req, color)
	}
	if resp, err = client.client.Do(req); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		code = "failed"
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= xhttp.StatusBadRequest {
		err = pkgerr.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, realURL(req))
		code = strconv.Itoa(resp.StatusCode)
		return
	}
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		return
	}
	return
}

// Do sends an HTTP request and returns an HTTP json response.
func (client *Client) Do(c context.Context, req *xhttp.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// JSON sends an HTTP request and returns an HTTP json response.
func (client *Client) JSON(c context.Context, req *xhttp.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// PB sends an HTTP request and returns an HTTP proto response.
func (client *Client) PB(c context.Context, req *xhttp.Request, res proto.Message, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		if err = proto.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

func (client *Client) onBreaker(breaker breaker.Breaker, err *error) {
	if err != nil && *err != nil {
		breaker.MarkFailed()
	} else {
		breaker.MarkSuccess()
	}
}

// sign calc appkey and appsecret sign.
func (client *Client) sign(params url.Values) (query string, err error) {
	client.mutex.RLock()
	key := client.conf.Key
	secret := client.conf.Secret
	client.mutex.RUnlock()
	if params == nil {
		params = url.Values{}
	}
	params.Set(_appKey, key)
	if params.Get(_appSecret) != "" {
		log.Warn("utils http get must not have parameter appSecret")
	}
	if params.Get(_ts) == "" {
		params.Set(_ts, strconv.FormatInt(time.Now().Unix(), 10))
	}
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(tmp)
	b.WriteString(secret)
	mh := md5.Sum(b.Bytes())
	// query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(hex.EncodeToString(mh[:]))
	query = qb.String()
	return
}

// realUrl return url with http://host/params.
func realURL(req *xhttp.Request) string {
	if req.Method == xhttp.MethodGet {
		return req.URL.String()
	} else if req.Method == xhttp.MethodPost {
		ru := req.URL.Path
		if req.Body != nil {
			rd, ok := req.Body.(io.Reader)
			if ok {
				buf := bytes.NewBuffer([]byte{})
				buf.ReadFrom(rd)
				ru = ru + "?" + buf.String()
			}
		}
		return ru
	}
	return req.URL.Path
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
