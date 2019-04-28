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
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_template = "%s\n%s\n\n%d\n"
	_method   = "PUT"
)

// Dao is bfs dao.
type Dao struct {
	c      *conf.Config
	client *http.Client
	bucket string
	url    string
	key    string
	secret string
}

//New bfs dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// http client
		client: &http.Client{
			Timeout: time.Duration(c.Bfs.Timeout),
		},
		bucket: c.Bfs.Bucket,
		url:    c.Bfs.Addr,
		key:    c.Bfs.Key,
		secret: c.Bfs.Secret,
	}
	return
}

// Upload upload bfs.
func (d *Dao) Upload(c context.Context, fileType string, body io.Reader) (location string, err error) {
	req, err := http.NewRequest(_method, d.url, body)
	if err != nil {
		log.Error("http.NewRequest error (%v) | fileType(%s) body(%v)", err, fileType, body)
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(d.key, d.secret, _method, d.bucket, expire)
	req.Header.Set("Host", d.url)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	log.Error("Authorization_:%v", authorization)
	// timeout
	c, cancel := context.WithTimeout(c, time.Duration(d.c.Bfs.Timeout))
	req = req.WithContext(c)
	defer cancel()
	resp, err := d.client.Do(req)
	if err != nil {
		log.Error("d.Client.Do error(%v) | _url(%s) req(%v)", err, d.url, req)
		err = ecode.BfsUploadServiceUnavailable
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Upload http.StatusCode nq http.StatusOK (%d) | url(%s)", resp.StatusCode, d.url)
		err = ecode.BfsUploadStatusErr
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		log.Error("strconv.Itoa err, code(%s) | url(%s)", code, d.url)
		err = ecode.BfsUploadCodeErr
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
