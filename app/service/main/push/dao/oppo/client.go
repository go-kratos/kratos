package oppo

import (
	"encoding/json"
	"io"
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

// Client huawei push http client.
type Client struct {
	Auth       *Auth
	HTTPClient *http.Client
	Stats      stat.Stat
	Activity   string
}

// NewClient new huawei push HTTP client.
func NewClient(a *Auth, activity string, timeout time.Duration) *Client {
	return &Client{
		Auth:       a,
		HTTPClient: &http.Client{Timeout: timeout},
		Stats:      prom.HTTPClient,
		Activity:   activity,
	}
}

// Message saves push message content.
func (c *Client) Message(m *Message) (response *Response, err error) {
	now := time.Now()
	if c.Stats != nil {
		defer func() {
			c.Stats.Timing(_apiMessage, int64(time.Since(now)/time.Millisecond))
			log.Info("oppo message stats timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(_apiMessage, "failed")
			}
		}()
	}
	params := url.Values{}
	params.Add("auth_token", c.Auth.Token)
	params.Add("title", m.Title)
	params.Add("content", m.Content)
	params.Add("click_action_type", strconv.Itoa(m.ActionType))
	params.Add("click_action_activity", c.Activity)
	params.Add("action_parameters", m.ActionParams)
	params.Add("off_line_ttl", strconv.Itoa(m.OfflineTTL))
	params.Add("call_back_url", m.CallbackURL)
	req, err := http.NewRequest(http.MethodPost, _apiMessage, strings.NewReader(params.Encode()))
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

// Push push notification.
func (c *Client) Push(msgID string, tokens []string) (response *Response, err error) {
	now := time.Now()
	if c.Stats != nil {
		defer func() {
			c.Stats.Timing(_apiPushBroadcast, int64(time.Since(now)/time.Millisecond))
			log.Info("oppo push stats timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(_apiPushBroadcast, "failed")
			}
		}()
	}
	params := url.Values{}
	params.Add("auth_token", c.Auth.Token)
	params.Add("message_id", msgID)
	params.Add("target_type", _pushTypeToken)
	params.Add("registration_ids", strings.Join(tokens, ";"))
	req, err := http.NewRequest(http.MethodPost, _apiPushBroadcast, strings.NewReader(params.Encode()))
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

// PushOne push single notification.
func (c *Client) PushOne(m *Message, token string) (response *Response, err error) {
	now := time.Now()
	if c.Stats != nil {
		defer func() {
			c.Stats.Timing(_apiPushUnicast, int64(time.Since(now)/time.Millisecond))
			log.Info("oppo pushOne stats timing: %v", int64(time.Since(now)/time.Millisecond))
			if err != nil {
				c.Stats.Incr(_apiPushUnicast, "failed")
			}
		}()
	}
	m.ActionActivity = c.Activity
	params := url.Values{}
	msg, _ := json.Marshal(&struct {
		TargetType   string   `json:"target_type"`
		Token        string   `json:"registration_id"`
		Notification *Message `json:"notification"`
	}{
		TargetType:   _pushTypeToken,
		Token:        token,
		Notification: m,
	})
	params.Add("auth_token", c.Auth.Token)
	params.Add("message", string(msg))
	req, err := http.NewRequest(http.MethodPost, _apiPushUnicast, strings.NewReader(params.Encode()))
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
