package jpush

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"go-common/library/stat"
	"go-common/library/stat/prom"
)

const (
	_charset         = "UTF-8"
	_contentTypeJSON = "application/json"

	_pushURL = "https://api.jpush.cn/v3/push"
	// _scheduleURL = "https://api.jpush.cn/v3/schedules"
	// _reportURL   = "https://report.jpush.cn/v3/received"
)

// PushResponse .
type PushResponse struct {
	SendNo        interface{} `json:"sendno,omitempty"`
	MsgID         interface{} `json:"msg_id,omitempty"`
	IllegalTokens []string    `json:"illegal_rids,omitempty"`
	Retry         bool        // 是否需要重试请求
	Error         struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Client jpush http client.
type Client struct {
	Auth    string
	Stats   stat.Stat
	Timeout time.Duration
}

// NewClient new client.
func NewClient(appKey, secret string, timeout time.Duration) *Client {
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(appKey+":"+secret))
	return &Client{
		Auth:    auth,
		Stats:   prom.HTTPClient,
		Timeout: timeout,
	}
}

// Push push notification.
func (cli *Client) Push(payload *Payload) (res *PushResponse, err error) {
	res = new(PushResponse)
	if cli.Stats != nil {
		now := time.Now()
		defer func() {
			cli.Stats.Timing(_pushURL, int64(time.Since(now)/time.Millisecond))
			// log.Info("jpush stats timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				cli.Stats.Incr(_pushURL, "failed")
			}
		}()
	}
	bs, err := payload.ToBytes()
	if err != nil {
		return
	}
	req, _ := http.NewRequest("POST", _pushURL, bytes.NewBuffer(bs))
	req.Header.Add("Charset", _charset)
	req.Header.Add("Authorization", cli.Auth)
	req.Header.Add("Content-Type", _contentTypeJSON)
	client := &http.Client{Timeout: cli.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		res.Retry = true
		return
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &res); err != nil {
		return
	}
	if res.Error.Code == ErrRetry || res.Error.Code == ErrInternal {
		res.Retry = true
	}
	return
}
