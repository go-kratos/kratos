package realname

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/account/conf"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// RealnameCaptchaGTRegister register a geetest apply
func (d *Dao) RealnameCaptchaGTRegister(c context.Context, mid int64, ip, clientType string, newCaptcha int) (challenge string, err error) {
	var (
		params = url.Values{}
		req    *http.Request
		url    = conf.Conf.Realname.Geetest.RegisterURL
	)
	params.Set("user_id", strconv.FormatInt(mid, 10))
	params.Set("new_captcha", strconv.Itoa(newCaptcha))
	params.Set("ip_address", ip)
	params.Set("client_type", clientType)
	params.Set("gt", conf.Conf.Realname.Geetest.CaptchaID)

	if req, err = http.NewRequest("GET", url+"?"+params.Encode(), nil); err != nil {
		err = errors.Wrapf(err, "dao.RealnameCaptchaGTRegister url(%s) params(%s)", url, params.Encode())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var challengeBytes []byte
	if challengeBytes, err = d.client.Raw(c, req); err != nil {
		return
	}
	challenge = string(challengeBytes)
	if len(challenge) != 32 {
		log.Error("dao.RealnameCaptchaGTRegister challenge : %s ,length not equate 32bit", challenge)
	}
	return
}

// RealnameCaptchaGTRegisterValidate recheck the challenge code and get to seccode
func (d *Dao) RealnameCaptchaGTRegisterValidate(c context.Context, challenge, seccode, clientType, ip, captchaID string, mid int64) (realSeccode string, err error) {
	var (
		params = url.Values{}
		req    *http.Request
		url    = conf.Conf.Realname.Geetest.ValidateURL
	)
	params.Set("seccode", seccode)
	params.Set("challenge", challenge)
	params.Set("captchaid", captchaID)
	params.Set("client_type", clientType)
	params.Set("ip_address", metadata.String(c, metadata.RemoteIP))
	params.Set("json_format", "1")
	params.Set("sdk", "golang_3.0.0")
	params.Set("user_id", strconv.FormatInt(mid, 10))
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	log.Info("gt validate url : %s , params : %s", url, params.Encode())

	if req, err = http.NewRequest("POST", url, strings.NewReader(params.Encode())); err != nil {
		err = errors.Wrapf(err, "dao.RealnameCaptchaGTRegister url(%s) params(%s)", url, params.Encode())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var resp struct {
		Seccode string `json:"seccode"`
	}
	if err = d.client.Do(c, req, &resp); err != nil {
		return
	}
	realSeccode = resp.Seccode
	return
}
