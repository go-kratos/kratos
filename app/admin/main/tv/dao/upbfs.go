package dao

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"net/http"
	"strconv"

	"bytes"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_uploadURL = "/bfs/%s/%s"
	_template  = "%s\n%s\n%s\n%d\n"
	_method    = "PUT"
)

// Upload uploads cover to bfs
func (d *Dao) Upload(c context.Context, fileName string, fileType string, timing int64, data []byte) (location string, err error) {
	bfs := d.c.Bfs
	var (
		req  *http.Request
		resp *http.Response
		code int
		url  = fmt.Sprintf(bfs.Host+_uploadURL, bfs.Bucket, fileName)
	)
	// prepare the data of the file and init the request
	buf := new(bytes.Buffer)
	_, err = buf.Write(data)
	if err != nil {
		log.Error("Upload.buf.Write.error(%v)", err)
		err = ecode.RequestErr
		return
	}
	if req, err = http.NewRequest(_method, url, buf); err != nil {
		log.Error("http.NewRequest() Upload(%v) error(%v)", url, err)
		return
	}
	// request setting
	authorization := authorize(bfs.Key, bfs.Secret, _method, bfs.Bucket, fileName, timing)
	req.Header.Set("Date", fmt.Sprint(timing))
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", fileType)
	resp, err = d.bfsClient.Do(req)
	// response treatment
	if err != nil {
		log.Error("Bfs client.Do(%s) error(%v)", url, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Bfs status code error:%v", resp.StatusCode)
		return
	}
	code, err = strconv.Atoi(resp.Header.Get("code"))
	if err != nil || code != 200 {
		err = fmt.Errorf("Bfs response code error:%v", code)
		return
	}
	location = resp.Header.Get("Location")
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
