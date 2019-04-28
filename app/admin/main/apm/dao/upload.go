package dao

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"net/http"
	"strconv"
)

func authorize(key, secret, method, bucket string, expire int64) (authorization string) {
	var (
		content   string
		hash2     hash.Hash
		signature string
		template  = "%s\n%s\n\n%d\n"
	)
	content = fmt.Sprintf(template, method, bucket, expire)
	hash2 = hmac.New(sha1.New, []byte(secret))
	hash2.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(hash2.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}

// UploadProxy upload file to bfs with no filename.
func (d *Dao) UploadProxy(c context.Context, fileType string, expire int64, body []byte) (url string, err error) {
	var (
		req       *http.Request
		resp      *http.Response
		header    http.Header
		code      string
		bfs       = d.c.Bfs
		uploadURL = "/%s"
	)
	url = fmt.Sprintf(bfs.Addr+uploadURL, bfs.Bucket)
	if req, err = http.NewRequest(http.MethodPut, url, bytes.NewReader(body)); err != nil {
		err = fmt.Errorf("dao.UploadProxy NewRequest error(%v)", err)
		return
	}
	authorization := authorize(bfs.Key, bfs.Secret, http.MethodPut, bfs.Bucket, expire)
	req.Header.Set("Host", bfs.Addr)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("dao.UploadProxy error return_status(%d)", resp.StatusCode)
		return
	}
	header = resp.Header
	code = header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		err = fmt.Errorf("dao.UploadProxy Upload url(%s) return_code(%s)", url, code)
		return
	}
	url = header.Get("Location")
	return
}
