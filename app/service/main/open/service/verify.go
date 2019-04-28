package service

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Verify .
func (s *Service) Verify(c *bm.Context) {
	if err := s.verify(c); err != nil {
		c.JSON(nil, err)
		c.Abort()
		return
	}
}

// Verify .
func (s *Service) verify(c *bm.Context) error {
	req := c.Request
	params := req.Form
	if req.Method == "POST" {
		// Give priority to sign in url query, otherwise check sign in post form.
		q := req.URL.Query()
		if q.Get("sign") != "" {
			params = q
		}
	}

	sign := params.Get("sign")
	params.Del("sign")
	defer params.Set("sign", sign)
	sappkey := params.Get("appkey")
	secret, ok := s.appsecrets[sappkey]
	if !ok {
		log.Error("appkey(%s) not found in cache", sappkey)
		return ecode.NothingFound
	}
	if hsign := s.sign(params, sappkey, secret, true); hsign != sign {
		if hsign1 := s.sign(params, sappkey, secret, false); hsign1 != sign {
			log.Error("Get sign: %s, expect %s", sign, hsign)
			return ecode.SignCheckErr
		}
	}
	return nil
}

// Sign is used to sign form params by given condition.
func (s *Service) sign(params url.Values, appkey, secret string, lower bool) string {
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
