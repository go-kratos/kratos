package mi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go-common/library/log"
	"go-common/library/stat"
	"go-common/library/stat/prom"
)

// Client Xiaomi http client.
type Client struct {
	Header     http.Header
	HTTPClient *http.Client
	Package    string
	URL        string
	Stats      stat.Stat
}

// NewClient returns a Client with request header and auth.
func NewClient(pkg, auth string, timeout time.Duration) *Client {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("Authorization", AuthPrefix+auth)
	// transport := &http.Transport{
	// 	Proxy: func(_ *http.Request) (*url.URL, error) {
	// 		return url.Parse("http://10.28.10.11:80")
	// 	},
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 		DualStack: true,
	// 	}).DialContext,
	// 	MaxIdleConns:          100,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// }
	return &Client{
		Header:     header,
		HTTPClient: &http.Client{Timeout: timeout},
		// HTTPClient: &http.Client{Timeout: timeout, Transport: transport},
		Package: pkg,
		Stats:   prom.HTTPClient,
	}
}

// SetProductionURL sets Production URL.
func (c *Client) SetProductionURL(url string) {
	c.URL = ProductionHost + url
}

// SetDevelopmentURL sets Production URL.
func (c *Client) SetDevelopmentURL(url string) {
	c.URL = DevHost + url
}

// SetVipURL sets VIP URL.
func (c *Client) SetVipURL(url string) {
	c.URL = VipHost + url
}

// SetStatusURL sets feedback URL.
func (c *Client) SetStatusURL() {
	c.URL = ProductionHost + StatusURL
}

// Push sends a Notification to Xiaomi push service.
func (c *Client) Push(xm *XMMessage) (response *Response, err error) {
	if c.Stats != nil {
		now := time.Now()
		defer func() {
			c.Stats.Timing(c.URL, int64(time.Since(now)/time.Millisecond))
			log.Info("mi stats timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(c.URL, "failed")
			}
		}()
	}
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, c.URL, bytes.NewBuffer([]byte(xm.xmuv.Encode()))); err != nil {
		log.Error("http.NewRequest() error(%v)", err)
		return
	}
	req.Header = c.Header
	var res *http.Response
	if res, err = c.HTTPClient.Do(req); err != nil {
		log.Error("HTTPClient.Do() error(%v)", err)
		return
	}
	defer res.Body.Close()
	response = &Response{}
	var bs []byte
	bs, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		return
	} else if len(bs) == 0 {
		return
	}
	if e := json.Unmarshal(bs, &response); e != nil {
		if e != io.EOF {
			log.Error("json decode body(%s) error(%v)", string(bs), e)
		}
	}
	return
}

// MockPush mock push.
func (c *Client) MockPush(xm *XMMessage) (response *Response, err error) {
	if c.Stats != nil {
		now := time.Now()
		defer func() {
			c.Stats.Timing(c.URL, int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(c.URL, "mi push mock")
			}
		}()
	}
	time.Sleep(200 * time.Millisecond)
	response = &Response{Code: ResultCodeOk, Result: ResultOk}
	return
}

// InvalidTokens get invalid tokens.
func (c *Client) InvalidTokens() (response *Response, err error) {
	if c.Stats != nil {
		now := time.Now()
		defer func() {
			c.Stats.Timing(c.URL, int64(time.Since(now)/time.Millisecond))
			log.Info("mi invalidTokens timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(c.URL, "failed")
			}
		}()
	}
	req, err := http.NewRequest(http.MethodGet, feedbackHost+feedbackURI, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v)", c.URL, err)
		return
	}
	req.Header = c.Header
	c.HTTPClient.Timeout = time.Minute
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Error("HTTPClient.Do() error(%v)", err)
		return
	}
	defer res.Body.Close()
	response = &Response{}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		return
	} else if len(bs) == 0 {
		return
	}
	if e := json.Unmarshal(bs, &response); e != nil {
		if e != io.EOF {
			log.Error("json decode body(%s) error(%v)", string(bs), e)
		}
	}
	return
}

// UninstalledTokens get uninstalled tokens.
func (c *Client) UninstalledTokens() (response *UninstalledResponse, err error) {
	if c.Stats != nil {
		now := time.Now()
		defer func() {
			c.Stats.Timing(c.URL, int64(time.Since(now)/time.Millisecond))
			log.Info("mi UninstalledTokens timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(c.URL, "mi uninstalled tokens")
			}
		}()
	}
	req, err := http.NewRequest(http.MethodGet, emqHost+uninstalledURI+"?package_name="+c.Package, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v)", c.URL, err)
		return
	}
	req.Header = c.Header
	c.HTTPClient.Timeout = time.Minute
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Error("HTTPClient.Do() error(%v)", err)
		return
	}
	defer res.Body.Close()
	response = &UninstalledResponse{}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		return
	} else if len(bs) == 0 {
		return
	}
	if e := json.Unmarshal(bs, &response); e != nil {
		if e != io.EOF {
			log.Error("json decode body(%s) error(%v)", string(bs), e)
		}
		return
	}
	for _, s := range response.Result {
		s = strings.Replace(s, `\`, "", -1)
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		ud := UninstalledData{}
		if e := json.Unmarshal([]byte(s), &ud); e != nil {
			log.Error("json unmarshal(%s) error(%v)", s, e)
			continue
		}
		if ud.Token == "" {
			continue
		}
		response.Data = append(response.Data, ud.Token)
	}
	return
}
