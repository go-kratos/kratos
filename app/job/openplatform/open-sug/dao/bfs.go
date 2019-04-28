package dao

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

	"go-common/library/ecode"
	"go-common/library/log"
)

const _method = "PUT"

// Upload upload bfs.
func (d *Dao) Upload(c context.Context, fileType string, fileName string, body io.Reader) (location string, err error) {
	req, err := http.NewRequest(_method, d.c.Bfs.Addr+fileName, body)
	if err != nil {
		log.Error("http.NewRequest error (%v) | fileType(%s) body(%v)", err, fileType, body)
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(d.c.Bfs.Key, d.c.Bfs.Secret, _method, d.c.Bfs.Bucket, fileName, expire)
	req.Header.Set("Host", d.c.Bfs.Addr)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	// timeout
	c, cancel := context.WithTimeout(c, time.Duration(d.c.Bfs.Timeout))
	req = req.WithContext(c)
	defer cancel()
	resp, err := d.client.Do(req)
	if err != nil {
		log.Error("d.Client.Do error(%v) | _url(%s) req(%v)", err, d.c.Bfs.Addr, req)
		err = ecode.BfsUploadServiceUnavailable
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Upload http.StatusCode nq http.StatusOK (%d) | url(%s)", resp.StatusCode, d.c.Bfs.Addr)
		err = ecode.BfsUploadStatusErr
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		log.Error("strconv.Itoa err, code(%s) | url(%s)", code, d.c.Bfs.Addr)
		err = ecode.BfsUploadCodeErr
		return
	}
	location = header.Get("Location")
	return
}

func authorize(key, secret, method, bucket string, fileName string, expire int64) (authorization string) {
	var (
		content   string
		mac       hash.Hash
		signature string
	)
	content = fmt.Sprintf("%s\n%s\n%s\n%d\n", method, bucket, fileName, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}
