package http

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"

	"go-common/app/service/main/passport-game/model"
	"go-common/app/service/main/passport-game/service"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// VerifySign verify sign.
func verifySign(c *bm.Context, s *service.Service) (res *model.App, err error) {
	var (
		r     = c.Request
		query = r.Form
	)
	if r.Method == "POST" {
		// Give priority to sign in url query, otherwise check sign in post form.
		p := c.Request.URL.Query()
		if p.Get("sign") != "" {
			query = p
		}
	}
	if query.Get("ts") == "" {
		err = ecode.RequestErr
		return
	}
	appKey := query.Get("appkey")
	if appKey == "" {
		err = ecode.RequestErr
		return
	}
	app, ok := s.APP(appKey)
	if !ok {
		err = ecode.AppKeyInvalid
		return
	}
	secret := app.AppSecret
	tmp := encodeQuery(query)
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(strings.ToLower(tmp) + secret))
	sign := query.Get("sign")
	if hex.EncodeToString(mh[:]) != sign {
		mh1 := md5.Sum([]byte(tmp + secret))
		if hex.EncodeToString(mh1[:]) != sign {
			err = ecode.SignCheckErr
		}
	}
	res = app
	return
}

// encodeQuery encodes the values into ``URL encoded'' form ("bar=baz&foo=quux") sorted by key.
// NOTE: sign ignored!!!
func encodeQuery(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if k == "sign" {
			continue
		}
		vs := v[k]
		prefix := url.QueryEscape(k) + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}
