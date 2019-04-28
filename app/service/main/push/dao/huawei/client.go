package huawei

// http://developer.huawei.com/consumer/cn/service/hms/catalog/huaweipush.html?page=hmssdk_huaweipush_api_reference_s2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
	"go-common/library/stat"
	"go-common/library/stat/prom"
)

const (
	_pushURL = "https://api.push.hicloud.com/pushsend.do"
	_nspSvc  = "openpush.message.api.send"
	_ver     = "1" // current SDK version
	// ResponseCodeSuccess success code
	ResponseCodeSuccess = "80000000"
	// ResponseCodeSomeTokenInvalid some tokens failed
	ResponseCodeSomeTokenInvalid = "80100000"
	// ResponseCodeAllTokenInvalid all tokens failed
	ResponseCodeAllTokenInvalid = "80100002"
	// ResponseCodeAllTokenInvalidNew .
	ResponseCodeAllTokenInvalidNew = "80300007"
)

var (
	// ErrLimit .
	ErrLimit = errors.New("触发华为系统级流控")
)

// Client huawei push http client.
type Client struct {
	Access     *Access
	HTTPClient *http.Client
	Stats      stat.Stat
	SDKCtx     string
	Package    string
}

// ver huawei push service version.
type ver struct {
	Ver   string `json:"ver"`
	AppID string `json:"appId"`
}

// NewClient new huawei push HTTP client.
func NewClient(pkg string, a *Access, timeout time.Duration) *Client {
	ctx, _ := json.Marshal(ver{Ver: _ver, AppID: a.AppID})
	return &Client{
		Access:     a,
		HTTPClient: &http.Client{Timeout: timeout},
		Stats:      prom.HTTPClient,
		SDKCtx:     string(ctx),
		Package:    pkg,
	}
}

/*
Push push notifications.
access_token: 必选，使用OAuth2进行鉴权时的ACCESSTOKEN
nsp_ts: 必选，服务请求时间戳，自GMT 时间 1970-1-1 0:0:0至今的秒数。如果传入的时间与服务器时间相 差5分钟以上，服务器可能会拒绝请求。
nsp_svc: 必选， 本接口固定为openpush.message.api.send
device_token_list: 以半角逗号分隔的华为PUSHTOKEN列表，单次最多只是1000个
expire_time: 格式ISO 8601[6]:2013-06-03T17:30，采用本地时间精确到分钟
payload: 描述投递消息的JSON结构体，描述PUSH消息的:类型、内容、显示、点击动作、报表统计和扩展信 息。具体参考下面的详细说明。
*/
func (c *Client) Push(payload *Message, tokens []string, expire time.Time) (response *Response, err error) {
	now := time.Now()
	if c.Stats != nil {
		defer func() {
			c.Stats.Timing(_pushURL, int64(time.Since(now)/time.Millisecond))
			log.Info("huawei stats timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(_pushURL, "failed")
			}
		}()
	}
	pl, _ := payload.SetPkg(c.Package).JSON()
	reqURL := _pushURL + "?nsp_ctx=" + url.QueryEscape(c.SDKCtx)
	tokenStr, _ := json.Marshal(tokens)
	params := url.Values{}
	params.Add("access_token", c.Access.Token)
	params.Add("nsp_ts", strconv.FormatInt(now.Unix(), 10))
	params.Add("nsp_svc", _nspSvc)
	params.Add("device_token_list", string(tokenStr))
	params.Add("expire_time", expire.Format("2006-01-02T15:04"))
	params.Add("payload", pl)
	req, err := http.NewRequest(http.MethodPost, reqURL, strings.NewReader(params.Encode()))
	if err != nil {
		log.Error("http.NewRequest() error(%v)", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "Keep-Alive")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Error("HTTPClient.Do() error(%v)", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusServiceUnavailable {
		return nil, ErrLimit
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("huawei Push http code(%d)", res.StatusCode)
		return
	}
	response = &Response{}
	nspStatus := res.Header.Get("NSP_STATUS")
	if nspStatus != "" {
		log.Error("push huawei system error, NSP_STATUS(%s)", nspStatus)
		response.Code = nspStatus
		response.Msg = "NSP_STATUS error"
		return
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		return
	}
	if err = json.Unmarshal(bs, &response); err != nil {
		log.Error("json decode body(%s) error(%v)", string(bs), err)
	}
	return
}
