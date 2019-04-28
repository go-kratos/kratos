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
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	errUpload = errors.New("Upload failed")
)

// Upload upload bfs.
func (d *Dao) Upload(c context.Context, fileType string, filename string, body io.Reader) (location string, err error) {
	req, err := http.NewRequest(d.c.BFS.Method, d.c.BFS.URL+filename, body)
	if err != nil {
		log.Error("http.NewRequest error (%v) | fileType(%s)", err, fileType)
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(d.c.BFS.Key, d.c.BFS.Secret, d.c.BFS.Method, d.c.BFS.Bucket, filename, expire)
	req.Header.Set("Host", d.c.BFS.Host)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	// timeout
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.BFS.Timeout))
	req = req.WithContext(ctx)
	defer cancel()
	resp, err := d.client.Do(req)
	if err != nil {
		log.Error("d.Client.Do error(%v) | url(%s)", err, d.c.BFS.URL+filename)
		err = ecode.BfsUploadServiceUnavailable
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Upload http.StatusCode nq http.StatusOK (%d) | url(%s)", resp.StatusCode, d.c.BFS.URL+filename)
		err = errUpload
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		log.Error("strconv.Itoa err, code(%s) | url(%s)", code, d.c.BFS.URL+filename)
		err = errUpload
		return
	}
	location = header.Get("Location")
	return
}

// authorize returns authorization for upload file to bfs
func authorize(key, secret, method, bucket string, fname string, expire int64) (authorization string) {
	var (
		content   string
		mac       hash.Hash
		signature string
	)
	content = fmt.Sprintf("%s\n%s\n%s\n%d\n", method, bucket, fname, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}
