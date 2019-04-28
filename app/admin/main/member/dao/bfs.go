package dao

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go-common/app/admin/main/member/conf"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// UploadImage upload bfs.
func (d *Dao) UploadImage(c context.Context, fileType string, bs []byte, bfsConf *conf.BFS) (location string, err error) {
	url := bfsConf.URL + bfsConf.Bucket + "/"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bs))
	if err != nil {
		err = errors.Wrap(err, " UploadImage http.NewRequest() failed")
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(bfsConf.Key, bfsConf.Secret, "PUT", bfsConf.Bucket, "", expire)
	req.Header.Set("Host", url)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	// timeout
	ctx, cancel := context.WithTimeout(c, time.Duration(bfsConf.Timeout))
	req = req.WithContext(ctx)
	defer cancel()
	resp, err := d.bfsclient.Do(req)
	if err != nil {
		err = errors.Wrapf(ecode.BfsUploadServiceUnavailable, "d.bfsclient.Do error(%v) | url(%s)", err, url)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("UploadImage failed,http code:%d", resp.StatusCode)
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		err = errors.Wrap(ecode.String(code), "UploadImage failed")
		return
	}
	location = header.Get("Location")
	return
}

// DelImage del bfs image.
func (d *Dao) DelImage(c context.Context, fileName string, bfsConf *conf.BFS) (err error) {
	url := bfsConf.URL + bfsConf.Bucket + "/" + fileName
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		err = errors.Wrap(err, "DelImage http.NewRequest() failed")
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(bfsConf.Key, bfsConf.Secret, "DELETE", bfsConf.Bucket, fileName, expire)
	req.Header.Set("Host", url)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	// timeout
	ctx, cancel := context.WithTimeout(c, time.Duration(bfsConf.Timeout))
	req = req.WithContext(ctx)
	defer cancel()
	resp, err := d.bfsclient.Do(req)
	if err != nil {
		err = errors.Wrapf(ecode.BfsUploadServiceUnavailable, "d.bfsclient.Do error(%v) | url(%s)", err, url)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("DelImage failed,http code:%d", resp.StatusCode)
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		err = errors.Errorf("DelImage failed,res code:%s", code)
		return
	}
	return
}

// authorize returns authorization for upload file to bfs
func authorize(key, secret, method, bucket, fileName string, expire int64) (authorization string) {
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

//Image image.
func (d *Dao) Image(url string) ([]byte, error) {
	resp, err := d.bfsclient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Image url(%s) resp.StatusCode code(%v)", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
