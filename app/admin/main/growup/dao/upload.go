package dao

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strconv"

	"go-common/library/log"
	"go-common/library/net/trace"
)

const (
	_uploadURL  = "/%s/%s"
	_moduleName = "http_client"
	_template   = "%s\n%s\n%s\n%d\n"
	_method     = "PUT"
)

var (
	errUpload = errors.New("Upload failed")
)

// Upload upload picture or log file to bfs
func (d *Dao) Upload(c context.Context, fileName, fileType string, expire int64, body io.Reader) (location string, err error) {
	var (
		url    string
		req    *http.Request
		resp   *http.Response
		header http.Header
		code   string
	)
	bfsConf := d.c.Bfs
	url = fmt.Sprintf(bfsConf.Addr+_uploadURL, bfsConf.Bucket, fileName)
	if req, err = http.NewRequest(_method, url, body); err != nil {
		log.Error("http.NewRequest() Upload(%v) error(%v)", url, err)
		return
	}
	authorization := authorize(bfsConf.Key, bfsConf.Secret, _method, bfsConf.Bucket, fileName, expire)
	req.Header.Set("Host", bfsConf.Addr)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	t, ok := trace.FromContext(c)
	if ok {
		t = t.Fork("http_client", "bfs_upload")
		t.SetTag(trace.String(trace.TagAddress, _moduleName), trace.String(trace.TagComment, req.URL.Path))
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Finish(&err)
		log.Error("httpClient.Do(%s) error(%v)", url, err)
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
