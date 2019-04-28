package report

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-view/conf"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_add      = "/videoup/archive/report"
	_timeout  = 800 * time.Millisecond
	_bucket   = "archive"
	_url      = "http://bfs.bilibili.co/bfs/archive/"
	_template = "%s\n%s\n\n%d\n"
	_method   = "PUT"
	_key      = "8d4e593ba7555502"
	_secret   = "0bdbd4c7caeeddf587c3c4daec0475"
)

// Dao is report dao
type Dao struct {
	client    *httpx.Client
	bfsClient *http.Client
	add       string
}

// New is appeal inital func .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:    httpx.NewClient(c.HTTPWrite),
		bfsClient: http.DefaultClient,
		add:       c.Host.Archive + _add,
	}
	return
}

// AddAppeal add appeal .
func (d *Dao) AddReport(c context.Context, mid, aid int64, mold int, ak, reason, pics string) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("access_key", ak)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("type", strconv.Itoa(mold))
	params.Set("reason", reason)
	params.Set("pics", pics)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.add, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.add+"?"+params.Encode())
	}
	return
}

// Upload imgage upload .
func (d *Dao) Upload(c context.Context, fileType string, body io.Reader) (location string, err error) {
	req, err := http.NewRequest(_method, _url, body)
	if err != nil {
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(_key, _secret, _method, _bucket, expire)
	req.Header.Set("Host", _url)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	// timeout
	c, cancel := context.WithTimeout(c, _timeout)
	req = req.WithContext(c)
	defer cancel()
	resp, err := d.bfsClient.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.Wrap(ecode.Int(resp.StatusCode), req.URL.String())
		return
	}
	code := resp.Header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		err = errors.Wrap(ecode.String(code), req.URL.String())
		return
	}
	location = resp.Header.Get("Location")
	return
}

// authorize returns authorization for upload file to bfs
func authorize(key, secret, method, bucket string, expire int64) (authorization string) {
	var (
		content   string
		mac       hash.Hash
		signature string
	)
	content = fmt.Sprintf(_template, method, bucket, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}
