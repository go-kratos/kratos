package jpush

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"time"

	"go-common/library/stat"
	"go-common/library/stat/prom"
)

const (
	pushURL = "https://api.jpush.cn/v3/push"
	// 如果创建的极光应用分配的北京机房，并且 API 调用方的服务器也位于北京，则比较适合调用极光北京机房的 API，可以提升一定的响应速度。
	// PUSH_URL       = "https://bjapi.push.jiguang.cn/v3/push"
	// VALIDATE_URL   = "https://api.jpush.cn/v3/push/validate"
	// GROUP_PUSH_URL = "https://api.jpush.cn/v3/grouppush"
)

// Client for JPush
type Client struct {
	Auth    string
	Stats   stat.Stat
	Timeout time.Duration
}

// New .
func New(appKey string, secretKey string, timeout time.Duration) *Client {
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(appKey+":"+secretKey))
	return &Client{
		Auth:    auth,
		Stats:   prom.HTTPClient,
		Timeout: timeout,
	}
}

// Push .
func (clt *Client) Push(b []byte) (resp []byte, err error) {
	if clt.Stats != nil {
		now := time.Now()
		defer func() {
			clt.Stats.Timing(pushURL, int64(time.Since(now)/time.Millisecond))
			// log.Info("jpush stats timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				clt.Stats.Incr(pushURL, "failed")
			}
		}()
	}
	req, err := http.NewRequest("POST", pushURL, bytes.NewBuffer(b))
	req.Header.Add("Charset", "UTF-8")
	req.Header.Add("Authorization", clt.Auth)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: clt.Timeout}
	httpResp, err := client.Do(req)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	return ioutil.ReadAll(httpResp.Body)
}

// GetTimeout .
func (clt *Client) GetTimeout() time.Duration {
	return clt.Timeout
}
