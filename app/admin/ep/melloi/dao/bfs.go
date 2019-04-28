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
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_bucketName = "melloi"
	_melloiURI  = "/bfs/" + _bucketName
)

// UploadImg upload img
func (d *Dao) UploadImg(c context.Context, imgContent []byte, imgName string) (location string, err error) {

	var (
		req    *http.Request
		reqURL = d.c.BfsConf.Host + _melloiURI + "/" + imgName
		resp   *http.Response
		code   int
	)
	buf := new(bytes.Buffer)
	if _, err = buf.Write(imgContent); err != nil {
		log.Error("Upload.buf.Write.error(%v)", err)
		return
	}

	if req, err = http.NewRequest("PUT", reqURL, buf); err != nil {
		err = ecode.RequestErr
		return
	}

	authorization := d.Authorize(d.c.BfsConf.AccessKey, d.c.BfsConf.AccessSecret, http.MethodPut, _bucketName, imgName, time.Now().Unix())
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "image/png")

	if resp, err = d.client.Do(req); err != nil {
		log.Error("Upload.client.Do.error(%v)", err)
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		err = ecode.RequestErr
		return
	case http.StatusUnauthorized:
		// 验证不通过
		err = ecode.BfsUploadAuthErr
		return
	case http.StatusRequestEntityTooLarge:
		err = ecode.FileTooLarge
		return
	case http.StatusNotFound:
		err = ecode.NothingFound
		return
	case http.StatusMethodNotAllowed:
		err = ecode.MethodNotAllowed
		return
	case http.StatusServiceUnavailable:
		err = ecode.BfsUploadServiceUnavailable
		return
	case http.StatusInternalServerError:
		err = ecode.ServerErr
		return
	default:
		err = ecode.BfsUploadStatusErr
		return
	}
	if code, err = strconv.Atoi(resp.Header.Get("code")); err != nil || code != 200 {
		err = ecode.BfsUploadCodeErr
		return
	}
	location = resp.Header.Get("location")
	return
}

// Authorize .
func (d *Dao) Authorize(key, secret, method, bucket, fileName string, expire int64) (authorization string) {
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
