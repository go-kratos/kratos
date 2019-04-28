package http

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"
	"sync"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Verify is is the verify model.
type Verify struct {
	keys map[string]string
	lock sync.RWMutex
	// TODO to table data.
	authAppkeyMap map[string]string
}

// Config is the verify config model.
type Config struct {
	AuthAppkeyMap map[string]string
}

// NewThirdVerify will create a verify middleware by given config.
func NewThirdVerify(conf *Config) *Verify {
	if conf == nil || len(conf.AuthAppkeyMap) == 0 {
		panic("conf can not be nil")
	}
	v := &Verify{
		keys:          make(map[string]string),
		authAppkeyMap: conf.AuthAppkeyMap,
	}
	return v
}

// Verify will inject into handler func as verify required
func (v *Verify) Verify(ctx *bm.Context) {
	if err := v.verify(ctx); err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
	}
}

func (v *Verify) verify(c *bm.Context) error {
	req := c.Request
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
		fetched, err := v.fetchSecret(c, sappkey)
		if err != nil {
			return err
		}
		v.lock.Lock()
		v.keys[sappkey] = fetched
		v.lock.Unlock()
		secret = fetched
	}
	url := req.URL.Path
	if hsign := Sign(url, params, sappkey, secret, true); hsign != sign {
		if hsign1 := Sign(url, params, sappkey, secret, false); hsign1 != sign {
			log.Error("Get sign: %s, expect %s", sign, hsign)
			return ecode.SignCheckErr
		}
	}
	return nil
}

// Sign is used to sign form params by given condition.
func Sign(url string, params url.Values, appkey string, secret string, lower bool) string {
	data := params.Encode()
	if strings.IndexByte(data, '+') > -1 {
		data = strings.Replace(data, "+", "%20", -1)
	}
	if lower {
		data = strings.ToLower(data)
	}
	digest := md5.Sum([]byte(url + "?" + data + secret))
	return hex.EncodeToString(digest[:])
}

// TODO auth apppkey to table.
func (v *Verify) fetchSecret(ctx *bm.Context, appkey string) (string, error) {
	return v.authAppkeyMap[appkey], nil
}
