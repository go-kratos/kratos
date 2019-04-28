package http

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_sobotAppKey    = "bcef69bb71499209"
	_sobotAppSecret = "ace486f144f1467eefdce1fe5dfc7b14"
	_sobotAPI       = "https://sso-api.bilibili.co/x/internal/workflow/sobot/user"
)

func sobotSign(handler func(*bm.Context)) func(*bm.Context) {
	return func(c *bm.Context) {
		req := c.Request
		query := req.Form
		if query.Get("ts") == "" {
			log.Error("ts is empty")
			c.JSON(nil, ecode.RequestErr)
			return
		}
		sign := query.Get("sign")
		query.Del("sign")
		sappkey := query.Get("appkey")
		if sappkey != _sobotAppKey {
			log.Error("appkey not matched")
			c.JSON(nil, ecode.RequestErr)
			return
		}
		query.Set("appsecret", _sobotAppSecret)
		tmp := query.Encode()
		if strings.IndexByte(tmp, '+') > -1 {
			tmp = strings.Replace(tmp, "+", "%20", -1)
		}
		mh := md5.Sum([]byte(_sobotAPI + "?" + strings.ToLower(tmp) + _sobotAppSecret))
		if hex.EncodeToString(mh[:]) != sign {
			mh1 := md5.Sum([]byte(_sobotAPI + "?" + tmp + _sobotAppSecret))
			if hex.EncodeToString(mh1[:]) != sign {
				log.Error("Get sign: %s, expect %x", sign, mh1)
				c.JSON(nil, ecode.SignCheckErr)
				return
			}
		}
		handler(c)
	}
}
