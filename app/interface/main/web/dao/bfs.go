package dao

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
	"go-common/library/net/trace"
)

const (
	_uploadURL = "/%s"
	_family    = "http_client"
	_template  = "%s\n%s\n\n%d\n"
	_method    = "PUT"
	_fileType  = "txt"
)

var (
	errUpload = errors.New("Upload failed")
)

// Upload upload picture or log file to bfs
func (d *Dao) Upload(c context.Context, content string, expire int64) (location string, err error) {
	var (
		url    string
		req    *http.Request
		resp   *http.Response
		header http.Header
		code   string
	)
	bfsConf := d.c.Bfs
	url = fmt.Sprintf(bfsConf.Addr+_uploadURL, bfsConf.Bucket)
	if req, err = http.NewRequest(_method, url, strings.NewReader(content)); err != nil {
		log.Error("http.NewRequest() Upload(%v) error(%v)", url, err)
		return
	}
	authorization := authorize(bfsConf.Key, bfsConf.Secret, _method, bfsConf.Bucket, expire)
	req.Header.Set("Host", bfsConf.Addr)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", _fileType)
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, req.URL.Path)
		defer t.Finish(&err)
	}
	c, cancel := context.WithTimeout(c, time.Duration(d.c.Bfs.Timeout))
	req = req.WithContext(c)
	defer cancel()
	resp, err = d.bfsClient.Do(req)
	if err != nil {
		log.Error("bfsClient.Do(%s) error(%v)", url, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Upload url(%s) http.statuscode:%d", url, resp.StatusCode)
		err = errUpload
		return
	}
	header = resp.Header
	code = header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		log.Error("Upload url(%s) code:%s", url, code)
		err = errUpload
		return
	}
	location = header.Get("Location")
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
