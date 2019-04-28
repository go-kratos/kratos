package reply

import (
	"context"
	"go-common/app/interface/main/reply/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"net/url"
)

// NewCaptchaDao NewCaptchaDao
func NewCaptchaDao(c *httpx.ClientConfig) *CaptchaDao {
	return &CaptchaDao{httpClient: httpx.NewClient(c)}
}

// CaptchaDao CaptchaDao
type CaptchaDao struct {
	httpClient *httpx.Client
}

// Captcha Captcha.
func (s *CaptchaDao) Captcha(c context.Context) (string, string, error) {
	params := url.Values{}
	params.Set("bid", "reply")
	res := &struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
			URL   string `json:"url"`
		} `json:"data"`
		Msg string `json:"message"`
		TTL int    `json:"ttl"`
	}{}
	ip := metadata.String(c, metadata.RemoteIP)
	if err := s.httpClient.Get(c, conf.Conf.Reply.CaptchaTokenURL, ip, params, res); err != nil {
		log.Error("s.httpClient.Get(%s) error(%v)", conf.Conf.Reply.CaptchaTokenURL+"?"+params.Encode(), err)
		return "", "", err
	}
	if res.Code != 0 {
		log.Error("s.httpClient.Get(%s?%s) code:%d", conf.Conf.Reply.CaptchaTokenURL, params.Encode(), res.Code)
		return "", "", ecode.Int(res.Code)
	}
	return res.Data.Token, res.Data.URL, nil
}

// Verify Verify.
func (s *CaptchaDao) Verify(c context.Context, token, code string) error {
	params := url.Values{}
	params.Set("token", token)
	params.Set("code", code)
	res := &struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		TTL  int    `json:"ttl"`
	}{}
	ip := metadata.String(c, metadata.RemoteIP)
	if err := s.httpClient.Post(c, conf.Conf.Reply.CaptchaVerifyURL, ip, params, res); err != nil {
		log.Error("s.httpClient.POST(%s) error(%v)", conf.Conf.Reply.CaptchaVerifyURL+"?"+params.Encode(), err)
		return err
	}
	if res.Code != 0 {
		log.Error("s.httpClient.POST(%s?%s) code:%d", conf.Conf.Reply.CaptchaVerifyURL, params.Encode(), res.Code)
		return ecode.Int(res.Code)
	}
	return nil
}
