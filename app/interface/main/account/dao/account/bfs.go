package account

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

	"go-common/app/interface/main/account/conf"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// UploadImage upload bfs.
func (d *Dao) UploadImage(c context.Context, fileType string, bs []byte, bfs *conf.BFS) (location string, err error) {
	url := bfs.URL + bfs.Bucket + "/"
	req, err := http.NewRequest(bfs.Method, url, bytes.NewBuffer(bs))
	if err != nil {
		log.Error("account-interface: http.NewRequest error (%v) | fileType(%s)", err, fileType)
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(bfs.Key, bfs.Secret, bfs.Method, bfs.Bucket, expire)
	req.Header.Set("Host", url)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	// timeout
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.BFS.Timeout))
	req = req.WithContext(ctx)
	defer cancel()
	resp, err := d.bfsClient.Do(req)
	if err != nil {
		log.Error("account-interface: d.Client.Do error(%v) | url(%s)", err, url)
		err = ecode.BfsUploadServiceUnavailable
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("account-interface: Upload http.StatusCode nq http.StatusOK (%d) | url(%s)", resp.StatusCode, d.c.BFS.URL)
		err = errors.New("Upload failed")
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		log.Error("account-interface: strconv.Itoa err, code(%s) | url(%s)", code, url)
		err = errors.New("Upload failed")
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
	content = fmt.Sprintf("%s\n%s\n\n%d\n", method, bucket, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}
