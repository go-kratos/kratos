package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"math/rand"
	"net/url"
	"path/filepath"

	"strings"
	"time"
)

// const .
const (
	_template     = "%s\n%s\n%s\n%d\n"
	BfsEasyPath   = "/bfs/mcn"
	TimeFormatSec = "2006-01-02 15:04:05"
	TimeFormatDay = "2006-01-02"
	AllActiveTid  = 65535 //mcn_data_summary表active_tid所有分区
	DefaultTyName = "默认"
)

// BuildBfsURL is.
func BuildBfsURL(raw string, appkey, secret, bucket, ePath string) string {
	if raw == "" {
		return ""
	}
	ori, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if strings.HasPrefix(ori.Path, ePath) {
		token := authorize(appkey, secret, "GET", bucket, filepath.Base(ori.Path), time.Now().Unix())
		p := url.Values{}
		p.Set("token", token)
		ori.RawQuery = p.Encode()
	}
	if ori.Hostname() == "" {
		ori.Host = fmt.Sprintf("i%d.hdslb.com", rand.Int63n(3))
		ori.Scheme = "http"
	}
	return ori.String()
}

// Authorize returns authorization for upload file to bfs
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
