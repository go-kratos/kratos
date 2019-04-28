package datacenter

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"go-common/library/log"

	pkgerr "github.com/pkg/errors"
	"go-common/library/stat"
	"strconv"
)

/*
  访问数据平台的http client，处理了签名、接口监控等
*/

//ClientConfig client config
type ClientConfig struct {
	Key       string
	Secret    string
	Dial      time.Duration
	Timeout   time.Duration
	KeepAlive time.Duration
}

//New new client
func New(c *ClientConfig) *HttpClient {
	return &HttpClient{
		client: &http.Client{},
		conf:   c,
	}
}

//HttpClient http client
type HttpClient struct {
	client *http.Client
	conf   *ClientConfig
	Debug  bool
}

//Response response
type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

const (
	keyAppKey     = "appKey"
	keyAppID      = "apiId"
	keyTimeStamp  = "timestamp"
	keySign       = "sign"
	keySignMethod = "signMethod"
	keyVersion    = "version"
	//TimeStampFormat time format in second
	TimeStampFormat = "2006-01-02 15:04:05"
)

var (
	clientStats = stat.HTTPClient
)

// Get issues a GET to the specified URL.
func (client *HttpClient) Get(c context.Context, uri string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(http.MethodGet, uri, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// NewRequest new http request with method, uri, ip, values and headers.
// TODO(zhoujiahui): param realIP should be removed later.
func (client *HttpClient) NewRequest(method, uri string, params url.Values) (req *http.Request, err error) {
	signStr, err := client.sign(params)
	if err != nil {
		err = pkgerr.Wrapf(err, "uri:%s,params:%v", uri, params)
		return
	}
	params.Add(keySign, signStr)
	enc := params.Encode()
	ru := uri
	if enc != "" {
		ru = uri + "?" + enc
	}
	if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, ru, nil)
	} else {
		req, err = http.NewRequest(http.MethodPost, uri, strings.NewReader(enc))
	}
	if err != nil {
		err = pkgerr.Wrapf(err, "method:%s,uri:%s", method, ru)
		return
	}
	const (
		_contentType = "Content-Type"
		_urlencoded  = "application/x-www-form-urlencoded"
	)
	if method == http.MethodPost {
		req.Header.Set(_contentType, _urlencoded)
	}

	return
}

// Do sends an HTTP request and returns an HTTP json response.
func (client *HttpClient) Do(c context.Context, req *http.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s, response:%s", req.URL.Host, realURL(req), string(bs))
		}
	}
	return
}

//Raw get from url
func (client *HttpClient) Raw(c context.Context, req *http.Request, v ...string) (bs []byte, err error) {
	var resp *http.Response
	var uri = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.Path)

	var now = time.Now()
	var code string
	defer func() {
		clientStats.Timing(uri, int64(time.Since(now)/time.Millisecond))
		if code != "" {
			clientStats.Incr(uri, code)
		}
	}()
	req = req.WithContext(c)
	if resp, err = client.client.Do(req); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		code = "failed"
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		err = pkgerr.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, realURL(req))
		code = strconv.Itoa(resp.StatusCode)
		return
	}
	if bs, err = readAll(resp.Body, 16*1024); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		return
	}
	if client.Debug {
		log.Info("reqeust: host:%s, url:%s, response body:%s", req.URL.Host, realURL(req), string(bs))
	}
	return
}

// sign calc appkey and appsecret sign.
// see http://info.bilibili.co/pages/viewpage.action?pageId=5410881#id-%E6%95%B0%E6%8D%AE%E7%9B%98%EF%BC%8D%E5%AE%89%E5%85%A8%E8%AE%A4%E8%AF%81-%E4%BA%8C%E7%AD%BE%E5%90%8D%E7%AE%97%E6%B3%95
func (client *HttpClient) sign(params url.Values) (sign string, err error) {
	key := client.conf.Key
	secret := client.conf.Secret
	if params == nil {
		params = url.Values{}
	}

	params.Set(keyAppKey, key)
	params.Set(keyVersion, "1.0")
	if params.Get(keyTimeStamp) == "" {
		params.Set(keyTimeStamp, time.Now().Format(TimeStampFormat))
	}
	params.Set(keySignMethod, "md5")

	var needSignParams = url.Values{}
	needSignParams.Add(keyAppKey, key)
	needSignParams.Add(keyTimeStamp, params.Get(keyTimeStamp))
	needSignParams.Add(keyVersion, params.Get(keyVersion))

	//tmp := params.Encode()
	var valueMap = map[string][]string(needSignParams)
	var buf bytes.Buffer
	// 开头与结尾加secret
	buf.Write([]byte(secret))
	keys := make([]string, 0, len(valueMap))
	for k := range valueMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := valueMap[k]
		prefix := k
		buf.WriteString(prefix)
		for _, v := range vs {
			buf.WriteString(v)
			break
		}
	}

	buf.Write([]byte(secret))
	var md5 = md5.New()
	md5.Write(buf.Bytes())
	sign = fmt.Sprintf("%X", md5.Sum(nil))
	return
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

// realUrl return url with http://host/params.
func realURL(req *http.Request) string {
	if req.Method == http.MethodGet {
		return req.URL.String()
	} else if req.Method == http.MethodPost {
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

// SetTransport set client transport
func (client *HttpClient) SetTransport(t http.RoundTripper) {
	client.client.Transport = t
}
