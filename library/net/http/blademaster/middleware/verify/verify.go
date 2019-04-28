package verify

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"
	"sync"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
)

const (
	_secretURI = "/api/getsecret"
)

// Verify is is the verify model.
type Verify struct {
	lock sync.RWMutex
	keys map[string]string

	secretURI string
	client    *bm.Client
}

// Config is the verify config model.
type Config struct {
	OpenServiceHost string
	HTTPClient      *bm.ClientConfig
}

var _defaultConfig = &Config{
	OpenServiceHost: "http://open.bilibili.co",
	HTTPClient: &bm.ClientConfig{
		App: &bm.App{
			Key:    "53e2fa226f5ad348",
			Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
		},
		Dial:      xtime.Duration(time.Millisecond * 100),
		Timeout:   xtime.Duration(time.Millisecond * 300),
		KeepAlive: xtime.Duration(time.Second * 60),
	},
}

// New will create a verify middleware by given config.
// panic on conf is nil.
func New(conf *Config) *Verify {
	if conf == nil {
		conf = _defaultConfig
	}
	v := &Verify{
		keys:   make(map[string]string),
		client: bm.NewClient(conf.HTTPClient),

		secretURI: conf.OpenServiceHost + _secretURI,
	}
	return v
}

func (v *Verify) verify(ctx *bm.Context) error {
	req := ctx.Request
	params := req.Form
	if req.Method == "POST" {
		// Give priority to sign in url query, otherwise check sign in post form.
		q := req.URL.Query()
		if q.Get("sign") != "" {
			params = q
		}
	}

	// check timestamp is not empty (TODO : Check if out of some seconds.., like 100s)
	if params.Get("ts") == "" {
		log.Error("ts is empty")
		return ecode.RequestErr
	}

	sign := params.Get("sign")
	params.Del("sign")
	defer params.Set("sign", sign)
	sappkey := params.Get("appkey")
	v.lock.RLock()
	secret, ok := v.keys[sappkey]
	v.lock.RUnlock()
	if !ok {
		fetched, err := v.fetchSecret(ctx, sappkey)
		if err != nil {
			return err
		}
		v.lock.Lock()
		v.keys[sappkey] = fetched
		v.lock.Unlock()
		secret = fetched
	}

	if hsign := Sign(params, sappkey, secret, true); hsign != sign {
		if hsign1 := Sign(params, sappkey, secret, false); hsign1 != sign {
			log.Error("Get sign: %s, expect %s", sign, hsign)
			return ecode.SignCheckErr
		}
	}
	return nil
}

// Verify will inject into handler func as verify required
func (v *Verify) Verify(ctx *bm.Context) {
	if err := v.verify(ctx); err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
	}
}

// VerifyUser is used to mark path as verify and mid required.
func (v *Verify) VerifyUser(ctx *bm.Context) {
	if err := v.verify(ctx); err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
	}

	var midReq struct {
		Mid int64 `form:"mid" validate:"required"`
	}
	if err := ctx.Bind(&midReq); err != nil {
		return
	}
	ctx.Set("mid", midReq.Mid)
	if md, ok := metadata.FromContext(ctx); ok {
		md[metadata.Mid] = midReq.Mid
	}
}

// Sign is used to sign form params by given condition.
func Sign(params url.Values, appkey string, secret string, lower bool) string {
	data := params.Encode()
	if strings.IndexByte(data, '+') > -1 {
		data = strings.Replace(data, "+", "%20", -1)
	}
	if lower {
		data = strings.ToLower(data)
	}
	digest := md5.Sum([]byte(data + secret))
	return hex.EncodeToString(digest[:])
}

func (v *Verify) fetchSecret(ctx *bm.Context, appkey string) (string, error) {
	params := url.Values{}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			AppSecret string `json:"app_secret"`
		} `json:"data"`
	}
	params.Set("sappkey", appkey)
	if err := v.client.Get(ctx, v.secretURI, metadata.String(ctx, metadata.RemoteIP), params, &resp); err != nil {
		return "", err
	}

	if resp.Code != 0 || resp.Data.AppSecret == "" {
		log.Error("Failed to fetch secret with request(%s, %s) code(%d)", v.secretURI, params.Encode(), resp.Code)
		if resp.Code != 0 {
			return "", ecode.Int(resp.Code)
		}
		return "", ecode.ServerErr
	}
	return resp.Data.AppSecret, nil
}
