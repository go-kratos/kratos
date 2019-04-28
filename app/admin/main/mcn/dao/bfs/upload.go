package bfs

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

const (
	_uploadURL = "/bfs/%s/%s"
	_template  = "%s\n%s\n%s\n%d\n"
	_method    = "PUT"
)

// Upload upload picture or log file to bfs
func (d *Dao) Upload(c context.Context, fileName, fileType string, expire int64, body io.Reader) (location string, err error) {
	var (
		url  string
		req  *http.Request
		resp *http.Response
		code int
	)
	client := &http.Client{}
	url = fmt.Sprintf(d.bfs+_uploadURL, d.bucket, fileName)
	if req, err = http.NewRequest(_method, url, body); err != nil {
		err = errors.Errorf("http.NewRequest(url(%s), body(%+v)) error(%+v)", url, body, err)
		return
	}
	authorization := authorize(d.key, d.secret, _method, d.bucket, fileName, expire)
	req.Header.Set("Host", d.bfs)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	resp, err = client.Do(req)
	if err != nil {
		err = errors.Errorf("resp.Do(%s) error(%+v)", url, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("status code(%d) error(%+v)", resp.StatusCode, err)
		return
	}
	code, err = strconv.Atoi(resp.Header.Get("code"))
	if err != nil || code != http.StatusOK {
		err = errors.Errorf("response code(%d) error(%+v)", resp.StatusCode, err)
		return
	}
	location = resp.Header.Get("Location")
	return
}

// Authorize returns authorization for upload file to bfs
func authorize(key, secret, method, bucket, file string, expire int64) (authorization string) {
	var (
		content   string
		mac       hash.Hash
		signature string
	)
	content = fmt.Sprintf(_template, method, bucket, file, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}
