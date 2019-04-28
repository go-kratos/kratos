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

	"go-common/app/admin/main/up/conf"
	"go-common/library/log"
)

const (
	uploadurl = "/%s/%s"
	template  = "%s\n%s\n%s\n%d\n"
	method    = "PUT"
)

var (
	errUpload = errors.New("Upload failed")
)

// Upload upload picture or log file to bfs
func Upload(c context.Context, fileName, fileType string, expire int64, body io.Reader, bfsConf *conf.Bfs) (location string, err error) {
	var (
		url    string
		req    *http.Request
		resp   *http.Response
		header http.Header
		code   string
	)

	url = fmt.Sprintf(bfsConf.Addr+uploadurl, bfsConf.Bucket, fileName)
	if req, err = http.NewRequest(method, url, body); err != nil {
		log.Error("http.NewRequest() Upload(%v) error(%v)", url, err)
		return
	}
	log.Info("upload url={%s}", url)
	authorization := Authorize(bfsConf.Key, bfsConf.Secret, method, bfsConf.Bucket, fileName, expire)
	req.Header.Set("Host", bfsConf.Addr)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
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

// Authorize authorize returns authorization for upload file to bfs
func Authorize(key, secret, method, bucket, file string, expire int64) (authorization string) {
	var (
		content   string
		mac       hash.Hash
		signature string
	)
	content = fmt.Sprintf(template, method, bucket, file, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}
